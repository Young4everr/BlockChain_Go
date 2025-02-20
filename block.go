package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
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

// 序列化，将区块转换成字节流
func (block *Block) Serialize() []byte {

	var buffer bytes.Buffer
	// 创建编码器
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(block)
	if err != nil {
		log.Panic("编码出错!")
	}

	return buffer.Bytes()
}

// 反序列化
func DeSerialize(data []byte) *Block {
	fmt.Printf("解码传入数据为：%x\n", data)

	var block Block
	// 创建解码器
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic("解码出错！")
	}

	return &block
}
