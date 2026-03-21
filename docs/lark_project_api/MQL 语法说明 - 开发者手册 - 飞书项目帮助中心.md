---
title: "MQL 语法说明 - 开发者手册 - 飞书项目帮助中心"
source: "https://project.feishu.cn/b/helpcenter/1p8d7djs/tsl2uj3i"
author:
published:
created: 2026-03-19
updated: 2026-03-19
description: "飞书项目 MQL (Meego Query Language) 语法说明，是一种专用的结构化查询语言，兼容 SQL 语法，包含数据类型、函数和查询示例"
tags: [飞书项目, OpenAPI, MQL, 查询语法, SQL, 函数, 数组函数, 相对时间]
category: 飞书项目
related_docs:
  - "[[飞书项目API开发者知识库]]"
  - "[[全量搜索参数格式及常用示例 - 开发者手册 - 飞书项目帮助中心]]"
  - "[[数据结构汇总 - 开发者手册 - 飞书项目帮助中心]]"
  - "[[飞书项目OpenAPI完整API列表]]"
---

Meego Query Language（简称 MQL）是对飞书项目系统中的数据进行结构化查询的一种专用查询语言。基于该语法，飞书项目 MCP Server 中的 `search_by_mql` tool 提供了 MQL 入参查询实例的能力。

MQL 兼容 SQL 语法，同时为了满足飞书项目特有功能而新增了部分数据类型及函数。

## 基础语法规则

```sql
SELECT fieldList                -- 指定查询的字段列表
FROM objectType                -- 指定要查询的数据来源
WHERE conditionExpression       -- 指定查询条件
[ORDER BY fieldOrderByList [{ASC|DESC}] ]  -- 排序（可选）
[LIMIT [offset,] row_count]    -- 限制返回行数（可选）
```

## 数据类型

| 数据类型 | 说明 |
| --- | --- |
| bool | 表示真假，取值可以是 TRUE、FALSE、1、0 四种 |
| bigint | 整数类型 |
| double | 浮点数类型 |
| varchar | 字符串类型 |
| date | 日期类型，格式支持 `YYYY-MM-DD` 或带时区格式 `YYYY-MM-DD+TZD` |
| datetime | 日期时间类型，格式支持 ISO8601 格式 |
| array | 数组类型，元素可以是 bool、bigint、double、varchar、date、datetime 等基础类型 |
| lambda expression | 特殊函数表达式，返回值类型是 bool |

### Lambda 表达式写法示例

- `x -> x > 10`
- `x -> x > 10 and x < 100`
- `x -> x in ('a', 'b')`

## 支持函数

### 数组函数

| 函数作用 | 函数签名 | 入参说明 | 返回值说明 | 使用示例 |
| --- | --- | --- | --- | --- |
| 判断 array 中是否所有 element 都满足条件 | `all_match(array_col, predicate)` | array_col：array 列<br>predicate：lambda 表达式 | bool：TRUE=全部满足，FALSE=不全部满足 | `all_match(ary_col, x -> x > 10)` |
| 判断 array 中是否有一个 element 满足条件 | `any_match(array_col, predicate)` | array_col：array 列<br>predicate：lambda 表达式 | bool：TRUE=至少一个满足，FALSE=全部不满足 | `any_match(ary_col, x -> x = '用户 A')` |
| 判断 array 中是否所有 element 都不满足条件 | `none_match(array_col, predicate)` | array_col：array 列<br>predicate：lambda 表达式 | bool：TRUE=全部不满足，FALSE=至少一个满足 | `none_match(ary_col, x -> x = '用户 A')` |
| 返回 array 的元素个数 | `array_cardinality(array_col)` | array_col：array 列 | bigint：元素个数 | `array_cardinality(负责人)` |
| 判断 array 中是否包含 element | `array_contains(array_col, element)` | array_col：array 列<br>element：元素值 | bool：TRUE=包含，FALSE=不包含 | `array_contains(优先级, 'P0')` |
| 根据条件过滤 array，返回新数组 | `array_filter(array_col, predicate)` | array_col：array 列<br>predicate：lambda 表达式 | array：过滤后的新数组 | `array_filter(ary_col, x -> x > 100)` |

### 用户与团队函数

| 函数作用 | 函数签名 | 入参说明 | 返回值说明 | 使用示例 |
| --- | --- | --- | --- | --- |
| 表示当前登录用户 | `current_login_user()` | 无 | varchar：当前登录用户的 userkey | `current_login_user()` |
| 表示指定团队的成员 | `team(include_manager, team_name)` | include_manager：bool，是否包含管理者<br>team_name：varchar，团队名称 | array(varchar)：团队成员的 userkey | `team(true, '后端开发团队')` |
| 表示所有参与角色 | `participate_roles()` | 无 | array(varchar)：所有参与角色的 rolekey | `participate_roles()` |
| 表示所有参与人员 | `all_participate_persons()` | 无 | array(varchar)：所有参与人的 userkey | `all_participate_persons()` |

### 相对时间函数

| 函数作用 | 函数签名 | 入参说明 | 返回值说明 | 使用示例 |
| --- | --- | --- | --- | --- |
| 等于相对时间 | `RELATIVE_DATETIME_EQ(col_name, 'date_para', ['days'])` | 见下方说明 | bool | 见下方示例 |
| 大于相对时间 | `RELATIVE_DATETIME_GT(col_name, 'date_para', ['days'])` | 见下方说明 | bool | 见下方示例 |
| 大于等于相对时间 | `RELATIVE_DATETIME_GE(col_name, 'date_para', ['days'])` | 见下方说明 | bool | 见下方示例 |
| 小于相对时间 | `RELATIVE_DATETIME_LT(col_name, 'date_para', ['days'])` | 见下方说明 | bool | 见下方示例 |
| 小于等于相对时间 | `RELATIVE_DATETIME_LE(col_name, 'date_para', ['days'])` | 见下方说明 | bool | 见下方示例 |
| 属于相对时间范围 | `RELATIVE_DATETIME_BETWEEN(col_name, 'date_para', ['days'])` | 见下方说明 | bool | 见下方示例 |

