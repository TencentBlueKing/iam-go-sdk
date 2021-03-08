![](docs/resource/img/bk_iam_en.png)
---

[![license](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat)](https://github.com/TencentBlueKing/iam-go-sdk/blob/master/LICENSE.txt) [![Release Version](https://img.shields.io/badge/release-0.0.4-brightgreen.svg)](https://github.com/TencentBlueKing/iam-go-sdk/releases) [![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](https://github.com/TencentBlueKing/iam-go-sdk/pulls)

## Overview

iam-go-sdk is the SDK of blueking IAM(BK-IAM), your system can use BK-IAM easily via SDK.

## Features

- Authentication: IsAllowed / IsAllowedWithCache
- Batch Resource - Single Action Authentication: BatchIsAllowed
- Single Resource - Multiple Actions Authentication: ResourceMultiActionsAllowed
- Batch Resoure - Multiple Actions Authentication: BatchResourceMultiActionsAllowed
- Generate permission application URL: GetApplyURL
- Generate no permission json data: GenPermissionApplyData
- Callback request basic auth: IsBasicAuthAllowed  / 以及basic auth middleware
- Get Token of the system: GetToken
- Support prometheus metrics
- A simple framework for callback:dispatcher/provider interface

## Getting started

### Installation

```
$ go get -u github.com/TencentBlueKing/iam-go-sdk
```

### Usage

- [usage doc](docs/usage.md)
- [examples](examples/)

## Benchmark

```
$ make bench

BenchmarkExprCellEqual-12         	 6156810	       193 ns/op	      32 B/op	       1 allocs/op
BenchmarkExprCellNotEqual-12      	 6286323	       190 ns/op	      32 B/op	       1 allocs/op
BenchmarkExprCellLess-12          	 8583621	       142 ns/op	      32 B/op	       1 allocs/op
BenchmarkExprCellStartsWith-12    	 9665572	       126 ns/op	      32 B/op	       1 allocs/op
BenchmarkExprCellIn-12            	 2872075	       411 ns/op	      64 B/op	       3 allocs/op
```


## Roadmap

- [release log](release.md)

## Support

- [bk forum](https://bk.tencent.com/s-mart/community)
- [bk DevOps online video tutorial(In Chinese)](https://cloud.tencent.com/developer/edu/major-100008)
- Contact us, technical exchange QQ group:

<img src="https://github.com/Tencent/bk-PaaS/raw/master/docs/resource/img/bk_qq_group.png" width="250" hegiht="250" align=center />


## BlueKing Community

- [BK-CI](https://github.com/Tencent/bk-ci)：a continuous integration and continuous delivery system that can easily present your R & D process to you.
- [BK-BCS](https://github.com/Tencent/bk-bcs)：a basic container service platform which provides orchestration and management for micro-service business.
- [BK-BCS-SaaS](https://github.com/Tencent/bk-bcs-saas)：a SaaS provides users with highly scalable, flexible and easy-to-use container products and services.
- [BK-PaaS](https://github.com/Tencent/bk-PaaS)：an development platform that allows developers to create, develop, deploy and manage SaaS applications easily and quickly.
- [BK-SOPS](https://github.com/Tencent/bk-sops)：an lightweight scheduling SaaS  for task flow scheduling and execution through a visual graphical interface. 
- [BK-CMDB](https://github.com/Tencent/bk-cmdb)：an enterprise-level configuration management platform for assets and applications.

## Contributing

If you have good ideas or suggestions, please let us know by Issues or Pull Requests and contribute to the Blue Whale Open Source Community.

## License

Based on the MIT protocol. Please refer to [LICENSE](LICENSE.txt)
