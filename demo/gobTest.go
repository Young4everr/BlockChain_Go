package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
)

type Person struct {
	Name string
	Age  uint64
}

func main() {

	Jim := Person{
		Name: "Jim",
		Age:  24,
	}

	var buffer bytes.Buffer

	// 创建编码器
	encoder := gob.NewEncoder(&buffer)
	// 编码
	err := encoder.Encode(&Jim)
	if err != nil {
		log.Panic("编码出错！")
	}
	fmt.Printf("编码后：%x\n", buffer.Bytes())

	var p Person

	// 创建解码器
	decoder := gob.NewDecoder(bytes.NewReader(buffer.Bytes()))
	err = decoder.Decode(&p)
	if err != nil {
		log.Panic("解码失败！")
	}
	fmt.Printf("解码后：%v\n", p)

}
