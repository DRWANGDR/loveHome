package main

import (
	"fmt"
	"github.com/boltdb/bolt"
)

func (cli *CLI) CreateChain(address string) {
	bc := NewBlockChain(address)
	defer func(db *bolt.DB) {
		err := db.Close()
		if err != nil {
			CheckErr("数据库关闭失败！！！", err)
		}
	}(bc.db)
	fmt.Println("create blockchain successfully.")
}

/*
func (cli *CLI) AddBlock(data string)  {
	cli.bc.AddBlock(data)
	fmt.Println("AddBlock Success!")
}
*/

func (cli *CLI) PrintChain() {
	bc := GetBlockChainHandler()
	it := bc.Iterator()
	for {
		block := it.Next() //取回当前hash指向的block，将hash值前移

		fmt.Printf("Transaction: %s\n", block.Transactions)
		fmt.Println("Version:", block.Version)
		fmt.Printf("PrevBlockHash: %x\n", block.PrevBlockHash)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Printf("TimeStamp: %d\n", block.TimeStamp)
		fmt.Printf("MerKelRoot: %x\n", block.MerKelRoot)
		fmt.Printf("Nonce: %d\n", block.Nonce)
		pow := NewProofOfWork(block)
		fmt.Printf("IsValid: %v\n", pow.IsValid())
		fmt.Println("---------------------------------------")
		if len(block.PrevBlockHash) == 0 {
			break
		}

	}
}

func (cli *CLI) GetBalance(address string) {
	bc := GetBlockChainHandler()
	defer func(db *bolt.DB) {
		err := db.Close()
		if err != nil {
			fmt.Println("Close occur error!!")
		}
	}(bc.db)
	var total float64

	utxos := bc.FindUTXOs(address)
	for _, utxo := range utxos {
		total += utxo.Value
	}
	fmt.Println("The balance of  ", address, "  is  ", total)
}

func (cli *CLI) Send(from, to string, amount float64) {
	bc := GetBlockChainHandler()
	tx := NewTransaction(from, to, amount, bc)
	bc.AddBlock([]*Transaction{tx})
	fmt.Println("send successfully!")
}
