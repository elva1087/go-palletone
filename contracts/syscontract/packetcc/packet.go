/*
 *
 * 	This file is part of go-palletone.
 * 	go-palletone is free software: you can redistribute it and/or modify
 * 	it under the terms of the GNU General Public License as published by
 * 	the Free Software Foundation, either version 3 of the License, or
 * 	(at your option) any later version.
 * 	go-palletone is distributed in the hope that it will be useful,
 * 	but WITHOUT ANY WARRANTY; without even the implied warranty of
 * 	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * 	GNU General Public License for more details.
 * 	You should have received a copy of the GNU General Public License
 * 	along with go-palletone.  If not, see <http://www.gnu.org/licenses/>.
 * /
 *
 *
 *  * @author PalletOne core developer  <dev@pallet.one>
 *  * @date 2018-2020
 *
 */

package packetcc

import (
	"encoding/hex"
	"math/rand"
	"sort"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/palletone/go-palletone/common"
	"github.com/palletone/go-palletone/common/crypto"
	"github.com/palletone/go-palletone/contracts/shim"
	"github.com/palletone/go-palletone/dag/modules"
	"github.com/shopspring/decimal"
)

// old packet
type OldPacket struct {
	PubKey          []byte         //红包对应的公钥，也是红包的唯一标识
	Creator         common.Address //红包发放人员地址
	Token           *modules.Asset //红包中的TokenID
	Amount          uint64         //红包总金额
	Count           uint32         //红包数，为0表示可以无限领取
	MinPacketAmount uint64         //单个红包最小额
	MaxPacketAmount uint64         //单个红包最大额,最大额最小额相同，则说明不是随机红包,0则表示完全随机
	ExpiredTime     uint64         //红包过期时间，0表示永不过期
	Remark          string         //红包的备注
	Constant        bool           //是否固定数额
}

// new Packet
type Packet struct {
	PubKey          []byte         //红包对应的公钥，也是红包的唯一标识
	Creator         common.Address //红包发放人员地址
	Tokens          []*Tokens      //红包中的TokenID
	Amount          uint64         //红包总金额
	Count           uint32         //红包数，为0表示可以无限领取
	MinPacketAmount uint64         //单个红包最小额
	MaxPacketAmount uint64         //单个红包最大额,最大额最小额相同，则说明不是随机红包,0则表示完全随机
	ExpiredTime     uint64         //红包过期时间，0表示永不过期
	Remark          string         //红包的备注
	Constant        bool           //是否固定数额
}

type Tokens struct {
	Amount        uint64         `json:"amount"` //数量
	Asset         *modules.Asset `json:"asset"`  //资产
	BalanceAmount uint64         //红包剩余额度
	BalanceCount  uint32         //红包剩余次数
}

type TokensJson struct {
	Amount        decimal.Decimal `json:"amount"` //数量
	Asset         string          `json:"asset"`  //资产
	BalanceAmount decimal.Decimal //红包剩余额度
	BalanceCount  uint32          //红包剩余次数
}

type RecordTokensJson struct {
	Amount decimal.Decimal `json:"amount"` //数量
	Asset  string          `json:"asset"`  //资产
}

type RecordTokens struct {
	Amount uint64         `json:"amount"` //数量
	Asset  *modules.Asset `json:"asset"`  //资产
}

type PacketJson struct {
	PubKey          string          //红包对应的公钥，也是红包的唯一标识
	Creator         common.Address  //红包发放人员地址
	Token           []*TokensJson   //红包中的TokenID
	TotalAmount     decimal.Decimal //红包总金额
	PacketCount     uint32          //红包数，为0表示可以无限领取
	MinPacketAmount decimal.Decimal //单个红包最小额
	MaxPacketAmount decimal.Decimal //单个红包最大额,最大额最小额相同，则说明不是随机红包
	ExpiredTime     string          //红包过期时间，0表示永不过期
	Remark          string          //红包的备注
	IsConstant      string          //是否固定数额
	BalanceAmount   decimal.Decimal //红包剩余额度
	BalanceCount    uint32          //红包剩余次数
}

