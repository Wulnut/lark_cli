# TUI Work Item Full Flow Design

## Context

当前代码库已经具备以下基础能力：

- `lark` 启动后的 Bubble Tea TUI 框架；
- 项目列表、项目详情、工作项类型列表的基本交互；
- 工作项查询的 CLI 与 search builder 初步状态；
- OpenAPI 用户、项目、工作项类型、搜索能力的部分封装。

当前缺口主要有两类：

1. **TUI 主链路不完整**：还没有完整覆盖“工作项列表 → 工作项详情 → 搜索 → 搜索结果 → 搜索结果详情”；
2. **CLI 侧存在明显可用性问题**：例如 `work-item-type list` 的参数体验与入口设计还不合理。

本轮先聚焦 **TUI 优先**，先把交互主链路打通，再回头系统修复 CLI。

## Goals

本阶段目标：

1. 支持两个入口进入同一套工作项 TUI 流程：
  - `lark`
  - `lark project work-item -t`
2. 在 TUI 中打通完整链路：
  - 项目列表
  - 项目详情
  - 工作项类型列表
  - 工作项列表
  - 工作项详情
  - Search Builder
  - Search Results
  - Search Result Detail
3. 工作项详情默认展示核心字段，同时支持查看原始 JSON / 全字段；
4. 保持“简约但不简单”的交互风格，统一键位和返回逻辑；
5. 不在本轮顺手修 CLI，避免同时处理两条主线。

## Non-Goals

本阶段不做：

- CLI 参数体系修复；
- 跨空间工作项聚合浏览；
- search builder 的复杂表达式编辑器；
- 结果页多维排序、批量操作、内联编辑；
- 过度抽象成多 tea.Model 架构。

## Entry Strategy

完整 TUI 流支持两个入口：

### Entry A: root TUI

用户执行 `lark` 后进入项目列表，再逐级钻取。

### Entry B: work-item TUI shortcut

用户执行 `lark project work-item -t` 后直接进入工作项相关主流程。

设计要求：

- 两个入口最终都落到同一个 `rootModel` 状态机；
- 不复制页面逻辑；
- 只改变初始 state 与上下文装配方式。

## Architecture Decision

采用 **方案 B**：保留单一 `rootModel`，但拆分模块文件。

### Why not A

继续把所有逻辑堆在 `internal/tui/root.go` 中虽然最快，但文件会继续膨胀，后续工作项列表、详情、搜索、结果详情接入后会更难维护。

### Why not C

直接拆成多个 tea model 虽然架构更纯，但当前阶段改动面过大，不利于优先完成完整链路。

### Chosen approach

- `rootModel` 继续作为唯一状态容器；
- `Update()` 继续作为统一事件分发入口；
- 各页面的渲染逻辑、页面级辅助行为拆到独立文件；
- 公共样式统一收口到样式文件。

## State Model

`rootModel` 保留统一状态，并扩展/整理以下 state：

- `stateProjectList`
- `stateProjectDetail`
- `stateWorkItemTypeList`
- `stateWorkItemList`
- `stateWorkItemDetail`
- `stateWorkItemSearchBuilder`
- `stateWorkItemSearchResults`
- `stateWorkItemSearchResultDetail`

### Additional rootModel data

为支持完整链路，需要增加以下上下文字段：

- 当前选中的 `projectKey`；
- 当前选中的 `workItemTypeKey`；
- 当前工作项列表数据；
- 当前选中的工作项摘要；
- 当前工作项详情数据；
- 当前 search results；
- 当前 detail 页面来源（list / search results）；
- 当前 detail 页面视图模式（summary / raw）。

## File Layout

建议文件边界如下：

- `internal/tui/root.go`
  - state 定义
  - `rootModel` 定义
  - 总入口 `Init()` / `Update()` / `View()` 分发
  - 通用 helper
  - `Run()` / 新入口运行函数
- `internal/tui/project_list.go`
  - 项目列表渲染与列表行为 helper
- `internal/tui/project_detail.go`
  - 项目详情渲染与项目详情页面行为 helper
- `internal/tui/work_item_type.go`
  - 工作项类型列表渲染、筛选、选择行为
- `internal/tui/work_item_list.go`
  - 工作项列表渲染、导航、过滤、进入详情
- `internal/tui/work_item_detail.go`
  - 工作项详情渲染、摘要/原始视图切换
- `internal/tui/work_item_search.go`
  - search builder
  - search results
  - result detail 相关逻辑
- `internal/tui/style.go`
  - 颜色、标题、卡片、footer、选中态、辅助样式

## User Flow

### Flow 1: browse by type

1. 进入项目列表；
2. 进入项目详情；
3. 打开工作项类型列表；
4. 选择某个 work item type；
5. 拉取并展示工作项列表；
6. 进入工作项详情；
7. 返回工作项列表。

