---
title: "API 列表 - 开发者手册 - 飞书项目帮助中心"
source: "https://project.feishu.cn/b/helpcenter/1p8d7djs/wlcwhshe"
author:
published:
created: 2026-03-19
updated: 2026-03-19
description: "飞书项目官方 OpenAPI 完整列表，包含用户、空间、工作项、配置、视图、评论、度量等模块共 82 个 API"
tags: [飞书项目, OpenAPI, 官方文档, API列表, 工作项, 配置, 视图, 评论, 度量]
category: 飞书项目
api_version: "2.0.0"
api_count: 82
related_docs:
  - "[[飞书项目OpenAPI完整API列表]]"
  - "[[飞书项目API开发者知识库]]"
---

## 用户 & 用户组

| API 名称 | 说明 |
| --- | --- |
| `POST /open_api/user/query` 获取用户详情 | 获取指定用户的详细信息 |
| `POST /open_api/user/search` 搜索租户内的用户列表 | 模糊搜索租户内的用户并返回其详细信息 |
| `POST /open_api/:project_key/user_group` 创建自定义用户组 | 创建自定义用户组 |
| `PATCH /open_api/:project_key/user_group/members` 更新用户组成员 | 更新用户组成员 |
| `POST /open_api/:project_key/user_groups/members/page` 查询用户组成员 | 查询用户组成员 |

## 空间

| API 名称 | 说明 |
| --- | --- |
| `POST /open_api/projects` 获取空间列表 | 获取指定用户有权限访问空间和插件安装空间的交集 |
| `POST /open_api/projects/detail` 获取空间详情 | 获取空间详情信息，包括管理员信息 |

## 工作项

### 工作项实例搜索

| API 名称 | 说明 |
| --- | --- |
| `POST /open_api/:project_key/work_item/filter` 获取指定的工作项列表（单空间） | 在指定空间中搜索符合条件的工作项实例列表 |
| `POST /open_api/work_items/filter_across_project` 获取指定的工作项列表（跨空间） | 跨多个空间搜索符合条件的工作项实例列表 |
| `POST /open_api/:project_key/work_item/:work_item_type_key/search/params` 获取指定的工作项列表（单空间-复杂传参） | 在指定空间搜索符合复杂筛选条件的工作项实例列表 |
| `POST /open_api/compositive_search` 获取指定的工作项列表（全局搜索） | 获取跨空间和工作项类型搜索符合条件的工作项实例列表 |
| `POST /open_api/:project_key/work_item/:work_item_type_key/:work_item_id/search_by_relation` 获取指定的关联工作项列表（单空间） | 获取与指定工作项存在关联的工作项实例列表 |

### 工作项实例读写

| API 名称 | 说明 |
| --- | --- |
| `POST /open_api/:project_key/work_item/:work_item_type_key/query` 获取工作项详情 | 获取指定工作项实例的详细信息 |
| `GET /open_api/:project_key/work_item/:work_item_type_key/meta` 获取创建工作项元数据 | 获取创建工作项实例的最小数据单元 |
| `POST /open_api/:project_key/work_item/create` 创建工作项 | 在指定空间和工作项类型下新增工作项实例 |
| `PUT /open_api/:project_key/work_item/:work_item_type_key/:work_item_id` 更新工作项 | 修改指定工作项实例 |
| `DELETE /open_api/:project_key/work_item/:work_item_type_key/:work_item_id` 删除工作项 | 删除指定工作项实例 |
| `PUT /open_api/:project_key/work_item/:work_item_type_key/:work_item_id/abort` 终止/恢复工作项 | 终止或恢复指定工作项实例 |
| `POST /open_api/op_record/work_item/list` 获取工作项操作记录 | 获取指定空间下的工作项操作记录 |
| `POST /open_api/work_item/finished/batch_query` 批量查询评审意见、评审结论 | 批量查询节点的评审意见和结论 |
| `POST /open_api/work_item/finished/update` 修改评审结论和评审意见 | 更新节点评审意见和结论 |
| `POST /open_api/work_item/finished/query_conclusion_option` 评审结论标签值查询 | 查询节点下配置的评审结论标签 |
| `POST /open_api/work_item/man_hour/records` 获取工作项的工时登记记录列表 | 仅在安装工时登记插件时有效 |
| `POST /open_api/:project_key/work_item/:work_item_type_key/:work_item_id/work_hour_record` 新增工时登记记录 | 在指定工作项下添加工时记录 |
| `PUT /open_api/:project_key/work_item/:work_item_type_key/:work_item_id/work_hour_record` 更新工时登记记录 | 更新指定工作项的工时记录 |
| `DELETE /open_api/:project_key/work_item/:work_item_type_key/:work_item_id/work_hour_record` 删除工时登记记录 | 删除指定工作项的工时记录 |
| `PUT /open_api/work_item/freeze` 冻结/解冻工作项 | 冻结或解冻工作项实例 |
| `POST /open_api/work_item/deliverable/batch_query` 交付物信息批量查询（WBS） | 查询交付物详情信息 |