//红包余额
type PacketBalance struct {
	Amount uint64
	Count  uint32
}

func (p *Packet) IsFixAmount() bool {
	return p.MinPacketAmount == p.MaxPacketAmount && p.MaxPacketAmount > 0
}

func (p *Packet) PubKeyAddress() common.Address {
	return crypto.PubkeyBytesToAddress(p.PubKey)
}

func (p *Packet) GetPullAmount(seed int64, amount uint64, count uint32) uint64 {
	if p.IsFixAmount() {
		return p.MaxPacketAmount
	}
	if count == 1 {
		if amount > p.MaxPacketAmount {
			return p.MaxPacketAmount
		}
		return amount
	}
	expect := amount / uint64(count)
	return NormalRandom(seed, expect, p.MinPacketAmount, p.MaxPacketAmount)
}

// 随机返回数额
func NormalRandom(seed int64, expect uint64, min, max uint64) uint64 {
	if expect < min {
		return min
	}
	if expect > max {
		return max
	}
	//计算标准差
	bzc1 := max - expect
	bzc2 := expect - min
	bzc := bzc1
	if bzc2 < bzc1 {
		bzc = bzc2
	}
	bzc = bzc / 3 //正态分布，3标准差内的概率>99%

	for i := int64(0); i < 100; i++ {
		rand.Seed(seed + i)
		number := rand.NormFloat64()*float64(bzc) + float64(expect)
		if number <= 0 {
			continue
		}
		n64 := uint64(number)
		if n64 >= min && n64 <= max {
			return n64
		}
	}
	return expect
}

// 保存红包
func savePacket(stub shim.ChaincodeStubInterface, p *Packet) error {
	key := PacketPrefix + hex.EncodeToString(p.PubKey)
	value, err := rlp.EncodeToBytes(p)
	if err != nil {
		return err
	}
	return stub.PutState(key, value)
}

// 获取红包
func getPacket(stub shim.ChaincodeStubInterface, pubKey []byte) (*Packet, error) {
	key := PacketPrefix + hex.EncodeToString(pubKey)
	value, err := stub.GetState(key)
	if err != nil {
		return nil, err
	}
	p := Packet{}
	err = rlp.DecodeBytes(value, &p)
	if err != nil {
		// 兼容
		op := OldPacket{}
		err = rlp.DecodeBytes(value, &op)
		if err != nil {
			return nil, err
		}
		// 转换
		balanceAmount, balanceCount, _ := getPacketBalance(stub, op.PubKey)
		np := OldPacket2New(&op, balanceAmount, balanceCount)
		p = *np
	}
	sort.Slice(p.Tokens, func(i, j int) bool {
		return p.Tokens[i].Amount > p.Tokens[j].Amount
	})
	return &p, nil
}

// 获取所有红包
func getPackets(stub shim.ChaincodeStubInterface) ([]*Packet, error) {
	value, err := stub.GetStateByPrefix(PacketPrefix)
	if err != nil {
		return nil, err
	}
	ps := []*Packet{}
	for _, pp := range value {
		p := Packet{}
		err = rlp.DecodeBytes(pp.Value, &p)
		if err != nil {
			// 兼容
			op := OldPacket{}
			err = rlp.DecodeBytes(pp.Value, &op)
			if err != nil {
				return nil, err
			}
			// 转换
			balanceAmount, balanceCount, _ := getPacketBalance(stub, op.PubKey)
			np := OldPacket2New(&op, balanceAmount, balanceCount)
			p = *np
		}
		sort.Slice(p.Tokens, func(i, j int) bool {
			return p.Tokens[i].Amount > p.Tokens[j].Amount
		})
		ps = append(ps, &p)
	}
	return ps, nil
}

