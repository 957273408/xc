# 战队信息管理页面实现计划

## 一、项目概述

基于现有的 gin-vue-admin 框架，开发一个功能完整的战队信息管理页面，包含战队数据存储、Excel 导入、WarId 管理、积分统计等功能。

## 二、技术栈

- **后端**: Go 1.24 + Gin + GORM + excelize
- **前端**: Vue 3 + Element Plus + Vite
- **数据库**: MySQL (当前配置)

## 三、实现步骤

### 3.1 后端开发

#### 3.1.1 新增数据模型 (Model)

**文件**: `server/model/example/competition_team.go`

创建 `CompetitionTeam` 模型，扩展现有 `ExaTeam`：

```go
type CompetitionTeam struct {
    global.GVA_MODEL
    TeamCode    string `json:"teamCode" form:"teamCode" gorm:"column:team_code;comment:战队标识(3-8位字母数字);size:8;not null"`
    TeamName    string `json:"teamName" form:"teamName" gorm:"column:team_name;comment:战队名称;size:50;not null"`
    TeamLogo    string `json:"teamLogo" form:"teamLogo" gorm:"column:team_logo;comment:战队Logo路径;size:500"`
    TotalScore  int    `json:"totalScore" form:"totalScore" gorm:"column:total_score;comment:总积分"`
}
```

**文件**: `server/model/example/team_score.go`

创建 `TeamScore` 模型，用于记录每场比赛积分：

```go
type TeamScore struct {
    global.GVA_MODEL
    TeamID      uint    `json:"teamId" form:"teamId" gorm:"column:team_id;comment:战队ID;not null;index"`
    WarID       string  `json:"warId" form:"warId" gorm:"column:war_id;comment:战场ID;size:100;not null;uniqueIndex"`
    Rank        int     `json:"rank" form:"rank" gorm:"column:rank;comment:排名"`
    KillCount   int     `json:"killCount" form:"killCount" gorm:"column:kill_count;comment:淘汰人数"`
    RankScore   int     `json:"rankScore" form:"rankScore" gorm:"column:rank_score;comment:排名分"`
    KillScore   int     `json:"killScore" form:"killScore" gorm:"column:kill_score;comment:淘汰分"`
    TotalScore  int     `json:"totalScore" form:"totalScore" gorm:"column:total_score;comment:总积分"`
    GameTime    float64 `json:"gameTime" form:"gameTime" gorm:"column:game_time;comment:游戏时长(秒)"`
    SettleTime  time.Time `json:"settleTime" form:"settleTime" gorm:"column:settle_time;comment:结算时间"`
}
```

#### 3.1.2 新增请求/响应结构

**文件**: `server/model/example/request/competition_team.go`

```go
// 创建战队请求
type CreateCompetitionTeamRequest struct {
    TeamCode string `json:"teamCode" binding:"required,alphanum,min=3,max=8"`
    TeamName string `json:"teamName" binding:"required,max=50"`
    TeamLogo string `json:"teamLogo"`
}

// 更新战队请求
type UpdateCompetitionTeamRequest struct {
    ID       uint   `json:"id" binding:"required"`
    TeamCode string `json:"teamCode" binding:"required,alphanum,min=3,max=8"`
    TeamName string `json:"teamName" binding:"required,max=50"`
    TeamLogo string `json:"teamLogo"`
}

// 添加WarId请求
type AddWarIDRequest struct {
    TeamID  uint   `json:"teamId" binding:"required"`
    WarID   string `json:"warId" binding:"required"`
}

// Excel导入请求
type ImportExcelRequest struct {
    Mode string `json:"mode" binding:"required,oneof=incremental full"` // incremental/full
}
```

**文件**: `server/model/example/response/competition_team.go`

```go
// 战队积分记录响应
type TeamScoreRecordResponse struct {
    ID         uint      `json:"id"`
    WarID      string    `json:"warId"`
    Rank       int       `json:"rank"`
    KillCount  int       `json:"killCount"`
    TotalScore int       `json:"totalScore"`
    SettleTime string    `json:"settleTime"`
}

// 战队详情响应
type TeamDetailResponse struct {
    Team      CompetitionTeam          `json:"team"`
    ScoreList []TeamScoreRecordResponse `json:"scoreList"`
    TotalScore int                     `json:"totalScore"`
}
```

#### 3.1.3 新增Service层

**文件**: `server/service/example/competition_team.go`

主要方法：
- `CreateCompetitionTeam(team)` - 创建战队
- `UpdateCompetitionTeam(team)` - 更新战队
- `DeleteCompetitionTeam(id)` - 删除战队
- `GetCompetitionTeam(id)` - 获取单个战队
- `GetCompetitionTeamList(pageInfo)` - 获取战队列表
- `ImportTeamsFromExcel(file, mode)` - 从Excel导入战队
- `AddWarID(teamID, warID)` - 添加WarId并获取积分
- `GetTeamScores(teamID)` - 获取战队所有积分记录
- `GetTeamRecentScores(teamID, limit)` - 获取最近N次积分记录
- `CalculateScore(rank, killCount)` - 根据赛事规则计算积分
- `RecalculateTeamTotalScore(teamID)` - 重新计算战队总积分

