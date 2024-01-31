package main

import (
	"crypto/sha256"
	"did-contract/model"
	"encoding/hex"
	"encoding/json"
	"strings"

	"chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"
)

// 此为存入数据库的世界状态key，故越短越好
const (
	keyDid          = "d"
	keyIndexPubKey  = "p"
	keyIndexAddress = "a"
	keyTrustIssuer  = "ti"
	keyRevokeVc     = "r"
	keyBlackList    = "b"
	keyVcTemplate   = "vt"
)

const (
	defaultSearchCount = 1000
	defaultSearchStart = 1
)

// Dal 数据库访问层
type Dal struct {
	didMethod string
}

func NewDal(didMethod string) *Dal {
	return &Dal{
		didMethod: didMethod,
	}
}

// Db 获取数据库实例
func (dal *Dal) Db() sdk.SDKInterface {
	return sdk.Instance
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

func (dal *Dal) putIndexPubKey(pubKey string, did string) error {
	//将索引存入数据库
	err := dal.Db().PutStateByte(keyIndexPubKey, pubKeyToDbKey(pubKey), []byte(did))
	if err != nil {
		return err
	}
	return nil
}

func (dal *Dal) deleteIndexPubKey(pubKey string) error {
	//从数据库中删除索引
	err := dal.Db().DelState(keyIndexPubKey, pubKeyToDbKey(pubKey))
	if err != nil {
		return err
	}
	return nil
}

func (dal *Dal) getDidByPubKey(pubKey string) (string, error) {
	//从数据库中获取索引
	did, err := dal.Db().GetStateByte(keyIndexPubKey, pubKeyToDbKey(pubKey))
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

func (dal *Dal) putVcTemplate(templateId string, templateName string, version string, vcTemplate string) error {
	//将VcTemplate存入数据库
	vcTemplateObj := model.VcTemplate{
		Id:       templateId,
		Name:     templateName,
		Template: vcTemplate,
		Version:  version,
	}
	value, _ := json.Marshal(vcTemplateObj)
	err := dal.Db().PutStateByte(keyVcTemplate, templateId, value)
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

func (dal *Dal) searchVcTemplate(templateNameSearch string, start int, count int) ([]string, error) {
	//从数据库中查询VcTemplate迭代器
	iter, err := dal.Db().NewIteratorPrefixWithKeyField(keyVcTemplate, templateNameSearch)
	if err != nil {
		return nil, err
	}

	defer iter.Close()

	var vcTemplateSlice []string

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

		vcTemplateSlice = append(vcTemplateSlice, string(value))
	}

	return vcTemplateSlice, nil
}

func (dal *Dal) didToDbKey(did string) string {
	didPrefix := "did:" + dal.didMethod + ":"
	return strings.TrimPrefix(did, didPrefix)
}

func pubKeyToDbKey(pubKey string) string {
	hash := sha256.Sum256([]byte(pubKey))
	return hex.EncodeToString(hash[:])
}

func vcIdToKey(vcID string) string {
	//vcid 是一个http url，为了存入数据库，需要将其转换为一个只有字母大小写、数字、下划线的字符串
	vcID = strings.ReplaceAll(vcID, ":", "_")
	vcID = strings.ReplaceAll(vcID, "/", "_")
	vcID = strings.ReplaceAll(vcID, ".", "_")
	vcID = strings.ReplaceAll(vcID, "-", "_")
	return vcID
}

func (dal *Dal) getSenderDid() (string, error) {
	senderPk, err := sdk.Instance.GetSenderPk()
	if err != nil {
		return "", err
	}

	return dal.getDidByPubKey(senderPk)
}
