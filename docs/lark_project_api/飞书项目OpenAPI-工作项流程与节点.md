---
title: 飞书项目OpenAPI-工作项流程与节点
source: https://project.feishu.cn/b/helpcenter/1p8d7djs
author:
published: 2026-03-17
created: 2026-03-17
updated: 2026-03-17
description: 飞书项目OpenAPI工作项流程与节点详解，包含获取工作流详情、WBS视图、节点操作、状态流转等API
tags: [飞书项目, OpenAPI, 工作项, 流程, 节点, 状态流转]
category: 飞书项目
related_docs:
  - "[[飞书项目API开发者知识库]]"
  - "[[飞书项目OpenAPI完整API列表]]"
  - "[[飞书项目OpenAPI字段类型汇总]]"
---


> 文档编号：3
> 更新时间：2026-03-17
---
## 工作项流程与节点 API 列表
| # | API | 方法 | 说明 |
|---|-----|------|------|
| 1 | 获取工作流详情 | POST | `/open_api/:project_key/work_item/:work_item_type_key/:work_item_id/workflow/query` |
| 2 | 获取工作流详情（WBS） | GET | `/open_api/:project_key/work_item/:work_item_type_key/:work_item_id/wbs_view` |
| 3 | 更新节点/排期 | PUT | `/open_api/:project_key/workflow/:work_item_type_key/:work_item_id/node/:node_id` |
| 4 | 节点完成/回滚 | POST | `/open_api/:project_key/workflow/:work_item_type_key/:work_item_id/node/:node_id/operate` |
| 5 | 状态流转 | POST | `/open_api/:project_key/workflow/:work_item_type_key/:work_item_id/node/state_change` |
| 6 | 获取指定节点/状态流转所需必填信息 | POST | `/open_api/work_item/transition_required_info/get` |
---
## 1. 获取工作流详情
**API**：`POST /open_api/:project_key/work_item/:work_item_type_key/:work_item_id/workflow/query`
**说明**：用于获取指定空间和工作项类型下的一个工作项实例的工作流信息，包括节点的状态、负责人、估分以及表单、子任务等
### 请求参数
| 参数 | 类型 | 必须 | 说明 |
|------|------|------|------|
| `project_key` | string | ✅ | 空间 key |
| `work_item_type_key` | string | ✅ | 工作项类型 key |
| `work_item_id` | int64 | ✅ | 工作项 ID |
### 返回参数
| 参数 | 类型 | 说明 |
|------|------|------|
| `workflow_type` | string | 工作流类型：node_flow（节点流）/ status_flow（状态流） |
| `nodes` | list | 节点列表 |
| `nodes[].node_id` | string | 节点 ID |
| `nodes[].node_name` | string | 节点名称 |
| `nodes[].status` | string | 节点状态：pending/running/completed |
| `nodes[].owners` | list | 节点负责人 |
| `nodes[].estimate` | float | 估分 |
| `subtasks` | list | 子任务列表 |
### 返回示例
```json
{
  "err": {},
  "err_code": 0,
  "data": {
    "workflow_type": "node_flow",
    "nodes": [
      {
        "node_id": "node_1",
        "node_name": "需求评审",
        "status": "completed",
        "owners": ["user_key_1"],
        "estimate": 3.0
      },
      {
        "node_id": "node_2",
        "node_name": "开发中",
        "status": "running",
        "owners": ["user_key_2"],
        "estimate": 5.0
      }
    ],
    "subtasks": []
  }
}
```
---
## 2. 获取工作流详情（WBS）
**API**：`GET /open_api/:project_key/work_item/:work_item_type_key/:work_item_id/wbs_view`
**说明**：用于获取行业专版中一个节点流工作项实例的WBS工作流信息
### 返回参数
| 参数 | 类型 | 说明 |
|------|------|------|
| `wbs_nodes` | list | WBS 节点列表 |
| `deliverables` | list | 交付物列表 |
---
## 3. 更新节点/排期
**API**：`PUT /open_api/:project_key/workflow/:work_item_type_key/:work_item_id/node/:node_id`
**说明**：用于更新一个工作项实例的指定节点信息（节点流），包括节点负责人、排期和表单信息等
### 请求参数
| 参数 | 类型 | 必须 | 说明 |
|------|------|------|------|
| `node_id` | string | ✅ | 节点 ID |
| `owner_ids` | list\<string\> | 可选 | 节点负责人 |
| `start_time` | int64 | 可选 | 开始时间（毫秒） |
| `end_time` | int64 | 可选 | 结束时间（毫秒） |
| `estimate` | float | 可选 | 估分 |
| `fields` | object | 可选 | 表单字段值 |
### 请求示例
```json
{
  "node_id": "node_2",
  "owner_ids": ["user_key_1", "user_key_2"],
  "start_time": 1704067200000,
  "end_time": 1706659200000,
  "estimate": 8.0
}
```
---
## 4. 节点完成/回滚
**API**：`POST /open_api/:project_key/workflow/:work_item_type_key/:work_item_id/node/:node_id/operate`
**说明**：用于完成或者回滚一个工作项实例的指定节点（节点流），同时更新节点信息
### 请求参数
| 参数 | 类型 | 必须 | 说明 |
|------|------|------|------|
| `node_id` | string | ✅ | 节点 ID |
| `operation` | string | ✅ | 操作类型：`complete`（完成）/ `rollback`（回滚） |
| `owner_ids` | list\<string\> | 可选 | 节点负责人 |
| `start_time` | int64 | 可选 | 开始时间 |
| `end_time` | int64 | 可选 | 结束时间 |
| `estimate` | float | 可选 | 估分 |
| `fields` | object | 可选 | 表单字段值 |
### 请求示例
```json
{
  "node_id": "node_2",
  "operation": "complete",
  "owner_ids": ["user_key_1"],
  "fields": {
    "field_key": "value"
  }
}
```
---
## 5. 状态流转
**API**：`POST /open_api/:project_key/workflow/:work_item_type_key/:work_item_id/node/state_change`
**说明**：用于流转一个工作项实例到指定状态（状态流），同时更新节点信息
### 请求参数
| 参数 | 类型 | 必须 | 说明 |
|------|------|------|------|
| `target_state` | string | ✅ | 目标状态 |
| `owner_ids` | list\<string\> | 可选 | 节点负责人 |
| `fields` | object | 可选 | 表单字段值 |
### 请求示例
```json
{
  "target_state": "done",
  "owner_ids": ["user_key_1"],
  "fields": {
    "field_key": "value"
  }
}
```
---
## 6. 获取指定节点/状态流转所需必填信息
**API**：`POST /open_api/work_item/transition_required_info/get`
**说明**：用于获取一个工作项实例指定节点流转所需的必填信息
### 请求参数
| 参数 | 类型 | 必须 | 说明 |
|------|------|------|------|
| `project_key` | string | ✅ | 空间 key |
| `work_item_type_key` | string | ✅ | 工作项类型 key |
| `work_item_id` | int64 | ✅ | 工作项 ID |
| `target_node_id` | string | 可选 | 目标节点 ID（节点流） |
| `target_state` | string | 可选 | 目标状态（状态流） |
### 返回参数
| 参数 | 类型 | 说明 |
|------|------|------|
| `required_fields` | list | 必填表单项 |
| `required_node_fields` | list | 必填节点字段 |
| `required_subtasks` | list | 必填子任务 |
| `required_deliverables` | list | 必填交付物 |
---
## 相关文档
- [工作项实例读写](飞书项目OpenAPI-工作项实例读写.md)
- [工作项实例搜索](飞书项目OpenAPI-工作项实例搜索.md)
- [子任务管理](飞书项目OpenAPI-子任务.md)
---
## 🏷️ 标签
#飞书项目 #OpenAPI #工作流 #节点 #API文档