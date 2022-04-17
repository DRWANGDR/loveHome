package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

const targetBits = 24

type ProofOfWork struct {
	block     *Block
	targetBit *big.Int
}

func NewProofOfWork(block *Block) *ProofOfWork {
	var IntTarget = big.NewInt(1) //00000000000001
	IntTarget.Lsh(IntTarget, uint(256-targetBits))

	return &ProofOfWork{block, IntTarget}
}

func (pow *ProofOfWork) PrepareRawData(nonce int64) []byte {
	block := pow.block
	tmp := [][]byte{
		IntToByte(block.Version),
		block.PrevBlockHash,
		IntToByte(block.TimeStamp),
		block.MerKelRoot,
		IntToByte(nonce),
		IntToByte(targetBits)}
	//block.Data
	//block.Transactions//TODO

	data := bytes.Join(tmp, []byte{})
	return data
}

func (pow *ProofOfWork) Run() (int64, []byte) {

	var nonce int64
	var hash [32]byte
	var HashInt big.Int

	fmt.Println(" Begin mining...(Run)")
	fmt.Printf("target hash: %x\n", pow.targetBit.Bytes())
	for nonce < math.MaxInt64 {
		data := pow.PrepareRawData(nonce)
		hash = sha256.Sum256(data)

		HashInt.SetBytes(hash[:])

		if HashInt.Cmp(pow.targetBit) == -1 {
			fmt.Printf("Found Hash: %x\n", hash)
			break
		} else {
			//fmt.Println("current Hash is :",hash)
			nonce++
		}
	}
	return nonce, hash[:]
}

func (pow *ProofOfWork) IsValid() bool {
	data := pow.PrepareRawData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	var IntHash big.Int
	IntHash.SetBytes(hash[:])
	return IntHash.Cmp(pow.targetBit) == -1
}
