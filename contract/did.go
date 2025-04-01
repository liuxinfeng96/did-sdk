/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) DCPS. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"did-contract/model"
	"errors"
	"fmt"
	"strings"
)

// DidMethod 获取DID Method
func (d *DidContract) DidMethod() (string, error) {
	return d.dal.getDidMethod()
}

// IsValidDid 判断DID URL是否合法
func (d *DidContract) IsValidDid(did string) (bool, error) {

	didMethod, err := d.dal.getDidMethod()
	if err != nil {
		return false, err
	}

	didPrefix := "did:" + didMethod + ":"

	ok := strings.HasPrefix(did, didPrefix)
	if !ok {
		return false, errors.New("invalid did method")
	}

	ok = d.dal.isInBlackList(did)
	if ok {
		return false, errors.New("the did in the black list")
	}

	return true, nil
}

// AddDidDocument 添加DID Document
func (d *DidContract) AddDidDocument(didDocument string) error {
	didDoc, err := model.NewDIDDocument(didDocument)
	if err != nil {
		return errors.New("invalid did document")
	}

	ok, err := d.IsValidDid(didDoc.Id)
	if !ok {
		return fmt.Errorf("invalid DID, err: [%s]", err.Error())
	}

	ok, err = didDoc.VerifyProof()
	if !ok {
		return fmt.Errorf("the DID doc proof verify failed, err: [%s]", err.Error())
	}

	//存储DID Document
	return d.addDidDocument(didDoc)
}

func (d *DidContract) addDidDocument(didDoc *model.DidDocument) error {
	//检查DID Document是否存在
	dbDidDoc, _ := d.dal.getDidDocument(didDoc.Id)
	if len(dbDidDoc) != 0 {
		return errors.New("did document already exists")
	}

	did, pubKeys, addresses := didDoc.ParsePubKeyAddress()

	//压缩DID Document，去掉空格和换行符
	compactDidDoc, err := didDoc.CompactDoc()
	if err != nil {
		return err
	}
	//Save did document
	err = d.dal.putDidDocument(did, compactDidDoc)
	if err != nil {
		return err
	}
	//Save pubkey index
	for _, pk := range pubKeys {
		err = d.dal.putIndexPubKey(pk, did)
		if err != nil {
			return err
		}
	}
	//save address index
	for _, addr := range addresses {
		err = d.dal.putIndexAddress(addr, did)
		if err != nil {
			return err
		}
	}

	// 发送事件
	emitSetDidDocumentEvent(did, string(didDoc.JsonRaw()))
	return nil
}

// GetDidDocument 获取DID Document
func (d *DidContract) GetDidDocument(did string) (string, error) {
	// check did valid
	valid, err := d.IsValidDid(did)
	if err != nil {
		return "", err
	}

	if !valid {
		return "", errors.New("invalid did")
	}

	didDoc, err := d.dal.getDidDocument(did)
	if err != nil {
		return "", err
	}

	return string(didDoc), nil
}

// UpdateDidDocument 更新DID Document
func (d *DidContract) UpdateDidDocument(didDocument string) error {

	didDoc, err := model.NewDIDDocument(didDocument)
	if err != nil {
		return errors.New("invalid did document")
	}

	// 检查old DID Document是否存在
	oldDocBytes, err := d.dal.getDidDocument(didDoc.Id)
	if err != nil || oldDocBytes == nil {
		return errors.New("did does not exist")
	}

	oldDoc, err := model.NewDIDDocument(string(oldDocBytes))
	if err != nil {
		return errors.New("invalid old did document")
	}

	senderDid, err := d.dal.getSenderDid()
	if err != nil {
		return err
	}

	var hasPermission bool

	if len(senderDid) != 0 {
		for _, c := range oldDoc.Controller {
			if senderDid == c {
				hasPermission = true
				break
			}
		}
	}

	if !hasPermission {
		ok, _ := isSenderAdmin(d)
		if !ok {
			return errors.New("no operation permission")
		}
	}

	ok, err := d.IsValidDid(didDoc.Id)
	if !ok {
		return fmt.Errorf("invalid DID, err: [%s]", err.Error())
	}

	ok, err = didDoc.VerifyProof()
	if !ok {
		return fmt.Errorf("the DID doc proof verify failed, err: [%s]", err.Error())
	}

	return d.updateDidDocument(didDoc, oldDoc)
}

