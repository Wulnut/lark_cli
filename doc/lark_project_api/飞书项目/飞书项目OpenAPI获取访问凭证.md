---
title: 飞书项目OpenAPI获取访问凭证
tags: [飞书项目, OpenAPI, 访问凭证]
category: 飞书项目
created: 2026-03-18
updated: 2026-03-18
---

# 获取访问凭证

> 文档：[获取访问凭证](https://project.feishu.cn/b/helpcenter/1p8d7djs/4id4bvnf)
> 更新时间：2026-03-18

---

## 概述

要调用飞书项目开放平台 OpenAPI 前，插件需要获取对应的访问凭证。访问凭证代表插件从平台/用户手中获取的授权，包含插件信息和调用者的身份信息。

---

## 访问凭证类型

| 凭证类型 | 说明 | 示例 | 获取方法 |
|----------|------|------|----------|
| **plugin_access_token** | 插件访问凭证。用于获取用户访问凭证或调用 OpenAPI（需与 X-User-Key 搭配使用） | `p-b0dc4a7d-...` | 获取插件访问凭证 |
| **virtual_plugin_token** | 虚拟插件访问凭证。仅供开发调试使用。空间范围受数据范围影响 | `v-b0dc4a7d-...` | 获取插件访问凭证 |
| **user_access_token** | 用户访问凭证。代表用户对插件的临时授权，插件将代表用户执行操作 | `u-f09b09b8-...` | 获取用户访问凭证 |

---

## 使用访问凭证

飞书项目开放平台提供两种调用方式：

### 方式一：使用 user_access_token

```bash
curl --location 'https://{域名}/open_api/projects/detail' \
--header 'X-Plugin-Token: {{user_token}}' \
--header 'Content-Type: application/json' \
--data '{
  "simple_names": ["test_new"],
  "user_key": "test"
}'
```

### 方式二：使用 plugin_access_token + X-User-Key

```bash
curl --location 'https://{域名}/open_api/projects/detail' \
--header 'X-Plugin-Token: {{plugin_token}}' \
--header 'X-User-Key: {{user_key}}' \
--header 'Content-Type: application/json' \
--data '{
  "simple_names": ["test_new"],
  "user_key": "test"
}'
```

> **注意**：开发调试时使用 `virtual_plugin_token`

---

## 获取插件访问凭证 (plugin_access_token)

### 接口信息

| 项目 | 值 |
|------|-----|
| 请求方式 | POST |
| 请求地址 | `https://{域名}/open_api/authen/plugin_token` |

### 请求头参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| Content-Type | string | ✅ | 固定值：`application/json` |

### 请求体参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| plugin_id | string | ✅ | 插件唯一标识，Plugin ID |
| plugin_secret | string | ✅ | 插件密钥，Plugin Secret |
| type | int | ❌ | 凭证类型：0=plugin_access_token（默认），1=virtual_plugin_token |

### 请求示例

```bash
curl --location 'https://{域名}/open_api/authen/plugin_token' \
--header 'Content-Type: application/json' \
--data '{
  "plugin_id": "MII_63E9D49B8C82****",
  "plugin_secret": "D01B5F1A191C8620D133CDC371C0****",
  "type": 0
}'
```

### 响应体参数

```json
{
  "data": {
    "expire_time": 7200,
    "token": "p-49257489-f7d7-4cd6-b34f-98c6b81d****"
  },
  "error": {
    "code": 0,
    "msg": "success"
  }
}
```

### 注意事项

1. `plugin_access_token` 有效期为 **7200 秒（2 小时）**
2. 过期后获取会返回新的 token
3. 开发者需要**缓存** plugin_access_token 用于后续调用

---

## 获取用户访问凭证 (user_access_token)

### 流程概述

```
用户授权 → 获取授权码 code → 通过 code 获取 user_access_token
```

### 步骤 1：获取授权码 code

- 使用客户端 API `getAuthCode` 获取授权码
- 授权码有效期 **5 分钟**，只能使用一次

```javascript
const { code } = await window.JSSDK.utils.getAuthCode();
console.log('code', code);
```

### 步骤 2：通过授权码获取 user_access_token

| 项目 | 值 |
|------|-----|
| 请求方式 | POST |
| 请求地址 | `https://{域名}/open_api/authen/user_plugin_token` |

### 请求头参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| Content-Type | string | ✅ | 固定值：`application/json` |
| X-Plugin-Token | string | ✅ | plugin_access_token |

### 请求体参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| code | string | ✅ | 授权码 |
| grant_type | string | ✅ | 固定值：`authorization_code` |

### 请求示例

```bash
curl --location 'https://{域名}/open_api/authen/user_plugin_token' \
--header 'X-Plugin-Token: p-694b3fcd-5eec-4380-a4e4-9c09699a****' \
--header 'Content-Type: application/json' \
--data '{
  "code": "7849d137b6d947a5bdee470eedd6****",
  "grant_type": "authorization_code"
}'
```

### 响应体参数

```json
{
  "data": {
    "token": "u-efe05d32-b733-4b15-a960-6acadf8e****",
    "expire_time": 7200,
    "refresh_token": "85cde42e-76c8-4677-a783-1e5761f7****",
    "refresh_token_expire_time": 1209600,
    "saas_tenant_key": "DeepCollaboration",
    "user_key": "he"
  },
  "error": {
    "code": 0,
    "msg": "success"
  }
}
```

---

## 刷新用户访问凭证

### 接口信息

| 项目 | 值 |
|------|-----|
| 请求方式 | POST |
| 请求地址 | `https://{域名}/open_api/authen/refresh_token` |

### 请求头参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| Content-Type | string | ✅ | 固定值：`application/json` |
| X-Plugin-Token | string | ✅ | plugin_access_token |

### 请求体参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| refresh_token | string | ✅ | 刷新 token |
| type | string | ✅ | 固定值：`1` |

### 响应体参数

```json
{
  "data": {
    "expire_time": 7200,
    "token": "u-d098fd08-61fc-4d6e-8864-f653785c****",
    "refresh_token": "8311cee4-1d47-4351-9553-300ae4a****",
    "refresh_token_expire_time": 1209600
  },
  "error": {
    "code": 0,
    "msg": "success"
  }
}
```

---

## 资源访问权限

OpenAPI 是否允许调用受以下因素影响：

1. **插件权限**：插件是否申请对应权限
2. **空间安装**：对应空间是否安装拥有对应权限的插件版本
3. **用户权限**：用户是否拥有对应空间或数据的访问权限

> **注意**：使用 `virtual_plugin_token` 不受第 2 点限制，但只能获取数据范围内设置的空间数据

---

## 相关文档

- [[飞书项目OpenAPI概述]] - Open API 概述
- [[飞书项目OpenAPI完整API列表]] - 完整 API 列表
- [[飞书项目OpenAPI鉴权流程]] - 鉴权流程
- [获取授权码 - 飞书开放平台](https://open.feishu.cn/document/authentication-management/access-token/obtain-oauth-code)

---

## 标签

#飞书项目 #OpenAPI #访问凭证 #plugin_token #user_token
