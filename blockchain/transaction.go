package blockchain

import (
	"fmt"
	"time"
)

type Transaction struct {
	From      string
	To        string
	Amount    int
	Timestamp time.Time
}

func (t Transaction) String() string {
	return fmt.Sprintf("from %s to %s amount %d at %s", t.From, t.To, t.Amount, t.Timestamp.String())
}

func CreateTransaction(from string, to string, amount int, timestamp time.Time) Transaction {
	return Transaction{
		From:      from,
		To:        to,
		Amount:    amount,
		Timestamp: timestamp,
	}
}
