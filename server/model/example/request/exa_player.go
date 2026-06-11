package request

type CreatePlayerRequest struct {
	PlayerName string  `json:"playerName" form:"playerName" binding:"required"`
	UID        string  `json:"uid" form:"uid"`
	TeamID     uint    `json:"teamId" form:"teamId" binding:"required"`
	Bounty     float64 `json:"bounty" form:"bounty"`
}

type UpdatePlayerRequest struct {
	ID         uint    `json:"id" form:"id" binding:"required"`
	PlayerName string  `json:"playerName" form:"playerName"`
	UID        string  `json:"uid" form:"uid"`
	TeamID     uint    `json:"teamId" form:"teamId"`
	Bounty     float64 `json:"bounty" form:"bounty"`
}

type AllocateBountyRequest struct {
	TeamID       uint     `json:"teamId" form:"teamId" binding:"required"`
	PlayerBounties []PlayerBounty `json:"playerBounties" form:"playerBounties" binding:"required"`
}

type PlayerBounty struct {
	PlayerID uint    `json:"playerId" form:"playerId" binding:"required"`
	Amount   float64 `json:"amount" form:"amount" binding:"required"`
}

type KillRequest struct {
	KillerID uint `json:"killerId" form:"killerId" binding:"required"`
	VictimID uint `json:"victimId" form:"victimId" binding:"required"`
}

type ReviveRequest struct {
	PlayerID uint `json:"playerId" form:"playerId" binding:"required"`
}
