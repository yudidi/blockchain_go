#!/usr/bin/env bash
chmod +x run-createblockchain.sh
run1-createblockchain.sh
date >> log
date >> log
echo '>>>> 发送币 ' >> log
./blockchain send -from yudidi -to Pedro -amount 6 >> log