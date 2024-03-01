# Console

## key

### 生成公私钥

```shell
$ ./console key gen \
--algo=SM2 \
--pk-path=./testdata/pk.pem \
--sk-path=./testdata/sk.pem
```



## did

### 获取DID方法

```shell
$ ./console did method \
--sdk-path=./testdata/sdk_config.yml
```



### 根据公钥生成DID

```shell
$ ./console did gen \
--pk-path=./testdata/pk.pem \
--sdk-path=./testdata/sdk_config.yml
```



### DID在链上是否有效

```shell
$ ./console did valid \
--did=did:cm:9h6JLhdJbDdPFGJrf2YaxQzj1UX2NmcWfzL65VhmvoUT \
--sdk-path=./testdata/sdk_config.yml 
```



### 获取DID

```shell
$ ./console did get \
--pk-path=./testdata/pk.pem \
# --address=5a53a9ddf58129140734a0e5ab3ade209ae487fb \
--sdk-path=./testdata/sdk_config.yml
```



## doc

### 生成DID文档

```shell
$ ./console doc gen \
--sks-path=./testdata/sk.pem \
--pks-path=./testdata/pk.pem \
--algos=SM2 \
--controller=did:cm:test1,did:cm:test2 \
--sdk-path=./testdata/sdk_config.yml \
--doc-path=./testdata/doc.json
```



### DID文档上链

```shell
$ ./console doc add \
--doc-path=./testdata/doc.json \
--sdk-path=./testdata/sdk_config.yml 
```



### 获取DID文档

```shell
$ ./console doc get \
--did=did:cm:9h6JLhdJbDdPFGJrf2YaxQzj1UX2NmcWfzL65VhmvoUT \
--sdk-path=./testdata/sdk_config.yml \
--doc-path=./testdata/doc.json
```



### 更新DID文档

```shell
$ ./console doc update-local  \
--sks-path=./testdata/sk.pem \
--pks-path=./testdata/pk.pem \
--algos=SM2 \
--controller=did:cm:test6 \
--old-doc-path=./testdata/doc.json \
--new-doc-path=./testdata/newdoc.json
```



### 链上更新DID文档

```shell
$ ./console doc update \
--doc-path=./testdata/newdoc.json \
--sdk-path=./testdata/sdk_config.yml 
```



## black

### 黑名单上链

```shell
$ ./console black add \
--dids=did:cm:test1,did:cm:9h6JLhdJbDdPFGJrf2YaxQzj1UX2NmcWfzL65VhmvoUT \
--sdk-path=./testdata/sdk_config.yml 
```



### 链上获取黑名单

```shell
$ ./console black list \
--start=1 \
--count=10 \
--sdk-path=./testdata/sdk_config.yml
```



### 链上删除黑名单

```shell
$ ./console black delete \
--dids=did:cm:test1,did:cm:9h6JLhdJbDdPFGJrf2YaxQzj1UX2NmcWfzL65VhmvoUT \
--sdk-path=./testdata/sdk_config.yml 
```



## Issuer

### 链上添加权威签发者

```shell
$ ./console issuer add \
--dids=did:cm:9h6JLhdJbDdPFGJrf2YaxQzj1UX2NmcWfzL65VhmvoUT \
--sdk-path=./testdata/sdk_config.yml
```



### 链上获取权威签发者列表

```shell
$ ./console issuer list \
--start=1 \
--count=10 \
--sdk-path=./testdata/sdk_config.yml
```



### 链上删除权威签发者

```shell
$ ./console issuer delete \
--dids=did:cm:test1 \
--sdk-path=./testdata/sdk_config.yml 
```



## vc

### 链上颁发VC

```shell
$ ./console vc issue \
--sk-path=./testdata/sk.pem \
--pk-path=./testdata/pk.pem \
--algo=SM2 \
--subject=./testdata/subject.json \
--expiration=2025-01-25 \
--id=vc001 \
--temp-id=12313213 \
--type=Identity \
--vc-path=./testdata/vc.json \
--sdk-path=./testdata/sdk_config.yml
```



### 本地颁发VC

```shell
$ ./console vc issue-local \
--sk-path=./testdata/sk.pem \
--algo=SM2 \
--subject=./testdata/subject.json \
--issuer=did:cm:admin \
--expiration=2025-01-25 \
--id=vc001 \
--temp-path=./testdata/template.json \
--type=Identity \
--vc-path=./testdata/vc.json 
```



### 链上验证VC的有效性

```shell
$ ./console vc verify \
--vc-path=./testdata/vc.json \
--sdk-path=./testdata/sdk_config.yml
```



## vc-revoke

### 链上吊销VC

```shell
$ ./console vc-revoke add \
--id=16516616 \
--sdk-path=./testdata/sdk_config.yml
```



### 链上获取吊销列表

```shell
$ ./console vc-revoke list \
--search=111515515 \
--start=1 \
--count=10 \
--sdk-path=./testdata/sdk_config.yml
```



## vc-template

### VC模板上链

```shell
$ ./console vc-template add \
--temp-id=temp001 \
--temp-name=模板1 \
--temp-version=v1.0.0 \
--temp-path=./testdata/temp.json \
--sdk-path=./testdata/sdk_config.yml
```



 ### 获取VC模板

```shell
$ ./console vc-template get \
--temp-id=temp001 \
--temp-path=./testdata/temp.json \
--sdk-path=./testdata/sdk_config.yml
```



### 获取VC模板列表

```shell
$ ./console vc-template list \
--start=1 \
--count=10 \
--sdk-path=./testdata/sdk_config.yml
```



### 生成VC模板

```shell
$ ./console vc-template gen \
--map-key=name,age,sex \
--map-value=liu,18,man \
--temp-path=./testdata/temp.json
```



## vp

### 生成VP

```shell
$ ./console vp gen \
--sk-path=./testdata/sk.pem \
--algo=SM2 \
--holder=did:cm:admin \
--id=1231232 \
--vc-list=./testdata/vc1.json,./testdata/vc2.json \
--type=Identity \
--vp-path=./testdata/vp.json
```



### 链上验证VP

```shell
$ ./console vp verify \
--vp-path=./testdata/vp.json \
--sdk-path=./testdata/sdk_config.yml
```

