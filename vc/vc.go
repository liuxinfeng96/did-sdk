/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) DCPS. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package vc

import (
	"did-sdk/did"
	"did-sdk/invoke"
	"did-sdk/proof"
	"did-sdk/utils"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"chainmaker.org/chainmaker/did-contract/model"
	"chainmaker.org/chainmaker/pb-go/v2/common"
	cmsdk "chainmaker.org/chainmaker/sdk-go/v2"
	jsonschema "github.com/xeipuuv/gojsonschema"
)

var ContextVC = []string{
	"https://www.w3.org/2018/credentials/v1",
	"https://www.w3.org/2018/credentials/examples/v1",
}

// IssueVC 颁发VC（需要链上校验）
// @params skPem: 私钥的PEM编码
// @params pkPem: 公钥的PEM编码
// @params keyIndex：公钥在DID文档中的索引
// @params subject：颁发信息主体，对应VC中的`credentialSubject`字段
// @params client：长安链客户端
// @params vcId：VC的`id`字段，可以根据业务自定义
// @params expirationDate：VC的到期时间
// @params vcTemplateId：VC的模板Id，在链上获取VC模板
// @params vcType：VC中的`type`字段，描述VC的类型信息（可变参数，默认会填写“VerifiableCredential”,可继续根据业务类型追加）
func IssueVC(skPem, pkPem []byte, keyIndex int, subject map[string]interface{}, client *cmsdk.ChainClient,
	vcId string, expirationDate int64, vcTemplateId string, vcType ...string) ([]byte, error) {

	// 获取sunject中的DID
	d, ok := subject["id"]
	if !ok {
		return nil, errors.New("the id field must be included in the subject")
	}

	didStr, ok := d.(string)
	if !ok {
		return nil, errors.New("the data type of the id is incorrect")
	}

	// 链上获取模板
	vcTemplate, err := GetVcTemplateFromChain(vcTemplateId, client)
	if err != nil {
		return nil, err
	}

	if len(vcTemplate) == 0 {
		return nil, errors.New("vc template not found on chain")
	}

	var template model.VcTemplate
	err = json.Unmarshal(vcTemplate, &template)
	if err != nil {
		return nil, err
	}

	// 验证subject是否符合VC模板规范
	ok, err = verifyCredentialSubject(subject, template.Template)
	if !ok {
		return nil, err
	}

	vcType = append(vcType, "VerifiableCredential")
	issuer, err := did.GenerateDidByPK(pkPem, client)
	if err != nil {
		return nil, err
	}

	issuanceDate := utils.ISO8601Time(time.Now().Unix())
	expirationDateStr := utils.ISO8601Time(expirationDate)

	vc := &model.VerifiableCredential{
		Context:           ContextVC,
		Id:                vcId,
		Type:              vcType,
		CredentialSubject: subject,
		Issuer:            issuer,
		IssuanceDate:      issuanceDate,
		ExpirationDate:    expirationDateStr,
		Template: &struct {
			ID   string "json:\"id\""
			Name string "json:\"name\""
		}{
			ID:   template.Id,
			Name: template.Name,
		},
	}

	vcBytes, err := json.Marshal(vc)
	if err != nil {
		return nil, err
	}

	msg, err := utils.CompactJson(vcBytes)
	if err != nil {
		return nil, err
	}

	keyId := issuer + did.VerificationMethodKeySuffix + strconv.Itoa(keyIndex)
	pf, err := proof.GenerateProofByKey(skPem, msg, keyId)
	if err != nil {
		return nil, err
	}

	vc.Proof = pf

	vcBytesJSON, err := json.Marshal(vc)
	if err != nil {
		return nil, err
	}

	// 在链上生成签发日志（会对Issuer, did, vcTemplate进行校验）
	err = AddVcIssueLogToChain(issuer, didStr, vcId, vcTemplateId, client)
	if err != nil {
		return nil, err
	}

	return vcBytesJSON, nil
}

