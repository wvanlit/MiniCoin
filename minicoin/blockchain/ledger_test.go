/**
 * @Author: Wessel van Lit
 * @Project: minicoin
 * @Date: 05-Jun-2021
 */

package blockchain

import "testing"

func TestCreateNewLedger_isEmpty(t *testing.T) {
	ledger := CreateNewLedger()
	if len(ledger.tokenLedger) != 0 {
		t.Fatalf("New ledger is not empty")
	}
}

func TestLedger_AddTokensToNewUser(t *testing.T) {
	user := "A"
	ledger := CreateNewLedger()
	ledger.AddTokensTo(user, HectoCoin)
}
