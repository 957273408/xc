# 选手数据选择与展示系统 - 实现计划

## 一、代码库研究结论

### 1.1 技术栈

* **后端**: Go 1.x + Gin Web框架 + GORM ORM (gin-vue-admin架构)

* **前端**: Vue 3 (Composition API) + Element Plus UI + Pinia 状态管理

* **数据库**: MySQL (GORM AutoMigrate自动建表)

* **代理服务**: 已实现签名认证的代理转发 (ProxyServiceApp.Forward)

### 1.2 现有架构模式

* **后端分层**: Model(数据模型) → Service(业务逻辑) → API(接口层) → Router(路由注册)

* **公开路由**: 通过 `PublicRouter` 组注册，跳过JWT鉴权（参考 `GetPublicWarScores`）

* **前端API调用**: `@/utils/request.js` 封装axios，`@/api/*.js` 定义接口函数

* **前端路由**: `web/src/router/index.js` 直接注册路径，使用懒加载组件

* **数据库初始化**: `server/initialize/ensure_tables.go` 的 `MigrateTable` 方法添加模型自动迁移

### 1.3 代理API调用方式

现有 `fetchWarInfoRaw` 函数通过 ProxyServiceApp.Forward 调用 `/xdc/get_info`：

```go
resp, err := ProxyServiceApp.Forward("GET", "/xdc/get_info", map[string]interface{}{
    "warId": warID,
})
```

响应数据包含 `playerInfoList` 数组，每个玩家对象包含 `nickName`、`totalKill` 等字段。

***

## 二、需要修改/新增的文件清单

### 2.1 后端文件

| 文件路径                                                | 操作     | 说明                  |
| --------------------------------------------------- | ------ | ------------------- |
| `server/model/example/player_selection.go`          | **新增** | 选手选择数据模型            |
| `server/model/example/request/player_selection.go`  | **新增** | 请求参数结构体             |
| `server/model/example/response/player_selection.go` | **新增** | 响应数据结构体             |
| `server/service/example/player_selection.go`        | **新增** | 选手选择业务逻辑服务          |
| `server/api/v1/example/player_selection.go`         | **新增** | API接口处理函数           |
| `server/router/example/player_selection.go`         | **新增** | 路由注册                |
| `server/model/example/enter.go`                     | 修改     | 注册模型变量引用            |
| `server/service/example/enter.go`                   | 修改     | 注册服务变量引用            |
| `server/api/v1/example/enter.go`                    | 修改     | 注册API变量引用           |
| `server/router/example/enter.go`                    | 修改     | 注册路由初始化器            |
| `server/initialize/ensure_tables.go`                | 修改     | 添加新模型到AutoMigrate列表 |
| `server/initialize/router_biz.go`                   | 修改     | 注册新路由组到初始化器         |

### 2.2 前端文件

| 文件路径                                             | 操作     | 说明              |
| ------------------------------------------------ | ------ | --------------- |
| `web/src/api/playerSelection.js`                 | **新增** | 前端API调用函数定义     |
| `web/src/view/example/playerSelection/index.vue` | **新增** | 主页面组件           |
| `web/src/pinia/modules/playerSelection.js`       | **新增** | Pinia状态管理模块（可选） |
| `web/src/router/index.js`                        | 修改     | 添加新页面路由         |

***

## 三、详细修改步骤

### 阶段1：后端数据模型设计

#### Step 1.1: 创建 PlayerSelection 模型

**文件**: `server/model/example/player_selection.go`

```go
type PlayerSelection struct {
    global.GVA_MODEL
    WarID            string    `json:"warId" gorm:"column:war_id;size:64;not null;index"`
    SelectedPlayerIDs string   `json:"selectedPlayerIds" gorm:"column:selected_player_ids;type:text;not null"` // JSON数组: ["id1","id2",...]
    ExtraStat        string    `json:"extraStat" gorm:"column:extra_stat;size:32;not null"` // 可选值: healing/movement/throwables/identity_card/longest_kill
    SessionKey       string    `json:"sessionKey" gorm:"column:session_key;size:64;not null;uniqueIndex"` // warId + 随机后缀，用于会话间恢复
}
```

