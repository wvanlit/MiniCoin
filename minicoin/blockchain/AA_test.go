/**
 * @Author: Wessel van Lit
 * @Project: minicoin
 * @Date: 04-Jun-2021
 * This files is called AA_test because Go loads files based on their names during testing
 */

package blockchain

import (
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestMain(m *testing.M) {
	// Make the PoW easier
	oldHashStart := HASH_START
	HASH_START = "c001"
	defer func() {
		HASH_START = oldHashStart
	}()

	log.SetLevel(log.DebugLevel)
	m.Run()
}
