---
title: 飞书项目OpenAPI-子任务
tags: [飞书项目, OpenAPI]
category: 飞书项目
created: 2026-03-17
updated: 2026-03-17

---


> 文档编号：4
> 更新时间：2026-03-17
---
## 子任务 API 列表
| # | API | 方法 | 说明 |
|---|-----|------|------|
| 1 | 获取指定的子任务列表（跨空间） | POST | `/open_api/work_item/subtask/search` |
| 2 | 获取子任务详情 | GET | `/open_api/:project_key/work_item/:work_item_type_key/:work_item_id/workflow/task` |
| 3 | 创建子任务 | POST | `/open_api/:project_key/work_item/:work_item_type_key/:work_item_id/workflow/task` |
| 4 | 更新子任务 | POST | `/open_api/:project_key/work_item/:work_item_type_key/:work_item_id/workflow/:node_id/task/:task_id` |
| 5 | 子任务完成/回滚 | POST | `/open_api/:project_key/work_item/:work_item_type_key/:work_item_id/subtask/modify` |
| 6 | 删除子任务 | DELETE | `/open_api/:project_key/work_item/:work_item_type_key/:work_item_id/task/:task_id` |
---
## 1. 获取指定的子任务列表（跨空间）
**API**：`POST /open_api/work_item/subtask/search`
**说明**：用于跨空间搜索符合传入条件的子任务
### 请求参数
| 参数 | 类型 | 必须 | 说明 |
|------|------|------|------|
| `project_keys` | list\<string\> | 可选 | 空间 key 列表 |
| `work_item_type_keys` | list\<string\> | 可选 | 工作项类型 key 列表 |
| `work_item_ids` | list\<int64\> | 可选 | 工作项 ID 列表 |
| `task_names` | list\<string\> | 可选 | 子任务名称 |
| `task_status` | list\<string\> | 可选 | 子任务状态 |
| `owner_ids` | list\<string\> | 可选 | 负责人 |
| `pagination` | object | 可选 | 分页信息 |
### 请求示例
```json
{
  "project_keys": ["space_key1"],
  "work_item_type_keys": ["story"],
  "task_status": ["pending"],
  "pagination": {
    "page_size": 50,
    "page_num": 1
  }
}
```
---
## 2. 获取子任务详情
**API**：`GET /open_api/:project_key/work_item/:work_item_type_key/:work_item_id/workflow/task`
**说明**：获取指定工作项实例上的子任务详细信息
### 请求参数
| 参数 | 类型 | 必须 | 说明 |
|------|------|------|------|
| `project_key` | string | ✅ | 空间 key（路径） |
| `work_item_type_key` | string | ✅ | 工作项类型 key（路径） |
| `work_item_id` | int64 | ✅ | 工作项 ID（路径） |
| `node_id` | string | 可选 | 节点 ID |
### 返回参数
| 参数 | 类型 | 说明 |
|------|------|------|
| `tasks` | list | 子任务列表 |
| `tasks[].task_id` | string | 子任务 ID |
| `tasks[].task_name` | string | 子任务名称 |
| `tasks[].status` | string | 状态 |
| `tasks[].owner_id` | string | 负责人 |
| `tasks[].start_time` | int64 | 开始时间 |
| `tasks[].end_time` | int64 | 结束时间 |
---
## 3. 创建子任务
**API**：`POST /open_api/:project_key/work_item/:work_item_type_key/:work_item_id/workflow/task`
**说明**：在一个工作项实例的指定节点上创建一个子任务
### 请求参数
| 参数 | 类型 | 必须 | 说明 |
|------|------|------|------|
| `task_name` | string | ✅ | 子任务名称 |
| `node_id` | string | ✅ | 节点 ID |
| `owner_ids` | list\<string\> | 可选 | 负责人 |
| `assignee_ids` | list\<string\> | 可选 | 被委托人次 |
| `role_assignee` | list\<object\> | 可选 | 角色负责人 |
| `start_time` | int64 | 可选 | 开始时间（毫秒） |
| `end_time` | int64 | 可选 | 结束时间（毫秒） |
| `description` | string | 可选 | 描述 |
### 请求示例
```json
{
  "task_name": "开发功能模块",
  "node_id": "node_1",
  "owner_ids": ["user_key_1"],
  "start_time": 1704067200000,
  "end_time": 1706659200000,
  "description": "完成用户登录功能"
}
```
---
## 4. 更新子任务
**API**：`POST /open_api/:project_key/work_item/:work_item_type_key/:work_item_id/workflow/:node_id/task/:task_id`
**说明**：更新工作项实例指定节点上的一个子任务详细信息
### 请求参数
| 参数 | 类型 | 必须 | 说明 |
|------|------|------|------|
| `task_name` | string | 可选 | 子任务名称 |
| `owner_ids` | list\<string\> | 可选 | 负责人 |
| `start_time` | int64 | 可选 | 开始时间 |
| `end_time` | int64 | 可选 | 结束时间 |
| `description` | string | 可选 | 描述 |
---
## 5. 子任务完成/回滚
**API**：`POST /open_api/:project_key/work_item/:work_item_type_key/:work_item_id/subtask/modify`
**说明**：用于完成或者回滚工作项实例指定节点上的一个子任务
### 请求参数
| 参数 | 类型 | 必须 | 说明 |
|------|------|------|------|
| `task_id` | string | ✅ | 子任务 ID |
| `node_id` | string | ✅ | 节点 ID |
| `operation` | string | ✅ | 操作类型：`confirm`（确认完成）/ `rollback`（回滚） |
### 请求示例
```json
{
  "task_id": "task_123",
  "node_id": "node_1",
  "operation": "confirm"
}
```
---
## 6. 删除子任务
**API**：`DELETE /open_api/:project_key/work_item/:work_item_type_key/:work_item_id/task/:task_id`
**说明**：用于删除指定工作项实例中的一个子任务
### 请求参数
| 参数 | 类型 | 必须 | 说明 |
|------|------|------|------|
| `project_key` | string | ✅ | 空间 key（路径） |
| `work_item_type_key` | string | ✅ | 工作项类型 key（路径） |
| `work_item_id` | int64 | ✅ | 工作项 ID（路径） |
| `task_id` | string | ✅ | 子任务 ID（路径） |
---
## 相关文档
- [工作项流程与节点](飞书项目OpenAPI-工作项流程与节点.md)
- [工作项实例读写](飞书项目OpenAPI-工作项实例读写.md)
---
## 🏷️ 标签
#飞书项目 #OpenAPI #子任务 #任务管理 #API文档