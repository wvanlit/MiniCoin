/**
 * @Author: Wessel van Lit
 * @Project: minicoin
 * @Date: 25-May-2021
 */

package blockchain

import (
	"crypto/ed25519"
	"crypto/sha1"
	"fmt"
	"time"
)

const MiningRewardUser = "MINING_REWARD"

type SignedTransaction struct {
	Signature []byte
	Transaction
}

func (st SignedTransaction) IsVerified() bool {
	if st.From == MiningRewardUser {
		return ed25519.Verify([]byte(st.To), st.Hash(), st.Signature)
	}
	return ed25519.Verify([]byte(st.From), st.Hash(), st.Signature)
}

func (t Transaction) SignTransaction(wallet *Wallet) SignedTransaction {
	signature := ed25519.Sign(wallet.privateKey, t.Hash())
	return SignedTransaction{Signature: signature, Transaction: t}
}

type Transaction struct {
	From      string
	To        string
	Amount    uint
	Timestamp time.Time
}

func (t Transaction) String() string {
	return fmt.Sprintf("from %s to %s amount %d at %s", t.From, t.To, t.Amount, t.Timestamp.String())
}

func (t Transaction) Hash() []byte {
	h := sha1.New()
	h.Write([]byte(t.String()))
	return h.Sum(nil)
}

func CreateTransaction(from *Wallet, to string, amount uint, timestamp time.Time) SignedTransaction {
	return Transaction{
		From:      from.GetPublicAddress(),
		To:        to,
		Amount:    amount,
		Timestamp: timestamp,
	}.SignTransaction(from)
}

func CreateMiningReward(to *Wallet) SignedTransaction {
	return Transaction{
		From:      MiningRewardUser,
		To:        to.GetPublicAddress(),
		Amount:    MINING_REWARD,
		Timestamp: time.Now(),
	}.SignTransaction(to)
}
