package did

import (
	"did-sdk/invoke"
	"encoding/json"
	"strconv"

	"chainmaker.org/chainmaker/did-contract/model"
	"chainmaker.org/chainmaker/pb-go/v2/common"
	cmsdk "chainmaker.org/chainmaker/sdk-go/v2"
)

// AddTrustIssuerListToChain
// @params dids：权威颁发者DID列表
// @params client: 长安链客户端
func AddTrustIssuerListToChain(dids []string, client *cmsdk.ChainClient) error {

	didsBytes, err := json.Marshal(dids)
	if err != nil {
		return err
	}

	params := make([]*common.KeyValuePair, 0)

	params = append(params, &common.KeyValuePair{
		Key:   "dids",
		Value: []byte(didsBytes),
	})

	_, err = invoke.InvokeContract(invoke.DIDContractName, model.Method_AddTrustIssuer, params, client)
	if err != nil {
		return err
	}

	return nil
}

// GetTrustIssuerListFromChain
// @params didSearch：搜索的DID关键字
// @params start：开始的索引，默认从0开始
// @params count：要获取的数量，默认0表示获取所有
// @params client：长安链客户端
func GetTrustIssuerListFromChain(didSearch string, start int, count int,
	client *cmsdk.ChainClient) ([]string, error) {
	params := make([]*common.KeyValuePair, 0)

	params = append(params, &common.KeyValuePair{
		Key:   "didSearch",
		Value: []byte(didSearch),
	})

	params = append(params, &common.KeyValuePair{
		Key:   "start",
		Value: []byte(strconv.Itoa(start)),
	})

	params = append(params, &common.KeyValuePair{
		Key:   "count",
		Value: []byte(strconv.Itoa(count)),
	})

	resp, err := invoke.InvokeContract(invoke.DIDContractName, model.Method_GetTrustIssuer, params, client)
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
		Key:   "dids",
		Value: []byte(didsBytes),
	})

	_, err = invoke.InvokeContract(invoke.DIDContractName, model.Method_DeleteTrustIssuer, params, client)
	if err != nil {
		return err
	}

	return nil
}
