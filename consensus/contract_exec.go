package consensus

import (
	"github.com/palletone/go-palletone/common"
	"github.com/palletone/go-palletone/common/event"
	"github.com/palletone/go-palletone/consensus/jury"
	"github.com/palletone/go-palletone/core/accounts/keystore"
	"github.com/palletone/go-palletone/dag/modules"
	"github.com/palletone/go-palletone/dag/txspool"
	"time"
)

type ContractInf interface {
	SubscribeContractEvent(ch chan<- jury.ContractEvent) event.Subscription
	ProcessContractEvent(event *jury.ContractEvent) error
	ProcessElectionEvent(event *jury.ElectionEvent) (result *jury.ElectionEvent,err error)
	ProcessAdapterEvent(event *jury.AdapterEvent) (result *jury.AdapterEvent,err error)

	AdapterFunRequest(reqId common.Hash, contractId common.Address, timeOut time.Duration) (interface{}, error)
	AddContractLoop(txpool txspool.ITxPool, addr common.Address, ks *keystore.KeyStore) error
	CheckContractTxValid(tx *modules.Transaction, execute bool) bool
	IsSystemContractTx(tx *modules.Transaction) bool
}