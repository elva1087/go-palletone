// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

// Package ptn implements the PalletOne protocol.
package ptn

import (
	"fmt"
	"sync"

	"github.com/palletone/go-palletone/common"
	"github.com/palletone/go-palletone/common/bloombits"
	"github.com/palletone/go-palletone/common/event"
	"github.com/palletone/go-palletone/common/log"
	"github.com/palletone/go-palletone/common/p2p"
	palletdb "github.com/palletone/go-palletone/common/ptndb"
	"github.com/palletone/go-palletone/common/rpc"
	"github.com/palletone/go-palletone/consensus"
	"github.com/palletone/go-palletone/consensus/jury"
	mp "github.com/palletone/go-palletone/consensus/mediatorplugin"
	"github.com/palletone/go-palletone/contracts"
	"github.com/palletone/go-palletone/core"
	"github.com/palletone/go-palletone/core/accounts"
	"github.com/palletone/go-palletone/core/accounts/keystore"
	"github.com/palletone/go-palletone/core/node"
	"github.com/palletone/go-palletone/dag"
	"github.com/palletone/go-palletone/dag/dagconfig"
	"github.com/palletone/go-palletone/dag/modules"
	"github.com/palletone/go-palletone/dag/storage"
	"github.com/palletone/go-palletone/dag/txspool"
	"github.com/palletone/go-palletone/internal/ptnapi"
	"github.com/palletone/go-palletone/ptn/downloader"
	"github.com/palletone/go-palletone/ptn/filters"
	"github.com/palletone/go-palletone/ptn/indexer"
	"github.com/palletone/go-palletone/ptnjson"
	"github.com/palletone/go-palletone/tokenengine"
	"github.com/shopspring/decimal"
)

type LesServer interface {
	Start(srvr *p2p.Server)
	Stop()
	Protocols() []p2p.Protocol
	SetBloomBitsIndexer(bbIndexer *indexer.ChainIndexer)
}

// PalletOne implements the PalletOne full node service.
type PalletOne struct {
	config *Config

	// Channel for shutting down the service
	shutdownChan chan bool // Channel for shutting down the PalletOne

	// Handlers
	txPool          txspool.ITxPool
	protocolManager *ProtocolManager
	lesServer       LesServer

	eventMux       *event.TypeMux
	engine         core.ConsensusEngine
	accountManager *accounts.Manager

	ApiBackend *PtnApiBackend

	networkId     uint64
	netRPCService *ptnapi.PublicNetAPI

	dag dag.IDag

	contract *contracts.Contract

	lock sync.RWMutex // Protects the variadic fields (e.g. gas price and etherbase)

	bloomRequests chan chan *bloombits.Retrieval // Channel receiving bloom data retrieval requests
	bloomIndexer  *indexer.ChainIndexer          // Bloom indexer operating during block imports

	// append by Albert·Gou
	mediatorPlugin    *mp.MediatorPlugin
	contractPorcessor *jury.Processor
}

func (p *PalletOne) AddLesServer(ls LesServer) {
	p.lesServer = ls
	ls.SetBloomBitsIndexer(p.bloomIndexer)
}

// New creates a new PalletOne object (including the
// initialisation of the common PalletOne object)
func New(ctx *node.ServiceContext, config *Config) (*PalletOne, error) {
	if !config.SyncMode.IsValid() {
		return nil, fmt.Errorf("invalid sync mode %d", config.SyncMode)
	}

	db, err := CreateDB(ctx, config)
	if err != nil {
		log.Error("PalletOne New", "CreateDB err:", err)
		return nil, err
	}
	dag, err := dag.NewDag(db)
	if err != nil {
		log.Error("PalletOne New", "NewDag err:", err)
		return nil, err
	}

	ptn := &PalletOne{
		config:         config,
		eventMux:       ctx.EventMux,
		accountManager: ctx.AccountManager,
		shutdownChan:   make(chan bool),
		networkId:      config.NetworkId,
		dag:            dag,
		bloomRequests:  make(chan chan *bloombits.Retrieval),
		bloomIndexer:   NewBloomIndexer(db, BloomBitsBlocks),
	}
	log.Info("Initialising PalletOne protocol", "versions", ProtocolVersions, "network", config.NetworkId)

	if config.TxPool.Journal != "" {
		config.TxPool.Journal = ctx.ResolvePath(config.TxPool.Journal)
	}
	pool := txspool.NewTxPool(config.TxPool, ptn.dag)
	ptn.txPool = pool

	//Test for P2P
	ptn.engine = consensus.New(dag, pool)

	ptn.mediatorPlugin, err = mp.NewMediatorPlugin(ptn, dag, &config.MediatorPlugin)
	if err != nil {
		log.Error("Initialize mediator plugin err:", "error", err)
		return nil, err
	}

	ptn.contractPorcessor, err = jury.NewContractProcessor(ptn, dag, nil, &config.Jury)
	if err != nil {
		log.Error("contract processor creat:", "error", err)
		return nil, err
	}

	aJury := &consensus.AdapterJury{ptn.contractPorcessor}
	ptn.contract, err = contracts.Initialize(ptn.dag, aJury, &config.Contract)
	if err != nil {
		log.Error("Contract Initialize err:", "error", err)
		return nil, err
	}
	ptn.contractPorcessor.SetContract(ptn.contract)

	genesis, err := ptn.dag.GetGenesisUnit()
	if err != nil {
		log.Error("PalletOne New", "get genesis err:", err)
		return nil, err
	}

	gasToken := config.Dag.GetGasToken()
	ptn.bloomIndexer.Start(dag, gasToken)
	if ptn.protocolManager, err = NewProtocolManager(config.SyncMode, config.NetworkId, gasToken, ptn.txPool,
		ptn.dag, ptn.eventMux, ptn.mediatorPlugin, genesis, ptn.contractPorcessor, ptn.engine); err != nil {
		log.Error("NewProtocolManager err:", "error", err)
		return nil, err
	}

	ptn.ApiBackend = &PtnApiBackend{ptn}
	return ptn, nil
}

