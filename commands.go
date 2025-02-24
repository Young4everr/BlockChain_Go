package main

import (
	"bytes"
	"fmt"
	"time"
)

// 添加区块
func (cli *CLI) addBlock(txs []*Transaction) {
	cli.bc.AddBlock(txs)
}

// 打印区块
func (cli *CLI) printBlock() {
	it := cli.bc.NewIterator()
	for {
		fmt.Printf("\n++++++++++++++++++++++++++++++++\n")
		block := it.Next()
		fmt.Printf("Version : %d\n", block.Version)
		fmt.Printf("PreBlockHash : %x\n", block.PreBlockHash)
		fmt.Printf("MerKleRoot : %x\n", block.MerKleRoot)

		timeFormat := time.Unix(int64(block.TimeStamp), 0).Format("2006-01-02 15:04:05")
		fmt.Printf("TimeStamp : %s\n", timeFormat)

		fmt.Printf("Difficulty : %d\n", block.Difficulty)
		fmt.Printf("Nonce : %d\n", block.Nonce)
		fmt.Printf("Data : %s\n", block.Transactions[0].TXInputs[0].Address)
		fmt.Printf("Hash : %x\n", block.Hash)

		pow := NewProofOfWork(block)
		fmt.Printf("IsValid: %v\n", pow.IsValid())

		if bytes.Equal(block.PreBlockHash, []byte{}) {
			break
		}
	}
}

func (cli *CLI) Send(from, to string, amount float64, miner string, data string) {
	// 创建挖矿交易
	coinbase := NewCoinbaseTx(miner, data)
	txs := []*Transaction{coinbase}

	// 创建普通交易
	tx := NewTransaction(from, to, amount, cli.bc)
	if tx != nil {
		txs = append(txs, tx)
	} else {
		fmt.Printf("发现无效交易，过滤！\n")
	}

	// 添加到区块
	cli.bc.AddBlock(txs)

	fmt.Printf("挖矿成功")
}
