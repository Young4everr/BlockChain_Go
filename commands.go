package main

import (
	"bytes"
	"fmt"
	"time"
)

// 添加区块
func (cli *CLI) addBlock(data string) {
	cli.bc.AddBlock(data)
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
		fmt.Printf("Data : %s\n", block.Data)
		fmt.Printf("Hash : %x\n", block.Hash)

		pow := NewProofOfWork(block)
		fmt.Printf("IsValid: %v\n", pow.IsValid())

		if bytes.Equal(block.PreBlockHash, []byte{}) {
			break
		}
	}
}
