package unchained

import (
	"bytes"
	"crypto/sha256"
	"strconv"
	"time"
)

type Block struct {
	Timestamp int64
	Data      []byte
	PrevHash  []byte
	Hash      []byte
}

func (b *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	data := bytes.Join([][]byte{b.PrevHash, b.Data, timestamp}, []byte{})
	hash := sha256.Sum256(data)

	b.Hash = hash[:]
}

func NewBlock(data string, prev_hash []byte) *Block {
	block := &Block{time.Now().Unix(), []byte(data), prev_hash, []byte{}}
	block.SetHash()

	return block
}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}
