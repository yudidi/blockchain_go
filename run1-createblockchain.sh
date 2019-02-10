#!/usr/bin/env bash
go build -o blockchain
rm -f log & echo 'delete log success'
rm -f blockchain.db & echo 'delete db success'
date >> log
date >> log
echo '>>>> Start shell !!! 创建一个区块链(创建创世区块) 暂时用"字符串"作为简易地址，接收coinbase的挖矿奖励 ' >> log
./blockchain createblockchain -address yudidi >> log