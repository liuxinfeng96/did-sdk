/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) DCPS. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package proof

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"did-sdk/utils"
	"encoding/base64"
	"errors"
	"time"

	"chainmaker.org/chainmaker/did-contract/model"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	bccrypto "github.com/liuxinfeng96/bc-crypto"
	bcecdsa "github.com/liuxinfeng96/bc-crypto/ecdsa"
	bcx509 "github.com/liuxinfeng96/bc-crypto/x509"
	"github.com/tjfoc/gmsm/sm2"
)

// GenerateProofByKey 通过私钥生成证明
// @params skPem：私钥的PEM编码
// @params msg：签名的信息
// @params verificationMethod did中的验证方法，通常是`[DID]#key-[i]`格式
func GenerateProofByKey(skPem, msg []byte, verificationMethod string) (*model.Proof, error) {

	// 使用bcx509包里的解析密钥方法，反序列化密钥，不采用[chainmaker common]包是为了支持Secp256k1公钥算法
	privateKey, err := bcx509.ParsePrivateKey(skPem)
	if err != nil {
		return nil, err
	}

	privKey, ok := privateKey.(crypto.Signer)
	if !ok {
		return nil, errors.New("private key does not implement crypto.Signer")
	}

	var (
		hashFunc  bccrypto.Hash
		algorithm string
	)

	switch sk := privateKey.(type) {

	case *bcecdsa.PrivateKey:

		switch sk.Curve {
		case elliptic.P224(), elliptic.P256(), secp256k1.S256():
			hashFunc = bccrypto.SHA256
			algorithm = model.ECDSAWithSHA256
		case elliptic.P384():
			hashFunc = bccrypto.SHA384
			algorithm = model.ECDSAWithSHA384
		case elliptic.P521():
			hashFunc = bccrypto.SHA512
			algorithm = model.ECDSAWithSHA512
		case sm2.P256Sm2():
			hashFunc = bccrypto.SM3
			algorithm = model.SM2WithSM3
		default:
			return nil, errors.New("x509: unknown elliptic curve")
		}

	case *ecdsa.PrivateKey:

		switch sk.Curve {
		case elliptic.P224(), elliptic.P256(), secp256k1.S256():
			hashFunc = bccrypto.SHA256
			algorithm = model.ECDSAWithSHA256
		case elliptic.P384():
			hashFunc = bccrypto.SHA384
			algorithm = model.ECDSAWithSHA384
		case elliptic.P521():
			hashFunc = bccrypto.SHA512
			algorithm = model.ECDSAWithSHA512
		case sm2.P256Sm2():
			hashFunc = bccrypto.SM3
			algorithm = model.SM2WithSM3
		default:
			return nil, errors.New("x509: unknown elliptic curve")
		}

	case *rsa.PrivateKey:

		hashFunc = bccrypto.SHA256
		algorithm = model.SHA256WithRSA

	default:
		return nil, errors.New("unknown publicKey algorithm")
	}

	// 国密算法哈希摘要在其签名里实现
	if hashFunc != bccrypto.SM3 {
		h := hashFunc.New()
		h.Write(msg)
		msg = h.Sum(nil)
	}

	var signerOpts crypto.SignerOpts = hashFunc

	// 对传入的信息进行签名
	signature, err := privKey.Sign(rand.Reader, msg, signerOpts)
	if err != nil {
		return nil, err
	}

	created := utils.ISO8601Time(time.Now().Unix())

	signBase64 := base64.StdEncoding.EncodeToString(signature)

	return &model.Proof{
		Type:               algorithm,
		Created:            created,
		ProofPurpose:       "assertionMethod",
		VerificationMethod: verificationMethod,
		ProofValue:         signBase64,
	}, nil
}

// VerifyPKProof 通过公钥验证证明
// @params msg：证明的信息
// @params pkPem：公钥的PEM编码格式
// @params proof：证明的结构（引自合约）
// @return bool 验证结果
func VerifyPKProof(msg, pkPem []byte, proof *model.Proof) (bool, error) {
	// 直接使用合约里的验证方法
	return proof.Verify(msg, pkPem)
}
