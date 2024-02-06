package main

import (
	"did-contract/model"

	"chainmaker.org/chainmaker/contract-sdk-go/v2/pb/protogo"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sandbox"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"
)

type DidContract struct {
	didMethod         string
	dal               *Dal
	enableTrustIssuer bool
}

func (d *DidContract) InitDidContract(didMethod string, enableTrustIssuer bool) {
	d.didMethod = didMethod
	d.dal = NewDal(didMethod)
	d.enableTrustIssuer = enableTrustIssuer
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

	d.InitDidContract(method, enableTrustIssuer)

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

	d.InitDidContract(method, enableTrustIssuer)

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
		return sdk.Success([]byte(d.DidMethod()))
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
		didSearch := args[model.Params_VcTemplateNameSearch]
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
		vcIdSearch := args[model.Params_VcTemplateNameSearch]
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
		start := OptionInt("start", 1)
		count := OptionInt("count", 1000)
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
	}

	if d.enableTrustIssuer {
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
