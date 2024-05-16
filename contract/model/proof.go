/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) DCPS. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package model

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"

	bccrypto "github.com/liuxinfeng96/bc-crypto"
	bcecdsa "github.com/liuxinfeng96/bc-crypto/ecdsa"
	bcx509 "github.com/liuxinfeng96/bc-crypto/x509"
	"github.com/square/go-jose"
)

// Proof DID文档或者凭证的证明
type Proof struct {
	Type               string `json:"type"`
	Created            string `json:"created"`
	ProofPurpose       string `json:"proofPurpose"`
	Challenge          string `json:"challenge,omitempty"`
	VerificationMethod string `json:"verificationMethod"`
	Jws                string `json:"jws,omitempty"`
	ProofValue         string `json:"proofValue,omitempty"`
}

// Verify 证明的验证
func (p *Proof) Verify(msg, pkPem []byte) (bool, error) {
	// 使用bcx509包里的解析公钥方法，反序列化公钥，不采用[chainmaker common]包是为了支持Secp256k1公钥算法
	publicKey, err := bcx509.ParsePublicKey(pkPem)
	if err != nil {
		return false, err
	}

	//如果是JWS签名
	if len(p.Jws) > 0 {
		// 解析JWS
		jwsObject, err := jose.ParseSigned(p.Jws)
		if err != nil {
			return false, err
		}
		// 验证JWS
		output, err := jwsObject.Verify(publicKey)
		if err != nil {
			return false, err
		}
		// 比较解析出的数据和原始数据是否一致
		return bytes.Equal(msg, output), nil
	}

	// 哈希类型转换
	cryptoHash, err := GetHashType(p.Type)
	if err != nil {
		return false, err
	}

	switch cryptoHash {
	case bccrypto.SM3:
		// 国密的哈希在验签里实现
		break
	default:
		if !cryptoHash.Available() {
			return false, errors.New("cannot verify signature: algorithm unimplemented")
		}
		h := cryptoHash.New()
		h.Write(msg)
		msg = h.Sum(nil)
	}

	// 将签名内容base64解码
	signature, err := base64.StdEncoding.DecodeString(p.ProofValue)
	if err != nil {
		return false, err
	}

	// 判断ParsePublicKey解析的公钥类型，使用不同类型的公钥算法验签
	switch pub := publicKey.(type) {
	case *rsa.PublicKey:

		err := rsa.VerifyPKCS1v15(pub, crypto.Hash(cryptoHash), msg, signature)
		if err != nil {
			return false, err
		}

	// 该类型在bcecdsa包里比标准库增加实现了Secp256k1、SM2的验签
	case *bcecdsa.PublicKey:

		if !bcecdsa.VerifyASN1(pub, msg, signature) {
			return false, errors.New("x509: BCECDSA verification failure")
		}

	case *ecdsa.PublicKey:

		if !ecdsa.VerifyASN1(pub, msg, signature) {
			return false, errors.New("x509: ECDSA verification failure")
		}

	default:
		return false, errors.New("cannot verify signature: algorithm unimplemented")
	}

	return true, nil
}

// CompactJson 压缩json字符串，去掉空格换行等
func CompactJson(raw []byte) ([]byte, error) {
	var buf bytes.Buffer
	err := json.Compact(&buf, raw)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
