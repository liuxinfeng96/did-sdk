package did

import (
	"crypto/sha256"
	"did-sdk/contract"
	"did-sdk/key"
	"did-sdk/proof"
	"did-sdk/utils"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"strconv"
	"time"

	"chainmaker.org/chainmaker/pb-go/v2/common"
	cmsdk "chainmaker.org/chainmaker/sdk-go/v2"
	"github.com/mr-tron/base58"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"

	bcx509 "github.com/liuxinfeng96/bc-crypto/x509"
)

const (
	DidPrefix  = "did:"
	DidContext = "https://www.w3.org/ns/did/v1"
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

// GenerateDidByPK did was generated using the public key PEM encoding
// @params pkPem: PK PEM
// @params client: ChainMaker SDK
// @return string: the did string
func GenerateDidByPK(pkPem []byte, client *cmsdk.ChainClient) (string, error) {
	// 从链上获取DID方法名
	didMethod, err := GetDidMethodFromChain(client)
	if err != nil {
		return "", err
	}

	pkHash := sha256Hash(pkPem)
	didSuffix := base58Encode(pkHash)

	did := fmt.Sprintf("%s:%s:%s", DidPrefix, didMethod, didSuffix)

	return did, nil
}

// GenerateDidDoc generate a DID document using a key
// @params keyInfo：密钥信息
// @params client：长安链客户端
// @params controller：父控制器，可变参数
func GenerateDidDoc(keyInfo []*key.KeyInfo, client *cmsdk.ChainClient, controller ...string) ([]byte, error) {

	// 密钥最少一把
	if len(keyInfo) == 0 {
		return nil, errors.New("key info cannot be empty")
	}

	// 通过公钥生成DID字符串
	did, err := GenerateDidByPK(keyInfo[0].PkPEM, client)
	if err != nil {
		return nil, err
	}

	verificationMethod := make([]*VerificationMethod, 0)
	authentication := make([]string, 0)
	controller = append(controller, did)

	for k, v := range keyInfo {
		keyId := did + "#keys-" + strconv.Itoa(k)

		vm, err := newVerificationMethod(keyId, v.Algorithm, did, v.PkPEM)
		if err != nil {
			return nil, err
		}

		verificationMethod = append(verificationMethod, vm)

		authentication = append(authentication, keyId)

	}

	created := utils.ISO8601Time(time.Now().Unix())

	doc := &DidDocument{
		Context:            DidContext,
		Id:                 did,
		Created:            created,
		Updated:            created,
		VerificationMethod: verificationMethod,
		Authentication:     authentication,
		Controller:         controller,
	}

	docByte, err := json.Marshal(doc)
	if err != nil {
		return nil, err
	}

	proofs := make([]*proof.PkProofJSON, 0)

	for k, v := range keyInfo {
		keyId := did + "#keys-" + strconv.Itoa(k)
		var hash string
		if v.Algorithm == "SM2" {
			hash = "SM3"
		} else {
			hash = "SHA-256"
		}
		pf, err := proof.GenerateProofByKey(v.SkPEM, docByte, keyId, v.Algorithm, hash)
		if err != nil {
			return nil, err
		}
		proofs = append(proofs, pf)
	}

	doc.Proof = proofs

	return json.Marshal(doc)
}

// AddDidDocToChain store the DID document on the blockchain
// @params doc：DID文档
// @params client：长安链客户端
func AddDidDocToChain(doc string, client *cmsdk.ChainClient) error {
	params := make([]*common.KeyValuePair, 0)

	params = append(params, &common.KeyValuePair{
		Key:   "did",
		Value: []byte(doc),
	})

	_, err := contract.InvokeContract(contract.Contract_Did, contract.Method_AddDidDocument, params, client)
	if err != nil {
		return err
	}

	return nil
}

func sha256Hash(str []byte) []byte {
	hash := sha256.Sum256(str)
	return hash[:]
}

func base58Encode(hash []byte) string {
	encoded := base58.Encode(hash)
	return encoded
}

func newVerificationMethod(id, algo, controller string, pkPem []byte) (*VerificationMethod, error) {

	var pkDer []byte

	block, rest := pem.Decode(pkPem)
	if block == nil {
		pkDer = rest
	} else {
		pkDer = block.Bytes
	}

	// 校验是否是公钥
	_, err := bcx509.ParsePublicKeyFromDER(pkDer)
	if err != nil {
		return nil, err
	}

	// 计算出Address
	bytesAddr := ethcrypto.Keccak256(pkDer)
	addr := hex.EncodeToString(bytesAddr)[24:]

	return &VerificationMethod{
		Id:           id,
		Type:         algo,
		Controller:   controller,
		PublicKeyPem: string(pkPem),
		Address:      addr,
	}, nil

}
