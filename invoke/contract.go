package invoke

import (
	"context"
	"fmt"
	"time"

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
	resp, err := client.InvokeContract(contractName, method, txId, params, -1, false)
	if err != nil {
		return nil, fmt.Errorf("[%s] send tx failed, err: [%s]", contractAndMethodName, err.Error())
	}

	// 先判断交易执行的状态
	if resp.Code != common.TxStatusCode_SUCCESS {
		return nil, fmt.Errorf("[%s] exec tx failed, TxId: [%s], TxStatusCode: [%s], Msg: [%s]",
			contractAndMethodName, resp.TxId, resp.Code.String(), resp.Message)
	}

	// 开启交易订阅
	txC, err := client.SubscribeTx(context.Background(), -1, -1, "", []string{txId})
	if err != nil {
		return nil, fmt.Errorf("[%s] subscribe tx failed, err: [%s]", txId, err.Error())
	}

	select {
	case tx, ok := <-txC:
		if !ok {
			return nil, fmt.Errorf("[%s] subscribe tx failed, tx chan is closed", txId)
		}

		txInfo, ok := tx.(*common.Transaction)
		if !ok {
			return nil, fmt.Errorf("[%s] subscribe tx failed, the tx type error", txId)
		}

		if txInfo.Result.Code != common.TxStatusCode_SUCCESS || txInfo.Result.ContractResult.Code != 0 {
			return nil,
				fmt.Errorf("[%s] exec contract failed, TxId: [%s], TxStatusCode: [%s], ContractCode: [%d], Result: [%s]",
					contractAndMethodName,
					txInfo.Payload.TxId,
					txInfo.Result.Code.String(),
					txInfo.Result.ContractResult.Code,
					string(txInfo.Result.ContractResult.Result))
		}

		return txInfo.Result.ContractResult.Result, nil

	// 10分钟超时设置
	case <-time.After(10 * time.Minute):
		return nil, fmt.Errorf("[%s] subscribe tx failed, after 10 min, timeout", txId)
	}
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

	if resp.Code != common.TxStatusCode_SUCCESS || resp.ContractResult.Code != 0 {
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
