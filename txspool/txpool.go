/*
   This file is part of go-palletone.
   go-palletone is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.
   go-palletone is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.
   You should have received a copy of the GNU General Public License
   along with go-palletone.  If not, see <http://www.gnu.org/licenses/>.
*/

/*
 * @author PalletOne core developers <dev@pallet.one>
 * @date 2018
 */

package txspool

import (
	"bytes"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/palletone/go-palletone/dag/dagconfig"
	"github.com/palletone/go-palletone/validator"

	"github.com/ethereum/go-ethereum/event"
	"github.com/palletone/go-palletone/common"
	"github.com/palletone/go-palletone/common/log"
	"github.com/palletone/go-palletone/dag/errors"
	"github.com/palletone/go-palletone/dag/modules"
	"github.com/palletone/go-palletone/dag/palletcache"
	"github.com/palletone/go-palletone/dag/parameter"
	"github.com/palletone/go-palletone/tokenengine"
	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/karalabe/cookiejar.v2/collections/prque"
)

var (
	// Time interval to check for evictable transactions
	evictionInterval = time.Minute
	// Time interval to report transaction pool stats
	statsReportInterval = 8 * time.Second
	//The minimum amount of time in between scans of the orphan pool to evict expired transactions.
	orphanExpireScanInterval = time.Minute * 5
)
var (
	ErrNotFound = errors.New("txpool: not found")
	// ErrInvalidSender is returned if the transaction contains an invalid signature.
	ErrInvalidSender = errors.New("invalid sender")

	// ErrNonceTooLow is returned if the nonce of a transaction is lower than the
	// one present in the local chain.
	ErrNonceTooLow = errors.New("nonce too low")

	// ErrTxFeeTooLow is returned if a transaction's tx_fee is below the value of TXFEE.
	ErrTxFeeTooLow = errors.New("txfee too low")

	// ErrUnderpriced is returned if a transaction's gas price is below the minimum
	// configured for the transaction pool.
	ErrUnderpriced = errors.New("transaction underpriced")

	// ErrReplaceUnderpriced is returned if a transaction is attempted to be replaced
	// with a different one without the required price bump.
	ErrReplaceUnderpriced = errors.New("replacement transaction underpriced")

	// ErrInsufficientFunds is returned if the total cost of executing a transaction
	// is higher than the balance of the user's account.
	ErrInsufficientFunds = errors.New("insufficient funds for gas * price + value")

	// ErrNegativeValue is a sanity error to ensure noone is able to specify a
	// transaction with a negative value.
	ErrNegativeValue = errors.New("negative value")

	// ErrOversizedData is returned if the input data of a transaction is greater
	// than some meaningful limit a user might use. This is not a consensus error
	// making the transaction invalid, rather a DOS protection.
	ErrOversizedData = errors.New("oversized data")
)

var (
	txValidPrometheus = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "prometheus:txpool:tx:valid",
		Help: "txpool tx valid",
	})
	txInvalidPrometheus = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "prometheus:txpool:tx:invalid",
		Help: "txpool tx invalid",
	})

	txAlreadyPrometheus = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "prometheus:txpool:tx:already",
		Help: "txpool tx already",
	})

	txOrphanKnownPrometheus = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "prometheus:txpool:tx:orphan:known",
		Help: "txpool tx orphan known",
	})
	txOrphanValidPrometheus = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "prometheus:txpool:tx:orphan:valid",
		Help: "txpool tx orphan valid",
	})

	txCoinbasePrometheus = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "prometheus:txpool:tx:coinbase",
		Help: "txpool tx coinbase",
	})
)

func init() {
	prometheus.MustRegister(txValidPrometheus)
	prometheus.MustRegister(txInvalidPrometheus)

	prometheus.MustRegister(txAlreadyPrometheus)

	prometheus.MustRegister(txOrphanKnownPrometheus)
	prometheus.MustRegister(txOrphanValidPrometheus)

	prometheus.MustRegister(txCoinbasePrometheus)
}

type TxPool struct {
	config      TxPoolConfig
	unit        IDag
	txFeed      event.Feed
	scope       event.SubscriptionScope
	txValidator IValidator
	journal     *txJournal // Journal of local transaction to back up to disk

	all                   sync.Map          // All transactions to allow lookups
	priority_sorted       *txPrioritiedList // All transactions sorted by price and priority
	outpoints             sync.Map          // utxo标记池  map[modules.OutPoint]*TxPoolTransaction
	orphans               sync.Map          // 孤儿交易缓存池
	outputs               sync.Map          // 缓存 交易的outputs
	reqOutputs            sync.Map          // 缓存 交易的outputs
	sequenTxs             *SequeueTxPoolTxs
	userContractRequests  map[common.Hash]*TxPoolTransaction //用户合约请求，只参与utxo运算，不会被打包
	basedOnRequestOrphans map[common.Hash]*TxPoolTransaction //依赖于userContractRequests的孤儿交易池

	mu             sync.RWMutex
	wg             sync.WaitGroup // for shutdown sync
	quit           chan struct{}  // used for exit
	nextExpireScan time.Time
	cache          palletcache.ICache
	tokenEngine    tokenengine.ITokenEngine
	//enableGasFee   bool
}

// NewTxPool creates a new transaction pool to gather, sort and filter inbound
// transactions from the network.
func NewTxPool(config TxPoolConfig, cachedb palletcache.ICache, unit IDag, enableGasFee bool) *TxPool {
	tokenEngine := tokenengine.Instance
	pool := NewTxPool4DI(config, cachedb, unit, tokenEngine, nil, enableGasFee)
	val := validator.NewValidate(unit, pool, unit, unit, nil, cachedb, false, enableGasFee)
	pool.txValidator = val
	pool.startJournal(config)
	return pool
}

//构造函数的依赖注入，主要用于UT
func NewTxPool4DI(config TxPoolConfig, cachedb palletcache.ICache, unit IDag,
	tokenEngine tokenengine.ITokenEngine, validator IValidator, enableGasFee bool) *TxPool { // chainconfig *params.ChainConfig,
	// Sanitize the input to ensure no vulnerable gas prices are set
	config = (&config).sanitize()
	// Create the transaction pool with its initial settings
	pool := &TxPool{
		config:                config,
		unit:                  unit,
		all:                   sync.Map{},
		sequenTxs:             new(SequeueTxPoolTxs),
		outpoints:             sync.Map{},
		nextExpireScan:        time.Now().Add(config.OrphanTTL),
		orphans:               sync.Map{},
		outputs:               sync.Map{},
		reqOutputs:            sync.Map{},
		basedOnRequestOrphans: make(map[common.Hash]*TxPoolTransaction),
		userContractRequests:  make(map[common.Hash]*TxPoolTransaction),
		cache:                 cachedb,
		tokenEngine:           tokenEngine,
	}
	pool.mu = sync.RWMutex{}
	pool.priority_sorted = newTxPrioritiedList(&pool.all)
	pool.txValidator = validator
	pool.startJournal(config)
	return pool
}
func (pool *TxPool) startJournal(config TxPoolConfig) {
	// If local transactions and journaling is enabled, load from disk
	if !config.NoLocals && config.Journal != "" {
		log.Info("Journal path:" + config.Journal)
		pool.journal = newTxJournal(config.Journal)

		if err := pool.journal.load(pool.addJournalTx); err != nil {
			log.Warn("Failed to load transaction journal", "err", err)
		}
		if err := pool.journal.rotate(pool.local()); err != nil {
			log.Warn("Failed to rotate transaction journal", "err", err)
		}
	}
	// Start the event loop and return
	pool.wg.Add(1)
	go pool.loop()
}

// return a utxo by the outpoint in txpool
func (pool *TxPool) GetUtxoFromAll(outpoint *modules.OutPoint) (*modules.Utxo, error) {
	return pool.GetUtxoEntry(outpoint)
}

