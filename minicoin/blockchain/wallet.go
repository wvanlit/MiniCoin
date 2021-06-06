/**
 * @Author: Wessel van Lit
 * @Project: minicoin
 * @Date: 06-Jun-2021
 */

package blockchain

import (
	"crypto/ed25519"
	"crypto/rand"
)

type Wallet struct {
	privateKey ed25519.PrivateKey
	publicKey  ed25519.PublicKey
}

func GenerateWallet() (*Wallet, error) {
	pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	return &Wallet{
		privateKey: privKey,
		publicKey:  pubKey,
	}, nil
}

func (w *Wallet) GetPublicAddress() string {
	return string(w.publicKey)
}
