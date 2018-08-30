package unchained

type Blockchain struct {
	blocks []*Block
}

func (bc *Blockchain) AddBlock(data string) {
	prev_block := bc.blocks[len(bc.blocks)-1]
	new_block := NewBlock(data, prev_block.Hash)
	bc.blocks = append(bc.blocks, new_block)
}

func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}}
}