* **要点**:

  * SessionKey 用于无鉴权场景下的数据标识（前端生成并保存到localStorage）

  * SelectedPlayerIDs 使用 JSON 数组字符串存储5个玩家ID

  * WarID 建立索引以支持快速查询

#### Step 1.2: 创建 Request/Response 结构体

**Request** (`server/model/example/request/player_selection.go`):

```go
type SavePlayerSelectionRequest struct {
    WarID             string   `json:"warId" binding:"required"`
    SelectedPlayerIDs []string `json:"selectedPlayerIds" binding:"required,len=5"`
    ExtraStat         string   `json:"extraStat" binding:"required,oneof=healing movement throwables identity_card longest_kill"`
    SessionKey        string   `json:"sessionKey"`
}

type GetPlayerSelectionRequest struct {
    SessionKey string `form:"sessionKey" binding:"required"`
}

type GetWarPlayersRequest struct {
    WarID string `form:"warId" binding:"required"`
}
```

**Response** (`server/model/example/response/player_selection.go`):

```go
type PlayerInfo struct {
    PlayerID       string  `json:"playerId"`
    NickName       string  `json:"nickName"`
    KillCount      int     `json:"killCount"`      // 淘汰数
    HeadshotRate   float64 `json:"headshotRate"`   // 爆头率
    AccuracyRate   float64 `json:"accuracyRate"`   // 命中率
    DamageAmount   int     `json:"damageAmount"`   // 伤害量
    // 可选附加属性
    HealingAmount    *int     `json:"healingAmount,omitempty"`
    MovementDistance *float64 `json:"movementDistance,omitempty"`
    ThrowablesUsed   *int     `json:"throwablesUsed,omitempty"`
    IdentityCardUsed *int     `json:"identityCardUsed,omitempty"`
    LongestKillDist  *float64 `json:"longestKillDist,omitempty"`
}

type WarPlayersResponse struct {
    WarID   string       `json:"warId"`
    Players []PlayerInfo `json:"players"`
}

type PlayerSelectionResponse struct {
    ID                uint     `json:"id"`
    WarID             string   `json:"warId"`
    SelectedPlayerIDs []string `json:"selectedPlayerIds"`
    ExtraStat         string   `json:"extraStat"`
    SessionKey        string   `json:"sessionKey"`
    CreatedAt         string   `json:"createdAt"`
}
```

#### Step 1.3: 注册模型到 enter.go 和 ensure\_tables.go

* `server/model/example/enter.go`: 添加 `var ExaPlayerSelection = &PlayerSelection{}`

* `server/initialize/ensure_tables.go`: 在 tables 数组中添加 `example.PlayerSelection{}`

***

### 阶段2：后端业务逻辑层 (Service)

**文件**: `server/service/example/player_selection.go`

#### Step 2.1: GetWarPlayers 方法（获取并解析玩家数据）

* 复用 `ProxyServiceApp.Forward` 调用 `/xdc/get_info`

* 从 `playerInfoList` 提取玩家数据

* **字段映射策略**（根据实际响应字段调整）：

  * `totalKill` → KillCount

  * `headshotRate` / `headshots/totalKill*100` → HeadshotRate

  * `hitRate` / `hits/shots*100` → AccuracyRate

  * `totalDamage` → DamageAmount

  * 附加字段按 ExtraStat 选项过滤提取

* 返回标准化的 PlayerInfo 数组

* 实现 Redis/内存缓存（key: `war_players:{warId}`，TTL 5分钟），避免重复调用代理API

#### Step 2.2: SavePlayerSelection 方法（保存选择）

* 前端校验：确保 exactly 5 个玩家ID

* SessionKey 处理：

  * 如果前端不传，生成 `warId + "_" + 8位随机字符串`

  * 如果前端传入，使用 Upsert 语义（更新该 session 下的记录）

* 并发安全：使用 GORM 事务 + 唯一索引冲突处理

