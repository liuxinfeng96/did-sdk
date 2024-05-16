/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) DCPS. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package did

import (
	"did-sdk/invoke"
	"encoding/json"
	"strconv"

	"chainmaker.org/chainmaker/did-contract/model"
	"chainmaker.org/chainmaker/pb-go/v2/common"
	cmsdk "chainmaker.org/chainmaker/sdk-go/v2"
)

// AddTrustIssuerListToChain 在链上添加信任颁发者
// @params dids：权威颁发者DID列表
// @params client: 长安链客户端
func AddTrustIssuerListToChain(dids []string, client *cmsdk.ChainClient) error {

	didsBytes, err := json.Marshal(dids)
	if err != nil {
		return err
	}

	params := make([]*common.KeyValuePair, 0)

	params = append(params, &common.KeyValuePair{
		Key:   model.Params_DidList,
		Value: []byte(didsBytes),
	})

	_, err = invoke.InvokeContract(invoke.DIDContractName, model.Method_AddTrustIssuer, params, client)
	if err != nil {
		return err
	}

	return nil
}

// GetTrustIssuerListFromChain 从链上获取权威签发者列表
// @params didSearch：要查找的DID编号（空字符串可以查找全部列表）
// @params start：开始的索引，0表示从第一个开始
// @params count：要获取的数量，0表示获取所有
// @params client：长安链客户端
func GetTrustIssuerListFromChain(didSearch string, start int, count int,
	client *cmsdk.ChainClient) ([]string, error) {
	params := make([]*common.KeyValuePair, 0)

	params = append(params, &common.KeyValuePair{
		Key:   model.Params_DidSearch,
		Value: []byte(didSearch),
	})

	params = append(params, &common.KeyValuePair{
		Key:   model.Params_SearchStart,
		Value: []byte(strconv.Itoa(start)),
	})

	params = append(params, &common.KeyValuePair{
		Key:   model.Params_SearchCount,
		Value: []byte(strconv.Itoa(count)),
	})

	resp, err := invoke.QueryContract(invoke.DIDContractName, model.Method_GetTrustIssuer, params, client)
	if err != nil {
		return nil, err
	}

	list := make([]string, 0)

	err = json.Unmarshal(resp, &list)
	if err != nil {
		return nil, err
	}

	return list, nil
}

// DeleteTrustIssuerListFromChain
// @params dids: 要删除的did列表
// @params client: 长安链客户端
func DeleteTrustIssuerListFromChain(dids []string, client *cmsdk.ChainClient) error {
	didsBytes, err := json.Marshal(dids)
	if err != nil {
		return err
	}

	params := make([]*common.KeyValuePair, 0)

	params = append(params, &common.KeyValuePair{
		Key:   model.Params_DidList,
		Value: []byte(didsBytes),
	})

	_, err = invoke.InvokeContract(invoke.DIDContractName, model.Method_DeleteTrustIssuer, params, client)
	if err != nil {
		return err
	}

	return nil
}