func (pool *TxPool) Clear() {
	pool.all = sync.Map{}
	pool.sequenTxs = new(SequeueTxPoolTxs)
	pool.outpoints = sync.Map{}
	pool.orphans = sync.Map{}
	pool.outputs = sync.Map{}
	pool.reqOutputs = sync.Map{}
}
func (pool *TxPool) GetUtxoEntry(outpoint *modules.OutPoint) (*modules.Utxo, error) {
	if inter, ok := pool.outputs.Load(*outpoint); ok {
		utxo := inter.(*modules.Utxo)
		return utxo, nil
	}
	if inter, ok := pool.reqOutputs.Load(*outpoint); ok {
		utxo := inter.(*modules.Utxo)
		return utxo, nil
	}
	return pool.unit.GetUtxoEntry(outpoint)
}

// return a stxo by the outpoint in txpool
func (pool *TxPool) GetStxoEntry(outpoint *modules.OutPoint) (*modules.Stxo, error) {
	return pool.unit.GetStxoEntry(outpoint)
}
func (pool *TxPool) GetTxOutput(outpoint *modules.OutPoint) (*modules.Utxo, error) {
	if inter, ok := pool.outputs.Load(*outpoint); ok {
		utxo := inter.(*modules.Utxo)
		return utxo, nil
	}
	return pool.unit.GetTxOutput(outpoint)
}

// loop is the transaction pool's main event loop, waiting for and reacting to
// outside blockchain events as well as for various reporting and transaction
// eviction events.
func (pool *TxPool) loop() {
	defer pool.wg.Done()
	// Start the stats reporting and transaction eviction tickers
	var prevPending, prevQueued int

	report := time.NewTicker(statsReportInterval)
	defer report.Stop()

	evict := time.NewTicker(evictionInterval)
	defer evict.Stop()

	journal := time.NewTicker(pool.config.Rejournal)
	defer journal.Stop()
	// delete txspool's confirmed tx
	deleteTxTimer := time.NewTicker(10 * time.Minute)
	defer deleteTxTimer.Stop()

	orphanExpireScan := time.NewTicker(orphanExpireScanInterval)
	defer orphanExpireScan.Stop()

	// Keep waiting for and reacting to the various events
	for {
		select {
		// Handle stats reporting ticks
		case <-report.C:
			pending, queued, _ := pool.stats()

			if pending != prevPending || queued != prevQueued {
				log.Debug("Transaction pool status report", "executable", pending, "queued", queued)
				prevPending, prevQueued = pending, queued
			}

			// Handle inactive account transaction eviction
		case <-evict.C:

			// Handle local transaction journal rotation ----- once a honr -----
		case <-journal.C:
			if pool.journal != nil {
				pool.mu.Lock()
				if err := pool.journal.rotate(pool.local()); err != nil {
					log.Warn("Failed to rotate local tx journal", "err", err)
				}
				pool.mu.Unlock()
			}
			// delete tx
		case <-deleteTxTimer.C:
			pool.DeleteTx()

			// quit
		case <-orphanExpireScan.C:
			pool.mu.Lock()
			pool.limitNumberOrphans()
			pool.mu.Unlock()
		case <-pool.quit:
			log.Info("txspool are quit now", "time", time.Now().String())
			return
		}
	}
}

// Stats retrieves the current pool stats, namely the number of pending and the
// number of queued (non-executable) transactions.
func (pool *TxPool) Status() (int, int, int) {
	return pool.stats()
}

// stats retrieves the current pool stats, namely the number of pending and the
// number of queued (non-executable) transactions.
func (pool *TxPool) stats() (int, int, int) {
	p_count, q_count := 0, 0
	poolTxs := pool.AllTxpoolTxs()
	orphanTxs := pool.AllOrphanTxs()
	seq_txs := pool.sequenTxs.All()
	for _, tx := range poolTxs {
		if tx.Pending {
			p_count++
		} else if !tx.IsOrphan {
			q_count++
		}
	}
	for _, tx := range seq_txs {
		if tx.Pending {
			p_count++
		} else {
			q_count++
		}
	}
	return p_count, q_count, len(orphanTxs)
}

// Content retrieves the data content of the transaction pool, returning all the
// pending as well as queued transactions, grouped by account and sorted by nonce.
func (pool *TxPool) Content() (map[common.Hash]*TxPoolTransaction, map[common.Hash]*TxPoolTransaction) {
	pending := make(map[common.Hash]*TxPoolTransaction)
	queue := make(map[common.Hash]*TxPoolTransaction)

	alltxs := pool.AllTxpoolTxs()
	orphanTxs := pool.AllOrphanTxs()
	for hash, tx := range alltxs {
		if tx.Pending {
			pending[hash] = tx
		}
		if !tx.Pending {
			queue[hash] = tx
		}
	}
	for hash, tx := range orphanTxs {
		if !tx.Pending {
			queue[hash] = tx
		}
	}
	return pending, queue
}

// Pending retrieves all currently processable transactions, groupped by origin
// account and sorted by priority level. The returned transaction set is a copy and can be
// freely modified by calling code.
func (pool *TxPool) Pending() (map[common.Hash][]*TxPoolTransaction, error) {
	return pool.pending()
}
func (pool *TxPool) pending() (map[common.Hash][]*TxPoolTransaction, error) {
	pending := make(map[common.Hash][]*TxPoolTransaction)
	txs := pool.AllTxpoolTxs()
	for _, tx := range txs {
		if tx.Pending {
			pending[tx.UnitHash] = append(pending[tx.UnitHash], tx)
		}
	}
	return pending, nil
}

// Queued txs
func (pool *TxPool) Queued() ([]*TxPoolTransaction, error) {
	queue := make([]*TxPoolTransaction, 0)
	txs := pool.AllTxpoolTxs()
	for _, tx := range txs {
		if !tx.Pending {
			queue = append(queue, tx)
		}
	}
	return queue, nil
}

func (pool *TxPool) AllLength() int {
	var count int
	pool.all.Range(func(k, v interface{}) bool {
		count++
		return true
	})
	return count
}
func (pool *TxPool) AllTxpoolTxs() map[common.Hash]*TxPoolTransaction {
	txs := make(map[common.Hash]*TxPoolTransaction)
	pool.all.Range(func(k, v interface{}) bool {
		hash := k.(common.Hash)
		tx := v.(*TxPoolTransaction)
		tx_hash := tx.Tx.Hash()
		if hash != tx_hash {
			pool.all.Delete(hash)
			pool.all.Store(tx_hash, tx)
		}
		txs[tx_hash] = tx
		return true
	})
	return txs
}
func (pool *TxPool) AllOrphanTxs() map[common.Hash]*TxPoolTransaction {
	txs := make(map[common.Hash]*TxPoolTransaction)
	pool.orphans.Range(func(k, v interface{}) bool {
		tx := v.(*TxPoolTransaction)
		txs[tx.Tx.Hash()] = tx
		return true
	})
	return txs
}

// local retrieves all currently known local transactions, groupped by origin
// account and sorted by price. The returned transaction set is a copy and can be
// freely modified by calling code.
func (pool *TxPool) local() map[common.Hash]*TxPoolTransaction {
	txs := make(map[common.Hash]*TxPoolTransaction)
	pending, _ := pool.pending()
	for _, list := range pending {
		for _, tx := range list {
			if tx != nil {
				txs[tx.Tx.Hash()] = tx
			}
		}
	}
	return txs
}

// validateTx checks whether a transaction is valid according to the consensus
// rules and adheres to some heuristic limits of the local node (price and size).
func (pool *TxPool) validateTx(tx *modules.Transaction, local bool) ([]*modules.Addition,
	validator.ValidationCode, error) {

	return pool.txValidator.ValidateTx(tx, !tx.IsOnlyContractRequest())
}

// This function MUST be called with the txpool lock held (for reads).
func (pool *TxPool) isTransactionInPool(hash common.Hash) bool {
	if _, exist := pool.all.Load(hash); exist {
		return true
	}
	if _, exist := pool.orphans.Load(hash); exist {
		return true
	}
	return false
}

// IsTransactionInPool returns whether or not the passed transaction already exists in the main pool.
func (pool *TxPool) IsTransactionInPool(hash common.Hash) bool {
	return pool.isTransactionInPool(hash)
}
func (pool *TxPool) setPriorityLvl(tx *TxPoolTransaction) {
	tx.Priority_lvl = tx.GetPriorityLvl()
}

