# Console

## admin

### 设置管理员

```shell
$ ./console admin add \
--admin-sdk-path=./testdata/sdk_config2.yml \
--sdk-path=./testdata/sdk_config.yml
```

```shell
## 要设置的管理员的长安链sdk配置
--admin-sdk-path 
## 合约创建者长安链sdk配置路径
--sdk-path
```



### 删除管理员

```shell
$ ./console admin delete \
--admin-sdk-path=./testdata/sdk_config2.yml \
--sdk-path=./testdata/sdk_config.yml
```

```shell
## 要设置的管理员的长安链sdk配置路径
--admin-sdk-path 
## 合约创建者长安链sdk
--sdk-path
```



### 查询是否拥有管理员权限

```shell
$ ./console admin auth \
--sdk-path=./testdata/sdk_config2.yml
```

```shell
## 查询对象使用的长安链sdk配置
--sdk-path
```



## key

### 生成公私钥

```shell
$ ./console key gen \
--algo=SM2 \
--pk-path=./testdata/pk.pem \
--sk-path=./testdata/sk.pem
```

```shell
## 公钥算法名称
--algo
## 生成的公钥存储路径
--pk-path
## 生成的私钥存储路径
--sk-path
```



## did

### 获取DID方法

```shell
$ ./console did method \
--sdk-path=./testdata/sdk_config.yml
```

```shell
## 长安链sdk配置路径
--sdk-path
```



### 根据公钥生成DID

```shell
$ ./console did gen \
--pk-path=./testdata/pk.pem \
--sdk-path=./testdata/sdk_config.yml
```

```shell
## 公钥PEM编码存储路径
--pk-path
## 长安链sdk配置路径
--sdk-path
```



### DID在链上是否有效

```shell
$ ./console did valid \
--did=did:cm:9h6JLhdJbDdPFGJrf2YaxQzj1UX2NmcWfzL65VhmvoUT \
--sdk-path=./testdata/sdk_config.yml
```

```shell
## 要查询的DID字符串
--did 
## 长安链sdk配置路径
--sdk-path
```



### 获取DID

```shell
$ ./console did get \
--pk-path=./testdata/pk.pem \
--sdk-path=./testdata/sdk_config.yml
```

```shell
## 公钥PEM编码存储路径
--pk-path
## 长安链sdk配置路径
--sdk-path
```



## doc

### 生成DID文档

```shell
$ ./console doc gen \
--sks-path=./testdata/sk.pem \
--pks-path=./testdata/pk.pem \
--controller=did:cm:test1,did:cm:test2 \
--sdk-path=./testdata/sdk_config.yml \
--doc-path=./testdata/doc.json
```

```shell
## DID文档中公钥对应的私钥路径（可配置多个，用 "," 隔开）
--sks-path
## DID文档中公钥路径（可配置多个，用 "," 隔开）
--pks-path
## DID文档中控制者DID字符串（如果不填，默认是其本身DID）
--controller
## 长安链sdk配置路径
--sdk-path
## 生成的DID文档路径
--doc-path
```



### DID文档上链

```shell
$ ./console doc add \
--doc-path=./testdata/doc.json \
--sdk-path=./testdata/sdk_config.yml 
```

```shell
## DID文档路径
--doc-path
## 长安链sdk配置路径
--sdk-path
```



### 获取DID文档

```shell
$ ./console doc get \
--did=did:cm:5hEKjps5VQjyVsSsugqQfEWoXh4qeJacHgnchH7cwWEf \
--sdk-path=./testdata/sdk_config.yml \
--doc-path=./testdata/doc.json
```

```shell
## DID字符串
-did
## DID文档路径
--doc-path
## 长安链sdk配置路径
--sdk-path
```



### 更新DID文档

```shell
$ ./console doc update-local  \
--sks-path=./testdata/sk.pem \
--pks-path=./testdata/pk.pem \
--controller=did:cm:test6 \
--old-doc-path=./testdata/doc.json \
--new-doc-path=./testdata/newdoc.json
```

