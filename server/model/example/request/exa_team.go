package request

type CreateTeamRequest struct {
	TeamName    string  `json:"teamName" form:"teamName" binding:"required"`
	TotalBounty float64 `json:"totalBounty" form:"totalBounty"`
}

type UpdateTeamRequest struct {
	ID          uint    `json:"id" form:"id" binding:"required"`
	TeamName    string  `json:"teamName" form:"teamName"`
	TotalBounty float64 `json:"totalBounty" form:"totalBounty"`
}

type SetTeamBountyRequest struct {
	TeamID uint    `json:"teamId" form:"teamId" binding:"required"`
	Bounty float64 `json:"bounty" form:"bounty"`
}
