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
