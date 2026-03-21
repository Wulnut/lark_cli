---
title: "获取工作流详情（WBS） - 开发者手册 - 飞书项目帮助中心"
source: "https://project.feishu.cn/b/helpcenter/1p8d7djs/1t3boa28"
author:
published:
created: 2026-03-22
updated: 2026-03-22
description: "获取行业专版中节点流工作项实例的 WBS 工作流信息，包括节点、子任务、交付物等"
tags:
  - "飞书项目"
  - "工作项"
  - "流程与节点"
  - "WBS"
  - "API列表"
  - "OpenAPI"
category: "飞书项目"
related_docs:
  - "[[获取工作项详情 - 开发者手册 - 飞书项目帮助中心]]"
  - "[[获取工作流详情 - 开发者手册 - 飞书项目帮助中心]]"
---

该接口用于获取行业专版中一个节点流工作项实例的 WBS 工作流信息（WBS 对应能力参见 WBS流程配置），对应的权限申请在权限管理-工作项实例分类下，相关功能介绍详见权限管理。

请求说明

- 请求方式：GET

- 请求地址：/open_api/:project_key/work_item/:work_item_type_key/:work_item_id/wbs_view

请求参数

Path

| 参数 | 类型 | 是否必须 | 说明 |
| --- | --- | --- | --- |
| project_key | string | 必须 | 空间 ID（project_key），或者空间域名（simple_name）。  - 空间 ID（project_key），在飞书项目空间双击空间名称获取。 - 空间域名（simple_name），simple_name 一般在飞书项目空间 URL 中获取。例如空间 URL 为 "https://project.feishu.cn/doc/overview"，则 simple_name为 "doc"。 |
| work_item_type_key | string | 必须 | 工作项类型，系统自动生成，可以从 [获取空间下工作项类型](https://project.feishu.cn/b/helpcenter/1p8d7djs/3pjp854w) 接口获取。 |
| work_item_id | string | 必须 | 工作项实例 ID，在工作项实例详情右上角展开 **···** ，单击 ID 信息即可获取。 |

Query

| 参数 | 类型 | 是否必须 | 说明 |
| --- | --- | --- | --- |
| need_union_deliverable | boolean | 可选 | 是否融合需要交付物，与请求体参数中的 need_union_deliverable含义相同，在查询参数中冗余。 |
| need_schedule_table_agg | boolean | 可选 | 扩展查询，是否需要计划表自定义列聚合字段。 |

Body

| 参数 | 类型 | 是否必须 | 说明 |
| --- | --- | --- | --- |
| expand | object | 可选 | 额外参数：   need_union_deliverable：是否融合需要交付物   need_wbs_relation_chain_entity: 是否需要 WBS 链路实例   need_wbs_relation_chain_path: 是否需要 WBS 链路层级 |

返回参数

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| data | object | 遵循 [WbsViewResponse](https://project.feishu.cn/b/helpcenter/1p8d7djs/1x1d372l#eca9aa28) 数据结构规范。 |
| err | object | 请求成功时返回空值，请求失败时返回实际错误信息。 |
| err_msg | string | 请求成功时返回空值，请求失败时返回实际错误信息。 |
| err_code | int32 | 请求成功时0，请求失败返回实际错误码。 |

请求示例

查看节点流工作项详情

查看状态流带扩展字段工作项

查询节点流工作项详情，并查看 need_union_deliverable 和 need_schedule_table_agg。

GET /open_api/66c2fb3e0ab52f52c92****/work_item/story/5566778/wbs_view?need_union_deliverable=true&need_schedule_table_agg=true

返回示例

查看节点流工作项详情

查看状态流带扩展字段工作项

```json
{
  "err_code": 0,
  "data": {
    "template_key": "123123123",
    "related_sub_work_items": [
      {
        "node_uuid": "1",
        "work_item_id": 5566778,
        "type": "node",
        "name": "开始",
        "wbs_status_map": {
          "status_key": "not_started",
          "status_name": "未开始"
        },
        "sub_work_item": [
          {
            "deliverable": [
              {
                "field_key": "field_b46006",
                "field_value": [
                  [
                    {
                      "field_key": "field_cf394a",
                      "field_value": "2323"
                    },
                    {
                      "field_key": "field_db2cd2",
                      "field_value": "4545"
                    },
                    {
                      "field_key": "field_ec3f89",
                      "field_value": "rd"
                    }
                  ]
                ]
              }
            ],
            "name": "任务1",
            "node_uuid": "3",
            "work_item_id": 5566712,
            "type": "sub_task",
            "role_owners": [
              {
                "role": "DA",
                "name": "DA",
                "owners": [
                  "70439708518xxxxxx06",
                  "70970434698xxxxxx04"
                ]
              },
              {
                "role": "PM",
                "name": "PM"
              }
            ],
            "wbs_status_map": {
              "status_key": "not_started",
              "status_name": "未开始"
            },
            "sub_work_item": []
          }
        ],
        "union_deliverable": {
          "field_deliverables": [
            {
              "field_info": {
                "field_key": "name",
                "field_value": "名称",
                "field_value_type": "name"
              },
              "placeholder": "占位符",
              "remark": "标记",
              "status": 1
            }
          ],
          "instance_deliverables": [
            {
              "name": "交付物工作项名称",
              "work_item_id": 112333,
              "associated_deliverables_template": 1123111,
              "template_resources": false,
              "deletable": true,
              "must_complete": false,
              "state_key": "started",
              "state_name": "开始",
              "owners": ["username"],
              "remark": "标记"
            }
          ]
        },
        "related_parent_work_item": {
          "is_top": false,
          "work_item_id": 123456,
          "name": "test",
          "template_key": "637deb0exxxxxxxxxxxxxxx:story:template_123456789",
          "template_version": 1,
          "work_item_type_key": "story",
          "template_id": "1234",
          "template_name": "子流程"
        },
        "connections": [
          {
            "source_state_key": "state_0",
            "target_state_key": "state_1"
          }
        ]
      }
    ],
    "err_msg": ""
  }
}
```

错误码

以下列出了该接口相关的错误码信息，若需查看更多通用错误码说明，参见 [Open API 错误码](https://project.feishu.cn/b/helpcenter/1p8d7djs/5aueo3jr) 。

| err_code | err_msg | err.code | err.msg | 说明 |
| --- | --- | --- | --- | --- |
| 30005 | WorkItem Not Found | 30005 | WorkItem Not Found | 工作项未找到，即 work_item_id 不正确。 |
| 20089 | Current WorkItem Is Not Wbs Or Wbs Children | 20089 | Current WorkItem Is Not Wbs Or Wbs Children | 当前工作项不是 WBS 工作项或者不是 WBS 子工作项。 |
| 20014 | Project And WorkItem Not Match | 20014 | Project And WorkItem Not Match | 空间和工作项不匹配。 |
