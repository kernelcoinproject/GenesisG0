# GenesisG0

Port of https://github.com/lhartikk/GenesisH0 to golang with parallelism and time estimates.

Same syntax, only sha and scrypt are supported. No x11/x13/x15.

## Examples

Bitcoin
```
./genesis -z "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks" -n 2083236893 -t 1231006505

04ffff001d0104455468652054696d65732030332f4a616e2f32303039204368616e63656c6c6f72206f6e206272696e6b206f66207365636f6e64206261696c6f757420666f722062616e6b73
algorithm: SHA256
merkle hash: 4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b
pszTimestamp: The Times 03/Jan/2009 Chancellor on brink of second bailout for banks
pubkey: 04678afdb0fe5548271967f1a67130b7105cd6a828e03909a67962e0ea1f61deb649f6bc3f4cef38c4f35504e51ec112de5c384df7ba0b8d578a4c702b6bf11d5f
time: 1231006505
bits: 0x1d00ffff
Searching for genesis hash with 8 workers..

genesis hash found!
nonce: 2083236893
genesis hash: 000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f
```
Litecoin
```
./genesis -a scrypt -z "NY Times 05/Oct/2011 Steve Jobs, Apple’s Visionary, Dies at 56" -p "040184710fa689ad5023690c80f3a49c8f13f8d45b8c857fbcbc8bc4a8e4d3eb4b10f4d4604fa08dce601aaf0f470216fe1b51850b4acf21b179c45070ac7b03a9" -t 1317972665 -n 2084524493

04ffff001d0104404e592054696d65732030352f4f63742f32303131205374657665204a6f62732c204170706c65e280997320566973696f6e6172792c2044696573206174203536
algorithm: scrypt
merkle hash: 97ddfbbae6be97fd6cdf3e7ca13232a3afff2353e29badfab7f73011edd4ced9
pszTimestamp: NY Times 05/Oct/2011 Steve Jobs, Apple’s Visionary, Dies at 56
pubkey: 040184710fa689ad5023690c80f3a49c8f13f8d45b8c857fbcbc8bc4a8e4d3eb4b10f4d4604fa08dce601aaf0f470216fe1b51850b4acf21b179c45070ac7b03a9
time: 1317972665
bits: 0x1e0ffff0
Searching for genesis hash with 8 workers..
8.79 kh/s, elapsed: 26s, estimate: 1m59s                    
genesis hash found!
nonce: 246683
genesis hash: b723793cef58bba25e46c2aecf61da1f7f14df80a4dcafc71c21a20592f6820d
```

## Syntax

```
Usage of ./genesis:
  -a string
    	short flag for algorithm
  -algorithm string
    	the PoW algorithm: [SHA256|scrypt] (default "SHA256")
  -b int
    	the target in compact representation, associated to a difficulty of 1
  -c int
    	number of CPU cores to use for mining (default auto)
  -n uint
    	short flag for nonce
  -nonce uint
    	the first value of the nonce that will be incremented when searching the genesis hash
  -p string
    	the pubkey found in the output script (default "04678afdb0fe5548271967f1a67130b7105cd6a828e03909a67962e0ea1f61deb649f6bc3f4cef38c4f35504e51ec112de5c384df7ba0b8d578a4c702b6bf11d5f")
  -t int
    	short flag for time
  -time int
    	the (unix) time when the genesisblock is created (default 1763286915)
  -timestamp string
    	the pszTimestamp found in the coinbase of the genesisblock (default "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks")
  -v int
    	the value in coins for the output, full value (default 5000000000)
  -z string
    	short flag for timestamp

```

## Speed

| Device  | SHA-256 Hashrate | Scrypt Hashrate |
|---------|-------------------|------------------|
| M2 Air | 15.7 MH/s        | 18.91 kH/s       |
