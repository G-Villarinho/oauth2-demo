package ecdsa

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"

	"github.com/aetheris-lab/aetheris-id/api/configs"
)

type EcdsaKeyPair interface {
	ParseECDSAPrivateKey() (*ecdsa.PrivateKey, error)
	ParseECDSAPublicKey() (*ecdsa.PublicKey, error)
}

type ecdsaKeyPair struct {
	config configs.Environment
}

func NewEcdsaKeyPair(config configs.Environment) EcdsaKeyPair {
	return &ecdsaKeyPair{
		config: config,
	}
}

func (e *ecdsaKeyPair) ParseECDSAPrivateKey() (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(e.config.Key.PrivateKey))
	if block == nil || block.Type != "EC PRIVATE KEY" {
		return nil, errors.New("failed to parse EC private key")
	}

	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func (e *ecdsaKeyPair) ParseECDSAPublicKey() (*ecdsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(e.config.Key.PublicKey))
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("failed to parse EC public key")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	ecdsaPubKey, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("key is not a valid ECDSA public key")
	}

	return ecdsaPubKey, nil
}
