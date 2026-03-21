---
title: 飞书项目OpenAPI-附件与群组
source: https://project.feishu.cn/b/helpcenter/1p8d7djs
author:
published: 2026-03-17
created: 2026-03-17
updated: 2026-03-17
description: 飞书项目OpenAPI附件与群组操作详解，包含添加附件、上传文件、下载附件、删除附件及拉机器人入群等API
tags: [飞书项目, OpenAPI, 附件, 文件, 群组, 机器人]
category: 飞书项目
related_docs:
  - "[[飞书项目API开发者知识库]]"
  - "[[飞书项目OpenAPI完整API列表]]"
---


> 文档编号：6
> 更新时间：2026-03-17
---
## 附件 API 列表
| # | API | 方法 | 说明 |
|---|-----|------|------|
| 1 | 添加附件 | POST | `/open_api/:project_key/work_item/:work_item_type_key/:work_item_id/file/upload` |
| 2 | 上传文件或富文本图片 | POST | `/open_api/:project_key/file/upload` |
| 3 | 下载附件 | POST | `/open_api/:project_key/work_item/:work_item_type_key/:work_item_id/file/download` |
| 4 | 删除附件 | POST | `/open_api/file/delete` |
---
## 群组 API 列表
| # | API | 方法 | 说明 |
|---|-----|------|------|
| 1 | 拉机器人入群 | POST | `/open_api/:project_key/work_item/:work_item_id/bot_join_chat` |
---
## 1. 添加附件
**API**：`POST /open_api/:project_key/work_item/:work_item_type_key/:work_item_id/file/upload`
**说明**：用于在指定工作项的一个"附件类型"字段中添加附件
### 请求参数
| 参数 | 类型 | 必须 | 说明 |
|------|------|------|------|
| `work_item_id` | int64 | ✅ | 工作项 ID |
| `field_key` | string | ✅ | 附件字段 key |
| `file` | file | ✅ | 文件（二进制） |
### 请求方式
需要使用 `multipart/form-data` 格式上传
---
## 2. 上传文件或富文本图片
**API**：`POST /open_api/:project_key/file/upload`
**说明**：通用的文件上传接口，会返回上传后的资源路径，主要用于富文本中上传图片
### 请求参数
| 参数 | 类型 | 必须 | 说明 |
|------|------|------|------|
| `file_name` | string | ✅ | 文件名 |
| `file` | file | ✅ | 文件内容 |
### 返回参数
| 参数 | 类型 | 说明 |
|------|------|------|
| `file_key` | string | 文件 key |
| `url` | string | 文件访问 URL |
---
## 3. 下载附件
**API**：`POST /open_api/:project_key/work_item/:work_item_type_key/:work_item_id/file/download`
**说明**：用于下载一个工作项下的指定附件
### 请求参数
| 参数 | 类型 | 必须 | 说明 |
|------|------|------|------|
| `work_item_id` | int64 | ✅ | 工作项 ID |
| `field_key` | string | ✅ | 附件字段 key |
| `file_key` | string | ✅ | 文件 key |
---
## 4. 删除附件
**API**：`POST /open_api/file/delete`
**说明**：用于在指定工作项的一个"附件类型"字段中删除附件
### 请求参数
| 参数 | 类型 | 必须 | 说明 |
|------|------|------|------|
| `work_item_id` | int64 | ✅ | 工作项 ID |
| `field_key` | string | ✅ | 附件字段 key |
| `file_key` | string | ✅ | 文件 key |
---
## 5. 拉机器人入群
**API**：`POST /open_api/:project_key/work_item/:work_item_id/bot_join_chat`
**说明**：用于将指定的飞书机器人拉入工作项关联群
### 请求参数
| 参数 | 类型 | 必须 | 说明 |
|------|------|------|------|
| `work_item_id` | int64 | ✅ | 工作项 ID |
| `bot_id` | string | ✅ | 机器人 ID |
---
## 相关文档
- [工作项实例读写](飞书项目OpenAPI-工作项实例读写.md)
- [完整 API 列表](飞书项目OpenAPI完整API列表.md)
---
## 🏷️ 标签
#飞书项目 #OpenAPI #附件 #文件上传 #API文档