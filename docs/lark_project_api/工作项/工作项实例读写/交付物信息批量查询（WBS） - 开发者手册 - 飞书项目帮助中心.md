---
title: "交付物信息批量查询（WBS） - 开发者手册 - 飞书项目帮助中心"
source: "https://project.feishu.cn/b/helpcenter/1p8d7djs/xe2145vg"
author:
published:
created: 2026-03-21
updated: 2026-03-21
description: "查询交付物详情信息，目前支持查询工作项交付物"
tags:
  - "飞书项目"
  - "工作项"
  - "工作项实例读写"
  - "API列表"
  - "OpenAPI"
category: "飞书项目"
related_docs:
  - "[[获取工作项详情 - 开发者手册 - 飞书项目帮助中心]]"
---

该接口用于查询交付物详情信息，目前支持查询工作项交付物（节点交付物待建设）。对应的权限申请在权限管理-工作项实例分类下，相关功能介绍详见 获取权限 。

请求说明

- 请求方式：POST

- 请求地址：/open_api/work_item/deliverable/batch_query

请求参数

Body

| 参数 | 类型 | 是否必须 | 说明 |
| --- | --- | --- | --- |
| project_key | string | 必须 | 空间 ID（project_key），在飞书项目空间双击空间名称获取。 |
| work_item_ids | list<int64> | 可选 | 交付物工作项 ID 列表，限制10个，当查询 **工作项交付物** 使用此字段查询； **不传入** 则查询结果为空。 |

返回参数

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| data | list<object> | 遵循 [Deliverable](https://project.feishu.cn/b/helpcenter/1p8d7djs/1x1d372l#871e29be) 结构规范 |
| err | object | 请求成功时返回空值，请求失败时返回实际错误信息。 |
| err_msg | string | 请求成功时返回空值，请求失败时返回实际错误信息。 |
| err_code | int32 | 请求成功时返回空值，请求失败时返回实际错误信息。 |

请求示例

```json
{
  "project_key": "64fe8d784adf35b81ee6****",
  "work_item_ids": [434483****]
}
```

返回示例

请求成功示例

错误返回示例

```json
{
  "err_code": 0,
  "err_msg": "",
  "err": {},
  "data": [{
    "deliverable_uuid": "80317000",
    "deliverable_type": "work_item_deliverable",
    "deliverable_info": {
      "name": "f5-a-交付1",
      "work_item_id": 80317000,
      "instance_linked_virtual_resource_workitem": 80082946,
      "template_resources": false,
      "template_type": "662a39ded9e35a2a810199e0_template_1716791427582200564",
      "deleted": false,
      "delivery_related_info": {
        "root_work_item": {
          "project_key": "662a39ded9e35a2a8101xxxx",
          "work_item_id": 80316997,
          "work_item_type_key": "6693cea148a00c27960dxxxx",
          "name": "f5-a"
        },
        "source_work_item": {
          "project_key": "662a39ded9e35a2a8101xxxx",
          "work_item_id": 80316999,
          "work_item_type_key": "6693cea148a00c27960dxxxx",
          "name": "sub"
        }
      }
    }
  }]
}
```

错误码

以下列出了该接口相关的错误码信息，若需查看更多通用错误码说明，参见 [Open API 错误码](https://project.feishu.cn/b/helpcenter/1p8d7djs/5aueo3jr) 。

| err_code | err_msg | err.code | err.msg | 说明 |
| --- | --- | --- | --- | --- |
| 1000053287 | 传入的工作项实例不是交付物类型，请传入交付物类型的工作项实例 | 1000053287 | 传入的工作项实例不是交付物类型，请传入交付物类型的工作项实例 | 实例非交付物类型 |
