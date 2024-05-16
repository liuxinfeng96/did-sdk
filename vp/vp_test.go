/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) DCPS. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package vp

import (
	"did-sdk/did"
	"did-sdk/key"
	"did-sdk/vc"
	"encoding/json"
	"testing"
	"time"

	"did-sdk/testdata"

	"chainmaker.org/chainmaker/did-contract/model"
	"github.com/test-go/testify/require"
)

func TestGenerateVP(t *testing.T) {
	fieldsMap := make(map[string]string)

	fieldsMap["name"] = "姓名6"
	fieldsMap["phoneNumber"] = "手机号6"
	fieldsMap["idNumber"] = "身份证号6"

	jsonSchema, err := vc.GenerateSimpleVcTemplate(fieldsMap)
	require.Nil(t, err)

	keyInfo, err := key.GenerateKey("SM2")
	require.Nil(t, err)

	subject := make(map[string]interface{})

	subject["name"] = "小明"
	subject["phoneNumber"] = "18705453XXX"
	subject["idNumber"] = "370687199X010200XX"
	subject["id"] = "did:cm:test1"

	e := time.Now().Local().Add(time.Hour * 48).Unix()
	vcBytes, err := vc.IssueVCLocal(keyInfo.SkPEM, 0, subject,
		"did:cmdid:0xadfwfkqwfmkqm", "vc1", e, jsonSchema)
	require.Nil(t, err)

	keyInfo2, err := key.GenerateKey("EC_Secp256k1")
	require.Nil(t, err)

	vpBytes, err := GenerateVP(keyInfo2.SkPEM, 0, "did:cmdid:vpholder", "vp1", []string{string(vcBytes)})
	require.Nil(t, err)

	println(string(vpBytes))
}

func TestVerifyVPOnChain(t *testing.T) {
	// VP验证完整流程测试
	// 生成签发者
	c, err := testdata.GetChainmakerClient(testdata.ConfigPath1)
	require.Nil(t, err)

	// 签发者密钥生成
	issuerKey, err := key.GenerateKey("EC_Secp256k1")
	require.Nil(t, err)

	// 签发者DID文档生成
	issuerDocJson, err := did.GenerateDidDoc([]*key.KeyInfo{issuerKey}, c)
	require.Nil(t, err)

	// 签发者DID文档上链
	err = did.AddDidDocToChain(string(issuerDocJson), c)
	require.Nil(t, err)

	var issuerDoc model.DidDocument

	err = json.Unmarshal(issuerDocJson, &issuerDoc)
	require.Nil(t, err)

	err = did.AddTrustIssuerListToChain([]string{issuerDoc.Id}, c)
	require.Nil(t, err)

	// 被签发者密钥生成
	holderKey, err := key.GenerateKey("SM2")
	require.Nil(t, err)

	holderDocBytes, err := did.GenerateDidDoc([]*key.KeyInfo{holderKey}, c)
	require.Nil(t, err)

	var holderDoc model.DidDocument
	err = json.Unmarshal(holderDocBytes, &holderDoc)
	require.Nil(t, err)

	// 被签发者DID文档上链
	err = did.AddDidDocToChain(string(holderDocBytes), c)
	require.Nil(t, err)

	fieldsMap := make(map[string]string)

	fieldsMap["id"] = "被签发者"
	fieldsMap["name"] = "姓名"
	fieldsMap["phoneNumber"] = "手机号"
	fieldsMap["idNumber"] = "身份证号"

	// 生成VC模板
	jsonSchema, err := vc.GenerateSimpleVcTemplate(fieldsMap)
	require.Nil(t, err)

	// 将VC模板添加到链上
	err = vc.AddVcTemplateToChain("template001", "模板1", "1.0.0", jsonSchema, c)
	require.Nil(t, err)

	sub := make(map[string]interface{})
	sub["id"] = holderDoc.Id
	sub["name"] = "XiaoMing"
	sub["phoneNumber"] = "18700001111"
	sub["idNumber"] = "3706871996060000XX"

	e := time.Now().Add(time.Hour * 24 * 365).Unix()

	// 颁发VC
	vcBytes, err := vc.IssueVC(issuerKey.SkPEM, issuerKey.PkPEM, 0, sub, c, "vc_001", e, "template001")
	require.Nil(t, err)

	// 被签发者生成VP
	vpBytes, err := GenerateVP(holderKey.SkPEM, 0, holderDoc.Id, "vp_001", []string{string(vcBytes)})
	require.Nil(t, err)

	// 链上验证VP
	ok, err := VerifyVPOnChain(string(vpBytes), c)
	require.Nil(t, err)
	require.Equal(t, ok, true)
}