// CreateDB creates the chain database.
func CreateDB(ctx *node.ServiceContext, config *Config /*, name string*/) (palletdb.Database, error) {
	//db, err := ptndb.NewLDBDatabase(ctx.config.resolvePath(name), cache, handles)
	//path := ctx.DatabasePath(name)
	path := dagconfig.DagConfig.DbPath

	//fit dag DefaultConfig
	//dagconfig.DagConfig.DbPath = path

	log.Debug("Open leveldb path:", "path", path)
	db, err := storage.Init(path, config.DatabaseCache, config.DatabaseHandles)
	if err != nil {
		return nil, err
	}
	//	if db, ok := db.(*palletdb.LDBDatabase); ok {
	//		db.Meter("eth/db/chaindata/")
	//	}
	return db, nil
}

// APIs returns the collection of RPC services the ethereum package offers.
// NOTE, some of these services probably need to be moved to somewhere else.
func (s *PalletOne) APIs() []rpc.API {
	apis := ptnapi.GetAPIs(s.ApiBackend)

	// append by Albert·Gou
	apis = append(apis, s.mediatorPlugin.APIs()...)

	// Append all the local APIs and return
	return append(apis, []rpc.API{
		{
			Namespace: "ptn",
			Version:   "1.0",
			Service:   NewPublicPalletOneAPI(s),
			Public:    true,
		},
		{
			Namespace: "ptn",
			Version:   "1.0",
			Service:   downloader.NewPublicDownloaderAPI(s.protocolManager.downloader, s.eventMux),
			Public:    true,
		},
		{
			Namespace: "ptn",
			Version:   "1.0",
			Service:   filters.NewPublicFilterAPI(s.ApiBackend, false),
			Public:    true,
		}, {
			Namespace: "admin",
			Version:   "1.0",
			//Service:   NewPrivateAdminAPI(s),
		}, {
			Namespace: "net",
			Version:   "1.0",
			Service:   s.netRPCService,
			Public:    true,
		},
	}...)
}

func (s *PalletOne) AccountManager() *accounts.Manager { return s.accountManager }
func (s *PalletOne) TxPool() txspool.ITxPool           { return s.txPool }
func (s *PalletOne) EventMux() *event.TypeMux          { return s.eventMux }

func (s *PalletOne) Engine() core.ConsensusEngine       { return s.engine }
func (s *PalletOne) IsListening() bool                  { return true } // Always listening
func (s *PalletOne) EthVersion() int                    { return int(s.protocolManager.SubProtocols[0].Version) }
func (s *PalletOne) NetVersion() uint64                 { return s.networkId }
func (s *PalletOne) Downloader() *downloader.Downloader { return s.protocolManager.downloader }
func (s *PalletOne) Dag() dag.IDag                      { return s.dag }

func (s *PalletOne) ContractProcessor() *jury.Processor { return s.contractPorcessor }
func (s *PalletOne) ProManager() *ProtocolManager       { return s.protocolManager }

func (s *PalletOne) MockContractLocalSend(event jury.ContractEvent) {
	s.protocolManager.ContractReqLocalSend(event)
}
func (s *PalletOne) ContractBroadcast(event jury.ContractEvent, local bool) {
	s.protocolManager.ContractBroadcast(event, local)
}
func (s *PalletOne) ElectionBroadcast(event jury.ElectionEvent) {
	s.protocolManager.ElectionBroadcast(event)
}
func (s *PalletOne) AdapterBroadcast(event jury.AdapterEvent) {
	s.protocolManager.AdapterBroadcast(event)
}

func (s *PalletOne) GetLocalMediators() []common.Address {
	return s.mediatorPlugin.LocalMediators()
}
func (s *PalletOne) IsLocalActiveMediator(addr common.Address) bool {
	return s.mediatorPlugin.IsLocalActiveMediator(addr)
}

// Protocols implements node.Service, returning all the currently configured
// network protocols to start.
func (s *PalletOne) Protocols() []p2p.Protocol {
	if s.lesServer == nil {
		return s.protocolManager.SubProtocols
	}
	return append(s.protocolManager.SubProtocols, s.lesServer.Protocols()...)
}