func (pool *TxPool) add_journalTx(ptx *TxPoolTransaction) error {
	hash := ptx.Tx.Hash()
	reqHash := ptx.Tx.RequestHash()
	exitsInDb, _ := pool.unit.IsTransactionExist(hash)
	if exitsInDb {
		log.Infof("Tx[%s] already exist in db", hash.String())
		return nil
	}
	if id, err := pool.unit.GetTxHashByReqId(reqHash); err == nil {
		log.Infof("Request[%s] already exist in db,txhash[%s]", reqHash.String(), id.String())
		return nil
	}
	// Don't accept the transaction if it already in the pool .
	if _, has := pool.all.Load(hash); has {
		txAlreadyPrometheus.Add(1)
		log.Trace("Discarding already known transaction", "hash", hash.String())
		return fmt.Errorf("known transaction: %s", hash.String())
	}
	if pool.isOrphanInPool(hash) {
		txOrphanKnownPrometheus.Add(1)
		return fmt.Errorf("know orphanTx: %s", hash.String())
	}

	if ptx.Tx.IsSystemContract() && !ptx.Tx.IsOnlyContractRequest() {
		log.Infof("[%s]tx[%s] is a full system contract invoke tx, don't support", reqHash.ShortStr(), hash.String())
		return errors.New("txpool: not support")
	}
	var deletedReq *TxPoolTransaction
	if ptx.Tx.IsUserContract() && !ptx.Tx.IsOnlyContractRequest() { //FullTx about user contract
		//delete request
		var ok bool
		deletedReq, ok = pool.userContractRequests[reqHash]
		if ok {
			delete(pool.userContractRequests, reqHash)
			log.Debugf("[%s]delete user contract request by hash:%s", reqHash.ShortStr(), hash.String())
		}
	}
	reverseDeleteReq := func() {
		if deletedReq != nil {
			pool.userContractRequests[deletedReq.TxHash] = deletedReq
			log.Debugf("[%s]reverse delete request %s", reqHash.ShortStr(), deletedReq.TxHash.String())
		}
	}

	// If the transaction pool is full, discard underpriced transactions
	length := pool.AllLength()
	if uint64(length) >= pool.config.GlobalSlots+pool.config.GlobalQueue {
		// If the new transaction is underpriced, don't accept it
		if pool.priority_sorted.Underpriced(ptx) {
			log.Trace("Discarding underpriced transaction", "hash", hash, "price", ptx.GetTxFee().Int64())
			return ErrUnderpriced
		}
		// New transaction is better than our worse ones, make room for it
		count := length - int(pool.config.GlobalSlots+pool.config.GlobalQueue-1)
		if count > 0 {
			drop := pool.priority_sorted.Discard(count)
			for _, tx := range drop {
				log.Trace("Discarding freshly underpriced transaction", "hash", hash, "price", tx.GetTxFee().Int64())
				pool.removeTransaction(tx, true)
				pool.removeTx(tx.Tx.Hash())
			}
		}
	}
	if ptx.Tx.IsUserContract() && ptx.Tx.IsOnlyContractRequest() {
		//user contract request
		log.Debugf("[%s]add tx[%s] to user contract request pool", reqHash.ShortStr(), hash.String())
		pool.userContractRequests[ptx.TxHash] = ptx
		pool.txFeed.Send(modules.TxPreEvent{Tx: ptx.Tx, IsOrphan: false})
	} else { //不是用户合约请求
		//有可能是连续的用户合约请求R1,R2，但是R2先被执行完，这个时候R1还在RequestPool里面，没办法被打包，所以R2应该被扔到basedOnReqOrphanPool
		//父交易还是Request，所以本Tx是Orphan
		if pool.isBasedOnRequestPool(ptx) {
			log.Debugf("Tx[%s]'s parent or ancestor is a request, not a full tx, add it to based on request pool",
				ptx.TxHash.String())
			if err := pool.addBasedOnReqOrphanTx(ptx); err != nil {
				log.Errorf("add tx[%s] to based on request pool error:%s", ptx.TxHash.String(), err.Error())
				reverseDeleteReq()
				return err
			}
		} else {
			//3. process normal tx
			go pool.priority_sorted.Put(ptx)

			pool.all.Store(hash, ptx)
			pool.addCache(ptx)
			txValidPrometheus.Add(1)
			// We've directly injected a replacement transaction, notify subsystems
			pool.txFeed.Send(modules.TxPreEvent{Tx: ptx.Tx, IsOrphan: false})
		}
	}

	// 更新一次孤儿交易池数据。
	pool.reflashOrphanTxs(ptx.Tx, pool.AllOrphanTxs(), false)
	return nil
}

// add validates a transaction and inserts it into the non-executable queue for
// later pending promotion and execution. If the transaction is a replacement for
// an already pending or queued one, it overwrites the previous and returns this
// so outer code doesn't uselessly call promote.
//
// If a newly added transaction is marked as local, its sending account will be
// whitelisted, preventing any associated transaction from being dropped out of
// the pool due to pricing constraints.
func (pool *TxPool) add(tx *modules.Transaction, local bool) (bool, error) {
	hash := tx.Hash()
	reqHash := tx.RequestHash()
	msgs := tx.Messages()
	if msgs[0].Payload.(*modules.PaymentPayload).IsCoinbase() {
		txCoinbasePrometheus.Add(1)
		return true, nil
	}
	exitsInDb, _ := pool.unit.IsTransactionExist(hash)
	if exitsInDb {
		log.Infof("Tx[%s] already exist in db", hash.String())
		return false, nil
	}
	if id, err := pool.unit.GetTxHashByReqId(reqHash); err == nil {
		log.Infof("Request[%s] already exist in db,txhash[%s]", reqHash.String(), id.String())
		return false, nil
	}
	// Don't accept the transaction if it already in the pool .
	if _, has := pool.all.Load(hash); has {
		txAlreadyPrometheus.Add(1)
		log.Trace("Discarding already known transaction", "hash", hash.String())
		return false, fmt.Errorf("known transaction: %s", hash.String())
	}
	if pool.isOrphanInPool(hash) {
		txOrphanKnownPrometheus.Add(1)
		return false, fmt.Errorf("know orphanTx: %s", hash.String())
	}

	if tx.IsSystemContract() && !tx.IsOnlyContractRequest() {
		log.Infof("[%s]tx[%s] is a full system contract invoke tx, don't support", reqHash.ShortStr(), hash.String())
		return false, errors.New("txpool: not support")
	}
	var deletedReq *TxPoolTransaction
	if tx.IsUserContract() && !tx.IsOnlyContractRequest() { //FullTx about user contract
		//delete request
		var ok bool
		deletedReq, ok = pool.userContractRequests[reqHash]
		if ok {
			delete(pool.userContractRequests, reqHash)
			log.Debugf("[%s]delete user contract request by hash:%s", reqHash.ShortStr(), hash.String())
		}
	}
	reverseDeleteReq := func() {
		if deletedReq != nil {
			pool.userContractRequests[deletedReq.TxHash] = deletedReq
			log.Debugf("[%s]reverse delete request %s", reqHash.ShortStr(), deletedReq.TxHash.String())
		}
	}

	// If the transaction fails basic validation, discard it
	addition, code, err := pool.validateTx(tx, !tx.IsOnlyContractRequest())
	if err != nil && code != validator.TxValidationCode_ORPHAN {
		reverseDeleteReq()
		return false, validator.NewValidateError(code)
	}
	if code == validator.TxValidationCode_ORPHAN {
		if ok, err := pool.ValidateOrphanTx(tx); ok {
			txOrphanValidPrometheus.Add(1)
			log.Debug("validated the orphanTx", "hash", hash.String())
			pool.addOrphan(tx)
			return true, nil
		} else if err != nil {
			return false, err
		} // 孤儿单元验证通过了，将其添加到待打包交易池。
		code = validator.TxValidationCode_VALID
	}
	if code != validator.TxValidationCode_VALID {
		reverseDeleteReq()
		return false, validator.NewValidateError(code)
	}
	ptx := pool.convertTx(tx, addition)
	// 计算优先级
	pool.setPriorityLvl(ptx)
	// If the transaction pool is full, discard underpriced transactions
	length := pool.AllLength()
	if uint64(length) >= pool.config.GlobalSlots+pool.config.GlobalQueue {
		// If the new transaction is underpriced, don't accept it
		if pool.priority_sorted.Underpriced(ptx) {
			log.Trace("Discarding underpriced transaction", "hash", hash, "price", ptx.GetTxFee().Int64())
			return false, ErrUnderpriced
		}
		// New transaction is better than our worse ones, make room for it
		count := length - int(pool.config.GlobalSlots+pool.config.GlobalQueue-1)
		if count > 0 {
			drop := pool.priority_sorted.Discard(count)
			for _, tx := range drop {
				log.Trace("Discarding freshly underpriced transaction", "hash", hash, "price", tx.GetTxFee().Int64())
				pool.removeTransaction(tx, true)
				pool.removeTx(tx.Tx.Hash())
			}
		}
	}

	if tx.IsUserContract() && tx.IsOnlyContractRequest() {
		//user contract request
		log.Debugf("[%s]add tx[%s] to user contract request pool", reqHash.ShortStr(), hash.String())
		pool.userContractRequests[ptx.TxHash] = ptx
		pool.addCache(ptx)
		pool.txFeed.Send(modules.TxPreEvent{Tx: tx, IsOrphan: false})
	} else { //不是用户合约请求
		//有可能是连续的用户合约请求R1,R2，但是R2先被执行完，这个时候R1还在RequestPool里面，没办法被打包，所以R2应该被扔到basedOnReqOrphanPool
		//父交易还是Request，所以本Tx是Orphan
		pool.addCache(ptx)
		if pool.isBasedOnRequestPool(ptx) {
			log.Debugf("Tx[%s]'s parent or ancestor is a request, not a full tx, add it to based on request pool",
				ptx.TxHash.String())
			if err := pool.addBasedOnReqOrphanTx(ptx); err != nil {
				log.Errorf("add tx[%s] to based on request pool error:%s", ptx.TxHash.String(), err.Error())
				reverseDeleteReq()
				return true, err
			}
		} else {
			//3. process normal tx
			pool.priority_sorted.Put(ptx)
			if local {
				go pool.journalTx(ptx)
			}
			pool.all.Store(hash, ptx)
			txValidPrometheus.Add(1)

			err = pool.checkBasedOnReqOrphanTxToNormal(hash, reqHash)
			if err != nil {
				return true, err
			}
			// We've directly injected a replacement transaction, notify subsystems
			pool.txFeed.Send(modules.TxPreEvent{Tx: tx, IsOrphan: false})
		}
	}

	// 更新一次孤儿交易池数据。
	pool.reflashOrphanTxs(tx, pool.AllOrphanTxs(), local)
	return true, nil
}