### Flow 2: search from project detail

1. 在项目详情页打开 search builder；
2. 编辑查询条件；
3. 执行搜索；
4. 浏览搜索结果；
5. 进入搜索结果详情；
6. 返回搜索结果或返回 builder refine。

### Flow 3: direct work-item entry

1. 从 `lark project work-item -t` 进入；
2. 装配必要 project/type 上下文；
3. 进入同一套工作项 TUI 流。

## Entry B Bootstrap Contract

`lark project work-item -t` 作为 TUI 快捷入口，只做 **入口接线**，不在本轮顺手做整套 CLI UX 重构。

### Required context

- `project_key` 对于“direct jump”是必需上下文；
- `work_item_type_key` 为可选上下文；
- 如果调用入口时未提供 `project_key`，则不视为命令错误，而是降级为 interactive bootstrap flow。

### Bootstrap behavior

- 若已提供 `project_key`：直接进入该项目下的工作项类型列表；
- 若同时提供 `work_item_type_key`：直接进入该类型对应的工作项列表；
- 若缺少 `project_key`：进入项目列表让用户选择，而不是直接报错退出；
- 若 `project_key` 无效：进入项目错误页，`q` 或 `r` 之后的确定性回退目标为项目列表；
- 若 `project_key` 有效但 `work_item_type_key` 无效：进入类型错误页，`q` 或 `r` 之后的确定性回退目标为该项目下的工作项类型列表。

### Scope guard

这一入口在本轮只负责让用户更快进入同一套 TUI 工作项流，不负责同时解决 CLI 参数命名、位置参数兼容、help 文案等系统性问题。

## Search Builder Input Semantics

Search Builder 明确分为两种模式：

### Structured mode

由以下字段生成查询：

- project key
- work item type key
- me
- statuses
- created from
- created to

行为：

- `me=true` 时编译到结构化查询；
- `statuses` 采用逗号分隔；
- 日期字段为空时不参与编译；
- 输入非法日期时，不执行查询，直接在当前页展示可诊断错误。

### Raw JSON mode

当 `raw json` 非空时，进入 raw 模式。

行为约定：

- **raw json 优先**；
- 为避免用户误判，本轮采用简单规则：`raw json` 非空时，structured 字段不再参与合并；
- 页面上明确显示当前是 `raw` 模式，避免误以为两者会自动 AND 合并；
- raw json 非法时，不执行查询，直接在 builder 页展示错误。

这样做的原因是：CLI 当前有双层表达，但 TUI 第一阶段优先强调可解释和稳定，避免在界面里引入隐式合并行为。

## Pagination & Performance

工作项列表与搜索结果都必须有分页策略，避免一次性加载过多数据导致 TUI 卡顿。

### Initial policy

- 默认页大小：`20`
- 首屏只加载第一页
- 页面底部展示 `page_num / total` 或“当前已加载条数”

### Navigation policy

- 先不做无限滚动；
- 使用显式键位加载下一页/上一页（如 `n` / `p`，具体键位在实现阶段确定并统一写入 footer）；
- 翻页后保持列表页状态一致，避免重置过滤条件；
- model 中显式保存 `pageNum`、`pageSize`、`total`（若接口提供）、`hasNextPage`。

### API pagination contract

- work item list 使用接口返回的 `pagination.page_num`、`pagination.page_size`、`pagination.total`；
- search results 同样优先使用接口返回分页信息；
- 若接口未提供 `total`，则按“返回条数是否达到 pageSize 且本次请求非空”推导 `hasNextPage`；
- retry 必须复用同一组分页参数，不能悄悄重置到第一页。

### Empty/end behavior

- 无结果：显示空态；
- 到最后一页：footer 明确提示没有更多数据；
- API 不返回 total 时，使用“是否还有下一页”语义展示。

## Update Ownership Boundaries

虽然保留单一 `rootModel`，但必须限制 `Update()` 的职责。

### root.go owns

- 全局 state 定义；
- 总消息分发；
- 页面间跳转；
- 通用共享 helper；
- 入口初始化。

### page modules own

每个页面文件负责：

- 本页面的 view 渲染；
- 本页面局部键位处理 helper；
- 本页面专属消息处理 helper；
- 本页面 state 的最小变更逻辑。

### Rule

`root.go` 不直接堆叠每个页面的细节分支。实现上应采用类似：

- `updateProjectList(...)`
- `updateWorkItemTypeList(...)`
- `updateWorkItemList(...)`
- `updateWorkItemDetail(...)`
- `updateWorkItemSearch(...)`

由 `root.Update()` 统一路由到这些 helper。

## State Restoration Rules

返回行为必须稳定且可预测。

### list -> detail -> list

