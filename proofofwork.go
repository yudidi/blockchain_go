package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
	"time"
)

var (
	maxNonce = math.MaxInt64
)

// 挖矿难度值 叫difficulty更直观
const targetBits = 4 //24

//todo 一个区块共5个属性
// 1. 传入一个区块(含3个属性)，然后计算出剩下的2个属性。
// 已有3个属性值【数据data，时间戳timestamp，副区块prehash】
// 计算2个属性值：然后由pow暴力计算一个nonce及满足条件的hash
// 2. 参与hash计算的属性是3个 + 难度值和nonce值

// ProofOfWork represents a proof-of-work
type ProofOfWork struct {
	block *Block //todo
	//当前难度值下,合法hash的上界。叫maxValidHash更直观
	target *big.Int
}

// NewProofOfWork builds and returns a ProofOfWork
func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	pow := &ProofOfWork{b, target}

	return pow
}

// todo 参与hash计算的数据有：
func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.Data,
			IntToHex(pow.block.Timestamp),
			IntToHex(int64(targetBits)),
			IntToHex(int64(nonce)),
		},
		[]byte{},
	)

	return data
}

// todo 切片到定义和初始化 http://www.runoob.com/go/go-slice.html
// Run performs a proof-of-work
func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte //todo 你可以声明一个未指定大小的数组来定义切片：
	nonce := 0

	fmt.Printf("Mining the block containing \"%s\"\n", pow.block.Data)
	for nonce < maxNonce {
		data := pow.prepareData(nonce)

		hash = sha256.Sum256(data)
		fmt.Printf("当前得到的hash %x \r", hash)
		fmt.Printf("此时的nonce值 %v \r ", nonce)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 {
			fmt.Printf("找到符合条件的hash值 %x", hash)
			time.Sleep(1000)
			break
		} else {
			nonce++
		}
	}
	fmt.Print("\n\n")

	return nonce, hash[:] //todo 返回切片的第一个元素到最后一个元素(所有元素)
}

// Validate validates block's PoW
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.target) == -1

	return isValid
}
