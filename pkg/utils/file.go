package utils

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"io"
	"os"

	"github.com/pkg/errors"
)

func ReadFile(f string) ([]byte, error) {
	kf, err := os.Open(f)
	if err != nil {
		return nil, errors.Wrapf(err, "read %s error", f)
	}

	return io.ReadAll(kf)
}

func ReadECDSAPublicKeyFile(f string) (*ecdsa.PublicKey, error) {
	key, err := ReadFile(f)
	if err != nil {
		return nil, err
	}

	// Parse PEM block
	var block *pem.Block
	if block, _ = pem.Decode(key); block == nil {
		return nil, errors.New("invalid PEM file format")
	}

	// Parse the key
	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKIXPublicKey(block.Bytes); err != nil {
		if cert, err := x509.ParseCertificate(block.Bytes); err == nil {
			parsedKey = cert.PublicKey
		} else {
			return nil, errors.WithStack(err)
		}
	}

	var pkey *ecdsa.PublicKey
	var ok bool
	if pkey, ok = parsedKey.(*ecdsa.PublicKey); !ok {
		return nil, errors.New("invalid ECDSA public key")
	}

	return pkey, nil
}

func ReadECDSAPrivateKey(f string) (*ecdsa.PrivateKey, error) {
	key, err := ReadFile(f)
	if err != nil {
		return nil, err
	}

	// Parse PEM block
	var block *pem.Block
	if block, _ = pem.Decode(key); block == nil {
		return nil, errors.New("invalid PEM file format")
	}

	// Parse the key
	var parsedKey interface{}
	if parsedKey, err = x509.ParseECPrivateKey(block.Bytes); err != nil {
		if parsedKey, err = x509.ParsePKCS8PrivateKey(block.Bytes); err != nil {
			return nil, err
		}
	}

	var pkey *ecdsa.PrivateKey
	var ok bool
	if pkey, ok = parsedKey.(*ecdsa.PrivateKey); !ok {
		return nil, errors.New("invalid ECDSA public key")
	}

	return pkey, nil
}
