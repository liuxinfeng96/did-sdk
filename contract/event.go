/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) DCPS. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"did-contract/model"

	"chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"
)

// 发送设置DID Document事件
func emitSetDidDocumentEvent(did string, didDocument string) {
	sdk.Instance.EmitEvent(model.Topic_SetDidDocument, []string{did, didDocument})
}

// 发送添加黑名单事件
func emitAddBlackListEvent(dids []string) {
	sdk.Instance.EmitEvent(model.Topic_AddBlackList, dids)
}

// 发送删除黑名单事件
func emitDeleteBlackListEvent(dids []string) {
	sdk.Instance.EmitEvent(model.Topic_DeleteBlackList, dids)
}

// 发送添加信任发行者事件
func emitAddTrustIssuerListEvent(dids []string) {
	sdk.Instance.EmitEvent(model.Topic_AddTrustIssuer, dids)
}

// 发送删除信任发行者事件
func emitDeleteTrustIssuerEvent(dids []string) {
	sdk.Instance.EmitEvent(model.Topic_DeleteTrustIssuer, dids)
}

// 发送撤销VC事件
func emitRevokeVcEvent(vcID string) {
	sdk.Instance.EmitEvent(model.Topic_RevokeVc, []string{vcID})
}

// 发送设置VC模板事件
func emitSetVcTemplateEvent(templateId string, vcTemplate []byte) {
	sdk.Instance.EmitEvent(model.Topic_SetVcTemplate, []string{templateId, string(vcTemplate)})
}

// 发送记录VC签发日志事件
func emitVcIssueLogEvent(vcId string, log []byte) {
	sdk.Instance.EmitEvent(model.Topic_VcIssueLog, []string{vcId, string(log)})
}
