package vc

import (
	"did-sdk/did"
	"did-sdk/invoke"
	"did-sdk/key"
	"did-sdk/proof"
	"did-sdk/utils"
	"encoding/json"
	"errors"
	"fmt"
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

// IssueVC 颁发VC
// @params keyInfo：颁发者的密钥信息
// @params subject：颁发信息主体，对应VC中的`credentialSubject`字段
// @params client：需要生成DID，并且在链上校验是否具有颁发资格
// @params vcId：VC的`id`字段，可以根据业务自定义
// @params expirationDate：VC的到期时间
// @params vcTemplateId：VC的模板Id，在链上获取VC模板
// @params vcType：VC中的`type`字段，描述VC的类型信息（可变参数，默认会填写“VerifiableCredential”,可继续根据业务类型追加）
func IssueVC(keyInfo *key.KeyInfo, subject map[string]interface{}, client *cmsdk.ChainClient,
	vcId string, expirationDate int64, vcTemplateId string, vcType ...string) ([]byte, error) {

	// 链上获取模板
	vcTemplate, err := GetVcTemplateFromChain(vcTemplateId, client)
	if err != nil {
		return nil, err
	}

	// 验证subject是否符合VC模板规范
	ok, err := verifyCredentialSubject(subject, vcTemplate)
	if !ok {
		return nil, err
	}

	vcType = append(vcType, "VerifiableCredential")
	issuer, err := did.GenerateDidByPK(keyInfo.PkPEM, client)
	if err != nil {
		return nil, err
	}

	// 校验Issuer是否具有颁发权限
	list, err := did.GetTrustIssuerListFromChain(issuer, 0, 0, client)
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, errors.New("not a trusted issuer")
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
	}

	vcBytes, err := json.Marshal(vc)
	if err != nil {
		return nil, err
	}

	msg, err := utils.CompactJson(vcBytes)
	if err != nil {
		return nil, err
	}

	keyId := issuer + "#keys-1"
	pf, err := proof.GenerateProofByKey(keyInfo.SkPEM, msg, keyId,
		keyInfo.Algorithm, utils.GetHashTypeByAlgorithm(keyInfo.Algorithm))
	if err != nil {
		return nil, err
	}

	vc.Proof = pf

	return json.Marshal(vc)

}

// IssueVCLocal 本地颁发VC（不经过链上计算和校验）
// @params skPem: 私钥的PEM编码
// @params algorithm: 公钥算法名称
// @params subject: 颁发信息主体，对应VC中的`credentialSubject`字段
// @params issuer: 颁发者的DID编号
// @params vcId：VC的`id`字段，可以根据业务自定义
// @params expirationDate：VC的到期时间
// @params vcTemplate：VC的模板内容，是一个JSON schema，一般存储在链上
// @params vcType：VC中的`type`字段，描述VC的类型信息（可变参数，默认会填写“VerifiableCredential”,可继续根据业务类型追加）
func IssueVCLocal(skPem []byte, algorithm string, subject map[string]interface{}, issuer string,
	vcId string, expirationDate int64, vcTemplate []byte, vcType ...string) ([]byte, error) {
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

	keyId := issuer + "#keys-1"
	pf, err := proof.GenerateProofByKey(skPem, msg, keyId,
		algorithm, utils.GetHashTypeByAlgorithm(algorithm))
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

// VerifyVCOnChain
func VerifyVCOnChain(vc string, client *cmsdk.ChainClient) (bool, error) {
	params := make([]*common.KeyValuePair, 0)

	params = append(params, &common.KeyValuePair{
		Key:   model.Params_VcJson,
		Value: []byte(vc),
	})

	_, err := invoke.InvokeContract(invoke.DIDContractName, model.Method_VerifyVc, params, client)
	if err != nil {
		return false, err
	}

	return true, nil
}
