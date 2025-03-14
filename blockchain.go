package main

import (
	"bytes"
	"crypto/ecdsa"
	"fmt"
	"log"
	"os"

	"github.com/boltdb/bolt"
	"github.com/btcsuite/btcutil/base58"
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
func CraeteBlockChain(miner string) *BlockChain {
	if IsFileExist(blockChainName) {
		fmt.Printf("区块链已存在，无需重复创建！")
		return nil
	}

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

		b, err := tx.CreateBucket([]byte(blockBucketName))
		if err != nil {
			log.Panic(err)
		}

		// bucket创建完成，开始添加创世区块
		// 创世快中只有一个挖矿交易
		coinbase := NewCoinbaseTx(miner, genesisInfo)
		genesisBlock := NewBlock([]*Transaction{coinbase}, []byte{})
		b.Put(genesisBlock.Hash, genesisBlock.Serialize() /*将区块序列化转化为字节流*/)
		b.Put([]byte(lastHashKey), []byte(genesisBlock.Hash))

		tail = genesisBlock.Hash

		return nil
	})

	return &BlockChain{db, tail}
}

// 实现获取区块链实例
func NewBlockChain() *BlockChain {
	if !IsFileExist(blockChainName) {
		fmt.Printf("区块链不存在，请先创建！")
		return nil
	}

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

	db.View(func(tx *bolt.Tx) error {
		// 获取桶(其实就是链)
		b := tx.Bucket([]byte(blockBucketName))
		// 当桶不存在时，创建桶，并添加创世区块
		if b == nil {
			fmt.Printf("区块链bucket为空，请检查！\n")
			os.Exit(1)
		}
		tail = b.Get([]byte(lastHashKey))

		return nil
	})

	return &BlockChain{db, tail}
}

// 添加区块
func (bc *BlockChain) AddBlock(txs []*Transaction) {
	// 矿工得到交易时，第一时间对交易进行验证
	// 矿工如果不验证，即使挖矿成功，广播区块后，其他验证矿工任然会校验每一笔交易
	validTxs := []*Transaction{}
	for _, tx := range txs {
		if bc.VerifyTransaction(tx) {
			fmt.Printf("发现有效的交易：%x\n", tx.TXid)
			validTxs = append(validTxs, tx)
		} else {
			fmt.Printf("发现无效的交易：%x\n", tx.TXid)
		}
	}

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
		block := NewBlock(validTxs, bc.tail)
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

// 定义UTXOInfo结构
type UTXOInfo struct {
	TXID   []byte   // 交易id
	Index  int64    // output的索引值
	Output TXOutput // output本身
}

// 查询我的比特币余额
// // 1. 遍历账本
// // 2. 遍历交易
// // 3. 遍历output
// 4. 找到属于我的output
func (bc *BlockChain) FindMyUtxos(pubKeyHash []byte) []UTXOInfo {
	fmt.Printf("FindMyUtxos\n")
	var UTXOInfos []UTXOInfo

	it := bc.NewIterator()

	//key是交易id，value是index数组
	spentUTXO := make(map[string][]int64)

	// 遍历账本
	for {
		block := it.Next()

		// 遍历交易
		for _, tx := range block.Transactions {
			if !tx.IsCoinbase() {
				// 遍历交易输入
				for _, input := range tx.TXInputs {
					// 判断当前被使用的input是否为目标地址所有
					if bytes.Equal(HashPubKey(input.PubKey), pubKeyHash) {
						fmt.Printf("找到消耗过的output! index : %d\n", input.Index)
						spentUTXO[string(input.TXID)] = append(spentUTXO[string(input.TXID)], input.Index)
					}
				}
			}

		OUTPUT:
			// 遍历output
			for i, output := range tx.TXOutputs {
				key := tx.TXid
				indexes := spentUTXO[string(key)]
				if len(indexes) != 0 {
					fmt.Printf("当前这笔交易中有被消耗过的output!")
					for _, j := range indexes {
						if int64(i) == j {
							continue OUTPUT
						}
					}
				}

				// 找到属于我的所有output
				if bytes.Equal(pubKeyHash, output.PubKeyHash) {
					utxoinfo := UTXOInfo{tx.TXid, int64(i), output}
					UTXOInfos = append(UTXOInfos, utxoinfo)
				}
			}
		}

		if len(block.PreBlockHash) == 0 {
			fmt.Printf("遍历区块链结束！\n")
			break
		}
	}

	return UTXOInfos
}

func (bc *BlockChain) GetMyBalance(address string) {
	hash := base58.Decode(address)
	pubKeyHash := hash[1 : len(hash)-4]

	utxoinfos := bc.FindMyUtxos(pubKeyHash)

	var total = 0.0

	for _, utxoinfo := range utxoinfos {
		total += utxoinfo.Output.Value
	}
	fmt.Printf("%s 的余额是：%f", address, total)
}

// 查找合理的utxo
func (bc *BlockChain) FindNeedUtxos(pubKeyHash []byte, amount float64) (map[string][]int64, float64) {
	needUtxos := make(map[string][]int64)
	var resValue float64

	UTXOInfos := bc.FindMyUtxos(pubKeyHash)
	for _, utxoinfo := range UTXOInfos {
		key := string(utxoinfo.TXID)
		needUtxos[key] = append(needUtxos[key], utxoinfo.Index)
		resValue += utxoinfo.Output.Value
		// 判断金额是否足够
		if resValue >= amount {
			break
		}
	}

	return needUtxos, resValue
}

func (bc *BlockChain) SignTransaction(tx *Transaction, privateKey *ecdsa.PrivateKey) {
	// 遍历账本找到所有应用的交易
	prevTXs := make(map[string]*Transaction)

	// 遍历inputs，通过id查找所引用的交易
	for _, input := range tx.TXInputs {
		prevTX := bc.FindTransaction(input.TXID)
		if prevTX == nil {
			fmt.Printf("没有找到交易：%s\n", input.TXID)
		} else {
			prevTXs[string(input.TXID)] = prevTX
		}
	}

	tx.Sign(privateKey, prevTXs)
}

func (bc *BlockChain) VerifyTransaction(tx *Transaction) bool {
	// 挖矿交易直接返回true
	if tx.IsCoinbase() {
		return true
	}
	// 遍历账本找到所有应用的交易
	prevTXs := make(map[string]*Transaction)

	// 遍历inputs，通过id查找所引用的交易
	for _, input := range tx.TXInputs {
		prevTX := bc.FindTransaction(input.TXID)
		if prevTX == nil {
			fmt.Printf("没有找到交易：%s\n", input.TXID)
		} else {
			prevTXs[string(input.TXID)] = prevTX
		}
	}

	return tx.Verify(prevTXs)
}

func (bc *BlockChain) FindTransaction(txid []byte) *Transaction {

	it := bc.NewIterator()
	// 遍历账本
	for {
		block := it.Next()
		// 遍历交易
		for _, tx := range block.Transactions {
			if bytes.Equal(tx.TXid, txid) {
				fmt.Printf("找到了所引用的交易：%x\n", tx.TXid)
				return tx
			}
		}

		if len(block.PreBlockHash) == 0 {
			break
		}
	}

	return nil
}
