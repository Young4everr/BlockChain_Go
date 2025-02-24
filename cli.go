package main

import (
	"fmt"
	"os"
	"strconv"
)

type CLI struct {
	// bc *BlockChain  CLI中不需要保存区块链实例，所有使用区块链的地方自己获取区块链实例
}

const Usage = `
	./blockchain createBlockChain 地址 "创建区块链"
	./blockchain printBlock
	./blockchain getBalance 地址
	./blockchain send FROM TO AMOUNT MINER DATA "转账命令"
	./blockchain createWallet "创建钱包"
`

func (cli *CLI) Run() {

	cmds := os.Args
	if len(cmds) < 2 {
		fmt.Print(Usage)
		os.Exit(1)
	}

	switch cmds[1] {
	case "createBlockChain":
		if len(cmds) != 3 {
			fmt.Print(Usage)
			os.Exit(1)
		}

		fmt.Printf("创建区块链命令被调用，数据：%s\n", cmds[2])
		address := cmds[2]
		cli.CreateBlockChain(address)

	case "printBlock":
		fmt.Printf("打印区块链命令被调用\n")
		cli.PrintBlock()

	case "getBalance":
		if len(cmds) != 3 {
			fmt.Print(Usage)
			os.Exit(1)
		}
		fmt.Printf("获取余额命令被调用\n")
		address := cmds[2]
		cli.GetMyBalance(address)

	case "send":
		if len(cmds) != 7 {
			fmt.Print(Usage)
			os.Exit(1)
		}
		fmt.Printf("转账命令被调用\n")
		from := cmds[2]
		to := cmds[3]
		amount, _ := strconv.ParseFloat(cmds[4], 64)
		miner := cmds[5]
		data := cmds[6]
		cli.Send(from, to, amount, miner, data)

	case "createWallet":
		fmt.Printf("创建钱包命令被调用\n")
		cli.CreateWallet()

	default:
		fmt.Printf("无效命令，请检查！")
		fmt.Print(Usage)
		os.Exit(1)
	}
}
