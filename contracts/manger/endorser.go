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

package manger

import (
	"time"

	"github.com/palletone/go-palletone/dag/dboperation"
	"github.com/pkg/errors"
	"golang.org/x/net/context"

	"github.com/palletone/go-palletone/common/log"
	"github.com/palletone/go-palletone/contracts/core"
	"github.com/palletone/go-palletone/contracts/shim"
	pb "github.com/palletone/go-palletone/core/vmContractPub/protos/peer"
	putils "github.com/palletone/go-palletone/core/vmContractPub/protos/utils"
	"github.com/palletone/go-palletone/dag/modules"
	"github.com/palletone/go-palletone/dag/rwset"
)

//var log = flogging.MustGetLogger("ccmanger")

// Support contains functions that the endorser requires to execute its tasks
type Support interface {
	IsSysCCAndNotInvokableExternal(name string) bool
	// GetTxSimulator returns the transaction simulator ,they are made unique
	// by way of the supplied txid
	GetTxSimulator(rwM rwset.TxManager, idag dboperation.IContractDag, txId string) (rwset.TxSimulator, error)

	IsSysCC(name string) bool

	Execute(contractid []byte, ctxt context.Context, cid, name, version, txid string, syscc bool, signedProp *pb.PtnSignedProposal, prop *pb.PtnProposal, spec interface{}, timeout time.Duration) (*pb.PtnResponse, *pb.PtnChaincodeEvent, error)
}

// Endorser provides the Endorser service ProcessProposal
type Endorser struct {
	s Support
}

// validateResult provides the result of endorseProposal verification
type validateResult struct {
	prop *pb.PtnProposal
	//hdrExt  *pb.ChaincodeHeaderExtension
	chainID string
	txid    string
	resp    *pb.PtnProposalResponse
}

// NewEndorserServer creates and returns a new Endorser server instance.
func NewEndorserServer(s Support) pb.EndorserServer {
	e := &Endorser{
		s: s,
	}
	return e
}

//call specified chaincode (system or user)
func (e *Endorser) callChaincode(contractid []byte, ctxt context.Context, chainID string, version string, txid string, signedProp *pb.PtnSignedProposal, prop *pb.PtnProposal, cis *pb.PtnChaincodeInvocationSpec, chaincodeName string, txsim rwset.TxSimulator, timeout time.Duration) (*pb.PtnResponse, *pb.PtnChaincodeEvent, error) {
	log.Debugf("call chain code enter")
	log.Debugf("[%s][%s] Entry chaincode: %s version: %s", chainID, shorttxid(txid), chaincodeName, version)
	defer log.Debugf("[%s][%s] Exit", chainID, shorttxid(txid))
	var err error
	var res *pb.PtnResponse
	var ccevent *pb.PtnChaincodeEvent

	if txsim != nil {
		ctxt = context.WithValue(ctxt, core.TXSimulatorKey, txsim)
	}

	scc := e.s.IsSysCC(chaincodeName)
	res, ccevent, err = e.s.Execute(contractid, ctxt, chainID, chaincodeName, version, txid, scc, signedProp, prop, cis, timeout)
	log.Debugf("execute")
	if err != nil {
		return res, nil, err
	}

	if res.Status >= shim.ERRORTHRESHOLD {
		return res, nil, nil
	}

	return res, ccevent, err
}

func (e *Endorser) simulateProposal(contractid []byte, ctx context.Context, chainID string, txid string, signedProp *pb.PtnSignedProposal, prop *pb.PtnProposal, cid *pb.PtnChaincodeID, txsim rwset.TxSimulator, tmout time.Duration) (*pb.PtnResponse, []byte, *pb.PtnChaincodeEvent, error) {
	log.Debugf("[%s][%s] Entry chaincode: %s", chainID, shorttxid(txid), cid)
	defer log.Debugf("[%s][%s] Exit", chainID, shorttxid(txid))

	cis, err := putils.GetChaincodeInvocationSpec(prop)
	if err != nil {
		log.Errorf("GetChaincodeInvocationSpec err:[%s][%s] Entry chaincode: %s", chainID, shorttxid(txid), cid)
		return nil, nil, nil, err
	}
	log.Infof("spec=%v", cis)

	//var cdLedger resourcesconfig.ChaincodeDefinition
	//
	//if !e.s.IsSysCC(cid.Name) {
	//	cdLedger, err = e.s.GetChaincodeDefinition(ctx, chainID, txid, signedProp, prop, cid.Name, txsim)
	//	if err != nil {
	//		return nil, nil, nil, nil, errors.WithMessage(err, fmt.Sprintf("make sure the chaincode %s has been successfully instantiated and try again", cid.Name))
	//	}
	//	version = cdLedger.CCVersion()
	//
	//	err = e.s.CheckInstantiationPolicy(cid.Name, version, cdLedger)
	//	if err != nil {
	//		return nil, nil, nil, nil, err
	//	}
	//} else {
	//	version = util.GetSysCCVersion()
	//}

	//---3. execute the proposal and get simulation results
	//var simResult *ledger.TxSimulationResults
	var simResBytes []byte
	var res *pb.PtnResponse
	var ccevent *pb.PtnChaincodeEvent
	res, ccevent, err = e.callChaincode(contractid, ctx, chainID, cid.Version, txid, signedProp, prop, cis, cid.Name, txsim, tmout)
	log.Debugf("call chain code")
	if err != nil {
		log.Errorf("[%s][%s] failed to invoke chaincode %s, error: %+v", chainID, shorttxid(txid), cid, err)
		return res, nil, nil, err
	}

	//if txsim != nil {
	//	//if simResult, err = txsim.GetTxSimulationResults(); err != nil {
	//	//	return  nil, nil, nil, err
	//	//}
	//}

	return res, simResBytes, ccevent, nil
}

