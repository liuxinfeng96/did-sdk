package did

import (
	"did-sdk/key"
	"did-sdk/testdata"
	"encoding/json"
	"testing"

	"chainmaker.org/chainmaker/did-contract/model"
	"github.com/test-go/testify/require"
)

func TestAddTrustIssuerListToChain(t *testing.T) {
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

	err = AddTrustIssuerListToChain([]string{document.Id}, c)
	require.Nil(t, err)
}

func TestGetTrustIssuerListFromChain(t *testing.T) {
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

	err = AddTrustIssuerListToChain([]string{document.Id}, c)
	require.Nil(t, err)

	list, err := GetTrustIssuerListFromChain(document.Id, 0, 0, c)
	require.Nil(t, err)
	require.Equal(t, list, []string{document.Id})
}

func TestDeleteTrustIssuerListFromChain(t *testing.T) {
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

	err = AddTrustIssuerListToChain([]string{document.Id}, c)
	require.Nil(t, err)

	list, err := GetTrustIssuerListFromChain(document.Id, 0, 0, c)
	require.Nil(t, err)
	require.Equal(t, list, []string{document.Id})

	err = DeleteTrustIssuerListFromChain([]string{document.Id}, c)
	require.Nil(t, err)

	list2, err := GetTrustIssuerListFromChain(document.Id, 0, 0, c)
	require.Nil(t, err)
	require.NotEqual(t, list2, []string{document.Id})
}
