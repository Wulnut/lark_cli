---
title: 飞书项目OpenAPI全量搜索参数格式及常用示例
tags: [飞书项目, OpenAPI]
category: 飞书项目
created: 2026-03-17
updated: 2026-03-17

---


> 文档：[全量搜索参数格式及常用示例](https://project.feishu.cn/b/helpcenter/1p8d7djs/w11hyb8w)
> 版本：2.0.0
> 更新时间：2026-03-17
---
## 📋 概述
本文档详细介绍飞书项目 OpenAPI 中**全量搜索**的参数格式。全量搜索支持跨空间搜索，可在多个空间中查询工作项。
**注意**：筛选参数不支持字段类型为**公式计算**的字段。
---
## 🔧 SearchGroup 结构
### 基本结构
```json
{
  "search_group": {
    "search_params": [
      {
        "param_key": "字段key",
        "value": "值",
        "operator": "操作符"
      }
    ],
    "conjunction": "AND",  // 或 "OR"
    "search_groups": [
      // 嵌套的 SearchGroup
    ]
  }
}
```
### 字段说明
| 字段 | 类型 | 说明 |
|------|------|------|
| `search_params` | list\<SearchParam\> | 筛选条件列表 |
| `conjunction` | string | AND（且）/ OR（或） |
| `search_groups` | list\<SearchGroup\> | 嵌套的筛选组合 |
---
## 🔧 SearchParam 结构
### 基本结构
```json
{
  "param_key": "字段key",
  "value": "值",
  "operator": "操作符",
  "pre_operator": "前置操作符",
  "value_search_groups": {
    // 复合字段的子字段筛选
  }
}
```
### 字段说明
| 字段 | 类型 | 说明 |
|------|------|------|
| `param_key` | string | 筛选参数（字段 key 或特殊参数） |
| `value` | interface{} | 搜索字段值 |
| `operator` | string | 操作符类型 |
| `pre_operator` | string | 前置操作符（复合字段用） |
| `value_search_groups` | SearchGroup | 复合字段子字段筛选 |
---
## 📊 操作符枚举值 (18种)
| 编号 | 操作符 | 常量 | 说明 |
|------|--------|------|------|
| 1 | `~` | Reg | 匹配（模糊搜索） |
| 2 | `!~` | NReg | 不匹配 |
| 3 | `=` | Eq | 等于 |
| 4 | `!=` | Ne | 不等于 |
| 5 | `<` | Lt | 小于 |
| 6 | `>` | Gt | 大于 |
| 7 | `<=` | Lte | 小于等于 |
| 8 | `>=` | Gte | 大于等于 |
| 9 | `HAS ANY OF` | HasAnyOf | 存在选项属于 |
| 10 | `HAS NONE OF` | HasNoneOf | 不存在选项属于 |
| 11 | `IS NULL` | IsNull | 为空 |
| 12 | `IS NOT NULL` | NotNull | 不为空 |
| 13 | `CONTAINS` | Contains | 包含 |
| 14 | `NOT CONTAINS` | NotContains | 不包含 |
| 15 | `MEET` | Meet | 满足（复合字段） |
| 16 | `NOT MEET` | NotMeet | 不满足（复合字段） |
---
## 📊 前置操作符枚举值
| 编号 | 操作符 | 说明 |
|------|--------|------|
| 1 | `EVERY` | 每一组 |
| 2 | `ANY` | 存在一组 |
---
## 📌 固定参数 (Special Parameters)
| 参数名 | param_key | 支持操作符 | 值类型 | 说明 |
|--------|-----------|------------|--------|------|
| 进行中节点 | `current_nodes` | `=`, `!=`, `HAS ANY OF`, `HAS NONE OF`, `CONTAINS`, `NOT CONTAINS` | list\<string\> | 节点名称列表 |
| 流程节点 | `all_states` | 同上 | list\<string\> | 所有节点名称列表 |
| 流程节点时间 | `feature_state_time` | `<`, `>`, `<=`, `>=`, `IS NULL`, `IS NOT NULL` | object | 含 state_name, state_timestamp, state_condition |
| 全部人员 | `people` | `=`, `!=`, `HAS ANY OF`, `HAS NONE OF`, `IS NULL`, `IS NOT NULL`, `CONTAINS`, `NOT CONTAINS` | list\<string\> | user_key 列表 |
| 创建时间 | `created_at` | `=`, `!=`, `<`, `>`, `<=`, `>=` | int64 | 毫秒时间戳 |
| 更新时间 | `updated_at` | 同上 | int64 | 毫秒时间戳 |
| 节点负责人 | `node_owners` | 同 people | list\<object\> | 含 state_name 和 owners |
| 工作项ID | `work_item_id` | `=`, `!=`, `<`, `>`, `<=`, `>=`, `HAS ANY OF`, `HAS NONE OF` | list\<int64\> | 工作项 ID 列表 |
| 工作项状态 | `work_item_status` | `=`, `!=`, `HAS ANY OF`, `HAS NONE OF` | list\<string\> | start/closed 等 |
| 模板ID | `template_id` | `=`, `!=`, `HAS ANY OF`, `HAS NONE OF`, `IS NULL`, `IS NOT NULL` | list\<int64\> | 模板 ID 列表 |
| 业务线 | `business` | 同上 | list\<string\> | 业务线 ID 列表 |
| 角色人员 | `role_owners` | `=`, `!=`, `HAS ANY OF`, `HAS NONE OF`, `IS NULL`, `IS NOT NULL`, `CONTAINS`, `NOT CONTAINS` | list\<object\> | 含 role 和 owners |
### 流程节点时间参数示例
```json
{
  "state_name": "Android开发估分",  // 节点名称
  "state_timestamp": 1702310399000,  // 节点筛选时间戳
  "state_condition": 1  // 筛选的节点状态，开始：0，结束：1
}
```
### 节点负责人参数示例
```json
[
  {
    "state_name": "节点名称",
    "owners": ["user_key"]
  }
]
```
### 角色人员参数示例
```json
[
  {
    "role": "role_id",
    "owners": ["user_key"]
  }
]
```
---
## 📝 自定义字段类型与操作符
| 参数类型 | field_type_key | 支持的操作符 |
|----------|----------------|--------------|
| 单行文本/多行文本 | `text` | `~`, `!~`, `=`, `!=`, `IS NULL`, `IS NOT NULL` |
| 数字 | `number` | `=`, `!=`, `<`, `>`, `<=`, `>=`, `IS NULL`, `IS NOT NULL` |
| 链接 | `link` | `~`, `!~`, `=`, `!=`, `IS NULL`, `IS NOT NULL` |
| 开关 | `bool` | `=`, `!=`, `IS NULL`, `IS NOT NULL` |
| 系统外信号 | `signal` | `=`, `!=`, `HAS ANY OF`, `HAS NONE OF` |
| 单选 | `select` | `=`, `!=`, `HAS ANY OF`, `HAS NONE OF`, `IS NULL`, `IS NOT NULL` |
| 单选按钮 | `radio` | 同上 |
| 多选 | `multi-select` | 同上 + `CONTAINS`, `NOT CONTAINS` |
| 级联单选 | `tree-select` | `=`, `!=`, `HAS ANY OF`, `HAS NONE OF`, `IS NULL`, `IS NOT NULL` |
| 级联多选 | `tree-multi-select` | 同上 |
| 单选人员 | `user` | 同 select |
| 多选人员 | `multi-user` | 同上 + `CONTAINS`, `NOT CONTAINS` |
| 富文本 | `multi-text` | `~`, `!~`, `IS NULL`, `IS NOT NULL` |
| 关联单选 | `workitem_related_select` | `=`, `!=`, `HAS ANY OF`, `HAS NONE OF`, `IS NULL`, `IS NOT NULL` |
| 关联多选 | `workitem_related_multi_select` | 同上 + `CONTAINS`, `NOT CONTAINS` |
| 日期时间 | `precise_date` | `<`, `>`, `<=`, `>=`, `IS NULL`, `IS NOT NULL` |
| 日期 | `date` | 同上 |
| 复合字段 | `compound_field` | `IS NULL`, `IS NOT NULL`, `MEET`, `NOT MEET` |
---
## 📝 常用查询示例
### 示例1：查询需求下指定状态，且创建时间在某个区间的工作项列表
```json
{
  "pagination": {
    "page_size": 50,
    "page_num": 1
  },
  "data_sources": [
    {
      "project_key": "空间key",
      "work_item_type_keys": "story"
    }
  ],
  "search_group": {
    "conjunction": "AND",
    "search_params": [
      {
        "param_key": "created_at",
        "value": 1654064482000,
        "operator": ">"
      },
      {
        "param_key": "created_at",
        "value": 1654063482000,
        "operator": "<"
      },
      {
        "param_key": "work_item_status",
        "value": ["start"],
        "operator": "="
      }
    ]
  }
}
```
---
### 示例2：查询需求下包含指定人员的所有工作项
```json
{
  "pagination": {
    "page_size": 50,
    "page_num": 1
  },
  "data_sources": [
    {
      "project_key": "空间key",
      "work_item_type_keys": "story"
    }
  ],
  "search_group": {
    "conjunction": "AND",
    "search_params": [
      {
        "param_key": "people",
        "value": ["user_key"],
        "operator": "HAS ANY OF"
      }
    ]
  }
}
```
---
### 示例3：通过需求 ID 查询指定的需求
```json
{
  "pagination": {
    "page_size": 50,
    "page_num": 1
  },
  "data_sources": [
    {
      "project_key": "空间key",
      "work_item_type_keys": "story"
    }
  ],
  "search_group": {
    "conjunction": "AND",
    "search_params": [
      {
        "param_key": "work_item_id",
        "value": [12345, 45678],
        "operator": "HAS ANY OF"
      }
    ]
  }
}
```
---
### 示例4：查询一个需求下关联的所有缺陷
```json
{
  "pagination": {
    "page_size": 50,
    "page_num": 1
  },
  "data_sources": [
    {
      "project_key": "空间key",
      "work_item_type_keys": "issue"  // 缺陷类型
    }
  ],
  "search_group": {
    "conjunction": "AND",
    "search_params": [
      {
        "param_key": "_field_linked_story",
        "value": [12345],
        "operator": "="
      }
    ]
  }
}
```
---
### 示例5：查询一段时间内更新的工作项
```json
{
  "pagination": {
    "page_size": 50,
    "page_num": 1
  },
  "data_sources": [
    {
      "project_key": "空间key",
      "work_item_type_keys": "story"
    }
  ],
  "search_group": {
    "conjunction": "AND",
    "search_params": [
      {
        "param_key": "updated_at",
        "value": 1654064482000,
        "operator": ">"
      },
      {
        "param_key": "updated_at",
        "value": 1657064482000,
        "operator": "<"
      }
    ]
  }
}
```
---
## ⚠️ 注意事项
1. **data_sources 参数**：全量搜索需要指定 `data_sources`，包含 `project_key` 和 `work_item_type_keys`
2. **跨空间搜索**：可以传入多个空间的 `project_key` 实现跨空间搜索
3. **时间戳格式**：必须是毫秒时间戳（13位数字）
4. **分页**：`page_size` 最大值 100，`page_num` 从 1 开始
---
## 📎 相关文档
- [搜索参数格式及常用示例](飞书项目OpenAPI搜索参数格式及常用示例.md)
- [字段与属性解析格式](飞书项目OpenAPI字段类型汇总.md)
- [Open API 错误码](飞书项目OpenAPI开发者手册汇总.md)
---
## 🏷️ 标签
#飞书项目 #OpenAPI #全量搜索 #跨空间搜索 #SearchParam #API文档