* 返回保存后的完整数据（含 sessionKey）

#### Step 2.3: GetPlayerSelection 方法（按SessionKey读取）

* 按 SessionKey 查询数据库

* 返回 404 风格错误（"selection not found"）如果不存在

#### Step 2.4: 注册服务到 enter.go

`server/service/example/enter.go`: 添加 `var PlayerSelectionServiceApp = new(PlayerSelectionService)`

***

### 阶段3：后端 API 接口层

**文件**: `server/api/v1/example/player_selection.go`

三个公开接口（全部注册到 PublicRouter，无需鉴权）：

| 方法   | 路径                                      | 功能                  | 响应码             |
| ---- | --------------------------------------- | ------------------- | --------------- |
| GET  | `/playerSelection/warPlayers?warId=xxx` | 获取某WarId下所有玩家列表     | 200 / 400       |
| POST | `/playerSelection/save`                 | 保存玩家选择（5人+附加统计）     | 200 / 400 / 500 |
| GET  | `/playerSelection/get?sessionKey=xxx`   | 按SessionKey读取已保存的选择 | 200 / 400 / 404 |

* 每个接口使用 ShouldBindQuery / ShouldBindJSON 参数绑定

* 统一使用 `response.OkWithDetailed` / `response.FailWithMessage` 返回格式

* 错误通过 `global.GVA_LOG.Error` 记录日志

* 注册API到 `server/api/v1/example/enter.go`

***

### 阶段4：后端路由注册

**文件**: `server/router/example/player_selection.go`

```go
func (p *PlayerSelectionRouter) InitPlayerSelectionRouter(Router, PublicRouter *gin.RouterGroup) {
    publicGroup := PublicRouter.Group("playerSelection")
    {
        publicGroup.GET("warPlayers", playerSelectionApi.GetWarPlayers)
        publicGroup.POST("save", playerSelectionApi.SavePlayerSelection)
        publicGroup.GET("get", playerSelectionApi.GetPlayerSelection)
    }
}
```

* 注册到 `server/router/example/enter.go`

* 在 `server/initialize/router_biz.go` 中调用初始化

***

### 阶段5：前端API封装

**文件**: `web/src/api/playerSelection.js`

```javascript
import service from '@/utils/request'

// 使用 axios 直连（不走 JWT 拦截器）或使用公开前缀
const publicService = service // 复用现有request（公开接口无需特殊处理，后端已跳过鉴权）

export const getWarPlayers = (params) => {
  return publicService({
    url: '/playerSelection/warPlayers',
    method: 'get',
    params
  })
}

export const savePlayerSelection = (data) => {
  return publicService({
    url: '/playerSelection/save',
    method: 'post',
    data
  })
}

export const getPlayerSelection = (params) => {
  return publicService({
    url: '/playerSelection/get',
    method: 'get',
    params
  })
}
```

***

### 阶段6：前端页面组件

**文件**: `web/src/view/example/playerSelection/index.vue`

#### 页面结构（参考设计图风格，4人卡片+数据统计）：

