package unchained

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

const subsidy = 10

type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}

type TXOutput struct {
	Value        int
	ScriptPubKey string
}

type TXInput struct {
	Txid      []byte
	Vout      int
	ScriptSig string
}

func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		fmt.Println("Reward to %s", to)
	}

	txin := TXInput{[]byte{}, -1, data}
	txout := TXOutput{subsidy, to}

	tx := Transaction{nil, []TXInput{txin}, []TXOutput{txout}}
	tx.setID()

	return &tx
}

func NewUTXOTransaction(from, to string, amount int, bc *Blockchain) *Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	acc, validOutputs := bc.FindSpendableOutputs(from, amount)

	if acc < amount {
		log.Panic("ERROR: Not enough funds")
	}

	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)

		if err != nil {
			log.Panic("Decoding Error")
		}

		for _, out := range outs {
			input := TXInput{txID, out, from}
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, TXOutput{amount, to})
	if acc > amount {
		outputs = append(outputs, TXOutput{acc - amount, from})
	}

	tx := Transaction{nil, inputs, outputs}
	tx.setID()

	return &tx
}

func (tx *Transaction) setID() {
	var encoded bytes.Buffer
	var hash [32]byte

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)

	if err != nil {
		log.Panic("Encoding Error.")
	}

	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

func (tx Transaction) isCoinbase() bool {
	return len(tx.Vin) == 0 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
}

func (in *TXInput) CanUnlockOutputWith(data string) bool {
	return in.ScriptSig == data
}

func (out *TXOutput) CanBeUnlockedWith(data string) bool {
	return out.ScriptPubKey == data
}
