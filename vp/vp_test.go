package vp

import (
	"did-sdk/key"
	"did-sdk/vc"
	"testing"
	"time"

	"github.com/test-go/testify/require"
)

func TestGenerateVP(t *testing.T) {
	fieldsMap := make(map[string]string)

	fieldsMap["name"] = "姓名"
	fieldsMap["phoneNumber"] = "手机号"
	fieldsMap["idNumber"] = "身份证号"

	jsonSchema, err := vc.GenerateSimpleVcTemplate(fieldsMap)
	require.Nil(t, err)

	keyInfo, err := key.GenerateKey("SM2")
	require.Nil(t, err)

	subject := make(map[string]interface{})

	subject["name"] = "小明"
	subject["phoneNumber"] = "18705453XXX"
	subject["idNumber"] = "370687199X010200XX"

	e := time.Now().Local().Add(time.Hour * 48).Unix()
	vcBytes, err := vc.IssueVCLocal(keyInfo.SkPEM, keyInfo.Algorithm, subject,
		"did:cmdid:0xadfwfkqwfmkqm", "vc1", e, jsonSchema)
	require.Nil(t, err)

	keyInfo2, err := key.GenerateKey("EC_Secp256k1")
	require.Nil(t, err)

	vpBytes, err := GenerateVP(keyInfo.SkPEM, keyInfo2.Algorithm, "did:cmdid:vpholder", "vp1", []string{string(vcBytes)})
	require.Nil(t, err)

	println(string(vpBytes))
}
