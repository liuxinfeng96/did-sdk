/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) DCPS. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package model

const (
	// Method_DidMethod method "DidMethod"
	Method_DidMethod = "DidMethod"
	// Method_IsValidDid method "IsValidDid"
	Method_IsValidDid = "IsValidDid"
	// Method_AddDidDocument method "AddDidDocument"
	Method_AddDidDocument = "AddDidDocument"
	// Method_UpdateDidDocument method "UpdateDidDocument"
	Method_UpdateDidDocument = "UpdateDidDocument"
	// Method_GetDidDocument method "GetDidDocument"
	Method_GetDidDocument = "GetDidDocument"
	// Method_GetDidByPubKey method "GetDidByPubKey"
	Method_GetDidByPubKey = "GetDidByPubKey"
	// Method_GetDidByAddress method "GetDidByAddress"
	Method_GetDidByAddress = "GetDidByAddress"
	// Method_AddBlackList method "AddBlackList"
	Method_AddBlackList = "AddBlackList"
	// Method_GetBlackList method "GetBlackList"
	Method_GetBlackList = "GetBlackList"
	// Method_DeleteBlackList method "DeleteBlackList"
	Method_DeleteBlackList = "DeleteBlackList"
	// Method_AddTrustIssuer method "AddTrustIssuer"
	Method_AddTrustIssuer = "AddTrustIssuer"
	// Method_GetTrustIssuer method "GetTrustIssuer"
	Method_GetTrustIssuer = "GetTrustIssuer"
	// Method_DeleteTrustIssuer method "DeleteTrustIssuer"
	Method_DeleteTrustIssuer = "DeleteTrustIssuer"
	// Method_RevokeVc method "RevokeVc"
	Method_RevokeVc = "RevokeVc"
	// Method_GetRevokedVcList method "GetRevokedVcList"
	Method_GetRevokedVcList = "GetRevokedVcList"
	// Method_SetVcTemplate method "SetVcTemplate"
	Method_SetVcTemplate = "SetVcTemplate"
	// Method_GetVcTemplate method "GetVcTemplate"
	Method_GetVcTemplate = "GetVcTemplate"
	// Method_GetVcTemplateList method "GetVcTemplateList"
	Method_GetVcTemplateList = "GetVcTemplateList"
	// Method_VerifyVc method "VerifyVc"
	Method_VerifyVc = "VerifyVc"
	// Method_VerifyVp method "VerifyVp"
	Method_VerifyVp = "VerifyVp"
	// Method_SetAdmin method "SetAdmin"
	Method_SetAdmin = "SetAdmin"
	// Method_DeleteAdmin method "DeleteAdmin"
	Method_DeleteAdmin = "DeleteAdmin"
	// Method_IsAdmin method "IsAdmin"
	Method_IsAdmin = "IsAdmin"
	// Method_VcIssueLog method "VcIssueLog"
	Method_VcIssueLog = "VcIssueLog"
	// Method_GetVcIssueLogs method "GetVcIssueLogs"
	Method_GetVcIssueLogs = "GetVcIssueLogs"
)

const (
	// Topic_SetDidDocument contract event topic "SetDidDocument"
	Topic_SetDidDocument = "DidTopic_SetDidDocument"
	// Topic_AddBlackList contract event topic "AddBlackList"
	Topic_AddBlackList = "DidTopic_AddBlackList"
	// Topic_DeleteBlackList contract event topic "DeleteBlackList"
	Topic_DeleteBlackList = "DidTopic_DeleteBlackList"
	// Topic_AddTrustIssuer contract event topic "AddTrustIssuer"
	Topic_AddTrustIssuer = "DidTopic_AddTrustIssuer"
	// Topic_DeleteTrustIssuer contract event topic "DeleteTrustIssuer"
	Topic_DeleteTrustIssuer = "DidTopic_DeleteTrustIssuer"
	// Topic_RevokeVc contract event topic "RevokeVc"
	Topic_RevokeVc = "DidTopic_RevokeVc"
	// Topic_SetVcTemplate contract event topic "SetVcTemplate"
	Topic_SetVcTemplate = "DidTopic_SetVcTemplate"
	// Topic_VcIssueLog event topic "VcIssueLog"
	Topic_VcIssueLog = "DidTopic_VcIssueLog"
)

const (
	// Params_DidMethod parameter of the contract method
	Params_DidMethod = "didMethod"
	// Params_EnableTrustIssuer parameter of the contract method
	Params_EnableTrustIssuer = "enableTrustIssuer"
	// Params_DidDocument parameter of the contract method
	Params_DidDocument = "didDocument"
	// Params_Did parameter of the contract method
	Params_Did = "did"
	// Params_DidList parameter of the contract method
	Params_DidList = "dids"
	// Params_DidSearch parameter of the contract method
	Params_DidSearch = "didSearch"
	// Params_SearchStart parameter of the contract method
	Params_SearchStart = "start"
	// Params_SearchCount parameter of the contract method
	Params_SearchCount = "count"
	// Params_DidPubkey parameter of the contract method
	Params_DidPubkey = "pubKey"
	// Params_DidAddress parameter of the contract method
	Params_DidAddress = "address"
	// Params_VcId parameter of the contract method
	Params_VcId = "vcId"
	// Params_VcIdSearch parameter of the contract method
	Params_VcIdSearch = "vcIdSearch"
	// Params_VcTemplateId parameter of the contract method
	Params_VcTemplateId = "vcTemplateId"
	// Params_VcTemplateName parameter of the contract method
	Params_VcTemplateName = "vcTemplateName"
	// Params_VcTemplate parameter of the contract method
	Params_VcTemplate = "vcTemplate"
	// Params_VcTemplateVersion parameter of the contract method
	Params_VcTemplateVersion = "vcTemplateVersion"
	// Params_VcTemplateNameSearch parameter of the contract method
	Params_VcTemplateNameSearch = "vcTemplateNameSearch"
	// Params_VcJson parameter of the contract method
	Params_VcJson = "vcJson"
	// Params_VpJson parameter of the contract method
	Params_VpJson = "vpJson"
	// Params_Ski parameter of the contract method
	Params_Ski = "ski"
	// Params_Issuer parameter of the contract method
	Params_Issuer = "issuer"
)
