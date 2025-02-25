package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"

	"github.com/btcsuite/btcutil/base58"
)

type TXInput struct {
	TXID []byte // 交易id

	Index int64 // output索引
	// Address string // 解锁脚本，先用地址模拟

	Signature []byte // 交易签名

	PubKey []byte // 公钥本身，不是公钥hash

}

type TXOutput struct {
	Value float64 // 转账金额
	// Address string  // 锁定脚本

	PubKeyHash []byte // 公钥hash，不是公钥本身
}

// 给定转账地址，得到这个地址的公钥hash，完成对output的锁定
func (output *TXOutput) Lock(address string) {
	// address -> public key hash
	decodeInfo := base58.Decode(address)

	pubKeyHash := decodeInfo[1 : len(decodeInfo)-4]

	output.PubKeyHash = pubKeyHash
}

func NewTXOutput(value float64, address string) TXOutput {
	output := TXOutput{Value: value}

	output.Lock(address)

	return output
}

type Transaction struct {
	TXid      []byte     // 交易id
	TXInputs  []TXInput  // 所有交易的inputs
	TXOutputs []TXOutput // 所有的outputs
}

func (tx *Transaction) SetTXID() {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	hash := sha256.Sum256(buffer.Bytes())
	tx.TXid = hash[:]
}

// 挖矿奖励
const reward = 3.125

// 实现挖矿交易
// 只有输出，没有有效输入
func NewCoinbaseTx(miner string, data string) *Transaction {
	inputs := []TXInput{{nil, -1, nil, []byte(data)}}
	// outputs := []TXOutput{{3.125, miner}}

	output := NewTXOutput(reward, miner)
	outputs := []TXOutput{output}

	tx := Transaction{nil, inputs, outputs}
	tx.SetTXID()

	return &tx
}

// 判断交易是否为挖矿交易
func (tx *Transaction) IsCoinbase() bool {
	// 1.只有一个input，2. 引用的id是nil 3. 引用的索引是-1
	inputs := tx.TXInputs
	if len(inputs) == 1 && inputs[0].TXID == nil && inputs[0].Index == -1 {
		return true
	}

	return false
}

// 创建普通交易
func NewTransaction(from, to string, amount float64, bc *BlockChain) *Transaction {
	// 1. 打开钱包
	ws := NewWallets()
	// 获取密钥对
	walllet := ws.WalletsMap[from]

	if walllet == nil {
		fmt.Printf("%s 的私钥不存在，交易创建失败！\n", from)
		return nil
	}

	// 2. 获取公钥私钥
	// privateKey := walllet.PrivateKey
	publicKey := walllet.PublicKey

	pubKeyHash := HashPubKey(publicKey)

	// 标识能用的utxo
	var utxos = make(map[string][]int64)
	// 标识这些utxo存储的金额
	var resValue float64

	// 1.遍历账本，找到属于付款人的合适的金额
	utxos, resValue = bc.FindNeedUtxos(pubKeyHash, amount)

	// 2.如果找到的钱不足以转账，则交易失败
	if resValue < amount {
		fmt.Printf("余额不足，交易失败！\n")
		return nil
	}

	var inputs []TXInput
	var outputs []TXOutput

	// 3.将outputs转成inputs
	for txid, indexes := range utxos {
		for _, index := range indexes {
			inputs = append(inputs, TXInput{[]byte(txid), index, nil, publicKey})
		}
	}

	// 4.生成outputs
	output := NewTXOutput(amount, to)
	outputs = append(outputs, output)

	// 5.如果有剩余，则转给自己
	if resValue > amount {
		// output := TXOutput{(resValue - amount), from}
		output1 := NewTXOutput(resValue-amount, from)
		outputs = append(outputs, output1)
	}

	// 创建交易
	tx := Transaction{nil, inputs, outputs}
	// 设置交易id
	tx.SetTXID()
	// 返回交易结构
	return &tx
}