返回 list 时恢复：

- cursor
- filter query
- 当前页码
- scroll offset

detail 视图模式恢复规则：

- 每次重新进入 detail 默认回到 summary 视图；
- 不把上一次 raw 视图状态泄漏到下一条 work item。

### builder -> results -> detail -> results

返回 results 时恢复：

- cursor
- filter query（如果 results 页后续支持）
- 当前页码
- scroll offset

### builder -> results -> builder

返回 builder 时恢复：

- 所有已输入字段
- 当前 focus
- 最近一次 builder 错误信息，直到用户继续编辑或重新提交

## Loading and Error UX Contract

所有异步 API 页面都要有明确 loading 与 error 契约。

### Loading

以下动作触发 loading view：

- 拉取工作项类型列表
- 拉取工作项列表
- 拉取工作项详情
- 执行 search
- 拉取 search result detail

要求：

- loading 文案必须说明当前在加载什么；
- loading 状态下保留最少必要上下文（例如当前 project/type）；
- 不显示空列表占位来伪装 loading。

### Error

错误页必须满足：

- 展示明确错误文案；
- 支持 `q` 返回上一层；
- 支持 `r` 重试当前请求；
- 不丢失触发本次请求所需的上下文。

### Retry request ownership

为支持 `r`，model 需要保留“最后一次可重试请求”的最小上下文：

- 请求类型（type list / work item list / detail / search / search detail）
- project key
- work item type key
- work item id（如有）
- search payload（如有）
- pageNum / pageSize

`r` 只能重放最近一次失败请求，不能猜测用户意图或改写参数。

## Keyboard Precedence Rules

输入模式下必须定义键位优先级，避免全局快捷键误伤文本编辑。

### Builder field editing

当焦点在可编辑字段上时：

- 普通字符输入优先写入字段；
- `backspace` 优先删除字段内容；
- `tab` 切换焦点；
- `enter` 提交查询；
- `q` 不作为文本输入，仍执行返回；
- `j/k` 在 builder 文本编辑模式下默认视为普通字符，除非当前焦点字段明确是非文本型布尔切换字段。

### Filter mode

当列表页进入 `/` 过滤模式后：

- 普通字符写入 filter query；
- `backspace` 删除 query；
- `enter` / `esc` 退出过滤模式；
- `j/k` 在过滤模式下优先用于导航还是输入，必须与现有实现保持一致；本轮延续当前约定：过滤模式中 `j/k` 继续保留导航语义，不写入 query。

## Screen Design

### Project Detail

作为导航枢纽。保留当前已有信息展示，并新增清晰入口提示：

- 进入工作项类型列表；
- 打开 search builder。

### Work Item Type List

保留当前已有：

- `j/k` / `↑/↓` 导航；
- `/` 过滤；
- `q` 返回。

新增：

- `enter` 进入当前 type 对应的工作项列表。

### Work Item List

列表行优先展示核心信息：

- `[id] 标题`
- 状态
- 负责人（若可从摘要取到）

行为：

- `j/k` 或 `↑/↓` 移动；
- `/` 本地过滤；
- `enter` 进入 detail；
- `q` 返回工作项类型列表或项目详情。

### Work Item Detail

采用双视图：

#### Summary view

优先展示：

- ID
- 标题
- 状态
- 负责人
- 创建时间
- 更新时间
- 描述

#### Raw view

展示原始 JSON / 全字段展开结果。

行为：

- `tab` 切换 summary/raw；
- `q` 返回来源页面；
- 从 list 进入就回 list；从 search results 进入就回 results。

### Search Builder

保留轻量编辑体验，不在本轮过度复杂化。支持字段：

- project key
- work item type key
- me
- statuses
- created from
- created to
- raw json

行为：

- `tab` 切焦点；
- 输入编辑当前字段；
- `backspace` 删除；
- `enter` 执行搜索；
- `q` 返回项目详情；
- 页面底部展示最终查询摘要，帮助用户理解本次搜索将发送什么。

### Search Results

展示：

- `[id] title`
- 状态
- 负责人（能取到则展示）

行为：

- `j/k` 导航；
- `enter` 进入 result detail；
- `q` 回 builder；
- 保持 cursor，不因返回而丢失上下文。

### Search Result Detail

不单独维护完全不同的 detail 逻辑。

策略：

- 复用 `Work Item Detail` 的渲染能力；
- 只记录 detail 来源是 `search results`；
- `q` 返回到 results。

## API Strategy

本阶段所需 API：

### 1. Work item list

优先接入单空间列表接口：

- `POST /open_api/:project_key/work_item/filter`

用途：

- 按 `work_item_type_keys` 拉取某一类型的工作项列表；
- 作为 TUI 浏览型入口。

### 2. Work item detail