// Start implements node.Service, starting all internal goroutines needed by the
// PalletOne protocol implementation.
func (s *PalletOne) Start(srvr *p2p.Server) error {
	// Start the bloom bits servicing goroutines
	s.startBloomHandlers()

	// Start the RPC service
	s.netRPCService = ptnapi.NewPublicNetAPI(srvr, s.NetVersion())

	// Figure out a max peers count based on the server limits
	maxPeers := srvr.MaxPeers

	// Start the networking layer and the light server if requested
	s.protocolManager.Start(srvr, maxPeers)

	s.contractPorcessor.Start(srvr)

	s.mediatorPlugin.Start(srvr)
	return nil
}

// Stop implements node.Service, terminating all internal goroutines used by the
// PalletOne protocol.
func (s *PalletOne) Stop() error {
	//s.bloomIndexer.Close()
	s.protocolManager.Stop()
	s.txPool.Stop()
	s.engine.Stop()
	s.eventMux.Stop()
	close(s.shutdownChan)

	s.contract.Close()

	// append by Albert·Gou
	s.mediatorPlugin.Stop()

	s.dag.Close()

	return nil
}

// @author Albert·Gou
func (p *PalletOne) GetKeyStore() *keystore.KeyStore {
	return p.AccountManager().Backends(keystore.KeyStoreType)[0].(*keystore.KeyStore)
}

func (p *PalletOne) SignGenericTransaction(from common.Address, tx *modules.Transaction) (*modules.Transaction, error) {
	inputpoints := make(map[modules.OutPoint][]byte)

	for i := 0; i < len(tx.TxMessages); i++ {
		// 1. 获取PaymentPayload
		msg := tx.TxMessages[i]
		if msg.App != modules.APP_PAYMENT {
			continue
		}

		//
		payload, ok := msg.Payload.(*modules.PaymentPayload)
		if !ok {
			log.Debug("PaymentPayload conversion error, does not match TxMessage'APP type!")
		}

		// 2. 查询每个 Input 的 PkScript
		for _, txin := range payload.Inputs {
			inpoint := txin.PreviousOutPoint
			utxo, err := p.dag.GetUtxoEntry(inpoint)
			if err != nil {
				return nil, err
			}

			inputpoints[*inpoint] = utxo.PkScript
		}
	}

	// 3. 使用tokenengine 和 KeyStore 给 tx 签名
	ks := p.GetKeyStore()
	_, err := tokenengine.SignTxAllPaymentInput(tx, tokenengine.SigHashAll, inputpoints, nil,
		ks.GetPublicKey, ks.SignHash)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// @author Albert·Gou
func (p *PalletOne) SignAndSendTransaction(addr common.Address, tx *modules.Transaction) error {
	// 3. 签名 tx
	tx, err := p.SignGenericTransaction(addr, tx)
	if err != nil {
		return err
	}

	// 4. 将 tx 放入 pool
	txPool := p.TxPool()
	err = txPool.AddLocal(txspool.TxtoTxpoolTx(txPool, tx))
	if err != nil {
		return err
	}
	return nil
}

// @author Albert·Gou
func (p *PalletOne) TransferPtn(from, to string, amount decimal.Decimal, text *string) (*mp.TxExecuteResult, error) {
	// 参数检查
	if from == to {
		return nil, fmt.Errorf("please don't transfer ptn to yourself: %v", from)
	}

	if amount.Cmp(decimal.New(0, 0)) != 1 {
		return nil, fmt.Errorf("the amount of the transfer must be greater than 0")
	}

	fromAdd, err := common.StringToAddress(from)
	if err != nil {
		return nil, fmt.Errorf("invalid account address: %v", from)
	}

	toAdd, err := common.StringToAddress(to)
	if err != nil {
		return nil, fmt.Errorf("invalid account address: %v", to)
	}

	// 判断本节点是否同步完成，数据是否最新
	//if !p.dag.IsSynced() {
	//	return nil, fmt.Errorf("the data of this node is not synced, and can't transfer now")
	//}

	// 1. 创建交易
	tx, fee, err := p.dag.GenTransferPtnTx(fromAdd, toAdd, ptnjson.Ptn2Dao(amount), text, p.txPool)
	if err != nil {
		return nil, err
	}

	// 2. 签名和发送交易
	err = p.SignAndSendTransaction(fromAdd, tx)
	if err != nil {
		return nil, err
	}

	// 3. 返回执行结果
	textStr := ""
	if text != nil {
		textStr = *text
	}

	res := &mp.TxExecuteResult{}
	res.TxContent = fmt.Sprintf("Account %s transfer %vPTN to account %s with message: '%s'",
		from, amount, to, textStr)
	res.TxHash = tx.Hash()
	res.TxSize = tx.Size().TerminalString()
	res.TxFee = fmt.Sprintf("%vdao", fee)
	res.Warning = mp.DefaultResult

	return res, nil
}
