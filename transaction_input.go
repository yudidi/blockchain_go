package main

import "bytes"

// TXInput represents a transaction input
type TXInput struct {
	Txid      []byte
	Vout      int
	Signature []byte
	PubKey    []byte
}

// UsesKey checks whether the address initiated the transaction
// todo 判断是否是pubKeyHash使用了TXInput中的哪些交易输入。
// 判断该pubKeyHash是否==该交易输入的解锁脚本中的公钥哈希。如果相等，则说明这个交易输入中的UTXO，是被pubKeyHash使用的。
func (in *TXInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := HashPubKey(in.PubKey)

	return bytes.Compare(lockingHash, pubKeyHash) == 0
}
