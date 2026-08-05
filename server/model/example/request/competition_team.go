package request

type CreateCompetitionTeamRequest struct {
	TeamCode string `json:"teamCode" form:"teamCode" binding:"required" label:"战队标识"`
	TeamName string `json:"teamName" form:"teamName" binding:"required,max=50" label:"战队名称"`
	TeamLogo string `json:"teamLogo" form:"teamLogo" label:"战队Logo"`
}

type UpdateCompetitionTeamRequest struct {
	ID       uint   `json:"id" form:"id" binding:"required" label:"ID"`
	TeamCode string `json:"teamCode" form:"teamCode" binding:"required" label:"战队标识"`
	TeamName string `json:"teamName" form:"teamName" binding:"required,max=50" label:"战队名称"`
	TeamLogo string `json:"teamLogo" form:"teamLogo" label:"战队Logo"`
}

type DeleteCompetitionTeamRequest struct {
	ID uint `json:"id" form:"id" binding:"required" label:"ID"`
}

type AddWarIDRequest struct {
	TeamID uint   `json:"teamId" form:"teamId" binding:"required" label:"战队ID"`
	WarID  string `json:"warId" form:"warId" binding:"required" label:"战场ID"`
}

type ImportExcelRequest struct {
	TeamCode string `json:"teamCode" form:"teamCode" binding:"required" label:"战队标识"`
	Mode     string `json:"mode" form:"mode" binding:"required,oneof=incremental full" label:"导入模式"`
}

type GetTeamScoresRequest struct {
	TeamID uint `json:"teamId" form:"teamId" binding:"required" label:"战队ID"`
	Limit  int  `json:"limit" form:"limit" label:"限制数量"`
}

type WarIDRequest struct {
	WarID string `json:"warId" form:"warId" binding:"required" label:"战场ID"`
}

type UpdateScoreRequest struct {
	ID         uint    `json:"id" form:"id" binding:"required" label:"积分记录ID"`
	TeamID     uint    `json:"teamId" form:"teamId" binding:"required" label:"战队ID"`
	Rank       int     `json:"rank" form:"rank" binding:"min=0,max=16" label:"排名"`
	KillCount  int     `json:"killCount" form:"killCount" binding:"min=0" label:"淘汰人数"`
	RankScore  int     `json:"rankScore" form:"rankScore" binding:"min=0" label:"排名分"`
	KillScore  int     `json:"killScore" form:"killScore" binding:"min=0" label:"淘汰分"`
	TotalScore int     `json:"totalScore" form:"totalScore" binding:"min=0" label:"总积分"`
	BountyCoin int64   `json:"bountyCoin" form:"bountyCoin" binding:"min=0" label:"赏金"`
	SettleTime string  `json:"settleTime" form:"settleTime" label:"结算时间"`
}
