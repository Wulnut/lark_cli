---
title: 飞书项目OpenAPI-工作项实例读写
tags: [飞书项目, OpenAPI]
category: 飞书项目
created: 2026-03-17
updated: 2026-03-17

---


> 文档编号：2
> 更新时间：2026-03-17
---
## 工作项实例读写 API 列表
| # | API | 方法 | 说明 |
|---|-----|------|------|
| 1 | 获取工作项详情 | POST | `/open_api/:project_key/work_item/:work_item_type_key/query` |
| 2 | 获取创建工作项元数据 | GET | `/open_api/:project_key/work_item/:work_item_type_key/meta` |
| 3 | 创建工作项 | POST | `/open_api/:project_key/work_item/create` |
| 4 | 更新工作项 | PUT | `/open_api/:project_key/work_item/:work_item_type_key/:work_item_id` |
| 5 | 删除工作项 | DELETE | `/open_api/:project_key/work_item/:work_item_type_key/:work_item_id` |
| 6 | 终止/恢复工作项 | PUT | `/open_api/:project_key/work_item/:work_item_type_key/:work_item_id/abort` |
| 7 | 获取工作项操作记录 | POST | `/open_api/op_record/work_item/list` |
| 8 | 批量查询评审意见、评审结论 | POST | `/open_api/work_item/finished/batch_query` |
| 9 | 修改评审结论和评审意见 | POST | `/open_api/work_item/finished/update` |
| 10 | 评审结论标签值查询 | POST | `/open_api/work_item/finished/query_conclusion_option` |
| 11 | 获取工作项的工时登记记录列表 | POST | `/open_api/work_item/man_hour/records` |
| 12 | 新增工时登记记录 | POST | `/open_api/:project_key/work_item/:work_item_type_key/:work_item_id/work_hour_record` |
| 13 | 更新工时登记记录 | PUT | `/open_api/:project_key/work_item/:work_item_type_key/:work_item_id/work_hour_record` |
| 14 | 删除工时登记记录 | DELETE | `/open_api/:project_key/work_item/:work_item_type_key/:work_item_id/work_hour_record` |
| 15 | 冻结/解冻工作项 | PUT | `/open_api/work_item/freeze` |
| 16 | 交付物信息批量查询（WBS） | POST | `/open_api/work_item/deliverable/batch_query` |
---
## 1. 获取工作项详情
**API**：`POST /open_api/:project_key/work_item/:work_item_type_key/query`
**说明**：获取指定空间和工作项类型下的一个工作项实例的详细信息
### 请求参数
| 参数 | 类型 | 必须 | 说明 |
|------|------|------|------|
| `project_key` | string | ✅ | 路径参数：空间 key |
| `work_item_type_key` | string | ✅ | 路径参数：工作项类型 key |
| `work_item_id` | int64 | ✅ | Body参数：工作项 ID |
| `select_all` | bool | 可选 | 是否返回所有字段 |
| `field_keys` | list\<string\> | 可选 | 返回的字段列表 |
### 请求示例
```json
{
  "work_item_id": 12345
}
```
### 返回参数
| 参数 | 类型 | 说明 |
|------|------|------|
| `work_item_id` | int64 | 工作项 ID |
| `work_item_type_key` | string | 工作项类型 |
| `name` | string | 工作项名称 |
| `status` | string | 状态 |
| `fields` | object | 字段值 |
---
## 2. 获取创建工作项元数据
**API**：`GET /open_api/:project_key/work_item/:work_item_type_key/meta`
**说明**：获取指定工作项类型的"元数据"，它是创建一个工作项实例的最小数据单元
### 返回参数
| 参数 | 类型 | 说明 |
|------|------|------|
| `required_fields` | list | 必填字段 |
| `optional_fields` | list | 可选字段 |
| `field_configs` | list | 字段配置详情 |
---
## 3. 创建工作项
**API**：`POST /open_api/:project_key/work_item/create`
**说明**：在指定空间和工作项类型下新增一个"工作项实例"
### 请求参数
| 参数 | 类型 | 必须 | 说明 |
|------|------|------|------|
| `project_key` | string | ✅ | 空间 key |
| `work_item_type_key` | string | ✅ | 工作项类型 key |
| `name` | string | ✅ | 工作项名称 |
| `fields` | object | 可选 | 字段值，key-value 格式 |
| `role_owners` | list\<object\> | ✅ | 角色人员配置 |
### 请求示例
```json
{
  "project_key": "空间key",
  "work_item_type_key": "story",
  "name": "新建工作项",
  "fields": {
    "field_key_1": "值1",
    "field_key_2": ["选项1", "选项2"]
  },
  "role_owners": [
    {
      "role": "owner",
      "owner_ids": ["user_key_1"]
    }
  ]
}
```
---
## 4. 更新工作项
**API**：`PUT /open_api/:project_key/work_item/:work_item_type_key/:work_item_id`
**说明**：修改指定空间和工作项类型下的一个"工作项实例"
### 请求参数
| 参数 | 类型 | 必须 | 说明 |
|------|------|------|------|
| `work_item_id` | int64 | ✅ | 工作项 ID |
| `name` | string | 可选 | 工作项名称 |
| `fields` | object | 可选 | 字段值 |
| `role_owners` | list\<object\> | 可选 | 角色人员 |
### 请求示例
```json
{
  "work_item_id": 12345,
  "name": "更新后的名称",
  "fields": {
    "field_key_1": "新值"
  }
}
```
---
## 5. 删除工作项
**API**：`DELETE /open_api/:project_key/work_item/:work_item_type_key/:work_item_id`
**说明**：删除指定空间和工作项类型下的一个"工作项实例"
### 请求参数
| 参数 | 类型 | 必须 | 说明 |
|------|------|------|------|
| `project_key` | string | ✅ | 空间 key（路径） |
| `work_item_type_key` | string | ✅ | 工作项类型 key（路径） |
| `work_item_id` | int64 | ✅ | 工作项 ID（路径） |
### 返回参数
| 参数 | 类型 | 说明 |
|------|------|------|
| `err` | object | 错误信息 |
| `err_code` | int | 0=成功 |
---
## 6. 终止/恢复工作项
**API**：`PUT /open_api/:project_key/work_item/:work_item_type_key/:work_item_id/abort`
**说明**：用于终止或者恢复指定空间和工作项类型下的一个"工作项实例"
### 请求参数
| 参数 | 类型 | 必须 | 说明 |
|------|------|------|------|
| `operation` | string | ✅ | 操作类型：`abort`（终止）或 `restore`（恢复） |
| `reason` | string | ✅ | 原因 |
### 请求示例
```json
{
  "operation": "abort",
  "reason": "项目取消"
}
```
---
## 7. 获取工作项操作记录
**API**：`POST /open_api/op_record/work_item/list`
**说明**：用于获取指定空间下的多个工作项实例的操作记录
### 请求参数
| 参数 | 类型 | 必须 | 说明 |
|------|------|------|------|
| `project_key` | string | ✅ | 空间 key |
| `work_item_ids` | list\<int64\> | ✅ | 工作项 ID 列表 |
| `start_time` | int64 | 可选 | 开始时间（毫秒） |
| `end_time` | int64 | 可选 | 结束时间（毫秒） |
---
## 8-10. 评审相关 API
### 批量查询评审意见、评审结论
**API**：`POST /open_api/work_item/finished/batch_query`
### 修改评审结论和评审意见
**API**：`POST /open_api/work_item/finished/update`
### 评审结论标签值查询
**API**：`POST /open_api/work_item/finished/query_conclusion_option`
---
## 11-14. 工时登记 API
### 获取工作项的工时登记记录列表
**API**：`POST /open_api/work_item/man_hour/records`
### 新增工时登记记录
**API**：`POST /open_api/:project_key/work_item/:work_item_type_key/:work_item_id/work_hour_record`
**请求参数**：
| 参数 | 类型 | 必须 | 说明 |
|------|------|------|------|
| `work_item_id` | int64 | ✅ | 工作项 ID |
| `user_key` | string | ✅ | 登记人 |
| `hours` | float | ✅ | 工时数 |
| `record_date` | int64 | ✅ | 登记日期 |
| `description` | string | 可选 | 描述 |
### 更新工时登记记录
**API**：`PUT /open_api/:project_key/work_item/:work_item_type_key/:work_item_id/work_hour_record`
### 删除工时登记记录
**API**：`DELETE /open_api/:project_key/work_item/:work_item_type_key/:work_item_id/work_hour_record`
---
## 15. 冻结/解冻工作项
**API**：`PUT /open_api/work_item/freeze`
**请求参数**：
| 参数 | 类型 | 必须 | 说明 |
|------|------|------|------|
| `work_item_ids` | list\<int64\> | ✅ | 工作项 ID 列表 |
| `freeze` | bool | ✅ | true=冻结，false=解冻 |
---
## 16. 交付物信息批量查询（WBS）
**API**：`POST /open_api/work_item/deliverable/batch_query`
**说明**：用于查询交付物详情信息
---
## 相关文档
- [工作项实例搜索](飞书项目OpenAPI-工作项实例搜索.md)
- [工作项流程与节点](飞书项目OpenAPI-工作项流程与节点.md)
- [搜索参数格式](飞书项目OpenAPI搜索参数格式及常用示例.md)
---
## 🏷️ 标签
#飞书项目 #OpenAPI #工作项 #CRUD #API文档