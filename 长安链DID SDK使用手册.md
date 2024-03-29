# 长安链DID SDK使用手册

## 介绍

​	2019年，W3C提出了去中心化身份标识(Decentralized Identifier，DID)的标准，虽然W3C没有明确规范DID的实现需要使用区块链技术，但是很显然，与区块链技术结合，可以更好的实现DID。

​	为了方便长安链的用户能够快速的搭建自己的去中心化数字身份应用，我们基于长安链社区的去中心化数字身份（DID）合约标准`CM-CS-231201-DID`，实现了一套比较完整的DID智能合约，并且提供SDK开发工具包。

## 快速体验

使用DID**命令行工具**，快速体验DID应用主要功能流程，了解`DID SDK`功能。

### 搭建长安链

[快速搭建长安链]()

### 安装DID智能合约

代码下载地址：`https://github.com/liuxinfeng96/did-sdk.git`

拉取代码：

```
git clone https://github.com/liuxinfeng96/did-sdk.git
```

编译智能合约：

```shell
$ cd did-sdk && cd contract
$ ./build.sh ChainMakerDid
```

**CMC工具安装合约**

进入长安链项目主目录：

```shell
$ cd chainmaker-go
```

进入`cmc`文件夹：

```shell
$ cd tools/cmc/
```

编译`cmc`:

```shell
$ go build
```

执行`cmc`安装合约命令：

```shell
$ ./cmc client contract user create \
--contract-name=ChainMakerDid \
--runtime-type=DOCKER_GO \
--byte-code-path=../../../did-sdk/contract/ChainMakerDid.7z \
--version=1.0 \
--sdk-conf-path=./testdata/sdk_config.yml \
--admin-key-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.tls.key,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.tls.key,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.tls.key,./testdata/crypto-config/wx-org4.chainmaker.org/user/admin1/admin1.tls.key \
--admin-crt-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.tls.crt,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.tls.crt,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.tls.crt,./testdata/crypto-config/wx-org4.chainmaker.org/user/admin1/admin1.tls.crt \
--sync-result=true \
--params="{\"didMethod\":\"cm\",\"enableTrustIssuer\":\"true\"}"
```

**测试脚本安装合约**

拷贝长安链证书密钥文件：

```shell
$ cd chainmaker-go/build && cp -r ./crypto-config/ ../../did-sdk/testdata/
```

执行测试脚本：

```shell
$ cd ../../did-sdk/testdata/
$ go test -v -run TestInstallDidContract
```

### 主要功能

拷贝长安链证书密钥文件至**命令行工具**文件夹：

```shell
$ cp -r ./testdata/crypto-config/ ./console/testdata/
```

编译控制台：

```shell
$ cd console && go build
```

申请密钥：

```shell
$ ./console key gen \
--algo=SM2 \
--pk-path=./testdata/pk.pem \
--sk-path=./testdata/sk.pem
```

生成DID文档：

```shell
$ ./console doc gen \
--sks-path=./testdata/sk.pem \
--pks-path=./testdata/pk.pem \
--controller=did:cm:test1,did:cm:test2 \
--sdk-path=./testdata/sdk_config.yml \
--doc-path=./testdata/doc.json
```

DID文档`doc.json`内容：

```json
{
    "@context": "https://www.w3.org/ns/did/v1",
    "id": "did:cm:AmGZacodAgjdnoEDpH9v2VToz1UMhi5FD6SkVL2fYUo3",
    "created": "2024-03-08T16:38:21+08:00",
    "updated": "2024-03-08T16:38:21+08:00",
    "verificationMethod": [
        {
            "id": "did:cm:AmGZacodAgjdnoEDpH9v2VToz1UMhi5FD6SkVL2fYUo3#keys-0",
            "type": "SM2-SM3",
            "controller": "did:cm:AmGZacodAgjdnoEDpH9v2VToz1UMhi5FD6SkVL2fYUo3",
            "publicKeyPem": "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoEcz1UBgi0DQgAEIFvBi8CbbD3fVwfn4QSIckx12mpr\nQGFJCbcAqsLkeQ7pJKOxF4vOOeE9JFvt8qXnqRPgQHbjfnyUcQCdi6SmJg==\n-----END PUBLIC KEY-----\n",
            "address": "8c6e1121ac269037b69d69ca5bdf5db812ebb55a"
        }
    ],
    "authentication": [
        "did:cm:AmGZacodAgjdnoEDpH9v2VToz1UMhi5FD6SkVL2fYUo3#keys-0"
    ],
    "controller": [
        "did:cm:test1",
        "did:cm:test2",
        "did:cm:AmGZacodAgjdnoEDpH9v2VToz1UMhi5FD6SkVL2fYUo3"
    ],
    "proof": {
        "type": "SM2-SM3",
        "created": "2024-03-08T16:38:21+08:00",
        "proofPurpose": "assertionMethod",
        "verificationMethod": "did:cm:AmGZacodAgjdnoEDpH9v2VToz1UMhi5FD6SkVL2fYUo3#keys-0",
        "proofValue": "MEUCIQC7mxQgOFCknMqrAv6Ix5Jaftm3POrdXusi0F5qtDqo3gIgfdO1yjfA2WXqHqPN6eGSbORJwS9T14aV8tVxHw08z9I="
    }
}
```

