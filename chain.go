package unchained

import (
	"encoding/hex"
	"fmt"
	"github.com/boltdb/bolt"
	"os"
)

const blocksBucket = "blocks"
const dbFile = "./blocks.db"

const genesisCoinbaseData = "Hello, world!"

type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

func dbExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}

func (bc *Blockchain) MineBlock(transactions []*Transaction) {
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

func (bc *Blockchain) FindUnspentTransactions(address string) []Transaction {
	var unspentTXs []Transaction
	spentTXOs := make(map[string][]int)
	bci := bc.Iterator()

	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Vout {
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}

				if out.CanBeUnlockedWith(address) {
					unspentTXs = append(unspentTXs, *tx)
				}
			}

			if tx.isCoinbase() == false {
				for _, in := range tx.Vin {
					if in.CanUnlockOutputWith(address) {
						inTxID := hex.EncodeToString(in.Txid)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
					}
				}
			}
		}

		if len(block.PrevHash) == 0 {
			break
		}
	}

	return unspentTXs
}

func (bc *Blockchain) FindUTXOs(address string) []TXOutput {
	var UTXOs []TXOutput
	unspentTXs := bc.FindUnspentTransactions(address)

	for _, tx := range unspentTXs {
		for _, out := range tx.Vout {
			if out.CanBeUnlockedWith(address) {
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs
}

func (bc *Blockchain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	unspentTXs := bc.FindUnspentTransactions(address)
	accumulated := 0

Work:
	for _, tx := range unspentTXs {
		txID := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.Vout {
			if out.CanBeUnlockedWith(address) && accumulated < amount {
				accumulated += out.Value
				unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	return accumulated, unspentOutputs
}

func NewBlockchain(address string) *Blockchain {
	if dbExists() == false {
		fmt.Println("No existing Blockchain! Create one first.")
		os.Exit(1)
	}

	var tip []byte
	db, _ := bolt.Open(dbFile, 0600, nil)

	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		tip = b.Get([]byte("l"))

		return nil
	})

	bc := Blockchain{tip, db}

	return &bc
}

func CreateBlockchain(address string) *Blockchain {
	if dbExists() {
		fmt.Println("There already is a DB!")
		os.Exit(1)
	}

	var tip []byte

	db, _ := bolt.Open(dbFile, 0600, nil)
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		cbtx := NewCoinbaseTX(address, genesisCoinbaseData)
		genesis := NewGenesisBlock(cbtx)
		b, _ = tx.CreateBucket([]byte(blocksBucket))
		b.Put(genesis.Hash, genesis.Serialize())
		b.Put([]byte("l"), genesis.Hash)
		tip = genesis.Hash

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
