# ChainMaker DID SDK

## 密钥相关

### GenerateKey

**功能**：生成公私钥

**参数说明**

- algorithm：公钥算法名称，主要支持`SM2`, `EC_Secp256k1`, `EC_NISTP224`, `EC_NISTP256`, `EC_NISTP384`, `EC_NISTP521`, `RSA512`, `RSA1024`, `RSA2048`, `RSA3072` 算法

**返回值说明**

- KeyInfo：密钥信息

```go
type KeyInfo struct {
	// 公钥的PEM编码
	PkPEM []byte
	// 私钥的PEM编码
	SkPEM []byte

	// 公钥算法名称
	Algorithm string
}
```

```go
func GenerateKey(algorithm string) (*KeyInfo, error)
```

### IsSupportAlgorithm

**功能**：判断公钥算法SDK是否支持

**参数说明**

- algo：公钥算法名称

**返回值说明**

- bool：SDK是否支持该算法

```go
func IsSupportAlgorithm(algo string) bool
```



## 证明相关

### GenerateProofByKey

**功能**：通过私钥生成证明

**参数说明**

- skPem：私钥的PEM编码
- msg：签名的信息
- verificationMethod：did中的验证方法，通常是`[DID]#key-[i]`格式
- algorithm：公钥算法（如果为空，可自行解析）
- hash：信息做摘要的哈希类型

**返回值说明**

- Proof：证明结构（引自DID合约）

```go
func GenerateProofByKey(skPem, msg []byte, verificationMethod, algorithm, hash string) (*model.Proof, error)
```

### VerifyPKProof

**功能**：通过公钥验证证明

**参数说明**

- msg：证明的信息
- pkPem：公钥的PEM编码格式
- proof： 证明的结构（引自合约）

```go
func VerifyPKProof(msg, pkPem []byte, proof *model.Proof) (bool, error)
```



## DID相关

### GetDidMethodFromChain

**功能**：在区块链上获取DID Method

**参数说明**

- client：长安链客户端

```go
func GetDidMethodFromChain(client *cmsdk.ChainClient) (string, error)
```

### GenerateDidByPK

**功能**：根据公钥生成DID

**参数说明**

- pkPem：公钥PEM编码
- client：长安链客户端

```go
func GenerateDidByPK(pkPem []byte, client *cmsdk.ChainClient) (string, error)
```

### GenerateDidDoc

**功能**：生成DID文档

**参数说明**

- keyInfo：密钥信息
- client：长安链客户端
- controller：父控制器，可变参数

```go
func GenerateDidDoc(keyInfo []*key.KeyInfo, client *cmsdk.ChainClient, controller ...string) ([]byte, error)
```

### AddDidDocToChain

**功能**：DID文档上链

**参数说明**

- doc：DID文档
- client：长安链客户端

```go
func AddDidDocToChain(doc string, client *cmsdk.ChainClient) error
```

### IsValidDidOnChain

**功能**：DID在链上是否有效

**参数说明**

- did：DID
- client：长安链客户端

```go
func IsValidDidOnChain(did string, client *cmsdk.ChainClient) (bool, error)
```

### GetDidDocFromChain

**功能**：通过DID在链上获取DID文档

**参数说明**

- did：DID
- client：长安链客户端

```go
GetDidDocFromChain(did string, client *cmsdk.ChainClient) ([]byte, error)
```

### GetDidByPkFromChain

**功能**：通过公钥在链上获取DID

**参数说明**

- pkPem：公钥的PEM编码
- client：长安链客户端

```go
func GetDidByPkFromChain(pkPem string, client *cmsdk.ChainClient) (string, error)
```

### GetDidByAddressFromChain

**功能**：通过地址在链上获取DID

**参数说明**

- address：地址
- client：长安链客户端

```go
func GetDidByAddressFromChain(address string, client *cmsdk.ChainClient) (string, error)
```

### UpdateDidDocToChain

**功能**：在链上更新DID文档

**参数说明**

- doc：DID文档
- client：长安链客户端

```go
func UpdateDidDocToChain(doc string, client *cmsdk.ChainClient) error
```

### UpdateDidDoc

**功能**：更新DID文档（本地生成）

**参数说明**

- oldDoc：原来的DID文档
- keyInfo：文档里的密钥信息
- controller：父控制器，可变参数

```go
func UpdateDidDoc(oldDoc model.DidDocument, keyInfo []*key.KeyInfo, controller ...string) ([]byte, error)
```



## DID黑名单相关

### AddDidBlackListToChain

**功能**：在链上添加DID黑名单

**参数说明**

- dids：did列表
- client：长安链客户端

```go
func AddDidBlackListToChain(dids []string, client *cmsdk.ChainClient) error
```

### GetDidBlackListFromChain

**功能**：从链上获取DID黑名单

**参数说明**

- didSearch：要查找的DID编号（空字符串可以查找全部列表）
- start：开始的索引，0表示从第一个开始
- count：要获取的数量，0表示获取所有
- client：长安链客户端

```go
func GetDidBlackListFromChain(didSearch string, start int, count int, client *cmsdk.ChainClient) ([]string, error)
```

### DeleteDidBlackListFromChain

**功能**：在链上删除DID黑名单

**参数说明**

- dids：要删除的did列表
- client：长安链客户端

```go
func DeleteDidBlackListFromChain(dids []string, client *cmsdk.ChainClient) error
```



## 权威签发者相关

### AddTrustIssuerListToChain

**功能**：在链上添加权威颁发者

**参数说明**

- dids：权威颁发者DID列表
- client：长安链客户端

