/**
 * @Author: Wessel van Lit
 * @Project: minicoin
 * @Date: 25-May-2021
 */

package blockchain

import (
	"crypto/sha1"
	"fmt"
	"math"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

var HASH_START = "c0ffee"

type Block struct {
	// These should not be hashed
	Index uint
	Hash  string

	// These should be included in the hash
	Transactions        []SignedTransaction
	HashOfPreviousBlock string
	Timestamp           time.Time
	Nonce               uint
}

func CreateNewBlock(blockChain *BlockChain, transactions []SignedTransaction) Block {
	chainLength := len(blockChain.blocks)
	prevBlock := &blockChain.blocks[chainLength-1]

	return Block{
		Index:               uint(chainLength),
		Transactions:        transactions,
		HashOfPreviousBlock: prevBlock.Hash,
		Timestamp:           time.Now(),
		Nonce:               0,
	}
}

func CreateGenesisBlock() Block {
	return Block{
		Index:               0,
		Transactions:        make([]SignedTransaction, 0),
		HashOfPreviousBlock: "",
		Timestamp:           time.Now(),
	}
}

func (b Block) GetHashableString(nonce uint) string {
	str := ""
	str += "=== Nonce ===\n"
	str += fmt.Sprintf("%d\n", nonce)
	str += "=== Transactions ===\n"
	for _, transaction := range b.Transactions {
		str += transaction.String() + "\n"
	}
	str += "=== Previous Hash ===\n"
	str += b.HashOfPreviousBlock + "\n"
	str += "=== Timestamp ===\n"
	str += b.Timestamp.String()
	return str
}

func HashInput(input string) string {
	h := sha1.New()
	h.Write([]byte(input))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func IsValidHash(hash string) bool {
	return strings.HasPrefix(hash, HASH_START)
}

func (b *Block) ComputeBlock() error {
	b.computeProofOfWork()
	if !b.IsValid() {
		return fmt.Errorf("computed nonce '%d' is not actually valid", b.Nonce)
	}
	b.Hash = HashInput(b.GetHashableString(b.Nonce))
	return nil
}

func (b *Block) computeProofOfWork() {
	var i uint
	for i = 0; i < math.MaxUint64; i++ {
		log.Tracef("Trying Nonce: %d", i)

		if b.isValidNonce(i) {
			log.Infof("Found PoW nonce for block: %d", i)
			b.Nonce = i
			return
		}
	}
	log.Error("Could not find nonce for block:", *b)
}

func (b *Block) IsValid() bool {
	hash := HashInput(b.GetHashableString(b.Nonce))
	return IsValidHash(hash)
}

func (b *Block) isValidNonce(i uint) bool {
	hash := HashInput(b.GetHashableString(i))
	return IsValidHash(hash)
}