DID文档上链：

```shell
$ ./console doc add \
--doc-path=./testdata/doc.json \
--sdk-path=./testdata/sdk_config.yml 
```

根据公钥生成DID：

```shell
$ ./console did gen \
--pk-path=./testdata/pk.pem \
--sdk-path=./testdata/sdk_config.yml
```

返回DID字符串：

```shell
did: [did:cm:AmGZacodAgjdnoEDpH9v2VToz1UMhi5FD6SkVL2fYUo3]
```

从链上获取DID文档，将上一步返回的did填入`--did`参数：

```shell
$ ./console doc get \
--did=did:cm:AmGZacodAgjdnoEDpH9v2VToz1UMhi5FD6SkVL2fYUo3 \
--sdk-path=./testdata/sdk_config.yml \
--doc-path=./testdata/doc.json
```

链上注册签发者：

```shell
$ ./console issuer add \
--dids=did:cm:AmGZacodAgjdnoEDpH9v2VToz1UMhi5FD6SkVL2fYUo3 \
--sdk-path=./testdata/sdk_config.yml
```

查询签发者列表：

```shell
$ ./console issuer list \
--start=1 \
--count=10 \
--sdk-path=./testdata/sdk_config.yml
```

返回信任签发者列表：

```shell
get the did list of issuer: [[did:cm:AmGZacodAgjdnoEDpH9v2VToz1UMhi5FD6SkVL2fYUo3]]
```

被签发者申请密钥：

```shell
$ ./console key gen \
--algo=EC_Secp256k1 \
--pk-path=./testdata/pk2.pem \
--sk-path=./testdata/sk2.pem
```

生成DID文档：

```shell
$ ./console doc gen \
--sks-path=./testdata/sk2.pem \
--pks-path=./testdata/pk2.pem \
--sdk-path=./testdata/sdk_config.yml \
--doc-path=./testdata/doc2.json
```

DID文档`doc2.json`内容：

```json
{
    "@context": "https://www.w3.org/ns/did/v1",
    "id": "did:cm:Hyi8uUrBPhsH2i1sd1zaZGYeQA9WdxGgiAu8YPgY8pxe",
    "created": "2024-03-08T16:42:30+08:00",
    "updated": "2024-03-08T16:42:30+08:00",
    "verificationMethod": [
        {
            "id": "did:cm:Hyi8uUrBPhsH2i1sd1zaZGYeQA9WdxGgiAu8YPgY8pxe#keys-0",
            "type": "ECDSA-SHA256",
            "controller": "did:cm:Hyi8uUrBPhsH2i1sd1zaZGYeQA9WdxGgiAu8YPgY8pxe",
            "publicKeyPem": "-----BEGIN PUBLIC KEY-----\nMFYwEAYHKoZIzj0CAQYFK4EEAAoDQgAEk8LPsj8/8OF1/P5E885/taW9Z5briwXI\n+adAsK8bVluJGFZSK/TOoLklAhTgxGj1TBKRh0TSayUTBwC+oH0gdw==\n-----END PUBLIC KEY-----\n",
            "address": "2fc542c747a1a26b977c5b1fb824e070cb41b181"
        }
    ],
    "authentication": [
        "did:cm:Hyi8uUrBPhsH2i1sd1zaZGYeQA9WdxGgiAu8YPgY8pxe#keys-0"
    ],
    "controller": [
        "did:cm:Hyi8uUrBPhsH2i1sd1zaZGYeQA9WdxGgiAu8YPgY8pxe"
    ],
    "proof": {
        "type": "ECDSA-SHA256",
        "created": "2024-03-08T16:42:30+08:00",
        "proofPurpose": "assertionMethod",
        "verificationMethod": "did:cm:Hyi8uUrBPhsH2i1sd1zaZGYeQA9WdxGgiAu8YPgY8pxe#keys-0",
        "proofValue": "MEYCIQDU/qVbPIT6yumUuIH4L2Fvxu/ZAfuGdD6Jb8DtomgIdgIhALW6lnslC0d1NW+9WEn5ffDz5DR4P5KJonnM9J7nclg8"
    }
}
```

