package main


import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"os"
)

const reward float64 = 12.5

type Transaction struct {
	TXID []byte
	TXInputs []Input
	TXOutputs []Output
}

type Input struct {
	Txid []byte
	ReferOutputIndex int64
	UnlockScript string
	//ScriptSig
}

type Output struct {
	Value float64
	LockScript string
	//ScriptPubKey
}

func (tx *Transaction)SetTXID() {
	//data := bytes.Join([][]byte{},[]byte{})
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(tx)
	CheckErr("SetTXID occur error:::",err)
	hash := sha256.Sum256(buffer.Bytes())
	tx.TXID = hash[:]
}

func (input *Input)CanUnlockUTXOByAddress(unlockdata string) bool {
	return  input.UnlockScript == unlockdata
}

func (output *Output)CanBeUnlockedByAddress(unlockdata string) bool {
	fmt.Println(output.LockScript,"+++++",unlockdata)
	return output.LockScript == unlockdata
}

func NewCoinbaseTX(address string,data string) *Transaction {
	if data == "" {
		fmt.Println(data,"Current Reward is : ", reward)
	}

	input := Input{nil,-1,data}
	output:= Output{reward,address}

	tx := Transaction{nil,[]Input{input},[]Output{output}}
	tx.SetTXID()
	return &tx
}

func NewTransaction(from,to string,amount float64,bc *BlockChain) *Transaction {
	counted,container := bc.FindSuitableUTXOs(from,amount)
	if counted < amount {
		fmt.Println(counted,"::::::",amount)
		fmt.Println("No Enough Founds!!!")
		os.Exit(1)
	}

	var inputs []Input
	var outputs []Output

	for txid,outputIndexs := range container {
		for _,index := range outputIndexs{
			input:= Input{[]byte(txid),index,from}
			inputs = append(inputs,input)
		}
	}

	output := Output{amount,to}
	outputs =append(outputs,output)

	//找零
	if counted >amount {
		outputs =append(outputs, Output{counted - amount,from})
	}



	tx := Transaction{nil,inputs,outputs}
	tx.SetTXID()
	return &tx
}


func (tx *Transaction)IsCoinbase() bool {
	if len(tx.TXOutputs) == 1 {
		if tx.TXInputs[0].Txid == nil && tx.TXInputs[0].ReferOutputIndex == -1 {
			return true
		}
	}
	return false
}



func (bc *BlockChain)FindUTXOs(address string) []Output {
	var outputs []Output
	txs:=bc.FindUnspendTransaction(address)
	for _,tx := range txs{
		for _,output := range tx.TXOutputs {
			if output.CanBeUnlockedByAddress(address) {
				outputs =append(outputs,output)
			}else {
				fmt.Println("FindUTXOs CanBeUnlockedByAddress err!!! ")
			}
		}
	}
	return outputs
}