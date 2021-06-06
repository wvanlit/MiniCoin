/**
 * @Author: Wessel van Lit
 * @Project: minicoin
 * @Date: 04-Jun-2021
 */

package blockchain

import (
	"reflect"
	"testing"
	"time"
)

func TestCreateNewBlockChain(t *testing.T) {
	usr, err := GenerateWallet()
	if err != nil {
		t.Fatalf("Error on generating wallet: %s", err)
	}

	blockChain, err := CreateNewBlockChain(usr)
	if err != nil {
		t.Fatalf("Error on creating blockchain: %s", err)
	}

	if blockChain.ValidateChain() != nil {
		t.Fatalf("Blockchain is not actually valid")
	}

	amount, err := blockChain.ledger.GetTokenTotalFrom(usr.GetPublicAddress())
	if err != nil || amount != HectoCoin {
		t.Fatalf("User did not get Mining Reward")
	}
}

func TestBlockChain_AppendBlockAfterPoW(t *testing.T) {
	usr, err := GenerateWallet()
	wA, wB, wC, err := CreateThreeWallets()
	if err != nil {
		t.Fatalf("Error on generating wallet: %s", err)
	}
	blockChain, err := CreateNewBlockChain(usr)
	if err != nil {
		t.Fatalf("Error on creating blockchain: %s", err)
	}

	blockChain.ledger.AddTokensTo(wA.GetPublicAddress(), MiniCoin)
	blockChain.ledger.AddTokensTo(wB.GetPublicAddress(), MiniCoin)
	blockChain.ledger.AddTokensTo(wC.GetPublicAddress(), MiniCoin)

	transactions := []SignedTransaction{
		CreateTransaction(wA, wB.GetPublicAddress(), MiniCoin, time.Now()),
		CreateTransaction(wB, wC.GetPublicAddress(), MiniCoin, time.Now()),
		CreateTransaction(wC, wA.GetPublicAddress(), MiniCoin, time.Now()),
	}

	block := CreateNewBlock(&blockChain, transactions)

	err = blockChain.AppendBlockAfterPoW(block, usr)
	if err != nil {
		t.Fatalf("Error on appending block to blockchain: %s", err)
	}

	if blockChain.blocks[1].Transactions[0].String() != block.Transactions[0].String() {
		t.Fatalf("Block not added to blockchain correctly!\n%v", blockChain.blocks)
	}

	if blockChain.ValidateChain() != nil {
		t.Fatalf("Chain is no longer valid")
	}

	amount, err := blockChain.ledger.GetTokenTotalFrom(usr.GetPublicAddress())
	if err != nil || amount != 2*MINING_REWARD {
		t.Fatalf("Miner did not get reward")
	}
}

func TestBlockChain_UpdateLedgerWithValidBlock(t *testing.T) {
	usr, err := GenerateWallet()
	wA, wB, wC, err := CreateThreeWallets()
	if err != nil {
		t.Fatalf("Error on generating wallet: %s", err)
	}

	blockChain, err := CreateNewBlockChain(usr)
	if err != nil {
		t.Fatalf("Error on creating blockchain: %s", err)
	}

	var startingAmount uint = DecaCoin

	blockChain.ledger.AddTokensTo(wA.GetPublicAddress(), startingAmount)
	blockChain.ledger.AddTokensTo(wB.GetPublicAddress(), startingAmount)
	blockChain.ledger.AddTokensTo(wC.GetPublicAddress(), startingAmount)

	var aSendsToB uint = MiniCoin
	var bSendsToC uint = MiniCoin * 2
	var cSendsToA uint = MiniCoin * 3

	transactions := []SignedTransaction{
		CreateTransaction(wA, wB.GetPublicAddress(), aSendsToB, time.Now()),
		CreateTransaction(wB, wC.GetPublicAddress(), bSendsToC, time.Now()),
		CreateTransaction(wC, wA.GetPublicAddress(), cSendsToA, time.Now()),
	}

	block := CreateNewBlock(&blockChain, transactions)
	oldLedger := copyLedger(blockChain.ledger.tokenLedger)

	err = blockChain.AppendBlockAfterPoW(block, usr)
	if err != nil {
		t.Fatalf("Error on appending block to blockchain: %s", err)
	}

	if reflect.DeepEqual(oldLedger, blockChain.ledger.tokenLedger) {
		t.Fatalf("Old Ledger and New Ledger should be different")
	}

	aTotal, errA := blockChain.ledger.GetTokenTotalFrom(wA.GetPublicAddress())
	bTotal, errB := blockChain.ledger.GetTokenTotalFrom(wB.GetPublicAddress())
	cTotal, errC := blockChain.ledger.GetTokenTotalFrom(wC.GetPublicAddress())

	if errA != nil || errB != nil || errC != nil {
		t.Fatalf("Could not get token total from existing users")
	}

	if aTotal != startingAmount-aSendsToB+cSendsToA {
		t.Fatalf("A does not have the correct amount")
	}

	if bTotal != startingAmount-bSendsToC+aSendsToB {
		t.Fatalf("B does not have the correct amount")
	}

	if cTotal != startingAmount-cSendsToA+bSendsToC {
		t.Fatalf("C does not have the correct amount")
	}
}

func TestBlockChain_UpdateLedgerWithInvalidBlock(t *testing.T) {
	usr, err := GenerateWallet()
	wA, wB, wC, err := CreateThreeWallets()
	if err != nil {
		t.Fatalf("Error on generating wallet: %s", err)
	}

	blockChain, err := CreateNewBlockChain(usr)
	if err != nil {
		t.Fatalf("Error on creating blockchain: %s", err)
	}

	var startingAmount uint = MiniCoin

	blockChain.ledger.AddTokensTo(wA.GetPublicAddress(), startingAmount)
	blockChain.ledger.AddTokensTo(wB.GetPublicAddress(), startingAmount)
	blockChain.ledger.AddTokensTo(wC.GetPublicAddress(), startingAmount)

	var aSendsToB uint = MiniCoin * 5
	var bSendsToC uint = MiniCoin * 2
	var cSendsToA uint = MiniCoin * 3

	transactions := []SignedTransaction{
		CreateTransaction(wA, wB.GetPublicAddress(), aSendsToB, time.Now()),
		CreateTransaction(wB, wC.GetPublicAddress(), bSendsToC, time.Now()),
		CreateTransaction(wC, wA.GetPublicAddress(), cSendsToA, time.Now()),
	}

	block := CreateNewBlock(&blockChain, transactions)
	oldLedger := copyLedger(blockChain.ledger.tokenLedger)

	err = blockChain.AppendBlockAfterPoW(block, usr)
	if err == nil {
		t.Fatalf("Should have error on appending an invalid block")
	}

	if !reflect.DeepEqual(oldLedger, blockChain.ledger.tokenLedger) {
		t.Fatalf("Old Ledger and New Ledger should be the same, since the transaction is not valid")
	}
}

func copyLedger(ledger map[string]uint) map[string]uint {
	ledgerCopy := make(map[string]uint)
	for s, u := range ledger {
		ledgerCopy[s] = u
	}
	return ledgerCopy
}
