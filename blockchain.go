package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"

	"github.com/boltdb/bolt"
)

const dbFile = "blockchain.db"
const blocksBucket = "blocks"
const genesisCoinbaseData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"

// Blockchain implements interactions with a DB
type Blockchain struct {
	// tip存储最后一个块的哈希。tip有尾部，尖端的意思
	tip []byte
	db  *bolt.DB
}

// BlockchainIterator is used to iterate over blockchain blocks
type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

// MineBlock mines a new block with the provided transactions
func (bc *Blockchain) MineBlock(transactions []*Transaction) {
	var lastHash []byte

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	newBlock := NewBlock(transactions, lastHash)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("l"), newBlock.Hash)
		if err != nil {
			log.Panic(err)
		}

		bc.tip = newBlock.Hash

		return nil
	})
}

// todo 找到支付给address的 所有未使用交易输出 所在的交易。=== 找到address拥有的所有UTXO所在的交易。
//  Q：为什么不能直接返回UTXO集合呢，非要返回其所在的交易。
// address：收款人
// FindUnspentTransactions returns a list of transactions containing unspent outputs
func (bc *Blockchain) FindUnspentTransactions(address string) []Transaction {
	var unspentTXs []Transaction
	//  ydd：存放所有交易已经用掉的交易输出。 key：交易ID  value：该交易已经用过的所有交易输出。即这些交易输出不能再被使用。
	spentTXOs := make(map[string][]int)
	bci := bc.Iterator()
	// ydd：通过迭代器遍历区块链中所有区块：block
	for {
		block := bci.Next()
		// ydd：遍历区块block中的所有交易
		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID) //txID：交易的ID

		Outputs:
			// todo 外循环：遍历该交易的所有输出 +
			//  内循环：检查交易输出是否被使用过，if使用过，跳出内循环，继续外循环，else 没有使用过，检查该交易输出是否是支付给当前用户的，如果是，则被当前交易放入结果集合
			for outIdx, out := range tx.Vout { // outIdx：该交易输出在交易输出数组的索引【spentTXOs这个map的value存放的就是该值】, out：这个交易输出本身
				// Was the output spent?
				if spentTXOs[txID] != nil { // ydd：判断这个交易是否有TXO被使用过。 ==nil，则说明交易输出都没有使用过;  !=nil, 说明有交易输出被使用过
					// ydd：遍历该交易中，被使用过的交易输出
					for _, spentOut := range spentTXOs[txID] { // spentTXOs[txID]：使用 spentOut：
						if spentOut == outIdx { // ydd：该交易输出out，已经被使用过，不放入结果中。
							continue Outputs // ydd：结束内层循环，直接回到外层循环
						}
					}
				}
				// ydd：验证该交易输出是支付给当前用户的。把这个交易输出所在的交易放入 unspentTXs中。
				if out.CanBeUnlockedWith(address) {
					unspentTXs = append(unspentTXs, *tx)
				}
			}

			// todo 把已经被用过的交易输出放入到spentTXOs中。【也就是那些被包含在】
			// ydd：如果交易是coinbase交易，
			if tx.IsCoinbase() == false {
				for _, in := range tx.Vin {
					if in.CanUnlockOutputWith(address) { // 【验证当前用户是否可以解锁UTXO】
						inTxID := hex.EncodeToString(in.Txid)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
					}
				}
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return unspentTXs
}

// FindUTXO finds and returns all unspent transaction outputs
func (bc *Blockchain) FindUTXO(address string) []TXOutput {
	// 1。 扫描区块链 + 找到属于该用户用的所有UTXO
	var UTXOs []TXOutput
	unspentTransactions := bc.FindUnspentTransactions(address)
	for _, tx := range unspentTransactions {
		for _, out := range tx.Vout {
			if out.CanBeUnlockedWith(address) {
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs
}

// FindSpendableOutputs finds and returns unspent outputs to reference in inputs
func (bc *Blockchain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	unspentTXs := bc.FindUnspentTransactions(address)
	accumulated := 0

Work:
	for _, tx := range unspentTXs {
		txID := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.Vout {
			if out.CanBeUnlockedWith(address) && accumulated < amount {
				accumulated += out.Value
				unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	return accumulated, unspentOutputs
}

// Iterator returns a BlockchainIterat
func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.tip, bc.db}

	return bci
}

// Next returns next block starting from the tip
func (i *BlockchainIterator) Next() *Block {
	var block *Block

	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(i.currentHash)
		block = DeserializeBlock(encodedBlock)

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	i.currentHash = block.PrevBlockHash

	return block
}

func dbExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}

// NewBlockchain creates a new Blockchain with genesis Block
func NewBlockchain(address string) *Blockchain {
	if dbExists() == false {
		fmt.Println("No existing blockchain found. Create one first.")
		os.Exit(1)
	}

	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		tip = b.Get([]byte("l"))

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	bc := Blockchain{tip, db}

	return &bc
}

// CreateBlockchain creates a new blockchain DB
func CreateBlockchain(address string) *Blockchain {
	if dbExists() {
		fmt.Println("Blockchain already exists.")
		os.Exit(1)
	}

	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	// 创建coinbase交易，然后
	err = db.Update(func(tx *bolt.Tx) error {
		// 创建coinbase交易
		cbtx := NewCoinbaseTX(address, genesisCoinbaseData)
		// 使用coinbase交易，创建创世区块
		genesis := NewGenesisBlock(cbtx)
		// 在数据库中创建一个名为"blocks"的bucket。可以理解为一张表，表中的每条记录是一个区块。
		b, err := tx.CreateBucket([]byte(blocksBucket))
		if err != nil {
			log.Panic(err)
		}
		// 把创世区块存取数据库
		err = b.Put(genesis.Hash, genesis.Serialize())
		if err != nil {
			log.Panic(err)
		}
		// 存入最后一个区块的hash值
		err = b.Put([]byte("l"), genesis.Hash)
		if err != nil {
			log.Panic(err)
		}
		tip = genesis.Hash

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	bc := Blockchain{tip, db}

	return &bc
}
