package main

import (
	"did-contract/model"
	"errors"
	"fmt"
	"strings"
)

func (d *DidContract) VerifyVp(vpJson string) (bool, error) {

	vp, err := model.NewVerifiablePresentation(vpJson)
	if err != nil {
		return false, fmt.Errorf("invalid vp: [%s]", err.Error())
	}

	// 检查持有者是否在黑名单中
	if d.dal.isInBlackList(vp.Holder) {
		return false, errors.New("vp owner is in black list")
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

	if signerDid != vp.Holder {
		return false, errors.New("the proof that does not belong to the holder")
	}

	pkPem, err := doc.GetPkPemByVerificationMethodId(vp.Proof.VerificationMethod)
	if err != nil {
		return false, fmt.Errorf("get pk from did doc failed, err: [%s]", err.Error())
	}

	return vp.Verify([]byte(pkPem))
}