```shell
## DID文档中公钥对应的私钥路径（可配置多个，用 "," 隔开）
--sks-path
## DID文档中公钥路径（可配置多个，用 "," 隔开）
--pks-path
## DID文档中控制者DID字符串（如果不填，默认是其本身DID）
--controller
## 更新前DID文档路径
--old-doc-path
## 更新后DID文档路径
--new-doc-path
```



### 链上更新DID文档

```shell
$ ./console doc update \
--doc-path=./testdata/newdoc.json \
--sdk-path=./testdata/sdk_config.yml 
```

```shell
## DID文档路径
--doc-path
## 长安链sdk配置路径
--sdk-path
```



## black

### 黑名单上链

```shell
$ ./console black add \
--dids=did:cm:test1,did:cm:9h6JLhdJbDdPFGJrf2YaxQzj1UX2NmcWfzL65VhmvoUT \
--sdk-path=./testdata/sdk_config.yml 
```

```shell
## 黑名单DID列表
--dids
## 长安链sdk配置路径
--sdk-path
```



### 链上获取黑名单

```shell
$ ./console black list \
--start=1 \
--count=10 \
--sdk-path=./testdata/sdk_config.yml
```

```shell
## 查询的列表的关键字，如果为空可查询全部
--search
## 查询的列表的起始索引，默认为1
--start
## 查询的列表的数量，默认为1000
--count
## 长安链sdk配置路径
--sdk-path
```



### 链上删除黑名单

```shell
$ ./console black delete \
--dids=did:cm:test1,did:cm:9h6JLhdJbDdPFGJrf2YaxQzj1UX2NmcWfzL65VhmvoUT \
--sdk-path=./testdata/sdk_config.yml
```

```shell
## 黑名单DID列表
--dids
## 长安链sdk配置路径
--sdk-path
```



## Issuer

### 链上添加权威签发者

```shell
$ ./console issuer add \
--dids=did:cm:5hEKjps5VQjyVsSsugqQfEWoXh4qeJacHgnchH7cwWEf \
--sdk-path=./testdata/sdk_config.yml
```

```shell
## 签发者DID列表
--dids
## 长安链sdk配置路径
--sdk-path
```



### 链上获取权威签发者列表

```shell
$ ./console issuer list \
--start=1 \
--count=10 \
--sdk-path=./testdata/sdk_config.yml
```

```shell
## 查询的列表的关键字，如果为空可查询全部
--search
## 查询的列表的起始索引，默认为1
--start
## 查询的列表的数量，默认为1000
--count
## 长安链sdk配置路径
--sdk-path
```



### 链上删除权威签发者

```shell
$ ./console issuer delete \
--dids=did:cm:test1 \
--sdk-path=./testdata/sdk_config.yml 
```

```shell
## 签发者DID列表
--dids
## 长安链sdk配置路径
--sdk-path
```



## vc

### 链上颁发VC

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

```shell
## 签发者私钥PEM编码文件路径
--sk-path
## 签发者公钥PEM编码文件路径
--pk-path
## 公钥索引，如果签发者DID文档中拥有多个公钥，需要指定索引，（可不填，默认为0）
--key-index
## 颁发主体内容的JSON文件路径
--subject
## 签发到期时间，采用（yy-mm-dd）格式
--expiration
## 自定义VcId，一般与业务相关
--id
## 模板ID编号（链上已有模板ID）
--temp-id
## VC的自定义类型，多个请用 "," 分隔，可不填，默认为`VerifiableCredential`
--type
## 生成的VC的JSON文件路径
--vc-path
## 长安链sdk配置路径
--sdk-path=./testdata/sdk_config.yml
```



### 本地颁发VC

```shell
$ ./console vc issue-local \
--sk-path=./testdata/sk.pem \
--subject=./testdata/subject.json \
--issuer=did:cm:admin \
--expiration=2025-01-25 \
--id=vc001 \
--temp-path=./testdata/temp.json \
--type=Identity \
--vc-path=./testdata/vc.json 
```

```shell
## 签发者私钥PEM编码文件路径
--sk-path
## 颁发主体内容的JSON文件路径
--subject
## 公钥索引，如果签发者DID文档中拥有多个公钥，需要指定索引，（可不填，默认为0）
--key-index
## 签发者的DID字符串
--issuer
## 签发到期时间，采用（yy-mm-dd）格式
--expiration
## 自定义VcId，一般与业务相关
--id
## 模板文件路径，VC模板通常是一个JSON schema
--temp-path
## VC的自定义类型，多个请用 "," 分隔，可不填，默认为`VerifiableCredential`
--type
## 生成的VC的JSON文件路径
--vc-path
```



