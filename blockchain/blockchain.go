package blockchain

import "time"

type Blockchain struct {
	genesisBlock Block
}

type Block struct {
	// These should not be hashed
	Prev  *Block
	Next  *Block
	Index int
	Hash  string
	Nonce int

	// These should be included in the hash
	Transactions        []Transaction
	HashOfPreviousBlock string
	Timestamp           time.Time
}

type MilestoneBlock struct {
	Ledger map[string]int
	Block
}

func (b Block) validateHash(nonce int) bool {
	return false
}

func (b Block) GetHashableString() string {
	str := ""
	str += "=== Transactions ===\n"
	for _, transaction := range b.Transactions {
		str += transaction.String() + "\n"
	}
	str += "=== Previous Hash ===\n"
	str += b.Hash + "\n"
	str += "=== Timestamp ===\n"
	str += b.Timestamp.String()
	return str
}
