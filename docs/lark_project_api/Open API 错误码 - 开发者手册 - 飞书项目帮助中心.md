---
title: "Open API 错误码 - 开发者手册 - 飞书项目帮助中心"
source: "https://project.feishu.cn/b/helpcenter/1p8d7djs/5aueo3jr"
author:
published:
created: 2026-03-19
updated: 2026-03-19
description: "飞书项目 Open API 错误码详解，包含 HTTP 状态码、err_code 错误码的具体含义、产生原因和排查建议"
tags: [飞书项目, OpenAPI, 错误码, err_code, HTTP状态码, 排查建议]
category: 飞书项目
related_docs:
  - "[[飞书项目API开发者知识库]]"
  - "[[API 列表 - 开发者手册 - 飞书项目帮助中心]]"
  - "[[飞书项目OpenAPI完整API列表]]"
---

本文档详细解释了 Open API 的各类错误码，包括其具体含义、产生原因和推荐的解决方法。

## 错误码总览

| HTTP 状态码 | 错误码 | 错误信息 | 可能原因 & 排查建议 |
| --- | --- | --- | --- |
| 403 | 10001 | No Permission | 没有操作权限<br>- 更新评论：当前操作人不是评论创建人<br>- 更新或删除视图：没有编辑视图的权限<br>- 创建子任务：子任务所在空间和路径空间不匹配<br>- 节点完成/状态流转：没有权限完成<br>- 获取指定的工作项列表：未在对应空间安装插件 |
| 403 | 10002 | Illegal Operation | 非法操作<br>- 节点完成/状态流转：当前节点/状态不能流转到指定节点/状态 |
| 403 | 10004 | Operation Failed | 操作失败<br>- 该操作会导致流程图节点消失，而流程图至少需存在一个节点 |
| 401 | 10021 | Token Not Exist | 请求头未传 plugin_token |
| 401 | 10022 | Check Token Failed | plugin_token 校验失败 |
| 403 | 10210 | code invalid | 获取用户访问凭证时，入参的 code(授权码)无效 |
| 403 | 10211 | Token Info Is Invalid | plugin_token 信息不合法<br>- 插件 token 错误，无法解析出具体信息<br>- 如果是插件级 token，未传 X-USER-KEY |
| 401 | 10301 | Check Token Perm Failed | plugin_token 权限校验未通过<br>- 没申请当前操作接口的 API 权限<br>- 申请权限了，但是未发布版本或重新发布版本<br>- 发布版本了，但是未安装或更新插件<br>- 空间 ID 不存在<br>- 操作者没有对应空间权限 |
| 401 | 10302 | Check User Error, User Is Resigned | 用户已离职 |
| 404 | 13001 | View not exist | 视图不存在 |
| 403 | 10404 | No Project Permission | 当前操作用户没有空间访问权限 |
| 429 | 10429 | API Request Frequency Limit | 接口请求过于频繁，超过同一 token 请求同一接口 15qps 限制 |
| 429 | 10430 | API Request Idempotent Limit | 请求幂等性限制 |

## 请求参数错误（HTTP 400）

