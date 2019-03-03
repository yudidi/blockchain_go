# 创建2个钱包Alice/Bob，也就是产生2组密钥对，2个用户对意思。
# 然后使用其中一个地址Alice创建初始区块链，Alice获得挖矿奖励10个币
# 然后Alice转账1个币给Bob。
# 最后分别查看Alice和Bob对余额。

# 1. blockchain_go createwallet

# 2. blockchain_go createwallet

# 3. blockchain_go createblockchain -address 13S3ChRLoaN3J2K4ozYKyDzBvdxug4WBs7

# 4. blockchain_go send -from 13S3ChRLoaN3J2K4ozYKyDzBvdxug4WBs7 -to 14YZ8TZQZxiCWKXJD1t4Bey9Cepehxxtx9 -amount 6

# 5. blockchain_go getbalance -address 13S3ChRLoaN3J2K4ozYKyDzBvdxug4WBs7
Balance of '1AmVdDvvQ977oVCpUqz7zAPUEiXKrX5avR': 9

# 6. blockchain_go getbalance -address 13S3ChRLoaN3J2K4ozYKyDzBvdxug4WBs7
Balance of '1NE86r4Esjf53EL7fR86CsfTZpNN42Sfab': 1