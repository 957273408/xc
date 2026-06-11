package response

import "github.com/flipped-aurora/gin-vue-admin/server/model/example"

type TeamResponse struct {
	Team example.ExaTeam `json:"team"`
}

type TeamListResponse struct {
	List  []example.ExaTeam `json:"list"`
	Total int64             `json:"total"`
}
