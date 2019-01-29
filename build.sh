#!/usr/bin/env bash
go build -o blockchain
rm -rf log
./blockchain > log