#!/usr/bin/env bash
go build -o blockchain
rm -f log
rm -f blockchain.db
date > log
./blockchain >> log