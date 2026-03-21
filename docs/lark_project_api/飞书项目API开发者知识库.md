---
title: 飞书项目 API 开发者知识库
source: 飞书项目开发者手册汇总
author:
published:
created: 2026-03-19
updated: 2026-03-19
description: 飞书项目 Open API 开发者手册知识库，包含格式说明、字段属性、操作符和错误码汇总
tags: [飞书项目, API, OpenAPI, 开发手册, 知识库, SearchParams, MQL, 错误码, 字段类型]
category: 飞书项目
aliases:
  - 飞书项目开发者知识库
  - 飞书项目API手册
  - 飞书项目知识库
related_docs:
  - "[[API 列表 - 开发者手册 - 飞书项目帮助中心]]"
  - "[[飞书项目OpenAPI完整API列表]]"
  - "[[全量搜索参数格式及常用示例 - 开发者手册 - 飞书项目帮助中心]]"
  - "[[搜索参数格式及常用示例 - 开发者手册 - 飞书项目帮助中心]]"
  - "[[字段与属性解析格式 - 开发者手册 - 飞书项目帮助中心]]"
  - "[[MQL 语法说明 - 开发者手册 - 飞书项目帮助中心]]"
  - "[[Open API 错误码 - 开发者手册 - 飞书项目帮助中心]]"
  - "[[数据结构汇总 - 开发者手册 - 飞书项目帮助中心]]"
---

# 飞书项目 Open API 开发者知识库

> 本文档汇总了飞书项目 Open API 的核心开发知识，包含搜索参数格式、字段属性、数据结构及常见错误码。

---

## 目录