接入详情接口：

- `POST /open_api/:project_key/work_item/:work_item_type_key/query`

用途：

- 根据 `work_item_ids` 获取单个或少量工作项详情；
- 为 detail 页提供完整字段。

### 3. Search

继续使用现有 search builder 对应搜索能力：

- 已有 `search/params` 路径封装
- 继续沿用当前 query builder

## Data Flow

### Browse list path

1. 选定 project；
2. 拉取 work item types；
3. 选定 type；
4. 调用 work item list API；
5. 用户在列表页选中某项；
6. 调用 detail API；
7. 保存 detail + `detailSource=list`。

### Search path

1. 在 builder 中构造查询条件；
2. 生成 payload；
3. 调用 search API；
4. 用户在 results 中选中某项；
5. 依据结果中的 `project_key` / `work_item_type_key` / `id` 请求 detail；
6. 保存 detail + `detailSource=searchResults`。

## Error Handling

采用简洁且可恢复的策略：

- 页面级请求失败：渲染当前页面对应 error view；
- 所有错误页保留 `q` 返回上一层；
- 空数据明确展示 `No work items found.` 或等价文案；
- detail 拉取失败时，不吞错、不静默降级，直接展示 detail error view；
- 不在本轮实现复杂重试、后台刷新或局部 toast 系统。

## Keyboard Consistency

统一交互规则：

- `j/k` 或 `↑/↓`：上下移动
- `enter`：进入/执行
- `q`：返回上一层，最外层退出
- `tab`：切换焦点 / 切换 detail 视图
- `/`：进入过滤模式（适用于列表页）

所有页面 footer 必须显式说明当前页面可执行操作。

## Testing Strategy

坚持 TDD，测试放在现有偏好目录下。

### TUI tests

优先补以下测试：

1. 从项目详情进入工作项类型列表；
2. 从工作项类型列表进入工作项列表；
3. 从工作项列表进入工作项详情；
4. 从 search results 进入 result detail；
5. detail 页 `tab` 切换 summary/raw；
6. detail 页 `q` 按来源返回正确页面；
7. work item list 的 `j/k` 导航正确；
8. work item list 的过滤行为稳定；
9. search results 返回后保持 cursor；
10. 空列表、空结果、请求失败时视图稳定。

### OpenAPI tests

新增/补充：

- work item list API path/payload/response parsing；
- work item detail API path/payload/response parsing；
- detail 查询字段选择模式；
- 异常响应的错误处理。

## Implementation Phases

### Phase 1

先整理 TUI 文件边界，不改变已有行为：

- 拆出 project/work item type/style 文件；
- 保持现有测试继续通过。

### Phase 2

新增 work item list：

- API 封装
- state
- 页面渲染
- 进入 detail 的导航

### Phase 3

新增 work item detail：

- summary/raw 双视图
- detail source 返回逻辑

### Phase 4

打通 search result detail：

- 从 search results 进入 detail
- 复用 detail 页面能力

### Phase 5

打磨视觉与回退逻辑：

- footer 统一
- 空态/错误态统一
- 选中态和详情卡片样式统一

## Risks and Mitigations

### Risk 1: rootModel 继续变胖

Mitigation:

- 即便保留单一 model，也要拆出页面文件与 helper；
- 把页面级逻辑留在模块文件中，减少 root.go 的膨胀速度。

### Risk 2: search result detail 缺少必要上下文

Mitigation:

- 在结果行中保留并解析 `project_key`、`work_item_type_key`、`id`；
- 若缺少必要字段，直接提示无法打开详情，而不是猜测。

### Risk 3: detail 字段结构过于动态

Mitigation:

- summary 仅抽取少量稳定字段；
- raw view 直接展示原始结构，避免过早强绑定完整模型。

## Validation

实现完成后，至少验证：

1. `go test ./internal/tui ./test/tuitest ./test/openapitest`
2. `go test ./...`
3. 手工路径：
  - `lark` → project list → detail → type list → item list → item detail
  - `lark` → project detail → search builder → results → result detail
  - `lark project work-item -t` → 进入同一工作项流

### Acceptance checklist

- 从两个入口都能进入同一条工作项 TUI 主链；
- 可从工作项类型列表进入工作项列表；
- 可从工作项列表进入工作项详情；
- 可从搜索结果进入详情；
- detail 支持 summary/raw 切换；
- `q` 返回时状态恢复符合设计；
- loading、error、empty state 都有明确表现；
- 列表与结果页不会一次性无界加载大量数据。

## Follow-up

TUI 主链路完成后，再单独开一轮修 CLI，重点处理：

- `work-item-type` 命令入口可用性；
- 参数与位置参数设计；
- `--project-key` 体验改进；
- `work-item` 相关 CLI 子命令整合。

