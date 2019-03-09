package main

import "bytes"

// TXInput represents a transaction input
type TXInput struct {
	PreTxid []byte
	// 该交易输入使用的pre交易输出的序号
	IndexOfPreOutput int
	Signature        []byte
	PubKey           []byte
}

// UsesKey checks whether the address initiated the transaction
func (in *TXInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := HashPubKey(in.PubKey)

	return bytes.Compare(lockingHash, pubKeyHash) == 0
}
