package main

import (
	"did-contract/model"
	"encoding/json"

	"chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"
)

// 发送设置DID Document事件
func emitSetDidDocumentEvent(did string, didDocument string) {
	sdk.Instance.EmitEvent(model.Topic_SetDidDocument, []string{did, didDocument})
}

// 发送添加黑名单事件
func emitAddBlackListEvent(dids []string) {
	value, _ := json.Marshal(dids)
	sdk.Instance.EmitEvent(model.Topic_AddBlackList, []string{string(value)})
}

// 发送删除黑名单事件
func emitDeleteBlackListEvent(dids []string) {
	value, _ := json.Marshal(dids)
	sdk.Instance.EmitEvent(model.Topic_DeleteBlackList, []string{string(value)})
}

// 发送添加信任发行者事件
func emitAddTrustIssuerListEvent(dids []string) {
	for _, did := range dids {
		sdk.Instance.EmitEvent(model.Topic_AddTrustIssuer, []string{did})
	}
}

// 发送删除信任发行者事件
func emitDeleteTrustIssuerEvent(dids []string) {
	for _, did := range dids {
		sdk.Instance.EmitEvent(model.Topic_DeleteTrustIssuer, []string{did})
	}
}

// 发送撤销VC事件
func emitRevokeVcEvent(vcID string) {
	sdk.Instance.EmitEvent(model.Topic_RevokeVc, []string{vcID})
}

// 发送设置VC模板事件
func emitSetVcTemplateEvent(templateId string, templateName, version string, vcTemplate string) {
	sdk.Instance.EmitEvent(model.Method_SetVcTemplate, []string{templateId, templateName, version, vcTemplate})
}
