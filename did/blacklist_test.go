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
	"testing"

	"chainmaker.org/chainmaker/did-contract/model"
	"github.com/test-go/testify/require"
)

func TestAddDidBlackListToChain(t *testing.T) {
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

	err = AddDidBlackListToChain([]string{document.Id}, c)
	require.Nil(t, err)

	ok, err := IsValidDidOnChain(document.Id, c)
	require.NotNil(t, err)
	require.Equal(t, false, ok)
}

func TestGetDidBlackListFromChain(t *testing.T) {
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

	err = AddDidBlackListToChain([]string{document.Id}, c)
	require.Nil(t, err)

	list, err := GetDidBlackListFromChain(document.Id, 0, 0, c)
	require.Nil(t, err)

	var isInBlacklist bool
	for _, v := range list {
		if v == document.Id {
			isInBlacklist = true
		}
	}

	require.Equal(t, true, isInBlacklist)
}

func TestDeleteDidBlackListFromChain(t *testing.T) {
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

	err = AddDidBlackListToChain([]string{document.Id}, c)
	require.Nil(t, err)

	ok, err := IsValidDidOnChain(document.Id, c)
	require.NotNil(t, err)
	require.Equal(t, false, ok)

	err = DeleteDidBlackListFromChain([]string{document.Id}, c)
	require.Nil(t, err)

	ok, err = IsValidDidOnChain(document.Id, c)
	require.Nil(t, err)
	require.Equal(t, true, ok)
}
