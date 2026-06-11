package example

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

type ExaBountyRecord struct {
	global.GVA_MODEL
	PlayerID    uint    `json:"playerId" form:"playerId" gorm:"comment:涉及选手ID"`
	PlayerName  string  `json:"playerName" form:"playerName" gorm:"comment:涉及选手姓名"`
	TeamID      uint    `json:"teamId" form:"teamId" gorm:"comment:所属战队ID"`
	TeamName    string  `json:"teamName" form:"teamName" gorm:"comment:所属战队名称"`
	ChangeType  string  `json:"changeType" form:"changeType" gorm:"comment:变动类型"`
	Amount      float64 `json:"amount" form:"amount" gorm:"comment:变动金额"`
	Balance     float64 `json:"balance" form:"balance" gorm:"comment:变动后余额"`
	Reason      string  `json:"reason" form:"reason" gorm:"comment:变动原因"`
	RelatedID   uint    `json:"relatedId" form:"relatedId" gorm:"comment:关联ID"`
	RelatedName string  `json:"relatedName" form:"relatedName" gorm:"comment:关联名称"`
}
