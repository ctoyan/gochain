package blockchain

import (
	"fmt"

	"github.com/ctoyan/gochain/internal/account"
	"github.com/ctoyan/gochain/internal/block"
	"github.com/ctoyan/gochain/internal/transaction"
	log "github.com/sirupsen/logrus"
)

const MINING_DIFFICULTY = 3

type Blockchain struct {
	txPool       []*transaction.Transaction
	chain        []*block.Block
	minerAddress *account.Account
}

func New(minerAddress *account.Account) *Blockchain {
	bc := new(Blockchain)
	bc.minerAddress = minerAddress
	bc.chain = append(bc.chain, block.New(0, [32]byte{}, []*transaction.Transaction{})) // genesis block
	return bc
}

func (bc *Blockchain) AddNewBlock() error {
	lastBlockHash, err := bc.GetLastBlock().Hash()
	if err != nil {
		log.WithError(err).Error("failed to generate block hash")
		return err
	}

	b := block.New(0, lastBlockHash, bc.txPool)
	nonce := 0
	for b.Mine(nonce, MINING_DIFFICULTY, bc.minerAddress) == nil {
		nonce += 1
	}

	bc.chain = append(bc.chain, b)

	bc.txPool = []*transaction.Transaction{}

	return nil
}

func (bc *Blockchain) AddTxToPool(tx *transaction.Transaction) error {
	if tx.Verify() == false {
		return fmt.Errorf("failed to verify transaction")
	}

	bc.txPool = append(bc.txPool, tx)

	return nil
}

func (bc *Blockchain) GetLastBlock() *block.Block {
	return bc.chain[len(bc.chain)-1]
}

func (bc *Blockchain) Print() {
	for i, b := range bc.chain {
		fmt.Printf("\n Block number  %d \n", i)
		b.Print()
	}
	bc.GetBalances()
}

func (bc *Blockchain) GetBalances() {
	balances := map[string]float64{}
	for _, b := range bc.chain {
		for _, tx := range b.Transactions {
			balances[tx.To.GetAddress()] += tx.Amount
		}
	}

	for k, v := range balances {
		fmt.Printf("%v -> %v\n", k, v)
	}

}
