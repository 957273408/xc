package example

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

type ExaTeam struct {
	global.GVA_MODEL
	TeamName   string  `json:"teamName" form:"teamName" gorm:"comment:战队名称"`
	TotalBounty float64 `json:"totalBounty" form:"totalBounty" gorm:"comment:总赏金"`
}