### 链上验证VC的有效性

```shell
$ ./console vc verify \
--vc-path=./testdata/vc.json \
--sdk-path=./testdata/sdk_config.yml
```

```shell
## VC的JSON文件路径
--vc-path
## 长安链sdk配置路径
--sdk-path
```



### 获取VC签发日志

```shell
$ ./console vc log \
--search=vc001 \
--start=1 \
--count=10 \
--sdk-path=./testdata/sdk_config.yml
```

```shell
## 查询的列表的关键字，如果为空可查询全部
--search
## 查询的列表的起始索引，默认为1
--start
## 查询的列表的数量，默认为1000
--count
## 长安链sdk配置路径
--sdk-path
```



## vc-revoke

### 链上吊销VC

```shell
$ ./console vc-revoke add \
--id=vc001 \
--sdk-path=./testdata/sdk_config.yml
```

```shell
## 吊销的VC编号
--id
## 长安链sdk配置路径
--sdk-path
```



### 链上获取吊销列表

```shell
$ ./console vc-revoke list \
--start=1 \
--count=10 \
--sdk-path=./testdata/sdk_config.yml
```

```shell
## 查询的列表的关键字，如果为空可查询全部
--search
## 查询的列表的起始索引，默认为1
--start
## 查询的列表的数量，默认为1000
--count
## 长安链sdk配置路径
--sdk-path
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

```shell
## VC模板编号（根据业务自定义）
--temp-id
## VC模板名称
--temp-name
## VC模板版本
--temp-version
## VC模板内容路径（一般为JSON schema）
--temp-path
## 长安链sdk配置路径
--sdk-path
```



 ### 获取VC模板

```shell
$ ./console vc-template get \
--temp-id=temp001 \
--temp-path=./testdata/temp.json \
--sdk-path=./testdata/sdk_config.yml
```

```shell
## VC模板编号（根据业务自定义）
--temp-id
## 获取写入VC内容的文件路径
--temp-path
## 长安链sdk配置路径
--sdk-path
```



### 获取VC模板列表

```shell
$ ./console vc-template list \
--start=1 \
--count=10 \
--sdk-path=./testdata/sdk_config.yml
```

```shell
## 查询的列表的关键字，如果为空可查询全部
--search
## 查询的列表的起始索引，默认为1
--start
## 查询的列表的数量，默认为1000
--count
## 长安链sdk配置路径
--sdk-path
```



### 生成VC模板

```shell
$ ./console vc-template gen \
--map-key=name,age,sex \
--map-value=姓名,年龄,性别 \
--temp-path=./testdata/temp.json
```

```shell
## 模板内字段的key值列表（多个，使用","分隔）
--map-key
## 模板内字段的含义，与key值所（多个，使用","分隔）
--map-value
## 生成模板的JSON文件路径
--temp-path
```



## vp

### 生成VP

```shell
$ ./console vp gen \
--sk-path=./testdata/sk.pem \s
--holder=did:cm:5hEKjps5VQjyVsSsugqQfEWoXh4qeJacHgnchH7cwWEf \
--id=vp001 \
--vc-list=./testdata/vc.json \
--type=Identity \
--vp-path=./testdata/vp.json
```

```shell
## 持有者私钥PEM编码的文件路径
--sk-path
## 持有者的DID
--holder
## VP的编号（根据业务自定义）
--id
## VP包含的VC文件路径列表（可配置多个，使用","分隔）
--vc-list
## VP的自定义类型，多个请用 "," 分隔，可不填，默认为`VerifiablePresentation`
--type
## 生成的VP的JSON文件路径
--vp-path
```



### 链上验证VP

```shell
$ ./console vp verify \
--vp-path=./testdata/vp.json \
--sdk-path=./testdata/sdk_config.yml
```

```shell
## VP的JSON文件路径
--vp-path
## 长安链sdk配置路径
--sdk-path
```

