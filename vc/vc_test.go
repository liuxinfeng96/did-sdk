package vc

import (
	"did-sdk/did"
	"did-sdk/key"
	"did-sdk/testdata"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"chainmaker.org/chainmaker/did-contract/model"
	"github.com/test-go/testify/require"
)

func TestIssueVCLocal(t *testing.T) {
	fieldsMap := make(map[string]string)

	fieldsMap["name"] = "姓名"
	fieldsMap["phoneNumber"] = "手机号"
	fieldsMap["idNumber"] = "身份证号"

	jsonSchema, err := GenerateSimpleVcTemplate(fieldsMap)
	require.Nil(t, err)

	keyInfo, err := key.GenerateKey("SM2")
	require.Nil(t, err)

	subject := make(map[string]interface{})

	subject["name"] = "小明"
	subject["phoneNumber"] = "18705453XXX"
	subject["idNumber"] = "370687199X010200XX"

	e := time.Now().Local().Add(time.Hour * 48).Unix()
	vcBytes, err := IssueVCLocal(keyInfo.SkPEM, keyInfo.Algorithm, 0, subject,
		"did:cmdid:0xadfwfkqwfmkqm", "vc1", e, jsonSchema)
	require.Nil(t, err)

	println(string(vcBytes))

	// 不符合VC模板测试

	subject2 := make(map[string]interface{})

	subject2["name"] = "小明"
	subject2["phoneNumber"] = "18705453XXX"
	subject2["idNumber"] = "370687199X010200XX"
	subject2["test"] = "test"

	e = time.Now().Local().Add(time.Hour * 48).Unix()
	_, err = IssueVCLocal(keyInfo.SkPEM, keyInfo.Algorithm, 0, subject2,
		"did:cmdid:0xadfwfkqwfmkqm", "vc1", e, jsonSchema)
	require.NotNil(t, err)
	println(err.Error())

}

func TestIssueVC(t *testing.T) {
	c, err := testdata.GetChainmakerClient()
	require.Nil(t, err)
	// 可信任颁发者上链
	keyInfo, err := key.GenerateKey("EC_Secp256k1")
	require.Nil(t, err)

	doc, err := did.GenerateDidDoc([]*key.KeyInfo{keyInfo}, c)
	require.Nil(t, err)

	err = did.AddDidDocToChain(string(doc), c)
	require.Nil(t, err)

	var document model.DidDocument

	err = json.Unmarshal(doc, &document)
	require.Nil(t, err)

	err = did.AddTrustIssuerListToChain([]string{document.Id}, c)
	require.Nil(t, err)

	fieldsMap := make(map[string]string)

	fieldsMap["name"] = "姓名"
	fieldsMap["phoneNumber"] = "手机号"
	fieldsMap["idNumber"] = "身份证号"

	jsonSchema, err := GenerateSimpleVcTemplate(fieldsMap)
	require.Nil(t, err)

	err = AddVcTemplateToChain("abc123", "模板1", "1.0.0", jsonSchema, c)
	require.Nil(t, err)

	sub := make(map[string]interface{})
	sub["name"] = "XiaoMing"
	sub["phoneNumber"] = "18700001111"
	sub["idNumber"] = "37068711112222000"

	e := time.Now().Add(time.Hour * 24 * 365).Unix()

	vcBytes, err := IssueVC(keyInfo, 0, sub, c, "vc_1111", e, "abc123")
	require.Nil(t, err)
	fmt.Println(string(vcBytes))
}

func TestRevokeVCOnChain(t *testing.T) {
	c, err := testdata.GetChainmakerClient()
	require.Nil(t, err)

	err = RevokeVCOnChain("vc_1111", c)
	require.Nil(t, err)
}

func TestGetVCRevokedListFromChain(t *testing.T) {
	c, err := testdata.GetChainmakerClient()
	require.Nil(t, err)

	err = RevokeVCOnChain("vc_1111", c)
	require.Nil(t, err)

	list, err := GetVCRevokedListFromChain("", 0, 0, c)
	require.Nil(t, err)
	require.Equal(t, []string{"vc_1111"}, list)
}

func TestVerifyVCOnChain(t *testing.T) {
	c, err := testdata.GetChainmakerClient()
	require.Nil(t, err)
	// 可信任颁发者上链
	keyInfo, err := key.GenerateKey("EC_Secp256k1")
	require.Nil(t, err)

	doc, err := did.GenerateDidDoc([]*key.KeyInfo{keyInfo}, c)
	require.Nil(t, err)
	fmt.Println(string(doc))

	err = did.AddDidDocToChain(string(doc), c)
	require.Nil(t, err)

	var document model.DidDocument

	json.Unmarshal(doc, &document)
	require.Nil(t, err)

	err = did.AddTrustIssuerListToChain([]string{document.Id}, c)
	require.Nil(t, err)

	fieldsMap := make(map[string]string)

	fieldsMap["id"] = "DID"
	fieldsMap["name"] = "姓名"
	fieldsMap["phoneNumber"] = "手机号"
	fieldsMap["idNumber"] = "身份证号"

	jsonSchema, err := GenerateSimpleVcTemplate(fieldsMap)
	require.Nil(t, err)

	err = AddVcTemplateToChain("abc12345", "模板1", "1.0.0", jsonSchema, c)
	require.Nil(t, err)

	sub := make(map[string]interface{})
	sub["id"] = "did:cm:test"
	sub["name"] = "XiaoMing"
	sub["phoneNumber"] = "18700001111"
	sub["idNumber"] = "37068711112222000"

	e := time.Now().Add(time.Hour * 24 * 365).Unix()

	vcBytes, err := IssueVC(keyInfo, 0, sub, c, "vc_2222", e, "abc12345")
	require.Nil(t, err)
	fmt.Println(string(vcBytes))

	ok, err := VerifyVCOnChain(string(vcBytes), c)
	require.Nil(t, err)
	require.Equal(t, true, ok)
}
