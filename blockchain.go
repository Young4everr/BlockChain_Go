package main

import (
	"blockchain/bolt"
	"fmt"
	"log"
	"os"
)

// v1-创建区块链，使用Block数组模拟
// v2-使用bolt进行改写，需要两个字段
//  1. bolt数据库的句柄
//  2. 最后一个区块的hash
type BlockChain struct {
	// v1
	// Blocks []*Block

	// v2 - 添加bolt存储区块
	db   *bolt.DB // 句柄
	tail []byte   // 最后一个区块的hash
}

const blockChainName = "blockChain.db"
const blockBucketName = "blockBucket"
const lastHashKey = "lastHashKey"

// 实现创建区块链的方法
func NewBlockChain(miner string) *BlockChain {
	// v1
	// genesisBlock := NewBlock(genesisInfo, []byte{0x0000000000000000})
	// bc := BlockChain{Blocks: []*Block{genesisBlock}}
	// return &bc

	// v2 - 添加bolt存储区块
	// 功能分析
	// 1. 获得数据库句柄，打开数据库，读写数据
	db, err := bolt.Open(blockChainName, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	var tail []byte

	db.Update(func(tx *bolt.Tx) error {
		// 获取桶(其实就是链)
		b := tx.Bucket([]byte(blockBucketName))
		// 当桶不存在时，创建桶，并添加创世区块
		if b == nil {
			b, err = tx.CreateBucket([]byte(blockBucketName))
			if err != nil {
				log.Panic(err)
			}

			// bucket创建完成，开始添加创世区块
			// 创世快中只有一个挖矿交易
			coinbase := NewCoinbaseTx(miner, genesisInfo)
			genesisBlock := NewBlock([]*Transaction{coinbase}, []byte{})
			b.Put(genesisBlock.Hash, genesisBlock.Serialize() /*将区块序列化转化为字节流*/)
			b.Put([]byte(lastHashKey), []byte(genesisBlock.Hash))

			// Test
			// block := b.Get(genesisBlock.Hash)
			// fmt.Printf("block: %s\n", block)

			tail = genesisBlock.Hash
		} else {
			tail = b.Get([]byte(lastHashKey))
		}

		return nil
	})

	return &BlockChain{db, tail}
}

// 添加区块
func (bc *BlockChain) AddBlock(txs []*Transaction) {
	// v1
	// 1.创建一个区块
	// bc.Blocks的最后一个区块的Hash值就是当前新区块的PrevBlockHash
	// lastBlock := bc.Blocks[len(bc.Blocks)-1]
	// preHash := lastBlock.Hash

	// block := NewBlock(data, preHash)

	// // 2.添加到bc.Blocks中
	// bc.Blocks = append(bc.Blocks, block)

	// v2 - 增加bolt存储区块
	// 1. 创建一个区块
	bc.db.Update(func(tx *bolt.Tx) error {
		// 获取桶(其实就是链)
		b := tx.Bucket([]byte(blockBucketName))
		// 当桶不存在时，报错
		if b == nil {
			log.Panic("bukcet not exist, please check!")
			os.Exit(1)
		}

		// 创建区块
		block := NewBlock(txs, bc.tail)
		b.Put(block.Hash, block.Serialize() /*将区块序列化转化为字节流*/)
		b.Put([]byte(lastHashKey), block.Hash)

		bc.tail = block.Hash

		return nil
	})
}

// 定义一个区块链迭代器，包含db, current
type BlockChainIterator struct {
	db      *bolt.DB
	current []byte //当前所指向的区块hash
}

// 创建迭代器，使用bc初始化
func (bc *BlockChain) NewIterator() *BlockChainIterator {
	return &BlockChainIterator{bc.db, bc.tail}
}

// 定义Next方法
func (it *BlockChainIterator) Next() *Block {

	var block Block

	it.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucketName))
		if b == nil {
			log.Panic("bucket not exist, please check!")
			os.Exit(1)
		}
		// 读取数据
		blockInfo /*字节流*/ := b.Get(it.current)
		block = *DeSerialize(blockInfo)

		it.current = block.PreBlockHash

		return nil
	})

	return &block
}

// 查询我的比特币余额
// // 1. 遍历账本
// // 2. 遍历交易
// // 3. 遍历output
// 4. 找到属于我的output
func (bc *BlockChain) FindMyUtxos(address string) []TXOutput {
	fmt.Printf("FindMyUtxos\n")
	var UTXOs []TXOutput

	it := bc.NewIterator()

	//key是交易id，value是index数组
	spentUTXO := make(map[string][]int64)

	// 遍历账本
	for {
		block := it.Next()

		// 遍历交易
		for _, tx := range block.Transactions {
			// 遍历交易输入
			for _, input := range tx.TXInputs {
				if input.Address == address {
					fmt.Printf("找到消耗过的output! index : %d\n", input.Index)
					spentUTXO[string(input.TXID)] = append(spentUTXO[string(input.TXID)], input.Index)
				}
			}

		OUTPUT:
			// 遍历output
			for i, output := range tx.TXOutputs {
				key := string(tx.TXid)
				indexes := spentUTXO[key]
				if len(indexes) != 0 {
					fmt.Printf("当前这笔交易中有被消耗过的output!")
					for _, j := range indexes {
						if int64(i) == j {
							continue OUTPUT
						}
					}
				}

				// 找到属于我的所有output
				if address == output.Address {
					fmt.Printf("找到属于 %s 的output, i : %d\n", address, i)
					UTXOs = append(UTXOs, output)
				}
			}
		}

		if len(block.PreBlockHash) == 0 {
			fmt.Printf("遍历区块链结束！\n")
			break
		}
	}

	return UTXOs
}

func (bc *BlockChain) GetMyBalance(address string) {

	utxos := bc.FindMyUtxos(address)

	var total = 0.0

	for _, utxo := range utxos {
		total += utxo.Value
	}
	fmt.Printf("%s 的余额是：%f", address, total)
}

// 查找合理的utxo
func (bc *BlockChain) FindNeedUtxos(from string, amount float64) (map[string][]int64, float64) {

	needUtxos := make(map[string][]int64)
	var resValue float64

	// 找到合理的utxo
	it := bc.NewIterator()

	spentUTXO := make(map[string][]int64)

	// 遍历账本
	for {
		block := it.Next()

		// 遍历交易
		for _, tx := range block.Transactions {
			// 遍历inputs
			for i, input := range tx.TXInputs {
				if input.Address == from {
					fmt.Printf("找到消耗过的output! index : %d\n", input.Index)
					spentUTXO[string(input.TXID)] = append(spentUTXO[string(input.TXID)], int64(i))
				}
			}

		OUTPUT:
			// 遍历output，找到合适的utxo
			for i, output := range tx.TXOutputs {
				key := tx.TXid
				if len(spentUTXO[string(key)]) != 0 {
					fmt.Printf("当前这笔交易中有被消耗过的output!\n")
					for _, index := range spentUTXO[string(key)] {
						if index == int64(i) {
							continue OUTPUT
						}
					}
				}
				if from == output.Address {
					fmt.Printf("找到属于 %s 的output, i : %d\n", from, i)
					// 添加到返回结构needUtxos
					needUtxos[string(key)] = append(needUtxos[string(key)], int64(i))
					resValue += output.Value

					// 判断金额是否足够，如足够则退出
					if resValue >= amount {
						return needUtxos, resValue
					}
					// 如不足，则继续上述循环
				}
			}
		}

		if len(block.PreBlockHash) == 0 {
			fmt.Printf("遍历区块链结束！\n")
			break
		}
	}

	return needUtxos, resValue
}
