/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) DCPS. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package did

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
	"crypto/sha256"
	"did-sdk/invoke"
	"did-sdk/key"
	"did-sdk/proof"
	"did-sdk/utils"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"chainmaker.org/chainmaker/pb-go/v2/common"
	cmsdk "chainmaker.org/chainmaker/sdk-go/v2"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/mr-tron/base58"
	"github.com/tjfoc/gmsm/sm2"

	bcx509 "github.com/liuxinfeng96/bc-crypto/x509"

	"chainmaker.org/chainmaker/common/v2/evmutils"
	"chainmaker.org/chainmaker/did-contract/model"
	bcecdsa "github.com/liuxinfeng96/bc-crypto/ecdsa"
)

const (
	DidPrefix                   = "did"
	DidContext                  = "https://www.w3.org/ns/did/v1"
	VerificationMethodKeySuffix = "#keys-"
)

// GetDidMethodFromChain query contract from chain
// @params client the chainmaker sdk client
// @return string the did method
func GetDidMethodFromChain(client *cmsdk.ChainClient) (string, error) {

	result, err := invoke.QueryContract(invoke.DIDContractName, model.Method_DidMethod, nil, client)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

// GenerateDidByPK did was generated using the public key PEM encoding
// @params pkPem: PK PEM
// @params client: ChainMaker SDK
// @return string: the did string
func GenerateDidByPK(pkPem []byte, client *cmsdk.ChainClient) (string, error) {
	// 从链上获取DID方法名
	didMethod, err := GetDidMethodFromChain(client)
	if err != nil {
		return "", err
	}

	pkHash := sha256Hash(pkPem)
	didSuffix := base58Encode(pkHash)

	did := fmt.Sprintf("%s:%s:%s", DidPrefix, didMethod, didSuffix)

	return did, nil
}

// GenerateDidDoc generate a DID document using a key
// @params keyInfo：密钥信息
// @params client：长安链客户端
// @params controller：父控制器，可变参数
func GenerateDidDoc(keyInfo []*key.KeyInfo, client *cmsdk.ChainClient, controller ...string) ([]byte, error) {

	// 密钥最少一把
	if len(keyInfo) == 0 {
		return nil, errors.New("key info cannot be empty")
	}

	// 通过公钥生成DID字符串
	did, err := GenerateDidByPK(keyInfo[0].PkPEM, client)
	if err != nil {
		return nil, err
	}

	// 验证方法构造
	verificationMethod := make([]*model.VerificationMethod, 0)
	authentication := make([]string, 0)
	controller = append(controller, did)

	for k, v := range keyInfo {
		keyId := did + VerificationMethodKeySuffix + strconv.Itoa(k)

		var vm *model.VerificationMethod
		vm, err = newVerificationMethod(keyId, did, v.PkPEM)
		if err != nil {
			return nil, err
		}

		verificationMethod = append(verificationMethod, vm)

		authentication = append(authentication, keyId)

	}

	created := utils.ISO8601Time(time.Now().Unix())

	doc := &model.DidDocument{
		Context:            DidContext,
		Id:                 did,
		Created:            created,
		Updated:            created,
		VerificationMethod: verificationMethod,
		Authentication:     authentication,
		Controller:         controller,
	}

	docBytes, err := json.Marshal(doc)
	if err != nil {
		return nil, err
	}

	// 压缩JSON
	msg, err := utils.CompactJson(docBytes)
	if err != nil {
		return nil, err
	}

	var proofBytes []byte

	if len(keyInfo) > 1 {

		proofs := make([]*model.Proof, 0)

		for k, v := range keyInfo {
			keyId := did + VerificationMethodKeySuffix + strconv.Itoa(k)

			// 生成证明
			var pf *model.Proof
			pf, err = proof.GenerateProofByKey(v.SkPEM, msg, keyId)
			if err != nil {
				return nil, err
			}
			proofs = append(proofs, pf)
		}

		proofBytes, err = json.Marshal(proofs)
		if err != nil {
			return nil, err
		}

	} else {

		keyId := did + VerificationMethodKeySuffix + "0"

		// 生成证明
		pf, err := proof.GenerateProofByKey(keyInfo[0].SkPEM, msg, keyId)
		if err != nil {
			return nil, err
		}

		proofBytes, err = json.Marshal(pf)
		if err != nil {
			return nil, err
		}

	}

	doc.Proof = proofBytes

	return json.Marshal(doc)
}

// AddDidDocToChain store the DID document on the blockchain
// @params doc：DID文档
// @params client：长安链客户端
func AddDidDocToChain(doc string, client *cmsdk.ChainClient) error {
	params := make([]*common.KeyValuePair, 0)

	params = append(params, &common.KeyValuePair{
		Key:   model.Params_DidDocument,
		Value: []byte(doc),
	})

	_, err := invoke.InvokeContract(invoke.DIDContractName, model.Method_AddDidDocument, params, client)
	if err != nil {
		return err
	}

	return nil
}

// IsValidDidOnChain 判断DID在链上是否有效（格式、是否在黑名单）
// @params did：DID
// @params client：长安链客户端
func IsValidDidOnChain(did string, client *cmsdk.ChainClient) (bool, error) {
	params := make([]*common.KeyValuePair, 0)

	params = append(params, &common.KeyValuePair{
		Key:   model.Params_Did,
		Value: []byte(did),
	})

	result, err := invoke.QueryContract(invoke.DIDContractName, model.Method_IsValidDid, params, client)
	if err != nil {
		return false, err
	}

	if string(result) == "true" {
		return true, nil
	}

	return false, nil
}

// GetDidDocFromChain 通过DID在链上获取DID文档
// @params did：DID
// @params client：长安链客户端
func GetDidDocFromChain(did string, client *cmsdk.ChainClient) ([]byte, error) {
	params := make([]*common.KeyValuePair, 0)

	params = append(params, &common.KeyValuePair{
		Key:   model.Params_Did,
		Value: []byte(did),
	})

	resp, err := invoke.QueryContract(invoke.DIDContractName, model.Method_GetDidDocument, params, client)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GetDidByPkFromChain 通过PK获取DID
// @params pkPem：公钥的PEM编码
// @params client：长安链客户端
func GetDidByPkFromChain(pkPem string, client *cmsdk.ChainClient) (string, error) {
	params := make([]*common.KeyValuePair, 0)

	params = append(params, &common.KeyValuePair{
		Key:   model.Params_DidPubkey,
		Value: []byte(pkPem),
	})

	resp, err := invoke.QueryContract(invoke.DIDContractName, model.Method_GetDidByPubKey, params, client)
	if err != nil {
		return "", err
	}

	return string(resp), nil
}

// GetDidByAddressFromChain 通过Address获取DID
// @params address：公钥的PEM编码
// @params client：长安链客户端
func GetDidByAddressFromChain(address string, client *cmsdk.ChainClient) (string, error) {
	params := make([]*common.KeyValuePair, 0)

	params = append(params, &common.KeyValuePair{
		Key:   model.Params_DidAddress,
		Value: []byte(address),
	})

	resp, err := invoke.QueryContract(invoke.DIDContractName, model.Method_GetDidByAddress, params, client)
	if err != nil {
		return "", err
	}

	return string(resp), nil
}

// UpdateDidDocToChain 在链上更新DID文档
// @params doc：DID文档
// @params client：长安链客户端
func UpdateDidDocToChain(doc string, client *cmsdk.ChainClient) error {
	params := make([]*common.KeyValuePair, 0)

	params = append(params, &common.KeyValuePair{
		Key:   model.Params_DidDocument,
		Value: []byte(doc),
	})

	_, err := invoke.InvokeContract(invoke.DIDContractName, model.Method_UpdateDidDocument, params, client)
	if err != nil {
		return err
	}

	return nil
}

// UpdateDidDoc 更新DID文档（本地生成）
// @params oldDoc：老的DID文档
// @params keyInfo：密钥信息
// @params controller：父控制器，可变参数
func UpdateDidDoc(oldDoc model.DidDocument, keyInfo []*key.KeyInfo, controller ...string) ([]byte, error) {

	var newDoc model.DidDocument

	newDoc.Authentication = oldDoc.Authentication
	newDoc.Context = oldDoc.Context
	newDoc.Controller = oldDoc.Controller
	newDoc.Created = oldDoc.Created
	newDoc.Updated = oldDoc.Updated
	newDoc.Id = oldDoc.Id
	newDoc.VerificationMethod = oldDoc.VerificationMethod
	newDoc.Service = oldDoc.Service

	if len(keyInfo) != 0 {

		verificationMethod := make([]*model.VerificationMethod, 0)
		authentication := make([]string, 0)

		if len(controller) != 0 {
			newDoc.Controller = append(newDoc.Controller, controller...)
		}

		// 多密钥生成VerificationMethod
		for k, v := range keyInfo {
			keyId := newDoc.Id + VerificationMethodKeySuffix + strconv.Itoa(k)

			vm, err := newVerificationMethod(keyId, newDoc.Id, v.PkPEM)
			if err != nil {
				return nil, err
			}

			verificationMethod = append(verificationMethod, vm)

			authentication = append(authentication, keyId)

		}

		newDoc.Authentication = authentication
		newDoc.VerificationMethod = verificationMethod
	}

	updated := utils.ISO8601Time(time.Now().Unix())

	newDoc.Updated = updated

	docBytes, err := json.Marshal(newDoc)
	if err != nil {
		return nil, err
	}

	msg, err := utils.CompactJson(docBytes)
	if err != nil {
		return nil, err
	}

	var proofBytes []byte

	if len(keyInfo) > 1 {

		proofs := make([]*model.Proof, 0)

		for k, v := range keyInfo {
			keyId := newDoc.Id + VerificationMethodKeySuffix + strconv.Itoa(k)

			var pf *model.Proof
			pf, err = proof.GenerateProofByKey(v.SkPEM, msg, keyId)
			if err != nil {
				return nil, err
			}
			proofs = append(proofs, pf)
		}

		proofBytes, err = json.Marshal(proofs)
		if err != nil {
			return nil, err
		}

	} else {

		keyId := newDoc.Id + VerificationMethodKeySuffix + "0"

		pf, err := proof.GenerateProofByKey(keyInfo[0].SkPEM, msg, keyId)
		if err != nil {
			return nil, err
		}

		proofBytes, err = json.Marshal(pf)
		if err != nil {
			return nil, err
		}

	}

	newDoc.Proof = proofBytes

	return json.Marshal(newDoc)
}

func sha256Hash(str []byte) []byte {
	hash := sha256.Sum256(str)
	return hash[:]
}

func base58Encode(hash []byte) string {
	encoded := base58.Encode(hash)
	return encoded
}

func newVerificationMethod(id, controller string, pkPem []byte) (*model.VerificationMethod, error) {

	// 校验是否是公钥
	pubKey, err := bcx509.ParsePublicKey(pkPem)
	if err != nil {
		return nil, err
	}

	var algorithm string

	switch pk := pubKey.(type) {

	case *bcecdsa.PublicKey:

		switch pk.Curve {
		case elliptic.P224(), elliptic.P256(), secp256k1.S256():
			algorithm = model.ECDSAWithSHA256
		case elliptic.P384():
			algorithm = model.ECDSAWithSHA384
		case elliptic.P521():
			algorithm = model.ECDSAWithSHA512
		case sm2.P256Sm2():
			algorithm = model.SM2WithSM3
		default:
			return nil, errors.New("x509: unknown elliptic curve")
		}

	case *ecdsa.PrivateKey:

		switch pk.Curve {
		case elliptic.P224(), elliptic.P256(), secp256k1.S256():
			algorithm = model.ECDSAWithSHA256
		case elliptic.P384():
			algorithm = model.ECDSAWithSHA384
		case elliptic.P521():
			algorithm = model.ECDSAWithSHA512
		case sm2.P256Sm2():
			algorithm = model.SM2WithSM3
		default:
			return nil, errors.New("x509: unknown elliptic curve")
		}

	case *rsa.PrivateKey:

		algorithm = model.SHA256WithRSA

	default:
		return nil, errors.New("unknown publicKey algorithm")
	}

	ski, err := bcx509.ComputeSKI(pubKey)
	if err != nil {
		return nil, err
	}

	// 根据长安链规则, 计算出Address
	bytesAddr := evmutils.Keccak256(ski)
	addr := hex.EncodeToString(bytesAddr)[24:]

	return &model.VerificationMethod{
		Id:           id,
		Type:         algorithm,
		Controller:   controller,
		PublicKeyPem: string(pkPem),
		Address:      addr,
	}, nil

}
