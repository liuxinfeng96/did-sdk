package vc

import (
	"did-sdk/testdata"
	"encoding/json"
	"fmt"
	"testing"

	"chainmaker.org/chainmaker/did-contract/model"
	"github.com/test-go/testify/require"
)

func TestGenerateSimpleVcTemplate(t *testing.T) {
	fieldsMap := make(map[string]string)

	fieldsMap["name"] = "姓名"
	fieldsMap["phoneNumber"] = "手机号"
	fieldsMap["idNumber"] = "身份证号"

	jsonSchema, err := GenerateSimpleVcTemplate(fieldsMap)
	require.Nil(t, err)

	fmt.Println(string(jsonSchema))
}

func TestAddVcTemplateToChain(t *testing.T) {

	c, err := testdata.GetChainmakerClient()
	require.Nil(t, err)
	fieldsMap := make(map[string]string)

	fieldsMap["name"] = "姓名"
	fieldsMap["phoneNumber"] = "手机号"
	fieldsMap["idNumber"] = "身份证号"

	jsonSchema, err := GenerateSimpleVcTemplate(fieldsMap)
	require.Nil(t, err)

	err = AddVcTemplateToChain("1231323", "模板1", "version", jsonSchema, c)
	require.Nil(t, err)
}

func TestGetVcTemplateFromChain(t *testing.T) {
	c, err := testdata.GetChainmakerClient()
	require.Nil(t, err)
	fieldsMap := make(map[string]string)

	fieldsMap["name"] = "姓名"
	fieldsMap["phoneNumber"] = "手机号"
	fieldsMap["idNumber"] = "身份证号"

	jsonSchema, err := GenerateSimpleVcTemplate(fieldsMap)
	require.Nil(t, err)

	err = AddVcTemplateToChain("3213", "模板1", "version", jsonSchema, c)
	require.Nil(t, err)

	v, err := GetVcTemplateFromChain("3213", c)
	require.Nil(t, err)

	var schema model.VcTemplate
	err = json.Unmarshal(v, &schema)
	require.Nil(t, err)

	require.Equal(t, jsonSchema, []byte(schema.Template))
}

func TestGetVcTemplateListFromChain(t *testing.T) {
	c, err := testdata.GetChainmakerClient()
	require.Nil(t, err)

	list, err := GetVcTemplateListFromChain("", 0, 0, c)
	require.Nil(t, err)

	fmt.Println(list)
}