## 流程与节点

| API 名称 | 说明 |
| --- | --- |
| `POST /open_api/:project_key/work_item/:work_item_type_key/:work_item_id/workflow/query` 获取工作流详情 | 获取工作项实例的工作流信息 |
| `GET /open_api/:project_key/work_item/:work_item_type_key/:work_item_id/wbs_view` 获取工作流详情（WBS） | 获取 WBS 工作流信息 |
| `PUT /open_api/:project_key/workflow/:work_item_type_key/:work_item_id/node/:node_id` 更新节点/排期 | 更新工作项实例的指定节点信息 |
| `POST /open_api/:project_key/workflow/:work_item_type_key/:work_item_id/node/:node_id/operate` 节点完成/回滚 | 完成或回滚工作项实例的指定节点 |
| `POST /open_api/:project_key/workflow/:work_item_type_key/:work_item_id/node/state_change` 状态流转 | 流转工作项实例到指定状态 |
| `POST /open_api/work_item/transition_required_info/get` 获取指定节点/状态流转所需必填信息 | 获取流转所需的必填信息 |

## 子任务

| API 名称 | 说明 |
| --- | --- |
| `POST /open_api/work_item/subtask/search` 获取指定的子任务列表（跨空间） | 跨空间搜索符合传入条件的子任务 |
| `GET /open_api/:project_key/work_item/:work_item_type_key/:work_item_id/workflow/task?node_id=:node_id` 获取子任务详情 | 获取工作项实例上的子任务详细信息 |
| `POST /open_api/:project_key/work_item/:work_item_type_key/:work_item_id/workflow/task` 创建子任务 | 在工作项实例的指定节点上创建子任务 |
| `POST /open_api/:project_key/work_item/:work_item_type_key/:work_item_id/workflow/:node_id/task/:task_id` 更新子任务 | 更新子任务详细信息 |
| `POST /open_api/:project_key/work_item/:work_item_type_key/:work_item_id/subtask/modify` 子任务完成/回滚 | 完成或回滚子任务 |
| `DELETE /open_api/:project_key/work_item/:work_item_type_key/:work_item_id/task/:task_id` 删除子任务 | 删除指定子任务 |

## 附件

| API 名称 | 说明 |
| --- | --- |
| `POST /open_api/:project_key/work_item/:work_item_type_key/:work_item_id/file/upload` 添加附件 | 在指定工作项的附件类型字段中添加附件 |
| `POST /open_api/:project_key/file/upload` 上传文件或富文本图片 | 通用文件上传接口 |
| `POST /open_api/:project_key/work_item/:work_item_type_key/:work_item_id/file/download` 下载附件 | 下载工作项下的指定附件 |
| `POST /open_api/file/delete` 删除附件 | 删除工作项下的附件 |

## 空间关联

| API 名称 | 说明 |
| --- | --- |
| `POST /open_api/:project_key/relation/rules` 获取空间关联规则列表 | 获取空间下配置的空间关联规则列表 |
| `POST /open_api/:project_key/relation/:work_item_type_key/:work_item_id/work_item_list` 获取空间关联下的关联工作项实例列表 | 获取有空间关联的工作项实例列表 |
| `POST /open_api/:project_key/relation/:work_item_type_key/:work_item_id/batch_bind` 绑定空间关联的关联工作项实例 | 建立空间关联绑定关系 |
| `DELETE /open_api/:project_key/relation/:work_item_type_key/:work_item_id` 解绑空间关联的关联工作项实例 | 解除空间关联绑定关系 |

## 群组

| API 名称 | 说明 |
| --- | --- |
| `POST /open_api/:project_key/work_item/:work_item_id/bot_join_chat` 拉机器人入群 | 将飞书机器人拉入工作项关联群 |

## 配置

### 空间配置

| API 名称 | 说明 |
| --- | --- |
| `GET /open_api/:project_key/business/all` 获取空间下业务线详情 | 获取空间的业务线信息 |
| `GET /open_api/:project_key/work_item/all-types` 获取空间下工作项类型 | 获取空间下所有工作项类型 |
| `GET /open_api/:project_key/teams/all` 获取空间下团队人员 | 获取团队详情信息，包括人员列表、管理员列表等 |

