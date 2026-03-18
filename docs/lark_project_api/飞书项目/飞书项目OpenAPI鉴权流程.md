---
title: 飞书项目OpenAPI鉴权流程
tags: [飞书项目, OpenAPI, 鉴权]
category: 飞书项目
created: 2026-03-17
updated: 2026-03-17
---

# 飞书项目 OpenAPI 鉴权流程

> 文档编号：附录1
> 更新时间：2026-03-17

---

## 概述

飞书项目 OpenAPI 使用 **OAuth 2.0** 协议进行身份验证和授权。

---

## 访问凭证类型

| 凭证类型 | 适用场景 | 说明 |
|----------|----------|------|
| **plugin_access_token** | 插件 API 调用 | 使用应用凭证获取，不依赖具体用户 |
| **user_access_token** | 用户代表调用 | 代表当前用户操作，需要用户授权 |

---

## Header 配置

| Header | 说明 | 示例 |
|--------|------|------|
| `Authorization` | 访问凭证 | `Bearer YOUR_ACCESS_TOKEN` |
| `Content-Type` | 请求内容类型 | `application/json; charset=utf-8` |
| `X-USER-KEY` | 用户标识 | 使用 user_access_token 时必填 |

---

## 获取 plugin_access_token

```
POST https://open.feishu.cn/open-apis/auth/v3/tenant_access_token/internal
```

### 请求体

```json
{
  "app_id": "cli_xxxxxxxxxxxxxx",
  "app_secret": "xxxxxxxxxxxxxxxxxxxxxx"
}
```

### 返回

```json
{
  "err": {},
  "err_code": 0,
  "tenant_access_token": "tGwt3bHZxxxxxxxxxxxxxx",
  "expire": 7200
}
```

> **注意**：`tenant_access_token` 有效期为 **2 小时**，需要缓存并定期刷新。

---

## 获取 user_access_token

### 第一步：获取授权码 (code)

用户访问以下 URL 进行授权：

```
https://open.feishu.cn/open-apis/oauth/authorize?app_id=APP_ID&redirect_uri=ENCODED_URL&state=STATE
```

### 第二步：换取 access_token

```
POST https://open.feishu.cn/open-apis/oauth/v3/access_token
```

### 请求体

```json
{
  "grant_type": "authorization_code",
  "code": "授权码"
}
```

---

## 常见鉴权错误

| 错误码 | 说明 | 解决方案 |
|--------|------|----------|
| 10021 | Token Not Exist | 请求头未传 access_token |
| 10022 | Check Token Failed | token 校验失败 |
| 10301 | Check Token Perm Failed | 权限不足 |
| 10302 | User Is Resigned | 用户已离职 |

---

## 相关文档

- [[飞书项目OpenAPI开发者手册汇总]] - 错误码详解
- [[飞书项目OpenAPI完整API列表]] - 主索引
- [[飞书项目OpenAPI依赖关系]] - API依赖关系

---

## 标签

#飞书项目 #OpenAPI #鉴权 #OAuth
