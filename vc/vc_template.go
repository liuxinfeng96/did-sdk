package vc

import "encoding/json"

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

	t := &SimpleVcTemplate{
		Schema:               "http://json-schema.org/draft-07/schema#",
		Type:                 "object",
		Properties:           properties,
		Required:             required,
		AdditionalProperties: false,
	}

	return json.Marshal(t)
}