// IssueVCLocal 本地颁发VC（不经过链上计算和校验）
// @params skPem: 私钥的PEM编码
// @params keyIndex：公钥在DID文档中的索引
// @params subject: 颁发信息主体，对应VC中的`credentialSubject`字段
// @params issuer: 颁发者的DID编号
// @params vcId：VC的`id`字段，可以根据业务自定义
// @params expirationDate：VC的到期时间
// @params vcTemplate：VC的模板内容，是一个JSON schema，一般存储在链上
// @params vcType：VC中的`type`字段，描述VC的类型信息（可变参数，默认会填写“VerifiableCredential”,可继续根据业务类型追加）
func IssueVCLocal(skPem []byte, keyIndex int, subject map[string]interface{}, issuer string,
	vcId string, expirationDate int64, vcTemplate []byte, vcType ...string) ([]byte, error) {
	// 获取sunject中的DID
	d, ok := subject["id"]
	if !ok {
		return nil, errors.New("the id field must be included in the subject")
	}

	_, ok = d.(string)
	if !ok {
		return nil, errors.New("the data type of the id is incorrect")
	}

	// 验证subject是否符合VC模板规范
	ok, err := verifyCredentialSubject(subject, vcTemplate)
	if !ok {
		return nil, err
	}

	vcType = append(vcType, "VerifiableCredential")

	issuanceDate := utils.ISO8601Time(time.Now().Unix())
	expirationDateStr := utils.ISO8601Time(expirationDate)

	vc := &model.VerifiableCredential{
		Context:           ContextVC,
		Id:                vcId,
		Type:              vcType,
		CredentialSubject: subject,
		Issuer:            issuer,
		IssuanceDate:      issuanceDate,
		ExpirationDate:    expirationDateStr,
	}

	vcBytes, err := json.Marshal(vc)
	if err != nil {
		return nil, err
	}

	msg, err := utils.CompactJson(vcBytes)
	if err != nil {
		return nil, err
	}

	keyId := issuer + did.VerificationMethodKeySuffix + strconv.Itoa(keyIndex)
	pf, err := proof.GenerateProofByKey(skPem, msg, keyId)
	if err != nil {
		return nil, err
	}

	vc.Proof = pf

	return json.Marshal(vc)
}

func verifyCredentialSubject(subject map[string]interface{}, vcTemplate []byte) (bool, error) {
	subjectBytes, err := json.Marshal(subject)
	if err != nil {
		return false, err
	}

	schemaLoader := jsonschema.NewBytesLoader(vcTemplate)
	subjectLoader := jsonschema.NewBytesLoader(subjectBytes)

	result, err := jsonschema.Validate(schemaLoader, subjectLoader)
	if err != nil {
		return false, err
	}

	if result.Valid() {
		return true, nil
	}

	errMsg := "invalid credentialSubject, errors:"
	for _, desc := range result.Errors() {
		errMsg += fmt.Sprintf("- %s\n", desc)
	}

	return false, fmt.Errorf(errMsg)
}

// VerifyVCOnChain 链上验证VC的有效性
// @params vc: VC的JSON字符串
// @params client：长安链客户端
func VerifyVCOnChain(vc string, client *cmsdk.ChainClient) (bool, error) {
	params := make([]*common.KeyValuePair, 0)

	params = append(params, &common.KeyValuePair{
		Key:   model.Params_VcJson,
		Value: []byte(vc),
	})

	_, err := invoke.QueryContract(invoke.DIDContractName, model.Method_VerifyVc, params, client)
	if err != nil {
		return false, err

	}

	return true, nil
}

// RevokeVCOnChain 在链上吊销VC
// @params vcId: vc的ID编号
// @params client：长安链客户端
func RevokeVCOnChain(vcId string, client *cmsdk.ChainClient) error {
	params := make([]*common.KeyValuePair, 0)

	params = append(params, &common.KeyValuePair{
		Key:   model.Params_VcId,
		Value: []byte(vcId),
	})

	_, err := invoke.InvokeContract(invoke.DIDContractName, model.Method_RevokeVc, params, client)
	if err != nil {
		return err
	}

	return nil
}

