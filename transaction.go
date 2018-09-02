package unchained

import "fmt"

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

func (tx *Transaction) setID() {

}
