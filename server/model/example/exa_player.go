package example

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

type ExaPlayer struct {
	global.GVA_MODEL
	PlayerName string  `json:"playerName" form:"playerName" gorm:"comment:选手姓名"`
	UID        string  `json:"uid" form:"uid" gorm:"comment:玩家ID"`
	TeamID     uint    `json:"teamId" form:"teamId" gorm:"comment:所属战队ID"`
	Team       ExaTeam `json:"team" form:"team" gorm:"foreignKey:TeamID"`
	Bounty     float64 `json:"bounty" form:"bounty" gorm:"comment:当前赏金"`
}
