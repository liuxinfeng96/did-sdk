package vp

import (
	"did-sdk/proof"
	"did-sdk/utils"
	"did-sdk/vc"
	"encoding/json"
)

var ContextVP = []string{
	"https://www.w3.org/2018/credentials/v1",
	"https://www.w3.org/2018/credentials/examples/v1",
}

// VerifiablePresentation VP
type VerifiablePresentation struct {
	Context              []string                   `json:"context"`
	Id                   string                     `json:"id"`
	Type                 []string                   `json:"type"`
	VerifiableCredential []*vc.VerifiableCredential `json:"verifiableCredential"`
	Holder               string                     `json:"holder"`
	Proof                *proof.PkProofJSON         `json:"proof"`
}

// GenerateVP 生成自己的VP
// @params skPem: 私钥的PEM编码
// @params algorithm: 公钥算法名称
// @params vpId：VP的`id`字段，可以根据业务自定义
// @params VP中包含的VC列表
// @params vpType：VP中的`type`字段，描述VP的类型信息（可变参数，默认会填写“VerifiablePresentation”,可继续根据业务类型追加）
func GenerateVP(skPem []byte, algorithm string, holder string,
	vpId string, vcList []string, vpType ...string) ([]byte, error) {

	var verifiablePresentation VerifiablePresentation

	vpType = append(vpType, "VerifiablePresentation")
	for _, v := range vcList {
		var verifiableCredential vc.VerifiableCredential
		err := json.Unmarshal([]byte(v), &verifiableCredential)
		if err != nil {
			return nil, err
		}

		verifiablePresentation.VerifiableCredential =
			append(verifiablePresentation.VerifiableCredential, &verifiableCredential)
	}

	verifiablePresentation.Context = ContextVP
	verifiablePresentation.Id = vpId
	verifiablePresentation.Type = vpType
	verifiablePresentation.Holder = holder

	vpBytes, err := json.Marshal(verifiablePresentation)
	if err != nil {
		return nil, err
	}

	keyId := holder + "#keys-1"

	pf, err := proof.GenerateProofByKey(skPem, vpBytes, keyId,
		algorithm, utils.GetHashTypeByAlgorithm(algorithm))
	if err != nil {
		return nil, err
	}

	verifiablePresentation.Proof = pf

	return json.Marshal(verifiablePresentation)

}
