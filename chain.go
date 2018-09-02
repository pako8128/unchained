package unchained

import "github.com/boltdb/bolt"

const blocksBucket = "blocks"
const dbFile = "./blocks.db"

const genesisCoinbaseData = "Hello, world!"

type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

func (bc *Blockchain) AddBlock(transactions []*Transaction) {
	var lastHash []byte

	bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))

		return nil
	})

	block := NewBlock(transactions, lastHash)

	bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		b.Put(block.Hash, block.Serialize())
		b.Put([]byte("l"), block.Hash)
		bc.tip = block.Hash

		return nil
	})
}

func NewBlockchain(address string) *Blockchain {
	var tip []byte
	db, _ := bolt.Open(dbFile, 0600, nil)

	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		if b == nil {
			cbtx := NewCoinbaseTX(address, genesisCoinbaseData)
			genesis := NewGenesisBlock(cbtx)
			b, _ = tx.CreateBucket([]byte(blocksBucket))
			b.Put(genesis.Hash, genesis.Serialize())
			b.Put([]byte("l"), genesis.Hash)
			tip = genesis.Hash
		} else {
			tip = b.Get([]byte("l"))
		}

		return nil
	})

	bc := Blockchain{tip, db}

	return &bc
}

func (bc *Blockchain) Close() {
	bc.db.Close()
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	bciter := BlockchainIterator{bc.tip, bc.db}

	return &bciter
}

type BlockchainIterator struct {
	current []byte
	db      *bolt.DB
}

func (iter *BlockchainIterator) Next() *Block {
	var block *Block

	iter.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		raw := b.Get(iter.current)
		block = Deserialize(raw)

		return nil
	})

	iter.current = block.PrevHash

	return block
}
