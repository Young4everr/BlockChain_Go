# 一、人类交易发展史

#### 1. 物物交换

#### 2. 实物货币(自然货币)

#### 3. 传统货币(人工货币)

- #### 金属货币(称量货币)

- #### 纸币

- #### 电子货币

  - #### 银行卡

- #### 数字货币

  

# 二、比特币诞生背景

#### 1. 纸币的风险

- #### 政府信用 - 津巴布韦

#### 2. 金融危机

#### 3. 密码朋克组织

- #### 加密电子邮件系统

#### 4. 比特币诞生

- #### 中本聪

  - #### [《比特币：一种点对点的电子现金系统》](https://bitcoin.org/bitcoin.pdf)

  - #### [创世区块](https://www.blockchain.com/explorer/blocks/btc/0)

#### 5. 背景

- #### 直接原因：金融危机，技术成熟。

- #### 根本原因：纸币受到国家的影响太大。

  

# 三、比特币详细介绍

#### 1. 是什么？

- #### 比特币是一个软件，是一个电子支付系统，它安全可靠，解决了纸币被政府随意超发的问题，避免通货膨胀

#### 2. 中心化与去中心化

- #### 中心化：寡头&垄断&单点故障

- #### 去中心化：平等&互联&防篡改

- #### **<u>比特币就是一个具有转账支付功能的去中心化系统</u>**

#### 3. 区块链

- #### 直观理解

  - #### 交易产生数据，数据存储在数据库中，它有自己的存储结构，这个数据块叫做区块。

  - #### 将所有的区块按照这个区块的哈希链接起来的链条，叫做区块链。

#### 4. 传统记账VS比特币记账

- #### 传统：账单 -> 账本

- #### 比特币：区块 -> 区块链

#### 5. 几个概念

- #### 钱包

  - #### 创建私钥公钥，保存私钥，相当于钱包，可以存放多个地址

  ```
  地址：类似于钱包中不同的银行卡 （由公钥生成）
  私钥：类似于每张银行卡有一个密码 （私钥可以解开地址）
  ```

  - #### 一个私钥对应一个地址

  - #### 一个钱包中可以放很多私钥对 (类似于一个钱包中可以放很多银行卡)

  - #### wallet.dat是真正的钱包，负责维护私钥和地址，但是一般把维护这个钱包的客户端，统称为钱包

  - #### 比特币转账每次都会自动生成新的地址，从而隐藏自己的资产

  - #### wallet.dat会自动帮助我们维护地址和私钥，保存好这个文件即可（如果有密码，一定要记住）

- #### 节点

  - #### 每一个运行挖矿软件的人变成为一个区块链网络的节点

  - #### 轻节点：不下载所有账本，只下载区块头，和自己有关的交易

  - #### 全节点：强调是全指的是账本，查看当前总容量

- #### 账本

  - #### 所有交易信息集合，使用数据库，每个节点都同步一个账本

- #### 挖矿

  - #### 节点间竞争记账权力的过程就叫做挖矿，竞争成功者获得记账的权力，即挖到矿。

  - #### 挖矿的过程就是比特币发行的过程。

  - #### 本质：对区块数据做哈希运算，寻找一个满足条件的随机数。

- #### 区块链演示

  - https://andersbrownworth.com/blockchain/blockchain
  - https://blockchaindemo.io/

- #### 算力

  - 

- #### 矿机

  - #### CPU -> GPU -> FPGA -> ASIC(专业矿机)

- #### 矿场

- #### 矿池

#### 6. 比特币系统参数

- #### 出块时间

  - #### 10min

- #### 出块奖励 

  - #### 最初50BTC，每21万个块奖励减半，10min，21万 * 10min / 60 / 24h / 365 = 4年，4年奖励衰减一次

- #### 比特币总量

  - #### 总共2100万个比特币，2140年挖完

    ```go
    import "fmt"
    
    const blockInterval = 21 //区块衰减区间，单位是万
    
    func main() {
        total := 0.0  //总量
        reward := 50.0  //最初奖励50个，每过21万个区块，奖励减半，大约是4年
        
        for reward > 0 {
            amount := reward * blockInterval  //在21万个块中产生比特币的数量
            total = total + amount
            // reward /= 2  //比特币奖励衰减一半
            reward *= 0.5  //除法效率低，使用乘法代替
        }
        
        fmt.Printf("比特币总量为：%f\n", total)
    }
    ```

