package example

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

type CompetitionTeam struct {
	global.GVA_MODEL
	TeamCode   string `json:"teamCode" form:"teamCode" gorm:"column:team_code;comment:战队标识(3-8位字母数字);size:8;not null;uniqueIndex"`
	TeamName   string `json:"teamName" form:"teamName" gorm:"column:team_name;comment:战队名称;size:50;not null"`
	TeamLogo   string `json:"teamLogo" form:"teamLogo" gorm:"column:team_logo;comment:战队Logo路径;size:500"`
	TotalScore int    `json:"totalScore" form:"totalScore" gorm:"column:total_score;comment:总积分;default:0"`
}

func (CompetitionTeam) TableName() string {
	return "competition_team"
}
