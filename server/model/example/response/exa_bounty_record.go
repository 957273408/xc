package response

import "github.com/flipped-aurora/gin-vue-admin/server/model/example"

type BountyRecordListResponse struct {
	List  []example.ExaBountyRecord `json:"list"`
	Total int64                     `json:"total"`
}

// ExaTeamBountyRankingItem 队伍赏金排行榜项（基于 ExaTeam 表）
type ExaTeamBountyRankingItem struct {
	ID          uint    `json:"id"`
	TeamName    string  `json:"teamName"`
	TotalBounty float64 `json:"totalBounty"`
	Rank        int     `json:"rank"`
}

// ExaTeamBountyRankingResponse 队伍赏金排行榜响应（基于 ExaTeam 表）
type ExaTeamBountyRankingResponse struct {
	List  []ExaTeamBountyRankingItem `json:"list"`
	Total int64                      `json:"total"`
	Page  int                        `json:"page"`
}

// PlayerBountyRankingItem 选手赏金排行榜项
type PlayerBountyRankingItem struct {
	ID         uint    `json:"id"`
	PlayerName string  `json:"playerName"`
	Avatar     string  `json:"avatar"`
	TeamName   string  `json:"teamName"`
	Bounty     float64 `json:"bounty"`
	Rank       int     `json:"rank"`
}

// PlayerBountyRankingResponse 选手赏金排行榜响应
type PlayerBountyRankingResponse struct {
	List  []PlayerBountyRankingItem `json:"list"`
	Total int64                     `json:"total"`
	Page  int                       `json:"page"`
}
