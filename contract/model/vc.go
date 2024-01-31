package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"
	"github.com/buger/jsonparser"
	jsonschema "github.com/xeipuuv/gojsonschema"
)

// VerifiableCredential VC的结构内容定义
type VerifiableCredential struct {
	rawData           json.RawMessage
	Context           []string               `json:"@context"`
	Id                string                 `json:"id"`
	Type              []string               `json:"type"`
	CredentialSubject map[string]interface{} `json:"credentialSubject"`
	Issuer            string                 `json:"issuer"`
	IssuanceDate      string                 `json:"issuanceDate"`
	ExpirationDate    string                 `json:"expirationDate"`
	Template          *struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"template,omitempty"`
	Proof *Proof `json:"proof,omitempty"`
}

// NewVerifiableCredential 根据VC凭证json字符串创建VC凭证
func NewVerifiableCredential(vcJson string) (*VerifiableCredential, error) {
	var vc VerifiableCredential
	err := json.Unmarshal([]byte(vcJson), &vc)
	if err != nil {
		return nil, err
	}
	vc.rawData = []byte(vcJson)
	return &vc, nil
}

// GetCredentialSubjectID 获取VC凭证的持有者DID
func (vc *VerifiableCredential) GetCredentialSubjectID() string {
	return vc.CredentialSubject["id"].(string)
}

func (vc *VerifiableCredential) Verify(pkPem, template []byte) (bool, error) {

	// Check if the VC type is correct
	if len(vc.Type) == 0 || vc.Type[0] != "VerifiableCredential" {
		return false, errors.New("invalid VC type")
	}

	issuanceDate, err := time.Parse(time.RFC3339, vc.IssuanceDate)
	if err != nil {
		return false, err
	}

	expirationDate, err := time.Parse(time.RFC3339, vc.ExpirationDate)
	if err != nil {
		return false, err
	}

	if issuanceDate.After(expirationDate) {
		return false, errors.New("issuance date is after the expiration date")
	}

	// 检查当前时间是否在有效期内
	myTime, err := getTxTime()
	if err != nil {
		return false, err
	}

	if myTime < issuanceDate.Unix() || myTime > expirationDate.Unix() {
		return false, errors.New("the verifiable credential has expired")
	}

	// 验证模板字段
	if template != nil {
		ok, err := vc.verifyCredentialSubject(template)
		if !ok {
			return false, fmt.Errorf("credential subject verified failed, err: [%s]", err.Error())
		}
	}

	// 删除proof字段
	withoutProof := jsonparser.Delete(vc.rawData, "proof")
	//去掉空格换行等
	withoutProof, err = CompactJson(withoutProof)
	if err != nil {
		return false, err
	}

	// 验签
	return vc.Proof.Verify(withoutProof, pkPem)
}

func getTxTime() (int64, error) {
	timestamp, err := sdk.Instance.GetTxTimeStamp()
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(timestamp, 10, 64)
}

// VcTemplate VC模板的结构内容定义
type VcTemplate struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Template string `json:"template"`
	Version  string `json:"version"`
}

func (vc *VerifiableCredential) verifyCredentialSubject(vcTemplate []byte) (bool, error) {

	subjectBytes, err := json.Marshal(vc.CredentialSubject)
	if err != nil {
		return false, err
	}

	var template VcTemplate
	err = json.Unmarshal(vcTemplate, &template)
	if err != nil {
		return false, fmt.Errorf("invalid VC template: [%s]", err.Error())
	}

	if template.Name != vc.Template.Name {
		return false, errors.New("invalid VC template name")
	}

	schemaLoader := jsonschema.NewBytesLoader(vcTemplate)
	subjectLoader := jsonschema.NewBytesLoader(subjectBytes)

	result, err := jsonschema.Validate(schemaLoader, subjectLoader)
	if err != nil {
		return false, err
	}

	if result.Valid() {
		return true, nil
	}

	errMsg := "invalid credentialSubject, errors:"
	for _, desc := range result.Errors() {
		errMsg += fmt.Sprintf("- %s\n", desc)
	}

	return false, fmt.Errorf(errMsg)
}
