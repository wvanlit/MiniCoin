/**
 * @Author: Wessel van Lit
 * @Project: minicoin
 * @Date: 25-May-2021
 */

package blockchain

import (
	"fmt"
	"time"
)

type Blockchain struct {
	blocks []Block
}

type Block struct {
	// These should not be hashed
	Prev  *Block
	Next  *Block
	Index uint
	Hash  string

	// These should be included in the hash
	Transactions        []Transaction
	HashOfPreviousBlock string
	Timestamp           time.Time
	Nonce               uint
}

type MilestoneBlock struct {
	Ledger map[string]int
	Block
}

func (b Block) GetHashableString(nonce uint) string {
	str := ""
	str += "=== Transactions ===\n"
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
