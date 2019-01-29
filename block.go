package main

import (
	"time"
)

// Block keeps block headers
type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	//todo 共识属性
	Hash  []byte
	Nonce int
}

// todo 使用pow把仅有数据属性的区块，加工为一个拥有完备属性(共识属性)的合法区块。
// NewBlock creates and returns Block
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}, 0}
	// todo block：打包好的区块，进行记账权的争夺。挖矿
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()
	//todo 通过pow，暴力尝试出一个nonce值，把nonce和对应的hash放入块中，使得其他人可以验证，
	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

// NewGenesisBlock creates and returns genesis Block
func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}
