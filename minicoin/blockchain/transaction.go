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

type Transaction struct {
	From      string
	To        string
	Amount    uint
	Timestamp time.Time
}

func (t Transaction) String() string {
	return fmt.Sprintf("from %s to %s amount %d at %s", t.From, t.To, t.Amount, t.Timestamp.String())
}

func CreateTransaction(from string, to string, amount uint, timestamp time.Time) Transaction {
	return Transaction{
		From:      from,
		To:        to,
		Amount:    amount,
		Timestamp: timestamp,
	}
}
