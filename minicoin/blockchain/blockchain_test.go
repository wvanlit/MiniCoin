/**
 * @Author: Wessel van Lit
 * @Project: minicoin
 * @Date: 25-May-2021
 */

package blockchain

import (
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	// Make the PoW easier
	oldHashStart := HashStart
	HashStart = "c001"
	defer func() {
		HashStart = oldHashStart
	}()

	m.Run()
}

func createBlock() Block {
	// This is necessary to have consistent hashes on different testing runs
	blocktime := time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local)

	return Block{
		Transactions: []Transaction{
			CreateTransaction("A", "B", MiniCoin, blocktime),
			CreateTransaction("A", "C", DecaCoin*2, blocktime),
			CreateTransaction("C", "B", HectoCoin, blocktime),
		},
		HashOfPreviousBlock: "ABCDEFG",
		Timestamp:           blocktime,
	}
}

func TestBlock_HashString(t *testing.T) {
	block := createBlock()
	hashString := block.GetHashableString(1)

	if hashString == "" || len(hashString) < 10 {
		t.Fail()
	}
}

func TestBlock_HashStringConsitency(t *testing.T) {
	block := createBlock()
	prev := block.GetHashableString(1)
	for i := 0; i < 10; i++ {
		curr := block.GetHashableString(1)

		if curr != prev {
			t.Fatalf("Hashing Function is not consistent:\ncurrent:'%s'\nprev'%s'", curr, prev)
		}
	}
}

func TestBlock_ComputeBlock(t *testing.T) {
	block := createBlock()
	err := block.ComputeBlock()

	if err != nil {
		t.Fatalf("Error on compute block: %s", err)
	}

	if !IsValidHash(&block, block.Nonce) {
		t.Fatalf("Hash is not valid on nonce: %d", block.Nonce)
	}
}

func TestBlock_ComputeBlockVerification(t *testing.T) {
	outputs := make([]string, 5)
	block := createBlock()
	for index := 0; index < 5; index++ {
		err := block.ComputeBlock()
		if err != nil {
			t.Fatalf("Error on compute block: %s", err)
		}
		if !IsValidHash(&block, block.Nonce) {
			t.Fatalf("Hash is not valid on nonce: %d", block.Nonce)
		}
		outputs[index] = block.Hash
		t.Logf("Got hash '%s' on nonce '%d'", block.Hash, block.Nonce)
	}

	for _, output := range outputs {
		for _, output2 := range outputs {
			if output != output2 {
				t.Fatalf("'%s' and '%s' aren't equal", output, output2)
			}
		}
	}
}

func BenchmarkBlock_ComputeBlock(b *testing.B) {
	block := createBlock()
	for i := 0; i < b.N; i++ {
		block.ComputeBlock()
	}
}
