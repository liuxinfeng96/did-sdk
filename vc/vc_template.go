/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) DCPS. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package vc

import (
	"did-sdk/invoke"
	"encoding/json"
	"strconv"

	"chainmaker.org/chainmaker/did-contract/model"
	"chainmaker.org/chainmaker/pb-go/v2/common"
	cmsdk "chainmaker.org/chainmaker/sdk-go/v2"
)

// SimpleVcTemplate 简易的JSON Schema的VC模板
type SimpleVcTemplate struct {
	Schema               string                          `json:"$schema"`
	Type                 string                          `json:"type"`
	Properties           map[string]*SimplePropertyField `json:"properties"`
	Required             []string                        `json:"required"`
	AdditionalProperties bool                            `json:"additionalProperties"`
}

// SimplePropertyField 简易的JSON Schema的VC模板里的Property定义
type SimplePropertyField struct {
	Type  string `json:"type"`
	Title string `json:"title"`
}

// GenerateSimpleVcTemplate 生成字段都是String类型的简易的VC模板
// 注：复杂模板需要根据业务灵活定义，需要符合JSON Schema规范
// @params fieldsMap: key: 字段名 value: 具体含义
// @return JSON Schema 模板
func GenerateSimpleVcTemplate(fieldsMap map[string]string) ([]byte, error) {

	properties := make(map[string]*SimplePropertyField)
	required := make([]string, 0)
	for k, v := range fieldsMap {
		properties[k] = &SimplePropertyField{
			Title: v,
			Type:  "string",
		}

		required = append(required, k)
	}

	// 默认添加id字段
	_, ok := properties["id"]
	if !ok {
		properties["id"] = &SimplePropertyField{
			Type:  "string",
			Title: "DID",
		}

		required = append(required, "id")
	}

	t := &SimpleVcTemplate{
		Schema:               "http://json-schema.org/draft-07/schema#",
		Type:                 "object",
		Properties:           properties,
		Required:             required,
		AdditionalProperties: false,
	}

	return json.Marshal(t)
}

// AddVcTemplateToChain VC模板上链
// @params id：模板ID
// @params name：模板名称
// @params version：模板版本
// @params template：模板内容，需要JSON schema格式
// @params client：长安链客户端
func AddVcTemplateToChain(id string, name string, version string, template []byte, client *cmsdk.ChainClient) error {
	params := make([]*common.KeyValuePair, 0)

	params = append(params, &common.KeyValuePair{
		Key:   model.Params_VcTemplateId,
		Value: []byte(id),
	})

	params = append(params, &common.KeyValuePair{
		Key:   model.Params_VcTemplateName,
		Value: []byte(name),
	})

	params = append(params, &common.KeyValuePair{
		Key:   model.Params_VcTemplateVersion,
		Value: []byte(version),
	})

	params = append(params, &common.KeyValuePair{
		Key:   model.Params_VcTemplate,
		Value: json.RawMessage(template),
	})

	_, err := invoke.InvokeContract(invoke.DIDContractName, model.Method_SetVcTemplate, params, client)
	if err != nil {
		return err
	}

	return nil
}

// GetVcTemplateFromChain 从链上获取VC模板
// @params id：模板ID
// @params client：长安链客户端
func GetVcTemplateFromChain(id string, client *cmsdk.ChainClient) ([]byte, error) {
	params := make([]*common.KeyValuePair, 0)

	params = append(params, &common.KeyValuePair{
		Key:   model.Params_VcTemplateId,
		Value: []byte(id),
	})

	return invoke.QueryContract(invoke.DIDContractName, model.Method_GetVcTemplate, params, client)
}

// GetVcTemplateListFromChain 从链上获取VC模板列表
// @params nameSearch：模板名称关键字（空字符串可以查找全部列表）
// @params start：开始的索引，0表示从第一个开始
// @params count：要获取的数量，0表示获取所有
// @params client：长安链客户端
func GetVcTemplateListFromChain(nameSearch string, start int, count int,
	client *cmsdk.ChainClient) ([]*model.VcTemplate, error) {
	params := make([]*common.KeyValuePair, 0)

	params = append(params, &common.KeyValuePair{
		Key:   model.Params_VcTemplateNameSearch,
		Value: []byte(nameSearch),
	})

	params = append(params, &common.KeyValuePair{
		Key:   model.Params_SearchStart,
		Value: []byte(strconv.Itoa(start)),
	})

	params = append(params, &common.KeyValuePair{
		Key:   model.Params_SearchCount,
		Value: []byte(strconv.Itoa(count)),
	})

	resp, err := invoke.QueryContract(invoke.DIDContractName, model.Method_GetVcTemplateList, params, client)
	if err != nil {
		return nil, err
	}

	var vcTemplateList []*model.VcTemplate

	err = json.Unmarshal(resp, &vcTemplateList)
	if err != nil {
		return nil, err
	}

	return vcTemplateList, nil
}
