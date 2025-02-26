package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/btcsuite/btcutil/base58"
)

type TXInput struct {
	TXID []byte // 交易id

	Index int64 // output索引
	// Address string // 解锁脚本，先用地址模拟

	Signature []byte // 交易签名

	PubKey []byte // 公钥本身，不是公钥hash

}

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

type Transaction struct {
	TXid      []byte     // 交易id
	TXInputs  []TXInput  // 所有交易的inputs
	TXOutputs []TXOutput // 所有的outputs
}

func (tx *Transaction) SetTXID() {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	hash := sha256.Sum256(buffer.Bytes())
	tx.TXid = hash[:]
}

// 挖矿奖励
const reward = 3.125

// 实现挖矿交易
// 只有输出，没有有效输入
func NewCoinbaseTx(miner string, data string) *Transaction {
	inputs := []TXInput{{nil, -1, nil, []byte(data)}}
	// outputs := []TXOutput{{3.125, miner}}

	output := NewTXOutput(reward, miner)
	outputs := []TXOutput{output}

	tx := Transaction{nil, inputs, outputs}
	tx.SetTXID()

	return &tx
}

// 判断交易是否为挖矿交易
func (tx *Transaction) IsCoinbase() bool {
	// 1.只有一个input，2. 引用的id是nil 3. 引用的索引是-1
	inputs := tx.TXInputs
	if len(inputs) == 1 && inputs[0].TXID == nil && inputs[0].Index == -1 {
		return true
	}

	return false
}

// 创建普通交易
func NewTransaction(from, to string, amount float64, bc *BlockChain) *Transaction {
	// 1. 打开钱包
	ws := NewWallets()
	// 获取密钥对
	walllet := ws.WalletsMap[from]

	if walllet == nil {
		fmt.Printf("%s 的私钥不存在，交易创建失败！\n", from)
		return nil
	}

	// 2. 获取公钥私钥
	privateKey := walllet.PrivateKey
	publicKey := walllet.PublicKey

	pubKeyHash := HashPubKey(publicKey)

	// 标识能用的utxo
	var utxos = make(map[string][]int64)
	// 标识这些utxo存储的金额
	var resValue float64

	// 1.遍历账本，找到属于付款人的合适的金额
	utxos, resValue = bc.FindNeedUtxos(pubKeyHash, amount)

	// 2.如果找到的钱不足以转账，则交易失败
	if resValue < amount {
		fmt.Printf("余额不足，交易失败！\n")
		return nil
	}

	var inputs []TXInput
	var outputs []TXOutput

	// 3.将outputs转成inputs
	for txid, indexes := range utxos {
		for _, index := range indexes {
			inputs = append(inputs, TXInput{[]byte(txid), index, nil, publicKey})
		}
	}

	// 4.生成outputs
	output := NewTXOutput(amount, to)
	outputs = append(outputs, output)

	// 5.如果有剩余，则转给自己
	if resValue > amount {
		// output := TXOutput{(resValue - amount), from}
		output1 := NewTXOutput(resValue-amount, from)
		outputs = append(outputs, output1)
	}

	// 创建交易
	tx := Transaction{nil, inputs, outputs}
	// 设置交易id
	tx.SetTXID()

	// 对交易进行签名
	bc.SignTransaction(&tx, privateKey)

	// 返回交易结构
	return &tx
}

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
