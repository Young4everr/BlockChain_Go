package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

// 定义一个工作量证明的结构ProofOfWork

// a. block
// b. 目标值

// 2. 提供创建POW的函数
// 	- NewProofOfWork(参数)

// 3. 提供计算不断计算hash的函数
// 	- Run()

// 4. 提供一个校验函数
// 	- IsValid()

type ProofOfWork struct {
	block *Block

	// 用来存储哈希值，它内置一些方法Cmp: 比较方法
	// SetBytes : 把 bytes 转为 big.int 类型
	// SetString : 把 string 转成 big.int 类型
	target *big.Int // 系统提供的，固定的
}

const Bits = 16

func NewProofOfWork(block *Block) *ProofOfWork {
	pow := ProofOfWork{
		block: block,
	}

	// 固定难度值
	// 16进制格式字符串
	// targetStr := "0001000000000000000000000000000000000000000000000000000000000000"
	// var bigIntTmp big.Int
	// bigIntTmp.SetString(targetStr, 16)
	bigIntTmp := big.NewInt(1)
	bigIntTmp.Lsh(bigIntTmp, 256-Bits)
	pow.target = bigIntTmp
	return &pow
}

func (pow *ProofOfWork) Run() ([]byte, uint64) {
	// 1. 获取block数据
	// 2. 拼接nonce
	// 3. sha256
	// 4. 与难度值比较
	// 5. 哈希值大于难度值，nonce++
	// 6. 哈希值小于难度值，挖矿成功，退出

	var nonce uint64

	var hash [32]byte

	for {
		fmt.Printf("%x\r", hash)
		hash = sha256.Sum256(pow.PrepareData(nonce))

		// 将hash (数组类型)转成big.int，然后与pow.target进行比较，需要引入局部变量
		var bigIntTmp big.Int
		bigIntTmp.SetBytes(hash[:])
		if bigIntTmp.Cmp(pow.target) == -1 {
			fmt.Printf("挖矿成功！nonce: %d, 哈希值为: %x\n", nonce, hash)
			break
		} else {
			nonce++
		}
	}
	return hash[:], nonce
}

// 准备数据，即计算区块的hash值
func (pow *ProofOfWork) PrepareData(nonce uint64) []byte {
	block := pow.block
	tmp := [][]byte{
		UintToByte(block.Version),
		block.PreBlockHash,
		block.MerKleRoot,
		UintToByte(block.TimeStamp),
		UintToByte(block.Difficulty),
		UintToByte(nonce),
	}

	data := bytes.Join(tmp, []byte{})
	return data
}

// 有效性证明
func (pow *ProofOfWork) IsValid() bool {

	hash := sha256.Sum256(pow.PrepareData(pow.block.Nonce))
	var bigIntTmp big.Int
	bigIntTmp.SetBytes(hash[:])

	return bigIntTmp.Cmp(pow.target) == -1
}
