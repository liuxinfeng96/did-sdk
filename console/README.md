# Console

## key

### 生成公私钥

```shell
$ ./didc key gen \
--algo=SM2 \
--pk-path=./testdata/pk.pem \
--sk-path=./testdata/sk.pem
```



## did

### 获取DID方法

```shell
$ ./didc did method \
--sdk-path=./testdata/sdk.yaml
```



### 根据公钥生成DID

```shell
$ ./didc did gen \
--pk-path=./testdata/pk.pem \
--sdk-path=./testdata/sdk.yaml
```



### DID在链上是否有效

```shell
$ ./didc did valid \
--did=did:cm:test1 \
--sdk-psath=./testdata/sdk.yaml 
```



### 获取DID

```shell
$ ./didc did get \
--pk-path=./testdata/pk.pem \
--address=dkalwdkladawldakdwa \
--sdk-path=./testdata/sdk.yaml
```



## doc

### 生成DID文档

```shell
$ ./didc doc gen \
--sk-path=./testdata/sk.pem \
--pk-path=./testdata/pk.pem \
--algo=SM2 \
--controller=did:cm:test1,did:cm:test2 \
--sdk-path=./testdata/sdk.yaml \
--doc-path=./testdata/doc.json
```



### DID文档上链

```shell
$ ./didc doc add \
--doc-path=./testdata/doc.json \
--sdk-path=./testdata/sdk.yaml 
```



### 获取DID文档

```shell
$ ./didc doc get \
--did=did:cm:test1 \
--sdk-path=./testdata/sdk.yaml \
--doc-path=./testdata/doc.json
```



### 更新DID文档

```shell
$ ./didc doc update-local  \
--sk-path=./testdata/sk.pem \
--pk-path=./testdata/pk.pem \
--algo=SM2 \
--controller=did:cm:test1,did:cm:test2 \
--old-doc-path=./testdata/doc.json \
--new-doc-path=./testdata/doc2.json
```



### 链上更新DID文档

```shell
$ ./didc doc update \
--doc-path=./testdata/doc.json \
--sdk-path=./testdata/sdk.yaml 
```



## black

### 黑名单上链

```shell
$ ./didc black add \
--dids=did:cm:test1,did:cm:test2 \
--sdk-path=./testdata/sdk.yaml 
```



### 链上获取黑名单

```shell
$ ./didc black list \
--search=did:cm:test1 \
--start=1 \
--count=10\
--sdk-path=./testdata/sdk.yaml
```



### 链上删除黑名单

```shell
$ ./didc black delete \
--dids=did:cm:test1,did:cm:test2 \
--sdk-path=./testdata/sdk.yaml 
```



## Issuer

### 链上添加权威签发者

```shell
$ ./didc issuer add \
--dids=did:cm:test1,did:cm:test2 \
--sdk-path=./testdata/sdk.yaml 
```



### 链上获取权威签发者列表

```shell
$ ./didc issuer list \
--search=did:cm:test1 \
--start=1 \
--count=10 \
--sdk-path=./testdata/sdk.yaml
```



### 链上删除权威签发者

```shell
$ ./didc issuer delete \
--dids=did:cm:test1,did:cm:test2 \
--sdk-path=./testdata/sdk.yaml 
```



## vc

### 链上颁发VC

```shell
$ ./didc vc issue \
--sk-path=./testdata/sk.pem \
--pk-path=./testdata/pk.pem \
--algo=SM2 \
--key-index=1 \
--subject=./testdata/subject.json \
--expiration=2025-01-25 \
--id=111233 \
--template=./testdata/template.json \
--type=Identity \
--vc-path=./testdata/vc.json \
--sdk-path=./testdata/sdk.yaml
```



### 本地颁发VC

```shell
$ ./didc vc issue-local \
--sk-path=./testdata/sk.pem \
--algo=SM2 \
--key-index=1 \
--subject=./testdata/subject.json \
--issuer=did:cm:admin \
--expiration=2025-01-25 \
--id=111233 \
--template=./testdata/template.json \
--type=Identity \
--vc-path=./testdata/vc.json 
```



### 链上验证VC的有效性

```shell
$ ./didc vc verify \
--vc-path=./testdata/vc.json \
--sdk-path=./testdata/sdk.yaml
```



## vc-revoke

### 链上吊销VC

```shell
$ ./didv vc-revoke add \
--id=16516616 \
--sdk-path=./testdata/sdk.yaml
```



### 链上获取吊销列表

```shell
$ ./didv vc-revoke list \
--search=111515515 \
--start=1 \
--count=10 \
--sdk-path=./testdata/sdk.yaml
```



## vc-template

### VC模板上链

```shell
$ ./didv vc-template add \
--id=1515 \
--name=模板1 \
--version=v1.0.0 \
--template=./testdata/template.json \
--sdk-path=./testdata/sdk.yaml
```



 ### 获取VC模板

```shell
$ ./didv vc-template get \
--id=151515 \
--sdk-path=./testdata/sdk.yaml
```



### 获取VC模板列表

```shell
$ ./didv vc-template list \
--search=模板1 \
--start=1 \
--count=10 \
--sdk-path=./testdata/sdk.yaml
```



## vp

### 生成VP

```shell
$ ./didv vp gen \
--sk-path=./testdata/sk.pem \
--algo=SM2 \
--key-index=1 \
--id=1231232 \
--vcs=./testdata/vc1.json,./testdata/vc2.json \
--type=Identity
```



### 链上验证VP

```shell
$ ./didv vp verify \
--vp-path=./testdata/vp.json \
--sdk-path=./testdata/sdk.yaml
```

