package request

type GetBountyRecordListRequest struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"pageSize" form:"pageSize"`
	PlayerID int `json:"playerId" form:"playerId"`
	TeamID   int `json:"teamId" form:"teamId"`
}

// GetTeamBountyRankingRequest 队伍赏金排行榜请求
type GetTeamBountyRankingRequest struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"pageSize" form:"pageSize"`
}

// GetPlayerBountyRankingRequest 选手赏金排行榜请求
type GetPlayerBountyRankingRequest struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"pageSize" form:"pageSize"`
}