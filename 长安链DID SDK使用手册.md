# 长安链DID SDK使用手册

## 介绍

​	2019年，W3C提出了去中心化身份标识(Decentralized Identifier，DID)的标准，虽然W3C没有明确规范DID的实现需要使用区块链技术，但是很显然，与区块链技术结合，可以更好的实现DID。

​	为了方便长安链的用户能够快速的搭建自己的去中心化数字身份应用，我们基于长安链社区的去中心化数字身份（DID）合约标准`CM-CS-231201-DID`，实现了一套比较完整的DID智能合约，并且提供SDK开发工具包。

## 快速体验

### 搭建长安链

[链接快速搭建]

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

#### 使用CMC工具安装

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



#### 代码安装




### 使用命令行工具

## SDK

## 命令行工具

