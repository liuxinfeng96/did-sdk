package proof

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"did-sdk/utils"
	"encoding/base64"
	"errors"
	"time"

	bccrypto "github.com/liuxinfeng96/bc-crypto"
	bcecdsa "github.com/liuxinfeng96/bc-crypto/ecdsa"
	bcx509 "github.com/liuxinfeng96/bc-crypto/x509"
	"github.com/tjfoc/gmsm/sm2"
)

// PkProofJSON the type of proof signed using a public key algorithm
type PkProofJSON struct {
	Type               string `json:"type"`
	Created            string `json:"created"`
	ProofPurpose       string `json:"proofPurpose"`
	VerificationMethod string `json:"verificationMethod"`
	ProofValue         string `json:"proofValue"`
}

// GenerateProofByKey 通过私钥生成证明
// @params skPem：私钥的PEM编码
// @params msg：签名的信息
// @params verificationMethod did中的验证方法，通常指向一个验证方法的详情
// @params algorithm：公钥算法（如果为空，可自行解析）
// @params hash：信息做摘要的哈希类型
// @return PkProofJSON 证明的结构
func GenerateProofByKey(skPem, msg []byte, verificationMethod, algorithm, hash string) (*PkProofJSON, error) {

	// 使用bcx509包里的解析密钥方法，反序列化密钥，不采用[chainmaker common]包是为了支持Secp256k1公钥算法
	privateKey, err := bcx509.ParsePrivateKey(skPem)
	if err != nil {
		return nil, err
	}

	key, ok := privateKey.(crypto.Signer)
	if !ok {
		return nil, errors.New("private key does not implement crypto.Signer")
	}

	cryptoHash := utils.HashStringToHashType(hash)

	var (
		setAlgo, isSm2 bool
	)

	// 如果传入的公钥算法名称为空，通过密钥类型设置算法名称
	if len(algorithm) == 0 {
		setAlgo = true
	}

	switch sk := privateKey.(type) {

	case *bcecdsa.PrivateKey:

		if sk.Curve == sm2.P256Sm2() {
			isSm2 = true
			if setAlgo {
				algorithm = "SM2"
			}
		} else {
			if setAlgo {
				algorithm = "ECDSA"
			}
		}

	case *ecdsa.PrivateKey:

		if setAlgo {
			algorithm = "ECDSA"
		}

	case *rsa.PrivateKey:

		if setAlgo {
			algorithm = "RSA"
		}

	default:
		return nil, errors.New("unknown publicKey algorithm")
	}

	// 国密算法哈希摘要在其签名里作实现
	if !isSm2 {
		h := cryptoHash.New()
		h.Write(msg)
		msg = h.Sum(nil)
	}

	var signerOpts crypto.SignerOpts = cryptoHash

	// 对传入的信息进行签名
	signature, err := key.Sign(rand.Reader, msg, signerOpts)
	if err != nil {
		return nil, err
	}

	created := utils.ISO8601Time(time.Now().Unix())

	signBase64 := base64.StdEncoding.EncodeToString(signature)

	return &PkProofJSON{
		Type:               algorithm,
		Created:            created,
		ProofPurpose:       "assertionMethod",
		VerificationMethod: verificationMethod,
		ProofValue:         signBase64,
	}, nil
}

// VerifyPKProof 通过公钥验证证明
// @params msg：证明的信息
// @params proof：证明的内容(序列化的“PkProofJSON”)
// @params pkPem：公钥的PEM编码格式
// @params hash：信息做摘要的哈希类型
// @return bool 验证结果
func VerifyPKProof(msg, pkPem []byte, proof *PkProofJSON, hash string) (bool, error) {

	// 使用bcx509包里的解析公钥方法，反序列化公钥，不采用[chainmaker common]包是为了支持Secp256k1公钥算法
	publicKey, err := bcx509.ParsePublicKey(pkPem)
	if err != nil {
		return false, err
	}

	// 哈希类型转换
	cryptoHash := utils.HashStringToHashType(hash)

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
	signature, err := base64.StdEncoding.DecodeString(proof.ProofValue)
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
