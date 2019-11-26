/*
 *
 *    This file is part of go-palletone.
 *    go-palletone is free software: you can redistribute it and/or modify
 *    it under the terms of the GNU General Public License as published by
 *    the Free Software Foundation, either version 3 of the License, or
 *    (at your option) any later version.
 *    go-palletone is distributed in the hope that it will be useful,
 *    but WITHOUT ANY WARRANTY; without even the implied warranty of
 *    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *    GNU General Public License for more details.
 *    You should have received a copy of the GNU General Public License
 *    along with go-palletone.  If not, see <http://www.gnu.org/licenses/>.
 * /
 *
 *  * @author PalletOne core developer <dev@pallet.one>
 *  * @date 2018-2019
 *
 */

package ptnjson

import (
	"encoding/hex"
	"github.com/palletone/go-palletone/common"
	crypto2 "github.com/palletone/go-palletone/common/crypto"
	"github.com/palletone/go-palletone/dag/modules"
	"time"
)

type ContractJson struct {
	//Contract Id 即Address，[20]byte，
	// 根据用户创建合约实例申请的RequestId截取其后20字节生成
	ContractId      string                `json:"contract_id"` //Hex
	ContractAddress string                `json:"contract_address"`
	TemplateId      string                `json:"tpl_id"`
	Name            string                `json:"contract_name"`
	Status          byte                  `json:"status"` // 合约状态
	Creator         string                `json:"creator"`
	CreationTime    time.Time             `json:"creation_time"` // creation date
	DuringTime      time.Time                `json:"during_time"`   // deploy during date
	Template        *ContractTemplateJson `json:"template"`
	AddrPubKey       []string `json:"jury_address"`
	Version string `json:"version"`
}

func ConvertContract2Json(contract *modules.Contract) *ContractJson {
	addr := common.NewAddress(contract.ContractId, common.ContractHash)
	creatorAddr := common.NewAddress(contract.Creator, common.PublicKeyHash)

	c := &ContractJson{
		ContractId:      hex.EncodeToString(contract.ContractId),
		ContractAddress: addr.String(),
		TemplateId:      hex.EncodeToString(contract.TemplateId),
		Name:            contract.Name,
		Status:          contract.Status,
		Creator:         creatorAddr.String(),
		CreationTime:    time.Unix(int64(contract.CreationTime), 0).UTC(),
		DuringTime:      time.Unix(int64(contract.DuringTime), 0).UTC(),
		Version:contract.Version,
	}
	for _,a := range contract.JuryPubkeys {
		c.AddrPubKey = append(c.AddrPubKey,crypto2.PubkeyBytesToAddress([]byte(a)).String())
	}
	return c
}

type ContractTemplateJson struct {
	TplId          string   `json:"tpl_id"`
	TplName        string   `json:"tpl_name"`
	TplDescription string   `json:"tpl_description"`
	Path           string   `json:"install_path"`
	Version        string   `json:"tpl_version"`
	Abi            string   `json:"abi"`
	Language       string   `json:"language"`
	AddrHash       []string `json:"addr_hash" rlp:"nil"`
	Size           uint16   `json:"size"`
	Creator        string   `json:"creator"`
	CreateTime time.Time `json:"create_time"`
}

func ConvertContractTemplate2Json(tpl *modules.ContractTemplate) *ContractTemplateJson {

	json := &ContractTemplateJson{
		TplId:          hex.EncodeToString(tpl.TplId),
		TplName:        tpl.TplName,
		TplDescription: tpl.TplDescription,
		Path:           tpl.Path,
		Version:        tpl.Version,
		Abi:            tpl.Abi,
		Language:       tpl.Language,
		Size:           tpl.Size,
		AddrHash:       []string{},
		Creator:        tpl.Creator,
		CreateTime:time.Unix(int64(tpl.CreateTime),0).UTC(),
	}
	for _, addH := range tpl.AddrHash {
		json.AddrHash = append(json.AddrHash, addH.String())
	}
	return json
}

