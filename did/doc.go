package did

import (
	"did-sdk/contract"
	"did-sdk/proof"
	"fmt"

	cmsdk "chainmaker.org/chainmaker/sdk-go/v2"
)

// DidDocument the JSON structure of the DID document
type DidDocument struct {
	Context            string                `json:"@context"`
	Id                 string                `json:"id"`
	Created            string                `json:"created"`
	Updated            string                `json:"updated"`
	VerificationMethod []*VerificationMethod `json:"verificationMethod"`
	Authentication     []string              `json:"authentication"`
	Controller         []string              `json:"controller"`
	Proof              []*proof.PkProofJSON  `json:"proof"`
}

// VerificationMethod the JSON structure of the DID document VerificationMethod
type VerificationMethod struct {
	Id           string `json:"id"`
	Type         string `json:"type"`
	Controller   string `json:"controller"`
	PublicKeyPem string `json:"publicKeyPem"`
	Address      string `json:"address"`
}

// GetDidMethodFromChain query contract from chain
// @params client the chainmaker sdk client
// @return string the did method
func GetDidMethodFromChain(client *cmsdk.ChainClient) (string, error) {

	resp, err := client.QueryContract(contract.Contract_Did, contract.Method_DidMethod, nil, -1)
	if err != nil {
		return "", fmt.Errorf("send tx failed, err: [%s]", err.Error())
	}

	result, err := contract.DealTxResponse(resp, contract.Contract_Did, contract.Method_DidMethod)
	if err != nil {
		return "", err
	}

	return string(result), nil
}
