package unchained

import (
	"bytes"
	"encoding/gob"
	"time"
)

type Block struct {
	Timestamp int64
	Data      []byte
	PrevHash  []byte
	Hash      []byte
	Nonce     int
}

func NewBlock(data string, prev_hash []byte) *Block {
	block := &Block{time.Now().Unix(), []byte(data), prev_hash, []byte{}, 0}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	encoder.Encode(b)

	return result.Bytes()
}

func Deserialize(d []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	decoder.Decode(&block)

	return &block
}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}
