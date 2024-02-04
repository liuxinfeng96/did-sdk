package main

import (
	"did-contract/model"
	"errors"
	"fmt"
	"strings"
)

// VerifyVc 验证VC的有效性
func (d *DidContract) VerifyVc(vcJson string) (bool, error) {

	vc, err := model.NewVerifiableCredential(vcJson)
	if err != nil {
		return false, fmt.Errorf("invalid vc: [%s]", err.Error())
	}

	//检查vc拥有者是否在黑名单中
	if d.dal.isInBlackList(vc.GetCredentialSubjectID()) {
		return false, errors.New("vc owner is in black list")
	}

	// 检查签发者是否可信任
	if d.isTrustIssuer(vc.Issuer) {
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
func (d *DidContract) RevokeVc(vcID string) error {
	ok, err := isSenderCreator()
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("no operation permission")
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
		ok, _ := isSenderCreator()
		if !ok {
			return errors.New("no operation permission")
		}
	}

	err := d.dal.putVcTemplate(id, name, version, template)
	if err != nil {
		return err
	}

	emitSetVcTemplateEvent(id, name, version, template)
	return nil
}

// GetVcTemplate 获取VC模板
func (d *DidContract) GetVcTemplate(id string) ([]byte, error) {
	return d.dal.getVcTemplate(id)
}

// GetVcTemplateList 获取VC模板列表
func (d *DidContract) GetVcTemplateList(templateNameSearch string, start int, count int) (
	[]string, error) {
	return d.dal.searchVcTemplate(templateNameSearch, start, count)
}
