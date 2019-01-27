package main

import (
	"crypto/sha256"
	"fmt"
	"math/big"
	"testing"
)

// 左移
func TestLsh(t *testing.T) {
	data1 := []byte("I like donuts")
	data2 := []byte("I like donutsca07ca")
	targetBits := 24
	maxTargetHash := big.NewInt(1)
	maxTargetHash.Lsh(maxTargetHash, uint(256-targetBits))
	fmt.Printf("%x\n", sha256.Sum256(data1)) //返回值是 32个字节 | 1个字节 == 8位 | 256位
	fmt.Printf("%x\n", maxTargetHash)        //%x 小写的十六进制 | 默认左对齐
	fmt.Printf("%32x\n", maxTargetHash)      //%x 小写的十六进制 |
	fmt.Printf("%64x\n", maxTargetHash)      //%x 小写的十六进制 | 1111->F 4位一个十六进制数 | 256位 -> 64个十六进制数
	fmt.Printf("%65x\n", maxTargetHash)      //todo 对齐方式？？？。 %nx： n是多少，就先产生多少个空格，然后开始从右往左填充数字。
	fmt.Printf("%x\n", sha256.Sum256(data2))

	fmt.Printf("%v\n", maxTargetHash)

}
