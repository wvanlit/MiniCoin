/**
 * @Author: Wessel van Lit
 * @Project: minicoin
 * @Date: 04-Jun-2021
 */

package blockchain

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

type Ledger struct {
	tokenLedger map[string]uint
}

func CreateNewLedger() Ledger {
	return Ledger{
		tokenLedger: make(map[string]uint),
	}
}

func (l Ledger) GetTokenTotalFrom(user string) (uint, error) {
	tokens, found := l.tokenLedger[user]
	if !found {
		return 0, fmt.Errorf("user '%s' not found in ledger", user)
	}
	return tokens, nil
}

func (l *Ledger) AddTokensTo(user string, amount uint) {
	currTokens, err := l.GetTokenTotalFrom(user)

	if err != nil {
		l.tokenLedger[user] = amount
	} else {
		l.tokenLedger[user] = currTokens + amount
	}
}

func (l *Ledger) RemoveTokensFrom(user string, amount uint) error {
	if user == MiningRewardUser {
		return nil
	}

	currTokens, err := l.GetTokenTotalFrom(user)
	if err != nil {
		return err
	}

	if currTokens < amount {
		return fmt.Errorf("user '%s' does not have enough tokens to remove '%d'", user, amount)
	}

	l.tokenLedger[user] = currTokens - amount
	return nil
}

func (l Ledger) IsValidTransaction(transaction SignedTransaction) bool {
	if !transaction.IsVerified() {
		log.Warnf("Transaction '%s' is not valid since it isn't signed correctly", transaction.String())
		return false
	}

	// Check if mining reward is correct
	if transaction.From == MiningRewardUser {
		if transaction.Amount == MINING_REWARD {
			return true
		} else {
			log.Warnf("Transaction '%s' is not valid since it doesn't have the correct amount for a Mining Reward", transaction.String())
			return false
		}
	}

	if transaction.To == MiningRewardUser {
		log.Warnf("Transaction '%s' is not valid since it is sent to ", transaction.String())
		return false
	}

	// Sender has enough money
	amount, err := l.GetTokenTotalFrom(transaction.From)
	if err != nil || amount < transaction.Amount {
		log.Warnf("Transaction '%s' is not valid since '%s' does not have enough tokens", transaction.String(), transaction.From)
		return false
	}
	// Sender is not sending money to themselves
	if transaction.From == transaction.To {
		log.Warnf("Transaction '%s' is not valid since '%s' is sending it to themselves", transaction.String(), transaction.From)
		return false
	}
	return true
}

func (l *Ledger) ProcessTransaction(transaction SignedTransaction) error {
	if l.IsValidTransaction(transaction) {
		err := l.RemoveTokensFrom(transaction.From, transaction.Amount)
		if err != nil {
			return err
		}
		l.AddTokensTo(transaction.To, transaction.Amount)
		return nil
	} else {
		return fmt.Errorf("transaction is not valid: '%s'", transaction)
	}
}
