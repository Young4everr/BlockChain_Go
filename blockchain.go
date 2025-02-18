package main

// 创建区块链，使用Block数组模拟
type BlockChain struct {
	Blocks []*Block
}

// 实现创建区块链的方法
func NewBlockChain() *BlockChain {
	genesisBlock := NewBlock(genesisInfo, []byte{0x0000000000000000})
	bc := BlockChain{Blocks: []*Block{genesisBlock}}
	return &bc
}

// 添加区块
func (bc *BlockChain) AddBlock(data string) {
	// 1.创建一个区块
	// bc.Blocks的最后一个区块的Hash值就是当前新区块的PrevBlockHash
	lastBlock := bc.Blocks[len(bc.Blocks)-1]
	preHash := lastBlock.Hash

	block := NewBlock(data, preHash)

	// 2.添加到bc.Blocks中
	bc.Blocks = append(bc.Blocks, block)

}
