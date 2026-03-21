---
title: 飞书项目OpenAPI概述
source: https://project.feishu.cn/b/helpcenter/1p8d7djs/4bsmoql6
author:
published: 2026-03-18
created: 2026-03-18
updated: 2026-03-18
description: 飞书项目OpenAPI概述，介绍开放平台提供的OpenAPI能力，实现数据的获取与操作，包含API结构、认证流程和基础概念
tags:
  - "飞书项目"
  - "OpenAPI"
  - "概述"
  - "认证"
category: "飞书项目"
related_docs:
  - "[[飞书项目API开发者知识库]]"
  - "[[飞书项目OpenAPI完整API列表]]"
  - "[[飞书项目OpenAPI获取访问凭证]]"
---

# Open API 概述

> 文档：[Open API 概述](https://project.feishu.cn/b/helpcenter/1p8d7djs/4bsmoql6)
> 更新时间：2026-03-18

---

## 概述

飞书项目开放平台提供 **Open API** 能力，可实现数据的获取与操作（与飞书项目页面的大部分操作等效）。开发者能够依据自身使用场景，丰富和优化业务流程，增强自身业务与现有工具之间的协同效应。

---

## API 结构介绍

### URL 结构

```
https://{Base URL}/open_api/:project_key/business/all
```

| 参数 | 说明 |
|------|------|
| `{Base URL}` | 实际访问域名，例如 `project.feishu.cn` |
| `:project_key` | 在飞书项目空间双击空间图标获取 |

### 完整示例

```
https://project.feishu.cn/open_api/projectkeyisme/work_item/create
```

---

## Header 参数

| 字段 | 必填 | 说明 |
|------|------|------|
| `Content-Type` | ✅ | `application/json` |
| `X-PLUGIN-TOKEN` | ✅ | 访问凭证，支持两种：<br>1. **插件身份凭证** (`plugin_token`)<br>2. **用户身份凭证** (`user_plugin_token`) |
| `X-USER-KEY` | 可选 | 当使用插件身份凭证时，指定接口调用的用户 user_key |
| `X-IDEM-UUID` | 可选 | 写类型接口的幂等串，用于防止重复提交 |
| `x-auth-mode` | 可选 | 权限校验模式：<br>- `1`：严格权限校验<br>- `0` 或不传：兼容模式 |

### x-auth-mode 说明

| 值 | 说明 |
|-----|------|
| `1` | 严格使用 X-User-Key 中的用户身份进行权限校验，不返回无权限资源 |
| `0` 或不传 | 遵循现有兼容性逻辑 |

---

## 服务端 API 调用限制

### 通用限制

- **频率限制**：**15 QPS**（每秒最多 15 次请求）

### 特殊限制接口

| 接口 | 限流阈值 |
|------|----------|
| `[POST]/open_api/view/v1/update_condition_view` | 15 QPS + 450 QPM |
| `[POST]/open_api/work_items/filter_across_project` | 15 QPS + 450 QPM |
| `[POST]/open_api/work_item/subtask/search` | 15 QPS + 450 QPM |
| `[POST]/open_api/:project_key/work_item/:work_item_type_key/search/params` | 15 QPS + 450 QPM |
| `[POST]/open_api/:project_key/work_item/filter` | 15 QPS + 450 QPM |
| `[POST]/open_api/:project_key/work_item/:work_item_type_key/:work_item_id/search_by_relation` | 10 QPS |
| `[POST]/open_api/work_item/actual_time/update` | 10 QPS |

---

## API 调用流程

```
1. 获取访问凭证
   ↓
2. 授权接口：在插件详情页授权需要调用的 Open API
   ↓
3. 调用 API
```

---

## 格式说明

### 相关文档

| 文档 | 说明 |
|------|------|
| [数据结构汇总](https://project.feishu.cn/b/helpcenter/1p8d7djs/1x1d372l) | 所有 Open API 的数据结构 |
| [字段与属性解析格式](https://project.feishu.cn/b/helpcenter/1p8d7djs/1tj6ggll) | 创建和查询的字段格式 |
| [Open API 错误码](https://project.feishu.cn/b/helpcenter/1p8d7djs/5aueo3jr) | 常见错误码及原因分析 |
| [搜索参数格式及常用示例](https://project.feishu.cn/b/helpcenter/1p8d7djs/1l8il0l6) | 搜索参数格式 |
| [全量搜索参数格式及常用示例](https://project.feishu.cn/b/helpcenter/1p8d7djs/w11hyb8w) | 全量搜索参数格式 |

---

## 相关文档

- [[飞书项目OpenAPI完整API列表]] - 完整 API 列表
- [[飞书项目OpenAPI依赖关系]] - API 依赖关系

---

## 标签

#飞书项目 #OpenAPI #概述 #API限制
