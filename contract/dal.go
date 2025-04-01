/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) DCPS. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"crypto/sha256"
	"did-contract/model"
	"encoding/hex"
	"encoding/json"
	"strings"

	"chainmaker.org/chainmaker/common/v2/evmutils"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"
)

// 此为存入数据库的世界状态key，故越短越好
const (
	keyDid           = "d"
	keyIndexPubKey   = "p"
	keyIndexAddress  = "a"
	keyTrustIssuer   = "ti"
	keyRevokeVc      = "r"
	keyBlackList     = "b"
	keyVcTemplate    = "vt"
	keyContractAdmin = "admin"
	keyVcIssueLog    = "l"

	// 合约状态数据，只存出一次，不需要很短的key来节省空间
	keyContractStatus       = "cs"
	failedDidMethod         = "didMethod"
	failedEnableTrustIssuer = "enableTrustIssuer"
)

const (
	defaultSearchCount = 1000
	defaultSearchStart = 1
)

// Dal 数据库访问层
type Dal struct {
}

func NewDal() *Dal {
	return &Dal{}
}

// Db 获取数据库实例
func (dal *Dal) Db() sdk.SDKInterface {
	return sdk.Instance
}

func (dal *Dal) putDidMethod(didMethod string) error {
	// 将DID Method存入数据库
	err := dal.Db().PutStateByte(keyContractStatus, failedDidMethod, []byte(didMethod))
	if err != nil {
		return err
	}
	return nil
}

func (dal *Dal) getDidMethod() (string, error) {
	// 从数据库中获取DID Method
	didMethod, err := dal.Db().GetStateByte(keyContractStatus, failedDidMethod)
	if err != nil {
		return "", err
	}
	return string(didMethod), nil
}

func (dal *Dal) putEnableTrustIssuer(enableTrustIssuer string) error {
	// 将EnableTrustIssuer 存入数据库
	err := dal.Db().PutStateByte(keyContractStatus, failedEnableTrustIssuer, []byte(enableTrustIssuer))
	if err != nil {
		return err
	}
	return nil
}

func (dal *Dal) getEnableTrustIssuer() (string, error) {
	// 从数据库中获取 EnableTrustIssuer
	enableTrustIssuer, err := dal.Db().GetStateByte(keyContractStatus, failedEnableTrustIssuer)
	if err != nil {
		return "", err
	}
	return string(enableTrustIssuer), nil
}

func (dal *Dal) putAdmin(ski string) error {
	err := dal.Db().PutStateByte(keyContractAdmin, ski, []byte(ski))
	if err != nil {
		return err
	}
	return nil
}

func (dal *Dal) deleteAdmin(ski string) error {
	err := dal.Db().DelState(keyContractAdmin, ski)
	if err != nil {
		return err
	}
	return nil
}

func (dal *Dal) isAdmin(ski string) bool {
	//从数据库中获取DID Document
	v, err := dal.Db().GetStateByte(keyContractAdmin, ski)
	if err != nil {
		return false
	}

	if len(v) == 0 {
		return false
	}

	return true
}

func (dal *Dal) putDidDocument(did string, didDocument []byte) error {
	//将DID Document存入数据库
	err := dal.Db().PutStateByte(keyDid, dal.didToDbKey(did), didDocument)
	if err != nil {
		return err
	}
	return nil
}

func (dal *Dal) getDidDocument(did string) ([]byte, error) {
	//从数据库中获取DID Document
	didDocument, err := dal.Db().GetStateByte(keyDid, dal.didToDbKey(did))
	if err != nil {
		return nil, err
	}
	return didDocument, nil
}

func (dal *Dal) isDidDocExisting(did string) bool {
	didDocument, err := dal.Db().GetStateByte(keyDid, dal.didToDbKey(did))
	if err != nil {
		return false
	}

	if len(didDocument) == 0 {
		return false
	}

	return true
}

func (dal *Dal) putIndexPubKey(pubKey string, did string) error {
	//将索引存入数据库
	err := dal.Db().PutStateByte(keyIndexPubKey, pubKeyToDbKey([]byte(pubKey)), []byte(did))
	if err != nil {
		return err
	}
	return nil
}

func (dal *Dal) deleteIndexPubKey(pubKey string) error {
	//从数据库中删除索引
	err := dal.Db().DelState(keyIndexPubKey, pubKeyToDbKey([]byte(pubKey)))
	if err != nil {
		return err
	}
	return nil
}

