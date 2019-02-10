#!/usr/bin/env bash
chmod +x run-createblockchain.sh
run1-createblockchain.sh
date >> log
date >> log
echo '>>>> 使用多个简易输出，完成 发送币 #通过2个帐户分别发送2个币给Alice，Alice共计有4个币，然后Alice发送3个币给Bob ' >> log
./blockchain send -from yudidi -to Pedro -amount 6 >> log