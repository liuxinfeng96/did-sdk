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
	"strconv"
	"strings"

	"chainmaker.org/chainmaker/contract-sdk-go/v2/pb/protogo"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"
)

// RequireString 必须要有参数 string 类型
// @param key
// @return string
// @return error
func RequireString(key string) (string, error) {
	args := sdk.Instance.GetArgs()
	b, ok := args[key]
	if !ok || len(b) == 0 {
		return "", fmt.Errorf("missing required parameters:'%s'", key)
	}
	return string(b), nil
}

// RequireBytes 必须要有参数 []bytes 类型
// @param key
// @return []byte
// @return error
func RequireBytes(key string) ([]byte, error) {
	args := sdk.Instance.GetArgs()
	b, ok := args[key]
	if !ok || len(b) == 0 {
		return nil, fmt.Errorf("missing required parameters:'%s'", key)
	}
	return b, nil
}

// RequireBool 必须要有参数 Bool类型
// @param key
// @return string
// @return error
func RequireBool(key string) (bool, error) {
	args := sdk.Instance.GetArgs()
	b, ok := args[key]
	if !ok || len(b) == 0 {
		return false, fmt.Errorf("missing required parameters:'%s'", key)
	}

	str := string(b)

	switch strings.ToLower(str) {
	case "false":
		return false, nil
	case "true":
		return true, nil
	default:
		return false, errors.New("parameter error of Boolean type")
	}
}

// RequireStringList 必须要有参数key1 单个string或者key2 []string类型
// @param key
// @return []string
// @return error
func RequireStringList(key1, key2 string) ([]string, error) {
	args := sdk.Instance.GetArgs()
	b, ok := args[key2]
	if !ok || len(b) == 0 {
		b1, ok1 := args[key1]
		if !ok1 || len(b1) == 0 {
			return nil, fmt.Errorf("missing required parameters:'%s' or '%s'", key1, key2)
		}
		return []string{string(b1)}, nil
	}

	var didSlice []string
	err := json.Unmarshal(b, &didSlice)
	if err != nil {
		return nil, err
	}

	return didSlice, nil
}

// Return 封装返回Bool类型为Response，如果有error则忽略bool，封装error
// @param err
// @return Response
func Return(err error) protogo.Response {
	if err != nil {
		return sdk.Error(err.Error())
	}
	return sdk.SuccessResponse
}

// ReturnString 封装返回string类型为Response，如果有error则忽略str，封装error
// @param str
// @param err
// @return Response
func ReturnString(str string, err error) protogo.Response {
	if err != nil {
		return sdk.Error(err.Error())
	}
	return sdk.Success([]byte(str))
}

// ReturnBytes 封装返回[]byte类型为Response，如果有error则忽略str，封装error
func ReturnBytes(str []byte, err error) protogo.Response {
	if err != nil {
		return sdk.Error(err.Error())
	}
	return sdk.Success(str)
}

// ReturnBool 封装返回bool类型为Response，如果有error则忽略bool，封装error
func ReturnBool(b bool, e error) protogo.Response {
	if e != nil {
		return sdk.Error(e.Error())
	}
	if b {
		return sdk.Success([]byte("true"))
	}
	return sdk.Success([]byte("false"))
}

// ReturnJson 封装返回interface类型为json string Response
// @param data
// @return Response
func ReturnJson(data interface{}, err error) protogo.Response {
	if err != nil {
		return sdk.Error(err.Error())
	}
	standardsBytes, err := json.Marshal(data)
	if err != nil {
		return sdk.Error(err.Error())
	}
	return sdk.Success(standardsBytes)
}

// OptionInt 获取可选参数 int类型，没有则返回defaultValue
func OptionInt(key string, defaultValue int) int {
	args := sdk.Instance.GetArgs()
	b, ok := args[key]
	if !ok {
		return defaultValue
	}
	num, err := strconv.Atoi(string(b))
	if err != nil {
		return defaultValue
	}
	return num
}

func isSenderCreator() (bool, error) {
	createrPk, err := sdk.Instance.GetCreatorPk()
	if err != nil {
		return false, err
	}

	senderPk, err := sdk.Instance.GetSenderPk()
	if err != nil {
		return false, err
	}

	if createrPk == senderPk {
		return true, nil
	}

	return false, nil
}

func isSenderAdmin(d *DidContract) (bool, error) {
	senderPk, err := sdk.Instance.GetSenderPk()
	if err != nil {
		return false, err
	}

	createrPk, err := sdk.Instance.GetCreatorPk()
	if err != nil {
		return false, err
	}

	return (d.IsAdmin(senderPk)) || (senderPk == createrPk), nil
}

func (d *DidContract) isSenderTrustIssuer() bool {
	enableTrustIssuer, _ := d.dal.getEnableTrustIssuer()
	if enableTrustIssuer != "true" {
		return true
	}

	did, err := d.dal.getSenderDid()
	if err != nil {
		return false
	}

	dbId, err := d.dal.Db().GetStateByte(keyTrustIssuer, d.dal.didToDbKey(did))
	if err != nil {
		return false
	}

	if len(dbId) == 0 {
		return false
	}

	return true
}

func (d *DidContract) isTrustIssuer(did string) bool {

	enableTrustIssuer, _ := d.dal.getEnableTrustIssuer()
	if enableTrustIssuer != "true" {
		return true
	}

	dbId, err := d.dal.Db().GetStateByte(keyTrustIssuer, d.dal.didToDbKey(did))
	if err != nil {
		return false
	}

	if len(dbId) == 0 {
		return false
	}

	return true
}

func isInList(str string, list []string) bool {
	for _, k := range list {
		if k == str {
			return true
		}
	}
	return false
}

func (d *DidContract) isSenderIssued(vcId string) (bool, error) {
	senderDid, err := d.dal.getSenderDid()
	if err != nil {
		return false, err
	}

	// 查找颁发记录
	log, err := d.dal.getVcIssueLog(vcId)
	if err != nil {
		return false, err
	}

	var issueLog model.VcIssueLog
	err = json.Unmarshal(log, &issueLog)
	if err != nil {
		return false, err
	}

	return senderDid == issueLog.Issuer, nil
}