// journalTx adds the specified transaction to the local disk journal if it is
// deemed to have been sent from a local account.
func (pool *TxPool) journalTx(tx *TxPoolTransaction) {
	// Only journal if it's enabled and the transaction is local
	if pool.config.NoLocals {
		return
	}
	if len(tx.From) > 0 && pool.journal == nil {
		log.Trace("Pool journal is nil.", "journal", pool.journal.path)
		return
	}
	pool.mu.Lock()
	defer pool.mu.Unlock()
	if err := pool.journal.insert(tx); err != nil {
		log.Warn("Failed to journal local transaction", "err", err)
	}
}

// promoteTx adds a transaction to the pending (processable) list of transactions.
//
// Note, this method assumes the pool lock is held!
func (pool *TxPool) promoteTx(hash, txhash common.Hash, tx *TxPoolTransaction, number, index uint64) {
	// Try to insert the transaction into the pending queue
	tx.Pending = true
	tx.Discarded = false
	tx.UnitHash = hash
	tx.UnitIndex = number
	tx.Index = index
	// delete utxo
	//pool.deletePoolUtxos(tx.Tx)
	pool.all.Store(txhash, tx)
}

// AddLocal enqueues a single transaction into the pool if it is valid, marking
// the sender as a local one in the mean time, ensuring it goes around the local
// pricing constraints.
func (pool *TxPool) AddLocal(tx *modules.Transaction) error {
	//删除请求交易，添加完整交易
	if tx.RequestHash() != tx.Hash() && pool.IsTransactionInPool(tx.RequestHash()) {
		pool.DeleteTxByHash(tx.RequestHash())
	}
	return pool.addLocal(tx)
}
func (pool *TxPool) addLocal(tx *modules.Transaction) error {
	pool.mu.Lock()
	defer pool.mu.Unlock()
	return pool.addTx(tx, !pool.config.NoLocals)
}
func (pool *TxPool) addJournalTx(ptx *TxPoolTransaction) error {
	pool.mu.Lock()
	defer pool.mu.Unlock()
	// Try to inject the transaction and update any state
	return pool.add_journalTx(ptx)
}

// AddRemote enqueues a single transaction into the pool if it is valid. If the
// sender is not among the locally tracked ones, full pricing constraints will
// apply.
func (pool *TxPool) AddRemote(tx *modules.Transaction) error {
	if tx.TxMessages()[0].Payload.(*modules.PaymentPayload).IsCoinbase() {
		return nil
	}
	return pool.addTx(tx, !pool.config.NoLocals)
}

func IsCoinBase(tx *modules.Transaction) bool {
	msgs := tx.TxMessages()
	if len(msgs) != 1 {
		return false
	}
	msg, ok := msgs[0].Payload.(*modules.PaymentPayload)
	if !ok {
		return false
	}
	return msg.IsCoinbase()
}

// addTx enqueues a single transaction into the pool if it is valid.
func (pool *TxPool) addTx(tx *modules.Transaction, local bool) error {
	// Try to inject the transaction and update any state
	replace, err := pool.add(tx, local)
	if err != nil {
		return err
	}
	// If we added a new transaction, run promotion checks and return
	if !replace {
		pool.promoteExecutables()
	}
	return nil
}

func (pool *TxPool) GetUnpackedTxsByAddr(addr common.Address) ([]*TxPoolTransaction, error) {
	pool.mu.RLock()
	defer pool.mu.RUnlock()
	return pool.getPoolTxsByAddr(addr, true)
}

func (pool *TxPool) getPoolTxsByAddr(addr common.Address, onlyUnpacked bool) ([]*TxPoolTransaction, error) {
	// 将交易按地址分类
	result := make([]*TxPoolTransaction, 0)
	poolTxs := pool.AllTxpoolTxs()
	for _, tx := range pool.getSequenTxs() {
		poolTxs[tx.Tx.Hash()] = tx
	}
	for hash, tx := range pool.userContractRequests {
		poolTxs[hash] = tx
	}
	for hash, tx := range pool.basedOnRequestOrphans {
		poolTxs[hash] = tx
	}
	for _, tx := range poolTxs {
		// 如果已被打包、被标记为discard，则忽略
		if tx.Pending || tx.Discarded {
			continue
		}
		if tx.IsFrom(addr) || tx.IsTo(addr) {
			result = append(result, tx)
		}
	}
	return result, nil
}

// Get returns a transaction if it is contained in the pool
// and nil otherwise.
func (pool *TxPool) GetTx(hash common.Hash) (*TxPoolTransaction, error) {
	interTx, has := pool.all.Load(hash)
	if has {
		tx := interTx.(*TxPoolTransaction)
		log.Debug("get tx info by hash in txpool... ", "unit_hash", tx.UnitHash, "p_tx", tx)
		return tx, nil
	}
	if otx, exist := pool.orphans.Load(hash); exist {
		tx := otx.(*TxPoolTransaction)
		log.Debug("get tx info by hash in orphan txpool... ", "txhash", tx.Tx.Hash(), "info", tx)
		return tx, nil
	}
	if tx, ok := pool.userContractRequests[hash]; ok {
		return tx, nil
	}
	if tx, ok := pool.basedOnRequestOrphans[hash]; ok {
		return tx, nil
	}

	//4个池都找不到
	return nil, ErrNotFound
}

// DeleteTx
func (pool *TxPool) DeleteTx() error {
	txs := pool.AllTxpoolTxs()
	for hash, tx := range txs {
		if tx.Discarded {
			// delete Discarded tx
			log.Debug("delete the status of Discarded tx.", "tx_hash", hash.String())
			pool.DeleteTxByHash(hash)
			continue
		}
		if tx.CreationDate.Add(pool.config.Removetime).Before(time.Now()) {
			// delete
			log.Debug("delete the confirmed tx.", "tx_hash", hash)
			pool.DeleteTxByHash(hash)
			continue
		}
		if tx.CreationDate.Add(pool.config.Lifetime).Before(time.Now()) {
			// delete
			log.Debug("delete the tx(overtime).", "tx_hash", hash)
			pool.DeleteTxByHash(hash)
			continue
		}
	}
	return nil
}

