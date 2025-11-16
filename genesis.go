package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"flag"
	"fmt"
	"math/big"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/crypto/scrypt"
)

type Options struct {
	Timestamp string
	PubKey    string
	Time      int64
	Bits      uint32
	Nonce     uint32
	Algorithm string
	Value     int64
	NumWorkers int
}

type HashResult struct {
	SHA256Hash []byte
	Hash       []byte
	Nonce      uint32
}

func main() {
	options := getArgs()

	if !isValidAlgorithm(options.Algorithm) {
		fmt.Fprintf(os.Stderr, "Error: Given algorithm must be one of: SHA256, scrypt\n")
		os.Exit(1)
	}

	inputScript := createInputScript(options.Timestamp)
	outputScript := createOutputScript(options.PubKey)

	tx := createTransaction(inputScript, outputScript, options)
	hashMerkleRoot := doubleSHA256(tx)

	printBlockInfo(options, hashMerkleRoot)

	blockHeader := createBlockHeader(hashMerkleRoot, uint32(options.Time), options.Bits, options.Nonce)
	genesisHash, nonce := generateHashParallel(blockHeader, options.Algorithm, options.Nonce, options.Bits, options.NumWorkers)

	announceFoundGenesis(genesisHash, nonce)
}

func getArgs() Options {
	timestamp := flag.String("timestamp", "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks",
		"the pszTimestamp found in the coinbase of the genesisblock")
	flagTime := flag.Int64("time", time.Now().Unix(),
		"the (unix) time when the genesisblock is created")
	nonceFlag := flag.Uint("nonce", 0,
		"the first value of the nonce that will be incremented when searching the genesis hash")
	algorithm := flag.String("algorithm", "SHA256",
		"the PoW algorithm: [SHA256|scrypt]")
	flagZ := flag.String("z", "",
		"short flag for timestamp")
	flagT := flag.Int64("t", 0,
		"short flag for time")
	flagN := flag.Uint("n", 0,
		"short flag for nonce")
	flagA := flag.String("a", "",
		"short flag for algorithm")
	flagP := flag.String("p", "04678afdb0fe5548271967f1a67130b7105cd6a828e03909a67962e0ea1f61deb649f6bc3f4cef38c4f35504e51ec112de5c384df7ba0b8d578a4c702b6bf11d5f",
		"the pubkey found in the output script")
	flagV := flag.Int64("v", 5000000000,
		"the value in coins for the output, full value")
	flagB := flag.Int("b", 0,
		"the target in compact representation, associated to a difficulty of 1")
	flagC := flag.Int("c", runtime.NumCPU(),
		"number of CPU cores to use for mining")

	flag.Parse()

	if *flagZ != "" {
		*timestamp = *flagZ
	}
	if *flagT != 0 {
		*flagTime = *flagT
	}
	if *flagN != 0 {
		*nonceFlag = *flagN
	}
	if *flagA != "" {
		*algorithm = *flagA
	}

	bits := uint32(*flagB)
	if bits == 0 {
		if *algorithm == "scrypt" {
			bits = 0x1e0ffff0
		} else {
			bits = 0x1d00ffff
		}
	}

	return Options{
		Timestamp:  *timestamp,
		PubKey:     *flagP,
		Time:       *flagTime,
		Bits:       bits,
		Nonce:      uint32(*nonceFlag),
		Algorithm:  *algorithm,
		Value:      *flagV,
		NumWorkers: *flagC,
	}
}

func isValidAlgorithm(algo string) bool {
	supported := []string{"SHA256", "scrypt"}
	for _, a := range supported {
		if algo == a {
			return true
		}
	}
	return false
}

func createInputScript(pszTimestamp string) []byte {
	timestampBytes := []byte(pszTimestamp)
	pszPrefix := ""

	if len(timestampBytes) > 76 {
		pszPrefix = "4c"
	}

	lengthInHex := fmt.Sprintf("%x", len(timestampBytes))
	scriptPrefix := "04ffff001d0104" + pszPrefix + lengthInHex
	inputScriptHex := scriptPrefix + hex.EncodeToString(timestampBytes)

	fmt.Println(inputScriptHex)

	result, err := hex.DecodeString(inputScriptHex)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error decoding input script: %v\n", err)
		os.Exit(1)
	}
	return result
}

func createOutputScript(pubkey string) []byte {
	scriptLen := "41"
	opChecksig := "ac"
	scriptHex := scriptLen + pubkey + opChecksig

	result, err := hex.DecodeString(scriptHex)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error decoding output script: %v\n", err)
		os.Exit(1)
	}
	return result
}

