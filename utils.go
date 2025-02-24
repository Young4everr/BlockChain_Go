package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"os"
)

//工具函数库

func UintToByte(num uint64) []byte {
	// 使用binary.Write来进行编码
	var buffer bytes.Buffer

	// 编码要进行错误检查，一定要做
	err := binary.Write(&buffer, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buffer.Bytes()
}

// 判断文件是否存在
func IsFileExist(fileName string) bool {
	// 使用os.Stat来判断
	_, err := os.Stat(fileName)

	return !os.IsNotExist(err)
}