DID文档上链：

```shell
$ ./console doc add \
--doc-path=./testdata/doc2.json \
--sdk-path=./testdata/sdk_config.yml
```

根据公钥生成DID：

```shell
$ ./console did gen \
--pk-path=./testdata/pk2.pem \
--sdk-path=./testdata/sdk_config.yml
```

返回DID字符串：

```shell
did: [did:cm:Hyi8uUrBPhsH2i1sd1zaZGYeQA9WdxGgiAu8YPgY8pxe]
```

生成`VC`模板：

```shell
$ ./console vc-template gen \
--map-key=name,age,sex \
--map-value=姓名,年龄,性别 \
--temp-path=./testdata/temp.json
```

`VC`模板`temp.json`内容：

```json
{
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "properties": {
        "age": {
            "type": "string",
            "title": "年龄"
        },
        "id": {
            "type": "string",
            "title": "DID"
        },
        "name": {
            "type": "string",
            "title": "姓名"
        },
        "sex": {
            "type": "string",
            "title": "性别"
        }
    },
    "required": [
        "name",
        "age",
        "sex",
        "id"
    ],
    "additionalProperties": false
}
```

`VC`模板上链：

```shell
$ ./console vc-template add \
--temp-id=temp001 \
--temp-name=模板1 \
--temp-version=v1.0.0 \
--temp-path=./testdata/temp.json \
--sdk-path=./testdata/sdk_config.yml
```

新建被签发者主体`subject.json`文件：

```shell
$ vim ./testdata/subject.json
```

将下面签发的主体内容复制并保存至`subject.json`文件，**将查询的被签发者DID填入`id`字段**：

```json
{
    "id": "did:cm:Hyi8uUrBPhsH2i1sd1zaZGYeQA9WdxGgiAu8YPgY8pxe",
    "name": "xiaoming",
    "age": "18",
    "sex": "男"
}
```

签发者颁发可验证凭证`VC`：

```shell
$ ./console vc issue \
--sk-path=./testdata/sk.pem \
--pk-path=./testdata/pk.pem \
--subject=./testdata/subject.json \
--expiration=2025-01-25 \
--id=vc001 \
--temp-id=temp001 \
--type=Identity \
--vc-path=./testdata/vc.json \
--sdk-path=./testdata/sdk_config.yml
```

可验证凭证 `vc.json`内容：

```json
{
    "@context": [
        "https://www.w3.org/2018/credentials/v1",
        "https://www.w3.org/2018/credentials/examples/v1"
    ],
    "id": "vc001",
    "type": [
        "Identity",
        "VerifiableCredential"
    ],
    "credentialSubject": {
        "age": "18",
        "id": "did:cm:Hyi8uUrBPhsH2i1sd1zaZGYeQA9WdxGgiAu8YPgY8pxe",
        "name": "xiaoming",
        "sex": "男"
    },
    "issuer": "did:cm:AmGZacodAgjdnoEDpH9v2VToz1UMhi5FD6SkVL2fYUo3",
    "issuanceDate": "2024-03-08T16:49:03+08:00",
    "expirationDate": "2025-01-25T00:00:00+08:00",
    "template": {
        "id": "temp001",
        "name": "模板1"
    },
    "proof": {
        "type": "SM2-SM3",
        "created": "2024-03-08T16:49:03+08:00",
        "proofPurpose": "assertionMethod",
        "verificationMethod": "did:cm:AmGZacodAgjdnoEDpH9v2VToz1UMhi5FD6SkVL2fYUo3#keys-0",
        "proofValue": "MEYCIQCh7l5EJckKjZ+yMbOabCBzfnfU47EiHE4cY2HJSYWmOQIhANmGwmvHas1d/jGa20ain5s/s7yV37BfPRkKIljtvAWk"
    }
}
```

