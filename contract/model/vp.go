/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) DCPS. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package model

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/buger/jsonparser"
)

// VerifiablePresentation VP
type VerifiablePresentation struct {
	rawData              json.RawMessage
	Context              []string                `json:"@context"`
	Id                   string                  `json:"id"`
	Type                 []string                `json:"type"`
	VerifiableCredential []*VerifiableCredential `json:"verifiableCredential"`
	Holder               string                  `json:"holder"`
	ExpirationDate       string                  `json:"expirationDate,omitempty"`
	Proof                *Proof                  `json:"proof,omitempty"`
}

// NewVerifiablePresentation 根据VP持有者展示的凭证json字符串创建VP持有者展示的凭证
func NewVerifiablePresentation(vpJson string) (*VerifiablePresentation, error) {
	var vp VerifiablePresentation
	err := json.Unmarshal([]byte(vpJson), &vp)
	if err != nil {
		return nil, err
	}
	vp.rawData = []byte(vpJson)
	return &vp, nil
}

func (vp *VerifiablePresentation) Verify(pkPem []byte) (bool, error) {
	// Check if the VC type is correct
	if len(vp.Type) == 0 {
		return false, errors.New("invalid VP type")
	} else {
		var isVpType bool
		for _, v := range vp.Type {
			if v == "VerifiablePresentation" {
				isVpType = true
			}
		}

		if !isVpType {
			return false, errors.New("invalid VP type")
		}
	}

	if len(vp.ExpirationDate) != 0 {
		// 检查当前时间是否在有效期内
		myTime, err := GetTxTime()
		if err != nil {
			return false, err
		}

		expirationDate, err := time.Parse(time.RFC3339, vp.ExpirationDate)
		if err != nil {
			return false, err
		}

		if myTime > expirationDate.Unix() {
			return false, errors.New("the verifiable presentation has expired")
		}

	}

	// 删除proof字段
	withoutProof := jsonparser.Delete(vp.rawData, "proof")
	//去掉空格换行等
	withoutProof, err := CompactJson(withoutProof)
	if err != nil {
		return false, err
	}

	// 验签
	return vp.Proof.Verify(withoutProof, pkPem)
}
