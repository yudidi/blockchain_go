package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

// subsidy 是奖励的数额。
// 在比特币中，实际并没有存储这个数字，而是基于区块总数进行计算而得：区块总数除以 210000 就是 subsidy. 挖出创世块的奖励是 50 BTC，每挖出 210000 个块后，奖励减半。
// 在我们的实现中，这个奖励值将会是一个常量（至少目前是）。
const subsidy = 10

// Transaction represents a Bitcoin transaction
type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}

// IsCoinbase checks whether the transaction is coinbase
func (tx Transaction) IsCoinbase() bool {
	// todo:
	//  Q：为什么coinbase交易只能有一个交易输入
	//  A：coinbase交易的交易输入没有UTXO，但是有一个锁定脚本。
	//  只有一个交易输入 && 没有UTXO(所以交易输入对应的2个属性没有被赋值)
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
}

// SetID sets ID of a transaction
func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

// TODO
//========================
// UTXO: 表示某交易中的未使用的某个输出
type UTXO struct {
	// 包含该UTXO的那个交易的ID
	Txid []byte
	// 该交易输出在该交易所有交易输出中的索引
	Index int
}

// 交易输入中包含一个UTXO和一个解锁脚本，所以结构体应该是这样
type TXIn struct {
	Utxo            UTXO
	unlockingScript string
}

//========================

// TXInput represents a transaction input
type TXInput struct {
	// TODO FAQ Q： 一个交易输入结构体TXInput应该包含1个还是多个UTXO呢
	// Txid 和 Vout 唯一标识一个UTXO
	Txid      []byte
	Vout      int
	ScriptSig string
}

// TXOutput represents a transaction output
type TXOutput struct {
	Value        int
	ScriptPubKey string
}

// CanUnlockOutputWith checks whether the address initiated the transaction
func (in *TXInput) CanUnlockOutputWith(unlockingData string) bool {
	return in.ScriptSig == unlockingData
}

// CanBeUnlockedWith checks if the output can be unlocked with the provided data
//todo P2PKH
// 验证这个交易输出是否是支付给unlockingData
// 比特币交易中：解锁脚本(签名+公钥) || 我们的简化交易中：解锁脚本(收款人)
// 对于P2PKH交易实现的理解，付款给公钥 简化为 付款给一个字符串
func (out *TXOutput) CanBeUnlockedWith(unlockingData string) bool {
	return out.ScriptPubKey == unlockingData //todo unlockingData -->  payee收款人。公钥就是比特币交易中的收款人
}

// todo：因为还没有实现奖励制度，所以目前只有创建创世区块时，会涉及到coinbase交易。
//   如果之后加入了奖励，那么在每个新块产生时，都会有奖励，也就有coinbase。
// todo 矿工挖出新块时，它会向新块中添加一个coinbase交易。这个交易 只有唯一到一个交易输入，并且这个交易输入没有UTXO，只有锁定脚本。
// NewCoinbaseTX creates a new coinbase transaction
func NewCoinbaseTX(to, lockingScript string) *Transaction {
	if lockingScript == "" {
		lockingScript = fmt.Sprintf("Reward to '%s'", to)
	}

	txin := TXInput{[]byte{}, -1, lockingScript} //todo:  coinbase交易没有输入UTXO，只需要传递一个锁定脚本即可
	txout := TXOutput{subsidy, to}
	tx := Transaction{nil, []TXInput{txin}, []TXOutput{txout}}
	tx.SetID()

	return &tx
}

// NewUTXOTransaction creates a new transaction
func NewUTXOTransaction(from, to string, amount int, bc *Blockchain) *Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	acc, validOutputs := bc.FindSpendableOutputs(from, amount)

	if acc < amount {
		log.Panic("ERROR: Not enough funds")
	}

	// Build a list of inputs
	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}

		for _, out := range outs {
			input := TXInput{txID, out, from}
			inputs = append(inputs, input)
		}
	}

	// Build a list of outputs
	outputs = append(outputs, TXOutput{amount, to})
	if acc > amount {
		outputs = append(outputs, TXOutput{acc - amount, from}) // a change
	}

	tx := Transaction{nil, inputs, outputs}
	tx.SetID()

	return &tx
}