查询`VC`签发日志：

```shell
$ ./console vc log \
--search=vc001 \
--start=1 \
--count=10 \
--sdk-path=./testdata/sdk_config.yml
```

返回签发日志内容：

```shell
&{Issuer:did:cm:AmGZacodAgjdnoEDpH9v2VToz1UMhi5FD6SkVL2fYUo3 Did:did:cm:Hyi8uUrBPhsH2i1sd1zaZGYeQA9WdxGgiAu8YPgY8pxe TemplateId:temp001 VcId:vc001 IssueTime:1709887743}
```

链上`VC`验证：

```shell
$ ./console vc verify \
--vc-path=./testdata/vc.json \
--sdk-path=./testdata/sdk_config.yml
```

返回校验结果：

```shell
the verification result of vc is: [true]
```

可验证表述`VP`的生成，**`--holder`参数填入被签发者DID，同`subject.json`文件里的`id`字段**：

```shell
$ ./console vp gen \
--sk-path=./testdata/sk2.pem \
--holder=did:cm:Hyi8uUrBPhsH2i1sd1zaZGYeQA9WdxGgiAu8YPgY8pxe \
--id=vp001 \
--vc-list=./testdata/vc.json \
--type=Identity \
--vp-path=./testdata/vp.json
```

可验证表述`vp.json`内容：

```json
{
    "@context": [
        "https://www.w3.org/2018/credentials/v1",
        "https://www.w3.org/2018/credentials/examples/v1"
    ],
    "id": "vp001",
    "type": [
        "Identity",
        "VerifiablePresentation"
    ],
    "verifiableCredential": [
        {
            "@context": [
                "https://www.w3.org/2018/credentials/v1",
                "https://www.w3.org/2018/credentials/examples/v1"
            ],
            "id": "vc001",
            "type": [
                "Identity",
                "VerifiableCredential"
            ],
            "credentialSubject": {
                "age": "18",
                "id": "did:cm:Hyi8uUrBPhsH2i1sd1zaZGYeQA9WdxGgiAu8YPgY8pxe",
                "name": "xiaoming",
                "sex": "男"
            },
            "issuer": "did:cm:AmGZacodAgjdnoEDpH9v2VToz1UMhi5FD6SkVL2fYUo3",
            "issuanceDate": "2024-03-08T16:49:03+08:00",
            "expirationDate": "2025-01-25T00:00:00+08:00",
            "template": {
                "id": "temp001",
                "name": "模板1"
            },
            "proof": {
                "type": "SM2-SM3",
                "created": "2024-03-08T16:49:03+08:00",
                "proofPurpose": "assertionMethod",
                "verificationMethod": "did:cm:AmGZacodAgjdnoEDpH9v2VToz1UMhi5FD6SkVL2fYUo3#keys-0",
                "proofValue": "MEYCIQCh7l5EJckKjZ+yMbOabCBzfnfU47EiHE4cY2HJSYWmOQIhANmGwmvHas1d/jGa20ain5s/s7yV37BfPRkKIljtvAWk"
            }
        }
    ],
    "holder": "did:cm:Hyi8uUrBPhsH2i1sd1zaZGYeQA9WdxGgiAu8YPgY8pxe",
    "proof": {
        "type": "ECDSA-SHA256",
        "created": "2024-03-08T16:52:31+08:00",
        "proofPurpose": "assertionMethod",
        "verificationMethod": "did:cm:Hyi8uUrBPhsH2i1sd1zaZGYeQA9WdxGgiAu8YPgY8pxe#keys-0",
        "proofValue": "MEQCIFUXSUz0TCBANreHVGfVR3+hhkqbiXSpURfErOXx1tXbAiBBR5vDKdfS2nhmKZ8otwqtFZQpcV1BjO6t5UTHJxD+xQ=="
    }
}
```

链上`VP`验证

```shell
$ ./console vp verify \
--vp-path=./testdata/vp.json \
--sdk-path=./testdata/sdk_config.yml
```

返回校验结果：

```shell
the verification result of vp is: [true]
```

### 其他功能