func (dal *Dal) getDidByPubKey(pubKey string) (string, error) {
	//从数据库中获取索引
	did, err := dal.Db().GetStateByte(keyIndexPubKey, pubKeyToDbKey([]byte(pubKey)))
	if err != nil {
		return "", err
	}
	return string(did), nil
}

func (dal *Dal) putIndexAddress(address string, did string) error {
	//将索引存入数据库
	err := dal.Db().PutStateByte(keyIndexAddress, address, []byte(did))
	if err != nil {
		return err
	}
	return nil
}

func (dal *Dal) deleteIndexAddress(address string) error {
	//从数据库中删除索引
	err := dal.Db().DelState(keyIndexAddress, address)
	if err != nil {
		return err
	}
	return nil
}

func (dal *Dal) getDidByAddress(address string) (string, error) {
	//从数据库中获取索引
	did, err := dal.Db().GetStateByte(keyIndexAddress, address)
	if err != nil {
		return "", err
	}
	return string(did), nil
}

func (dal *Dal) putBlackList(did string) error {
	//将BlackList存入数据库
	err := dal.Db().PutStateByte(keyBlackList, dal.didToDbKey(did), []byte(did))
	if err != nil {
		return err
	}
	return nil
}

func (dal *Dal) isInBlackList(did string) bool {
	//从数据库中获取BlackList
	dbId, err := dal.Db().GetStateByte(keyBlackList, dal.didToDbKey(did))
	if err != nil || len(dbId) == 0 {
		return false
	}
	return true
}

func (dal *Dal) deleteBlackList(did string) error {
	//从数据库中删除BlackList
	err := dal.Db().DelState(keyBlackList, dal.didToDbKey(did))
	if err != nil {
		return err
	}
	return nil
}

func (dal *Dal) searchBlackList(didSearch string, start int, count int) ([]string, error) {
	//从数据库中查询BlackList迭代器
	iter, err := dal.Db().NewIteratorPrefixWithKeyField(keyBlackList, dal.didToDbKey(didSearch))
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	var didSlice []string

	if count == 0 {
		count = defaultSearchCount
	}

	if start == 0 {
		start = defaultSearchStart
	}

	for i := 1; iter.HasNext(); i++ {
		_, _, value, err := iter.Next()
		if err != nil {
			return nil, err
		}

		if i >= start+count {
			break
		}

		if i < start {
			continue
		}

		didSlice = append(didSlice, string(value))
	}

	return didSlice, nil
}

func (dal *Dal) putTrustIssuer(did string) error {
	//将TrustIssuer存入数据库
	err := dal.Db().PutStateByte(keyTrustIssuer, dal.didToDbKey(did), []byte(did))
	if err != nil {
		return err
	}
	return nil
}

func (dal *Dal) deleteTrustIssuer(did string) error {
	//从数据库中删除TrustIssuer
	err := dal.Db().DelState(keyTrustIssuer, dal.didToDbKey(did))
	if err != nil {
		return err
	}
	return nil
}

func (dal *Dal) searchTrustIssuer(didSearch string, start int, count int) ([]string, error) {
	//从数据库中查询RevokeVc迭代器
	iter, err := dal.Db().NewIteratorPrefixWithKeyField(keyTrustIssuer, dal.didToDbKey(didSearch))
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	var didSlice []string

	if count == 0 {
		count = defaultSearchCount
	}

	if start == 0 {
		start = defaultSearchStart
	}

	for i := 1; iter.HasNext(); i++ {
		_, _, value, err := iter.Next()
		if err != nil {
			return nil, err
		}

		if i >= start+count {
			break
		}

		if i < start {
			continue
		}

		didSlice = append(didSlice, string(value))
	}

	return didSlice, nil
}

func (dal *Dal) putRevokeVc(vcID string) error {
	//将RevokeVc存入数据库
	err := dal.Db().PutStateByte(keyRevokeVc, vcIdToKey(vcID), []byte(vcID))
	if err != nil {
		return err
	}
	return nil
}

func (dal *Dal) getRevokeVc(vcID string) (string, error) {
	//从数据库中获取RevokeVc
	vcIDUrl, err := dal.Db().GetStateByte(keyRevokeVc, vcIdToKey(vcID))
	if err != nil {
		return "", err
	}
	return string(vcIDUrl), nil
}