func (pool *TxPool) DeleteTxByHash(hash common.Hash) error {
	inter, has := pool.all.Load(hash)
	if !has {
		if inter, has = pool.orphans.Load(hash); !has {
			return errors.New(fmt.Sprintf("the tx(%s) isn't exist in pool.", hash.String()))
		}
	}
	tx := inter.(*TxPoolTransaction)
	pool.all.Delete(hash)
	pool.orphans.Delete(hash)
	pool.priority_sorted.Removed()

	if tx != nil {
		for i, msg := range tx.Tx.TxMessages() {
			if msg.App == modules.APP_PAYMENT {
				payment, ok := msg.Payload.(*modules.PaymentPayload)
				if ok {
					// delete outputs's utxo
					preout := modules.OutPoint{TxHash: hash}
					for j := range payment.Outputs {
						preout.MessageIndex = uint32(i)
						preout.OutIndex = uint32(j)
						pool.deleteOrphanTxOutputs(preout)
					}
				}
			}
		}
	}
	return nil
}

// removeTx removes a single transaction from the queue, moving all subsequent
// transactions back to the future queue.
func (pool *TxPool) removeTx(hash common.Hash) {
	// Fetch the transaction we wish to delete
	interTx, has := pool.all.Load(hash)
	if !has {
		return
	}
	tx, ok := interTx.(*TxPoolTransaction)
	if !ok {
		return
	}
	tx.Discarded = true
	pool.all.Store(hash, tx)

	for i, msg := range tx.Tx.TxMessages() {
		if msg.App != modules.APP_PAYMENT {
			continue
		}
		payment, ok := msg.Payload.(*modules.PaymentPayload)
		if !ok {
			continue
		}
		for _, input := range payment.Inputs {
			// 排除手续费的输入为nil
			if input.PreviousOutPoint != nil {
				pool.outpoints.Delete(*input.PreviousOutPoint)
			}
		}
		// delete outputs's utxo
		preout := modules.OutPoint{TxHash: hash}
		for j := range payment.Outputs {
			preout.MessageIndex = uint32(i)
			preout.OutIndex = uint32(j)
			pool.deleteOrphanTxOutputs(preout)
		}
	}
}

// 标记该交易及引用的交易为discard。
func (pool *TxPool) removeTransaction(tx *TxPoolTransaction, removeRedeemers bool) {
	hash := tx.Tx.Hash()
	if removeRedeemers {
		// 删除所有引用该交易的交易。
		for i, msgcopy := range tx.Tx.TxMessages() {
			if msgcopy.App == modules.APP_PAYMENT {
				if msg, ok := msgcopy.Payload.(*modules.PaymentPayload); ok {
					for j := uint32(0); j < uint32(len(msg.Outputs)); j++ {
						preout := modules.OutPoint{TxHash: hash, MessageIndex: uint32(i), OutIndex: j}
						if pooltxRedeemer, exist := pool.outpoints.Load(preout); exist {
							pool.removeTransaction(pooltxRedeemer.(*TxPoolTransaction), true)
						}
					}
				}
			}
		}
	}
	// Remove the transaction if needed.
	_, has := pool.all.Load(hash)
	if !has {
		return
	}
	// mark the referenced outpoints as unspent by the pool.
	for _, msgcopy := range tx.Tx.TxMessages() {
		if msgcopy.App == modules.APP_PAYMENT {
			if msg, ok := msgcopy.Payload.(*modules.PaymentPayload); ok {
				for _, input := range msg.Inputs {
					pool.outpoints.Delete(*input.PreviousOutPoint)
				}
			}
		}
	}
	tx.Discarded = true
	pool.all.Store(hash, tx)
}

// promoteExecutables moves transactions that have become processable from the
// future queue to the set of pending transactions. During this process, all
// invalidated transactions (low nonce, low balance) are deleted.
func (pool *TxPool) promoteExecutables() {
	// If the pending limit is overflown, start equalizing allowances
	pendingTxs := make([]*TxPoolTransaction, 0)
	poolTxs := pool.AllTxpoolTxs()
	for _, tx := range poolTxs {
		if tx.Pending {
			continue
		}
		pendingTxs = append(pendingTxs, tx)
	}
	pending := len(pendingTxs)
	if uint64(pending) > pool.config.GlobalSlots {
		// Assemble a spam order to penalize large transactors first
		spammers := prque.New()
		for i, tx := range pendingTxs {
			// Only evict transactions from high rollers
			spammers.Push(tx.Tx.Hash(), float32(i))
		}
		// Gradually drop transactions from offenders
		offenders := []common.Hash{}
		for uint64(pending) > pool.config.GlobalQueue && !spammers.Empty() {
			// Retrieve the next offender if not local address
			offender, _ := spammers.Pop()
			offenders = append(offenders, offender.(common.Hash))

			// Equalize balances until all the same or below threshold
			if len(offenders) > 1 {
				// Iteratively reduce all offenders until below limit or threshold reached
				for uint64(pending) > pool.config.GlobalQueue {
					for i := 0; i < len(offenders)-1; i++ {
						for _, tx := range pendingTxs {
							hash := tx.Tx.Hash()
							if offenders[i].String() == hash.String() {
								// Drop the transaction from the global pools too
								pool.all.Delete(hash)
								pool.priority_sorted.Removed()
								log.Trace("Removed fairness-exceeding pending transaction", "hash", hash)
								pending--
								break
							}
						}
					}
				}
			}
		}
		// If still above threshold, reduce to limit or min allowance
		if uint64(pending) > pool.config.GlobalQueue && len(offenders) > 0 {
			for uint64(pending) > pool.config.GlobalQueue {
				for _, addr := range offenders {
					for _, tx := range pendingTxs {
						hash := tx.Tx.Hash()
						if addr.String() == hash.String() {
							pool.all.Delete(hash)
							pool.priority_sorted.Removed()
							log.Trace("Removed fairness-exceeding pending transaction", "hash", hash)
							pending--
							break
						}
					}
				}
			}
		}
	}
}

// Stop terminates the transaction pool.
func (pool *TxPool) Stop() {
	pool.scope.Close()
	// pool.wg.Wait()
	if pool.journal != nil {
		pool.journal.close()
	}
	log.Info("Transaction pool stopped")
}

// 打包后的没有被最终确认的交易，废弃处理
func (pool *TxPool) DiscardTxs(txs []*modules.Transaction) error {
	pool.mu.Lock()
	defer pool.mu.Unlock()
	for _, tx := range txs {
		requestHash := tx.RequestHash()
		if tx.IsContractTx() {
			if tx.IsSystemContract() {
				err := pool.discardTx(requestHash)
				if err != nil && err != ErrNotFound {
					log.Warnf("Req[%s] discard error:%s", requestHash.String(), err.Error())
				}
			}
			pool.orphans.Delete(requestHash)
			//删除对应的Request,可能有后续Tx在孤儿池，添加回来
			if _, ok := pool.userContractRequests[requestHash]; ok {
				log.Debugf("Request[%s] already packed into unit, delete it from request pool", requestHash.String())
				delete(pool.userContractRequests, requestHash)
				pool.checkBasedOnReqOrphanTxToNormal(tx.Hash(), requestHash)

			}
		}
		err := pool.discardTx(tx.Hash())
		if err != nil && err != ErrNotFound {
			log.Warnf("Tx[%s] discard error:%s", tx.Hash().String(), err.Error())
		}
		// 删除孤儿 txhash
		pool.orphans.Delete(tx.Hash())
	}
	return nil
}

