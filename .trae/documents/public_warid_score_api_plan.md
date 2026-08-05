# 公开WarId积分查询API实现计划

## 概述

新增一个无需权限验证的公开API接口，通过传入WarId参数，实时计算并返回该场比赛中所有战队的积分数据，按积分降序排列。

## 当前架构分析

- **公开路由组**：`PublicGroup` 在 `router.go:65` 定义，无JWT/Casbin中间件
- **现有注册方式**：`exampleRouter.InitCompetitionTeamRouter(PrivateGroup)` 仅注册到私有组
- **积分计算逻辑**：`CalculateWarIDForAllTeams` 已实现完整的实时计算（获取战队→调用代理接口→匹配选手→计算排名分+淘汰分）
- **响应结构**：`BatchWarIDCalcResponse` 已包含所有需要的字段

## 实现方案

### 1. 新增响应结构
**文件**：`server/model/example/response/competition_team.go`

新增 `PublicWarScoreItem` 和 `PublicWarScoreResponse`：
```go
type PublicWarScoreItem struct {
    Rank           int    `json:"rank"`           // 排名（按积分降序的序号）
    TeamName       string `json:"teamName"`       // 战队名称
    TeamLogo       string `json:"teamLogo"`       // 战队logo URL
    TotalScore     int    `json:"totalScore"`     // 当场比赛积分
    IsTopKiller    bool   `json:"isTopKiller"`    // 是否淘汰数最多
}

type PublicWarScoreResponse struct {
    WarID  string               `json:"warId"`
    Items  []PublicWarScoreItem `json:"items"`
}
```

### 2. 新增Service方法
**文件**：`server/service/example/competition_team.go`

新增 `GetPublicWarScores(warID string)` 方法：
- 复用 `CalculateWarIDForAllTeams` 的核心逻辑（获取战队→fetchWarInfoRaw→matchTeamPlayers→calculateTeamRanks→CalculateScore）
- 过滤掉未匹配的战队（`Matched == false`）
- 按积分降序排序
- 标记淘汰数最多的战队 `IsTopKiller = true`
- 映射为 `PublicWarScoreItem` 返回

### 3. 新增API方法
**文件**：`server/api/v1/example/competition_team.go`

新增 `GetPublicWarScores(c *gin.Context)`：
- 从 query 参数获取 `warId`
- 调用 service 层获取数据
- 返回标准响应格式

### 4. 新增公开路由注册
**文件**：`server/router/example/competition_team.go`

在 `InitCompetitionTeamRouter` 方法中新增 `PublicRouter *gin.RouterGroup` 参数：
```go
func (c *CompetitionTeamRouter) InitCompetitionTeamRouter(Router, PublicRouter *gin.RouterGroup) {
    // ... 现有私有路由 ...
    // 公开路由（无鉴权）
    publicGroup := PublicRouter.Group("competitionTeam")
    {
        publicGroup.GET("public/warScores", competitionTeamApi.GetPublicWarScores)
    }
}
```

### 5. 更新路由注册
**文件**：`server/initialize/router.go`

将 `InitCompetitionTeamRouter` 调用改为传入 `PrivateGroup` 和 `PublicGroup`：
```go
exampleRouter.InitCompetitionTeamRouter(PrivateGroup, PublicGroup)
```

## 接口定义

- **路径**：`GET /competitionTeam/public/warScores?warId=xxx`
- **鉴权**：无
- **参数**：`warId`（query，必填）
- **响应示例**：
```json
{
  "code": 0,
  "msg": "获取成功",
  "data": {
    "warId": "0fe848a9b617b6e3f0ccd5fe73a5594d_1",
    "items": [
      {
        "rank": 1,
        "teamName": "流氓兔",
        "teamLogo": "https://...",
        "totalScore": 18,
        "isTopKiller": true
      }
    ]
  }
}
```

## 验证步骤

1. 重启后端服务
2. 调用 `GET /competitionTeam/public/warScores?warId=0fe848a9b617b6e3f0ccd5fe73a5594d_1`（不带token）
3. 验证返回数据按积分降序排列
4. 验证 `isTopKiller` 标记在淘汰数最多的战队上为 `true`
5. 验证未匹配战队不包含在返回结果中
