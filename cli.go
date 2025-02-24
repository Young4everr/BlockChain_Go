package main

import (
	"fmt"
	"os"
	"strconv"
)

type CLI struct {
	bc *BlockChain
}

const Usage = `
	./blockchain addBlock "xxx"
	./blockchain printBlock
	./blockchain getBalance 地址
	./blockchain send FROM TO AMOUNT MINER DATA "转账命令"
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

	case "send":
		if len(cmds) != 7 {
			fmt.Printf(Usage)
			os.Exit(1)
		}
		fmt.Printf("转账命令被调用\n")
		from := cmds[2]
		to := cmds[3]
		amount, _ := strconv.ParseFloat(cmds[4], 64)
		miner := cmds[5]
		data := cmds[6]
		cli.Send(from, to, amount, miner, data)

	default:
		fmt.Printf("无效命令，请检查！")
		fmt.Printf(Usage)
		os.Exit(1)
	}
}
