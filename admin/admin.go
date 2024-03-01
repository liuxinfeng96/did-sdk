package admin

import (
	"did-sdk/invoke"
	"encoding/hex"

	cmcert "chainmaker.org/chainmaker/common/v2/cert"
	cmcrypto "chainmaker.org/chainmaker/common/v2/crypto"
	"chainmaker.org/chainmaker/did-contract/model"
	"chainmaker.org/chainmaker/pb-go/v2/common"
	cmsdk "chainmaker.org/chainmaker/sdk-go/v2"
)

// SetAdminForDidContract 为DID合约设置管理员（仅合约创建者有权限）
// @params pubKey：公钥
// @params hash：哈希算法（一般与链保持一致）
// @params client：长安链客户端
func SetAdminForDidContract(pubKey interface{}, hash cmcrypto.HashType, client *cmsdk.ChainClient) error {

	ski, err := cmcert.ComputeSKI(hash, pubKey)
	if err != nil {
		return err
	}

	skiStr := hex.EncodeToString(ski)

	params := make([]*common.KeyValuePair, 0)

	params = append(params, &common.KeyValuePair{
		Key:   model.Params_Ski,
		Value: []byte(skiStr),
	})

	_, err = invoke.InvokeContract(invoke.DIDContractName, model.Method_SetAdmin, params, client)
	if err != nil {
		return err
	}

	return nil
}

// DeleteAdminForDidContract 为DID合约删除管理员（仅合约创建者有权限）
// @params pubKey：公钥
// @params hash：哈希算法（一般与链保持一致）
// @params client：长安链客户端
func DeleteAdminForDidContract(pubKey interface{}, hash cmcrypto.HashType, client *cmsdk.ChainClient) error {
	ski, err := cmcert.ComputeSKI(hash, pubKey)
	if err != nil {
		return err
	}

	skiStr := hex.EncodeToString(ski)

	params := make([]*common.KeyValuePair, 0)

	params = append(params, &common.KeyValuePair{
		Key:   model.Params_Ski,
		Value: []byte(skiStr),
	})

	_, err = invoke.InvokeContract(invoke.DIDContractName, model.Method_DeleteAdmin, params, client)
	if err != nil {
		return err
	}

	return nil
}

// IsAdminOfDidContract 查询是否拥有合约管理员权限
// @params pubKey：公钥
// @params hash：哈希算法（一般与链保持一致）
// @params client：长安链客户端
func IsAdminOfDidContract(pubKey interface{}, hash cmcrypto.HashType, client *cmsdk.ChainClient) (bool, error) {
	ski, err := cmcert.ComputeSKI(hash, pubKey)
	if err != nil {
		return false, err
	}

	skiStr := hex.EncodeToString(ski)

	params := make([]*common.KeyValuePair, 0)

	params = append(params, &common.KeyValuePair{
		Key:   model.Params_Ski,
		Value: []byte(skiStr),
	})

	resp, err := invoke.QueryContract(invoke.DIDContractName, model.Method_IsAdmin, params, client)
	if err != nil {
		return false, err
	}

	if string(resp) == "true" {
		return true, nil
	}

	return false, nil
}
