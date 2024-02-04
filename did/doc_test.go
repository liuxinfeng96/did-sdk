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

	c, err := testdata.GetChainmakerClient()
	require.Nil(t, err)

	method, err := GetDidMethodFromChain(c)
	require.Nil(t, err)

	fmt.Println(method)
}

func TestGenerateDidDoc(t *testing.T) {
	keyInfo, err := key.GenerateKey("SM2")
	require.Nil(t, err)

	c, err := testdata.GetChainmakerClient()
	require.Nil(t, err)

	doc, err := GenerateDidDoc([]*key.KeyInfo{keyInfo}, c)
	require.Nil(t, err)

	fmt.Println(string(doc))
}

func TestAddDidDocToChain(t *testing.T) {
	keyInfo, err := key.GenerateKey("SM2")
	require.Nil(t, err)

	c, err := testdata.GetChainmakerClient()
	require.Nil(t, err)

	doc, err := GenerateDidDoc([]*key.KeyInfo{keyInfo}, c)
	require.Nil(t, err)

	fmt.Println(string(doc))

	err = AddDidDocToChain(string(doc), c)
	require.Nil(t, err)
}

func TestIsValidDidOnChain(t *testing.T) {
	c, err := testdata.GetChainmakerClient()
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
	c, err := testdata.GetChainmakerClient()
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
	c, err := testdata.GetChainmakerClient()
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
	c, err := testdata.GetChainmakerClient()
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
	c, err := testdata.GetChainmakerClient()
	require.Nil(t, err)

	keyInfo, err := key.GenerateKey("EC_Secp256k1")
	require.Nil(t, err)

	doc, err := GenerateDidDoc([]*key.KeyInfo{keyInfo}, c)
	require.Nil(t, err)

	err = AddDidDocToChain(string(doc), c)
	require.Nil(t, err)

	var oldDoc model.DidDocument

	json.Unmarshal(doc, &oldDoc)
	require.Nil(t, err)

	newKeyInfo, err := key.GenerateKey("SM2")
	require.Nil(t, err)

	newDoc, err := UpdateDidDoc(oldDoc, []*key.KeyInfo{newKeyInfo}, "did:cm:admin")
	require.Nil(t, err)

	err = UpdateDidDocToChain(string(newDoc), c)
	require.Nil(t, err)

	getDoc, err := GetDidDocFromChain(oldDoc.Id, c)
	require.Nil(t, err)

	require.Equal(t, newDoc, getDoc)

	fmt.Println(string(getDoc))
}