| 错误码 | 错误信息 | 说明 |
| --- | --- | --- |
| 20001 | Param Request Limit | 获取空间详情：传入的空间 id 超过最大 100 的限制 |
| 20002 | Page Size Limit | 查询评论：page_size 超过最大 200 的限制 |
| 20003 | Wrong WorkItemType Param | 获取指定的工作项列表：work_item_type_keys 未填 |
| 20004 | Search User Limit | 获取指定的工作项列表：user_keys 超过最大 10 的限制 |
| 20005 | Missing Param | 必填的请求参数未填 |
| 20006 | Invalid Param | 请求参数不合法，请检查参数与字段格式是否匹配 |
| 20007 | WorkItem Is Already Aborted | 工作项已经被终止 |
| 20008 | WorkItem Is Already Restored | 工作项已经被恢复 |
| 20009 | Abort Or Restore WorkItem No Reason | 终止/恢复工作项，缺失了原因 |
| 20010 | WorkItemType Is Not Same | 创建、更新视图：传入的工作项 id 列表不是同一种工作项类型 |
| 20011 | Input View Is Not Fix View | 删除固定视图：想要删除的视图不是固定视图 |
| 20012 | View Is Not In The Input Project | 视图的 id 所属空间不属于参数中的空间 |
| 20013 | Invalid Time Interval | 时间相关参数，不是毫秒时间戳(13位数字) |
| 20014 | Project And WorkItem Not Match | 工作项所属空间和传入的空间不匹配 |
| 20015 | Field Mix With '-' And Without '-' | 存在相同的字段出现在需要返回列表和不需要返回列表 |
| 20016 | Node Is Not Arrived, Could Not Be operated | 节点未到达，无法进行节点流转 |
| 20017 | Node Is Completed, Could Not Be Completed | 节点已经完成，无法再次完成 |
| 20018 | Node ID Not Exist In Workflow | 节点不存在当前工作项的节点流配置中 |
| 20019 | Invite Bot Limit 5 | 拉取群的机器人数量超过 5 个 |
| 20020 | Bot App_ids Empty | 未填拉取群的机器人 |
| 20021 | ChatID Not Belong WorkItem | 群 id 不属于参数中的工作项 |
| 20024 | Uploaded File Size Limit 100M | 上传文件大小限制 100M |
| 20025 | Field_key/Field_alias Missing | 字段的 key 和对接标识都缺失 |
| 20026 | FlowType Is Error, Please Convert To Status FlowType | 查询的工作项是状态流工作项 |
| 20027 | FlowType Is Error, Please Convert To Node FlowType | 查询的工作项是节点流工作项 |
| 20028 | Workitem Ids Limit 50 | 传入的工作项 id 数量限制 50 |
| 20029 | Unsupported Field Type | 不支持更新的字段类型 |
| 20032 | DifferentSchedule Set Owner Invalid | 差异化排期未指定用户更新排期，或者指定的用户超过一个 |
| 20033 | Update Field Invalid | 不能更新该字段 |
| 20037 | Node Is not completed, Could Not Be Rollback | 节点未完成，无法回滚 |
| 20038 | Required Field Is Not Set | 节点完成/状态流转时，必填字段未填写 |
| 20039 | - | 使用的是应用 token，请求头中必须带上 X-USER-KEY |
| 20040 | Request Form Is Null | 上传附件时，传入的表单是空 |
| 20041 | Field RoleOwner Must Be Set In FieldValuePairs | 创建工作项必须传入 role_owners 字段 |
| 20042 | X-User-Key Is Wrong, please Check First | 未填 X-User-Key 或者传入的 user_key 错误 |
| 20043 | View Ids Limit 10 | 查询视图列表时，最多只能传入 10 个视图 id |
| 20044 | WorkItem Has Been Disabled | 工作项已被禁用，无法查询到元数据 |
| 20045 | Comment And WorkItem Not Match | 评论不属于指定的工作项 |
| 20046 | Task ID Not Exist In Workflow | 工作流中不存在该子任务 |
| 20047 | Role_Assignee and Assignee Can't Be Set Together | 更新或创建子任务，role_assignee 和 assignee 不能同时传入 |
| 20048 | Role_Assignee's Role Is Not Match | 更新或创建子任务时，传入的 role_assignee 里面的 role 与节点绑定的 role 不匹配 |
| 20049 | TenantGroupID Is Wrong | 仅渠道用户使用，填入的 TenantGroupId 错误 |
| 20050 | 「field_key」Field Option Value Is Wrong | 更新或创建工作项时，传入的字段选项值错误 |
| 20051 | FieldLinkedStory Value Is Wrong | 填入的 field_linked_story 字段值错误 |
| 20052 | IssueOperator Value Is Wrong | 缺陷的 operator 角色负责人填写错误 |
| 20053 | IssueReporter Value Is Wrong | 缺陷的 reporter 角色负责人填写错误 |
| 20055 | Search Result Bigger Than 2000, Please Check Your Search Params | 查询结果超过 2000 个，请重新设置筛选条件 |
| 20056 | Only WorkFlow Mode Can Be Aborted Or Restored | 只有工作流模式可以终止和恢复 |
| 20057 | Search ProjectKeys And SimpleNames Limit 10 | 搜索时，传入的 ProjectKeys 和 SimpleNames 的并集不能超过 10 个 |
| 20058 | SearchUser.Role And SearchUser.FieldKey Can't appear together | 搜索时，Role 和 FieldKey 不能同时出现 |
| 20059 | SearchUser.FieldKey Or SearchUser.Role Can't Appear Alone | 搜索时，UserKeys 如果为空，Role 或 FieldKey 不能单独传入 |
| 20060 | WorkItemTypeKey Or RelationWorkItemTypeKey Is Not Match In Field Configuration | 工作项类型或关联工作项类型，与关联关系的配置不匹配 |
| 20061 | RelationKey Type Is Not Relation In Configuration | 指定关联关系的字段类型，不是关联类型 |
| 20062 | RelationType Error | 指定的 RelationType 不存在，目前只支持 0=字段key，1=字段alias |
| 20063 | Search Operator Error | 搜索的操作错误，不同的参数可使用的操作符不同 |
| 20064 | Search Option Size Too Large | 搜索指定的选项个数超过限制，最多可传入 50 个选项值 |
| 20065 | Search Param Key Not Support StateFlow | 当前参数不支持状态流筛选 |
| 20066 | Search Signal Only Support len=1 When Operator Is 「=」 or 「!=」 | 搜索系统外信号，当操作符是 = 或 != 时，数组长度只能是 1 |
| 20067 | Search Signal Not Support Value | 搜索系统外信号，传入的值不支持筛选 |
| 20068 | Search Param Is Not Support | 指定的参数不支持筛选 |
| 20069 | Search Param Value Error | 搜索传入的参数值异常 |
| 20070 | Field is InValid | 上传附件时指定的字段已失效 |
| 20071 | Search People Not Support Issue | 搜索指定参数是 people 时，不支持缺陷工作项 |
| 20072 | Conjunction Value Only Support 「AND」、「OR」 | 搜索的 Conjunction 仅支持且、或 |
| 20080 | Query length must be larger than 0, smaller than 200 | 综搜查询中 query 必填同时长度限制小于 200 |
| 20081 | Query type is not supported | 综搜查询中目前仅支持查询工作项和视图 |
| 20082 | Action is not supported | 子任务状态更新接口中只有回滚和确认操作 |
| 20083 | Duplication Field Exist, FieldKey | 创建工作项时，传入的字段重复 |
| 20090 | Your request has been intercepted | 请求或操作被插件拦截了 |