```
┌─────────────────────────────────────────────────────────────┐
│  页面头部：SMC 2026 品牌装饰 + 标题                          │
├─────────────────────────────────────────────────────────────┤
│  WarId 输入区                                                │
│  [输入框: 请输入WarId]  [获取数据按钮(loading)]              │
├─────────────────────────────────────────────────────────────┤
│  附加统计下拉选择                                            │
│  [治疗量/移动距离/投掷物数/身份卡数/最远击杀距离]             │
├─────────────────────────────────────────────────────────────┤
│  玩家列表区（表格/卡片，支持多选，限制最多5人）                │
│  [x] 玩家A  淘汰:235  爆头率:65%  命中率:42%  伤害:12000    │
│  [✓] 玩家B  ...                                              │
│  ...                                                         │
│  已选择: 3/5  [保存按钮(disabled until 5人)]                │
├─────────────────────────────────────────────────────────────┤
│  预览/已保存结果区域（5张卡片，对应设计图4张风格 + 1张）      │
│  ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐   │
│  │ 头像框 │ │ 头像框 │ │ 头像框 │ │ 头像框 │ │ 头像框 │   │
│  │ PLAYER │ │ PLAYER │ │ PLAYER │ │ PLAYER │ │ PLAYER │   │
│  │ 淘汰数 │ │ 淘汰数 │ │ 淘汰数 │ │ 淘汰数 │ │ 淘汰数 │   │
│  │ 爆头率 │ │ 爆头率 │ │ 爆头率 │ │ 爆头率 │ │ 爆头率 │   │
│  │ 命中率 │ │ 命中率 │ │ 命中率 │ │ 命中率 │ │ 命中率 │   │
│  │ 伤害量 │ │ 伤害量 │ │ 伤害量 │ │ 伤害量 │ │ 伤害量 │   │
│  │附加项4 │ │附加项4 │ │附加项4 │ │附加项4 │ │附加项4 │   │
│  └────────┘ └────────┘ └────────┘ └────────┘ └────────┘   │
│                                                             │
│              TEAM [战队名称]  CHAMPIONSHIP                   │
└─────────────────────────────────────────────────────────────┘
```

#### 核心逻辑：

1. **加载流程**:

   * 页面加载时检查 localStorage 中的 sessionKey

   * 如果有 sessionKey → 自动调用 `getPlayerSelection` 恢复数据

   * 恢复时自动填充 WarId、附加统计选项、选中玩家

2. **获取玩家数据**:

   * 输入 WarId → 点击获取 → 显示 loading → 调用 `getWarPlayers`

   * 错误处理：400参数错误 / 代理API超时 / 数据为空

   * 成功后渲染玩家列表表格，重置已选择

3. **玩家选择**:

   * 表格多选模式（el-table type="selection"）

   * `@selection-change` 事件中限制 selection.length <= 5

   * 超过5个时提示并自动回滚最后一次选择

   * 实时显示选中数量计数器 `x/5`

4. **附加统计选择**:

   * el-select 下拉组件，选项：

     * `healing` → 治疗量

     * `movement` → 移动距离

     * `throwables` → 投掷物使用数

     * `identity_card` → 身份卡使用数

     * `longest_kill` → 最远击杀距离

   * 选择后实时刷新预览卡片的第5行数据

5. **保存功能**:

   * 校验 exactly 5 人 → 否则禁用按钮并提示

   * 调用 `savePlayerSelection`，如果无 sessionKey 则后端自动生成

   * 成功后：

     * 保存 sessionKey 到 localStorage（key: `playerSelectionSession_${warId}`）

     * 显示成功提示 + 展示预览卡片

     * 启用"分享链接"功能（复制 `?sessionKey=xxx` URL）

6. **响应式设计**:

   * 桌面端：5张卡片横排 (`grid-template-columns: repeat(5, 1fr)`)

   * 平板端：3+2 两行布局

   * 移动端：单列滚动 (`flex-direction: column`)

   * 玩家表格列数根据屏幕宽度动态隐藏次要列

7. **Loading & Error状态**:

   * 骨架屏组件（卡片加载时的灰色占位）

   * el-message / el-alert 错误提示

   * API失败时提供重试按钮

#### 状态管理（Pinia可选，复杂场景使用）：

**文件**: `web/src/pinia/modules/playerSelection.js`

```javascript
// state: warId, playersList, selectedIds, extraStat, sessionKey, loading states
// actions: fetchPlayers, saveSelection, restoreFromSession
```

#### 路由注册：

**文件**: `web/src/router/index.js`

```javascript
{
  path: '/playerSelection',
  name: 'PlayerSelection',
  meta: { title: '选手数据选择', client: true }, // client: true 表示无需登录即可访问
  component: () => import('@/view/example/playerSelection/index.vue')
}
```

***

### 阶段7：会话间数据持久化与恢复

**localStorage 存储策略**:

```
Key: playerSelectionSession_{warId}
Value: { sessionKey: "xxx", savedAt: 时间戳 }

Key: playerSelectionLastWarId
Value: "上次访问的warId"
```

