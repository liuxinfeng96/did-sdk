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

	var hash string

	if p.Type == "SM2" {
		hash = "SM3"
	} else {
		hash = "SHA-256"
	}

	// 哈希类型转换
	cryptoHash := hashStringToHashType(hash)

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
	signature, err := base64.StdEncoding.DecodeString(p.ProofPurpose)
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

func hashStringToHashType(h string) bccrypto.Hash {
	switch h {
	case "MD4":
		return bccrypto.MD4
	case "MD5":
		return bccrypto.MD5
	case "SHA-1":
		return bccrypto.SHA1
	case "SHA-224":
		return bccrypto.SHA224
	case "SHA-256":
		return bccrypto.SHA256
	case "SHA-384":
		return bccrypto.SHA384
	case "SHA-512":
		return bccrypto.SHA512
	case "MD5+SHA1":
		return bccrypto.MD5SHA1
	case "RIPEMD-160":
		return bccrypto.RIPEMD160
	case "SHA3-224":
		return bccrypto.SHA3_224
	case "SHA3-256":
		return bccrypto.SHA3_256
	case "SHA3-384":
		return bccrypto.SHA3_384
	case "SHA3-512":
		return bccrypto.SHA3_512
	case "SHA-512/224":
		return bccrypto.SHA512_224
	case "SHA-512/256":
		return bccrypto.SHA512 / 256
	case "BLAKE2s-256":
		return bccrypto.BLAKE2s_256
	case "BLAKE2b-256":
		return bccrypto.BLAKE2b_256
	case "BLAKE2b-384":
		return bccrypto.BLAKE2b_384
	case "BLAKE2b-512":
		return bccrypto.BLAKE2b_512
	case "SM3":
		return bccrypto.SM3
	default:
		return bccrypto.Hash(0)
	}
}

// compactJson 压缩json字符串，去掉空格换行等
func CompactJson(raw []byte) ([]byte, error) {
	var buf bytes.Buffer
	err := json.Compact(&buf, raw)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
