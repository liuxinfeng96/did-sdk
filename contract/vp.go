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

// VerifyVp 验证vp
// @params vpJson vp的json字符串
func (d *DidContract) VerifyVp(vpJson string) (bool, error) {

	vp, err := model.NewVerifiablePresentation(vpJson)
	if err != nil {
		return false, fmt.Errorf("invalid vp: [%s]", err.Error())
	}

	// 检查持有者是否在黑名单中
	if d.dal.isInBlackList(vp.Holder) {
		return false, errors.New("vp owner is in black list")
	}

	// 验证VP中的VC
	for _, v := range vp.VerifiableCredential {

		subId, err := v.GetCredentialSubjectID()
		if err != nil {
			return false, err
		}

		if vp.Holder != subId {
			return false, errors.New("the holder is different from the VC's subject ID")
		}

		vcString, err := json.Marshal(v)
		if err != nil {
			return false, err
		}

		ok, err := d.VerifyVc(string(vcString))
		if !ok {
			return false, fmt.Errorf("vc verify failed, err: [%s]", err.Error())
		}

	}
	// 获取签发者DID公钥
	didDoc, err := d.dal.getDidDocument(vp.Holder)
	if err != nil {
		return false, fmt.Errorf("get did document of the holder failed, err: [%s]", err.Error())
	}

	doc, err := model.NewDIDDocument(string(didDoc))
	if err != nil {
		return false, fmt.Errorf("get did document of the holder failed, err: [%s]", err.Error())
	}

	signerDid := vp.Proof.VerificationMethod[0:strings.Index(vp.Proof.VerificationMethod, "#")]

	// 判断证明是不是属于持有者
	if signerDid != vp.Holder {
		return false, errors.New("the proof that does not belong to the holder")
	}

	pkPem, err := doc.GetPkPemByVerificationMethodId(vp.Proof.VerificationMethod)
	if err != nil {
		return false, fmt.Errorf("get pk from did doc failed, err: [%s]", err.Error())
	}

	return vp.Verify([]byte(pkPem))
}