```go
func AddTrustIssuerListToChain(dids []string, client *cmsdk.ChainClient) error
```

### GetTrustIssuerListFromChain

**功能**：从链上获取权威签发者列表

**参数说明**

- didSearch：要查找的DID编号（空字符串可以查找全部列表）
- start：开始的索引，0表示从第一个开始
- count：要获取的数量，0表示获取所有
- client：长安链客户端

```go
func GetTrustIssuerListFromChain(didSearch string, start int, count int, client *cmsdk.ChainClient) ([]string, error)
```

### DeleteTrustIssuerListFromChain

**功能**：从链上删除权威签发者列表

**参数说明**

- dids：要删除的did列表
- client：长安链客户端

```go
func DeleteTrustIssuerListFromChain(dids []string, client *cmsdk.ChainClient) error
```



## VC相关

### IssueVC

**功能**：颁发VC（需要链上校验）

**参数说明**

- keyInfo：颁发者的密钥信息
- keyIndex：公钥在DID文档中的索引
- subject：颁发信息主体，对应VC中的`credentialSubject`字段
- client：长安链客户端
- vcId：VC的`id`字段，可以根据业务自定义
- expirationDate：VC的到期时间
- vcTemplate：VC的模板内容，是一个JSON schema，一般存储在链上
- vcType：VC中的`type`字段，描述VC的类型信息（可变参数，默认会填写“VerifiableCredential”,可继续根据业务类型追加）

```go
func IssueVC(keyInfo *key.KeyInfo, keyIndex int, subject map[string]interface{}, client *cmsdk.ChainClient, vcId string, expirationDate int64, vcTemplateId string, vcType ...string) ([]byte, error)
```

### IssueVCLocal

**功能**：本地颁发VC（不经过链上计算和校验）

**参数说明**

- skPem：私钥的PEM编码
- algorithm：公钥算法名称
- keyIndex：公钥在DID文档中的索引
- subject：颁发信息主体，对应VC中的`credentialSubject`字段
- issuer：颁发者的DID编号
- vcId：VC的`id`字段，可以根据业务自定义
- expirationDate：VC的到期时间
- vcTemplate：VC的模板内容，是一个JSON schema，一般存储在链上
- vcType：VC中的`type`字段，描述VC的类型信息（可变参数，默认会填写“VerifiableCredential”,可继续根据业务类型追加）

```go
func IssueVCLocal(skPem []byte, algorithm string, keyIndex int, subject map[string]interface{}, issuer string, vcId string, expirationDate int64, vcTemplate []byte, vcType ...string) ([]byte, error)
```

### VerifyVCOnChain

**功能**：链上验证VC的有效性

**参数说明**

- vc：vc的JSON字符串
- client：长安链客户端

```go
func VerifyVCOnChain(vc string, client *cmsdk.ChainClient) (bool, error)
```

### RevokeVCOnChain

**功能**：在链上吊销VC

**参数说明**

- vcId：要撤销的VC编号
- client：长安链客户端

```go
func RevokeVCOnChain(vcId string,client *cmsdk.ChainClient) error
```

### GetVCRevokedListFromChain

**功能**：从链上获取VC吊销列表

**参数说明**

- vcIdSearch：要查找的vc编号（空字符串可以查找全部列表）
- start：开始的索引，0表示从第一个开始
- count：要获取的数量，0表示获取所有
- client：长安链客户端

```go
func GetVCRevokedListFromChain(vcIdSearch string, start int, count int, client *cmsdk.ChainClient) ([]string, error)
```

### GenerateSimpleVcTemplate

**功能**：生成字段都是String类型的简易的VC模板

**参数说明**

- fieldsMap：key: 字段名 value: 具体含义

```go
func GenerateSimpleVcTemplate(fieldsMap map[string]string) ([]byte, error)
```

### AddVcTemplateToChain

**功能**：VC模板上链

**参数说明**

- id：模板ID
- name：模板名称
- version：模板版本
- template：模板内容，需要JSON schema格式
- client：长安链客户端

```go
func AddVcTemplateToChain(id string, name string, version string, template []byte, client *cmsdk.ChainClient) error
```

### GetVcTemplateFromChain

**功能**：从链上获取VC模板

**参数说明**

- id：模板ID
- client：长安链客户端

```go
func GetVcTemplateFromChain(id string, client *cmsdk.ChainClient) ([]byte, error)
```

### GetVcTemplateListFromChain

**功能**：从链上获取VC模板列表

**参数说明**

- nameSearch：模板名称关键字（空字符串可以查找全部列表）
- start：开始的索引，0表示从第一个开始
- count：要获取的数量，0表示获取所有
- client：长安链客户端

```go
func GetVcTemplateListFromChain(nameSearch string, start int, count int, client *cmsdk.ChainClient) ([]string, error)
```



## VP相关

### GenerateVP

**功能**：生成VP

**参数说明**

- skPem：私钥的PEM编码
- algorithm: 公钥算法名称
- keyIndex：公钥在DID文档中的索引
- vpId：VP的`id`字段，可以根据业务自定义
- vcList：VP中包含的VC列表
- VP中的`type`字段，描述VP的类型信息（可变参数，默认会填写“VerifiablePresentation”,可继续根据业务类型追加）

```go
func GenerateVP(skPem []byte, algorithm string, keyIndex int, holder string, vpId string, vcList []string, vpType ...string) ([]byte, error)
```

### VerifyVPOnChain

**功能**：在链上验证VP的有效性

**参数说明**

- vp：vp的JSON字符串
- client：长安链客户端

```go
func VerifyVPOnChain(vp string, client *cmsdk.ChainClient) (bool, error)
```

