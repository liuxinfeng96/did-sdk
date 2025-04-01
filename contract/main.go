/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) DCPS. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"did-contract/model"

	"chainmaker.org/chainmaker/contract-sdk-go/v2/pb/protogo"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sandbox"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"
)

type DidContract struct {
	dal *Dal
}

func (d *DidContract) InitDidContract(didMethod string, enableTrustIssuer bool) error {
	d.dal = NewDal()

	// 需要存储，防止虚拟机重启内存丢失
	err := d.dal.putDidMethod(didMethod)
	if err != nil {
		return err
	}

	if enableTrustIssuer {
		err = d.dal.putEnableTrustIssuer("true")
		if err != nil {
			return err
		}
	} else {
		err = d.dal.putEnableTrustIssuer("false")
		if err != nil {
			return err
		}
	}

	return nil
}

// InitContract install contract func
func (d *DidContract) InitContract() protogo.Response {
	method, err := RequireString(model.Params_DidMethod)
	if err != nil {
		return sdk.Error(err.Error())
	}

	enableTrustIssuer, err := RequireBool(model.Params_EnableTrustIssuer)
	if err != nil {
		return sdk.Error(err.Error())
	}

	err = d.InitDidContract(method, enableTrustIssuer)
	if err != nil {
		return sdk.Error(err.Error())
	}

	return sdk.SuccessResponse
}

// UpgradeContract upgrade contract func
func (d *DidContract) UpgradeContract() protogo.Response {
	method, err := RequireString(model.Params_DidMethod)
	if err != nil {
		return sdk.Error(err.Error())
	}

	enableTrustIssuer, err := RequireBool(model.Params_EnableTrustIssuer)
	if err != nil {
		return sdk.Error(err.Error())
	}

	err = d.InitDidContract(method, enableTrustIssuer)
	if err != nil {
		return sdk.Error(err.Error())
	}

	return sdk.SuccessResponse
}

