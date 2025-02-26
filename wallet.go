package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"log"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

// 1. 创建一个结构WalletKeyPair密钥对，保存公钥和私钥
// 2. 给这个结构提供一个方法GetAddress: 私钥->公钥->地址

type WalletKeyPair struct {
	PrivateKey *ecdsa.PrivateKey

	// type PublicKey struct {
	// 	elliptic.Curve
	// 	X, Y *big.Int
	// }
	// 将公钥的X, Y进行字节流拼接后传输，在对端再进行切割还原，好处是可以方便后面的编码
	PublicKey []byte
}

func NewWalletKeyPairr() *WalletKeyPair {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	if err != nil {
		log.Panic(err)
	}

	publicKeyRaw := privateKey.PublicKey
	publicKey := append(publicKeyRaw.X.Bytes(), publicKeyRaw.Y.Bytes()...)
	return &WalletKeyPair{PrivateKey: privateKey, PublicKey: publicKey}
}

// 生成地址
func (w *WalletKeyPair) GetAddress() string {
	publicHash := HashPubKey(w.PublicKey)
	version := 0x00

	// 21个字节数据
	payload := append([]byte{byte(version)}, publicHash...)

	checksum := CheckSum(payload)

	// 25个字节
	payload = append(payload, checksum...)

	address := base58.Encode(payload)

	return address
}

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

func HashPubKey(pubKey []byte) []byte {
	hash := sha256.Sum256(pubKey)

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

	return publicHash
}

func CheckSum(payload []byte) []byte {
	first := sha256.Sum256(payload)
	second := sha256.Sum256(first[:])

	// 4个字节
	checksum := second[0:4]

	return checksum
}
