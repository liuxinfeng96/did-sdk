package vc

import (
	"did-sdk/key"
	"testing"
	"time"

	"github.com/test-go/testify/require"
)

func TestIssueVCLocal(t *testing.T) {
	fieldsMap := make(map[string]string)

	fieldsMap["name"] = "姓名"
	fieldsMap["phoneNumber"] = "手机号"
	fieldsMap["idNumber"] = "身份证号"

	jsonSchema, err := GenerateSimpleVcTemplate(fieldsMap)
	require.Nil(t, err)

	keyInfo, err := key.GenerateKey("SM2")
	require.Nil(t, err)

	subject := make(map[string]interface{})

	subject["name"] = "小明"
	subject["phoneNumber"] = "18705453XXX"
	subject["idNumber"] = "370687199X010200XX"

	e := time.Now().Local().Add(time.Hour * 48).Unix()
	vcBytes, err := IssueVCLocal(keyInfo.SkPEM, keyInfo.Algorithm, subject,
		"did:cmdid:0xadfwfkqwfmkqm", "vc1", e, jsonSchema)
	require.Nil(t, err)

	println(string(vcBytes))

	// 不符合VC模板测试

	subject2 := make(map[string]interface{})

	subject2["name"] = "小明"
	subject2["phoneNumber"] = "18705453XXX"
	subject2["idNumber"] = "370687199X010200XX"
	subject2["test"] = "test"

	e = time.Now().Local().Add(time.Hour * 48).Unix()
	_, err = IssueVCLocal(keyInfo.SkPEM, keyInfo.Algorithm, subject2,
		"did:cmdid:0xadfwfkqwfmkqm", "vc1", e, jsonSchema)
	require.NotNil(t, err)
	println(err.Error())

}
