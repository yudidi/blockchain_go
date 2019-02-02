#!/usr/bin/env bash
go build -o blockchain
rm -f log
./blockchain > log