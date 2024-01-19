package contract

import (
	"context"
	"fmt"
	"time"

	"chainmaker.org/chainmaker/pb-go/v2/common"
	cmsdk "chainmaker.org/chainmaker/sdk-go/v2"
	sdkutils "chainmaker.org/chainmaker/sdk-go/v2/utils"
)

const (
	// Contract_Did DID's contract name
	Contract_Did = "ChainMakerDid"
	// Method_DidMethod method "DidMethod"
	Method_DidMethod = "DidMethod"
	// Method_AddDidDocument method "AddDidDocument"
	Method_AddDidDocument = "AddDidDocument"
	// Method_IsValidDid method "IsValidDid"
	Method_IsValidDid = "IsValidDid"
	// Method_GetDidDocument method "GetDidDocument"
	Method_GetDidDocument = "GetDidDocument"
	// Method_GetDidByPubkey method "GetDidByPubkey"
	Method_GetDidByPubkey = "GetDidByPubkey"
	// Method_GetDidByAddress method "GetDidByAddress"
	Method_GetDidByAddress = "GetDidByAddress"
)

// DealTxResponse parse the response returned by the transaction
// @params resp: the chainmaker tx response
// @params contractName
// @params method
// @return the contract result from the response
func DealTxResponse(resp *common.TxResponse, contractName, method string) ([]byte, error) {
	contractAndMethodName := contractName + "-" + method

	// 先判断交易执行的状态
	if resp.Code != common.TxStatusCode_SUCCESS {
		return nil, fmt.Errorf("[%s] exec tx failed, TxId: [%s], TxStatusCode: [%s], Msg: [%s]",
			contractAndMethodName, resp.TxId, resp.Code.String(), resp.Message)
	}

	// 再判断合约执行的状态（状态码在合约里定义）
	if resp.ContractResult.Code != 0 {
		return nil, fmt.Errorf("[%s] exec contract failed, TxId: [%s], ContractCode: [%d], Msg: [%s], Result: [%s]",
			contractAndMethodName, resp.TxId, resp.ContractResult.Code,
			resp.ContractResult.Message, string(resp.ContractResult.Result))
	}

	return resp.ContractResult.Result, nil
}

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

		// 先判断交易执行的状态
		if txInfo.Result.Code != common.TxStatusCode_SUCCESS {
			return nil, fmt.Errorf("[%s] exec tx failed, TxId: [%s], TxStatusCode: [%s], Msg: [%s]",
				contractAndMethodName, resp.TxId, resp.Code.String(), resp.Message)
		}

		// 再判断合约执行的状态（状态码在合约里定义）
		if txInfo.Result.ContractResult.Code != 0 {
			return nil, fmt.Errorf("[%s] exec contract failed, TxId: [%s], ContractCode: [%d], Msg: [%s], Result: [%s]",
				contractAndMethodName, resp.TxId, resp.ContractResult.Code,
				resp.ContractResult.Message, string(resp.ContractResult.Result))
		}

		return txInfo.Result.ContractResult.Result, nil

	// 10分钟超时设置
	case <-time.After(10 * time.Minute):
		return nil, fmt.Errorf("[%s] subscribe tx failed, after 10 min, timeout", txId)
	}
}