- #### 区块容量

  - #### 考虑同步效率

  - #### 1M大约容纳4000条交易 （已经扩容）

  - #### 1M / 每笔交易的字节数 = 交易数 （1024 * 1024  / 223 = 4200）

- #### 每秒成交量

  - #### 4200 / 600s = 7笔/秒

- #### 单位

  - #### 1BTC = 10**8sat(聪)

    

# 四、比特币依赖技术

#### 1. 交易是一切数据的来源

#### 2. 交易流程

- #### 各个节点交易可以不相同，即包含的交易数可以不相同

#### 3. 技术

1. #### 密码学

   1. ##### 对称加密

   2. ##### 非对称加密

      - ###### 公钥：加密，保护隐私

      - ###### 私钥：签名，1、保证数据来源，2、保证数据未被篡改，3、签名人无法否认是自己签的

      - ![image](https://github.com/user-attachments/assets/7f35ec76-e186-4b16-87e4-d5bd4633e6eb)


2. ### P2P网络

3. #### 工作量证明（挖矿）

   1. ![image-20250217001518895](https://github.com/user-attachments/assets/397e511c-bc63-4fde-9bd9-eedd8d47e880)
   2. ```go
      import (
          "crypto/sha256"
          "fmt"
      )
      
      func main() {
          var data = "blockchain"
          for i := 0; i < 1000000; i++ {
              hash := sha256.Sum256([]byte(data + string(i)))
              fmt.Printf("hash : %x, d%\n", string(hash[:]), i)
          }
      }
      ```

4. #### base58

   1. ##### 比特币地址生成

      1. ![image-20250217001518895](https://github.com/user-attachments/assets/ec646fb1-bbde-497e-aabb-a601b1edd4f1)


      2. ![image-20250217003518300](https://github.com/user-attachments/assets/c84ebb23-8900-4f9e-8a7c-ec02f657f985)


      3. ###### [演示网址](https://gobittest.appspot.com/Address)



# 五、区块结构

![image-20250217004852181](https://github.com/user-attachments/assets/51e5616f-cbff-4eba-9dcb-436a5dea0a65)


​	完整结构

​	![image-20250217011522866](https://github.com/user-attachments/assets/b2252bce-feec-49fc-9c20-9ce61bc66277)


#### 1. 区块头（Block Header）

- ![image-20250217004957819](https://github.com/user-attachments/assets/d1ca9e98-c5ce-461b-8055-8725de1ed821)


#### 2. 区块体

1. ##### Coinbase交易

   - 第一条交易，挖矿奖励矿工。永远是第一条，没有输入（钱的来源），只有输出（钱的流向）

2. ##### 普通转账交易

   - 每笔交易包括付款方、收款方、付款金额、手续费等



# 六、代码实现

#### 1. 区块定义和区块链定义

1. ##### 创建区块方法

2. ##### 创建区块链方法

#### 2. 工作量证明结构定义

1. ##### Run方法

#### 3. 区块存储

1. ##### 采用bolt数据库

   1. https://github.com/boltdb/bolt

2. ##### bucket中存储两种数据

   1. ###### 区块，区块的hash作为key，区块的字节流作为value

      1. block.Hash -> block.toBytes()

   2. ###### 最后一个区块的hash值

      1. key使用固定的字符串：[]byte("lastHashKey"), value就是最后一个区块的hash

   3. ![image-20250220155056815](https://github.com/user-attachments/assets/1e92aa96-0c6f-4e6f-8fc9-0844c5ef1c07)


3. ##### 创建区块链迭代器

   1. ###### 定义迭代器结构

      1. ```go
         // 定义一个区块链迭代器，包含db, current
         type BlockChainIterator struct {
         	db      *bolt.DB
         	current []byte //当前所指向的区块hash
         }
         ```

   2. ###### 初始化

      1. ```go
         // 创建迭代器，使用bc初始化
         func (bc *BlockChain) NewIterator() *BlockChainIterator {
         	return &BlockChainIterator{bc.db, bc.tail}
         }
         ```

   3. ###### 实现Next方法

      1. ```go
         // 定义Next方法
         func (it BlockChainIterator) Next() *Block {
         
         	var block Block
         	it.db.View(func(tx *bolt.Tx) error {
         		b := tx.Bucket([]byte(blockBucketName))
         		if b != nil {
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
         ```

#### 4. 命令行

1. ##### 操作

   1. ###### addBlock

   2. ###### printChain

2. ##### cli.go

3. ##### commands.go

#### 5. 比特币余额

1. ##### 交易

   1. 一对一
   2. 多对一
   3. 多对多
   4. 一对多

2. ##### 锁定、解锁

   1. ![image-20250220155056815](https://github.com/user-attachments/assets/dbac2d3c-0edd-4fc7-9434-80cedd9ece00)


3. ##### UTXO（未消费输出，unspent transaction output）



#### 6. 交易结构

```go
type TXInput struct {
	TXID    []byte // 交易id
	Index   int64  // output索引
	Address string // 解锁脚本，先用地址模拟
}

type TXOutput struct {
	Value   float64 // 转账金额
	Address string  // 锁定脚本
}

type Transaction struct {
	TXid      []byte     // 交易id
	TXInputs  []TXInput  // 所有交易的inputs
	TXOutputs []TXOutput // 所有的outputs
}
```

#### 7. 我有多少可用比特币

1. ##### 遍历账本，找到能够支配的utxo

2. ##### 剔除已经花费过的output

3. ```go
   func (bc *BlockChain) FindMyUtxos(address string) []TXOutput {
   	fmt.Printf("FindMyUtxos\n")
   	var UTXOs []TXOutput
   
   	it := bc.NewIterator()
   
   	// 遍历账本
   	for {
   		block := it.Next()
   
   		//key是交易id，value是index数组
   		spentUTXO := make(map[string][]int64)
   
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
   ```

   


#### 8. 创建普通交易

```go
// 创建普通交易
func NewTransaction(from, to string, amount float64, bc *BlockChain) *Transaction {
	// 标识能用的utxo
	utxos := make(map[string][]int64)
	// 标识这些utxo存储的金额
	var resValue float64

	// 1.遍历账本，找到属于付款人的合适的金额
	utxos, resValue = bc.FindNeedUtxos(from, amount)

	// 2.如果找到的钱不足以转账，则交易失败
	if resValue < amount {
		fmt.Printf("余额不足，交易失败！")
		return nil
	}

	var inputs []TXInput
	var outputs []TXOutput

	// 3.将outputs转成inputs
	for txid, indexes := range utxos {
		for _, index := range indexes {
			inputs = append(inputs, TXInput{[]byte(txid), index, from})
		}
	}

	// 4.生成outputs
	outputs = append(outputs, TXOutput{amount, to})

	// 5.如果有剩余，则转给自己
	if resValue > amount {
		output := TXOutput{(resValue - amount), from}
		outputs = append(outputs, output)
	}

	// 创建交易
	tx := Transaction{nil, inputs, outputs}
	// 设置交易id
	tx.SetTXID()
	// 返回交易结构
	return &tx
}
```



#### 9. 生成新的比特币地址

```go
// 生成地址
func (w *WalletKeyPair) GetAddress() string {
	hash := sha256.Sum256(w.PublicKey)

	// 创建一个hash160对象
	// 向hash160中write数据
	// 做hash运算
	rip160Hasher := ripemd160.New()
	_, err := rip160Hasher.Write(hash[:])

	if err != nil {
		log.Panic(err)
	}

	// Sum函数会把我们的结果与Sum参数append到一起，然后返回，传入nil防止数据污染
	publicHash := rip160Hasher.Sum(nil)

	version := 0x00

	// 21个字节数据
	payload := append([]byte{byte(version)}, publicHash...)

	first := sha256.Sum256(payload)
	second := sha256.Sum256(first[:])

	// 4个字节
	checksum := second[0:4]

	// 25个字节
	payload = append(payload, checksum...)

	address := base58.Encode(payload)

	return address

}
```

#### 10. 创建钱包

```go
package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"math/big"
)

const WalletName = "wallet.dat"

type Wallets struct {
	WalletsMap map[string]*WalletKeyPair
}

// SerializableWallet Solve gob: type elliptic.p256Curve has no exported fields
type SerializableWallet struct {
	D         *big.Int
	X, Y      *big.Int
	PublicKey []byte
}

// 创建Wallets，返回Wallets实例
func NewWallets() *Wallets {
	var ws Wallets
	ws.WalletsMap = make(map[string]*WalletKeyPair)
	// 1. 把所有的钱包从本地加载出来
	if !ws.LoadFromFile() {
		fmt.Printf("加载钱包数据失败！\n")
	}

	// 2. 把实例返回
	return &ws
}

// Wallets是对外的，WalletKeyPair是对内的
// Wallets调用WalletKeyPair

func (ws Wallets) CreateWallet() string {
	// 调用WalletKeyPair
	wallet := NewWalletKeyPairr()
	// 将返回的WalletKeyPair添加到WalletsMap中
	address := wallet.GetAddress()

	ws.WalletsMap[address] = wallet
	// 保存到本地文件
	res := ws.SaveToFile()
	if !res {
		fmt.Printf("创建钱包失败！\n")
		return ""
	}

	return address
}

// 保存钱包到文件
func (ws *Wallets) SaveToFile() bool {
	var buffer bytes.Buffer

	// 将接口类型明确注册一下，否则gob编码失败
	gob.Register(SerializableWallet{})

	wallets := make(map[string]SerializableWallet)
	for k, v := range ws.WalletsMap {
		wallets[k] = SerializableWallet{
			D:         v.PrivateKey.D,
			X:         v.PrivateKey.PublicKey.X,
			Y:         v.PrivateKey.PublicKey.Y,
			PublicKey: v.PublicKey,
		}
	}

	encoder := gob.NewEncoder(&buffer)

	err := encoder.Encode(wallets)

	if err != nil {
		fmt.Printf("钱包序列化失败: %v\n", err)
		return false
	}

	content := buffer.Bytes()
	err = ioutil.WriteFile(WalletName, content, 0600)

	return err == nil
}

func (ws *Wallets) LoadFromFile() bool {
	if !IsFileExist(WalletName) {
		fmt.Printf("钱包文件不存在，请先创建！\n")
		return true
	}

	// 读取文件
	content, err := ioutil.ReadFile(WalletName)
	if err != nil {
		return false
	}

	// gob解码
	var wallets map[string]SerializableWallet
	// gob.Register(elliptic.P256())
	gob.Register(SerializableWallet{})
	decoder := gob.NewDecoder(bytes.NewReader(content))
	err = decoder.Decode(&wallets)

	if err != nil {
		return false
	}

	// 赋值给ws
	ws.WalletsMap = make(map[string]*WalletKeyPair)
	for k, v := range wallets {
		ws.WalletsMap[k] = &WalletKeyPair{
			PrivateKey: &ecdsa.PrivateKey{
				PublicKey: ecdsa.PublicKey{
					Curve: elliptic.P256(),
					X:     v.X,
					Y:     v.Y,
				},
				D: v.D,
			},
			PublicKey: v.PublicKey,
		}
	}
	return true
}

// 获取所有钱包地址
func (ws *Wallets) ListAddress() []string {

	var addresses []string

	for address, _ := range ws.WalletsMap {
		addresses = append(addresses, address)
	}

	return addresses
}

```

#### 11. 改写交易

1. TXInput

   1. ```go
      type TXInput struct {
      	TXID []byte // 交易id
      
      	Index int64 // output索引
      	// Address string // 解锁脚本，先用地址模拟
      
      	Signature []byte // 交易签名
      
      	PubKey []byte // 公钥本身，不是公钥hash
      
      }
      ```

2. TXOutput

   1. ```go
      type TXOutput struct {
      	Value float64 // 转账金额
      	// Address string  // 锁定脚本
      
      	PubKeyHash []byte // 公钥hash，不是公钥本身
      }
      
      // 给定转账地址，得到这个地址的公钥hash，完成对output的锁定
      func (output *TXOutput) Lock(address string) {
      	// address -> public key hash
      	decodeInfo := base58.Decode(address)
      
      	pubKeyHash := decodeInfo[1 : len(decodeInfo)-4]
      
      	output.PubKeyHash = pubKeyHash
      }
      
      func NewTXOutput(value float64, address string) TXOutput {
      	output := TXOutput{Value: value}
      
      	output.Lock(address)
      
      	return output
      }
      ```

#### 12. 地址有效性校验

1. ```go
   func IsValidAddress(address string) bool {
   	decodeInfo := base58.Decode(address)
   	if len(decodeInfo) != 25 {
   		return false
   	}
   
   	payload := decodeInfo[0 : len(decodeInfo)-4]
   	// 自行计算出的校验码
   	checkSum1 := CheckSum(payload)
   
   	// 解出的原始校验码
   	checkSum2 := decodeInfo[len(decodeInfo)-4:]
   
   	// 对比是否一致
   	return bytes.Equal(checkSum1, checkSum2)
   
   }
   ```

#### 13. 交易签名及验证

1. ```go
   // 第一个参数是私钥
   // 第二个参数是这个交易input所引用的所有的交易
   func (tx *Transaction) Sign(privKey *ecdsa.PrivateKey, prevTXs map[string]*Transaction) {
   	fmt.Printf("对交易进行签名！\n")
   
   	// 挖矿交易
   	if tx.IsCoinbase() {
   		return
   	}
   
   	// 1.拷贝一份交易txCopy
   	//    做相应剪裁：把每个input的Sig和pubkey设置为nil
   	//    output不做改变
   	txCopy := tx.TrimmedCopy()
   
   	// 2.遍历txCopy.inputs, 把这个inpit所引用的output的公钥hash拿过来，赋值给pubkey
   	for i, input := range txCopy.TXInputs {
   		prevTX := prevTXs[string(input.TXID)]
   		pubKeyHash := prevTX.TXOutputs[input.Index].PubKeyHash
   		txCopy.TXInputs[i].PubKey = pubKeyHash
   
   		// 3.生成要签名的数据
   		txCopy.SetTXID()
   		signData := txCopy.TXid
   
   		// 4.对数据进行签名r, s
   		r, s, err := ecdsa.Sign(rand.Reader, privKey, signData)
   		if err != nil {
   			fmt.Printf("交易签名失败，err: %v\n", err)
   		}
   
   		// 5.拼接r, s为字节流，赋值给原始的交易Signature字段
   		signature := append(r.Bytes(), s.Bytes()...)
   		tx.TXInputs[i].Signature = signature
   
   		txCopy.TXInputs[i].PubKey = nil
   	}
   }
   
   // 构造裁剪tx
   func (tx *Transaction) TrimmedCopy() Transaction {
   	var inputs []TXInput
   	var outputs []TXOutput
   
   	for _, input := range tx.TXInputs {
   		input1 := TXInput{input.TXID, input.Index, nil, nil}
   		inputs = append(inputs, input1)
   	}
   
   	outputs = tx.TXOutputs
   	tx1 := Transaction{tx.TXid, inputs, outputs}
   
   	return tx1
   
   }
   
   // 验证交易
   func (tx *Transaction) Verify(prevTXs map[string]*Transaction) bool {
   	fmt.Printf("对交易进行验证！\n")
   
   	txCopy := tx.TrimmedCopy()
   	for i, input := range tx.TXInputs {
   
   		// 构造签名数据
   		prevTX := prevTXs[string(input.TXID)]
   		pubKeyHash := prevTX.TXOutputs[input.Index].PubKeyHash
   		txCopy.TXInputs[i].PubKey = pubKeyHash
   		txCopy.SetTXID()
   		verifyData := txCopy.TXid
   		txCopy.TXInputs[i].PubKey = nil
   
   		// 还原签名r, s
   		signature := input.Signature
   		r := big.Int{}
   		s := big.Int{}
   
   		rData := signature[:len(signature)/2]
   		sData := signature[len(signature)/2:]
   		r.SetBytes(rData)
   		s.SetBytes(sData)
   
   		// 还原pubKey
   		pubKeyBytes := input.PubKey
   		x := big.Int{}
   		y := big.Int{}
   		xData := pubKeyBytes[:len(pubKeyBytes)/2]
   		yData := pubKeyBytes[len(pubKeyBytes)/2:]
   		x.SetBytes(xData)
   		y.SetBytes(yData)
   		curve := elliptic.P256()
   		publicKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}
   
   		// 校验
   		if !ecdsa.Verify(&publicKey, verifyData, &r, &s) {
   			return false
   		}
   	}
   
   	return true
   }
   ```

#### 14. 打印

1. ```go
   func (tx Transaction) String() string {
   	var lines []string
   
   	lines = append(lines, fmt.Sprintf("--- Transaction %x:", tx.TXid))
   
   	for i, input := range tx.TXInputs {
   
   		lines = append(lines, fmt.Sprintf("     Input %d:", i))
   		lines = append(lines, fmt.Sprintf("       TXID:      %x", input.TXID))
   		lines = append(lines, fmt.Sprintf("       Out:       %d", input.Index))
   		lines = append(lines, fmt.Sprintf("       Signature: %x", input.Signature))
   		lines = append(lines, fmt.Sprintf("       PubKey:    %x", input.PubKey))
   	}
   
   	for i, output := range tx.TXOutputs {
   		lines = append(lines, fmt.Sprintf("     Output %d:", i))
   		lines = append(lines, fmt.Sprintf("       Value:  %f", output.Value))
   		lines = append(lines, fmt.Sprintf("       Script: %x", output.PubKeyHash))
   	}
   
   	return strings.Join(lines, "\n")
   }
   ```

   
