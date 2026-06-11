package request

type GetBountyRecordListRequest struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"pageSize" form:"pageSize"`
	PlayerID int `json:"playerId" form:"playerId"`
	TeamID   int `json:"teamId" form:"teamId"`
}