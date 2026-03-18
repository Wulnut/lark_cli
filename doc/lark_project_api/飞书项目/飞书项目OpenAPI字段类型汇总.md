---
title: 飞书项目OpenAPI字段类型汇总
tags: [飞书项目, OpenAPI]
category: 飞书项目
created: 2026-03-17
updated: 2026-03-17

---


> 原文：[工作项 API 汇总](https://project.feishu.cn/b/helpcenter/1p8d7djs/1tj6ggll)
> 版本：2.0.0
> 更新时间：2026-03-17
---
## 📋 概述
本文档详细介绍飞书项目 (Project) OpenAPI 中各类字段的数据结构和传参规则。
---
## 📝 字段类型详解
### 1. 文本类
| 类型 | field_type_key | 说明 | 示例 |
|------|----------------|------|------|
| 单行文本 | `text` | 单行文本内容 | `"文本"` |
| 多行富文本 | `multi_text` | 支持富文本格式 | 见富文本格式部分 |
---
### 2. 选项类
| 类型 | field_type_key | 说明 | 数据结构 |
|------|----------------|------|----------|
| 单选 | `select` | 单个选项 | `{"label": "选项1", "value": "8lheuaepp"}` |
| 多选 | `multi_select` | 多个选项 | `[{"label": "选项1", "value": "b0gzgge5o"}, ...]` |
| 单选树 | `tree_select` | 单层树形选择 | 含 `children` 字段 |
| 多选树 | `tree_multi_select` | 多层树形选择 | 含多层 `children` |
| 单选按钮 | `radio` | 单选按钮 | `{"label": "选项1", "value": "zgw2edjby"}` |
---
### 3. 用户类
| 类型 | field_type_key | 说明 | 数据结构 |
|------|----------------|------|----------|
| 用户 | `user` | 单个用户 | `"735679528XXXXX"` (user_key) |
| 多用户 | `multi_user` | 多个用户 | `["user_key1", "user_key2", ...]` |
---
### 4. 日期时间类
| 类型 | field_type_key | 说明 | 数据结构 |
|------|----------------|------|----------|
| 日期 | `date` | 天精度日期 | 毫秒时间戳 `1722182400000` |
| 日期时间 | `datetime` | 精确到秒 | 毫秒时间戳 `1722220183000` |
| 日期区间 | `schedule` | 开始+结束 | `{"start_time": xxx, "end_time": xxx}` |
---
### 5. 数值类
| 类型 | field_type_key | 说明 | 示例 |
|------|----------------|------|------|
| 数字 | `number` | 浮点数 | `11.11111111111111` |
| 布尔 | `bool` | true/false | `true` / `false` |
---
### 6. 附件类
| 类型 | field_type_key | 说明 | 数据结构 |
|------|----------------|------|----------|
| 附件 | `file` | 单个文件 | 含 `uid`, `url`, `name`, `size` |
| 多附件 | `multi_file` | 多个文件 | 文件数组 |
| 图片 | `image` | 图片文件 | 含 `type`, `size`, `uid`, `url` |
---
### 7. 关联类
| 类型 | field_type_key | 说明 | 数据结构 |
|------|----------------|------|----------|
| 关联单选 | `work_item_related_select` | 关联单个工作项 | 工作项 ID |
| 关联多选 | `work_item_related_multi_select` | 关联多个工作项 | `[id1, id2, ...]` |
---
### 8. 其他类型
| 类型 | field_type_key | 说明 |
|------|----------------|------|
| 链接 | `link` | URL 链接 |
| 云文档 | `link_cloud_doc` | 飞书云文档链接 |
| 电话 | `telephone` | 电话号码 |
| 邮箱 | `email` | 邮箱地址 |
| 业务线 | `business` | 业务线 ID |
| 群组 | `chat_group` / `group_id` | 飞书群组 ID |
| 角色 | `role_owners` | 角色+成员列表 |
---
## 🔗 富文本格式 (multi_text)
### 支持的格式元素
| 类型 | 说明 | 结构 |
|------|------|------|
| `paragraph` | 段落 | 文本内容数组 |
| `text` | 文本 | 含 `text` 和 `attrs` |
| `horizontalLine` | 水平线 | - |
| `checklist` | 检查清单 | 含 `content` 和 `lineAttrs` |
| `table` | 表格 | 含 `tableInfo` 和 `cellList` |
| `ul` | 无序列表 | - |
| `ol` | 有序列表 | 含 `list` 编号 |
| `blockquote` | 引用 | - |
### 文本属性 (attrs)
| 属性 | 说明 | 值 |
|------|------|-----|
| `bold` | 加粗 | `"true"` / `"false"` |
| `italic` | 斜体 | `"true"` / `"false"` |
| `underline` | 下划线 | `"true"` / `"false"` |
| `strikethrough` | 删除线 | `"true"` / `"false"` |
| `fontColor` | 字体颜色 | 颜色名称或 RGB |
| `backgroundColor` | 背景颜色 | 颜色名称或 RGB |
| `fontSize` | 字体大小 | `"h1"`, `"h2"`, `"h3"` 等 |
### 链接
```json
{
  "type": "hyperlink",
  "attrs": {
    "title": "测试链接",
    "url": "https://project.feishu.cn/b/helpcenter/1p8d7djs/16ejynnw"
  }
}
```
---
## 📤 传参规则
### 基本字段结构
```json
{
  "field_key": "field_xxx",      // 字段 key
  "field_value": "值",            // 字段值
  "field_type_key": "text",       // 字段类型
  "field_alias": ""               // 字段对接标识（可选）
}
```
### 注意事项
1. **用户字段**：传 user_key（飞书用户 ID）
2. **日期字段**：传毫秒时间戳，天精度建议传入 `00:00:00`
3. **选项类字段**：
   - 入参时 `label` 可选，`value` 必填
   - 出参时 `label` 和 `value` 都有
4. **富文本字段**：需要构建复杂的 JSON 结构
---
## 📌 常见 FAQ
### Q1: 如何更新富文本字段（文本+图片）？
需要构建完整的富文本 JSON 结构，包含 `type: "paragraph"` 和 `type: "image"` 元素。
### Q2: 超链接如何更新？
```json
{
  "type": "hyperlink",
  "attrs": {
    "title": "链接标题",
    "url": "https://example.com"
  }
}
```
### Q3: 关注人（watchers）字段是什么类型？
类型为 `multi_user`，存储用户 user_key 数组。
---
## 🏷️ 标签
#飞书项目 #OpenAPI #字段类型 #API文档