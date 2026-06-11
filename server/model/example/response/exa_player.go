package response

import "github.com/flipped-aurora/gin-vue-admin/server/model/example"

type PlayerResponse struct {
	Player example.ExaPlayer `json:"player"`
}

type PlayerListResponse struct {
	List  []example.ExaPlayer `json:"list"`
	Total int64               `json:"total"`
}

type BountyChangeResponse struct {
	KillerBounty float64 `json:"killerBounty"`
	VictimBounty float64 `json:"victimBounty"`
	Amount       float64 `json:"amount"`
}
