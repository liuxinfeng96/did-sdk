/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) DCPS. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"did-contract/model"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// VerifyVc 验证VC的有效性
// @params vcJson vc的json字符串
func (d *DidContract) VerifyVc(vcJson string) (bool, error) {

	vc, err := model.NewVerifiableCredential(vcJson)
	if err != nil {
		return false, fmt.Errorf("invalid vc: [%s]", err.Error())
	}

	subId, err := vc.GetCredentialSubjectID()
	if err != nil {
		return false, err
	}

	//检查vc拥有者是否在黑名单中
	if d.dal.isInBlackList(subId) {
		return false, errors.New("vc owner is in black list")
	}

	// 检查签发者是否可信任
	if !d.isTrustIssuer(vc.Issuer) {
		return false, errors.New("the issuer of VC is not a trusted issuer on the chain")
	}

	// 检查VC撤销状态
	if d.isInRevokeVcList(vc.Id) {
		return false, errors.New("the VC is revoked")
	}

	//检查VC模板
	var vcTemplateBytes []byte
	if vc.Template != nil {
		vcTemplateBytes, err = d.dal.getVcTemplate(vc.Template.ID)
		if err != nil {
			return false, err
		}

		if vcTemplateBytes == nil {
			return false, errors.New("VC template was not found")
		}
	}

	// 获取签发者DID公钥
	didDoc, err := d.dal.getDidDocument(vc.Issuer)
	if err != nil {
		return false, fmt.Errorf("get did document of the issuer failed, err: [%s]", err.Error())
	}

	doc, err := model.NewDIDDocument(string(didDoc))
	if err != nil {
		return false, fmt.Errorf("get did document of the issuer failed, err: [%s]", err.Error())
	}

	signerDid := vc.Proof.VerificationMethod[0:strings.Index(vc.Proof.VerificationMethod, "#")]

	if signerDid != vc.Issuer {
		return false, errors.New("the proof that does not belong to the issuer")
	}

	pkPem, err := doc.GetPkPemByVerificationMethodId(vc.Proof.VerificationMethod)
	if err != nil {
		return false, fmt.Errorf("get pk from did doc failed, err: [%s]", err.Error())
	}

	return vc.Verify([]byte(pkPem), vcTemplateBytes)
}

// RevokeVc 撤销VC
// @params vcID VC业务编号
func (d *DidContract) RevokeVc(vcID string) error {
	// 判断是不是管理员
	ok, err := isSenderAdmin(d)
	if err != nil {
		return err
	}
	if !ok {
		// 判断是不是签发者本人
		isIssuer, _ := d.isSenderIssued(vcID)
		if !isIssuer {
			return errors.New("no operation permission")
		}
	}

	err = d.dal.putRevokeVc(vcID)
	if err != nil {
		return err
	}
	emitRevokeVcEvent(vcID)
	return nil
}

// GetRevokedVcList 获取撤销VC列表
func (d *DidContract) GetRevokedVcList(vcIDSearch string, start int, count int) ([]string, error) {
	return d.dal.searchRevokeVc(vcIDSearch, start, count)
}

func (d *DidContract) isInRevokeVcList(id string) bool {
	dbId, err := d.dal.getRevokeVc(id)
	if err != nil || len(dbId) == 0 {
		return false
	}
	return true
}

// SetVcTemplate 设置VC模板
func (d *DidContract) SetVcTemplate(id string, name string, version string, template string) error {
	// 判读是否有权限
	if !d.isSenderTrustIssuer() {
		ok, _ := isSenderAdmin(d)
		if !ok {
			return errors.New("no operation permission")
		}
	}

	value, _ := d.GetVcTemplate(id)
	if len(value) != 0 {
		return errors.New("the VC template already exists")
	}

	// 需要校验一下模板里面是否包含ID字段
	var tempJson model.VcTemplateJSONSchema
	err := json.Unmarshal([]byte(template), &tempJson)
	if err != nil {
		return errors.New("the template does not conform to the json Schema specification")
	}

	var isIncludedId bool

	for _, v := range tempJson.Required {
		if v == "id" {
			isIncludedId = true
			break
		}
	}

	if !isIncludedId {
		return errors.New("the template must contain the `id` subfield")
	}

	vcTemp := &model.VcTemplate{
		Id:       id,
		Name:     name,
		Version:  version,
		Template: json.RawMessage(template),
	}

	tempBytes, err := json.Marshal(vcTemp)
	if err != nil {
		return err
	}

	err = d.dal.putVcTemplate(id, tempBytes)
	if err != nil {
		return err
	}

	emitSetVcTemplateEvent(id, tempBytes)
	return nil
}

// GetVcTemplate 获取VC模板
func (d *DidContract) GetVcTemplate(id string) ([]byte, error) {
	return d.dal.getVcTemplate(id)
}

// GetVcTemplateList 获取VC模板列表
func (d *DidContract) GetVcTemplateList(templateNameSearch string, start int, count int) (
	[]*model.VcTemplate, error) {
	return d.dal.searchVcTemplate(templateNameSearch, start, count)
}

// VcIssueLog 存储签发日志
func (d *DidContract) VcIssueLog(issuer, did, templateId, vcId string) error {
	// 校验签发者是否具有签发资格
	if !d.isTrustIssuer(issuer) {
		return errors.New("the issuer is not in trust issuer list")
	}

	// 校验被签发者是否合格
	if d.dal.isInBlackList(did) {
		return errors.New("the did is in black list")
	}

	if !d.dal.isDidDocExisting(did) {
		return fmt.Errorf("the did's doc not found on chain, did: [%s]", did)
	}

	// 检查模板是否在链上
	temp, err := d.GetVcTemplate(templateId)
	if err != nil {
		return err
	}

	if len(temp) == 0 {
		return errors.New("the vc template not found on chain")
	}

	myTime, err := model.GetTxTime()
	if err != nil {
		return err
	}

	vcIssueLog := model.NewVcIssueLog(issuer, did, templateId, vcId, myTime)

	v, err := json.Marshal(vcIssueLog)
	if err != nil {
		return err
	}

	emitVcIssueLogEvent(vcId, v)

	return d.dal.putVcIssueLog(vcId, v)
}

// GetVcIssueLogs 获取签发日志列表
func (d *DidContract) GetVcIssueLogs(vcIdSearch string, start int, count int) (
	[]*model.VcIssueLog, error) {
	return d.dal.searchVcIssueLogs(vcIdSearch, start, count)
}
