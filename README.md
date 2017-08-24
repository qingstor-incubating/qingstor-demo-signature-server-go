# QingStor Demo Signature Server (Go)
[![Build Status](https://travis-ci.org/Colin0114/qingstor-demo-signature-server-go.svg?branch=add_travis_ci)](https://travis-ci.org/Colin0114/qingstor-demo-signature-server-go)

English | [中文](./docs/zh_CN/README.md)

This project demonstrates how to use qingstor-sdk-go to create the QingStor Demo Signature Server,
and a demo client written in Go is also provided.

View [QingStor Demo Signature Server API Specs](https://github.com/yunify/qingstor-demo-signature-server-api-specs) for details about QingStor Signature Server.

## Getting Started

**Notices:**

_1\. Go & Glide is required_

_2\. The following QingStor Demo Signature Server cannot run on production environment.
Otherwise anyone can visit the server and get signed._

**Run QingStor Demo Signature Server**

```bash
$ git clone https://github.com/yunify/qingstor-demo-signature-server-go.git
$ cd qingstor-demo-signature-server-go
# Modify config_server.yaml and replace the access_key_id and secret_access_key with yours
$ cp config_server.yaml.example config_server.yaml
$ glide init
$ glide install
$ make run
```

**Run Client**

```bash
$ cd your/path/to/qingstor-demo-signature-server-go
# Modify config_client.yaml with your bucket information
$ cp config_client.yaml.example config_client.yaml
$ make test
```
## Reference Documentation
* [QingStor Object Storage Signature Verification](https://docs.qingcloud.com/qingstor/api/common/signature.html)

## LICENSE
[The Apache License (Version 2.0, January 2004).](http://www.apache.org/licenses/LICENSE-2.0.html)
