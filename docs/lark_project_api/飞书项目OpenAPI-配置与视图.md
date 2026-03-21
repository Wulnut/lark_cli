---
title: 飞书项目OpenAPI-配置与视图
source: https://project.feishu.cn/b/helpcenter/1p8d7djs
author:
published: 2026-03-17
created: 2026-03-17
updated: 2026-03-17
description: 飞书项目OpenAPI配置与视图详解，包含空间配置、工作项配置、自定义字段、关系列表、视图管理等API
tags: [飞书项目, OpenAPI, 配置, 视图, 字段, 工作项类型]
category: 飞书项目
related_docs:
  - "[[飞书项目API开发者知识库]]"
  - "[[飞书项目OpenAPI完整API列表]]"
  - "[[飞书项目OpenAPI字段类型汇总]]"
---


> 文档编号：7
> 更新时间：2026-03-17
---
## 配置 API 总览
### 空间配置 (3 APIs)
| # | API | 方法 | 说明 |
|---|-----|------|------|
| 1 | 获取空间下业务线详情 | GET | `/open_api/:project_key/business/all` |
| 2 | 获取空间下工作项类型 | GET | `/open_api/:project_key/work_item/all-types` |
| 3 | 获取空间下团队人员 | GET | `/open_api/:project_key/teams/all` |
### 工作项配置 (9 APIs)
| # | API | 方法 | 说明 |
|---|-----|------|------|
| 1 | 获取工作项基础信息配置 | GET | `/open_api/:project_key/work_item/type/:work_item_type_key` |
| 2 | 更新工作项基础信息配置 | PUT | `/open_api/:project_key/work_item/type/:work_item_type_key` |
| 3 | 获取字段信息 | POST | `/open_api/:project_key/field/all` |
| 4 | 创建自定义字段 | POST | `/open_api/:project_key/field/:work_item_type_key/create` |
| 5 | 更新自定义字段 | PUT | `/open_api/:project_key/field/:work_item_type_key` |
| 6 | 获取工作项关系列表 | GET | `/open_api/:project_key/work_item/relation` |
| 7 | 新增工作项关系 | POST | `/open_api/work_item/relation/create` |
| 8 | 更新工作项关系 | POST | `/open_api/work_item/relation/update` |
| 9 | 删除工作项关系 | DELETE | `/open_api/work_item/relation/delete` |
### 流程模板配置 (5 APIs)
| # | API | 方法 | 说明 |
|---|-----|------|------|
| 1 | 获取工作项下的流程模板列表 | GET | `/open_api/:project_key/template_list/:work_item_type_key` |
| 2 | 获取流程模板配置详情 | GET | `/open_api/:project_key/template_detail/:template_id` |
| 3 | 新增流程模板 | POST | `/open_api/template/v2/create_template` |
| 4 | 更新流程模板 | PUT | `/open_api/template/v2/update_template` |
| 5 | 删除流程模板 | DELETE | `/open_api/template/v2/delete_template/:project_key/:template_id` |
### 流程角色配置 (1 API)
| # | API | 方法 | 说明 |
|---|-----|------|------|
| 1 | 获取流程角色配置详情 | GET | `/open_api/:project_key/flow_roles/:work_item_type_key` |
---
## 视图 API 列表 (8 APIs)
| # | API | 方法 | 说明 |
|---|-----|------|------|
| 1 | 获取视图列表及配置信息 | POST | `/open_api/:project_key/view_conf/list` |
| 2 | 获取视图下工作项列表 | GET | `/open_api/:project_key/fix_view/:view_id` |
| 3 | 获取视图下工作项列表（全景视图） | POST | `/open_api/:project_key/view/:view_id` |
| 4 | 创建固定视图 | POST | `/open_api/:project_key/:work_item_type_key/fix_view` |
| 5 | 更新固定视图 | POST | `/open_api/:project_key/:work_item_type_key/fix_view/:view_id` |
| 6 | 创建条件视图 | POST | `/open_api/view/v1/create_condition_view` |
| 7 | 更新条件视图 | POST | `/open_api/view/v1/update_condition_view` |
| 8 | 删除视图 | DELETE | `/open_api/:project_key/fix_view/:view_id` |
---
## 租户 API 列表 (2 APIs)
| # | API | 方法 | 说明 |
|---|-----|------|------|
| 1 | 获取租户信息 | GET | `/open_api/tenant/info` |
| 2 | 获取租户安装空间列表 | GET | `/open_api/tenant/installed_projects` |
---
## 空间关联 API 列表 (4 APIs)
| # | API | 方法 | 说明 |
|---|-----|------|------|
| 1 | 获取空间关联规则列表 | POST | `/open_api/:project_key/relation/rules` |
| 2 | 获取空间关联下的关联工作项实例列表 | POST | `/open_api/:project_key/relation/:work_item_type_key/:work_item_id/work_item_list` |
| 3 | 绑定空间关联的关联工作项实例 | POST | `/open_api/:project_key/relation/:work_item_type_key/:work_item_id/batch_bind` |
| 4 | 解绑空间关联的关联工作项实例 | DELETE | `/open_api/:project_key/relation/:work_item_type_key/:work_item_id` |
---
## 常用配置 API 详解
### 1. 获取空间下工作项类型
**API**：`GET /open_api/:project_key/work_item/all-types`
**说明**：获取空间下所有的工作项类型
### 返回参数
| 参数 | 类型 | 说明 |
|------|------|------|
| `work_item_types` | list | 工作项类型列表 |
| `work_item_types[].key` | string | 类型 key |
| `work_item_types[].name` | string | 类型名称 |
| `work_item_types[].type` | string | 类型（story/task/issue 等） |
---
### 2. 获取字段信息
**API**：`POST /open_api/:project_key/field/all`
**说明**：获取指定空间或一个工作项类型下所有字段的基础信息
### 请求参数
| 参数 | 类型 | 必须 | 说明 |
|------|------|------|------|
| `work_item_type_key` | string | 可选 | 工作项类型 key |
| `field_keys` | list\<string\> | 可选 | 字段 key 列表 |
### 返回参数
| 参数 | 类型 | 说明 |
|------|------|------|
| `fields` | list | 字段列表 |
| `fields[].field_key` | string | 字段 key |
| `fields[].field_name` | string | 字段名称 |
| `fields[].field_type_key` | string | 字段类型 |
| `fields[].required` | bool | 是否必填 |
---
### 3. 获取视图列表
**API**：`POST /open_api/:project_key/view_conf/list`
**说明**：用于在指定空间，搜索符合请求参数中传入条件的视图列表及相关配置信息
### 请求参数
| 参数 | 类型 | 必须 | 说明 |
|------|------|------|------|
| `view_ids` | list\<string\> | 可选 | 视图 ID 列表 |
| `view_type` | string | 可选 | 视图类型 |
---
## 相关文档
- [完整 API 列表](飞书项目OpenAPI完整API列表.md)
- [字段类型汇总](飞书项目OpenAPI字段类型汇总.md)
---
## 🏷️ 标签
#飞书项目 #OpenAPI #配置 #视图 #租户 #API文档