func (pool *TxPool) discardTx(hash common.Hash) error {
	// in all pool
	interTx, has := pool.all.Load(hash)
	if !has {
		return nil
	}
	tx := interTx.(*TxPoolTransaction)
	tx.Discarded = true
	pool.deletePoolUtxos(tx.Tx)
	pool.all.Store(hash, tx)
	return nil
}
func (pool *TxPool) SetPendingTxs(unit_hash common.Hash, num uint64, txs []*modules.Transaction) error {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	for i, tx := range txs {
		if i == 0 { // coinbase
			continue
		}
		err := pool.setPendingTx(unit_hash, tx, num, uint64(i))
		if err != nil {
			return err
		}
	}
	if len(txs) > 1 {
		pool.priority_sorted.Removed()
	}
	return nil
}
func (pool *TxPool) setPendingTx(unit_hash common.Hash, tx *modules.Transaction, number, index uint64) error {
	// convert
	p_tx := pool.convertBaseTx(tx)
	//如果是系统合约，那么需要按RequestHash去查找并改变状态
	if tx.IsSystemContract() {
		if !pool.isTransactionInPool(tx.RequestHash()) {
			//如果有交易没有出现在交易池中，则直接补充
			e := pool.addLocal(tx.GetRequestTx())
			if e != nil {
				return e
			}
		}
		// 更新交易的状态
		pool.promoteTx(unit_hash, tx.RequestHash(), p_tx, number, index)
	} else {
		if !pool.isTransactionInPool(tx.Hash()) {
			//如果有交易没有出现在交易池中，则直接补充
			e := pool.addLocal(tx)
			if e != nil {
				return e
			}
		}
		// 更新交易的状态
		pool.promoteTx(unit_hash, tx.Hash(), p_tx, number, index)
	}

	return nil
}

func (pool *TxPool) addCache(tx *TxPoolTransaction) {
	if tx == nil {
		return
	}
	txHash := tx.Tx.Hash()
	reqHash := tx.Tx.RequestHash()
	for i, msgcopy := range tx.Tx.Messages() {
		if msgcopy.App != modules.APP_PAYMENT {
			continue
		}
		msg, ok := msgcopy.Payload.(*modules.PaymentPayload)
		if !ok {
			continue
		}
		for _, txin := range msg.Inputs {
			if txin.PreviousOutPoint != nil {
				pool.outpoints.Store(*txin.PreviousOutPoint, tx)
			}
		}
		// add  outputs
		preout := modules.OutPoint{TxHash: txHash}
		for j, out := range msg.Outputs {
			preout.MessageIndex = uint32(i)
			preout.OutIndex = uint32(j)
			utxo := &modules.Utxo{Amount: out.Value, Asset: &modules.Asset{
				AssetId: out.Asset.AssetId, UniqueId: out.Asset.UniqueId},
				PkScript: out.PkScript[:]}
			pool.outputs.Store(preout, utxo)
			if txHash != reqHash {
				preout.TxHash = reqHash
				pool.reqOutputs.Store(preout, utxo)
			}
		}
	}
}
func (pool *TxPool) ResetPendingTxs(txs []*modules.Transaction) error {
	pool.mu.Lock()
	defer pool.mu.Unlock()
	for i, tx := range txs {
		if i == 0 { //coinbase
			continue
		}
		pool.resetPendingTx(tx)
	}
	return nil
}
func (pool *TxPool) resetPendingTx(tx *modules.Transaction) error {
	hash := tx.Hash()
	pool.DeleteTxByHash(hash)

	_, err := pool.add(tx, !pool.config.NoLocals)
	return err
}

/******  end utxoSet  *****/
// GetSortedTxs returns 根据优先级返回list
func (pool *TxPool) GetSortedTxs() ([]*TxPoolTransaction, error) {
	t0 := time.Now()
	canbe_packaged := false
	var total common.StorageSize
	list := make([]*TxPoolTransaction, 0)

	gasAsset := dagconfig.DagConfig.GetGasToken()
	_, chainindex, err := pool.unit.GetNewestUnit(gasAsset)
	if err != nil {
		return nil, err
	}
	unithigh := int64(chainindex.Index)
	map_pretxs := make(map[common.Hash]int)
	// get sequenTxs
	stxs := pool.GetSequenTxs()
	poolTxs := pool.AllTxpoolTxs()
	orphanTxs := pool.AllOrphanTxs()
	unit_size := common.StorageSize(parameter.CurrentSysParameters.UnitMaxSize)
	for _, tx := range stxs {
		list = append(list, tx)
		total += tx.Tx.Size()
	}
	for {
		if time.Since(t0) > time.Millisecond*800 {
			log.Infof("get sorted timeout spent times: %s , count: %d ", time.Since(t0), len(list))
			break
		}
		if total >= unit_size {
			break
		}
		tx := pool.priority_sorted.Get()
		if tx == nil {
			log.Debugf("The task of txspool get priority_pricedtx has been finished,count:%d", len(list))
			break
		} else {
			if !tx.Pending {
				if has, _ := pool.unit.IsTransactionExist(tx.Tx.Hash()); has {
					continue
				}
				// add precusorTxs 获取该交易的前驱交易列表
				p_txs := pool.getPrecusorTxs(tx, poolTxs)
				for _, p_tx := range p_txs {
					if _, has := map_pretxs[p_tx.Tx.Hash()]; !has {
						map_pretxs[p_tx.Tx.Hash()] = len(list)
						if !p_tx.Pending {
							list = append(list, p_tx)
							total += p_tx.Tx.Size()
						}
					}
				}
			}
		}
	}
	//  验证孤儿交易
	or_list := make(orList, 0)
	for _, tx := range orphanTxs {
		or_list = append(or_list, tx)
	}
	// 按入池时间排序
	if len(or_list) > 1 {
		sort.Sort(or_list)
	}
	// pool rlock
	pool.mu.RLock()
	defer pool.mu.RUnlock()
	for _, tx := range or_list {
		txhash := tx.Tx.Hash()
		if has, _ := pool.unit.IsTransactionExist(txhash); has {
			pool.orphans.Delete(txhash)
			continue
		}
		locktime := tx.Tx.GetLocktime()
		if locktime > 0 {
			if locktime < 500000000 && unithigh >= locktime {
				canbe_packaged = true
			} else if locktime < 500000000 && unithigh < locktime {
				canbe_packaged = false
			}
			if (locktime >= 500000000 && locktime-time.Now().Unix() < 0) || canbe_packaged {
				tx.IsOrphan = false
				pool.all.Store(txhash, tx)
				pool.orphans.Delete(txhash)
				list = append(list, tx)
				total += tx.Tx.Size()
				if total > unit_size {
					break
				}
			}
			continue
		}
		ok, err := pool.ValidateOrphanTx(tx.Tx)
		if !ok && err == nil {
			//  更改孤儿交易的状态
			tx.IsOrphan = false
			pool.all.Store(txhash, tx)
			pool.orphans.Delete(txhash)
			list = append(list, tx)
			total += tx.Tx.Size()
			if total > unit_size {
				break
			}
		}
	}

	// 	去重
	m := make(map[common.Hash]*TxPoolTransaction)
	indexL := make(map[int]common.Hash)
	for i, tx := range list {
		hash := tx.Tx.Hash()
		tx.Index = uint64(i)
		indexL[i] = hash
		m[hash] = tx
	}
	list = make([]*TxPoolTransaction, 0)
	for i := 0; i < len(indexL); i++ {
		t_hash := indexL[i]
		if tx, has := m[t_hash]; has {
			delete(m, t_hash)
			list = append(list, tx)
		}
	}

	toSorted := make(map[common.Hash]*modules.Transaction)
	for _, tx := range list {
		toSorted[tx.Tx.Hash()] = tx.Tx
	}
	sorted, orphans, double := modules.SortTxs(toSorted, pool.GetUtxoEntry)
	for _, tx := range orphans {
		pool.addOrphan(tx)
	}
	for _, tx := range double {
		pool.DeleteTxByHash(tx.Hash())
	}
	list = make([]*TxPoolTransaction, 0)
	for _, tx := range sorted {
		list = append(list, pool.convertBaseTx(tx))
	}

	log.DebugDynamic(func() string {
		var str string
		for i, tx := range sorted {
			str += fmt.Sprintf("index:%d hash:%s\n", i, tx.Hash().String())
		}
		return str
	})
	log.Debugf("get sorted and rm Orphan txs spent times: %s , count: %d ,sorted: %d , txs_size %s,  "+
		"total_size %s", time.Since(t0), len(list), len(sorted), total.String(), unit_size.String())
	return list, nil
}
func (pool *TxPool) getPrecusorTxs(tx *TxPoolTransaction, poolTxs map[common.Hash]*TxPoolTransaction) []*TxPoolTransaction {
	pretxs := make([]*TxPoolTransaction, 0)
	for _, op := range tx.Tx.GetSpendOutpoints() {
		// 交易池做了utxo的缓存，包括request交易的缓存utxo，不能用pool.GetUtxoEntry
		_, err := pool.unit.GetUtxoEntry(op)
		if err == nil {
			continue
		}
		//  若该utxo在db里找不到,try to find it in pool and ophans txs
		queue_tx, has := poolTxs[op.TxHash]
		if !has {
		poolloop:
			for _, otx := range poolTxs {
				if otx.Tx.RequestHash() != op.TxHash {
					continue
				}
				for i, msg := range otx.Tx.Messages() {
					if msg.App != modules.APP_PAYMENT {
						continue
					}
					payment := msg.Payload.(*modules.PaymentPayload)
					for j := range payment.Outputs {
						if op.OutIndex == uint32(j) && op.MessageIndex == uint32(i) {
							queue_tx = otx
							break poolloop
						}
					}
				}
			}
		}
		if queue_tx == nil {
			continue
		}
		if queue_tx != nil || queue_tx.Pending {
			continue
		}
		//if find precusor tx  ,and go on to find its
		log.Info("find in precusor tx.", "hash", queue_tx.Tx.Hash().String(), "ohash", op.TxHash.String(),
			"pending", tx.Pending)
		list := pool.getPrecusorTxs(queue_tx, poolTxs)
		for _, p_tx := range list {
			pretxs = append(pretxs, p_tx)
			delete(poolTxs, p_tx.Tx.Hash())
		}
	}

	pretxs = append(pretxs, tx)
	return pretxs
}
func (pool *TxPool) GetSequenTxs() []*TxPoolTransaction {
	return pool.getSequenTxs()
}
func (pool *TxPool) getSequenTxs() []*TxPoolTransaction {
	return pool.sequenTxs.All()
}

