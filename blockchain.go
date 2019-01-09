package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

// Blockchain list
var Blockchain []Block

// Block structure
type Block struct {
	Index     int    // Position of the data record in the blockchain
	Timestamp string // Time the data is written
	Data      string // Data of Block
	Hash      string // SHA256 identifier representing this data record
	PrevHash  string // SHA256 identifier of the previous record in the chain
	Suffix    string // Value added to hash
	Hardness  int    // Number of beginning "0" in hash
}

// Payload structure
type Payload struct {
	Data string `json:"data"`
}

func (b *Block) String() string {
	return fmt.Sprintf("index: %d - timestamp: %s - prevhash: %s - data: %s - suffix: %s", b.Index, b.Timestamp, b.PrevHash, b.Data, b.Suffix)
}

func calculateHash(block Block) string {
	h := sha256.New()
	h.Write([]byte(block.String()))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// Create a new Block
func generateBlock(oldBlock Block, data string) Block {
	var newBlock Block
	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = time.Now().String()
	newBlock.Data = data
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hardness = int(oldBlock.Index/256 + 1)

	trials := 0
	for i := 0; ; i++ {
		newBlock.Suffix = fmt.Sprintf("%x", i)
		potentialHash := calculateHash(newBlock)
		if !isHashValid(potentialHash, newBlock.Hardness) {
			trials++
		} else {
			fmt.Println(potentialHash, "- valid hash after", trials, "trials. Hardness set to", newBlock.Hardness)
			newBlock.Hash = potentialHash
			trials = 0
			break
		}
	}

	return newBlock
}

// Replace Blockchain if provided chain is more up to date
func replaceChain(newChain []Block) {
	if len(newChain) > len(Blockchain) {
		Blockchain = newChain
	}
}

// Check if Block is valid
func isBlockValid(oldBlock Block, newBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}
	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}
	if newBlock.Hash != calculateHash(newBlock) {
		return false
	}
	return true
}

func isHashValid(hash string, hardness int) bool {
	prefix := strings.Repeat("0", hardness)
	return strings.HasPrefix(hash, prefix)
}