## 资源未找到（HTTP 404）

| 错误码 | 错误信息 | 说明 |
| --- | --- | --- |
| 30005 | WorkItem Not Found | 工作项未找到，可能是工作项已删除、查询的工作项 id 不正确、工作项类型和工作项 id 不匹配 |
| 30006 | User Not Found | 未查到指定用户；使用了虚拟 token，只能查到插件协作者相关信息；确认插件是否共享到插件市场 |
| 30007 | Workflow Not Found | 工作项中未找到节点流 |
| 30008 | Business Not Found | 空间下的业务线未找到 |
| 30009 | Field Not Found | 字段未在字段配置中，无法更新或创建 |
| 30010 | Stateflow Not Found | 工作项中未找到状态流 |
| 30011 | Node Not Found In Workflow | 节点在工作项的节点流配置中未找到 |
| 30012 | State Not Found In Stateflow | 状态在工作项的状态流配置中未找到 |
| 30015 | Record Not Found | 更新评论失败，评论 ID 不存在 |

## 服务端错误（HTTP 500）

| 错误码 | 错误信息 | 说明 |
| --- | --- | --- |
| 50006 | RPC Call Error This type does not support | 未知错误，请联系技术支持 |
| 50006 | openapi system err, please try again later | 服务瞬时超载，可尝试重试恢复 |

## 其他错误

| HTTP 状态码 | 错误码 | 错误信息 | 说明 |
| --- | --- | --- | --- |
| 400 | 9999 | Invalid Param | 数据结构不匹配，检查入参 |
| 400 | 1000051942 | commercial usage exceeded | 商业化用量超限 |
