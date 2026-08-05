package example

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

type TeamScore struct {
	global.GVA_MODEL
	TeamID      uint      `json:"teamId" form:"teamId" gorm:"column:team_id;comment:战队ID;not null;index:idx_team_war,unique"`
	WarID       string    `json:"warId" form:"warId" gorm:"column:war_id;comment:战场ID;size:100;not null;index:idx_team_war,unique"`
	Rank        int       `json:"rank" form:"rank" gorm:"column:rank;comment:排名"`
	KillCount   int       `json:"killCount" form:"killCount" gorm:"column:kill_count;comment:淘汰人数"`
	RankScore   int       `json:"rankScore" form:"rankScore" gorm:"column:rank_score;comment:排名分"`
	KillScore   int       `json:"killScore" form:"killScore" gorm:"column:kill_score;comment:淘汰分"`
	TotalScore  int       `json:"totalScore" form:"totalScore" gorm:"column:total_score;comment:总积分"`
	BountyCoin  int64     `json:"bountyCoin" form:"bountyCoin" gorm:"column:bounty_coin;comment:赏金"`
	GameTime    float64   `json:"gameTime" form:"gameTime" gorm:"column:game_time;comment:游戏时长(秒)"`
	SettleTime  time.Time `json:"settleTime" form:"settleTime" gorm:"column:settle_time;comment:结算时间"`
}

func (TeamScore) TableName() string {
	return "team_score"
}