func createTransaction(inputScript, outputScript []byte, options Options) []byte {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, uint32(1))
	binary.Write(buf, binary.LittleEndian, uint8(1))
	buf.Write(make([]byte, 32))
	binary.Write(buf, binary.LittleEndian, uint32(0xFFFFFFFF))
	binary.Write(buf, binary.LittleEndian, uint8(len(inputScript)))
	buf.Write(inputScript)
	binary.Write(buf, binary.LittleEndian, uint32(0xFFFFFFFF))
	binary.Write(buf, binary.LittleEndian, uint8(1))
	binary.Write(buf, binary.LittleEndian, options.Value)
	binary.Write(buf, binary.LittleEndian, uint8(0x43))
	buf.Write(outputScript)
	binary.Write(buf, binary.LittleEndian, uint32(0))

	return buf.Bytes()
}

func createBlockHeader(hashMerkleRoot [32]byte, timestamp uint32, bits, nonce uint32) []byte {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, uint32(1))
	buf.Write(make([]byte, 32))
	buf.Write(hashMerkleRoot[:])
	binary.Write(buf, binary.LittleEndian, timestamp)
	binary.Write(buf, binary.LittleEndian, bits)
	binary.Write(buf, binary.LittleEndian, nonce)

	return buf.Bytes()
}

func doubleSHA256(data []byte) [32]byte {
	hash1 := sha256.Sum256(data)
	hash2 := sha256.Sum256(hash1[:])
	return hash2
}

func reverseBytes(b []byte) []byte {
	for i := 0; i < len(b)/2; i++ {
		j := len(b) - 1 - i
		b[i], b[j] = b[j], b[i]
	}
	return b
}

func generateHashParallel(baseHeader []byte, algorithm string, startNonce, bits uint32, numWorkers int) (string, uint32) {
	fmt.Printf("Searching for genesis hash with %d workers..\n", numWorkers)

	target := calculateTarget(bits)
	resultChan := make(chan *HashResult, 1)
	stopChan := make(chan bool)
	var hashCount int64
	lastUpdated := time.Now()
	var wg sync.WaitGroup

	// Start worker goroutines
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			worker(workerID, baseHeader, algorithm, startNonce, target, resultChan, stopChan, &hashCount, numWorkers)
		}(i)
	}

	// Monitor results and print hashrate
	startMonitor := time.Now()
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			select {
			case <-stopChan:
				return
			default:
				now := time.Now()
				elapsed := now.Sub(lastUpdated).Seconds()
				if elapsed > 0 {
					currentCount := atomic.LoadInt64(&hashCount)
					hashrate := float64(currentCount) / elapsed

					var hashRateStr string
					if hashrate >= 1000000 {
						hashRateStr = fmt.Sprintf("%.2f mh/s", hashrate/1000000)
					} else if hashrate >= 1000 {
						hashRateStr = fmt.Sprintf("%.2f kh/s", hashrate/1000)
					} else {
						hashRateStr = fmt.Sprintf("%.2f h/s", hashrate)
					}

					// Calculate expected hashes based on difficulty
					// The difficulty is determined by how the target compares to the max possible target
					// For this block's bits value, we use it as our baseline
					// Expected hashes = (2^256) / target gives us roughly how many hashes needed
					targetBig := new(big.Int).SetBytes(target)
					var expectedHashes float64
					
					if targetBig.Sign() > 0 {
						// Use 2^256 as the theoretical maximum search space
						// This gives us expected attempts = 2^256 / target
						maxPossible := new(big.Int)
						maxPossible.Lsh(big.NewInt(1), 256)
						
						result := new(big.Int).Div(maxPossible, targetBig)
						expectedHashes, _ = new(big.Float).SetInt(result).Float64()
					} else {
						expectedHashes = float64(1 << 32)
					}

					estimatedSeconds := expectedHashes / hashrate
					estimatedTime := time.Duration(int64(estimatedSeconds)) * time.Second

					elapsedTime := now.Sub(startMonitor)
					var elapsedStr string
					if elapsedTime.Hours() >= 1 {
						elapsedStr = fmt.Sprintf("%.1fh", elapsedTime.Hours())
					} else if elapsedTime.Minutes() >= 1 {
						mins := int(elapsedTime.Minutes())
						secs := int(elapsedTime.Seconds()) % 60
						elapsedStr = fmt.Sprintf("%dm%ds", mins, secs)
					} else {
						elapsedStr = fmt.Sprintf("%ds", int64(elapsedTime.Seconds()))
					}

					fmt.Fprintf(os.Stderr, "\r%s, elapsed: %s, estimate: %v                    ", hashRateStr, elapsedStr, estimatedTime)
					atomic.StoreInt64(&hashCount, 0)
					lastUpdated = now
				}
			}
		}
	}()

	// Wait for result
	result := <-resultChan
	close(stopChan)
	wg.Wait()

	if result != nil && len(result.SHA256Hash) > 0 {
		return hex.EncodeToString(result.SHA256Hash), result.Nonce
	}

	return "", 0
}

