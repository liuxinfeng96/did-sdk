package testdata

import (
	"context"
	"did-sdk/invoke"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"time"

	cmcrypto "chainmaker.org/chainmaker/common/v2/crypto"
	"chainmaker.org/chainmaker/common/v2/crypto/asym"
	"chainmaker.org/chainmaker/did-contract/model"
	acpb "chainmaker.org/chainmaker/pb-go/v2/accesscontrol"
	"chainmaker.org/chainmaker/pb-go/v2/common"
	cmsdk "chainmaker.org/chainmaker/sdk-go/v2"
	sdkutils "chainmaker.org/chainmaker/sdk-go/v2/utils"
)

const (
	ConfigPath           = "../testdata/sdk_config.yml"
	ContractByteCodePath = "../contract/ChainMakerDid.7z"
	HashType             = "SHA256"
	ContractVersion      = "1.0.0"
)

type ChainMakerTestUser struct {
	OrgId        string
	SignCertPath string
	SignKeyPath  string
}

var testUser = []*ChainMakerTestUser{
	&ChainMakerTestUser{
		OrgId:        "wx-org1.chainmaker.org",
		SignCertPath: "../testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.crt",
		SignKeyPath:  "../testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.key",
	},
	&ChainMakerTestUser{
		OrgId:        "wx-org2.chainmaker.org",
		SignCertPath: "../testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.crt",
		SignKeyPath:  "../testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.key",
	},
	&ChainMakerTestUser{
		OrgId:        "wx-org3.chainmaker.org",
		SignCertPath: "../testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.crt",
		SignKeyPath:  "../testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.key",
	},
	&ChainMakerTestUser{
		OrgId:        "wx-org4.chainmaker.org",
		SignCertPath: "../testdata/crypto-config/wx-org4.chainmaker.org/user/admin1/admin1.sign.crt",
		SignKeyPath:  "../testdata/crypto-config/wx-org4.chainmaker.org/user/admin1/admin1.sign.key",
	},
}

func GetChainmakerClient() (*cmsdk.ChainClient, error) {
	return cmsdk.NewChainClient(
		cmsdk.WithConfPath(ConfigPath),
	)
}

func InstallDidContract(client *cmsdk.ChainClient) error {
	params := make([]*common.KeyValuePair, 0)

	params = append(params, &common.KeyValuePair{
		Key:   model.Params_DidMethod,
		Value: []byte("cm"),
	})

	params = append(params, &common.KeyValuePair{
		Key:   model.Params_EnableTrustIssuer,
		Value: []byte("true"),
	})

	return InstallContract(params, client)
}

func InstallContract(kvs []*common.KeyValuePair, client *cmsdk.ChainClient) error {
	payload, err := client.CreateContractCreatePayload(invoke.DIDContractName, ContractVersion,
		ContractByteCodePath, common.RuntimeType_DOCKER_GO, kvs)
	if err != nil {
		return err
	}

	sender, err := client.SignPayload(payload)
	if err != nil {
		return err
	}

	endors := make([]*common.EndorsementEntry, 0)

	for _, u := range testUser {

		sk, err := os.ReadFile(u.SignKeyPath)
		if err != nil {
			return err
		}

		cert, err := os.ReadFile(u.SignCertPath)
		if err != nil {
			return err
		}

		sign, err := createSignatureForPayload(sk, HashType, payload)
		if err != nil {
			return err
		}

		end, err := createBlockChainIdentity(acpb.MemberType_CERT.String(), u.OrgId, cert, sign)
		if err != nil {
			return err
		}

		endors = append(endors, end)
	}

	_, err = sendContractTxRequest(payload, sender, endors, client)
	if err != nil {
		return err
	}

	return nil

}

func createSignatureForPayload(privateKey []byte, hashType string, payload *common.Payload) ([]byte, error) {
	var (
		cmKey cmcrypto.PrivateKey
		err   error
	)
	block, rest := pem.Decode(privateKey)
	if block == nil {
		cmKey, err = asym.PrivateKeyFromDER(rest)
	} else {
		cmKey, err = asym.PrivateKeyFromDER(block.Bytes)
	}
	if err != nil {
		return nil, err
	}

	hash := cmcrypto.HashAlgoMap[hashType]

	return sdkutils.SignPayloadWithHashType(cmKey, hash, payload)
}

