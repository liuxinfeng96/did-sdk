/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) DCPS. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/
package did

import (
	"did-sdk/key"
	"did-sdk/testdata"
	"encoding/json"
	"fmt"
	"testing"

	"chainmaker.org/chainmaker/did-contract/model"
	"github.com/test-go/testify/require"
)

func TestGetDidMethodFromChain(t *testing.T) {

	// 获取测试长安链客户端
	c, err := testdata.GetChainmakerClient(testdata.ConfigPath1)
	require.Nil(t, err)

	// 链上获取DID Method
	method, err := GetDidMethodFromChain(c)
	require.Nil(t, err)

	fmt.Println(method)
}

func TestGenerateDidDoc(t *testing.T) {
	keyInfo, err := key.GenerateKey("SM2")
	require.Nil(t, err)

	c, err := testdata.GetChainmakerClient(testdata.ConfigPath1)
	require.Nil(t, err)

	doc, err := GenerateDidDoc([]*key.KeyInfo{keyInfo}, c)
	require.Nil(t, err)

	fmt.Println(string(doc))
}

func TestAddDidDocToChain(t *testing.T) {
	keyInfo, err := key.GenerateKey("SM2")
	require.Nil(t, err)

	c, err := testdata.GetChainmakerClient(testdata.ConfigPath1)
	require.Nil(t, err)

	doc, err := GenerateDidDoc([]*key.KeyInfo{keyInfo}, c)
	require.Nil(t, err)

	fmt.Println(string(doc))

	err = AddDidDocToChain(string(doc), c)
	require.Nil(t, err)
}

func TestIsValidDidOnChain(t *testing.T) {
	c, err := testdata.GetChainmakerClient(testdata.ConfigPath1)
	require.Nil(t, err)

	keyInfo, err := key.GenerateKey("SM2")
	require.Nil(t, err)

	did, err := GenerateDidByPK(keyInfo.PkPEM, c)
	require.Nil(t, err)

	ok, err := IsValidDidOnChain(did, c)
	require.Nil(t, err)
	require.Equal(t, true, ok)
}

func TestGetDidDocFromChain(t *testing.T) {
	c, err := testdata.GetChainmakerClient(testdata.ConfigPath1)
	require.Nil(t, err)

	keyInfo, err := key.GenerateKey("EC_Secp256k1")
	require.Nil(t, err)

	doc, err := GenerateDidDoc([]*key.KeyInfo{keyInfo}, c)
	require.Nil(t, err)

	err = AddDidDocToChain(string(doc), c)
	require.Nil(t, err)

	var document model.DidDocument

	json.Unmarshal(doc, &document)
	require.Nil(t, err)

	getDoc, err := GetDidDocFromChain(document.Id, c)
	require.Nil(t, err)

	require.Equal(t, doc, getDoc)
}

func TestGetDidByAddressFromChain(t *testing.T) {
	c, err := testdata.GetChainmakerClient(testdata.ConfigPath1)
	require.Nil(t, err)

	keyInfo, err := key.GenerateKey("EC_Secp256k1")
	require.Nil(t, err)

	doc, err := GenerateDidDoc([]*key.KeyInfo{keyInfo}, c)
	require.Nil(t, err)

	err = AddDidDocToChain(string(doc), c)
	require.Nil(t, err)

	var document model.DidDocument

	json.Unmarshal(doc, &document)
	require.Nil(t, err)

	did, err := GetDidByAddressFromChain(document.VerificationMethod[0].Address, c)
	require.Nil(t, err)

	require.Equal(t, document.Id, did)
}

func TestGetDidByPkFromChain(t *testing.T) {
	c, err := testdata.GetChainmakerClient(testdata.ConfigPath1)
	require.Nil(t, err)

	keyInfo, err := key.GenerateKey("EC_Secp256k1")
	require.Nil(t, err)

	doc, err := GenerateDidDoc([]*key.KeyInfo{keyInfo}, c)
	require.Nil(t, err)

	err = AddDidDocToChain(string(doc), c)
	require.Nil(t, err)

	var document model.DidDocument

	json.Unmarshal(doc, &document)
	require.Nil(t, err)

	did, err := GetDidByPkFromChain(document.VerificationMethod[0].PublicKeyPem, c)
	require.Nil(t, err)

	require.Equal(t, document.Id, did)
}

func TestUpdateDidDocToChain(t *testing.T) {
	c, err := testdata.GetChainmakerClient(testdata.ConfigPath1)
	require.Nil(t, err)

	c2, err := testdata.GetChainmakerClient(testdata.ConfigPath2)
	require.Nil(t, err)
	c2pk, err := c2.GetPublicKey().String()
	require.Nil(t, err)
	c2sk, err := c2.GetPrivateKey().String()
	require.Nil(t, err)

	c2KeyInfo := &key.KeyInfo{
		PkPEM: []byte(c2pk),
		SkPEM: []byte(c2sk),
	}

	c2doc, err := GenerateDidDoc([]*key.KeyInfo{c2KeyInfo}, c)
	require.Nil(t, err)

	err = AddDidDocToChain(string(c2doc), c)
	require.Nil(t, err)

	var c2Doc model.DidDocument

	json.Unmarshal(c2doc, &c2Doc)
	require.Nil(t, err)

	c2Did, err := GetDidByAddressFromChain(c2Doc.VerificationMethod[0].Address, c)
	require.Nil(t, err)

	keyInfo, err := key.GenerateKey("EC_Secp256k1")
	require.Nil(t, err)

	doc, err := GenerateDidDoc([]*key.KeyInfo{keyInfo}, c, c2Did)
	require.Nil(t, err)

	err = AddDidDocToChain(string(doc), c)
	require.Nil(t, err)

	var oldDoc model.DidDocument

	json.Unmarshal(doc, &oldDoc)
	require.Nil(t, err)

	newKeyInfo, err := key.GenerateKey("SM2")
	require.Nil(t, err)

	newDoc, err := UpdateDidDoc(oldDoc, []*key.KeyInfo{newKeyInfo})
	require.Nil(t, err)

	// 使用c2更新，测试权限逻辑
	err = UpdateDidDocToChain(string(newDoc), c2)
	require.Nil(t, err)

	getDoc, err := GetDidDocFromChain(oldDoc.Id, c)
	require.Nil(t, err)

	require.Equal(t, newDoc, getDoc)

	fmt.Println(string(getDoc))
}