type orList []*TxPoolTransaction

func (ol orList) Len() int {
	return len(ol)
}
func (ol orList) Swap(i, j int) {
	ol[i], ol[j] = ol[j], ol[i]
}
func (ol orList) Less(i, j int) bool {
	return ol[i].CreationDate.Unix() < ol[j].CreationDate.Unix()
}

// SubscribeTxPreEvent registers a subscription of TxPreEvent and
// starts sending event to the given channel.
func (pool *TxPool) SubscribeTxPreEvent(ch chan<- modules.TxPreEvent) event.Subscription {
	return pool.scope.Track(pool.txFeed.Subscribe(ch))
}

func (pool *TxPool) limitNumberOrphans() {
	// scan the orphan pool and remove any expired orphans when it's time.
	orphanTxs := pool.AllOrphanTxs()
	if now := time.Now(); now.After(pool.nextExpireScan) {
		originNum := len(orphanTxs)
		for _, tx := range orphanTxs {
			if now.After(tx.Expiration) {
				// remove
				pool.removeOrphan(tx, true)
			}
			ok, err := pool.ValidateOrphanTx(tx.Tx)
			if !ok && err == nil {
				pool.add(tx.Tx, !pool.config.NoLocals)
			}
		}
		// set next expireScan
		pool.nextExpireScan = time.Now().Add(pool.config.OrphanTTL)
		numOrphans := len(pool.AllOrphanTxs())

		if numExpied := originNum - numOrphans; numExpied > 0 {
			log.Debug(fmt.Sprintf("Expired %d %s (remaining: %d)", numExpied, pickNoun(numExpied,
				"orphan", "orphans"), numOrphans))
		}
	}
	// nothing to do if adding another orphan will not cause the pool to exceed the limit
	if len(pool.AllOrphanTxs())+1 <= pool.config.MaxOrphanTxs {
		return
	}

	// remove a random entry from the map.
	for _, tx := range orphanTxs {
		pool.removeOrphan(tx, false)
		break
	}
}

// pickNoun returns the singular or plural form of a noun depending
// on the count n.
func pickNoun(n int, singular, plural string) string {
	if n == 1 {
		return singular
	}
	return plural
}

func (pool *TxPool) addOrphan(otx *modules.Transaction) {
	if pool.config.MaxOrphanTxs <= 0 {
		return
	}
	tx := pool.convertBaseTx(otx)
	tx.IsOrphan = true
	pool.orphans.Store(tx.Tx.Hash(), tx)
	log.Debugf("Stored orphan tx's hash:[%s] (total: %d)", tx.Tx.Hash().String(), len(pool.AllOrphanTxs()))
}

func (pool *TxPool) removeOrphan(tx *TxPoolTransaction, reRedeemers bool) {
	hash := tx.Tx.Hash()
	orphanTxs := pool.AllOrphanTxs()
	otx, has := orphanTxs[hash]
	if !has {
		return
	}

	for _, msg := range otx.Tx.TxMessages() {
		if msg.App == modules.APP_PAYMENT {
			payment, ok := msg.Payload.(*modules.PaymentPayload)
			if ok {
				for _, in := range payment.Inputs {
					pool.deleteOrphanTxOutputs(*in.PreviousOutPoint)
				}
			}
		}
	}
	// remove any orphans that redeem outputs from this one if requested.
	if !reRedeemers {
		pool.orphans.Delete(hash)
		return
	}
	prevOut := modules.OutPoint{TxHash: hash}
	for i, msg := range tx.Tx.TxMessages() {
		if msg.App != modules.APP_PAYMENT {
			continue
		}
		payment, ok := msg.Payload.(*modules.PaymentPayload)
		if !ok {
			continue
		}
		for j := range payment.Outputs {
			prevOut.MessageIndex = uint32(i)
			prevOut.OutIndex = uint32(j)
			pool.outputs.Delete(prevOut)
		}
	}
	// remove the transaction from the orphan pool.
	pool.orphans.Delete(hash)
}

// isOrphanInPool returns whether or not the passed transaction already exists
// in the orphan pool.
//
// This function MUST be called with the mempool lock held (for reads).
func (pool *TxPool) isOrphanInPool(hash common.Hash) bool {
	if _, exists := pool.orphans.Load(hash); exists {
		return true
	}
	return false
}

// validate tx is an orphanTx or not.
func (pool *TxPool) ValidateOrphanTx(tx *modules.Transaction) (bool, error) {
	gasToken := dagconfig.DagConfig.GetGasToken()

	_, chainindex, err := pool.unit.GetNewestUnit(gasToken)
	if err != nil {
		return false, errors.New("can not get GetNewestUnit.")
	}
	unithigh := int64(chainindex.Index)

	var isOrphan bool
	for _, msg := range tx.Messages() {
		if isOrphan {
			break
		}
		if msg.App != modules.APP_PAYMENT {
			continue
		}
		payment, ok := msg.Payload.(*modules.PaymentPayload)
		if !ok {
			continue
		}
		if payment.LockTime > 500000000 && (int64(payment.LockTime)-time.Now().Unix()) < 0 {
			isOrphan = false
			break
		} else if payment.LockTime > 500000000 && (int64(payment.LockTime)-time.Now().Unix()) >= 0 {
			isOrphan = true
			break
		} else if payment.LockTime > 0 && payment.LockTime < 500000000 && (int64(payment.LockTime) < unithigh) {
			// if persent unit is high than lock unit ,not Orphan
			isOrphan = false
			break
		} else if payment.LockTime > 0 && payment.LockTime < 500000000 && (int64(payment.LockTime) > unithigh) {
			// if persent unit is low than lock unit ,not Orphan
			isOrphan = true
			break
		}

		for _, in := range payment.Inputs {
			if _, err := pool.unit.GetUtxoEntry(in.PreviousOutPoint); err != nil {
				_, err = pool.GetUtxoEntry(in.PreviousOutPoint)
				if err != nil {
					log.Debugf("get utxo failed,%s", in.PreviousOutPoint.String())
					return true, err
				}
			}
		}
	}

	return isOrphan, nil
}

func (pool *TxPool) deleteOrphanTxOutputs(outpoint modules.OutPoint) {
	pool.outputs.Delete(outpoint)
	// 删除缓存的req utxo
	pool.reqOutputs.Delete(outpoint)
}

