package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"time"
)

type Block struct {
	Version int64
	PrevBlockHash []byte
	Hash []byte  //简化,应该不带的
	TimeStamp int64
	TargetBits int64
	Nonce int64
	MerKelRoot []byte

	//Data []byte
	Transactions []*Transaction
}

// NewBlock /*初始化区块的各个字段的数据*/
func NewBlock(transactions []*Transaction,prevBlockHash []byte)  *Block {
	block :=  &Block{
		Version: 1,
		PrevBlockHash: prevBlockHash,
		//Hash:
		TimeStamp:  time.Now().Unix(),
		TargetBits: targetBits,
		Nonce:      0,
		MerKelRoot: []byte{},
		//Data: []byte(data)
		Transactions: transactions}
	//block.SetHash()//设置区块的哈希值
	pow := NewProofOfWork(block)
	nonce,hash := pow.Run()
	block.Nonce = nonce
	block.Hash = hash

	return block
}


//func (block *Block)SetHash()  {
//	tmp := [][]byte{
//		IntToByte(block.Version),
//		block.PrevBlockHash,
//		IntToByte(block.TimeStamp),
//		block.MerKelRoot,
//		IntToByte(block.Nonce),
//		block.Data}
//
//	data := bytes.Join(tmp,[]byte{})
//	//对区块进行sha256啊哈希算法，返回值为[32]byte数值，不是切片
//	hash := sha256.Sum256(data)
//	block.Hash=hash[:]
//}

func (block *Block)Serialize()[]byte {
	var buffer bytes.Buffer

	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(block)
	CheckErr("Serialize", err)
	return buffer.Bytes()
}
func  Deserialize (data []byte) *Block {
	if len(data)==0 {
		fmt.Println("Deserialize-----data is empty!")
		os.Exit(1)
	}
	decoder := gob.NewDecoder(bytes.NewReader(data))
	var block Block
	err:=decoder.Decode(&block)
	CheckErr("deserialize",err)
	return &block
}
// NewGenesisBlock /*创建创世块*/
func NewGenesisBlock(coinbase *Transaction)  *Block {
	//return NewBlock("Genesis Block!",[]byte{})
	return  NewBlock([]*Transaction{coinbase},[]byte{})
}