package main

import (
	"bytes"
	"fmt"
	"time"
)

func main() {
	fmt.Print("hello world!\n")

	// block := NewBlock(genesisInfo, []byte{0x0000000000000000})
	bc := NewBlockChain()
	defer bc.db.Close()

	bc.AddBlock("The Second Block....")
	bc.AddBlock("The Third Block....")
	bc.AddBlock("The fourth Block....")

	it := bc.NewIterator()
	for {
		block := it.Next()
		fmt.Printf("\n++++++++++++++++++++++++++++++++\n")
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

	// for i, block := range bc.Blocks {
	// 	fmt.Printf("\n++++++++++++++  %d ++++++++++++++++++\n", i)
	// 	fmt.Printf("Version : %d\n", block.Version)
	// 	fmt.Printf("PreBlockHash : %x\n", block.PreBlockHash)
	// 	fmt.Printf("MerKleRoot : %x\n", block.MerKleRoot)

	// 	timeFormat := time.Unix(int64(block.TimeStamp), 0).Format("2006-01-02 15:04:05")
	// 	fmt.Printf("TimeStamp : %s\n", timeFormat)

	// 	fmt.Printf("Difficulty : %d\n", block.Difficulty)
	// 	fmt.Printf("Nonce : %d\n", block.Nonce)
	// 	fmt.Printf("Data : %s\n", block.Data)
	// 	fmt.Printf("Hash : %x\n", block.Hash)

	// 	pow := NewProofOfWork(block)
	// 	fmt.Printf("IsValid: %v\n", pow.IsValid())
	// }
}
