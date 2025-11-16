# GenesisG0

Port of https://github.com/lhartikk/GenesisH0 to golang with parallelism and time estimates.

Same syntax

## Examples

Bitcoin
```
./genesis -z "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks" -n 2083236893 -t 1231006505
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