// InvokeContract the entry func of invoke contract func
func (d *DidContract) InvokeContract(method string) (result protogo.Response) { //nolint
	// 记录异常结果日志
	defer func() {
		if result.Status != 0 {
			sdk.Instance.Warnf(result.Message)
		}
	}()

	switch method {
	case model.Method_DidMethod:
		return ReturnString(d.DidMethod())
	case model.Method_IsValidDid:
		did, err := RequireString(model.Params_Did)
		if err != nil {
			return sdk.Error(err.Error())
		}
		return ReturnBool(d.IsValidDid(did))
	case model.Method_AddDidDocument:
		didDocument, err := RequireString(model.Params_DidDocument)
		if err != nil {
			return sdk.Error(err.Error())
		}
		return Return(d.AddDidDocument(didDocument))
	case model.Method_GetDidDocument:
		did, err := RequireString(model.Params_Did)
		if err != nil {
			return sdk.Error(err.Error())
		}
		return ReturnString(d.GetDidDocument(did))
	case model.Method_UpdateDidDocument:
		didDocument, err := RequireString(model.Params_DidDocument)
		if err != nil {
			return sdk.Error(err.Error())
		}
		return Return(d.UpdateDidDocument(didDocument))
	case model.Method_GetDidByPubKey:
		pubKey, err := RequireString(model.Params_DidPubkey)
		if err != nil {
			return sdk.Error(err.Error())
		}
		return ReturnString(d.GetDidByPubkey(pubKey))
	case model.Method_GetDidByAddress:
		address, err := RequireString(model.Params_DidAddress)
		if err != nil {
			return sdk.Error(err.Error())
		}
		return ReturnString(d.GetDidByAddress(address))
	case model.Method_AddBlackList:
		dids, err := RequireStringList(model.Params_Did, model.Params_DidList)
		if err != nil {
			return sdk.Error(err.Error())
		}
		return Return(d.AddBlackList(dids))
	case model.Method_DeleteBlackList:
		dids, err := RequireStringList(model.Params_Did, model.Params_DidList)
		if err != nil {
			return sdk.Error(err.Error())
		}
		return Return(d.DeleteBlackList(dids))
	case model.Method_GetBlackList:
		args := sdk.Instance.GetArgs()
		didSearch := args[model.Params_DidSearch]
		start := OptionInt(model.Params_SearchStart, 1)
		count := OptionInt(model.Params_SearchCount, 1000)
		return ReturnJson(d.GetBlackList(string(didSearch), start, count))
	case model.Method_RevokeVc:
		vcId, err := RequireString(model.Params_VcId)
		if err != nil {
			return sdk.Error(err.Error())
		}
		return Return(d.RevokeVc(vcId))
	case model.Method_GetRevokedVcList:
		args := sdk.Instance.GetArgs()
		vcIdSearch := args[model.Params_VcIdSearch]
		start := OptionInt(model.Params_SearchStart, 1)
		count := OptionInt(model.Params_SearchCount, 1000)
		return ReturnJson(d.GetRevokedVcList(string(vcIdSearch), start, count))
	case model.Method_SetVcTemplate:
		templateId, err := RequireString(model.Params_VcTemplateId)
		if err != nil {
			return sdk.Error(err.Error())
		}
		templateName, err := RequireString(model.Params_VcTemplateName)
		if err != nil {
			return sdk.Error(err.Error())
		}
		vcTemplate, err := RequireString(model.Params_VcTemplate)
		if err != nil {
			return sdk.Error(err.Error())
		}
		version, err := RequireString(model.Params_VcTemplateVersion)
		if err != nil {
			return sdk.Error(err.Error())
		}
		return Return(d.SetVcTemplate(templateId, templateName, version, vcTemplate))
	case model.Method_GetVcTemplate:
		templateId, err := RequireString(model.Params_VcTemplateId)
		if err != nil {
			return sdk.Error(err.Error())
		}
		return ReturnBytes(d.GetVcTemplate(templateId))
	case model.Method_GetVcTemplateList:
		args := sdk.Instance.GetArgs()
		nameSearch := args[model.Params_VcTemplateNameSearch]
		start := OptionInt(model.Params_SearchStart, 1)
		count := OptionInt(model.Params_SearchCount, 1000)
		return ReturnJson(d.GetVcTemplateList(string(nameSearch), start, count))
	case model.Method_VerifyVc:
		vcJson, err := RequireString(model.Params_VcJson)
		if err != nil {
			return sdk.Error(err.Error())
		}
		return ReturnBool(d.VerifyVc(vcJson))
	case model.Method_VerifyVp:
		vpJson, err := RequireString(model.Params_VpJson)
		if err != nil {
			return sdk.Error(err.Error())
		}
		return ReturnBool(d.VerifyVp(vpJson))
	case model.Method_SetAdmin:
		ski, err := RequireString(model.Params_Ski)
		if err != nil {
			return sdk.Error(err.Error())
		}
		return Return(d.SetAdmin(ski))
	case model.Method_DeleteAdmin:
		ski, err := RequireString(model.Params_Ski)
		if err != nil {
			return sdk.Error(err.Error())
		}
		return Return(d.DeleteAdmin(ski))
	case model.Method_IsAdmin:
		ski, err := RequireString(model.Params_Ski)
		if err != nil {
			return sdk.Error(err.Error())
		}
		ok := d.IsAdmin(ski)
		return ReturnBool(ok, nil)
	case model.Method_VcIssueLog:
		issuer, err := RequireString(model.Params_Issuer)
		if err != nil {
			return sdk.Error(err.Error())
		}

		did, err := RequireString(model.Params_Did)
		if err != nil {
			return sdk.Error(err.Error())
		}

		vcId, err := RequireString(model.Params_VcId)
		if err != nil {
			return sdk.Error(err.Error())
		}

		vcTemplateId, err := RequireString(model.Params_VcTemplateId)
		if err != nil {
			return sdk.Error(err.Error())
		}

		return Return(d.VcIssueLog(issuer, did, vcTemplateId, vcId))

	case model.Method_GetVcIssueLogs:
		args := sdk.Instance.GetArgs()
		vcIdSearch := args[model.Params_VcIdSearch]
		start := OptionInt(model.Params_SearchStart, 1)
		count := OptionInt(model.Params_SearchCount, 1000)
		return ReturnJson(d.GetVcIssueLogs(string(vcIdSearch), start, count))
	}

	enableTrustIssuer, err := d.dal.getEnableTrustIssuer()
	if err != nil {
		return sdk.Error(err.Error())
	}

	if enableTrustIssuer == "true" {
		switch method {
		case model.Method_AddTrustIssuer:
			dids, err := RequireStringList(model.Params_Did, model.Params_DidList)
			if err != nil {
				return sdk.Error(err.Error())
			}
			return Return(d.AddTrustIssuerList(dids))
		case model.Method_DeleteTrustIssuer:
			dids, err := RequireStringList(model.Params_Did, model.Params_DidList)
			if err != nil {
				return sdk.Error(err.Error())
			}
			return Return(d.DeleteTrustIssuer(dids))
		case model.Method_GetTrustIssuer:
			args := sdk.Instance.GetArgs()
			didSearch := args[model.Params_DidSearch]
			start := OptionInt(model.Params_SearchStart, 1)
			count := OptionInt(model.Params_SearchCount, 1000)
			return ReturnJson(d.GetTrustIssuer(string(didSearch), start, count))
		}
	}

	return sdk.Error("invalid method")
}

func main() {
	err := sandbox.Start(new(DidContract))
	if err != nil {
		sdk.Instance.Errorf(err.Error())
	}
}