// GetVCRevokedListFromChain 获取链上VC的吊销列表
// @params vcIdSearch：要查找的vc编号（空字符串可以查找全部列表）
// @params start：开始的索引，0表示从第一个开始
// @params count：要获取的数量，0表示获取所有
// @params client：长安链客户端
func GetVCRevokedListFromChain(vcIdSearch string, start int, count int, client *cmsdk.ChainClient) ([]string, error) {
	params := make([]*common.KeyValuePair, 0)

	params = append(params, &common.KeyValuePair{
		Key:   model.Params_VcTemplateNameSearch,
		Value: []byte(vcIdSearch),
	})

	params = append(params, &common.KeyValuePair{
		Key:   model.Params_SearchStart,
		Value: []byte(strconv.Itoa(start)),
	})

	params = append(params, &common.KeyValuePair{
		Key:   model.Params_SearchCount,
		Value: []byte(strconv.Itoa(count)),
	})

	resp, err := invoke.QueryContract(invoke.DIDContractName, model.Method_GetRevokedVcList, params, client)
	if err != nil {
		return nil, err
	}

	var revokedList []string

	err = json.Unmarshal(resp, &revokedList)
	if err != nil {
		return nil, err
	}

	return revokedList, nil
}

// AddVcIssueLogToChain 在链上添加VC签发日志
// @params issuer：签发者DID（需要在链上被认可）
// @params did：被签发者did
// @params vcId：VC业务编号
// @params vcTemplateId：模板编号
// @params client：长安链客户端
func AddVcIssueLogToChain(issuer, did, vcId, vcTemplateId string,
	client *cmsdk.ChainClient) error {
	params := make([]*common.KeyValuePair, 0)

	params = append(params, &common.KeyValuePair{
		Key:   model.Params_Issuer,
		Value: []byte(issuer),
	})

	params = append(params, &common.KeyValuePair{
		Key:   model.Params_Did,
		Value: []byte(did),
	})

	params = append(params, &common.KeyValuePair{
		Key:   model.Params_VcId,
		Value: []byte(vcId),
	})

	params = append(params, &common.KeyValuePair{
		Key:   model.Params_VcTemplateId,
		Value: []byte(vcTemplateId),
	})

	_, err := invoke.InvokeContract(invoke.DIDContractName, model.Method_VcIssueLog, params, client)
	if err != nil {
		return err
	}

	return nil
}

// GetVcIssueLogListFromChain 从链上获取VC签发日志列表
// @params vcIdSearch：VC编号关键字（空字符串可以查找全部列表）
// @params start：开始的索引，0表示从第一个开始
// @params count：要获取的数量，0表示获取所有
// @params client：长安链客户端
func GetVcIssueLogListFromChain(vcIdSearch string, start int, count int,
	client *cmsdk.ChainClient) ([]*model.VcIssueLog, error) {
	params := make([]*common.KeyValuePair, 0)

	params = append(params, &common.KeyValuePair{
		Key:   model.Params_VcIdSearch,
		Value: []byte(vcIdSearch),
	})

	params = append(params, &common.KeyValuePair{
		Key:   model.Params_SearchStart,
		Value: []byte(strconv.Itoa(start)),
	})

	params = append(params, &common.KeyValuePair{
		Key:   model.Params_SearchCount,
		Value: []byte(strconv.Itoa(count)),
	})

	resp, err := invoke.QueryContract(invoke.DIDContractName, model.Method_GetVcIssueLogs, params, client)
	if err != nil {
		return nil, err
	}

	var vcIssueLogs []*model.VcIssueLog

	err = json.Unmarshal(resp, &vcIssueLogs)
	if err != nil {
		return nil, err
	}

	return vcIssueLogs, nil
}