//endorse the proposal
//func (e *Endorser) endorseProposal(ctx context.Context, chainID string, txid string, signedProp *pb.PtnSignedProposal, proposal *pb.PtnProposal, response *pb.Response, simRes []byte, event *pb.PtnChaincodeEvent, visibility []byte, ccid *pb.PtnChaincodeID, txsim rwset.TxSimulator) (*pb.PtnProposalResponse, error) {
//	log.Debugf("[%s][%s] Entry chaincode: %s", chainID, shorttxid(txid), ccid)
//	defer log.Debugf("[%s][%s] Exit", chainID, shorttxid(txid))
//
//	return nil, nil
//}

//preProcess checks the tx proposal headers, uniqueness and ACL
func (e *Endorser) validateProcess(signedProp *pb.PtnSignedProposal) (*validateResult, error) {
	vr := &validateResult{}

	// extract the Proposal message from signedProp
	prop, err := putils.GetProposal(signedProp.ProposalBytes)
	if err != nil {
		return nil, err
	}

	// 1) look at the ProposalHeader
	hdr, err := putils.GetHeader(prop.Header)
	if err != nil {
		return nil, err
	}

	//TODO validate the header

	//if err != nil {
	//	vr.resp = &pb.PtnProposalResponse{Response: &pb.Response{Status: 500, Message: err.Error()}}
	//	return vr, err
	//}

	chdr, err := putils.UnmarshalChannelHeader(hdr.ChannelHeader)
	if err != nil {
		vr.resp = &pb.PtnProposalResponse{Response: &pb.PtnResponse{Status: 500, Message: err.Error()}}
		return vr, err
	}
	//shdr, err := putils.GetSignatureHeader(hdr.SignatureHeader)
	//if err != nil {
	//	vr.resp = &pb.PtnProposalResponse{Response: &pb.Response{Status: 500, Message: err.Error()}}
	//	return vr, err
	//}

	vr.prop, vr.chainID, vr.txid = prop, chdr.ChannelId, chdr.TxId

	return vr, nil
}

// ProcessProposal process the Proposal
//func (e *Endorser) ProcessProposal(ctx context.Context, signedProp *pb.PtnSignedProposal) (*pb.PtnProposalResponse, error) {
func (e *Endorser) ProcessProposal(rwM rwset.TxManager, idag dboperation.IContractDag, deployId []byte, ctx context.Context,
	signedProp *pb.PtnSignedProposal, prop *pb.PtnProposal, chainID string, cid *pb.PtnChaincodeID, tmout time.Duration) (
	*pb.PtnProposalResponse, *modules.ContractInvokeResult, error) {
	log.Debugf("process proposal enter")
	var txsim rwset.TxSimulator

	//addr := util.ExtractRemoteAddress(ctx)
	//log.Debug("Entering: Got request from", addr)
	//defer log.Debugf("Exit: request from", addr)

	//0 -- check and validate
	result, err := e.validateProcess(signedProp)
	log.Debugf("validate process")
	if err != nil {
		log.Debugf("validate signedProp err:%s", err)
		return nil, nil, err
	}
	txid := result.txid
	if chainID != "" {
		if txsim, err = e.s.GetTxSimulator(rwM, idag, txid); err != nil {
			return &pb.PtnProposalResponse{Response: &pb.PtnResponse{Status: 500, Message: err.Error()}}, nil, err
		}
		//defer txsim.Done()
	}
	if err != nil {
		return &pb.PtnProposalResponse{Response: &pb.PtnResponse{Status: 500, Message: err.Error()}}, nil, err
	}

	//1 -- simulate
	res, _, _, err := e.simulateProposal(deployId, ctx, chainID, txid, signedProp, prop, cid, txsim, tmout)
	log.Debugf("simulate proposal")
	if err != nil {
		txsim.Rollback()
		return &pb.PtnProposalResponse{Response: &pb.PtnResponse{Status: 500, Message: err.Error()}}, nil, err
	}
	if res != nil {
		if res.Status >= shim.ERROR {
			log.Infof("[%s][%s] simulateProposal() resulted in chaincode, response status %d for txid %s:%s",
				chainID, shorttxid(txid), res.Status, txid, res.Message)

			resp := &pb.PtnProposalResponse{
				Payload:  nil,
				Response: &pb.PtnResponse{Status: 500, Message: res.Message}}
			txsim.Rollback()
			return resp, nil, errors.New("Chaincode Error:" + res.Message)
		}
	} else {
		log.Error("simulateProposal response is nil")
		txsim.Rollback()
		return &pb.PtnProposalResponse{
				Payload: nil, Response: &pb.PtnResponse{Status: 500, Message: "simulateProposal response is nil"}}, nil,
			errors.New("Chaincode Error:simulateProposal response is nil")
	}

	//2 -- endorse and get a marshaled ProposalResponse message
	pResp := &pb.PtnProposalResponse{Response: res}
	cis, err := putils.GetChaincodeInvocationSpec(prop)
	if err != nil {
		txsim.Rollback()
		return nil, nil, err
	}
	invokeResult, err := RwTxResult2DagInvokeUnit(txsim, txid, cis.ChaincodeSpec.ChaincodeId.Name, deployId,
		cis.ChaincodeSpec.Input.Args, tmout)
	if err != nil {
		log.Errorf("chainID[%s] converRwTxResult2DagUnit failed", chainID)
		txsim.Rollback()
		return nil, nil, errors.New("Conver RwSet to dag unit fail")
	}

	pResp.Response.Payload = res.Payload
	invokeResult.Payload = res.Payload
	txsim.Done()
	return pResp, invokeResult, nil
}
