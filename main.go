package main

func main() {
	bc := NewBlockChain("miner01")
	defer bc.db.Close()
	cli := CLI{bc}
	cli.Run()
}