1. [[#一、核心概念与数据结构]]
2. [[#二、搜索参数格式 (SearchParams)]]
3. [[#三、字段与属性 (Fields & Properties)]]
4. [[#四、MQL 查询语法]]
5. [[#五、错误码参考]]
6. [[#六、关系图谱]]

---

## 一、核心概念与数据结构

### 1.1 数据模型概述

飞书项目的数据模型以 **工作项 (WorkItem)** 为核心，包含以下关键实体：

| 实体 | 说明 | 关键字段 |
|------|------|---------|
| **Project** | 空间 | project_key, name, business |
| **WorkItem** | 工作项 | id, name, work_item_type_key, status |
| **Field** | 字段 | field_key, field_type_key, field_value |
| **Template** | 模板 | template_id, version |
| **Comment** | 评论 | id, content, created_by |
| **Business** | 业务线 | id, name, role_owners |

### 1.2 SearchGroup 结构 (搜索筛选)

```
SearchGroup (筛选组合)
├── search_params: List<SearchParam>  // 筛选条件列表
├── conjunction: "AND" | "OR"          // 逻辑关系
└── search_groups: List<SearchGroup>  // 嵌套筛选组
```

### 1.3 SearchParam 结构 (筛选参数)

```
SearchParam (筛选参数)
├── param_key: string           // 字段key或固定参数key
├── value: interface{}          // 筛选值
├── operator: string             // 操作符
├── pre_operator: string         // 前置操作符 (复合字段用)
└── value_search_groups: SearchGroup  // 嵌套条件 (复合字段用)
```

---

## 二、搜索参数格式 (SearchParams)

### 2.1 操作符枚举值

| 操作符 | 含义 | 适用类型 |
|--------|------|---------|
| `=` | 等于 | 文本/数字/选项 |
| `!=` | 不等于 | 文本/数字/选项 |
| `~` | 匹配 (模糊) | 文本/链接 |
| `!~` | 不匹配 | 文本/链接 |
| `<` / `>` / `<=` / `>=` | 比较 | 数字/日期/时间 |
| `HAS ANY OF` | 存在选项属于 | 多选项/人员 |
| `HAS NONE OF` | 不存在选项属于 | 多选项/人员 |
| `CONTAINS` | 包含 | 多选项/多选人员 |
| `NOT CONTAINS` | 不包含 | 多选项/多选人员 |
| `IS NULL` | 为空 | 所有类型 |
| `IS NOT NULL` | 不为空 | 所有类型 |
| `MEET` | 满足 (复合字段) | 复合字段 |
| `NOT MEET` | 不满足 (复合字段) | 复合字段 |

### 2.2 前置操作符 (复合字段专用)

| 前置操作符 | 含义 |
|-----------|------|
| `EVERY` | 每一组 |
| `ANY` | 存在一组 |

### 2.3 固定参数及取值

| 参数名 | param_key | 支持的操作符 | value类型 |
|--------|-----------|-------------|---------|
| 进行中节点 | `current_nodes` | `=`, `!=`, `HAS ANY OF`, `HAS NONE OF`, `CONTAINS`, `NOT CONTAINS` | `list<string>` |
| 流程节点 | `all_states` | `=`, `!=`, `HAS ANY OF`, `HAS NONE OF`, `CONTAINS`, `NOT CONTAINS` | `list<string>` |
| 流程节点时间 | `feature_state_time` | `<`, `>`, `<=`, `>=`, `IS NULL`, `IS NOT NULL` | `object` |
| 全部人员 | `people` | `=`, `!=`, `HAS ANY OF`, `HAS NONE OF`, `IS NULL`, `IS NOT NULL`, `CONTAINS`, `NOT CONTAINS` | `list<string>` |
| 创建时间 | `created_at` | `=`, `!=`, `<`, `>`, `<=`, `>=` | `int64` (毫秒时间戳) |
| 节点负责人 | `node_owners` | `=`, `!=`, `HAS ANY OF`, `HAS NONE OF`, `IS NULL`, `IS NOT NULL`, `CONTAINS`, `NOT CONTAINS` | `list<object>` |
| 工作项ID | `work_item_id` | `=`, `!=`, `<`, `>`, `<=`, `>=`, `HAS ANY OF`, `HAS NONE OF` | `list<int64>` |
| 工作项状态 | `work_item_status` | `=`, `!=`, `HAS ANY OF`, `HAS NONE OF` | `list<string>` |
| 模板ID | `template_id` | `=`, `!=`, `HAS ANY OF`, `HAS NONE OF`, `IS NULL`, `IS NOT NULL` | `list<int64>` |
| 业务线 | `business` | `=`, `!=`, `HAS ANY OF`, `HAS NONE OF`, `IS NULL`, `IS NOT NULL` | `list<string>` |
| 角色人员 | `role_owners` | `=`, `!=`, `HAS ANY OF`, `HAS NONE OF`, `IS NULL`, `IS NOT NULL`, `CONTAINS`, `NOT CONTAINS` | `list<object>` |

### 2.4 常用查询示例

#### 查询指定时间范围和状态的工作项

```json
{
  "page_size": 50,
  "page_num": 1,
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
  }
}
```

#### 查询指定人员的所有工作项

```json
{
  "search_group": {
    "conjunction": "AND",
    "search_params": [
      {
        "param_key": "people",
        "value": ["user_key_xxx"],
        "operator": "HAS ANY OF"
      }
    ]
  }
}
```

#### 查询多个角色都包含指定人员的工作项

```json
{
  "search_group": {
    "conjunction": "AND",
    "search_params": [
      {
        "param_key": "role_owners",
        "value": [
          {"role": "RD", "owners": ["user_key_1"]},
          {"role": "QA", "owners": ["user_key_2"]}
        ],
        "operator": "HAS ANY OF"
      }
    ]
  }
}
```

#### 复合字段筛选 - 每一组子字段都满足条件

```json
{
  "search_group": {
    "conjunction": "AND",
    "search_params": [
      {
        "param_key": "field_ba8f2e",      // 复合字段-父字段
        "operator": "MEET",
        "pre_operator": "EVERY",
        "value_search_groups": {
          "conjunction": "AND",
          "search_params": [
            {
              "param_key": "field_4396f8", // 复合字段-子字段
              "operator": "=",
              "value": "1"
            }
          ]
        }
      }
    ]
  }
}
```

---

## 三、字段与属性 (Fields & Properties)

> **属性 (Properties)**: 工作项的原始系统信息，不可自定义
> **字段 (Fields)**: 可配置的自定义或系统字段

### 3.1 属性 (Properties) - 系统内置

| 属性名 | API Key | 入参 | 出参类型 | 说明 |
|--------|---------|------|---------|------|
| 工作项ID | `id` | 不支持 | `int64` | 唯一标识 |
| 工作项名称 | `name` | `string` | `string` | - |
| 工作项状态 | `work_item_status` | `struct` | `struct` | 含 state_key, is_archived_state |
| 工作项模板ID | `template_id` | `int64` | `int64` | - |
| 工作项模式 | `pattern` | 不支持 | `string` | Node/State |
| 当前节点 | `current_nodes` | 系统计算 | `array` | 含 id, name, owners |
| 创建者 | `created_by` | 系统计算 | `string` | user_key |
| 更新者 | `updated_by` | 系统计算 | `string` | user_key |
| 创建时间 | `created_at` | 系统计算 | `int64` | 毫秒时间戳 |
| 更新时间 | `updated_at` | 系统计算 | `int64` | 毫秒时间戳 |
| 所属空间 | `project_key` | 系统计算 | `string` | - |

### 3.2 字段 (Fields) - 自定义类型

#### 文本类

| 字段类型 | field_type_key | 入参类型 | 出参类型 | 支持操作符 |
|---------|----------------|---------|---------|-----------|
| 单行文本 | `text` | `string` | `string` | `~`, `!~`, `=`, `!=`, `IS NULL`, `IS NOT NULL` |
| 多行文本 | `text` | `string` | `string` | `~`, `!~`, `=`, `!=`, `IS NULL`, `IS NOT NULL` |
| 富文本 | `multi_text` | `array` | `string/array` | `~`, `!~`, `IS NULL`, `IS NOT NULL` |
| URL链接 | `link` | `string` | `string` | `~`, `!~`, `=`, `!=`, `IS NULL`, `IS NOT NULL` |
| 电话号码 | `telephone` | `string` | `string` | `~`, `!~`, `=`, `!=`, `IS NULL`, `IS NOT NULL` |
| 电子邮件 | `email` | `string` | `string` | `~`, `!~`, `=`, `!=`, `IS NULL`, `IS NOT NULL` |

#### 选项类

| 字段类型 | field_type_key | 入参类型 | 出参类型 | 支持操作符 |
|---------|----------------|---------|---------|-----------|
| 单选 | `select` | `struct` | `struct` | `=`, `!=`, `HAS ANY OF`, `HAS NONE OF`, `IS NULL`, `IS NOT NULL` |
| 多选 | `multi_select` | `array` | `array` | `=`, `!=`, `HAS ANY OF`, `HAS NONE OF`, `IS NULL`, `IS NOT NULL`, `CONTAINS`, `NOT CONTAINS` |
| 单选按钮 | `radio` | `struct` | `struct` | `=`, `!=`, `HAS ANY OF`, `HAS NONE OF`, `IS NULL`, `IS NOT NULL` |
| 级联单选 | `tree_select` | `struct` | `struct` | `=`, `!=`, `HAS ANY OF`, `HAS NONE OF`, `IS NULL`, `IS NOT NULL` |
| 级联多选 | `tree_multi_select` | `array` | `array` | `=`, `!=`, `HAS ANY OF`, `HAS NONE OF`, `IS NULL`, `IS NOT NULL` |

#### 人员类

| 字段类型 | field_type_key | 入参类型 | 出参类型 | 支持操作符 |
|---------|----------------|---------|---------|-----------|
| 单选人员 | `user` | `string` | `string` | `=`, `!=`, `HAS ANY OF`, `HAS NONE OF`, `IS NULL`, `IS NOT NULL` |
| 多选人员 | `multi_user` | `array<string>` | `array<string>` | `=`, `!=`, `HAS ANY OF`, `HAS NONE OF`, `IS NULL`, `IS NOT NULL`, `CONTAINS`, `NOT CONTAINS` |

#### 日期时间类

| 字段类型 | field_type_key | 入参类型 | 出参类型 | 支持操作符 |
|---------|----------------|---------|---------|-----------|
| 日期 | `date` | `int64` | `int64` | `<`, `>`, `<=`, `>=`, `IS NULL`, `IS NOT NULL` |
| 日期时间 | `date` (秒精度) | `int64` | `int64` | `<`, `>`, `<=`, `>=`, `IS NULL`, `IS NOT NULL` |
| 日期区间 | `schedule` | `struct` | `struct` | - |

#### 关联类

| 字段类型 | field_type_key | 入参类型 | 出参类型 | 支持操作符 |
|---------|----------------|---------|---------|-----------|
| 单选关联工作项 | `work_item_related_select` | `int64` | `int64` | `=`, `!=`, `HAS ANY OF`, `HAS NONE OF`, `IS NULL`, `IS NOT NULL` |
| 多选关联工作项 | `work_item_related_multi_select` | `array<int64>` | `array<int64>` | `=`, `!=`, `HAS ANY OF`, `HAS NONE OF`, `IS NULL`, `IS NOT NULL`, `CONTAINS`, `NOT CONTAINS` |

#### 其他类型

| 字段类型 | field_type_key | 入参类型 | 出参类型 | 支持操作符 |
|---------|----------------|---------|---------|-----------|
| 数字 | `number` | `float` | `float` | `=`, `!=`, `<`, `>`, `<=`, `>=`, `IS NULL`, `IS NOT NULL` |
| 开关 | `bool` | `bool` | `bool` | `=`, `!=`, `IS NULL`, `IS NOT NULL` |
| 系统外信号 | `signal` | `bool/null` | `bool/null` | `=`, `!=`, `HAS ANY OF`, `HAS NONE OF` |
| 附件 | `multi_file` | 不支持 | `array` | - |
| 复合字段 | `compound_field` | `array` | `array` | `IS NULL`, `IS NOT NULL`, `MEET`, `NOT MEET` |

### 3.3 字段值格式示例

#### 单选字段

```json
{
  "field_key": "field_312244",
  "field_type_key": "select",
  "field_value": {
    "label": "选项1",
    "value": "8lheuaepp"
  }
}
```

#### 多选人员字段

```json
{
  "field_key": "field_8b18fd",
  "field_type_key": "multi_user",
  "field_value": [
    "user_key_1",
    "user_key_2"
  ]
}
```

#### 复合字段 (多行数据)

```json
{
  "field_key": "field_5a711c",
  "field_type_key": "compound_field",
  "field_value": [
    [
      {"field_key": "field_text1", "field_value": "内容1", "field_type_key": "text"},
      {"field_key": "field_date1", "field_value": 1722355200000, "field_type_key": "date"}
    ],
    [
      {"field_key": "field_text1", "field_value": "内容2", "field_type_key": "text"},
      {"field_key": "field_date1", "field_value": 1722441600000, "field_type_key": "date"}
    ]
  ]
}
```

### 3.4 富文本格式

#### 结构

```json
[
  {
    "type": "paragraph",           // paragraph/blank/checklist/ul/ol/horizontalLine/table
    "content": [
      {
        "type": "text",            // text/hyperlink/img/linkPreview
        "text": "文本内容",
        "attrs": {
          "fontColor": "blue",
          "bold": "true",
          "italic": "true",
          "underline": "true",
          "fontSize": "h1"
        }
      }
    ],
    "lineAttrs": {
      "align": "left",
      "indent": "1",
      "blockquote": "true"
    }
  }
]
```

#### 颜色值

| 颜色名 | RGB值 |
|-------|-------|
| black | rgb(26, 26, 26) |
| white | rgb(255, 255, 255) |
| blue | rgb(76, 136, 255) |
| green | rgb(84, 194, 72) |
| purple | rgb(127, 59, 245) |
| yellow | rgb(250, 200, 35) |
| red | rgb(240, 91, 86) |

---

## 四、MQL 查询语法

### 4.1 基础语法

```sql
SELECT fieldList
FROM objectType
WHERE conditionExpression
[ORDER BY fieldOrderByList [{ASC|DESC}]]
[LIMIT [offset,] row_count]
```

### 4.2 数据类型

| 类型 | 说明 |
|------|------|
| `bool` | TRUE/FALSE/1/0 |
| `bigint` | 整数 |
| `double` | 浮点数 |
| `varchar` | 字符串 |
| `date` | 日期 (YYYY-MM-DD) |
| `datetime` | 日期时间 (ISO8601) |
| `array` | 数组 |
| `lambda expression` | Lambda表达式 |

### 4.3 常用函数

| 函数 | 说明 | 示例 |
|------|------|------|
| `all_match(array, predicate)` | 所有元素满足条件 | `all_match(ary_col, x -> x > 10)` |
| `any_match(array, predicate)` | 存在元素满足条件 | `any_match(ary_col, x -> x = '用户A')` |
| `none_match(array, predicate)` | 所有元素都不满足 | `none_match(ary_col, x -> x = 'A')` |
| `array_contains(array, element)` | 数组包含元素 | `array_contains(\`优先级\`, 'P0')` |
| `array_cardinality(array)` | 数组元素个数 | `array_cardinality(\`负责人\`)` |
| `current_login_user()` | 当前登录用户 | `current_login_user()` |
| `team(include_manager, name)` | 团队成员 | `team(true, '开发团队')` |
| `RELATIVE_DATETIME_EQ(col, date_para)` | 相对时间等于 | `RELATIVE_DATETIME_EQ(\`创建时间\`, 'today')` |
| `RELATIVE_DATETIME_BETWEEN(col, date_para)` | 相对时间区间 | `RELATIVE_DATETIME_BETWEEN(\`创建时间\`, 'past', '30d')` |

### 4.4 相对时间参数

| date_para | 说明 | 可配days |
|-----------|------|---------|
| `today` | 今天 | 是 (如 '3d', '-3d') |
| `tomorrow` | 明天 | 否 |
| `yesterday` | 昨天 | 否 |
| `current_week` | 本周 | 否 |
| `next_week` | 下周 | 否 |
| `last_week` | 上周 | 否 |
| `current_month` | 本月 | 否 |
| `next_month` | 下月 | 否 |
| `last_month` | 上月 | 否 |
| `future` | 未来 | 是 |
| `past` | 过去 | 是 |

### 4.5 MQL 示例

#### 查询需求

```sql
SELECT `工作项id`, `优先级`
FROM `空间x`.`需求`
WHERE array_contains(`负责人`, '用户1')
  AND `创建时间` between '2025-01-01' and '2025-10-01'
  AND `优先级` = 'P0'
  AND `是否完结` = true
ORDER BY `优先级` DESC
LIMIT 100
```

#### 查询缺陷

```sql
SELECT `工作项id`
FROM `空间x`.`缺陷`
WHERE array_contains(`__RD`, '用户1')
  AND RELATIVE_DATETIME_EQ(`创建时间`, 'yesterday')
  AND any_match(`处理人`, x -> x in (team(true, '开放平台团队')))
  AND `__开发周期_开始时间` > '2025-01-01'
  AND `__开发周期_结束时间` < '2025-01-31'
```

---

## 五、错误码参考

### 5.1 HTTP 4xx 错误码

| HTTP状态码 | 错误码 | 错误信息 | 说明 |
|-----------|-------|---------|------|
| 400 | 20001 | Param Request Limit | 参数超过限制 (如空间ID>100) |
| 400 | 20002 | Page Size Limit | page_size超过200 |
| 400 | 20003 | Wrong WorkItemType Param | 工作项类型参数错误 |
| 400 | 20005 | Missing Param | 必填参数缺失 |
| 400 | 20006 | Invalid Param | 参数不合法 |
| 400 | 20029 | Unsupported Field Type | 不支持的字段类型 |
| 400 | 20041 | Field RoleOwner Must Be Set | 必须传入role_owners |
| 400 | 20055 | Search Result > 2000 | 搜索结果超限 |
| 400 | 20063 | Search Operator Error | 操作符不支持 |
| 400 | 20072 | Conjunction Value Error | conjunction仅支持AND/OR |
| 403 | 10001 | No Permission | 无操作权限 |
| 403 | 10002 | Illegal Operation | 非法操作 |
| 403 | 10301 | Check Token Perm Failed | token权限校验失败 |
| 404 | 13001 | View not exist | 视图不存在 |
| 404 | 30005 | WorkItem Not Found | 工作项未找到 |
| 404 | 30009 | Field Not Found | 字段未找到 |
| 429 | 10429 | API Frequency Limit | 请求超限 (15qps) |

### 5.2 常见场景错误

| 场景 | 错误码 | 说明 |
|------|-------|------|
| 更新评论 | 10001 | 非评论创建人无法更新 |
| 删除视图 | 10001 | 无编辑视图权限 |
| 节点流转 | 10002 | 当前节点不能流转到指定节点 |
| 创建子任务 | 10001 | 空间不匹配 |
| 字段更新 | 20029 | 不支持更新的字段类型 |
| 字段选项 | 20050 | 字段选项值错误 |
| 工作项终止 | 20056 | 只有工作流模式可终止 |

---

## 六、关系图谱

### 6.1 数据结构关系图

```
Project (空间)
├── Business (业务线)
│   ├── role_owners
│   └── labels
├── WorkItemType (工作项类型)
│   ├── story (需求)
│   ├── issue (缺陷)
│   └── custom (自定义)
├── WorkItem (工作项)
│   ├── id / name / status
│   ├── fields[]
│   │   ├── text / number / select
│   │   ├── user / multi_user
│   │   ├── date / schedule
│   │   └── compound_field
│   ├── comments[]
│   ├── sub_tasks[]
│   └── linked_work_items[]
├── Template (模板)
│   └── workflow
└── View (视图)
    └── filters
```

### 6.2 搜索参数层级关系

```
SearchRequest
└── search_group: SearchGroup
    ├── conjunction: "AND" | "OR"
    ├── search_params: SearchParam[]
    │   ├── param_key (字段key或固定key)
    │   ├── operator
    │   ├── value
    │   ├── pre_operator (复合字段)
    │   └── value_search_groups (复合字段嵌套)
    └── search_groups: SearchGroup[] (嵌套组)
        └── ... 递归
```

### 6.3 字段类型与操作符映射

```
┌─────────────────────────────────────────────────────────┐
│                    字段类型                              │
├─────────────────┬───────────────────────────────────────┤
│  文本类          │  ~, !~, =, !=, IS NULL, IS NOT NULL   │
│  (text, link)   │                                       │
├─────────────────┼───────────────────────────────────────┤
│  数字类          │  =, !=, <, >, <=, >=, IS NULL        │
│  (number)       │                                       │
├─────────────────┼───────────────────────────────────────┤
│  选项类          │  =, !=, HAS ANY OF, HAS NONE OF,      │
│  (select, radio)│  IS NULL, IS NOT NULL                 │
├─────────────────┼───────────────────────────────────────┤
│  多选项类        │  =, !=, HAS ANY OF, HAS NONE OF,      │
│  (multi-select) │  CONTAINS, NOT CONTAINS, IS NULL       │
├─────────────────┼───────────────────────────────────────┤
│  人员类          │  =, !=, HAS ANY OF, HAS NONE OF,       │
│  (user)         │  IS NULL, IS NOT NULL                 │
├─────────────────┼───────────────────────────────────────┤
│  日期类          │  <, >, <=, >=, IS NULL, IS NOT NULL   │
│  (date)         │                                       │
├─────────────────┼───────────────────────────────────────┤
│  复合字段        │  IS NULL, IS NOT NULL, MEET, NOT MEET  │
│  (compound)     │                                       │
└─────────────────┴───────────────────────────────────────┘
```

---

## 相关文档

本知识库整合了以下飞书项目开发者手册内容：

### 搜索相关
- [[全量搜索参数格式及常用示例 - 开发者手册 - 飞书项目帮助中心]] - 跨空间搜索参数
- [[搜索参数格式及常用示例 - 开发者手册 - 飞书项目帮助中心]] - 搜索参数详解
- [[字段与属性解析格式 - 开发者手册 - 飞书项目帮助中心]] - 字段类型与操作符

### 查询语法
- [[MQL 语法说明 - 开发者手册 - 飞书项目帮助中心]] - MQL 查询语言

### 数据与错误
- [[数据结构汇总 - 开发者手册 - 飞书项目帮助中心]] - 数据结构定义
- [[Open API 错误码 - 开发者手册 - 飞书项目帮助中心]] - 错误码参考

### API 总览
- [[飞书项目OpenAPI完整API列表]] - 完整 API 列表 (82 个)
- [[API 列表 - 开发者手册 - 飞书项目帮助中心]] - 官方 API 列表

### 图谱文件
- [[飞书项目/飞书项目API结构图谱.canvas]] - 数据结构关系图
- [[飞书项目/飞书项目字段与操作符图谱.canvas]] - 字段类型映射图

---

## 标签

#飞书项目 #OpenAPI #API列表 #开发手册 #知识库 #SearchParams #MQL #错误码 #字段类型 #搜索参数

---

> **最后更新**: 2026-03-19
> **来源**: 飞书项目开发者手册
