/**
 * @Author: Wessel van Lit
 * @Project: minicoin
 * @Date: 25-May-2021
 */

package blockchain

import (
	"crypto/sha1"
	"fmt"
	"math"
	"strings"

	log "github.com/sirupsen/logrus"
)

var HashStart = "c0ffee"

func HashInput(input string) string {
	h := sha1.New()
	h.Write([]byte(input))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func IsValidHash(block *Block, nonce uint) bool {
	hashString := HashInput(block.GetHashableString(nonce))
	return strings.HasPrefix(hashString, HashStart)
}

func (b *Block) ComputeBlock() error {
	b.computeProofOfWork()
	if !IsValidHash(b, b.Nonce) {
		return fmt.Errorf("computed nonce '%d' is not actually valid", b.Nonce)
	}
	b.Hash = HashInput(b.GetHashableString(b.Nonce))
	return nil
}

func (b *Block) computeProofOfWork() {
	var i uint
	for i = 0; i < math.MaxUint64; i++ {
		log.Tracef("Trying Nonce: %d", i)

		if IsValidHash(b, i) {
			log.Infof("Found PoW nonce for block: %d", i)
			b.Nonce = i
			return
		}
	}
	log.Error("Could not find nonce for block:", *b)
}