积分计算规则（根据赛事规则图）：
- 排名分: #1=16, #2=12, #3=10, #4=8, #5=6, #6=5, #7=4, #8=3, #9=2, #10=1, #11-16=0
- 淘汰分: 每淘汰1人得1分
- 总积分 = 排名分 + 淘汰分

**文件**: `server/service/example/excel_import.go`

Excel导入逻辑：
- 解析Excel文件（使用excelize库）
- 字段映射：Excel列名 -> 数据库字段
- 数据验证：必填字段、格式校验
- 支持两种模式：
  - incremental（增量）：按teamCode匹配，已存在则跳过
  - full（全量）：按teamCode匹配，已存在则更新
- 返回导入结果（成功数、失败数、失败原因）

#### 3.1.4 新增API层

**文件**: `server/api/v1/example/competition_team.go`

主要API端点：
- `POST /competitionTeam` - 创建战队
- `PUT /competitionTeam` - 更新战队
- `DELETE /competitionTeam` - 删除战队
- `GET /competitionTeam` - 获取单个战队
- `GET /competitionTeamList` - 获取战队列表
- `POST /competitionTeam/importExcel` - Excel导入战队
- `POST /competitionTeam/addWarID` - 添加WarId
- `GET /competitionTeam/scores` - 获取战队积分记录
- `GET /competitionTeam/recentScores` - 获取最近积分记录
- `GET /competitionTeam/allScores` - 获取所有战队积分汇总

#### 3.1.5 新增路由配置

**文件**: `server/router/example/competition_team.go`

注册所有新API路由。

#### 3.1.6 数据库迁移

**文件**: `server/initialize/ensure_tables.go`

在 `MigrateTable` 和 `TableCreated` 方法中添加新模型：
```go
example.CompetitionTeam{},
example.TeamScore{},
```

### 3.2 前端开发

#### 3.2.1 新增API接口

**文件**: `web/src/api/competitionTeam.js`

封装所有后端API调用。

#### 3.2.2 新增页面组件

**文件**: `web/src/view/example/competitionTeam/index.vue`

主页面组件，包含：
- 战队列表展示区
- 新增/编辑战队弹窗
- Excel导入功能
- WarId管理面板
- 积分统计展示
- 响应式布局

#### 3.2.3 子组件

**文件**: `web/src/view/example/competitionTeam/TeamCard.vue`

战队卡片组件，展示：
- Logo图片
- 战队名称
- 战队标识
- 总积分
- 操作按钮

**文件**: `web/src/view/example/competitionTeam/WarIDManager.vue`

WarId管理组件：
- WarId输入框
- 验证/提交按钮
- 调用xdc/get_info接口
- 显示接口调用状态
- 积分计算展示

**文件**: `web/src/view/example/competitionTeam/ScoreChart.vue`

积分统计组件：
- 进度条展示总积分
- 最近4次记录列表
- 时间戳显示
- 排名变化

**文件**: `web/src/view/example/competitionTeam/ImportExcelDialog.vue`

Excel导入弹窗：
- 文件上传
- 模式选择（增量/全量）
- 导入进度显示
- 结果反馈

#### 3.2.4 样式设计

符合赛事主题风格：
- 主色调：深红/黑色/金色（电竞风格）
- 响应式布局：支持桌面端和移动端
- 视觉层次分明
- 交互反馈及时

### 3.3 Excel 文件格式

根据赛事文档结构，Excel应包含以下列：
- 战队名称
- 战队标识（缩写）
- 战队Logo（URL或文件名）
- 其他战队相关信息

## 四、文件清单

### 后端新增文件
1. `server/model/example/competition_team.go` - 战队模型
2. `server/model/example/team_score.go` - 积分记录模型
3. `server/model/example/request/competition_team.go` - 请求结构
4. `server/model/example/response/competition_team.go` - 响应结构
5. `server/service/example/competition_team.go` - 业务逻辑
6. `server/service/example/excel_import.go` - Excel导入
7. `server/api/v1/example/competition_team.go` - API接口
8. `server/router/example/competition_team.go` - 路由配置

### 后端修改文件
1. `server/initialize/ensure_tables.go` - 添加新表迁移
2. `server/router/example/enter.go` - 注册新路由

### 前端新增文件
1. `web/src/api/competitionTeam.js` - API接口
2. `web/src/view/example/competitionTeam/index.vue` - 主页面
3. `web/src/view/example/competitionTeam/TeamCard.vue` - 战队卡片
4. `web/src/view/example/competitionTeam/WarIDManager.vue` - WarId管理
5. `web/src/view/example/competitionTeam/ScoreChart.vue` - 积分展示
6. `web/src/view/example/competitionTeam/ImportExcelDialog.vue` - Excel导入

## 五、注意事项

1. **数据验证**: 战队标识需为3-8位字母数字组合，战队名称最长50字符
2. **Logo上传**: 支持JPG/PNG格式，最大2MB
3. **接口调用**: WarId验证需调用xdc/get_info接口，处理好错误和超时
4. **积分计算**: 严格按照赛事规则图中的积分规则实现
5. **响应式设计**: 确保在不同设备上良好显示
6. **缓存优化**: 积分统计数据可考虑使用前端缓存

## 六、实施顺序

1. 后端模型和数据库迁移
2. 后端Service层和API层
3. 后端路由配置
4. 前端API接口封装
5. 前端页面组件开发
6. 样式优化和响应式适配
7. 功能测试和联调
