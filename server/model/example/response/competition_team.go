package response

import (
	"encoding/json"

	"github.com/flipped-aurora/gin-vue-admin/server/model/example"
)

// FlexInt 自定义整数类型，Valid=false 时 JSON 输出 "-"，Valid=true 时输出数字（包括0）
type FlexInt struct {
	Valid bool
	Value int
}

func (v FlexInt) MarshalJSON() ([]byte, error) {
	if !v.Valid {
		return []byte(`"-"`), nil
	}
	return json.Marshal(v.Value)
}

func (v *FlexInt) UnmarshalJSON(data []byte) error {
	s := string(data)
	if s == `"-"` || s == `""` || s == "null" {
		v.Valid = false
		v.Value = 0
		return nil
	}
	var i int
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}
	v.Valid = true
	v.Value = i
	return nil
}

// FlexInt64 自定义 int64 类型，Valid=false 时 JSON 输出 "-"，Valid=true 时输出数字（包括0）
type FlexInt64 struct {
	Valid bool
	Value int64
}

func (v FlexInt64) MarshalJSON() ([]byte, error) {
	if !v.Valid {
		return []byte(`"-"`), nil
	}
	return json.Marshal(v.Value)
}

func (v *FlexInt64) UnmarshalJSON(data []byte) error {
	s := string(data)
	if s == `"-"` || s == `""` || s == "null" {
		v.Valid = false
		v.Value = 0
		return nil
	}
	var i int64
	if err := json.Unmarshal(data, &i); err != nil {
		var i32 int
		if err2 := json.Unmarshal(data, &i32); err2 != nil {
			return err
		}
		v.Valid = true
		v.Value = int64(i32)
		return nil
	}
	v.Valid = true
	v.Value = i
	return nil
}

type TeamScoreRecordResponse struct {
	ID         uint    `json:"id"`
	WarID      string  `json:"warId"`
	Rank       int     `json:"rank"`
	KillCount  int     `json:"killCount"`
	RankScore  int     `json:"rankScore"`
	KillScore  int     `json:"killScore"`
	TotalScore int     `json:"totalScore"`
	BountyCoin int64   `json:"bountyCoin"`
	GameTime   float64 `json:"gameTime"`
	SettleTime string  `json:"settleTime"`
}

type TeamDetailResponse struct {
	Team       example.CompetitionTeam   `json:"team"`
	ScoreList  []TeamScoreRecordResponse `json:"scoreList"`
	TotalScore int                       `json:"totalScore"`
}

type ImportResultResponse struct {
	SuccessCount int      `json:"successCount"`
	FailCount    int      `json:"failCount"`
	Errors       []string `json:"errors"`
}

type ScoreSummaryResponse struct {
	TeamID     uint   `json:"teamId"`
	TeamName   string `json:"teamName"`
	TeamCode   string `json:"teamCode"`
	TotalScore int    `json:"totalScore"`
	MatchCount int    `json:"matchCount"`
	LastRank   int    `json:"lastRank"`
}

// BatchWarIDScoreItem 单个战队的批量积分计算结果
type BatchWarIDScoreItem struct {
	TeamID      uint   `json:"teamId"`
	TeamCode    string `json:"teamCode"`
	TeamName    string `json:"teamName"`
	PlayerCount int    `json:"playerCount"`
	KillCount   int    `json:"killCount"`
	Rank        int    `json:"rank"`
	RankScore   int    `json:"rankScore"`
	KillScore   int    `json:"killScore"`
	TotalScore  int    `json:"totalScore"`
	BountyCoin  int64  `json:"bountyCoin"`
	Matched     bool   `json:"matched"`
	Message     string `json:"message"`
}

// BatchWarIDCalcResponse 批量计算WarId积分的响应
type BatchWarIDCalcResponse struct {
	WarID      string                `json:"warId"`
	TotalTeams int                   `json:"totalTeams"`
	MatchedNum int                   `json:"matchedNum"`
	Items      []BatchWarIDScoreItem `json:"items"`
}

// ConfirmWarIDScoresRequest 确认批量积分的请求
type ConfirmWarIDScoresRequest struct {
	WarID   string  `json:"warId" binding:"required"`
	TeamIDs []uint  `json:"teamIds" binding:"required"`
}

// ConfirmWarIDResultResponse 确认批量积分的响应
type ConfirmWarIDResultResponse struct {
	WarID         string   `json:"warId"`
	SuccessCount  int      `json:"successCount"`
	FailCount     int      `json:"failCount"`
	Errors        []string `json:"errors"`
}

