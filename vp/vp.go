/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) DCPS. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package vp

import (
	"did-sdk/did"
	"did-sdk/invoke"
	"did-sdk/proof"
	"did-sdk/utils"
	"encoding/json"
	"strconv"

	"chainmaker.org/chainmaker/did-contract/model"
	"chainmaker.org/chainmaker/pb-go/v2/common"
	cmsdk "chainmaker.org/chainmaker/sdk-go/v2"
)

var ContextVP = []string{
	"https://www.w3.org/2018/credentials/v1",
	"https://www.w3.org/2018/credentials/examples/v1",
}

// GenerateVP 生成自己的VP
// @params skPem: 私钥的PEM编码
// @params keyIndex：公钥在DID文档中的索引
// @params vpId：VP的`id`字段，可以根据业务自定义
// @params VP中包含的VC列表
// @params vpType：VP中的`type`字段，描述VP的类型信息（可变参数，默认会填写“VerifiablePresentation”,可继续根据业务类型追加）
func GenerateVP(skPem []byte, keyIndex int, holder string,
	vpId string, vcList []string, vpType ...string) ([]byte, error) {

	var verifiablePresentation model.VerifiablePresentation

	vpType = append(vpType, "VerifiablePresentation")
	for _, v := range vcList {
		var verifiableCredential model.VerifiableCredential
		err := json.Unmarshal([]byte(v), &verifiableCredential)
		if err != nil {
			return nil, err
		}

		verifiablePresentation.VerifiableCredential =
			append(verifiablePresentation.VerifiableCredential, &verifiableCredential)
	}

	verifiablePresentation.Context = ContextVP
	verifiablePresentation.Id = vpId
	verifiablePresentation.Type = vpType
	verifiablePresentation.Holder = holder

	vpBytes, err := json.Marshal(verifiablePresentation)
	if err != nil {
		return nil, err
	}

	msg, err := utils.CompactJson(vpBytes)
	if err != nil {
		return nil, err
	}

	keyId := holder + did.VerificationMethodKeySuffix + strconv.Itoa(keyIndex)

	pf, err := proof.GenerateProofByKey(skPem, msg, keyId)
	if err != nil {
		return nil, err
	}

	verifiablePresentation.Proof = pf

	return json.Marshal(verifiablePresentation)

}

// VerifyVPOnChain 在链上验证VP的有效性
// @params vc: VP的JSON字符串
// @params client：长安链客户端
func VerifyVPOnChain(vp string, client *cmsdk.ChainClient) (bool, error) {
	params := make([]*common.KeyValuePair, 0)

	params = append(params, &common.KeyValuePair{
		Key:   model.Params_VpJson,
		Value: []byte(vp),
	})

	_, err := invoke.QueryContract(invoke.DIDContractName, model.Method_VerifyVp, params, client)
	if err != nil {
		return false, err
	}

	return true, nil
}
