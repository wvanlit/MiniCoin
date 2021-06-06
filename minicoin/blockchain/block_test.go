/**
 * @Author: Wessel van Lit
 * @Project: minicoin
 * @Date: 25-May-2021
 */

package blockchain

import (
	"math/rand"
	"testing"
	"time"
)

func createBlock(walletA *Wallet, walletB *Wallet, walletC *Wallet) Block {
	// This is necessary to have consistent hashes on different testing runs
	blocktime := time.Date(
		rand.Int(), time.Month(rand.Intn(12)), rand.Intn(25), rand.Intn(23), rand.Intn(60), rand.Intn(60), 0, time.Local)

	return Block{
		Transactions: []SignedTransaction{
			CreateTransaction(walletA, walletB.GetPublicAddress(), uint(rand.Uint32()), blocktime),
			CreateTransaction(walletA, walletC.GetPublicAddress(), uint(rand.Uint32()), blocktime),
			CreateTransaction(walletB, walletB.GetPublicAddress(), uint(rand.Uint32()), blocktime),
		},
		HashOfPreviousBlock: "X",
		Timestamp:           blocktime,
	}
}

func TestBlock_HashString(t *testing.T) {
	wA, wB, wC, err := CreateThreeWallets()
	if err != nil {
		t.Fatalf("Could not create wallets")
	}
	block := createBlock(wA, wB, wC)
	hashString := block.GetHashableString(1)

	if hashString == "" || len(hashString) < 10 {
		t.Fail()
	}
}

func TestBlock_HashStringConsistency(t *testing.T) {
	wA, wB, wC, err := CreateThreeWallets()
	if err != nil {
		t.Fatalf("Could not create wallets")
	}
	block := createBlock(wA, wB, wC)
	prev := block.GetHashableString(1)
	for i := 0; i < 10; i++ {
		curr := block.GetHashableString(1)

		if curr != prev {
			t.Fatalf("Hashing Function is not consistent:\ncurrent:'%s'\nprev'%s'", curr, prev)
		}
	}
}

func TestBlock_ComputeBlock(t *testing.T) {
	wA, wB, wC, err := CreateThreeWallets()
	if err != nil {
		t.Fatalf("Could not create wallets")
	}
	block := createBlock(wA, wB, wC)
	err = block.ComputeBlock()

	if err != nil {
		t.Fatalf("Error on compute block: %s", err)
	}

	if !block.IsValid() {
		t.Fatalf("Hash is not valid on nonce: %d", block.Nonce)
	}
}

func TestBlock_ComputeBlockVerification(t *testing.T) {
	n := 10
	outputs := make([]string, n)

	wA, wB, wC, err := CreateThreeWallets()
	if err != nil {
		t.Fatalf("Could not create wallets")
	}
	block := createBlock(wA, wB, wC)

	for index := 0; index < len(outputs); index++ {
		err := block.ComputeBlock()
		if err != nil {
			t.Fatalf("Error on compute block: %s", err)
		}
		if !block.IsValid() {
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

func TestCreateNewBlock(t *testing.T) {
	wA, wB, wC, err := CreateThreeWallets()
	if err != nil {
		t.Fatalf("Could not create wallets")
	}

	blockChain, err := CreateNewBlockChain(wA)
	if err != nil {
		t.Fatalf("Error on creating blockchain: %s", err)
	}

	transactions := []SignedTransaction{
		CreateTransaction(wA, wB.GetPublicAddress(), MiniCoin, time.Now()),
		CreateTransaction(wB, wC.GetPublicAddress(), 2*MiniCoin, time.Now()),
		CreateTransaction(wC, wA.GetPublicAddress(), 3*MiniCoin, time.Now()),
	}

	block := CreateNewBlock(&blockChain, transactions)

	if blockChain.blocks[0].Hash != block.HashOfPreviousBlock {
		t.Fatalf("Hash of previous block not added correctly")
	}

	if int(block.Index) != len(blockChain.blocks) {
		t.Fatalf("Block index not equal to the end of the blocks slice")
	}
}

func BenchmarkBlock_ComputeBlock(b *testing.B) {
	wA, wB, wC, err := CreateThreeWallets()
	if err != nil {
		b.Fatalf("Could not create wallets")
	}
	block := createBlock(wA, wB, wC)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		block.ComputeBlock()
	}
}
