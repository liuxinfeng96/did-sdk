/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) DCPS. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

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

	fieldsMap["name"] = "姓名4"
	fieldsMap["phoneNumber"] = "手机号4"
	fieldsMap["idNumber"] = "身份证号4"

	jsonSchema, err := GenerateSimpleVcTemplate(fieldsMap)
	require.Nil(t, err)

	keyInfo, err := key.GenerateKey("SM2")
	require.Nil(t, err)

	subject := make(map[string]interface{})

	subject["name"] = "小明1"
	subject["phoneNumber"] = "1870545XXXX"
	subject["idNumber"] = "3706871996010200XX"
	subject["id"] = "did:cm:test66"

	e := time.Now().Local().Add(time.Hour * 48).Unix()
	vcBytes, err := IssueVCLocal(keyInfo.SkPEM, 0, subject,
		"did:cmdid:0xadfwfkqwfmkqm", "vc1", e, jsonSchema)
	require.Nil(t, err)

	println(string(vcBytes))

	// 不符合VC模板测试

	subject2 := make(map[string]interface{})

	subject2["name1"] = "小明2"
	subject2["phoneNumber"] = "18705453XXX"
	subject2["idNumber"] = "370687199X010200XX"
	subject2["id"] = "did:cm:test1"

	e = time.Now().Local().Add(time.Hour * 48).Unix()
	_, err = IssueVCLocal(keyInfo.SkPEM, 0, subject2,
		"did:cmdid:0xadfwfkqwfmkqm", "vc1", e, jsonSchema)
	require.NotNil(t, err)
}

func TestIssueVC(t *testing.T) {
	c, err := testdata.GetChainmakerClient(testdata.ConfigPath1)
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

	fieldsMap["name"] = "姓名5"
	fieldsMap["phoneNumber"] = "手机号5"
	fieldsMap["idNumber"] = "身份证号5"

	jsonSchema, err := GenerateSimpleVcTemplate(fieldsMap)
	require.Nil(t, err)

	err = AddVcTemplateToChain("abc12312", "模板1", "1.0.0", jsonSchema, c)
	require.Nil(t, err)

	// 被签发者上链
	keyInfo2, err := key.GenerateKey("EC_Secp256k1")
	require.Nil(t, err)

	doc2, err := did.GenerateDidDoc([]*key.KeyInfo{keyInfo2}, c)
	require.Nil(t, err)

	err = did.AddDidDocToChain(string(doc2), c)
	require.Nil(t, err)

	var doc2Struct model.DidDocument

	err = json.Unmarshal(doc2, &doc2Struct)
	require.Nil(t, err)

	sub := make(map[string]interface{})
	sub["name"] = "XiaoMing1"
	sub["phoneNumber"] = "18700001111"
	sub["idNumber"] = "37068711112222000"
	sub["id"] = doc2Struct.Id

	e := time.Now().Add(time.Hour * 24 * 365).Unix()

	vcBytes, err := IssueVC(keyInfo.SkPEM, keyInfo.PkPEM, 0, sub, c, "vc_1111", e, "abc12312")
	require.Nil(t, err)
	fmt.Println(string(vcBytes))

	// 获取颁发日志
	list, err := GetVcIssueLogListFromChain("vc_1111", 0, 0, c)
	require.Nil(t, err)
	for _, v := range list {
		fmt.Printf("%+v\n", v)
	}
}

func TestRevokeVCOnChain(t *testing.T) {
	c, err := testdata.GetChainmakerClient(testdata.ConfigPath1)
	require.Nil(t, err)

	err = RevokeVCOnChain("vc_1111", c)
	require.Nil(t, err)

	c2, err := testdata.GetChainmakerClient(testdata.ConfigPath2)
	require.Nil(t, err)

	err = RevokeVCOnChain("vc_1111", c2)
	require.NotNil(t, err)
}

func TestGetVCRevokedListFromChain(t *testing.T) {
	c, err := testdata.GetChainmakerClient(testdata.ConfigPath1)
	require.Nil(t, err)

	err = RevokeVCOnChain("vc_1111", c)
	require.Nil(t, err)

	list, err := GetVCRevokedListFromChain("", 0, 0, c)
	require.Nil(t, err)
	require.Equal(t, []string{"vc_1111"}, list)
}

func TestVerifyVCOnChain(t *testing.T) {
	c, err := testdata.GetChainmakerClient(testdata.ConfigPath1)
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

	// 被签发者上链
	keyInfo2, err := key.GenerateKey("EC_Secp256k1")
	require.Nil(t, err)

	doc2, err := did.GenerateDidDoc([]*key.KeyInfo{keyInfo2}, c)
	require.Nil(t, err)

	err = did.AddDidDocToChain(string(doc2), c)
	require.Nil(t, err)

	var doc2Struct model.DidDocument

	err = json.Unmarshal(doc2, &doc2Struct)
	require.Nil(t, err)

	fieldsMap := make(map[string]string)

	fieldsMap["name"] = "姓名6"
	fieldsMap["phoneNumber"] = "手机号6"
	fieldsMap["idNumber"] = "身份证号6"

	jsonSchema, err := GenerateSimpleVcTemplate(fieldsMap)
	require.Nil(t, err)

	err = AddVcTemplateToChain("abc12345", "模板1", "1.0.0", jsonSchema, c)
	require.Nil(t, err)

	sub := make(map[string]interface{})
	sub["id"] = doc2Struct.Id
	sub["name"] = "XiaoMing"
	sub["phoneNumber"] = "18700001111"
	sub["idNumber"] = "37068711112222000"

	e := time.Now().Add(time.Hour * 24 * 365).Unix()

	vcBytes, err := IssueVC(keyInfo.SkPEM, keyInfo.PkPEM, 0, sub, c, "vc_2222", e, "abc12345")
	require.Nil(t, err)
	fmt.Println(string(vcBytes))

	ok, err := VerifyVCOnChain(string(vcBytes), c)
	require.Nil(t, err)
	require.Equal(t, true, ok)
}
