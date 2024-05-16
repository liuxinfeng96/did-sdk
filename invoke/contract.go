/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) DCPS. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package invoke

import (
	"fmt"

	"chainmaker.org/chainmaker/pb-go/v2/common"
	cmsdk "chainmaker.org/chainmaker/sdk-go/v2"
	sdkutils "chainmaker.org/chainmaker/sdk-go/v2/utils"
)

// DIDContractName this contract name
const DIDContractName = "ChainMakerDid"

// InvokeContract 基于ChainMakerSDK包装的合约调用接口，使用监听交易的方式拿到交易结果
// @params contractName: 合约名称
// @params method: 方法名称
// @params params: 合约调用参数
// @params client: 长安链SDK客户端连接
// @return 交易里的结果
func InvokeContract(contractName, method string,
	params []*common.KeyValuePair, client *cmsdk.ChainClient) ([]byte, error) {

	// 生成交易ID
	txId := sdkutils.GetRandTxId()

	contractAndMethodName := contractName + "-" + method

	// 调用SDK不同步结果Invoke接口
	resp, err := client.InvokeContract(contractName, method, txId, params, -1, true)
	if err != nil {
		return nil, fmt.Errorf("[%s] send tx failed, err: [%s]", contractAndMethodName, err.Error())
	}

	if resp.Code != common.TxStatusCode_SUCCESS {
		if resp.ContractResult == nil {
			return nil,
				fmt.Errorf("[%s] exec contract failed, TxId: [%s], TxStatusCode: [%s], TxMsg: [%s]",
					contractAndMethodName,
					resp.TxId,
					resp.Code.String(),
					resp.Message)
		}

		return nil,
			fmt.Errorf("[%s] exec contract failed, TxId: [%s], TxStatusCode: [%s], ContractCode: [%d], Result: [%s]",
				contractAndMethodName,
				resp.TxId,
				resp.Code.String(),
				resp.ContractResult.Code,
				string(resp.ContractResult.Result))
	}

	return resp.ContractResult.Result, nil

}

// QueryContract 基于ChainMakerSDK包装的合约调用接口，仅查询使用，不落块
// @params contractName: 合约名称
// @params method: 方法名称
// @params params: 合约调用参数
// @params client: 长安链SDK客户端连接
// @return 交易里的结果
func QueryContract(contractName, method string,
	params []*common.KeyValuePair, client *cmsdk.ChainClient) ([]byte, error) {

	contractAndMethodName := contractName + "-" + method

	resp, err := client.QueryContract(contractName, method, params, -1)
	if err != nil {
		return nil, fmt.Errorf("[%s] send tx failed, err: [%s]", contractAndMethodName, err.Error())
	}

	if resp.Code != common.TxStatusCode_SUCCESS {
		if resp.ContractResult == nil {
			return nil,
				fmt.Errorf("[%s] exec contract failed, TxId: [%s], TxStatusCode: [%s], TxMsg: [%s]",
					contractAndMethodName,
					resp.TxId,
					resp.Code.String(),
					resp.Message)
		}

		return nil,
			fmt.Errorf("[%s] exec contract failed, TxId: [%s], TxStatusCode: [%s], ContractCode: [%d], Result: [%s]",
				contractAndMethodName,
				resp.TxId,
				resp.Code.String(),
				resp.ContractResult.Code,
				string(resp.ContractResult.Result))
	}

	return resp.ContractResult.Result, nil
}