### 工作项配置

| API 名称 | 说明 |
| --- | --- |
| `GET /open_api/:project_key/work_item/type/:work_item_type_key` 获取工作项基础信息配置 | 获取指定工作项类型的基础信息配置 |
| `PUT /open_api/:project_key/work_item/type/:work_item_type_key` 更新工作项基础信息配置 | 更新指定工作项类型的基础信息配置 |
| `POST /open_api/:project_key/field/all` 获取字段信息 | 获取空间或工作项类型下所有字段的基础信息 |
| `POST /open_api/:project_key/field/:work_item_type_key/create` 创建自定义字段 | 创建新的自定义字段 |
| `PUT /open_api/:project_key/field/:work_item_type_key` 更新自定义字段 | 更新自定义字段的配置信息 |
| `GET /open_api/:project_key/work_item/relation` 获取工作项关系列表 | 获取空间下的工作项关联关系列表 |
| `POST /open_api/work_item/relation/create` 新增工作项关系 | 新增工作项关联关系 |
| `POST /open_api/work_item/relation/update` 更新工作项关系 | 更新工作项关联关系的配置信息 |
| `DELETE /open_api/work_item/relation/delete` 删除工作项关系 | 删除工作项关联关系 |

### 流程模板配置

| API 名称 | 说明 |
| --- | --- |
| `GET /open_api/:project_key/template_list/:work_item_type_key` 获取工作项下的流程模板列表 | 获取工作项类型下所有流程模板列表 |
| `GET /open_api/:project_key/template_detail/:template_id` 获取流程模板配置详情 | 获取流程模板的配置信息详情 |
| `POST /open_api/template/v2/create_template` 新增流程模板 | 创建新的流程模板 |
| `PUT /open_api/template/v2/update_template` 更新流程模板 | 更新流程模板的配置信息 |
| `DELETE /open_api/template/v2/delete_template/:project_key/:template_id` 删除流程模板 | 删除指定的流程模板 |

### 流程角色配置

| API 名称 | 说明 |
| --- | --- |
| `GET /open_api/:project_key/flow_roles/:work_item_type_key` 获取流程角色配置详情 | 获取工作项类型下所有角色与人员的配置信息 |

## 视图

| API 名称 | 说明 |
| --- | --- |
| `POST /open_api/:project_key/view_conf/list` 获取视图列表及配置信息 | 搜索符合条件的所有视图列表及相关配置信息 |
| `GET /open_api/:project_key/fix_view/:view_id?page_size=:page_size&page_num=:page_num` 获取视图下工作项列表（非全景视图） | 获取指定视图中的工作项实例列表 |
| `POST /open_api/:project_key/view/:view_id` 获取视图下工作项列表（全景视图） | 获取全景视图中的工作项实例列表和详情信息 |
| `POST /open_api/:project_key/:work_item_type_key/fix_view` 创建固定视图 | 新增固定视图 |
| `POST /open_api/:project_key/:work_item_type_key/fix_view/:view_id` 更新固定视图 | 对指定固定视图增/删工作项实例 |
| `POST /open_api/view/v1/create_condition_view` 创建条件视图 | 新增条件视图 |
| `POST /open_api/view/v1/update_condition_view` 更新条件视图 | 更新条件视图的筛选条件和协作者信息 |
| `DELETE /open_api/:project_key/fix_view/:view_id` 删除视图 | 删除指定空间的视图 |

## 评论

| API 名称 | 说明 |
| --- | --- |
| `POST /open_api/:project_key/work_item/:work_item_type_key/:work_item_id/comment/create` 添加评论 | 在指定工作项下添加评论 |
| `GET /open_api/:project_key/work_item/:work_item_type_key/:work_item_id/comment` 获取评论列表 | 获取指定工作项下的所有评论信息 |
| `PUT /open_api/:project_key/work_item/:work_item_type_key/:work_item_id/comment/:comment_id` 更新评论 | 更新指定评论的内容 |
| `DELETE /open_api/:project_key/work_item/:work_item_type_key/:work_item_id/comment/:comment_id` 删除评论 | 删除工作项下指定的评论 |

## 度量

| API 名称 | 说明 |
| --- | --- |
| `GET /open_api/:project_key/measure/:chart_id` 获取度量图表明细数据 | 获取指定度量图表的明细数据 |