const Deposit_ABI = `[{"constant":true,"inputs":[{"name":"address","type":"string"}],"name":"getMediatorDeposit","outputs":[{"components":[{"components":[{"name":"ApplyEnterTime","type":"string"},{"name":"ApplyQuitTime","type":"string"},{"name":"Status","type":"string"},{"name":"AgreeTime","type":"string"}],"name":"MediatorDepositExtra","type":"tuple"},{"components":[{"name":"Balance","type":"Decimal"},{"name":"EnterTime","type":"string"},{"name":"Role","type":"string"}],"name":"DepositBalanceJson","type":"tuple"}],"name":"","type":"tuple"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"address","type":"string"}],"name":"getJuryDeposit","outputs":[{"components":[{"components":[{"name":"Balance","type":"Decimal"},{"name":"EnterTime","type":"string"},{"name":"Role","type":"string"}],"name":"DepositBalanceJson","type":"tuple"},{"components":[{"name":"PublicKey","type":"string"}],"name":"JurorDepositExtraJson","type":"tuple"},{"name":"Address","type":"string"}],"name":"","type":"tuple"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"address","type":"string"}],"name":"getNodeBalance","outputs":[{"components":[{"name":"Balance","type":"Decimal"},{"name":"EnterTime","type":"string"},{"name":"Role","type":"string"}],"name":"","type":"tuple"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"address","type":"string"}],"name":"isInDeveloperList","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"getListForDeveloper","outputs":[{"name":"","type":"map[string]bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"address","type":"string"}],"name":"isInJuryCandidateList","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"getListForJuryCandidate","outputs":[{"name":"","type":"map[string]bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"address","type":"string"}],"name":"isInMediatorCandidateList","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"getListForMediatorCandidate","outputs":[{"name":"","type":"map[string]bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"address","type":"string"}],"name":"isInForfeitureList","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"getListForForfeitureApplication","outputs":[],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"address","type":"string"}],"name":"isInQuitList","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"getQuitApplyList","outputs":[],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"address","type":"string"}],"name":"isInAgreeList","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"getAgreeForBecomeMediatorList","outputs":[{"name":"","type":"map[string]bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"address","type":"string"}],"name":"isInBecomeList","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"getBecomeMediatorApplyList","outputs":[{"name":"","type":"map[string]bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"mediatorCreateArgs","type":"string"}],"name":"applyBecomeMediator","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[],"name":"mediatorPayToDepositContract","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[],"name":"mediatorApplyQuit","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"mediatorUpdateArgs","type":"string"}],"name":"updateMediatorInfo","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"pubkey","type":"string"}],"name":"juryPayToDepositContract","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[],"name":"juryApplyQuit","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[],"name":"developerPayToDepositContract","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[],"name":"devApplyQuit","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"address","type":"string"},{"name":"okOrNo","type":"string"}],"name":"handleForApplyBecomeMediator","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"address","type":"string"},{"name":"okOrNo","type":"string"}],"name":"handleForApplyQuitMediator","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"address","type":"string"},{"name":"okOrNo","type":"string"}],"name":"handleForApplyQuitJury","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"address","type":"string"},{"name":"okOrNo","type":"string"}],"name":"handleForApplyQuitDev","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"address","type":"string"},{"name":"okOrNo","type":"string"}],"name":"handleForForfeitureApplication","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"address","type":"string"}],"name":"handleNodeRemoveFromAgreeList","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"forfeitureAddress","type":"string"},{"name":"role","type":"string"},{"name":"reason","type":"string"}],"name":"applyForForfeitureDeposit","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[],"name":"processPledgeDeposit","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"amount","type":"string"}],"name":"processPledgeWithdraw","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[],"name":"handlePledgeReward","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"address","type":"string"}],"name":"queryPledgeStatusByAddr","outputs":[{"components":[{"name":"NewDepositAmount","type":"Decimal"},{"name":"PledgeAmount","type":"Decimal"},{"name":"WithdrawApplyAmount","type":"string"},{"name":"OtherAmount","type":"Decimal"}],"name":"","type":"tuple"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"queryAllPledgeHistory","outputs":[{"components":[{"name":"TotalAmount","type":"uint64"},{"name":"Date","type":"string"},{"components":[{"name":"Address","type":"string"},{"name":"Amount","type":"uint64"},{"name":"Reward","type":"uint64"}],"name":"Members","type":"tuple[]"}],"name":"","type":"tuple[]"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"queryPledgeList","outputs":[{"components":[{"name":"TotalAmount","type":"uint64"},{"name":"Date","type":"string"},{"components":[{"name":"Address","type":"string"},{"name":"Amount","type":"uint64"},{"name":"Reward","type":"uint64"}],"name":"Members","type":"tuple[]"}],"name":"","type":"tuple"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"date","type":"string"}],"name":"queryPledgeListByDate","outputs":[{"components":[{"name":"TotalAmount","type":"uint64"},{"name":"Date","type":"string"},{"components":[{"name":"Address","type":"string"},{"name":"Amount","type":"uint64"},{"name":"Reward","type":"uint64"}],"name":"Members","type":"tuple[]"}],"name":"","type":"tuple"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"queryPledgeWithdraw","outputs":[{"components":[{"name":"Address","type":"string"},{"name":"Amount","type":"uint64"}],"name":"","type":"tuple[]"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"addresses","type":"string[]"}],"name":"handleMediatorInCandidateList","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"addresses","type":"string[]"}],"name":"handleJuryInCandidateList","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"addresses","type":"string[]"}],"name":"handleDevInList","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"getAllMediator","outputs":[],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"getAllNode","outputs":[],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"getAllJury","outputs":[],"payable":false,"stateMutability":"view","type":"function"}]`
const PRC20_ABI = `[{"constant":false,"inputs":[{"name":"name","type":"string"},{"name":"symbol","type":"string"},{"name":"decimals","type":"int"},{"name":"totalSupply","type":"uint64"},{"name":"supplyAddress","type":"string"}],"name":"createToken","outputs":[{"name":"","type":"byte[]"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"symbol","type":"string"},{"name":"supplyDecimal","type":"Decimal"}],"name":"supplyToken","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"symbol","type":"string"},{"name":"newSupplyAddr","type":"string"}],"name":"changeSupplyAddr","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"symbol","type":"string"}],"name":"frozenToken","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"symbol","type":"string"}],"name":"getTokenInfo","outputs":[{"components":[{"name":"Symbol","type":"string"},{"name":"CreateAddr","type":"string"},{"name":"TotalSupply","type":"uint64"},{"name":"Decimals","type":"uint64"},{"name":"SupplyAddr","type":"string"},{"name":"AssetID","type":"string"}],"name":"","type":"tuple"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"getAllTokenInfo","outputs":[{"components":[{"name":"Symbol","type":"string"},{"name":"CreateAddr","type":"string"},{"name":"TotalSupply","type":"uint64"},{"name":"Decimals","type":"uint64"},{"name":"SupplyAddr","type":"string"},{"name":"AssetID","type":"string"}],"name":"","type":"tuple[]"}],"payable":false,"stateMutability":"view","type":"function"}]`
const Vote_ABI = `[{"constant":false,"inputs":[{"name":"name","type":"string"},{"name":"voteType","type":"string"},{"name":"totalSupply","type":"uint64"},{"name":"voteEndTime","type":"string"},{"name":"voteContentJSON","type":"string"}],"name":"createToken","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"supportRequestJSON","type":"string"}],"name":"support","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"assetID","type":"string"}],"name":"getVoteResult","outputs":[{"components":[{"name":"IsVoteEnd","type":"bool"},{"name":"CreateAddr","type":"string"},{"name":"TotalSupply","type":"uint64"},{"components":[{"name":"TopicIndex","type":"uint64"},{"name":"TopicTitle","type":"string"},{"components":[{"name":"SelectOption","type":"string"},{"name":"Num","type":"uint64"}],"name":"VoteResults","type":"tuple[]"}],"name":"SupportResults","type":"tuple[]"},{"name":"AssetID","type":"string"}],"name":"","type":"tuple"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"assetID","type":"string"}],"name":"getVoteInfo","outputs":[{"components":[{"name":"Name","type":"string"},{"name":"CreateAddr","type":"string"},{"name":"VoteType","type":"byte"},{"name":"TotalSupply","type":"uint64"},{"name":"VoteEndTime","type":"string"},{"components":[{"name":"TopicIndex","type":"uint64"},{"name":"TopicTitle","type":"string"},{"name":"SelectOptions","type":"string[]"},{"name":"SelectMax","type":"uint64"}],"name":"VoteTopics","type":"tuple[]"},{"name":"AssetID","type":"string"}],"name":"","type":"tuple"}],"payable":false,"stateMutability":"view","type":"function"}]`
const SysConfig_ABI = `[{"constant":true,"inputs":[],"name":"getWithoutVoteResult","outputs":[{"name":"","type":"byte[]"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"getVotesResult","outputs":[{"components":[{"name":"CreateAddr","type":"string"},{"name":"TotalSupply","type":"uint64"},{"name":"LeastNum","type":"uint64"},{"name":"AssetID","type":"string"},{"name":"CreateTime","type":"int64"},{"name":"IsVoteEnd","type":"bool"},{"components":[{"name":"TopicIndex","type":"uint64"},{"name":"TopicTitle","type":"string"},{"components":[{"name":"SelectOption","type":"string"},{"name":"Num","type":"uint64"}],"name":"VoteResults","type":"tuple[]"}],"name":"SupportResults","type":"tuple[]"}],"name":"","type":"tuple"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"name","type":"string"},{"name":"totalSupply","type":"uint64"},{"name":"leastNum","type":"uint64"},{"name":"voteEndTime","type":"string"},{"name":"voteContentJSON","type":"string"}],"name":"createVotesTokens","outputs":[{"name":"","type":"byte[]"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"supportRequestJson","type":"string"}],"name":"nodesVote","outputs":[{"name":"","type":"byte[]"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"field","type":"string"},{"name":"value","type":"string"}],"name":"updateSysParamWithoutVote","outputs":[{"name":"","type":"byte[]"}],"payable":false,"stateMutability":"nonpayable","type":"function"}]`
const CoinBaseABI = `[{"constant":true,"inputs":[],"name":"queryGenerateUnitReward","outputs":[{"components":[{"name":"Address","type":"string"},{"name":"Amount","type":"Decimal"},{"name":"Token","type":"Asset"}],"name":"","type":"tuple[]"}],"payable":false,"stateMutability":"view","type":"function"}]`
const BlackList_ABI = `[{"constant":false,"inputs":[{"name":"blackAddr","type":"Address"},{"name":"reason","type":"string"}],"name":"addBlacklist","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"getBlacklistRecords","outputs":[{"components":[{"name":"Address","type":"Address"},{"name":"Reason","type":"string"},{"name":"FreezeToken","type":"string"}],"name":"","type":"tuple[]"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"getBlacklistAddress","outputs":[{"name":"","type":"[]Address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"addr","type":"Address"},{"name":"amount","type":"Decimal"},{"name":"asset","type":"Asset"}],"name":"payout","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"addr","type":"Address"}],"name":"queryIsInBlacklist","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"}]`
const PRC721_ABI = `[{"constant":false,"inputs":[{"name":"name","type":"string"},{"name":"symbol","type":"string"},{"name":"UIDType","type":"string"},{"name":"totalSupply","type":"uint64"},{"name":"tokenIDMetas","type":"string"},{"name":"supplyAddress","type":"string"}],"name":"createToken","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"symbol","type":"string"},{"name":"supplyAmount","type":"uint64"},{"name":"tokenIDMetas","type":"string"}],"name":"supplyToken","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"symbol","type":"string"},{"name":"supplyAddress","type":"string"}],"name":"changeSupplyAddr","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"assetTokenID","type":"string"}],"name":"existTokenID","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"assetTokenID","type":"string"},{"name":"tokenURI","type":"string"}],"name":"setTokenURI","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"assetTokenID","type":"string"}],"name":"getTokenURI","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"symbol","type":"string"}],"name":"getOneTokenInfo","outputs":[{"components":[{"name":"Symbol","type":"string"},{"name":"CreateAddr","type":"string"},{"name":"TokenType","type":"uint8"},{"name":"TotalSupply","type":"uint64"},{"name":"SupplyAddr","type":"string"},{"name":"AssetID","type":"string"},{"name":"TokenIDs","type":"string[]"}],"name":"","type":"tuple"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"getAllTokenInfo","outputs":[{"components":[{"name":"Symbol","type":"string"},{"name":"CreateAddr","type":"string"},{"name":"TokenType","type":"uint8"},{"name":"TotalSupply","type":"uint64"},{"name":"SupplyAddr","type":"string"},{"name":"AssetID","type":"string"},{"name":"TokenIDs","type":"string[]"}],"name":"","type":"tuple[]"}],"payable":false,"stateMutability":"view","type":"function"}]`
const DigitalID_ABI = `[{"constant":false,"inputs":[{"name":"certHolder","type":"string"},{"name":"certStr","type":"string"}],"name":"addServerCert","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"certHolder","type":"string"},{"name":"certStr","type":"string"}],"name":"addMemberCert","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"certHolder","type":"string"},{"name":"certStr","type":"string"},{"name":"isServer","type":"bool"}],"name":"addCert","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"certIDOriginal","type":"string"}],"name":"addCRLCert","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"holderAddr","type":"string"}],"name":"getAddressCertIDs","outputs":[],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"issuerAddr","type":"string"}],"name":"getIssuerCertsInfo","outputs":[{"components":[{"name":"Holder","type":"string"},{"name":"IsServer","type":"bool"},{"name":"CertID","type":"string"}],"name":"","type":"tuple[]"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"certID","type":"string"}],"name":"getCertFormateInfo","outputs":[],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"certID","type":"string"}],"name":"getCertBytes","outputs":[],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"certID","type":"string"}],"name":"getCertHolder","outputs":[],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"getRootCAHolder","outputs":[],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"issuerAddr","type":"string"}],"name":"getIssuerCRL","outputs":[],"payable":false,"stateMutability":"view","type":"function"}]`
const Partition_ABI = `[{"constant":false,"inputs":[{"name":"genesisHeaderRlp","type":"string"},{"name":"forkUnitHash","type":"string"},{"name":"forkUnitHeight","type":"string"},{"name":"gasToken","type":"string"},{"name":"status","type":"string"},{"name":"syncModel","type":"string"},{"name":"networkId","type":"string"},{"name":"version","type":"string"},{"name":"stableThreshold","type":"string"},{"name":"crossChainToken","type":"string"},{"name":"peers","type":"string[]"}],"name":"registerPartition","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[],"name":"listPartition","outputs":[{"components":[{"name":"GenesisHeaderRlp","type":"byte[]"},{"name":"ForkUnitHash","type":"Hash"},{"name":"ForkUnitHeight","type":"uint64"},{"name":"GasToken","type":"AssetId"},{"name":"Status","type":"byte"},{"name":"SyncModel","type":"byte"},{"name":"NetworkId","type":"uint64"},{"name":"Version","type":"uint64"},{"name":"StableThreshold","type":"uint32"},{"name":"Peers","type":"string[]"},{"name":"CrossChainTokens","type":"[]AssetId"}],"name":"","type":"tuple[]"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"genesisHeaderRlp","type":"string"},{"name":"forkUnitHash","type":"string"},{"name":"forkUnitHeight","type":"string"},{"name":"gasToken","type":"string"},{"name":"status","type":"string"},{"name":"syncModel","type":"string"},{"name":"networkId","type":"string"},{"name":"version","type":"string"},{"name":"stableThreshold","type":"string"},{"name":"crossChainToken","type":"string"},{"name":"peers","type":"string[]"}],"name":"updatePartition","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"genesisHeaderHex","type":"string"},{"name":"gasToken","type":"string"},{"name":"status","type":"string"},{"name":"syncModel","type":"string"},{"name":"networkId","type":"string"},{"name":"version","type":"string"},{"name":"stableThreshold","type":"string"},{"name":"crossChainToken","type":"string"},{"name":"peers","type":"string[]"}],"name":"setMainChain","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"getMainChain","outputs":[{"components":[{"name":"GenesisHeaderRlp","type":"byte[]"},{"name":"Status","type":"byte"},{"name":"SyncModel","type":"byte"},{"name":"GasToken","type":"AssetId"},{"name":"NetworkId","type":"uint64"},{"name":"Version","type":"uint64"},{"name":"StableThreshold","type":"uint32"},{"name":"Peers","type":"string[]"},{"name":"CrossChainTokens","type":"[]AssetId"}],"name":"","type":"tuple"}],"payable":false,"stateMutability":"view","type":"function"}]`
const Debug_ABI = `[{"constant":false,"inputs":[],"name":"error","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"a","type":"int"},{"name":"b","type":"int"}],"name":"add","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"addr","type":"string"}],"name":"getbalance","outputs":[],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"getRequesterCert","outputs":[],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[],"name":"checkRequesterCert","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"getRootCABytes","outputs":[],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"account","type":"string"},{"name":"amount","type":"string"}],"name":"addBalance","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"account","type":"string"}],"name":"getBalance","outputs":[],"payable":false,"stateMutability":"view","type":"function"}]`