**URL参数恢复**:

* 支持 `?sessionKey=xxx` 直接访问 → 自动读取并恢复

* 支持 `?warId=xxx` 直接访问 → 自动填充并获取数据

***

### 阶段8：测试用例

#### 后端单元测试（可选，放到 `_test.go` 文件）：

1. **TestSavePlayerSelection\_Valid5Players**: 保存正好5个玩家 → 成功
2. **TestSavePlayerSelection\_InvalidCount**: 保存4个或6个 → 返回400校验错误
3. **TestSavePlayerSelection\_InvalidExtraStat**: 非法附加统计选项 → 400
4. **TestGetPlayerSelection\_NotFound**: 不存在的sessionKey → 返回错误
5. **TestGetWarPlayers\_CacheHit**: 连续两次调用同一WarId → 缓存生效（可通过日志验证）
6. **TestUpsertSemantics**: 同sessionKey保存两次 → 更新而非创建新记录

#### 前端集成测试点（手动验证清单）：

1. WarId输入为空时 → 按钮禁用
2. 选择玩家超过5个 → 阻止并提示
3. 选择正好5个后 → 保存按钮启用
4. 保存成功后刷新页面 → 数据自动恢复
5. 带sessionKey参数直接访问 → 自动恢复选择
6. 切换附加统计选项 → 预览卡片数据实时更新
7. 调整浏览器窗口大小 → 卡片布局自适应（5列→2列→1列）
8. API返回错误 → 显示正确错误信息，可重试

***

## 四、潜在依赖与注意事项

### 4.1 字段名称映射风险

* `/xdc/get_info` 实际响应的玩家字段名未知（如爆头率字段可能是 `headshotRate` 或 `hsRate` 或需要计算 `headshots/kills`）

* **缓解措施**: 在 Service 层 GetWarPlayers 方法中实现多字段兼容探测（使用 switch/case 检查多种可能的字段名），并在数据解析失败时返回字段缺失的明确错误信息

### 4.2 并发写入一致性

* 多个用户同时使用同一sessionKey保存 → 后写覆盖先写（最后写入生效）

* **缓解措施**: 数据库层 sessionKey 唯一索引 + GORM `Clauses(clause.OnConflict{...})` Upsert 语法，确保原子性

### 4.3 代理API响应延迟

* `/xdc/get_info` 可能响应慢（>5秒）导致用户体验差

* **缓解措施**:

  * 前端超时设置（15秒）+ 加载动画 + 取消请求（AbortController）

  * 后端缓存策略（相同WarId 5分钟内直接返回缓存）

### 4.4 无鉴权数据安全

* 公开接口任何人可读写保存的选择数据

* **缓解措施**:

  * SessionKey 使用足够随机性（crypto/rand 16字节 → hex编码），难以被猜测

  * 不存储敏感信息（仅玩家ID和展示统计）

  * 不依赖sessionKey做权限校验，仅用于标识和恢复

### 4.5 前端直接访问公开接口

* 现有 `@/utils/request.js` 默认会在拦截器中附加 JWT token

* **确认**: 检查 request.js 拦截器逻辑，如果无token时不会报错则无需修改；否则创建公开专用的 axios 实例（不带 Auth 头）

***

## 五、风险处理汇总

| 风险          | 概率 | 影响      | 处理方案                              |
| ----------- | -- | ------- | --------------------------------- |
| 代理API字段名不匹配 | 高  | 玩家数据全为空 | Service层多字段兼容探测 + 错误日志包含原始响应key列表 |
| 保存失败会话丢失    | 中  | 用户重复操作  | 保存前先写入localStorage草稿，失败可恢复        |
| 并发选择冲突      | 低  | 数据不一致   | 唯一索引 + Upsert                     |
| 移动端5卡片过窄    | 中  | 显示不友好   | 响应式断点 + 横向滚动容器                    |
| 数据库连接失败     | 低  | 保存/读取失败 | 统一错误处理 + 重试机制（最多2次）               |