func (d *DidContract) updateDidDocument(didDoc, oldDoc *model.DidDocument) error {

	did, pubKeys, addresses := didDoc.ParsePubKeyAddress()

	_, oldPks, oldAddresses := oldDoc.ParsePubKeyAddress()

	// 如果oldPubKeys在新的pubKeys中不存在，则删除
	for _, oldPk := range oldPks {
		if !isInList(oldPk, pubKeys) {
			err := d.dal.deleteIndexPubKey(oldPk)
			if err != nil {
				return err
			}
		}
	}
	// 如果oldAddresses在新的addresses中不存在，则删除
	for _, oldAddr := range oldAddresses {
		if !isInList(oldAddr, addresses) {
			err := d.dal.deleteIndexAddress(oldAddr)
			if err != nil {
				return err
			}
		}
	}

	// 压缩DID Document，去掉空格和换行符
	compactDidDoc, err := didDoc.CompactDoc()
	if err != nil {
		return err
	}

	// 保存新的DID Document
	err = d.dal.putDidDocument(did, compactDidDoc)
	if err != nil {
		return err
	}

	// 保存新的pubKeys
	for _, pk := range pubKeys {
		err = d.dal.putIndexPubKey(pk, did)
		if err != nil {
			return err
		}
	}
	// 保存新的addresses
	for _, addr := range addresses {
		err = d.dal.putIndexAddress(addr, did)
		if err != nil {
			return err
		}
	}

	// 发送事件
	emitSetDidDocumentEvent(did, string(didDoc.JsonRaw()))
	return nil
}

// GetDidByPubkey 根据公钥获取DID
func (d *DidContract) GetDidByPubkey(pk string) (string, error) {
	return d.dal.getDidByPubKey(pk)
}

// GetDidByAddress 根据地址获取DID
func (d *DidContract) GetDidByAddress(address string) (string, error) {
	return d.dal.getDidByAddress(address)
}

// AddBlackList 添加黑名单
func (d *DidContract) AddBlackList(dids []string) error {

	ok, err := isSenderAdmin(d)
	if err != nil {
		return err
	}

	if !ok {
		return errors.New("no operation permission")
	}

	for _, did := range dids {
		err := d.dal.putBlackList(did)
		if err != nil {
			return err
		}
	}

	emitAddBlackListEvent(dids)
	return nil
}

// DeleteBlackList 删除黑名单
func (d *DidContract) DeleteBlackList(dids []string) error {
	ok, err := isSenderAdmin(d)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("no operation permission")
	}

	for _, did := range dids {
		err := d.dal.deleteBlackList(did)
		if err != nil {
			return err
		}
	}

	emitDeleteBlackListEvent(dids)
	return nil
}

// GetBlackList 获取黑名单
func (d *DidContract) GetBlackList(didSearch string, start int, count int) ([]string, error) {
	return d.dal.searchBlackList(didSearch, start, count)
}

// AddTrustIssuerList 添加信任发行者
func (d *DidContract) AddTrustIssuerList(dids []string) error {
	ok, err := isSenderAdmin(d)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("no operation permission")
	}

	for _, did := range dids {
		// 判断DID Doc是否已经上链
		ok := d.dal.isDidDocExisting(did)
		if !ok {
			return fmt.Errorf("the did's doc not found on chain, did: [%s]", did)
		}

		err := d.dal.putTrustIssuer(did)
		if err != nil {
			return err
		}
	}
	emitAddTrustIssuerListEvent(dids)
	return nil
}

// DeleteTrustIssuer 删除信任发行者
func (d *DidContract) DeleteTrustIssuer(dids []string) error {
	ok, err := isSenderAdmin(d)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("no operation permission")
	}

	for _, did := range dids {
		err := d.dal.deleteTrustIssuer(did)
		if err != nil {
			return err
		}
	}

	emitDeleteTrustIssuerEvent(dids)
	return nil
}

// GetTrustIssuer 获取信任发行者
func (e *DidContract) GetTrustIssuer(didSearch string, start int, count int) ([]string, error) {
	return e.dal.searchTrustIssuer(didSearch, start, count)
}