var (
	PalletOneABI = map[string]string{
		"PCGTta3M4t3yXu8uRgkKvaWd2d8DR32W9vM": Deposit_ABI,
		"PCGTta3M4t3yXu8uRgkKvaWd2d8DREThG43": PRC20_ABI,
		"PCGTta3M4t3yXu8uRgkKvaWd2d8DRLGbeyd": Vote_ABI,
		"PCGTta3M4t3yXu8uRgkKvaWd2d8DRS71ZEM": SysConfig_ABI,
		"PCGTta3M4t3yXu8uRgkKvaWd2d8DRUp5qmM": CoinBaseABI,
		"PCGTta3M4t3yXu8uRgkKvaWd2d8DRdWEXJF": BlackList_ABI,
		"PCGTta3M4t3yXu8uRgkKvaWd2d8DRijspoq": PRC721_ABI,
		"PCGTta3M4t3yXu8uRgkKvaWd2d8DRv2vsEk": DigitalID_ABI,
		"PCGTta3M4t3yXu8uRgkKvaWd2d8DRxVdGDZ": Partition_ABI,
		"PCGTta3M4t3yXu8uRgkKvaWd2d8DSfQdUHf": Debug_ABI,
	}
)

func GetSysContractABI(addr string) *ContractTemplateJson {
	addrABI, exist := PalletOneABI[addr]
	if !exist {
		return nil
	}
	json := &ContractTemplateJson{
		//TplId:          "",
		//TplName:        "PRC20",
		//TplDescription: "Fungible Token",
		//Path:           "",
		//Version:        "v1.0.0",
		Abi: addrABI,
		//Language:       "Golang",
		//Size:           0,
		//AddrHash:       []string{},
		//Creator:        "",
	}
	return json
}
