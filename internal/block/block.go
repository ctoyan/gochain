package block

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"strings"
	"time"

	"github.com/ctoyan/gochain/internal/account"
	"github.com/ctoyan/gochain/internal/transaction"
	log "github.com/sirupsen/logrus"
)

type Block struct {
	hash          [32]byte
	nonce         int
	Transactions  []*transaction.Transaction
	prevBlockHash [32]byte
	timestamp     int64
}

func New(nonce int, previousBlockHash [32]byte, transactions []*transaction.Transaction) *Block {
	return &Block{
		timestamp:     time.Now().UnixNano(),
		prevBlockHash: previousBlockHash,
		nonce:         nonce,
		Transactions:  transactions,
	}
}

// Hash returns the hash of a block
// Includes block nonce, timestamp, previous block hash and all transactions hash
func (b *Block) Hash() ([32]byte, error) {
	var strToHash bytes.Buffer
	for _, t := range b.Transactions {
		_, err := strToHash.WriteString(fmt.Sprintf("%v", t.Hash()))
		if err != nil {
			return [32]byte{}, err
		}
	}

	_, err := strToHash.WriteString(fmt.Sprintf("%v%v%v", b.nonce, b.timestamp, b.prevBlockHash))
	if err != nil {
		return [32]byte{}, err
	}

	hash := sha256.New()
	hash.Write(strToHash.Bytes())
	byteSum := (*[32]byte)(hash.Sum(nil))

	return *byteSum, nil
}

func (b *Block) Mine(nonce, difficulty int, rewardRecipient *account.Account) *Block {
	b.nonce = nonce
	zeros := strings.Repeat("0", difficulty)

	blockHash, err := b.Hash()
	if err != nil {
		log.WithError(err).Error("failed to generate block hash")
		return nil
	}

	guessHash := fmt.Sprintf("%x", blockHash)

	if guessHash[:difficulty] == zeros {
		b.Transactions = append(b.Transactions, transaction.New(account.GetRootAccount(), rewardRecipient, 1))
		b.hash, err = b.Hash()
		if err != nil {
			log.WithError(err).Error("failed to generate block hash")
			return nil
		}

		return b
	}

	return nil
}

func (b *Block) Print() {
	fmt.Printf("hash:             %x\n", b.hash)
	fmt.Printf("timestamp:        %v\n", b.timestamp)
	fmt.Printf("nonce:            %v\n", b.nonce)
	fmt.Printf("prevBlockHash:    %x\n", b.prevBlockHash)
	fmt.Printf("transactions:")
	for _, t := range b.Transactions {
		fmt.Printf("\t%v -> %v - %v\n", t.From.GetAddress(), t.To.GetAddress(), t.Amount)
	}
}
