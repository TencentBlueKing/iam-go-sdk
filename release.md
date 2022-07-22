版本日志
===============

## IAM v0.0.9

- [NEW] 支持StringContains 操作符 #25
- [OTHERS] 规范化所有操作符左值/右值, 并增加校验(校验失败直接False)

## IAM v0.0.8

- [BUGFIX] compareTwoValues will panic if got nil as input

## IAM v0.0.7

- [CHANGE] rename ValueEquals to ValueEqual

## IAM v0.0.6

- [NEW] 支持不同类型数值比较
- [NEW] 增加ValueEquals, 用于不同类型数值比较(原 Equals 只支持相同类型对象比较)

## IAM v0.0.5

- [NEW] support call apis via APIGateway

## IAM v0.0.4

- [NEW] add more IAM backend apis
    - GET /ping
    - GET /api/v1/model/systems/{system_id}/token
    - POST /api/v1/policy/query
    - POST /api/v1/policy/query_by_actions
    - POST /api/v1/policy/auth
    - POST /api/v1/policy/auth_by_resources
    - POST /api/v1/policy/auth_by_actions
    - GET /api/v1/systems/{system_id}/policies/{policy_id}
    - GET /api/v1/systems/{system_id}/policies
    - GET /api/v1/systems/{system_id}/policies/-/subjects

## IAM v0.0.3

- [NEW] support GetApplyURL, for generate apply url from iam
- [NEW] support IsAllowedWithCache, for cache the result for ttl
- [NEW] support BatchIsAllowed, one request, batch resources
- [NEW] support IsBasicAuthAllowed, for callback auth check
- [NEW] support GetToken, for show the system token
- [NEW] add logging module, support debug details

## IAM v0.0.2

- [NEW] support IsAllowed, basic expression eval

## IAM v0.0.1

- [NEW] init the project


