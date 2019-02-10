#!/usr/bin/env bash
date >> log
date >> log
echo '>>>> 获取yudidi的可用UTXO(余额) ' >> log
#./blockchain getbalance -address Ivan >> log
./blockchain getbalance -address yudidi >> log