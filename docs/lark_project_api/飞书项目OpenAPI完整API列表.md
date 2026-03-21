---
title: 飞书项目 OpenAPI 完整 API 列表
source: https://project.feishu.cn/b/helpcenter/1p8d7djs/wlcwhshe
author:
published: 2026-03-17
created: 2026-03-17
updated: 2026-03-19
description: 飞书项目 OpenAPI 完整 API 列表，包含 82 个 API，覆盖用户、空间、工作项、配置、视图、评论、度量等模块
tags: [飞书项目, OpenAPI, API列表, 工作项, 配置, 视图, 评论, 度量, 搜索]
category: 飞书项目
api_version: "2.0.0"
api_count: 82
related_docs:
  - "[[飞书项目API开发者知识库]]"
  - "[[API 列表 - 开发者手册 - 飞书项目帮助中心]]"
---

# 飞书项目 OpenAPI 完整 API 列表

> 文档：[API 列表](https://project.feishu.cn/b/helpcenter/1p8d7djs/wlcwhshe)
> 版本：2.0.0
> 更新时间：2026-03-17

---

## 📋 API 分类总览

| 分类 | 子类 | API 数量 |
|------|------|----------|
| **用户&用户组** | 用户组管理 | 5 个 |
| **空间** | 空间管理 | 2 个 |
| **工作项** | 实例搜索 | 5 个 |
| | 实例读写 | 16 个 |
| | 流程与节点 | 6 个 |
| | 子任务 | 6 个 |
| | 附件 | 4 个 |
| | 空间关联 | 4 个 |
| | 群组 | 1 个 |
| **配置** | 空间配置 | 3 个 |
| | 工作项配置 | 9 个 |
| | 流程模板配置 | 5 个 |
| | 流程角色配置 | 1 个 |
| **视图** | 视图管理 | 8 个 |
| **评论** | 评论管理 | 4 个 |
| **度量** | 度量管理 | 1 个 |
| **租户** | 租户管理 | 2 个 |

---

## 📁 详细文档分册

| 分册 | 文档 | 内容 |
|------|------|------|
| 1 | [工作项实例搜索](./飞书项目OpenAPI-工作项实例搜索.md) | 5 个搜索 API |
| 2 | [工作项实例读写](./飞书项目OpenAPI-工作项实例读写.md) | 16 个 CRUD API |
| 3 | [工作项流程与节点](./飞书项目OpenAPI-工作项流程与节点.md) | 6 个流程 API |
| 4 | [子任务](./飞书项目OpenAPI-子任务.md) | 6 个子任务 API |
| 5 | [用户与空间](./飞书项目OpenAPI-用户与空间.md) | 5 个用户 + 2 个空间 API |
| 6 | [附件与群组](./飞书项目OpenAPI-附件与群组.md) | 4 个附件 + 1 个群组 API |
| 7 | [配置与视图](./飞书项目OpenAPI-配置与视图.md) | 18 个配置 + 8 个视图 + 4 个空间关联 API |
| **附录** | [鉴权流程](./飞书项目OpenAPI鉴权流程.md) | 认证、OAuth、Header配置 |
| **8** | [依赖关系图](./飞书项目OpenAPI依赖关系.md) | 77 个 API 依赖关系 |
| **9** | [评论管理](./飞书项目OpenAPI-评论.md) | 4 个评论 API |
| **10** | [度量图表](./飞书项目OpenAPI-度量图表.md) | 1 个度量 API |

---

## 👤 用户&用户组

| API | 方法 | 说明 |
|-----|------|------|
| 获取用户详情 | POST | `/open_api/user/query` |
| 搜索租户内的用户列表 | POST | `/open_api/user/search` |
| 创建自定义用户组 | POST | `/open_api/:project_key/user_group` |
| 更新用户组成员 | PATCH | `/open_api/:project_key/user_group/members` |
| 查询用户组成员 | POST | `/open_api/:project_key/user_groups/members/page` |

---

## 🏠 空间

| API | 方法 | 说明 |
|-----|------|------|
| 获取空间列表 | POST | `/open_api/projects` |
| 获取空间详情 | POST | `/open_api/projects/detail` |

---

## 📋 工作项

### 工作项实例搜索

| API | 方法 | 说明 |
|-----|------|------|
| 获取指定的工作项列表（单空间） | POST | `/open_api/:project_key/work_item/filter` |
| 获取指定的工作项列表（跨空间） | POST | `/open_api/work_items/filter_across_project` |
| 获取指定的工作项列表（单空间-复杂传参） | POST | `/open_api/:project_key/work_item/:work_item_type_key/search/params` |
| 获取指定的工作项列表（全局搜索） | POST | `/open_api/compositive_search` |
| 获取指定的关联工作项列表（单空间） | POST | `/open_api/:project_key/work_item/:work_item_type_key/:work_item_id/search_by_relation` |

### 工作项实例读写

| API | 方法 | 说明 |
|-----|------|------|
| 获取工作项详情 | POST | `/open_api/:project_key/work_item/:work_item_type_key/query` |
| 获取创建工作项元数据 | GET | `/open_api/:project_key/work_item/:work_item_type_key/meta` |
| 创建工作项 | POST | `/open_api/:project_key/work_item/create` |
| 更新工作项 | PUT | `/open_api/:project_key/work_item/:work_item_type_key/:work_item_id` |
| 删除工作项 | DELETE | `/open_api/:project_key/work_item/:work_item_type_key/:work_item_id` |
| 终止/恢复工作项 | PUT | `/open_api/:project_key/work_item/:work_item_type_key/:work_item_id/abort` |
| 获取工作项操作记录 | POST | `/open_api/op_record/work_item/list` |
| 批量查询评审意见、评审结论 | POST | `/open_api/work_item/finished/batch_query` |
| 修改评审结论和评审意见 | POST | `/open_api/work_item/finished/update` |
| 评审结论标签值查询 | POST | `/open_api/work_item/finished/query_conclusion_option` |
| 获取工作项的工时登记记录列表 | POST | `/open_api/work_item/man_hour/records` |
| 新增工时登记记录 | POST | `/open_api/:project_key/work_item/:work_item_type_key/:work_item_id/work_hour_record` |
| 更新工时登记记录 | PUT | `/open_api/:project_key/work_item/:work_item_type_key/:work_item_id/work_hour_record` |
| 删除工时登记记录 | DELETE | `/open_api/:project_key/work_item/:work_item_type_key/:work_item_id/work_hour_record` |
| 冻结/解冻工作项 | PUT | `/open_api/work_item/freeze` |
| 交付物信息批量查询（WBS） | POST | `/open_api/work_item/deliverable/batch_query` |

### 流程与节点

| API | 方法 | 说明 |
|-----|------|------|
| 获取工作流详情 | POST | `/open_api/:project_key/work_item/:work_item_type_key/:work_item_id/workflow/query` |
| 获取工作流详情（WBS） | GET | `/open_api/:project_key/work_item/:work_item_type_key/:work_item_id/wbs_view` |
| 更新节点/排期 | PUT | `/open_api/:project_key/workflow/:work_item_type_key/:work_item_id/node/:node_id` |
| 节点完成/回滚 | POST | `/open_api/:project_key/workflow/:work_item_type_key/:work_item_id/node/:node_id/operate` |
| 状态流转 | POST | `/open_api/:project_key/workflow/:work_item_type_key/:work_item_id/node/state_change` |
| 获取指定节点/状态流转所需必填信息 | POST | `/open_api/work_item/transition_required_info/get` |

### 子任务

| API | 方法 | 说明 |
|-----|------|------|
| 获取指定的子任务列表（跨空间） | POST | `/open_api/work_item/subtask/search` |
| 获取子任务详情 | GET | `/open_api/:project_key/work_item/:work_item_type_key/:work_item_id/workflow/task` |
| 创建子任务 | POST | `/open_api/:project_key/work_item/:work_item_type_key/:work_item_id/workflow/task` |
| 更新子任务 | POST | `/open_api/:project_key/work_item/:work_item_type_key/:work_item_id/workflow/:node_id/task/:task_id` |
| 子任务完成/回滚 | POST | `/open_api/:project_key/work_item/:work_item_type_key/:work_item_id/subtask/modify` |
| 删除子任务 | DELETE | `/open_api/:project_key/work_item/:work_item_type_key/:work_item_id/task/:task_id` |

### 附件

| API | 方法 | 说明 |
|-----|------|------|
| 添加附件 | POST | `/open_api/:project_key/work_item/:work_item_type_key/:work_item_id/file/upload` |
| 上传文件或富文本图片 | POST | `/open_api/:project_key/file/upload` |
| 下载附件 | POST | `/open_api/:project_key/work_item/:work_item_type_key/:work_item_id/file/download` |
| 删除附件 | POST | `/open_api/file/delete` |

### 空间关联

| API | 方法 | 说明 |
|-----|------|------|
| 获取空间关联规则列表 | POST | `/open_api/:project_key/relation/rules` |
| 获取空间关联下的关联工作项实例列表 | POST | `/open_api/:project_key/relation/:work_item_type_key/:work_item_id/work_item_list` |
| 绑定空间关联的关联工作项实例 | POST | `/open_api/:project_key/relation/:work_item_type_key/:work_item_id/batch_bind` |
| 解绑空间关联的关联工作项实例 | DELETE | `/open_api/:project_key/relation/:work_item_type_key/:work_item_id` |

### 群组

| API | 方法 | 说明 |
|-----|------|------|
| 拉机器人入群 | POST | `/open_api/:project_key/work_item/:work_item_id/bot_join_chat` |

---

## ⚙️ 配置

### 空间配置

| API | 方法 | 说明 |
|-----|------|------|
| 获取空间下业务线详情 | GET | `/open_api/:project_key/business/all` |
| 获取空间下工作项类型 | GET | `/open_api/:project_key/work_item/all-types` |
| 获取空间下团队人员 | GET | `/open_api/:project_key/teams/all` |

### 工作项配置

| API | 方法 | 说明 |
|-----|------|------|
| 获取工作项基础信息配置 | GET | `/open_api/:project_key/work_item/type/:work_item_type_key` |
| 更新工作项基础信息配置 | PUT | `/open_api/:project_key/work_item/type/:work_item_type_key` |
| 获取字段信息 | POST | `/open_api/:project_key/field/all` |
| 创建自定义字段 | POST | `/open_api/:project_key/field/:work_item_type_key/create` |
| 更新自定义字段 | PUT | `/open_api/:project_key/field/:work_item_type_key` |
| 获取工作项关系列表 | GET | `/open_api/:project_key/work_item/relation` |
| 新增工作项关系 | POST | `/open_api/work_item/relation/create` |
| 更新工作项关系 | POST | `/open_api/work_item/relation/update` |
| 删除工作项关系 | DELETE | `/open_api/work_item/relation/delete` |

### 流程模板配置

| API | 方法 | 说明 |
|-----|------|------|
| 获取工作项下的流程模板列表 | GET | `/open_api/:project_key/template_list/:work_item_type_key` |
| 获取流程模板配置详情 | GET | `/open_api/:project_key/template_detail/:template_id` |
| 新增流程模板 | POST | `/open_api/template/v2/create_template` |
| 更新流程模板 | PUT | `/open_api/template/v2/update_template` |
| 删除流程模板 | DELETE | `/open_api/template/v2/delete_template/:project_key/:template_id` |

### 流程角色配置

| API | 方法 | 说明 |
|-----|------|------|
| 获取流程角色配置详情 | GET | `/open_api/:project_key/flow_roles/:work_item_type_key` |

---

## 📊 视图

| API | 方法 | 说明 |
|-----|------|------|
| 获取视图列表及配置信息 | POST | `/open_api/:project_key/view_conf/list` |
| 获取视图下工作项列表 | GET | `/open_api/:project_key/fix_view/:view_id` |
| 获取视图下工作项列表（全景视图） | POST | `/open_api/:project_key/view/:view_id` |
| 创建固定视图 | POST | `/open_api/:project_key/:work_item_type_key/fix_view` |
| 更新固定视图 | POST | `/open_api/:project_key/:work_item_type_key/fix_view/:view_id` |
| 创建条件视图 | POST | `/open_api/view/v1/create_condition_view` |
| 更新条件视图 | POST | `/open_api/view/v1/update_condition_view` |
| 删除视图 | DELETE | `/open_api/:project_key/fix_view/:view_id` |

---

## 💬 评论

| API | 方法 | 说明 |
|-----|------|------|
| 添加评论 | POST | `/open_api/:project_key/work_item/:work_item_type_key/:work_item_id/comment/create` |
| 获取评论列表 | POST | `/open_api/:project_key/work_item/:work_item_type_key/:work_item_id/comment/list` |
| 更新评论 | PUT | `/open_api/:project_key/work_item/:work_item_type_key/:work_item_id/comment/:comment_id` |
| 删除评论 | DELETE | `/open_api/:project_key/work_item/:work_item_type_key/:work_item_id/comment/:comment_id` |

---

## 📊 度量

| API | 方法 | 说明 |
|-----|------|------|
| 获取度量图表明细数据 | GET | `/open_api/:project_key/measure/:chart_id` |

---

## 🏢 租户

| API | 方法 | 说明 |
|-----|------|------|
| 获取租户信息 | GET | `/open_api/tenant/info` |
| 获取租户安装空间列表 | GET | `/open_api/tenant/installed_projects` |

---

## 📎 相关文档

### 官方文档
- [[API 列表 - 开发者手册 - 飞书项目帮助中心]] - 官方 API 列表
- [[全量搜索参数格式及常用示例 - 开发者手册 - 飞书项目帮助中心]] - 搜索参数详解
- [[搜索参数格式及常用示例 - 开发者手册 - 飞书项目帮助中心]] - 搜索参数基础
- [[字段与属性解析格式 - 开发者手册 - 飞书项目帮助中心]] - 字段类型与属性
- [[MQL 语法说明 - 开发者手册 - 飞书项目帮助中心]] - MQL 查询语法
- [[Open API 错误码 - 开发者手册 - 飞书项目帮助中心]] - 错误码参考
- [[数据结构汇总 - 开发者手册 - 飞书项目帮助中心]] - 数据结构定义

### 知识库整合
- [[飞书项目API开发者知识库]] - 开发者手册知识库汇总
- [[飞书项目/飞书项目API结构图谱.canvas]] - API 结构关系图

### 分册文档
- [[飞书项目OpenAPI-工作项实例搜索]] - 工作项搜索API (5个)
- [[飞书项目OpenAPI-工作项实例读写]] - 工作项CRUD API (16个)
- [[飞书项目OpenAPI-工作项流程与节点]] - 流程与节点API (6个)
- [[飞书项目OpenAPI-子任务]] - 子任务API (6个)
- [[飞书项目OpenAPI-用户与空间]] - 用户与空间API (7个)
- [[飞书项目OpenAPI-附件与群组]] - 附件与群组API (5个)
- [[飞书项目OpenAPI-配置与视图]] - 配置与视图API (26个)
- [[飞书项目OpenAPI-评论]] - 评论管理API (4个)
- [[飞书项目OpenAPI-度量图表]] - 度量图表API (1个)
- [[飞书项目OpenAPI鉴权流程]] - OAuth认证流程
- [[飞书项目OpenAPI依赖关系]] - API 依赖关系图

---

## 🏷️ 标签

#飞书项目 #OpenAPI #API列表 #完整API #工作项 #配置 #评论 #度量
