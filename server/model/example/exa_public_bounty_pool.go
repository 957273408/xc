package example

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

type ExaPublicBountyPool struct {
	global.GVA_MODEL
	PoolName    string  `json:"poolName" form:"poolName" gorm:"comment:池子名称;default:公共赏金池"`
	TotalAmount float64 `json:"totalAmount" form:"totalAmount" gorm:"comment:池子总额"`
	Remark      string  `json:"remark" form:"remark" gorm:"comment:备注"`
}

func (ExaPublicBountyPool) TableName() string {
	return "exa_public_bounty_pools"
}
