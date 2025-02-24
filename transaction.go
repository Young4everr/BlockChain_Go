package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
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
func NewCoinbaseTx(miner string) *Transaction {
	inputs := []TXInput{{nil, -1, genesisInfo}}
	outputs := []TXOutput{{3.125, miner}}

	tx := Transaction{nil, inputs, outputs}
	tx.SetTXID()

	return &tx
}