DID智能合约在`DID文档的更新`、`黑名单的管理`、`权威签发者的管理`和`VC的吊销`等操作需要一定的操作权限限制。

合约的`创建者（creator）`拥有合约最大权限，`creator`可以为合约设置管理员`admin`，添加合约管理员需要使用管理员的`公钥`。

操作权限具体看下表：

|                        | creator | admin | issuer | controller | other |
| :--------------------: | :-----: | :---: | :----: | :--------: | :---: |
|   管理员的设置、删除   |    Y    |   N   |   N    |     N      |   N   |
|   黑名单的添加、删除   |    Y    |   Y   |   N    |     N      |   N   |
| 权威签发者的添加、删除 |    Y    |   Y   |   N    |     N      |   N   |
|       凭证的吊销       |    Y    |   Y   |   Y    |     N      |   N   |
|       模板的添加       |    Y    |   Y   |   Y    |     N      |   N   |
|     DID文档的更新      |    Y    |   Y   |   Y    |     Y      |   N   |

**管理员的管理**

查询是否拥有管理员权限：

```shell
$ ./console admin auth \
--sdk-path=./testdata/sdk_config2.yml
```

返回管理员权限的结果：

```shell
Is admin: [false]
```

增加管理员：

```shell
$ ./console admin add \
--admin-sdk-path=./testdata/sdk_config2.yml \
--sdk-path=./testdata/sdk_config.yml
```

查询是否拥有管理员权限：

```shell
$ ./console admin auth \
--sdk-path=./testdata/sdk_config2.yml
```

返回管理员权限的结果：

```shell
Is admin: [true]
```

删除管理员：

```shell
$ ./console admin delete \
--admin-sdk-path=./testdata/sdk_config2.yml \
--sdk-path=./testdata/sdk_config.yml
```

查询是否拥有管理员权限：

```shell
$ ./console admin auth \
--sdk-path=./testdata/sdk_config2.yml
```

返回管理员权限的结果：

```shell
Is admin: [false]
```

**黑名单的管理**

查询DID在链上是否有效：

```shell
$ ./console did valid \
--did=did:cm:test1 \
--sdk-path=./testdata/sdk_config.yml
```

返回验证结果：

```shell
whether the did is valid: [true]
```

添加DID黑名单：

```shell
$ ./console black add \
--dids=did:cm:test1,did:cm:9h6JLhdJbDdPFGJrf2YaxQzj1UX2NmcWfzL65VhmvoUT \
--sdk-path=./testdata/sdk_config.yml
```

查询DID在链上是否有效：

```shell
$ ./console did valid \
--did=did:cm:test1 \
--sdk-path=./testdata/sdk_config.yml
```

返回验证结果：

```shell
whether the did is valid: [false], err: [[ChainMakerDid-IsValidDid] exec contract failed, TxId:[17baa93f7fc48952ca0379f1af81e3b4c33456a171a04e5ebfe1ecc14a45c631], TxStatusCode: [CONTRACT_FAIL], ContractCode: [1], Result: [the did in the black list]]
```

链上获取黑名单列表：

```shell
$ ./console black list \
--start=1 \
--count=10 \
--sdk-path=./testdata/sdk_config.yml
```

返回黑名单结果：

```shell
get the did black list: [[did:cm:9h6JLhdJbDdPFGJrf2YaxQzj1UX2NmcWfzL65VhmvoUT did:cm:test1]]
```

删除DID黑名单：

```shell
./console black delete \
--dids=did:cm:test1 \
--sdk-path=./testdata/sdk_config.yml
```

查询DID在链上是否有效：

```shell
$ ./console did valid \
--did=did:cm:test1 \
--sdk-path=./testdata/sdk_config.yml
```

返回验证结果：

```shell
whether the did is valid: [true]
```

**VC吊销**

吊销VC：

```shell
$ ./console vc-revoke add \
--id=vc001 \
--sdk-path=./testdata/sdk_config.yml
```

获取吊销列表：

```shell
$ ./console vc-revoke list \
--start=1 \
--count=10 \
--sdk-path=./testdata/sdk_config.yml
```

返回吊销列表：

```shell
get the vc revoke list: [[vc001]]
```



## DID SDK

[SDK接口详情](./README.md)

## 命令行工具

[命令行工具命令详情](./console/README.md)

