/**
 * @Author: Wessel van Lit
 * @Project: minicoin
 * @Date: 06-Jun-2021
 */

package blockchain

import (
	"testing"
	"time"
)

func TestTransaction_SignTransaction_withValidWallet(t *testing.T) {
	wallet1, err1 := GenerateWallet()
	wallet2, err2 := GenerateWallet()
	if err1 != nil || err2 != nil {
		t.Fatalf("Cannot generate wallets")
	}

	transaction := CreateTransaction(wallet1, wallet2.GetPublicAddress(), MiniCoin, time.Now())
	signedTransaction := transaction.SignTransaction(wallet1)

	if !signedTransaction.IsVerified() {
		t.Fatalf("Cannot verify transaction")
	}
}

func TestTransaction_SignTransaction_withInvalidWallet(t *testing.T) {
	wallet1, err1 := GenerateWallet()
	wallet2, err2 := GenerateWallet()
	if err1 != nil || err2 != nil {
		t.Fatalf("Cannot generate wallets")
	}

	transaction := CreateTransaction(wallet1, wallet2.GetPublicAddress(), MiniCoin, time.Now())
	signedTransaction := transaction.SignTransaction(wallet2)

	if signedTransaction.IsVerified() {
		t.Fatalf("Should not be able to verify transaction")
	}
}
