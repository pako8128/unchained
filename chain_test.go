package unchained

import (
	"fmt"
	"strconv"
	"testing"
)

func TestChain(t *testing.T) {
	bc := NewBlockchain()

	bc.AddBlock("Send 1 BTC to Ivan")
	bc.AddBlock("Send 2 more BTC to Ivan")

	iter := bc.Iterator()
	block := iter.Next()
	for string(block.Data) != "Genesis Block" {
		fmt.Printf("Prev. hash: %x\n", block.PrevHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
		block = iter.Next()
	}
}
