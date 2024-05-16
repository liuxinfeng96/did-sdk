/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) DCPS. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package model

import (
	"encoding/json"
	"errors"

	"github.com/buger/jsonparser"
)

// DidDocument the JSON structure of the DID document
type DidDocument struct {
	rawData            json.RawMessage
	Context            string                `json:"@context"`
	Id                 string                `json:"id"`
	Created            string                `json:"created"`
	Updated            string                `json:"updated"`
	VerificationMethod []*VerificationMethod `json:"verificationMethod"`
	Service            []struct {
		ID              string `json:"id"`
		Type            string `json:"type"`
		ServiceEndpoint string `json:"serviceEndpoint"`
	} `json:"service,omitempty"`

	Authentication []string        `json:"authentication"`
	Controller     []string        `json:"controller"`
	Proof          json.RawMessage `json:"proof,omitempty"`
}

// VerificationMethod the JSON structure of the DID document VerificationMethod
type VerificationMethod struct {
	Id           string `json:"id"`
	Type         string `json:"type"`
	Controller   string `json:"controller"`
	PublicKeyPem string `json:"publicKeyPem"`
	Address      string `json:"address"`
}

// NewDIDDocument 根据DID文档json字符串创建DID文档
func NewDIDDocument(didDocumentJson string) (*DidDocument, error) {
	var didDocument DidDocument
	err := json.Unmarshal([]byte(didDocumentJson), &didDocument)
	if err != nil {
		return nil, err
	}
	didDocument.rawData = []byte(didDocumentJson)
	return &didDocument, nil
}

// GetPkPemByVerificationMethodId 通过VerificationMethod的ID获取公钥PEM编码
func (d *DidDocument) GetPkPemByVerificationMethodId(id string) (string, error) {
	for _, vm := range d.VerificationMethod {
		if vm.Id == id {
			return vm.PublicKeyPem, nil
		}
	}

	return "", errors.New("the verification method was not found")
}

// VerifyProof 证明的验证
func (d *DidDocument) VerifyProof() (bool, error) {
	//删除proof字段
	withoutProof := jsonparser.Delete(d.rawData, "proof")
	//去掉空格换行等
	msg, err := CompactJson(withoutProof)
	if err != nil {
		return false, err
	}

	// 反序列化判断是一把密钥还是多把密钥
	// 如果是一把密钥，默认验证索引为0的公钥
	var pf Proof
	if err := json.Unmarshal(d.Proof, &pf); err == nil {
		return pf.Verify(msg, []byte(d.VerificationMethod[0].PublicKeyPem))
	}

	pfs := make([]*Proof, 0)

	if err := json.Unmarshal(d.Proof, &pfs); err == nil {
		for index, p := range pfs {
			ok, err := p.Verify(msg, []byte(d.VerificationMethod[index].PublicKeyPem))
			if !ok {
				return false, err
			}
		}

	}

	return true, nil
}

// ParsePubKeyAddress 从DOC里获取公钥列表和地址列表
func (d *DidDocument) ParsePubKeyAddress() (didUrl string, pubKeys []string, addresses []string) {
	pubKeys = make([]string, 0)
	addresses = make([]string, 0)
	for _, pk := range d.VerificationMethod {
		pubKeys = append(pubKeys, pk.PublicKeyPem)
		addresses = append(addresses, pk.Address)
	}
	return d.Id, pubKeys, addresses
}

// CompactDoc 返回压缩的JSON DOC
func (d *DidDocument) CompactDoc() ([]byte, error) {
	return CompactJson(d.rawData)
}

// JsonRaw 返回JSON原文
func (d *DidDocument) JsonRaw() []byte {
	return d.rawData
}