**date_para 可取值：**

| 取值 | 说明 |
| --- | --- |
| today | 当天 |
| tomorrow | 明天 |
| yesterday | 昨天 |
| current_week | 当周 |
| next_week | 下周 |
| last_week | 上周 |
| current_month | 当月 |
| next_month | 下月 |
| last_month | 上月 |
| future | 未来 |
| past | 过去 |

**days 可选参数：** 仅在 date_para 等于 `future`、`past`、`today` 时有效，表示偏移天数

**相对时间函数示例：**

- `RELATIVE_DATETIME_EQ(创建时间, 'today')` - 创建时间等于今天
- `RELATIVE_DATETIME_EQ(创建时间, 'today', '3d')` - 创建时间等于今天后3天
- `RELATIVE_DATETIME_EQ(创建时间, 'today', '-3d')` - 创建时间等于今天前3天
- `RELATIVE_DATETIME_EQ(创建时间, 'tomorrow')` - 创建时间等于明天
- `RELATIVE_DATETIME_EQ(创建时间, 'current_week')` - 创建时间等于当周
- `RELATIVE_DATETIME_EQ(创建时间, 'next_month')` - 创建时间等于下月
- `RELATIVE_DATETIME_BETWEEN(创建时间, 'past', '30d')` - 创建时间在过去30天内

## 语义说明及使用示例

在语义表达上，MQL 的核心设计理念是将飞书项目内的各种元素抽象为统一的【对象】和【对象属性】：

- **对象**：指代工作项类型、工作项节点和子任务等实体
- **对象属性**：指代【对象】拥有的属性，例如系统字段、自定义字段和角色人员

## 查询需求工作项示例

假设需求工作项有如下属性：

| 飞书项目字段 | 对应数据类型 |
| --- | --- |
| 工作项 ID | bigint，需求的唯一标识 |
| 负责人 | array(varchar) |
| 优先级 | varchar |
| 估算耗时 | bigint |
| 是否完结 | bool |
| 角色 RD | array(varchar)，表示研发角色的人员 |
| 角色 QA | array(varchar)，表示测试角色的人员 |

### 查询示例

**筛选出用户1负责的、在 2025-01-01 到 2025-10-01 之间创建的、优先级是 P0 的、已经完结的需求：**

```sql
SELECT `工作项id`
FROM `空间x`.`需求`
WHERE array_contains(`负责人`, '用户1')
  AND `创建时间` between '2025-01-01' and '2025-10-01'
  AND `优先级` = 'P0'
  AND `是否完结` = true
```

或使用英文 key：

```sql
SELECT `work_item_id`
FROM `664abc167e11cb9b1f******`.`story`
WHERE array_contains(`owner`, '用户1')
  AND `create_time` between '2025-01-01' and '2025-10-01'
  AND `priority` = 'P0'
  AND `finish_status` = true
```

**筛选出在过去 30 天内创建、有 RD 和 QA 角色参与、且估算耗时大于 30 天的需求，按优先级降序排序，仅返回前 100 条：**

```sql
SELECT `优先级`, `工作项id`
FROM `空间x`.`需求`
WHERE RELATIVE_DATETIME_BETWEEN(`创建时间`, 'past', '30d')
  AND array_contains(participate_roles(), 'RD', 'QA')
  AND `估算耗时` > 30
ORDER BY `优先级` DESC
LIMIT 100
```

或使用英文 key：

```sql
SELECT `priority`, `work_item_id`
FROM `664abc167e11cb9b1f******`.`story`
WHERE RELATIVE_DATETIME_BETWEEN(`create_time`, 'past', '30d')
  AND array_contains(participate_roles(), 'RD', 'QA')
  AND `estimation_time` > 30
ORDER BY `priority` DESC
LIMIT 100
```

## 查询缺陷工作项示例

假设缺陷工作项有如下属性：

| 飞书项目字段 | 对应数据类型 |
| --- | --- |
| 工作项 ID | bigint，缺陷的唯一标识 |
| 缺陷名 | varchar |
| 创建时间 | datetime |
| 处理人 | array(varchar) |
| 开发周期_开始时间 | datetime，开发周期的起始时间 |
| 开发周期_结束时间 | datetime，开发周期的结束时间 |
| 角色 RD | array(varchar)，表示研发角色的人员 |

### 查询示例

**筛选出研发包括用户2的、昨天创建的、至少有一位处理人属于开放平台团队的、开发周期在 2025-01-01 到 2025-01-31 之间的缺陷：**

```sql
SELECT `工作项id`
FROM `空间x`.`缺陷`
WHERE array_contains(__RD, '用户1')
  AND RELATIVE_DATETIME_EQ(`创建时间`, 'yesterday')
  AND any_match(`处理人`, x -> x in (team(true, '开放平台团队')))
  AND `__开发周期_开始时间` > '2025-01-01'
  AND `__开发周期_结束时间` < '2025-01-31'
```

**筛选出缺陷名不包含 '后端性能问题'、参与人中包含 '用户2'，且处理人是当前登录用户的缺陷，按创建时间升序排列，仅返回前 10 条：**

```sql
SELECT `工作项id`
FROM `空间x`.`缺陷`
WHERE `缺陷名` not like '%后端性能问题%'
  AND array_contains(all_participate_persons(), '用户2')
  AND `处理人` = current_login_user()
ORDER BY `创建时间` ASC
LIMIT 10
```
