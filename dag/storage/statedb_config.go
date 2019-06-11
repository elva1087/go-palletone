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
 *  * @date 2018
 *
 */

//SysConfig来自于系统合约SysConfigContractAddress的状态数据

package storage

import (
	"encoding/json"

	"github.com/palletone/go-palletone/common/log"
	"github.com/palletone/go-palletone/contracts/syscontract"
	"github.com/palletone/go-palletone/dag/dagconfig"
	"github.com/palletone/go-palletone/dag/modules"
)

//var CONF_PREFIX = append(constants.CONTRACT_STATE_PREFIX, scc.SysConfigContractAddress.Bytes()...)
func (statedb *StateDb) SaveSysConfig(key string, val []byte, ver *modules.StateVersion) error {
	//SaveContractState(id []byte, name string, value interface{}, version *modules.StateVersion)
	id := syscontract.SysConfigContractAddress.Bytes()
	err := saveContractState(statedb.db, id, key, val, ver)
	if err != nil {
		return err
	}
	return nil
}

/**
获取配置信息
get config information
*/
func (statedb *StateDb) GetSysConfig(name string) ([]byte, *modules.StateVersion, error) {
	id := syscontract.SysConfigContractAddress.Bytes()
	return statedb.GetContractState(id, name)
}

//func (statedb *StateDb) GetAllSysConfig() (map[string]*modules.ContractStateValue, error) {
//	id := syscontract.SysConfigContractAddress.Bytes()
//	return statedb.GetContractStatesById(id)
//}

func (statedb *StateDb) GetMinFee() (*modules.AmountAsset, error) {
	assetId := dagconfig.DagConfig.GetGasToken()
	return &modules.AmountAsset{Amount: 0, Asset: assetId.ToAsset()}, nil
}
func (statedb *StateDb) GetPartitionChains() ([]*modules.PartitionChain, error) {
	id := syscontract.PartitionContractAddress.Bytes()
	rows, err := statedb.GetContractStatesByPrefix(id, "PC")
	result := []*modules.PartitionChain{}
	if err != nil {
		return result, nil
	}

	for _, v := range rows {
		partition := &modules.PartitionChain{}
		json.Unmarshal(v.Value, &partition)
		result = append(result, partition)
	}
	return result, nil
}
func (statedb *StateDb) GetMainChain() (*modules.MainChain, error) {
	id := syscontract.PartitionContractAddress.Bytes()
	data, _, err := statedb.GetContractState(id, "MainChain")
	if err != nil {
		return nil, err
	}
	mainChain := &modules.MainChain{}
	err = json.Unmarshal(data, mainChain)
	if err != nil {
		return nil, err
	}
	return mainChain, nil
}

func (statedb *StateDb) GetSysParamWithoutVote() (map[string]string, error) {
	var res map[string]string

	val, _, err := statedb.GetSysConfig(modules.DesiredSysParamsWithoutVote)
	if err != nil {
		log.Debugf(err.Error())
		return nil, err
	}

	err = json.Unmarshal(val, &res)
	if err != nil {
		log.Debugf(err.Error())
		return nil, err
	}

	return res, nil
}

func (statedb *StateDb) GetSysParamsWithVotes() (*modules.SysTokenIDInfo, error) {
	val, _, err := statedb.GetSysConfig(modules.DesiredSysParamsWithVote)
	if err != nil {
		return nil, err
	}
	info := &modules.SysTokenIDInfo{}
	if val == nil {
		return nil, err
	} else if len(val) > 0 {
		err := json.Unmarshal(val, info)
		if err != nil {
			return nil, err
		}
		return info, nil
	} else {
		return nil, nil
	}
}

func (statedb *StateDb) UpdateSysParams(version *modules.StateVersion) error {
	//基金会单独修改的
	var err error
	modifies, err := statedb.GetSysParamWithoutVote()
	if err != nil {
		return err
	}
	//基金会发起投票的
	info, err := statedb.GetSysParamsWithVotes()
	if err != nil {
		return err
	}
	if modifies == nil && info == nil {
		return nil
	}
	//获取当前的version
	if len(modifies) > 0 {
		for k, v := range modifies {
			err = statedb.SaveSysConfig(k, []byte(v), version)
			if err != nil {
				return err
			}
		}
		//将基金会当前单独修改的重置为nil
		err = statedb.SaveSysConfig(modules.DesiredSysParamsWithoutVote, nil, version)
		if err != nil {
			return err
		}
	}
	if info == nil {
		return nil
	}
	//foundAddr, _, err := statedb.GetSysConfig(modules.FoundationAddress)
	//if err != nil {
	//	return err
	//}
	//if info.CreateAddr != string(foundAddr) {
	//	return fmt.Errorf("only foundation can call this function")
	//}
	if !info.IsVoteEnd {
		return nil
	}
	for _, v1 := range info.SupportResults {
		for _, v2 := range v1.VoteResults {
			//TODO
			if v2.Num >= info.LeastNum {
				err = statedb.SaveSysConfig(v1.TopicTitle, []byte(v2.SelectOption), version)
				if err != nil {
					return err
				}
				break
			}
		}
	}
	//将基金会当前投票修改的重置为nil
	err = statedb.SaveSysConfig(modules.DesiredSysParamsWithVote, nil, version)
	if err != nil {
		return err
	}
	return nil
}