// 保存红包余额和个数
func savePacketBalance(stub shim.ChaincodeStubInterface, pubKey []byte, balanceAmt uint64, balanceCount uint32) error {
	key := PacketBalancePrefix + hex.EncodeToString(pubKey)
	value, err := rlp.EncodeToBytes(PacketBalance{Amount: balanceAmt, Count: balanceCount})
	if err != nil {
		return err
	}
	return stub.PutState(key, value)
}

// 获取红包余额和个数
func getPacketBalance(stub shim.ChaincodeStubInterface, pubKey []byte) (uint64, uint32, error) {
	key := PacketBalancePrefix + hex.EncodeToString(pubKey)
	value, err := stub.GetState(key)
	if err != nil {
		return 0, 0, err
	}
	b := PacketBalance{}
	err = rlp.DecodeBytes(value, &b)
	if err != nil {
		return 0, 0, err
	}
	return b.Amount, b.Count, nil
}

func convertPacket2Json(packet *Packet, balanceAmount uint64, balanceCount uint32) *PacketJson {
	js := &PacketJson{
		PubKey:          hex.EncodeToString(packet.PubKey),
		Creator:         packet.Creator,
		TotalAmount:     packet.Tokens[0].Asset.DisplayAmount(packet.Amount),
		MinPacketAmount: packet.Tokens[0].Asset.DisplayAmount(packet.MinPacketAmount),
		MaxPacketAmount: packet.Tokens[0].Asset.DisplayAmount(packet.MaxPacketAmount),
		PacketCount:     packet.Count,
		Remark:          packet.Remark,
		IsConstant:      strconv.FormatBool(packet.Constant),
		BalanceAmount:   packet.Tokens[0].Asset.DisplayAmount(balanceAmount),
		BalanceCount:    balanceCount,
	}
	if packet.ExpiredTime != 0 {
		js.ExpiredTime = time.Unix(int64(packet.ExpiredTime), 0).String()
	}
	js.Token = make([]*TokensJson, len(packet.Tokens))
	for i, t := range packet.Tokens {
		js.Token[i] = &TokensJson{}
		js.Token[i].Amount = t.Asset.DisplayAmount(t.Amount)
		js.Token[i].Asset = t.Asset.String()
		js.Token[i].BalanceAmount = t.Asset.DisplayAmount(t.BalanceAmount)
		js.Token[i].BalanceCount = t.BalanceCount
	}
	return js
}

//  判断是否基金会发起的
func isFoundationInvoke(stub shim.ChaincodeStubInterface) bool {
	//  判断是否基金会发起的
	invokeAddr, err := stub.GetInvokeAddress()
	if err != nil {
		return false
	}
	//  获取
	gp, err := stub.GetSystemConfig()
	if err != nil {
		return false
	}
	foundationAddress := gp.ChainParameters.FoundationAddress
	// 判断当前请求的是否为基金会
	if invokeAddr.String() != foundationAddress {
		return false
	}
	return true
}

// 是否红包已被领了
func isPulledPacket(stub shim.ChaincodeStubInterface, pubKey []byte, message string) bool {
	key := PacketAllocationRecordPrefix + hex.EncodeToString(pubKey) + "-" + message
	byte, _ := stub.GetState(key)
	if byte == nil {
		return false
	}
	return true
}

func OldPacket2New(old *OldPacket, BalanceAmount uint64, BalanceCount uint32) *Packet {
	return &Packet{
		PubKey:  old.PubKey,
		Creator: old.Creator,
		Tokens: []*Tokens{
			{
				Amount:        old.Amount,
				Asset:         old.Token,
				BalanceCount:  BalanceCount,
				BalanceAmount: BalanceAmount,
			},
		},
		Amount:          old.Amount,
		Count:           old.Count,
		MinPacketAmount: old.MinPacketAmount,
		MaxPacketAmount: old.MaxPacketAmount,
		ExpiredTime:     old.ExpiredTime,
		Remark:          old.Remark,
		Constant:        old.Constant,
	}
}