// searchRevokeVc 根据vcID前缀查询RevokeVc,start为起始位置从0开始，count为查询数量
func (dal *Dal) searchRevokeVc(vcIDSearch string, start int, count int) ([]string, error) {
	//从数据库中查询RevokeVc迭代器
	iter, err := dal.Db().NewIteratorPrefixWithKeyField(keyRevokeVc, vcIdToKey(vcIDSearch))
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	var vcIdSlice []string

	if count == 0 {
		count = defaultSearchCount
	}

	if start == 0 {
		start = defaultSearchStart
	}

	for i := 1; iter.HasNext(); i++ {
		_, _, value, err := iter.Next()
		if err != nil {
			return nil, err
		}

		if i >= start+count {
			break
		}

		if i < start {
			continue
		}

		vcIdSlice = append(vcIdSlice, string(value))
	}

	return vcIdSlice, nil
}

func (dal *Dal) putVcTemplate(templateId string, template []byte) error {
	err := dal.Db().PutStateByte(keyVcTemplate, templateId, template)
	if err != nil {
		return err
	}
	return nil
}

func (dal *Dal) getVcTemplate(templateId string) ([]byte, error) {
	//从数据库中获取VcTemplate
	value, err := dal.Db().GetStateByte(keyVcTemplate, templateId)
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (dal *Dal) searchVcTemplate(templateNameSearch string, start int, count int) ([]*model.VcTemplate, error) {
	//从数据库中查询VcTemplate迭代器
	iter, err := dal.Db().NewIteratorPrefixWithKeyField(keyVcTemplate, templateNameSearch)
	if err != nil {
		return nil, err
	}

	defer iter.Close()

	var vcTemplateSlice []*model.VcTemplate

	if count == 0 {
		count = defaultSearchCount
	}

	if start == 0 {
		start = defaultSearchStart
	}

	for i := 1; iter.HasNext(); i++ {
		_, _, value, err := iter.Next()
		if err != nil {
			return nil, err
		}

		if i >= start+count {
			break
		}

		if i < start {
			continue
		}

		var temp model.VcTemplate
		err = json.Unmarshal(value, &temp)
		if err != nil {
			return nil, err
		}

		vcTemplateSlice = append(vcTemplateSlice, &temp)
	}

	return vcTemplateSlice, nil
}

func (dal *Dal) putVcIssueLog(vcId string, log []byte) error {
	err := dal.Db().PutStateByte(keyVcIssueLog, vcIdToKey(vcId), log)
	if err != nil {
		return err
	}

	return nil
}

func (dal *Dal) getVcIssueLog(vcId string) ([]byte, error) {
	return dal.Db().GetStateByte(keyVcIssueLog, vcIdToKey(vcId))
}

func (dal *Dal) searchVcIssueLogs(searchVcId string, start, count int) ([]*model.VcIssueLog, error) {
	iter, err := dal.Db().NewIteratorPrefixWithKeyField(keyVcIssueLog, vcIdToKey(searchVcId))
	if err != nil {
		return nil, err
	}

	defer iter.Close()

	var issueLogSlice []*model.VcIssueLog

	if count == 0 {
		count = defaultSearchCount
	}

	if start == 0 {
		start = defaultSearchStart
	}

	for i := 1; iter.HasNext(); i++ {
		_, _, value, err := iter.Next()
		if err != nil {
			return nil, err
		}

		if i >= start+count {
			break
		}

		if i < start {
			continue
		}

		var issueLog model.VcIssueLog
		err = json.Unmarshal(value, &issueLog)
		if err != nil {
			return nil, err
		}

		issueLogSlice = append(issueLogSlice, &issueLog)
	}

	return issueLogSlice, nil
}

func (dal *Dal) didToDbKey(did string) string {
	didMethod, _ := dal.getDidMethod()
	didPrefix := "did:" + didMethod + ":"
	return strings.TrimPrefix(did, didPrefix)
}

func pubKeyToDbKey(pubKey []byte) string {
	hash := sha256.Sum256(pubKey)
	return hex.EncodeToString(hash[:])
}

func vcIdToKey(vcID string) string {
	//vcid 是一个http url，为了存入数据库，需要将其转换为一个只有字母大小写、数字、下划线的字符串
	if len(vcID) == 0 {
		return vcID
	}

	vcID = strings.ReplaceAll(vcID, ":", "_")
	vcID = strings.ReplaceAll(vcID, "/", "_")
	vcID = strings.ReplaceAll(vcID, ".", "_")
	vcID = strings.ReplaceAll(vcID, "-", "_")
	return vcID
}

func (dal *Dal) getSenderDid() (string, error) {
	ski, err := sdk.Instance.GetSenderPk()
	if err != nil {
		return "", err
	}

	skiBytes, err := hex.DecodeString(ski)
	if err != nil {
		return "", err
	}

	bytesAddr := evmutils.Keccak256(skiBytes)
	addr := hex.EncodeToString(bytesAddr)[24:]

	return dal.getDidByAddress(addr)
}
