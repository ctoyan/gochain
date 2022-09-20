package main

import (
	"fmt"

	"github.com/ctoyan/gochain/internal/account"
	"github.com/ctoyan/gochain/internal/blockchain"
	"github.com/ctoyan/gochain/internal/transaction"
	log "github.com/sirupsen/logrus"
)

func main() {
	accountA, err := account.New()
	if err != nil {
		log.Error("failed to create an account")
	}
	accountB, err := account.New()
	if err != nil {
		log.Error("failed to create an account")
	}
	accountC, err := account.New()
	if err != nil {
		log.Error("failed to create an account")
	}

	miner, err := account.New()
	if err != nil {
		log.Error("failed to create a miner account")
	}

	bc := blockchain.New(miner)
	fmt.Println("miner address", miner.GetAddress())

	generateFakeBlock(bc, accountA, accountB, accountC)
	generateFakeBlock(bc, accountA, accountB, accountC)
	generateFakeBlock(bc, accountA, accountB, accountC)
	// generateFakeBlock(bc, accountA, accountB, accountC)

	bc.Print()
}

func generateFakeBlock(bc *blockchain.Blockchain, accountA, accountB, accountC *account.Account) {
	txs := []*transaction.Transaction{
		transaction.New(accountA, accountB, 1),
		transaction.New(accountA, accountB, 1),
		transaction.New(accountA, accountB, 1),
		transaction.New(accountB, accountA, 1),
		// transaction.New(accountB, accountC, 1),
	}

	for _, tx := range txs {
		err := bc.AddTxToPool(tx)
		if err != nil {
			log.WithError(err).Error("failed to add tx to blockchain tx pool")
		}
	}

	err := bc.AddNewBlock()
	if err != nil {
		log.WithError(err).Error("failed to add a new block")
	}

}
