package main

import (
	"bytes"
	"fmt"
	"time"
)

// 创建区块链
func (cli *CLI) CreateBlockChain(address string) {
	if !IsValidAddress(address) {
		fmt.Printf("%s 是无效地址！\n", address)
		return
	}
	bc := CraeteBlockChain(address)
	if bc == nil {
		return
	}

	bc.db.Close()

}

// 打印区块
func (cli *CLI) PrintBlock() {
	bc := NewBlockChain()
	if bc == nil {
		return
	}
	defer bc.db.Close()

	it := bc.NewIterator()
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
		fmt.Printf("Data : %s\n", block.Transactions[0].TXInputs[0].PubKey)
		fmt.Printf("Hash : %x\n", block.Hash)

		pow := NewProofOfWork(block)
		fmt.Printf("IsValid: %v\n", pow.IsValid())

		if bytes.Equal(block.PreBlockHash, []byte{}) {
			break
		}
	}
}

func (cli *CLI) GetMyBalance(address string) {
	if !IsValidAddress(address) {
		fmt.Printf("%s 是无效地址！\n", address)
		return
	}

	bc := NewBlockChain()
	if bc == nil {
		return
	}
	defer bc.db.Close()

	bc.GetMyBalance(address)
}

func (cli *CLI) Send(from, to string, amount float64, miner string, data string) {
	if !IsValidAddress(from) {
		fmt.Printf("%s 是无效地址！\n", from)
		return
	}

	if !IsValidAddress(to) {
		fmt.Printf("%s 是无效地址！\n", to)
		return
	}

	if !IsValidAddress(miner) {
		fmt.Printf("%s 是无效地址！\n", miner)
		return
	}

	bc := NewBlockChain()
	if bc == nil {
		return
	}
	defer bc.db.Close()

	// 创建挖矿交易
	coinbase := NewCoinbaseTx(miner, data)
	txs := []*Transaction{coinbase}

	// 创建普通交易
	tx := NewTransaction(from, to, amount, bc)
	if tx != nil {
		txs = append(txs, tx)
	} else {
		fmt.Printf("发现无效交易，过滤！\n")
	}

	// 添加到区块
	bc.AddBlock(txs)

	fmt.Printf("挖矿成功")
}

// 创建钱包地址
func (cli *CLI) CreateWallet() {
	ws := NewWallets()
	address := ws.CreateWallet()
	fmt.Printf("新创建的钱包地址为：%s\n", address)
}

// 打印地址
func (cli *CLI) ListAddresses() {
	ws := NewWallets()
	addresses := ws.ListAddress()
	for _, address := range addresses {
		fmt.Printf("address : %s\n", address)
	}
}

// 打印交易
func (cli *CLI) PrintTX() {
	bc := NewBlockChain()

	defer bc.db.Close()

	it := bc.NewIterator()
	for {
		block := it.Next()
		fmt.Printf("\n+++++++++++++++block %x++++++++++++++++\n", block.Hash)
		for _, tx := range block.Transactions {
			fmt.Printf("tx : %v\n", tx)
		}

		if len(block.PreBlockHash) == 0 {
			break
		}
	}
}
