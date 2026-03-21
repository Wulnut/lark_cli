---
title: 飞书项目OpenAPI开发者手册汇总
source: https://project.feishu.cn/b/helpcenter/1p8d7djs
author:
published: 2026-03-17
created: 2026-03-17
updated: 2026-03-17
description: 飞书项目OpenAPI开发者手册文档索引，汇总工作项API、字段属性、搜索参数、数据结构、MQL语法及错误码等官方文档
tags:
  - "飞书项目"
  - "OpenAPI"
  - "开发手册"
  - "知识库"
  - "索引"
category: "飞书项目"
related_docs:
  - "[[飞书项目API开发者知识库]]"
  - "[[飞书项目OpenAPI完整API列表]]"
---


> 汇总日期：2026-03-17
> 版本：2.0.0
> 来源：飞书项目开发者中心
---
## 📋 文档索引
| # | 文档 | 状态 |
|---|------|------|
| 1 | [工作项 API 汇总](https://project.feishu.cn/b/helpcenter/1p8d7djs/1tj6ggll) | ✅ |
| 2 | [字段与属性解析格式](https://project.feishu.cn/b/helpcenter/1p8d7djs/1tj6ggll#110a33af) | ✅ |
| 3 | [搜索参数格式及常用示例](https://project.feishu.cn/b/helpcenter/1p8d7djs/1l8il0l6) | ✅ |
| 4 | [全量搜索参数格式及常用示例](https://project.feishu.cn/b/helpcenter/1p8d7djs/w11hyb8w) | ✅ |
| 5 | [数据结构汇总](https://project.feishu.cn/b/helpcenter/1p8d7djs/1x1d372l) | ✅ |
| 6 | [MQL 语法说明](https://project.feishu.cn/b/helpcenter/1p8d7djs/tsl2uj3i) | ✅ 新发现 |
| 7 | [Open API 错误码](https://project.feishu.cn/b/helpcenter/1p8d7djs/5aueo3jr) | ✅ |
---
## 🔗 API 关联关系图
```
工作项 (Work Item)
      │
      ├── 工时 (Hours)
      ├── 资源库 (Resource)
      ├── 流程与节点 (Workflow)
      │      │
      │      ├── 评审 (Review)
      │      ├── 状态流转 (State)
      │      └── 子任务 (Subtask)
      │
      ├── 评论 (Comment)
      ├── 附件 (File)
      ├── 视图 (View)
      └── 度量 (Measure)
```
---
## 📦 主要 API 模块
### 1. 工作项 CRUD
| API | 方法 | 端点 |
|-----|------|------|
| 获取工作项列表 | GET | `/open_api/:project_key/work_item` |
| 获取工作项详情 | GET | `/open_api/:project_key/work_item/:work_item_id` |
| 创建工作项 | POST | `/open_api/:project_key/work_item` |
| 更新工作项 | PUT | `/open_api/:project_key/work_item/:work_item_id` |
| 删除工作项 | DELETE | `/open_api/:project_key/work_item/:work_item_type_key/:work_item_id` |
| 批量更新 | POST | `/open_api/task_result` |
| 冻结/解冻 | POST | `/open_api/work_item/freeze` |
| 终止/恢复 | POST | `/open_api/:project_key/work_item/:work_item_type_key/:work_item_id/abort` |
---
### 2. 工作项工时
| API | 说明 |
|-----|------|
| 获取工时列表 | `/work_item/:work_item_id/work_hour_record` |
| 新增工时 | POST 创建工时记录 |
| 更新工时 | PUT 修改工时 |
| 删除工时 | DELETE |
---
### 3. 流程与节点
| API | 说明 |
|-----|------|
| 获取工作流详情 | 获取工作流配置 |
| 获取工作流详情（WBS） | WBS 模式 |
| 更新节点/排期 | 修改节点信息 |
| 节点完成/回滚 | `/node/:node_id/operate` |
| 状态流转 | 工作项状态变更 |
| 获取必填信息 | 获取流转所需字段 |
---
### 4. 子任务
| API | 说明 |
|-----|------|
| 获取子任务列表 | 跨空间查询 |
| 获取子任务详情 | 获取详情 |
| 创建子任务 | 新建 |
| 更新子任务 | 修改 |
| 子任务完成/回滚 | 状态管理 |
| 删除子任务 | 删除 |
---
### 5. 附件
| API | 说明 |
|-----|------|
| 添加附件 | 上传文件 |
| 上传富文本图片 | 图片上传 |
| 下载附件 | 文件下载 |
| 删除附件 | 删除 |
---
### 6. 评论
| API | 方法 | 说明 |
|-----|------|------|
| 查询评论 | GET | 获取列表 |
| 添加评论 | POST | 新建 |
| 更新评论 | PUT | 修改 |
| 删除评论 | DELETE | 删除 |
---
### 7. 配置 API
#### 空间配置
- 获取空间下工作项类型
- 获取空间下业务线详情
#### 字段配置
- 获取字段信息
- 创建自定义字段
- 更新自定义字段
#### 流程配置
- 获取流程模板列表
- 获取流程模板详情
- 新增/更新/删除流程模板
#### 角色配置
- 创建/更新/删除流程角色
- 获取角色详情
---
## 📝 搜索参数格式 (SearchParam)
### SearchGroup 结构
| 字段 | 类型 | 说明 |
|------|------|------|
| search_params | list\<SearchParam\> | 筛选条件 |
| conjunction | string | AND/OR，对应"且"和"或" |
| search_groups | list\<SearchGroup\> | 筛选组合 |
### SearchParam 结构
| 字段 | 类型 | 说明 |
|------|------|------|
| param_key | string | 筛选参数（字段 key 或特殊参数） |
| value | interface{} | 搜索字段值 |
| operator | string | 操作符类型 |
| pre_operator | string | 前置操作符（复合字段用） |
| value_search_groups | SearchGroup | 复合字段子字段筛选 |
---
### 操作符枚举值 (18种)
| 编号 | 操作符 | 说明 |
|------|--------|------|
| 1 | `~` | 匹配 |
| 2 | `!~` | 不匹配 |
| 3 | `=` | 等于 |
| 4 | `!=` | 不等于 |
| 5 | `<` | 小于 |
| 6 | `>` | 大于 |
| 7 | `<=` | 小于等于 |
| 8 | `>=` | 大于等于 |
| 9 | `HAS ANY OF` | 存在选项属于 |
| 10 | `HAS NONE OF` | 全部选项均不属于 |
| 11 | `IS NULL` | 为空 |
| 12 | `IS NOT NULL` | 不为空 |
| 13 | `CONTAINS` | 包含 |
| 14 | `NOT CONTAINS` | 不包含 |
| 15 | `MEET` | 满足 |
| 16 | `NOT MEET` | 不满足 |
---
### 前置操作符枚举值
| 操作符 | 说明 |
|--------|------|
| `EVERY` | 每一组 |
| `ANY` | 存在一组 |
---
### 固定参数及取值
| 参数名 | param_key | 操作符 | 值类型 | 说明 |
|--------|-----------|--------|--------|------|
| 进行中节点 | `current_nodes` | `=`, `!=`, `HAS ANY OF`, `HAS NONE OF`, `CONTAINS`, `NOT CONTAINS` | list\<string\> | 节点名称列表 |
| 流程节点 | `all_states` | 同上 | list\<string\> | 节点名称列表 |
| 流程节点时间 | `feature_state_time` | `<`, `>`, `<=`, `>=`, `IS NULL`, `IS NOT NULL` | object | 含 state_name, state_timestamp, state_condition |
| 全部人员 | `people` | `=`, `!=`, `HAS ANY OF`, `HAS NONE OF`, `IS NULL`, `IS NOT NULL`, `CONTAINS`, `NOT CONTAINS` | list\<string\> | user_key 列表 |
| 创建时间 | `created_at` | `=`, `!=`, `<`, `>`, `<=`, `>=` | int64 | 毫秒时间戳 |
| 节点负责人 | `node_owners` | 同上 | list\<object\> | 含 state_name 和 owners |
| 工作项ID | `work_item_id` | `=`, `!=`, `<`, `>`, `<=`, `>=`, `HAS ANY OF`, `HAS NONE OF` | list\<int64\> | 工作项 ID 列表 |
| 工作项状态 | `work_item_status` | `=`, `!=`, `HAS ANY OF`, `HAS NONE OF` | list\<string\> | start/closed 等 |
| 模板ID | `template_id` | `=`, `!=`, `HAS ANY OF`, `HAS NONE OF`, `IS NULL`, `IS NOT NULL` | list\<int64\> | 模板 ID 列表 |
| 业务线 | `business` | 同上 | list\<string\> | 业务线 ID 列表 |
| 角色人员 | `role_owners` | `=`, `!=`, `HAS ANY OF`, `HAS NONE OF`, `IS NULL`, `IS NOT NULL`, `CONTAINS`, `NOT CONTAINS` | list\<object\> | 含 role 和 owners |
---
## 📝 字段类型汇总
### 自定义类型及取值
| 参数类型 | field_type_key | 操作符 | 值类型 | 说明 |
|----------|----------------|--------|--------|------|
| 单行文本/多行文本 | `text` | `~`, `!~`, `=`, `!=`, `IS NULL`, `IS NOT NULL` | string | 前后模糊匹配 |
| 数字 | `number` | `=`, `!=`, `<`, `>`, `<=`, `>=`, `IS NULL`, `IS NOT NULL` | float64 | |
| 链接 | `link` | `~`, `!~`, `=`, `!=`, `IS NULL`, `IS NOT NULL` | string | |
| 开关 | `bool` | `=`, `!=`, `IS NULL`, `IS NOT NULL` | bool | |
| 单值系统外信号 | `signal` | `=`, `!=`, `HAS ANY OF`, `HAS NONE OF` | list\<string\> | undefined/null/true/false |
| 单选 | `select` | `=`, `!=`, `HAS ANY OF`, `HAS NONE OF`, `IS NULL`, `IS NOT NULL` | list\<string\> | 选项 value 值 |
| 单选按钮 | `radio` | 同上 | list\<string\> | |
| 多选 | `multi-select` | 同上 + CONTAINS/NOT CONTAINS | list\<string\> | |
| 级联单选 | `tree-select` | `=`, `!=`, `HAS ANY OF`, `HAS NONE OF`, `IS NULL`, `IS NOT NULL` | list\<string\> | |
| 级联多选 | `tree-multi-select` | 同上 | list\<string\> | |
| 单选人员 | `user` | 同上 | list\<string\> | user_key |
| 多选人员 | `multi-user` | 同上 + CONTAINS/NOT CONTAINS | list\<string\> | user_key |
| 富文本 | `multi-text` | `~`, `!~`, `IS NULL`, `IS NOT NULL` | string | 无格式匹配 |
| 单选关联工作项 | `workitem_related_select` | `=`, `!=`, `HAS ANY OF`, `HAS NONE OF`, `IS NULL`, `IS NOT NULL` | list\<int64\> | 工作项 ID |
| 多选关联工作项 | `workitem_related_multi_select` | 同上 + CONTAINS/NOT CONTAINS | list\<int64\> | |
| 日期时间 | `precise_date` | `<`, `>`, `<=`, `>=`, `IS NULL`, `IS NOT NULL` | int64 | 毫秒时间戳 |
| 日期 | `date` | 同上 | int64 | |
| 复合字段-父字段 | `compound_field` | `IS NULL`, `IS NOT NULL`, `MEET`, `NOT MEET` | - | 需配合 value_search_groups |
| 电话号码 | `telephone` | `~`, `!~`, `=`, `!=`, `IS NULL`, `IS NOT NULL` | string | |
| 电子邮件 | `email` | 同上 | string | |
---
## 📌 常用查询示例
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
### 示例2：查询包含指定人员的工作项
```json
{
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
### 示例3：通过 ID 查询工作项
```json
{
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
---
## 📝 MQL 语法说明 (Meego Query Language)
MQL 是飞书项目的专用查询语言，兼容 SQL 语法。
### 基础语法规则
```sql
SELECT fieldList          -- 指定查询的字段列表
FROM objectType          -- 指定要查询的数据来源
WHERE conditionExpression -- 指定查询条件
[ORDER BY fieldOrderByList]  -- 排序（可选）
[LIMIT [offset,] row_count]  -- 限制返回行数
```
### 数据类型
| 数据类型 | 说明 |
|----------|------|
| `bool` | 布尔值：TRUE、FALSE、1、0 |
| `bigint` | 整数类型 |
| `double` | 浮点数类型 |
| `varchar` | 字符串类型 |
| `date` | 日期格式：YYYY-MM-DD |
| `datetime` | 日期时间格式：ISO8601 |
| `array` | 数组类型 |
| `lambda expression` | Lambda 表达式 |
### 支持函数
| 函数 | 说明 |
|------|------|
| `all_match(array, predicate)` | 判断数组是否所有元素都满足条件 |
| `any_match(array, predicate)` | 判断数组是否有一个元素满足条件 |
| `none_match(array, predicate)` | 判断数组是否所有元素都不满足条件 |
| `array_cardinality(array)` | 返回数组元素个数 |
| `array_contains(array, element)` | 判断数组是否包含元素 |
| `array_filter(array, predicate)` | 过滤数组 |
| `current_login_user()` | 返回当前登录用户 |
| `team(include_manager, team_name)` | 返回团队成员 |
| `participate_roles()` | 返回所有参与角色 |
| `all_participate_persons()` | 返回所有参与人员 |
| `RELATIVE_DATETIME_EQ(col, 'date_para')` | 等于相对时间 |
| `RELATIVE_DATETIME_GT(col, 'date_para')` | 大于相对时间 |
| `RELATIVE_DATETIME_BETWEEN(col, 'date_para')` | 属于相对时间范围 |
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
### 示例
```sql
-- 查询用户1负责的、P0优先级的需求
SELECT work_item_id 
FROM space.story 
WHERE array_contains(owner, '用户1') 
  AND priority = 'P0'
  AND finish_status = true
-- 查询过去30天创建的需求
SELECT * 
FROM space.story 
WHERE RELATIVE_DATETIME_BETWEEN(create_time, 'past', '30d')
ORDER BY priority DESC
LIMIT 100
```
---
## 📋 Open API 错误码详解
> 文档：[Open API 错误码](https://project.feishu.cn/b/helpcenter/1p8d7djs/5aueo3jr)
本文档详细解释了 Open API 的各类错误码，包括其具体含义、产生原因和推荐的解决方法。
### HTTP 状态码一览
| 状态码 | 说明 | 错误码范围 |
|--------|------|-----------|
| **400** | Bad Request | 20001 - 20090, 9999, 1000051942 |
| **401** | Unauthorized | 10021 - 10302 |
| **403** | Forbidden | 10001 - 10404 |
| **404** | Not Found | 13001, 30005 - 30015 |
| **429** | Too Many Requests | 10429 - 10430 |
| **500** | Internal Server Error | 50006 |
---
### 1. 权限与认证错误 (401/403)
| HTTP | 错误码 | 说明 | 排查建议 |
|------|--------|------|----------|
| 403 | **10001** | No Permission - 没有操作权限 | • 更新评论：当前操作人不是评论创建人 <br>• 更新或删除视图：没有编辑视图的权限 <br>• 创建子任务：子任务所在空间和路径空间不匹配 <br>• 节点完成/状态流转：没有权限完成 <br>• 获取指定的工作项列表：未在对应空间安装插件 <br>• 工作项配置修改：操作者没有工作项配置权限 |
| 403 | **10002** | Illegal Operation - 非法操作 | • 节点完成/状态流转：当前节点/状态不能流转到指定节点/状态 |
| 403 | **10004** | Operation Failed - 操作失败 | • 该操作会导致流程图节点消失，而流程图至少需存在一个节点 |
| 401 | **10021** | Token Not Exist - 请求头未传plugin_token | • 请求头未传plugin_token |
| 401 | **10022** | Check Token Failed - plugin_token校验失败 | • plugin_token校验失败 |
| 403 | **10210** | code invalid - 授权码无效 | • 获取用户访问凭证时，入参的code(授权码)无效 |
| 403 | **10211** | Token Info Is Invalid - plugin_token信息不合法 | • 插件token错误，无法解析出具体信息 <br>• 如果是插件级token，未传X-USER-KEY |
| 401 | **10301** | Check Token Perm Failed - plugin_token权限校验未通过 | • 没申请当前操作接口的api权限 <br>• 申请权限了，但是未发布版本或重新发布版本 <br>• 发布版本了，但是未安装或更新插件 <br>• 空间ID不存在 <br>• 操作者没有对应空间权限 |
| 401 | **10302** | User Is Resigned - 用户已离职 | • 用户已离职 |
| 403 | **10404** | No Project Permission - 没有空间访问权限 | • 当前操作用户没有空间访问权限 |
---
### 2. 频率限制 (429)
| HTTP | 错误码 | 说明 | 排查建议 |
|------|--------|------|----------|
| 429 | **10429** | API Request Frequency Limit - 请求过于频繁 | • 接口请求过于频繁，超过同一token请求同一接口 15qps 限制 |
| 429 | **10430** | API Request Idempotent Limit - 幂等限制 | • 接口请求幂等限制，可排查header中X-IDEM-UUID幂等串是否冲突 |
---
### 3. 参数错误 (400)
#### 3.1 限制类错误
| 错误码 | 说明 | 排查建议 |
|--------|------|----------|
| **20001** | Param Request Limit | • 获取空间详情：传入的空间id超过最大100的限制 |
| **20002** | Page Size Limit | • 查询评论：page_size超过最大200的限制 |
| **20004** | Search User Limit | • 获取指定的工作项列表：user_keys超过最大10的限制 <br>• 获取用户详情：最多一次可查询100个用户详情 |
| **20019** | Invite Bot Limit 5 | • 拉取群的机器人数量超过5个 |
| **20020** | Bot App_ids Empty | • 未填拉取群的机器人 |
| **20024** | Uploaded File Size Limit 100M | • 上传文件大小限制100M |
| **20028** | Workitem Ids Limit 50 | • 传入的工作项id数量限制50 |
| **20043** | View Ids Limit 10 | • 查询视图列表时，最多只能传入10个视图id |
| **20055** | Search Result Bigger Than 2000 | • 查询结果超过2000个，请重新设置筛选条件 |
| **20057** | Search ProjectKeys And SimpleNames Limit 10 | • 搜索时，传入的ProjectKeys和SimpleNames的并集不能超过10个 |
| **20064** | Search Option Size Too Large | • 搜索指定的选项个数超过限制，最多可传入50个选项值 |
#### 3.2 参数格式类错误
| 错误码 | 说明 | 排查建议 |
|--------|------|----------|
| **20003** | Wrong WorkItemType Param | • 获取指定的工作项列表：work_item_type_keys未填 <br>• 获取工作项详情：目标工作项为状态流，但设置了"expand":{"need_workflow":true} |
| **20005** | Missing Param | • 必填的请求参数未填 |
| **20006** | Invalid Param | • 请求参数不合法，请检查参数与字段与属性解析格式是否匹配 <br>• 创建子任务：节点负责人和角色绑定，需通过role_assignee字段指定负责人 |
| **20013** | Invalid Time Interval | • 时间相关参数，不是毫秒时间戳(13位数字) |
| **20025** | Field_key/Field_alias Missing | • 字段的key和对接标识都缺失 |
| **20026** | FlowType Is Error (Status Flow) | • 查询的工作项是状态流工作项 |
| **20027** | FlowType Is Error (Node Flow) | • 查询的工作项是节点流工作项 |
| **20029** | Unsupported Field Type | • 不支持更新的字段类型(字段key会在msg中返回) |
| **20039** | X-USER-KEY Required | • 使用的是应用token，请求头中必须带上X-USER-KEY |
| **20040** | Request Form Is Null | • 上传附件时，传入的表单是空(content-type未选择multipart/form-data) |
| **20042** | X-User-Key Is Wrong | • 未填X-User-Key或者传入的user_key错误 |
| **20070** | Field is InValid | • 上传附件时指定的字段已失效 |
| **20080** | Query Length Error | • 综搜查询中query必填同时长度限制小于200 |
| **20081** | Query Type Not Supported | • 综搜查询中目前仅支持查询工作项和视图 |
#### 3.3 业务状态类错误
| 错误码 | 说明 | 排查建议 |
|--------|------|----------|
| **20007** | WorkItem Is Already Aborted | • 工作项已经被终止 |
| **20008** | WorkItem Is Already Restored | • 工作项已经被恢复 |
| **20009** | Abort Or Restore WorkItem No Reason | • 终止/恢复工作项，缺失了原因 |
| **20016** | Node Is Not Arrived | • 节点未到达，无法进行节点流转 |
| **20017** | Node Is Completed | • 节点已经完成，无法再次完成 |
| **20018** | Node ID Not Exist In Workflow | • 节点不存在当前工作项的节点流配置中 <br>• 解决方案：请检查工作项配置中是否存在该节点id |
| **20037** | Node Is Not Completed (Rollback) | • 节点未完成，无法回滚 |
| **20038** | Required Field Is Not Set | • 节点完成/状态流转时，必填字段未填写 |
| **20041** | Field RoleOwner Must Be Set | • 创建工作项必须传入role_owners字段 |
| **20044** | WorkItem Has Been Disabled | • 工作项已被禁用，无法查询到元数据 |
| **20046** | Task ID Not Exist In Workflow | • 工作流中不存在该子任务 |
| **20056** | Only WorkFlow Mode Can Be Aborted | • 只有工作流模式可以终止和恢复 |
| **20082** | Action Not Supported | • 子任务状态更新接口中只有回滚和确认操作 |
| **20083** | Duplication Field Exist | • 创建工作项时，传入的字段重复 |
| **20090** | Request Intercepted | • 请求或操作被插件拦截了 |
#### 3.4 角色与人员类错误
| 错误码 | 说明 | 排查建议 |
|--------|------|----------|
| **20047** | Role_Assignee and Assignee Conflict | • 更新或创建子任务，role_assignee和assignee不能同时传入 |
| **20048** | Role_Assignee's Role Not Match | • 更新或创建子任务时，传入的role_assignee里面的role与节点绑定的role不匹配 |
| **20050** | Field Option Value Is Wrong | • 更新或创建工作项时，传入的字段选项值错误 |
| **20051** | FieldLinkedStory Value Is Wrong | • 填入的field_linked_story字段值错误 |
| **20052** | IssueOperator Value Is Wrong | • 缺陷的operator角色负责人填写错误 |
| **20053** | IssueReporter Value Is Wrong | • 缺陷的reporter角色负责人填写错误 |
| **20058** | SearchUser.Role And FieldKey Conflict | • 搜索时，在SearchUser这个结构中，Role和FieldKey不能同时出现 |
| **20059** | SearchUser.UserKeys Must Appear | • 搜索时，UserKeys如果为空，Role或FieldKey不能单独传入 |
| **20071** | Search People Not Support Issue | • 搜索指定参数是people时，不支持缺陷工作项 |
#### 3.5 关联关系类错误
| 错误码 | 说明 | 排查建议 |
|--------|------|----------|
| **20010** | WorkItemType Is Not Same | • 创建、更新视图：传入的工作项id列表不是同一种工作项类型 |
| **20011** | Input View Is Not Fix View | • 删除固定视图：目前只支持删除固定视图 |
| **20012** | View Is Not In The Input Project | • 视图的id所属空间不属于参数中的空间 |
| **20014** | Project And WorkItem Not Match | • 工作项所属空间和传入的空间不匹配 |
| **20015** | Field Mix With '-' And Without '-' | • 存在相同的字段出现在需要返回列表和不需返回列表 |
| **20021** | ChatID Not Belong WorkItem | • 群id不属于参数中的工作项 |
| **20032** | DifferentSchedule Set Owner Invalid | • 差异化排期未指定用户更新排期，或者指定的用户超过一个 |
| **20033** | Update Field Invalid | • 不能更新该字段 |
| **20045** | Comment And WorkItem Not Match | • 评论不属于指定的工作项 |
| **20060** | WorkItemTypeKey Not Match | • 工作项类型或关联工作项类型，与关联关系的配置不匹配 |
| **20061** | RelationKey Type Not Relation | • 指定关联关系的字段类型，不是关联类型 |
| **20062** | RelationType Error | • 指定的RelationType不存在，目前只支持0：字段key,1：字段alias |
| **20063** | Search Operator Error | • 搜索的操作错误，不同的参数可使用的操作符不同 |
| **20065** | Search Param Key Not Support StateFlow | • 当前参数不支持状态流筛选 |
| **20066** | Search Signal Only Support len=1 | • 搜索系统外信号，当操作符是=或者!=时，数组长度只能是1 |
| **20067** | Search Signal Not Support Value | • 搜素系统外信号，传入的值不支持筛选 |
| **20068** | Search Param Is Not Support | • 指定的参数不支持筛选 |
| **20069** | Search Param Value Error | • 搜索传入的参数值异常 |
| **20072** | Conjunction Value Only Support AND/OR | • 搜索的Conjunction仅支持且、或 |
---
### 4. 资源未找到 (404)
| 错误码 | 说明 | 排查建议 |
|--------|------|----------|
| **13001** | View not exist | • 视图不存在 |
| **30005** | WorkItem Not Found | • 工作项已删除 <br>• 查询的工作项id不正确 <br>• 查询的工作项类型和工作项id不匹配 |
| **30006** | User Not Found | • header中传入的user_key未找到对应用户 <br>• 未查到指定用户 <br>• 使用了虚拟token，只能查到插件协作者相关信息 <br>• 确认插件是否共享到插件市场 |
| **30007** | Workflow Not Found | • 工作项中未找到节点流 |
| **30008** | Business Not Found | • 空间下的业务线未找到 |
| **30009** | Field Not Found | • 字段未在字段配置中，无法更新或创建 |
| **30010** | Stateflow Not Found | • 工作项中未找到状态流 |
| **30011** | Node Not Found In Workflow | • 节点在工作项的节点流配置中未找到 |
| **30012** | State Not Found In Stateflow | • 状态在工作项的状态流配置中未找到 |
| **30015** | Record Not Found | • 更新评论失败，评论ID不存在，需检查参数评论ID |
---
### 5. 其他错误
| 错误码 | 说明 | 排查建议 |
|--------|------|----------|
| **9999** | Invalid Param | • 数据结构不匹配，检查入参 |
| **1000051942** | commercial usage exceeded | • 商业化用量超限 |
---
### 6. 服务器内部错误 (500)
| HTTP | 错误码 | 说明 | 排查建议 |
|------|--------|------|----------|
| 500 | **50006** | RPC Call Error - This type does not support | • 未知错误，请联系技术支持 |
| 500 | **50006** | openapi system err, please try again later | • 服务瞬时超载，可尝试重试恢复 |
---
## 相关文档
- [[飞书项目OpenAPI完整API列表]] - 主索引
- [[飞书项目OpenAPI概述]] - API 概述
- [[搜索参数格式及常用示例 - 开发者手册 - 飞书项目帮助中心]] - 搜索参数
---
## 🏷️ 标签
#飞书项目 #OpenAPI #开发者手册 #工作项 #API