func worker(workerID int, baseHeader []byte, algorithm string, startNonce uint32, target []byte, resultChan chan *HashResult, stopChan chan bool, hashCount *int64, numWorkers int) {
	nonce := startNonce + uint32(workerID)
	headerCopy := make([]byte, len(baseHeader))

	for {
		select {
		case <-stopChan:
			return
		default:
			copy(headerCopy, baseHeader)
			nonceBytes := make([]byte, 4)
			binary.LittleEndian.PutUint32(nonceBytes, nonce)
			copy(headerCopy[len(headerCopy)-4:], nonceBytes)

			sha256Hash, headerHash := generateHashesFromBlock(headerCopy, algorithm)
			atomic.AddInt64(hashCount, 1)

			if bytesLessThan(headerHash, target) {
				select {
				case resultChan <- &HashResult{
					SHA256Hash: sha256Hash,
					Hash:       headerHash,
					Nonce:      nonce,
				}:
					return
				case <-stopChan:
					return
				}
			}

			nonce += uint32(numWorkers)
		}
	}
}

func calculateTarget(bits uint32) []byte {
	mantissa := bits & 0xffffff
	exponent := (bits >> 24) & 0xff

	target := make([]byte, 32)

	if exponent >= 1 && exponent <= 32 {
		pos := 32 - int(exponent)

		if pos >= 0 && pos < 32 {
			target[pos] = byte((mantissa >> 16) & 0xff)
			if pos+1 < 32 {
				target[pos+1] = byte((mantissa >> 8) & 0xff)
			}
			if pos+2 < 32 {
				target[pos+2] = byte(mantissa & 0xff)
			}
		}
	}

	return target
}

func reverseBytesCopy(b []byte) []byte {
	reversed := make([]byte, len(b))
	copy(reversed, b)
	reverseBytes(reversed)
	return reversed
}

func generateHashesFromBlock(dataBlock []byte, algorithm string) ([]byte, []byte) {
	sha256Hash := doubleSHA256(dataBlock)
	sha256HashReversed := make([]byte, 32)
	copy(sha256HashReversed, sha256Hash[:])
	reverseBytes(sha256HashReversed)

	var headerHash []byte

	switch algorithm {
	case "SHA256":
		headerHash = make([]byte, 32)
		copy(headerHash, sha256HashReversed)

	case "scrypt":
		hash, err := scrypt.Key(dataBlock, dataBlock, 1024, 1, 1, 32)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error computing scrypt hash: %v\n", err)
			os.Exit(1)
		}
		headerHash = make([]byte, 32)
		copy(headerHash, hash)
		reverseBytes(headerHash)

	default:
		fmt.Fprintf(os.Stderr, "%s algorithm not yet implemented\n", algorithm)
		os.Exit(1)
	}

	return sha256HashReversed, headerHash
}

func bytesLessThan(a, b []byte) bool {
	for i := 0; i < len(a); i++ {
		if a[i] < b[i] {
			return true
		}
		if a[i] > b[i] {
			return false
		}
	}
	return false
}

func printBlockInfo(options Options, hashMerkleRoot [32]byte) {
	reversed := make([]byte, 32)
	copy(reversed, hashMerkleRoot[:])
	reverseBytes(reversed)

	fmt.Printf("algorithm: %s\n", options.Algorithm)
	fmt.Printf("merkle hash: %s\n", hex.EncodeToString(reversed))
	fmt.Printf("pszTimestamp: %s\n", options.Timestamp)
	fmt.Printf("pubkey: %s\n", options.PubKey)
	fmt.Printf("time: %d\n", options.Time)
	fmt.Printf("bits: %s\n", fmt.Sprintf("0x%x", options.Bits))
}

func announceFoundGenesis(genesisHash string, nonce uint32) {
	fmt.Println("\ngenesis hash found!")
	fmt.Printf("nonce: %d\n", nonce)
	fmt.Printf("genesis hash: %s\n", genesisHash)
}
