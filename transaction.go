package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
)

type TXInput struct {
	TXID    []byte // 交易id
	Index   int64  // output索引
	Address string // 解锁脚本，先用地址模拟
}

type TXOutput struct {
	Value   float64 // 转账金额
	Address string  // 锁定脚本
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

// 实现挖矿交易
// 只有输出，没有有效输入
func NewCoinbaseTx(miner string, data string) *Transaction {
	inputs := []TXInput{{nil, -1, data}}
	outputs := []TXOutput{{3.125, miner}}

	tx := Transaction{nil, inputs, outputs}
	tx.SetTXID()

	return &tx
}

// 创建普通交易
func NewTransaction(from, to string, amount float64, bc *BlockChain) *Transaction {
	// 标识能用的utxo
	var utxos = make(map[string][]int64)
	// 标识这些utxo存储的金额
	var resValue float64

	// 1.遍历账本，找到属于付款人的合适的金额
	utxos, resValue = bc.FindNeedUtxos(from, amount)

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
			inputs = append(inputs, TXInput{[]byte(txid), index, from})
		}
	}

	// 4.生成outputs
	outputs = append(outputs, TXOutput{amount, to})

	// 5.如果有剩余，则转给自己
	if resValue > amount {
		output := TXOutput{(resValue - amount), from}
		outputs = append(outputs, output)
	}

	// 创建交易
	tx := Transaction{nil, inputs, outputs}
	// 设置交易id
	tx.SetTXID()
	// 返回交易结构
	return &tx
}
