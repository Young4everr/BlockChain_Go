package main

func main() {
	bc := NewBlockChain("The miner")
	defer bc.db.Close()
	cli := CLI{bc}
	cli.Run()
}
