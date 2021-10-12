package main

import (
	"fmt"
	"goBlockchain/blockchain"
	"strconv"
)

type CommandLine struct {
	blockchain *blockchain.BlockChain
}

func main() {
	chain := blockchain.InitBlockChain()

	for _, block := range chain.Blocks {
		fmt.Printf("Previous Hash: %x\n", block.PrevHash)
		fmt.Printf("Data in Block: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)

		pow := blockchain.NewProof(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

	}
}
