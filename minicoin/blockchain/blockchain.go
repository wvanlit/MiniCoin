/**
 * @Author: Wessel van Lit
 * @Project: minicoin
 * @Date: 25-May-2021
 */

package blockchain

import "fmt"

const MINING_REWARD = HectoCoin

type BlockChain struct {
	blocks []Block
	ledger *Ledger
}

func CreateNewBlockChain(user *Wallet) (BlockChain, error) {
	ledger := CreateNewLedger()
	blockChain := BlockChain{
		blocks: make([]Block, 0),
		ledger: &ledger,
	}
	err := blockChain.AppendBlockAfterPoW(CreateGenesisBlock(), user)
	return blockChain, err
}

func (b *BlockChain) AppendBlockAfterPoW(block Block, miner *Wallet) error {
	block.Transactions = append(block.Transactions, CreateMiningReward(miner))

	err := block.ComputeBlock()
	if err != nil {
		return err
	}

	err = b.UpdateLedgerWithBlock(&block)
	if err != nil {
		return err
	}

	b.blocks = append(b.blocks, block)

	return b.ValidateChain()
}

func (b BlockChain) ValidateChain() error {
	for i, block := range b.blocks {
		if !block.IsValid() {
			return fmt.Errorf("block at index %d is not valid", i)
		}
	}
	return nil
}

func (b *BlockChain) UpdateLedgerWithBlock(block *Block) error {
	ledgerCopy := *b.ledger

	for _, transaction := range block.Transactions {
		err := b.ledger.ProcessTransaction(transaction)
		if err != nil {
			// Revert to old ledger state
			b.ledger = &ledgerCopy
			return err
		}
	}
	return nil
}
