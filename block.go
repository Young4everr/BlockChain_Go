package main

import (
	"bytes"
	"crypto/sha256"
	"time"
)

// 1. 定义结构（区块头的字段比正常的少）
//     1. 前区块哈希
//     2. 当前区块哈希
//     3. 数据
// 2. 创建区块
// 3. 生成哈希
// 4. 引入区块链
// 5. 添加区块
// 6. 重构代码

const genesisInfo = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"

type Block struct {
	Version uint64 //区块版本号

	PreBlockHash []byte //前区块哈希

	MerKleRoot []byte //先填写为空，后续v4的时候使用

	TimeStamp uint64 //从1970.1.1至今的秒数

	Difficulty uint64 //挖矿的难度值，v2时使用

	Nonce uint64 //随机数，挖矿找的就是它

	Data []byte //数据，目前使用字节流，v4开始使用交易代替

	Hash []byte //当前区块哈希，区块中不存在的字段，为了方便我们添加进来
}

// 创建区块，对Block的每个字段填充数据即可
func NewBlock(data string, preBlockHash []byte) *Block {
	block := Block{
		Version: 00,

		PreBlockHash: preBlockHash,

		MerKleRoot: []byte{},

		TimeStamp: uint64(time.Now().Unix()),

		Difficulty: Bits, //随便写的，v2在调整

		// Nonce: 10, //同Difficulty

		Data: []byte(data),

		Hash: []byte{}, //先填充为空，后续会填充数据
	}

	// block.SetHash()
	pow := NewProofOfWork(&block)
	hash, nonce := pow.Run()

	block.Hash = hash
	block.Nonce = nonce

	return &block
}

// 为了生成区块哈希，我们实现一个简单的函数，来计算哈希值，没有难度值
func (block *Block) SetHash() {
	// var data []byte
	// data = append(data, UintToByte(block.Version)...)
	// data = append(data, block.PreBlockHash...)
	// data = append(data, block.MerKleRoot...)
	// data = append(data, UintToByte(block.TimeStamp)...)
	// data = append(data, UintToByte(block.Difficulty)...)
	// data = append(data, UintToByte(block.Nonce)...)
	// data = append(data, block.Data...)

	tmp := [][]byte{
		UintToByte(block.Version),
		block.PreBlockHash,
		block.MerKleRoot,
		UintToByte(block.TimeStamp),
		UintToByte(block.Difficulty),
		UintToByte(block.Nonce),
		block.Data,
	}

	data := bytes.Join(tmp, []byte{})

	hash /*[32]byte*/ := sha256.Sum256(data)
	block.Hash = hash[:]
}
