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

// AddDidBlackListToChain
// @params dids: did列表
// @params client: 长安链客户端
func AddDidBlackListToChain(dids []string, client *cmsdk.ChainClient) error {

	didsBytes, err := json.Marshal(dids)
	if err != nil {
		return err
	}

	params := make([]*common.KeyValuePair, 0)

	params = append(params, &common.KeyValuePair{
		Key:   model.Params_DidList,
		Value: []byte(didsBytes),
	})

	_, err = invoke.InvokeContract(invoke.DIDContractName, model.Method_AddBlackList, params, client)
	if err != nil {
		return err
	}

	return nil
}

// GetDidBlackListFromChain
// @params didSearch：要查找的DID编号（空字符串可以查找全部列表）
// @params start：开始的索引，0表示从第一个开始
// @params count：要获取的数量，0表示获取所有
// @params client：长安链客户端
func GetDidBlackListFromChain(didSearch string, start int, count int,
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

	resp, err := invoke.QueryContract(invoke.DIDContractName, model.Method_GetBlackList, params, client)
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

// DeleteDidBlackListFromChain
// @params dids: did列表
// @params client: 长安链客户端
func DeleteDidBlackListFromChain(dids []string, client *cmsdk.ChainClient) error {
	didsBytes, err := json.Marshal(dids)
	if err != nil {
		return err
	}

	params := make([]*common.KeyValuePair, 0)

	params = append(params, &common.KeyValuePair{
		Key:   model.Params_DidList,
		Value: []byte(didsBytes),
	})

	_, err = invoke.InvokeContract(invoke.DIDContractName, model.Method_DeleteBlackList, params, client)
	if err != nil {
		return err
	}

	return nil
}
