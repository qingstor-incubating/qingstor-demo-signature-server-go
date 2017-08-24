# 签名服务器样例 (Go)
[![Build Status](https://travis-ci.org/Colin0114/qingstor-demo-signature-server-go.svg?branch=add_travis_ci)](https://travis-ci.org/Colin0114/qingstor-demo-signature-server-go)

[English](../../README.md) | 中文

本文档描述了如何使用 qingstor-sdk-go 来创建签名服务器样例，同时提供了 Go 版本的签名客户端样例。

关于签名服务器详见 [QingStor Demo Signature Server API Specs](https://github.com/yunify/qingstor-demo-signature-server-api-specs) 。

## 测试签名服务器样例

**注意：**

_1\. 运行环境要求安装 Go 和 Glide 。_

_2\. 签名服务器样例不适合运行在生产环境中，否则任何人均可以访问并进行签名验证。_

**运行签名服务器样例**

```bash
$ git clone https://github.com/yunify/qingstor-demo-signature-server-go.git
$ cd qingstor-demo-signature-server-go
# 修改 config_server.yaml 并配置你自己的密钥
$ cp config_server.yaml.example config_server.yaml
$ glide init
$ glide install
$ make run
```

**运行客户端样例**

```
$ cd your/path/to/qingstor-demo-signature-server-go
# 修改 config_client.yaml 并配置你自己的 Bucket 信息
$ cp config_client.yaml.example config_client.yaml
$ make test
```
## 参考文档
* [签名验证](https://docs.qingcloud.com/qingstor/api/common/signature.html)

## LICENSE
[The Apache License (Version 2.0, January 2004).](http://www.apache.org/licenses/LICENSE-2.0.html)
