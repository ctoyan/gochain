package account

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"

	"github.com/btcsuite/btcd/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

type Account struct {
	address     string
	totalAmount float64
	PublicKey   *ecdsa.PublicKey
	PrivateKey  *ecdsa.PrivateKey
}

func New() (*Account, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	publicKey := &privateKey.PublicKey

	address, err := generateAddress(privateKey, publicKey)
	if err != nil {
		return nil, err
	}

	return &Account{
		address:     address,
		totalAmount: 0.0,
		PrivateKey:  privateKey,
		PublicKey:   publicKey,
	}, nil
}

func (a *Account) GetTotalAmount() float64 {
	return a.totalAmount
}

func (a *Account) GetAddress() string {
	return a.address
}

func GetRootAccount() *Account {
	return &Account{
		address: "0000000000000000000000000000000000",
	}
}

// generateAddress generates a version one bitcoin address as per this article
// https://en.bitcoin.it/wiki/Technical_background_of_version_1_Bitcoin_addresses
func generateAddress(privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey) (string, error) {
	shaPublic := sha256.New()
	_, err := shaPublic.Write(publicKey.X.Bytes())
	if err != nil {
		return "", err
	}
	shaPublic.Write(publicKey.Y.Bytes())
	if err != nil {
		return "", err
	}
	shaPublicDigest := shaPublic.Sum(nil)

	ripemdPubKeyDigest := ripemd160.New()
	_, err = ripemdPubKeyDigest.Write(shaPublicDigest)
	if err != nil {
		return "", err
	}
	ripemdDigest := ripemdPubKeyDigest.Sum(nil)

	// Add version byte to ripemd (0x00 for Mainnet)
	ripemdDigestWithVersionByte := make([]byte, 21)
	ripemdDigestWithVersionByte[0] = 0x00
	copy(ripemdDigestWithVersionByte[1:], ripemdDigest[:])

	shaRipemd := sha256.New()
	_, err = shaRipemd.Write(ripemdDigestWithVersionByte)
	if err != nil {
		return "", err
	}
	shaRipemdDigest := shaRipemd.Sum(nil)

	shaShaRipemd := sha256.New()
	_, err = shaShaRipemd.Write(shaRipemdDigest)
	if err != nil {
		return "", err
	}
	shaShaRipemdDigest := shaShaRipemd.Sum(nil)

	// Take the first 4 bytes of shaShaRipemdDigest hash for checksum
	checksum := shaShaRipemdDigest[:4]

	// Add that 4 byte checksum to the ripemdDigestWithVersionByte
	ripemdFinal := make([]byte, 25)
	copy(ripemdFinal[:21], ripemdDigestWithVersionByte)
	copy(ripemdFinal[21:], checksum[:])

	// Covert the final ripemd hash to base58
	address := base58.Encode(ripemdFinal)

	return address, nil
}
