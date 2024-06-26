/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) DCPS. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package key

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	bcecdsa "github.com/liuxinfeng96/bc-crypto/ecdsa"
	bcx509 "github.com/liuxinfeng96/bc-crypto/x509"
	"github.com/tjfoc/gmsm/sm2"
)

const (
	// PEMPrivateKeyTypeStr a string suffix for the Type field of the PEM Block when used as a private key
	PEMPrivateKeyTypeStr = "PRIVATE KEY"
	// PEMPublicKeyTypeStr a string suffix for the Type field of the PEM Block when used as a public key
	PEMPublicKeyTypeStr = "PUBLIC KEY"
)

// KeyInfo the asymmetric encryption key information
// It includes the public key algorithm name, public key REM code and private key PEM code
type KeyInfo struct {
	// 公钥的PEM编码
	PkPEM []byte
	// 私钥的PEM编码
	SkPEM []byte
}

// SupportAlgorithm list of public key algorithms currently supported by this project
var SupportAlgorithm = []string{
	// GMSM_SM2
	"SM2",
	// ECDSA
	"EC_Secp256k1",
	"EC_NISTP224",
	"EC_NISTP256",
	"EC_NISTP384",
	"EC_NISTP521",
	// RSA
	"RSA2048",
	"RSA3072",
}

// IsSupportAlgorithm check whether the public key algorithm is supported
func IsSupportAlgorithm(algo string) bool {
	for _, v := range SupportAlgorithm {
		if algo == v {
			return true
		}
	}
	return false
}

// GenerateKey the public and private keys are generated by the name of the public key algorithm
// @params algorithm 公钥算法名称，支持的类型请见全局变量“SupportAlgorithm”
// @return *KeyInfo 密钥详细信息，包含私钥和公钥的PEM编码以及公钥算法名称
func GenerateKey(algorithm string) (*KeyInfo, error) {
	switch algorithm {
	case "EC_Secp256k1":
		key, err := bcecdsa.GenerateKey(secp256k1.S256(), rand.Reader)
		if err != nil {
			return nil, err
		}

		return bcEcdsaKeyMarshal(key)
	case "SM2":
		key, err := bcecdsa.GenerateKey(sm2.P256Sm2(), rand.Reader)
		if err != nil {
			return nil, err
		}

		return bcEcdsaKeyMarshal(key)
	case "EC_NISTP224":
		key, err := ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
		if err != nil {
			return nil, err
		}

		return ecdsaKeyMarshal(key)
	case "EC_NISTP256":
		key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			return nil, err
		}

		return ecdsaKeyMarshal(key)
	case "EC_NISTP384":
		key, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
		if err != nil {
			return nil, err
		}

		return ecdsaKeyMarshal(key)
	case "EC_NISTP521":
		key, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
		if err != nil {
			return nil, err
		}

		return ecdsaKeyMarshal(key)
	case "RSA2048":
		key, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return nil, err
		}

		return rsaKeyMarshal(key)
	case "RSA3072":
		key, err := rsa.GenerateKey(rand.Reader, 3072)
		if err != nil {
			return nil, err
		}

		return rsaKeyMarshal(key)
	default:
		return nil, errors.New("the public key algorithm curve is unknown")
	}

}

func spliceSkPEMBlockType(algo string) string {
	return algo + " " + PEMPrivateKeyTypeStr
}

func ecdsaKeyMarshal(key *ecdsa.PrivateKey) (*KeyInfo, error) {
	skDer, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		return nil, err
	}

	skBlock := &pem.Block{
		Type:  spliceSkPEMBlockType("EC"),
		Bytes: skDer,
	}

	skBuf := new(bytes.Buffer)
	if err = pem.Encode(skBuf, skBlock); err != nil {
		return nil, err
	}

	pkDer, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		return nil, err
	}

	pkBlock := &pem.Block{
		Type:  PEMPublicKeyTypeStr,
		Bytes: pkDer,
	}

	pkBuf := new(bytes.Buffer)
	if err = pem.Encode(pkBuf, pkBlock); err != nil {
		return nil, err
	}

	return &KeyInfo{
		SkPEM: skBuf.Bytes(),
		PkPEM: pkBuf.Bytes(),
	}, nil
}

func bcEcdsaKeyMarshal(key *bcecdsa.PrivateKey) (*KeyInfo, error) {
	skDer, err := bcx509.MarshalECPrivateKey(key)
	if err != nil {
		return nil, err
	}

	skBlock := &pem.Block{
		Type:  spliceSkPEMBlockType("EC"),
		Bytes: skDer,
	}

	skBuf := new(bytes.Buffer)
	if err = pem.Encode(skBuf, skBlock); err != nil {
		return nil, err
	}

	pkDer, err := bcx509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		return nil, err
	}

	pkBlock := &pem.Block{
		Type:  PEMPublicKeyTypeStr,
		Bytes: pkDer,
	}

	pkBuf := new(bytes.Buffer)
	if err = pem.Encode(pkBuf, pkBlock); err != nil {
		return nil, err
	}

	return &KeyInfo{
		SkPEM: skBuf.Bytes(),
		PkPEM: pkBuf.Bytes(),
	}, nil
}

func rsaKeyMarshal(key *rsa.PrivateKey) (*KeyInfo, error) {
	skDer, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return nil, err
	}

	skBlock := &pem.Block{
		Type:  spliceSkPEMBlockType("RSA"),
		Bytes: skDer,
	}

	skBuf := new(bytes.Buffer)
	if err = pem.Encode(skBuf, skBlock); err != nil {
		return nil, err
	}

	pkDer := x509.MarshalPKCS1PublicKey(&key.PublicKey)

	pkBlock := &pem.Block{
		Type:  PEMPublicKeyTypeStr,
		Bytes: pkDer,
	}

	pkBuf := new(bytes.Buffer)
	if err := pem.Encode(pkBuf, pkBlock); err != nil {
		return nil, err
	}

	return &KeyInfo{
		SkPEM: skBuf.Bytes(),
		PkPEM: pkBuf.Bytes(),
	}, nil
}
