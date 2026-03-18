---
title: 飞书项目OpenAPI搜索参数格式及常用示例
tags: [飞书项目, OpenAPI]
category: 飞书项目
created: 2026-03-17
updated: 2026-03-17

---


> 文档：[搜索参数格式及常用示例](https://project.feishu.cn/b/helpcenter/1p8d7djs/1l8il0l6)
> 版本：2.0.0
> 更新时间：2026-03-17
---
## 📋 概述
本文档详细介绍飞书项目 OpenAPI 中搜索参数的格式规范，支持复杂的筛选条件组合。
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
## 📊 操作符枚举值
### 基础操作符 (16种)
| 编号 | 操作符 | 说明 | 适用字段类型 |
|------|--------|------|-------------|
| 1 | `~` | 匹配（模糊搜索） | 文本类 |
| 2 | `!~` | 不匹配 | 文本类 |
| 3 | `=` | 等于 | 所有类型 |
| 4 | `!=` | 不等于 | 所有类型 |
| 5 | `<` | 小于 | 数字、日期 |
| 6 | `>` | 大于 | 数字、日期 |
| 7 | `<=` | 小于等于 | 数字、日期 |
| 8 | `>=` | 大于等于 | 数字、日期 |
| 9 | `HAS ANY OF` | 存在选项属于 | 选项类、用户类 |
| 10 | `HAS NONE OF` | 全部选项均不属于 | 选项类、用户类 |
| 11 | `IS NULL` | 为空 | 所有类型 |
| 12 | `IS NOT NULL` | 不为空 | 所有类型 |
| 13 | `CONTAINS` | 包含 | 文本、多选类 |
| 14 | `NOT CONTAINS` | 不包含 | 文本、多选类 |
| 15 | `MEET` | 满足 | 复合字段 |
| 16 | `NOT MEET` | 不满足 | 复合字段 |
### 前置操作符
| 操作符 | 说明 |
|--------|------|
| `EVERY` | 每一组（复合字段） |
| `ANY` | 存在一组（复合字段） |
---
## 📌 固定参数 (Special Parameters)
固定参数是系统级的筛选条件，不依赖于自定义字段。
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
---
## 📝 自定义字段类型与操作符
### 字段类型与支持的操作符
| 字段类型 | field_type_key | 支持的操作符 |
|----------|----------------|--------------|
| 单行文本 | `text` | `~`, `!~`, `=`, `!=`, `IS NULL`, `IS NOT NULL` |
| 多行文本 | `multi_text` | `~`, `!~`, `IS NULL`, `IS NOT NULL` |
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
| 电话 | `telephone` | `~`, `!~`, `=`, `!=`, `IS NULL`, `IS NOT NULL` |
| 邮箱 | `email` | 同上 |
### 特殊字段类型
| 字段类型 | 说明 | 注意事项 |
|----------|------|----------|
| `signal` (系统外信号) | 用于筛选不在系统内的状态 | 值只能是 `undefined`、`null`、`true` |
| `compound_field` (复合字段) | 包含多个子字段的字段 | 需要配合 `value_search_groups` |
---
## 📊 时间筛选 special
### 相对时间参数
| 值 | 说明 |
|-----|------|
| `today` | 今天 |
| `yesterday` | 昨天 |
| `tomorrow` | 明天 |
| `current_week` | 当周 |
| `next_week` | 下周 |
| `last_week` | 上周 |
| `current_month` | 当月 |
| `next_month` | 下月 |
| `last_month` | 上月 |
| `future` | 未来 |
| `past` | 过去 |
### 时间筛选函数
| 函数 | 说明 |
|------|------|
| `RELATIVE_DATETIME_EQ(col, 'date_para')` | 等于相对时间 |
| `RELATIVE_DATETIME_GT(col, 'date_para')` | 大于相对时间 |
| `RELATIVE_DATETIME_BETWEEN(col, 'date_para')` | 属于相对时间范围 |
---
## 📝 常用查询示例
### 示例1：查询指定状态且创建时间在某个区间的工作项
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
**说明**：查询创建时间在 2022年6月1日 到 2022年6月2日 之间，且状态为"进行中"的工作项
---
### 示例2：查询包含指定人员的工作项
```json
{
  "search_group": {
    "conjunction": "AND",
    "search_params": [
      {
        "param_key": "people",
        "value": ["user_key_123"],
        "operator": "HAS ANY OF"
      }
    ]
  }
}
```
**说明**：查询包含指定人员参与的工作项（负责人、关注人、创建人等）
---
### 示例3：通过 ID 列表查询工作项
```json
{
  "search_group": {
    "conjunction": "AND",
    "search_params": [
      {
        "param_key": "work_item_id",
        "value": [12345, 45678, 90000],
        "operator": "HAS ANY OF"
      }
    ]
  }
}
```
**说明**：通过工作项 ID 列表批量查询工作项
---
### 示例4：查询关联的工作项
```json
{
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
**说明**：查询关联了指定工作项的工作项（通过关联字段）
---
### 示例5：查询多个状态的工作项
```json
{
  "search_group": {
    "conjunction": "AND",
    "search_params": [
      {
        "param_key": "work_item_status",
        "value": ["start", "pending"],
        "operator": "HAS ANY OF"
      }
    ]
  }
}
```
**说明**：查询状态为"进行中"或"待处理"的工作项
---
### 示例6：查询某个节点的工作项
```json
{
  "search_group": {
    "conjunction": "AND",
    "search_params": [
      {
        "param_key": "current_nodes",
        "value": ["开发中", "测试中"],
        "operator": "HAS ANY OF"
      }
    ]
  }
}
```
**说明**：查询当前处于"开发中"或"测试中"节点的工作项
---
### 示例7：使用 OR 组合多个条件
```json
{
  "search_group": {
    "conjunction": "OR",
    "search_params": [
      {
        "param_key": "priority",
        "value": ["P0"],
        "operator": "="
      },
      {
        "param_key": "priority",
        "value": ["P1"],
        "operator": "="
      }
    ]
  }
}
```
**说明**：查询优先级为 P0 或 P1 的工作项
---
### 示例8：复杂组合查询
```json
{
  "page_size": 20,
  "search_group": {
    "conjunction": "AND",
    "search_params": [
      {
        "param_key": "work_item_status",
        "value": ["start"],
        "operator": "="
      },
      {
        "param_key": "people",
        "value": ["user_key_1", "user_key_2"],
        "operator": "HAS ANY OF"
      },
      {
        "param_key": "created_at",
        "value": 1704067200000,  // 2024-01-01
        "operator": ">="
      },
      {
        "param_key": "priority",
        "value": ["P0", "P1"],
        "operator": "HAS ANY OF"
      }
    ]
  }
}
```
**说明**：查询同时满足以下条件的工作项：
- 状态为"进行中"
- 包含指定人员
- 创建时间在 2024-01-01 之后
- 优先级为 P0 或 P1
---
### 示例9：复合字段查询
```json
{
  "search_group": {
    "conjunction": "AND",
    "search_params": [
      {
        "param_key": "复合字段key",
        "value": "",
        "operator": "MEET",
        "pre_operator": "ANY",
        "value_search_groups": {
          "conjunction": "AND",
          "search_params": [
            {
              "param_key": "子字段key",
              "value": "值",
              "operator": "="
            }
          ]
        }
      }
    ]
  }
}
```
**说明**：复合字段需要使用 `MEET` 操作符配合 `value_search_groups`
---
### 示例10：模糊搜索文本字段
```json
{
  "search_group": {
    "conjunction": "AND",
    "search_params": [
      {
        "param_key": "title",
        "value": "bug",
        "operator": "~"
      }
    ]
  }
}
```
**说明**：搜索标题中包含"bug"的工作项（模糊匹配）
---
## ⚠️ 注意事项
### 1. 时间戳格式
- 必须是**毫秒时间戳**（13位数字）
- 例如：`1654064482000` = `2022-06-01 10:21:22`
### 2. 操作符选择
- 不同字段类型支持不同的操作符
- 使用不支持的操作符会返回错误
### 3. 数组参数
- 使用 `HAS ANY OF` 时，value 为数组
- 数组最多支持 50 个元素
### 4. 分页
- `page_size` 最大值：100
- `page_num` 从 1 开始
### 5. 搜索限制
- 查询结果最多返回 2000 条
- 超过时需要缩小筛选范围
---
## 🔧 常见问题
### Q1: 如何查询空值字段？
```json
{
  "param_key": "字段key",
  "value": null,
  "operator": "IS NULL"
}
```
### Q2: 如何查询非空字段？
```json
{
  "param_key": "字段key",
  "value": null,
  "operator": "IS NOT NULL"
}
```
### Q3: 如何查询日期区间？
```json
{
  "param_key": "created_at",
  "value": 开始时间戳,
  "operator": ">="
},
{
  "param_key": "created_at",
  "value": 结束时间戳,
  "operator": "<="
}
```
### Q4: 搜索结果为空怎么办？
1. 检查 `param_key` 是否正确（字段 key 不是字段名称）
2. 检查操作符是否匹配字段类型
3. 检查时间戳格式是否正确（13位毫秒）
---
## 📎 相关文档
- [字段与属性解析格式](飞书项目OpenAPI字段类型汇总.md)
- [全量搜索参数格式及常用示例](./飞书项目OpenAPI全量搜索参数格式.md)
- [Open API 错误码](飞书项目OpenAPI开发者手册汇总.md)
---
## 相关文档
- [[飞书项目OpenAPI字段类型汇总]] - 字段类型详解
- [[飞书项目OpenAPI全量搜索参数格式及常用示例]] - 跨空间搜索
- [[飞书项目OpenAPI开发者手册汇总]] - 错误码
- [[飞书项目OpenAPI完整API列表]] - 主索引
---
## 🏷️ 标签
#飞书项目 #OpenAPI #搜索参数 #SearchParam #API文档