// PublicWarScoreItem 公开接口的单个战队积分项
type PublicWarScoreItem struct {
	Rank        int       `json:"rank"`        // 排名（按积分降序序号）
	TeamName    string    `json:"teamName"`    // 战队名称
	TeamLogo    string    `json:"teamLogo"`    // 战队logo URL
	GroupName   string    `json:"groupName"`   // 组别
	TotalScore  FlexInt   `json:"totalScore"`  // 当场比赛积分
	RankScore   FlexInt   `json:"rankScore"`   // 排名所获得的积分
	KillCount   FlexInt   `json:"killCount"`   // 淘汰数
	BountyCoin  FlexInt64 `json:"bountyCoin"`  // 赏金
	IsTopKiller bool      `json:"isTopKiller"` // 是否淘汰数最多
	RankOne     string    `json:"rankOne"`     // 排名第一图片URL，非第一为空
}

// PublicWarScoreResponse 公开接口的WarId积分响应
type PublicWarScoreResponse struct {
	WarID string               `json:"warId"`
	Items []PublicWarScoreItem `json:"items"`
}

// PublicTeamItem 公开接口的战队信息
type PublicTeamItem struct {
	TeamCode string `json:"teamCode"`
	TeamName string `json:"teamName"`
	TeamLogo string `json:"teamLogo"`
}

// PublicTeamListResponse 公开接口的战队列表响应
type PublicTeamListResponse struct {
	Teams []PublicTeamItem `json:"teams"`
}

// TeamBountyItem 战队赏金分配项
type TeamBountyItem struct {
	Rank        int      `json:"rank"`        // 排名（组内按总赏金降序）
	TeamName    string   `json:"teamName"`    // 战队名称
	TeamLogo    string   `json:"teamLogo"`    // 战队logo
	GroupName   string   `json:"groupName"`   // 组别
	TotalBounty int64    `json:"totalBounty"` // 总赏金
	Player1     int64    `json:"player1"`     // 选手1赏金
	Player2     int64    `json:"player2"`     // 选手2赏金
	Player3     int64    `json:"player3"`     // 选手3赏金
	Player4     int64    `json:"player4"`     // 选手4赏金
	IsTopBounty bool     `json:"isTopBounty"` // 是否组内赏金最高
}

// PublicWarBountyResponse 公开接口的WarId赏金分配响应
type PublicWarBountyResponse struct {
	WarID string           `json:"warId"`
	Items []TeamBountyItem `json:"items"`
}

// TeamRankingItem 积分排名项
type TeamRankingItem struct {
	Rank           int       `json:"rank"`
	TeamID         uint      `json:"teamId"`
	TeamName       string    `json:"teamName"`
	TeamCode       string    `json:"teamCode"`
	TeamLogo       string    `json:"teamLogo"`
	GroupName      string    `json:"groupName"`
	TotalScore     FlexInt   `json:"totalScore"`
	TotalKills     FlexInt   `json:"totalKills"`  // 总淘汰数
	TotalBounty    FlexInt64 `json:"totalBounty"` // 当前赏金总数
	MatchCount     FlexInt   `json:"matchCount"`
	LastRank       FlexInt   `json:"lastRank"`
	ScoreHistory1  FlexInt   `json:"scoreHistory1"`
	ScoreHistory2  FlexInt   `json:"scoreHistory2"`
	ScoreHistory3  FlexInt   `json:"scoreHistory3"`
	ScoreHistory4  FlexInt   `json:"scoreHistory4"`
}

// TeamScoreRankingResponse 战队积分排名响应
type TeamScoreRankingResponse struct {
	Items []TeamRankingItem `json:"items"`
}

// TeamBountyRankingItem 赏金排名项
type TeamBountyRankingItem struct {
	Rank          int       `json:"rank"`
	TeamID        uint      `json:"teamId"`
	TeamName      string    `json:"teamName"`
	TeamCode      string    `json:"teamCode"`
	TeamLogo      string    `json:"teamLogo"`
	GroupName     string    `json:"groupName"`
	TotalBounty   FlexInt64 `json:"totalBounty"`
	TotalScore    FlexInt   `json:"totalScore"`
	TotalKills    FlexInt   `json:"totalKills"`
	MatchCount    FlexInt   `json:"matchCount"`
	ScoreHistory1 FlexInt   `json:"scoreHistory1"`
	ScoreHistory2 FlexInt   `json:"scoreHistory2"`
	ScoreHistory3 FlexInt   `json:"scoreHistory3"`
	ScoreHistory4 FlexInt   `json:"scoreHistory4"`
}

// TeamBountyRankingResponse 战队赏金排名响应
type TeamBountyRankingResponse struct {
	Items []TeamBountyRankingItem `json:"items"`
}
