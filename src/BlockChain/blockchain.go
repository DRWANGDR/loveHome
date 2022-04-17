package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"os"
)

const dbfile = "blockChainDb.db"
const blockBucket = "block"
const lasthash = "lastHash"
const genesisBlockInfo = "The Times 03/Jan/2009 Chancellor on brink of second bailout forbanks."

// BlockChain /*构建区块链结构，使用数组来存储所有区块*/
type BlockChain struct {
	//blocks []*Block
	db       *bolt.DB
	lastHash []byte
}

// NewBlockChain /*创建区块链实例，并添加创世块*/
func NewBlockChain(address string) *BlockChain {
	//return &BlockChain{[]*Block{NewGenesisBlock()}}
	if IsBlockChainExit() {
		fmt.Println("Block Chain already exist!")
		os.Exit(1)
	}

	db, err := bolt.Open(dbfile, 0600, nil)
	CheckErr("NewBlockChain1", err)

	var lastHash []byte
	//db.View()
	err0 := db.Update(func(tx *bolt.Tx) error {
		coinbase := NewCoinbaseTX(address, genesisBlockInfo)
		genesis := NewGenesisBlock(coinbase)
		bucket, err := tx.CreateBucket([]byte(blockBucket))
		CheckErr("NewBlockChain2", err)

		err = bucket.Put(genesis.Hash, genesis.Serialize())
		CheckErr("NewBlockChain3", err)
		err = bucket.Put([]byte(lasthash), genesis.Hash)
		CheckErr("NewBlockChain4", err)
		lastHash = genesis.Hash

		return nil
	})
	CheckErr("err0:  ", err0)
	return &BlockChain{db, lastHash}
}

func GetBlockChainHandler() *BlockChain {
	if !IsBlockChainExit() {
		fmt.Println("Block Chain not exist , pls create it!")
		os.Exit(1)
	}

	db, err := bolt.Open(dbfile, 0600, nil)
	CheckErr("GetBlockChainHandler err1:", err)

	var lastHash []byte
	//db.View()
	err0 := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket != nil {
			//读取
			lastHash = bucket.Get([]byte(lasthash))
		}
		return nil
	})
	CheckErr("GetBlockChainHandler err2:  ", err0)
	return &BlockChain{db, lastHash}
}

// AddBlock /*添加区块操作*/
//func (bc *BlockChain)AddBlock(data string)  {
func (bc *BlockChain) AddBlock(transactions []*Transaction) {
	var prevBlockHash []byte
	err := bc.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		lastHash := bucket.Get([]byte(lasthash))
		prevBlockHash = lastHash
		return nil
	})
	CheckErr("AddBlock", err)

	block := NewBlock(transactions, prevBlockHash)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		err := bucket.Put(block.Hash, block.Serialize())
		CheckErr("AddBlock2", err)
		err = bucket.Put([]byte(lasthash), block.Hash)
		CheckErr("AddBlock3", err)
		bc.lastHash = block.Hash
		return nil
	})
	CheckErr("AddBlock4", err)
}

type BlockChainIterator struct {
	db          *bolt.DB
	currentHash []byte
}

func (bc *BlockChain) Iterator() *BlockChainIterator {
	return &BlockChainIterator{bc.db, bc.lastHash}
}

func (it *BlockChainIterator) Next() *Block {
	var block *Block
	err := it.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			os.Exit(1)
		}
		blockTmp := bucket.Get(it.currentHash)

		block = Deserialize(blockTmp)

		it.currentHash = block.PrevBlockHash
		return nil
	})
	CheckErr("Next", err)
	return block
}

func (bc *BlockChain) FindUnspendTransaction(address string) []Transaction {
	var transactions []Transaction
	var spentUTXOs = make(map[string] /*交易的txid*/ []int64)
	bci := bc.Iterator()
	for {
		block := bci.Next()

		for _, currTx := range block.Transactions {
			txid := string(currTx.TXID)

			//遍历当前交易的inputs，找到当前地址消耗的utxos
			for _, input := range currTx.TXInputs {
				if currTx.IsCoinbase() == false {
					if input.CanUnlockUTXOByAddress(address) {
						//map[txid] = []int64/*output的index*/
						spentUTXOs[string(input.Txid)] = append(spentUTXOs[txid], input.ReferOutputIndex)
					}
				}
			}
		LABEL1:
			//遍历当前交易的outputs，通过output的解锁条件，确定满足条件的交易
			for outputIndex, output := range currTx.TXOutputs {
				if spentUTXOs[txid] != nil {
					for _, usedIndex := range spentUTXOs[txid] {
						if int64(outputIndex) == usedIndex {
							fmt.Println("used,no need to add again!!!")
							continue LABEL1
						}
					}
				}

				if output.CanBeUnlockedByAddress(address) {
					transactions = append(transactions, *currTx)
				}
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
	return transactions
}

func IsBlockChainExit() bool {
	_, err := os.Stat(dbfile)
	if os.IsNotExist(err) {
		return false
	}
	return true
	//return !os.IsExist(err)
}

func (bc *BlockChain) FindSuitableUTXOs(address string, amount float64) (float64, map[string][]int64) {
	txs := bc.FindUnspendTransaction(address)
	var countTotal float64
	var container = make(map[string][]int64)

LABEL2:
	for _, tx := range txs {
		for index, output := range tx.TXOutputs {
			fmt.Println("开始遍历tx.TXOutputs")
			if countTotal < amount {
				fmt.Println(countTotal, "<<<", amount)
				if output.CanBeUnlockedByAddress(address) {
					countTotal += output.Value
					fmt.Println("output.Value:::", output.Value)
					container[string(tx.TXID)] = append(container[string(tx.TXID)], int64(index))
				} else {
					fmt.Println("break发生")
					break LABEL2
				}
			}
		}
	}
	fmt.Println("countTotal::::", countTotal)
	return countTotal, container
}
