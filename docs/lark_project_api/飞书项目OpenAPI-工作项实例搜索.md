---
title: 飞书项目OpenAPI-工作项实例搜索
source: https://project.feishu.cn/b/helpcenter/1p8d7djs
author:
published: 2026-03-17
created: 2026-03-17
updated: 2026-03-17
description: 飞书项目OpenAPI工作项实例搜索详解，包含单空间搜索、跨空间搜索、全局搜索及关联工作项搜索等API
tags: [飞书项目, OpenAPI, 工作项, 搜索, SearchParams]
category: 飞书项目
related_docs:
  - "[[飞书项目API开发者知识库]]"
  - "[[飞书项目OpenAPI完整API列表]]"
  - "[[飞书项目OpenAPI搜索参数格式及常用示例]]"
  - "[[飞书项目OpenAPI全量搜索参数格式及常用示例]]"
---


> 文档编号：1
> 更新时间：2026-03-17
---
## 工作项实例搜索 API 列表
| # | API | 方法 | 说明 |
|---|-----|------|------|
| 1 | 获取指定的工作项列表（单空间） | POST | `/open_api/:project_key/work_item/filter` |
| 2 | 获取指定的工作项列表（跨空间） | POST | `/open_api/work_items/filter_across_project` |
| 3 | 获取指定的工作项列表（单空间-复杂传参） | POST | `/open_api/:project_key/work_item/:work_item_type_key/search/params` |
| 4 | 获取指定的工作项列表（全局搜索） | POST | `/open_api/compositive_search` |
| 5 | 获取指定的关联工作项列表（单空间） | POST | `/open_api/:project_key/work_item/:work_item_type_key/:work_item_id/search_by_relation` |
---
## 1. 获取指定的工作项列表（单空间）
**API**：`POST /open_api/:project_key/work_item/filter`
**说明**：在指定一个空间，搜索符合请求参数中传入条件的工作项实例列表
### 请求参数 (Body)
| 参数 | 类型 | 必须 | 说明 |
|------|------|------|------|
| `project_key` | string | ✅ | 空间 ID 或 simple_name |
| `work_item_type_keys` | list\<string\> | ✅ | 工作项类型 key 列表 |
| `search_group` | object | 可选 | 搜索条件，详见搜索参数格式文档 |
| `pagination` | object | 可选 | 分页信息 |
| `page_size` | int | 可选 | 每页数量（默认 100，最大 100） |
| `page_num` | int | 可选 | 页码（默认 1） |
| `orders` | list\<object\> | 可选 | 排序规则 |
| `select_all` | bool | 可选 | 是否返回所有字段（默认 false） |
| `field_keys` | list\<string\> | 可选 | 返回的字段 key 列表 |
### 请求示例
```json
{
  "project_key": "空间key",
  "work_item_type_keys": ["story"],
  "search_group": {
    "conjunction": "AND",
    "search_params": [
      {
        "param_key": "work_item_status",
        "value": ["start"],
        "operator": "="
      }
    ]
  },
  "pagination": {
    "page_size": 50,
    "page_num": 1
  }
}
```
### 返回参数
| 参数 | 类型 | 说明 |
|------|------|------|
| `err` | object | 错误信息 |
| `err_code` | int | 错误码 |
| `err_msg` | string | 错误信息描述 |
| `data` | object | 返回数据 |
| `data.items` | list | 工作项列表 |
| `data.total` | int | 总数 |
### 返回示例
```json
{
  "err": {},
  "err_code": 0,
  "data": {
    "items": [
      {
        "work_item_id": 12345,
        "work_item_type_key": "story",
        "name": "工作项名称",
        "status": "start"
      }
    ],
    "total": 100
  }
}
```
---
## 2. 获取指定的工作项列表（跨空间）
**API**：`POST /open_api/work_items/filter_across_project`
**说明**：跨多个空间，搜索符合请求参数中传入条件的工作项实例列表
### 请求参数 (Body)
| 参数 | 类型 | 必须 | 说明 |
|------|------|------|------|
| `project_keys` | list\<string\> | ✅ | 空间 key 列表 |
| `work_item_type_keys` | list\<string\> | ✅ | 工作项类型 key 列表 |
| `search_group` | object | 可选 | 搜索条件 |
| `pagination` | object | 可选 | 分页信息 |
### 请求示例
```json
{
  "project_keys": ["空间key1", "空间key2"],
  "work_item_type_keys": ["story", "task"],
  "search_group": {
    "conjunction": "AND",
    "search_params": [
      {
        "param_key": "work_item_status",
        "value": ["start"],
        "operator": "="
      }
    ]
  },
  "pagination": {
    "page_size": 50,
    "page_num": 1
  }
}
```
### 返回参数
同"获取指定的工作项列表（单空间）"
---
## 3. 获取指定的工作项列表（单空间-复杂传参）
**API**：`POST /open_api/:project_key/work_item/:work_item_type_key/search/params`
**说明**：在指定一个空间，搜索符合"复杂筛选条件"的工作项实例列表
### 请求参数 (Body)
| 参数 | 类型 | 必须 | 说明 |
|------|------|------|------|
| `project_key` | string | ✅ | 空间 key |
| `work_item_type_key` | string | ✅ | 工作项类型 key |
| `search_group` | object | ✅ | 搜索条件（复杂参数） |
| `pagination` | object | 可选 | 分页信息 |
| `select_all` | bool | 可选 | 是否返回所有字段 |
| `field_keys` | list\<string\> | 可选 | 返回的字段 |
### 请求示例
```json
{
  "project_key": "空间key",
  "work_item_type_key": "story",
  "search_group": {
    "conjunction": "AND",
    "search_params": [
      {
        "param_key": "created_at",
        "value": 1654064482000,
        "operator": ">"
      },
      {
        "param_key": "work_item_status",
        "value": ["start"],
        "operator": "="
      }
    ]
  },
  "pagination": {
    "page_size": 50,
    "page_num": 1
  }
}
```
---
## 4. 获取指定的工作项列表（全局搜索）
**API**：`POST /open_api/compositive_search`
**说明**：获取跨空间和工作项类型搜索符合条件的工作项实例列表
### 请求参数 (Body)
| 参数 | 类型 | 必须 | 说明 |
|------|------|------|------|
| `data_sources` | list\<object\> | ✅ | 数据源列表 |
| `data_sources[].project_key` | string | ✅ | 空间 key |
| `data_sources[].work_item_type_keys` | list\<string\> | ✅ | 工作项类型列表 |
| `search_group` | object | 可选 | 搜索条件 |
| `pagination` | object | 可选 | 分页信息 |
### 请求示例
```json
{
  "data_sources": [
    {
      "project_key": "空间key1",
      "work_item_type_keys": ["story"]
    },
    {
      "project_key": "空间key2",
      "work_item_type_keys": ["task"]
    }
  ],
  "search_group": {
    "conjunction": "AND",
    "search_params": [
      {
        "param_key": "work_item_status",
        "value": ["start"],
        "operator": "="
      }
    ]
  },
  "pagination": {
    "page_size": 50,
    "page_num": 1
  }
}
```
---
## 5. 获取指定的关联工作项列表（单空间）
**API**：`POST /open_api/:project_key/work_item/:work_item_type_key/:work_item_id/search_by_relation`
**说明**：获取与指定工作项实例存在工作项关联的工作项实例列表
### 请求参数 (Body)
| 参数 | 类型 | 必须 | 说明 |
|------|------|------|------|
| `project_key` | string | ✅ | 空间 key |
| `work_item_type_key` | string | ✅ | 工作项类型 key |
| `work_item_id` | int64 | ✅ | 工作项 ID |
| `relation_key` | string | ✅ | 关联关系 key |
| `search_group` | object | 可选 | 搜索条件 |
| `pagination` | object | 可选 | 分页信息 |
### 请求示例
```json
{
  "project_key": "空间key",
  "work_item_type_key": "story",
  "work_item_id": 12345,
  "relation_key": "relates_to",
  "pagination": {
    "page_size": 50,
    "page_num": 1
  }
}
```
---
## 相关文档
- [搜索参数格式及常用示例](飞书项目OpenAPI搜索参数格式及常用示例.md)
- [工作项实例读写](飞书项目OpenAPI-工作项实例读写.md)
- [工作项流程与节点](飞书项目OpenAPI-工作项流程与节点.md)
---
## 🏷️ 标签
#飞书项目 #OpenAPI #工作项搜索 #API文档