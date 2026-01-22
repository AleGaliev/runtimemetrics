package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
)

type Crypto struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

func loadPublicKey(path string) (*rsa.PublicKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not an RSA public key")
	}

	return rsaPub, nil
}

func loadPrivateKey(path string) (*rsa.PrivateKey, error) {

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}

	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaPriv, ok := priv.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("not an RSA private key")
	}

	return rsaPriv, nil
}

func NewCrypto(pathPrivateKey, pathPublicKey string) (*Crypto, error) {

	var (
		publicKey  *rsa.PublicKey
		privateKey *rsa.PrivateKey
		err        error
	)

	if pathPublicKey != "" {
		publicKey, err = loadPublicKey(pathPublicKey)
		if err != nil {
			return nil, err
		}
	}
	if pathPrivateKey != "" {
		privateKey, err = loadPrivateKey(pathPrivateKey)
		if err != nil {
			return nil, err
		}
	}
	return &Crypto{
		publicKey:  publicKey,
		privateKey: privateKey,
	}, nil
}

func (c *Crypto) Encrypt(data []byte) ([]byte, error) {
	if c.publicKey == nil {
		return data, nil
	}

	hash := sha256.New()
	encrypted, err := rsa.EncryptOAEP(
		hash,
		rand.Reader,
		c.publicKey,
		data,
		nil,
	)
	if err != nil {
		return nil, err
	}
	return encrypted, nil
}

func (c *Crypto) Decrypt(data []byte) ([]byte, error) {
	if c.privateKey == nil {
		return data, nil
	}

	hash := sha256.New()
	decrypted, err := rsa.DecryptOAEP(
		hash,
		rand.Reader,
		c.privateKey,
		data,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return decrypted, nil
}