func (pool *TxPool) deletePoolUtxos(tx *modules.Transaction) {
	for _, msg := range tx.Messages() {
		if msg.App == modules.APP_PAYMENT {
			payment, ok := msg.Payload.(*modules.PaymentPayload)
			if ok {
				for _, in := range payment.Inputs {
					pool.deleteOrphanTxOutputs(*in.PreviousOutPoint)
				}
			}
		}
		payment, ok := msg.Payload.(*modules.PaymentPayload)
		if !ok {
			continue
		}
		for _, in := range payment.Inputs {
			pool.deleteOrphanTxOutputs(*in.PreviousOutPoint)
			// 删除缓存的req utxo
			pool.reqOutputs.Delete(*in.PreviousOutPoint)
		}
	}
}
func (pool *TxPool) checkBasedOnReqOrphanTxToNormal(ori_hash, ori_reqhash common.Hash) error {
	for hash, otx := range pool.basedOnRequestOrphans {
		if !pool.isBasedOnRequestPool(otx) && (otx.IsDependOnTx(ori_hash) || otx.IsDependOnTx(ori_reqhash)) {
			//满足Normal的条件了
			log.Debugf("move tx[%s] from based on request orphans to normals", otx.TxHash.String())
			delete(pool.basedOnRequestOrphans, hash)         //从孤儿池删除
			err := pool.addTx(otx.Tx, !pool.config.NoLocals) //因为之前孤儿交易没有手续费，UTXO等，所以需要重新计算
			if err != nil {
				log.Warnf("add tx[%s] to pool fail:%s", hash.String(), err.Error())
			}
		}
	}
	return nil
}

func (pool *TxPool) reflashOrphanTxs(tx *modules.Transaction, orphans map[common.Hash]*TxPoolTransaction, local bool) {
	tx_hash := tx.Hash()
	req_hash := tx.RequestHash()
	for hash, otx := range orphans {
		isOrphan := false
		for _, op := range otx.Tx.GetSpendOutpoints() {
			if _, err := pool.unit.GetUtxoEntry(op); err != nil {
				if _, err := pool.GetUtxoEntry(op); err != nil {
					if op.TxHash != tx_hash && op.TxHash != req_hash {
						isOrphan = true
						break
					}
				}
			}
		}
		if !isOrphan { //该交易不再是孤儿交易，使之变为有效交易。
			pool.orphans.Delete(hash)
			if err := pool.addTx(otx.Tx, local); err != nil {
				log.Debugf("addlocal failed,error:%s,hash:%s", err.Error(), otx.Tx.Hash().String())
				pool.orphans.Store(hash, otx)
			}
		}
		////该交易不再是孤儿交易，使之变为有效交易。
		//log.Infof("reflash orphan tx[%s] goto packaged.", hash.String())
		//pool.priority_sorted.Put(otx)
		//pool.orphans.Delete(hash)
		//pool.all.Store(hash, otx)
		//pool.addCache(otx)
	}
}
func (pool *TxPool) GetAddrUtxos(addr common.Address, token *modules.Asset) (
	map[modules.OutPoint]*modules.Utxo, error) {
	dbUtxos, dbReqTxMapping, err := pool.unit.GetAddrUtxoAndReqMapping(addr, token)
	if err != nil {
		return nil, err
	}
	log.DebugDynamic(func() string {
		utxoKeys := ""
		for o := range dbUtxos {
			utxoKeys += o.String() + ";"
		}
		mapping := ""
		for req, tx := range dbReqTxMapping {
			mapping += req.String() + ":" + tx.String() + ";"
		}
		return "db utxo outpoints:" + utxoKeys + " req:tx mapping :" + mapping
	})
	txs, err := pool.GetUnpackedTxsByAddr(addr)
	if err != nil {
		return nil, err
	}
	log.DebugDynamic(func() string {
		txHashs := ""
		for _, tx := range txs {
			txHashs += "[tx:" + tx.Tx.Hash().String() + "-req:" + tx.Tx.RequestHash().String() + "];"
		}
		return "txpool unpacked tx:" + txHashs
	})
	poolUtxo, poolReqTxMapping, poolSpend := parseTxUtxo(txs, addr, token)
	for k, v := range dbUtxos {
		poolUtxo[k] = v
	}
	for k, v := range dbReqTxMapping {
		poolReqTxMapping[k] = v
	}
	for spend := range poolSpend {
		delete(poolUtxo, spend)
		//删除引用request 的utxo
		if txHash, ok := poolReqTxMapping[spend.TxHash]; ok {
			spend2 := modules.OutPoint{
				TxHash:       txHash,
				MessageIndex: spend.MessageIndex,
				OutIndex:     spend.OutIndex,
			}
			delete(poolUtxo, spend2)
		}
	}
	//删除poolutxo里已重复引用的request utxo。
	for outpoint := range poolUtxo {
		if txHash, ok := poolReqTxMapping[outpoint.TxHash]; ok {
			if txHash != outpoint.TxHash {
				delete(poolUtxo, outpoint)
			}
		}
	}
	return poolUtxo, nil
}

func parseTxUtxo(txs []*TxPoolTransaction, addr common.Address, token *modules.Asset) (
	map[modules.OutPoint]*modules.Utxo, map[common.Hash]common.Hash, map[modules.OutPoint]bool) {
	dbUtxos := make(map[modules.OutPoint]*modules.Utxo)
	spendUtxo := make(map[modules.OutPoint]bool)
	dbReqTxMapping := make(map[common.Hash]common.Hash)
	lockScript := tokenengine.Instance.GenerateLockScript(addr)
	for _, tx := range txs {
		for k, v := range tx.Tx.GetNewUtxos() {
			if !bytes.Equal(lockScript, v.PkScript) {
				continue
			}
			if token != nil && v.Asset.Equal(token) {
				dbUtxos[k] = v
			}
		}
		for _, so := range tx.Tx.GetSpendOutpoints() {
			spendUtxo[*so] = true
		}
		if tx.TxHash != tx.ReqHash {
			dbReqTxMapping[tx.ReqHash] = tx.TxHash
		}
	}
	return dbUtxos, dbReqTxMapping, spendUtxo
}

func (pool *TxPool) convertBaseTx(tx *modules.Transaction) *TxPoolTransaction {
	dependOnTxs := make(map[common.Hash]bool)
	for _, o := range tx.GetSpendOutpoints() {
		dependOnTxs[o.TxHash] = false
	}
	txAddr, _ := tx.GetToAddrs(pool.tokenEngine.GetAddressFromScript)
	return &TxPoolTransaction{
		Tx:           tx,
		TxHash:       tx.Hash(),
		ReqHash:      tx.RequestHash(),
		CreationDate: time.Now(),
		DependOnTxs:  dependOnTxs,
		From:         tx.GetSpendOutpoints(),
		ToAddr:       txAddr,
		//IsSysContractRequest: tx.IsOnlyContractRequest() && tx.IsSystemContract(),
		//IsUserContractFullTx: tx.IsUserContract() && !tx.IsOnlyContractRequest(),
	}
}

func (pool *TxPool) convertTx(tx *modules.Transaction, fee []*modules.Addition) *TxPoolTransaction {
	fromAddr, _ := tx.GetFromAddrs(pool.GetUtxoEntry, pool.tokenEngine.GetAddressFromScript)
	tx2 := pool.convertBaseTx(tx)
	tx2.TxFee = fee
	tx2.FromAddr = fromAddr
	return tx2
}

func (pool *TxPool) isBasedOnRequestPool(tx *TxPoolTransaction) bool {
	for h := range tx.DependOnTxs {
		if _, ok := pool.userContractRequests[h]; ok {
			return true
		}
		if _, ok := pool.basedOnRequestOrphans[h]; ok {
			return true
		}
	}
	return false
}

func (pool *TxPool) addBasedOnReqOrphanTx(tx *TxPoolTransaction) error {
	log.Debugf("add tx[%s] to based on request orphan pool", tx.TxHash.String())
	tx.Status = TxPoolTxStatus_Orphan
	pool.basedOnRequestOrphans[tx.TxHash] = tx
	pool.txFeed.Send(modules.TxPreEvent{Tx: tx.Tx, IsOrphan: false})
	return nil
}
