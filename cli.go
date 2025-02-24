package main

import (
	"fmt"
	"os"
)

type CLI struct {
	bc *BlockChain
}

const Usage = `
	./blockchain addBlock "xxx"
	./blockchain printBlock
	./blockchain getBalance 地址
`

func (cli *CLI) Run() {

	cmds := os.Args
	if len(cmds) < 2 {
		fmt.Printf(Usage)
		os.Exit(1)
	}

	switch cmds[1] {
	case "addBlock":
		if len(cmds) != 3 {
			fmt.Printf(Usage)
			os.Exit(1)
		}

		fmt.Printf("添加区块命令被调用，数据：%s\n", cmds[2])
		// TODO
		// data := cmds[2]
		// cli.addBlock(data)

	case "printBlock":
		fmt.Printf("打印区块链命令被调用\n")
		cli.printBlock()

	case "getBalance":
		if len(cmds) != 3 {
			fmt.Printf(Usage)
			os.Exit(1)
		}
		fmt.Printf("获取余额命令被调用\n")
		address := cmds[2]
		cli.bc.GetMyBalance(address)

	default:
		fmt.Printf("无效命令，请检查！")
		os.Exit(1)
	}
}
