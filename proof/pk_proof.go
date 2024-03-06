package proof

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"did-sdk/key"
	"did-sdk/utils"
	"encoding/base64"
	"errors"
	"time"

	"chainmaker.org/chainmaker/did-contract/model"
	bcecdsa "github.com/liuxinfeng96/bc-crypto/ecdsa"
	bcx509 "github.com/liuxinfeng96/bc-crypto/x509"
	"github.com/tjfoc/gmsm/sm2"
)

// GenerateProofByKey 通过私钥生成证明
// @params skPem：私钥的PEM编码
// @params msg：签名的信息
// @params verificationMethod did中的验证方法，通常是`[DID]#key-[i]`格式
// @params algorithm：公钥算法（如果为空，可自行解析）
// @params hash：信息做摘要的哈希类型
func GenerateProofByKey(skPem, msg []byte, verificationMethod, algorithm, hash string) (*model.Proof, error) {

	// 使用bcx509包里的解析密钥方法，反序列化密钥，不采用[chainmaker common]包是为了支持Secp256k1公钥算法
	privateKey, err := bcx509.ParsePrivateKey(skPem)
	if err != nil {
		return nil, err
	}

	privKey, ok := privateKey.(crypto.Signer)
	if !ok {
		return nil, errors.New("private key does not implement crypto.Signer")
	}

	cryptoHash := key.HashStringToHashType(hash)

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
				algorithm = key.PkAlgorithmSM2
			}
		} else {
			if setAlgo {
				algorithm = key.PkAlgorithmECDSA
			}
		}

	case *ecdsa.PrivateKey:

		if setAlgo {
			algorithm = key.PkAlgorithmECDSA
		}

	case *rsa.PrivateKey:

		if setAlgo {
			algorithm = key.PkAlgorithmRSA
		}

	default:
		return nil, errors.New("unknown publicKey algorithm")
	}

	// 国密算法哈希摘要在其签名里实现
	if !isSm2 {
		h := cryptoHash.New()
		h.Write(msg)
		msg = h.Sum(nil)
	}

	var signerOpts crypto.SignerOpts = cryptoHash

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
