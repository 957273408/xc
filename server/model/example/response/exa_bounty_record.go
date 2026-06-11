package response

import "github.com/flipped-aurora/gin-vue-admin/server/model/example"

type BountyRecordListResponse struct {
	List  []example.ExaBountyRecord `json:"list"`
	Total int64                     `json:"total"`
}
