package blockchain

import (
	"github.com/wvanlit/minicoin/blockchain"
	"testing"
	"time"
)

func createBlock() blockchain.Block {
	return blockchain.Block{
		Transactions: []blockchain.Transaction{
			blockchain.CreateTransaction("A", "B", 10, time.Now()),
			blockchain.CreateTransaction("A", "C", 5, time.Now()),
			blockchain.CreateTransaction("C", "B", 20, time.Now()),
		},
		HashOfPreviousBlock: "ABCDEFG",
		Timestamp:           time.Now(),
	}
}

func TestBlock_HashString(t *testing.T) {
	block := createBlock()
	hashString := block.GetHashableString()

	if hashString == "" || len(hashString) < 10 {
		t.Fail()
	}
}
