#!/usr/bin/env bash
go build -o blockchain
rm -f log
date > log
./blockchain >> log