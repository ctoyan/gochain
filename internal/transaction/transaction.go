package transaction

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
	"time"

	"github.com/ctoyan/gochain/internal/account"
	log "github.com/sirupsen/logrus"
)

type Transaction struct {
	hash      []byte
	From      *account.Account
	To        *account.Account
	Amount    float64
	timestamp int64
	signature *Signature
}

type Signature struct {
	R *big.Int
	S *big.Int
}

func New(from, to *account.Account, amount float64) *Transaction {
	tx := &Transaction{
		From:      from,
		To:        to,
		Amount:    amount,
		timestamp: time.Now().UnixNano(),
	}

	tx.hash = tx.Hash()

	txSignature, err := tx.Sign()
	if err != nil {
		log.WithError(err).Error("failed to generate tx signature")
	}

	tx.signature = txSignature

	return tx
}

func (tx *Transaction) Hash() []byte {
	txHash := sha256.New()
	txHash.Write([]byte(fmt.Sprintf(
		"%v%v%v%v",
		tx.From.GetAddress(),
		tx.To.GetAddress(),
		tx.Amount,
		tx.timestamp,
	)))

	return txHash.Sum(nil)
}

func (tx *Transaction) Verify() bool {
	// rewards to miners are not signed, so no need to verify
	if tx.From.GetAddress() == account.GetRootAccount().GetAddress() {
		return true
	}

	return ecdsa.Verify(tx.From.PublicKey, tx.hash, tx.signature.R, tx.signature.S)
}

func (tx *Transaction) Sign() (*Signature, error) {
	// no need to sign if it's a reward to a miner
	if tx.From.GetAddress() == account.GetRootAccount().GetAddress() {
		return nil, nil
	}

	r, s, err := ecdsa.Sign(rand.Reader, tx.From.PrivateKey, tx.hash)
	if err != nil {
		return nil, err
	}

	return &Signature{
		R: r,
		S: s,
	}, nil
}