func createBlockChainIdentity(idType string, orgId string,
	idInfo []byte, signature []byte) (*common.EndorsementEntry, error) {

	switch idType {
	case acpb.MemberType_CERT.String():
		return &common.EndorsementEntry{
			Signer: &acpb.Member{
				OrgId:      orgId,
				MemberInfo: idInfo,
				MemberType: acpb.MemberType_CERT,
			},
			Signature: signature,
		}, nil
	case acpb.MemberType_CERT_HASH.String():
		return &common.EndorsementEntry{
			Signer: &acpb.Member{
				OrgId:      orgId,
				MemberInfo: idInfo,
				MemberType: acpb.MemberType_CERT_HASH,
			},
			Signature: signature,
		}, nil
	case acpb.MemberType_PUBLIC_KEY.String():
		return &common.EndorsementEntry{
			Signer: &acpb.Member{
				OrgId:      orgId,
				MemberInfo: idInfo,
				MemberType: acpb.MemberType_PUBLIC_KEY,
			},
			Signature: signature,
		}, nil
	case acpb.MemberType_ALIAS.String():
		return &common.EndorsementEntry{
			Signer: &acpb.Member{
				OrgId:      orgId,
				MemberInfo: idInfo,
				MemberType: acpb.MemberType_ALIAS,
			},
			Signature: signature,
		}, nil
	}

	return nil, errors.New("the identity type is unknown")

}

func sendContractTxRequest(payload *common.Payload, sender *common.EndorsementEntry,
	endorsers []*common.EndorsementEntry,
	client *cmsdk.ChainClient) ([]byte, error) {

	txReq := &common.TxRequest{
		Payload:   payload,
		Endorsers: endorsers,
		Sender:    sender,
	}

	resp, err := client.SendTxRequest(txReq, -1, false)
	if err != nil {
		return nil, err
	}

	// 先判断交易执行的状态
	if resp.Code != common.TxStatusCode_SUCCESS {
		return nil, fmt.Errorf("exec tx failed, TxId: [%s], TxStatusCode: [%s], Msg: [%s]",
			resp.TxId, resp.Code.String(), resp.Message)
	}

	// 开启交易订阅
	txC, err := client.SubscribeTx(context.Background(), -1, -1, "", []string{resp.TxId})
	if err != nil {
		return nil, fmt.Errorf("[%s] subscribe tx failed, err: [%s]", resp.TxId, err.Error())
	}

	select {
	case tx, ok := <-txC:
		if !ok {
			return nil, fmt.Errorf("[%s] subscribe tx failed, tx chan is closed", resp.TxId)
		}

		txInfo, ok := tx.(*common.Transaction)
		if !ok {
			return nil, fmt.Errorf("[%s] subscribe tx failed, the tx type error", resp.TxId)
		}

		// 先判断交易执行的状态
		if txInfo.Result.Code != common.TxStatusCode_SUCCESS {
			return nil, fmt.Errorf("exec tx failed, TxId: [%s], TxStatusCode: [%s], Msg: [%s]",
				resp.TxId, resp.Code.String(), resp.Message)
		}

		// 再判断合约执行的状态（状态码在合约里定义）
		if txInfo.Result.ContractResult.Code != 0 {
			return nil, fmt.Errorf("exec contract failed, TxId: [%s], ContractCode: [%d], Msg: [%s], Result: [%s]",
				resp.TxId, resp.ContractResult.Code,
				resp.ContractResult.Message, string(resp.ContractResult.Result))
		}

		return txInfo.Result.ContractResult.Result, nil

	// 10分钟超时设置
	case <-time.After(10 * time.Minute):
		return nil, fmt.Errorf("[%s] subscribe tx failed, after 10 min, timeout", resp.TxId)
	}
}
