package example

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

type PlayerSelection struct {
	global.GVA_MODEL
	WarID             string `json:"warId" form:"warId" gorm:"column:war_id;comment:战场ID;size:64;not null;index"`
	WarIDs            string `json:"warIds" form:"warIds" gorm:"column:war_ids;type:text;comment:多战场ID(JSON数组,多场汇总时使用)"`
	SelectedPlayerIDs string `json:"selectedPlayerIds" form:"selectedPlayerIds" gorm:"column:selected_player_ids;type:text;comment:选中的玩家ID(JSON数组)"`
	ExtraStat1        string `json:"extraStat1" form:"extraStat1" gorm:"column:extra_stat_1;comment:附加统计项1(第4项);size:32;not null;default:''"`
	ExtraStat2        string `json:"extraStat2" form:"extraStat2" gorm:"column:extra_stat_2;comment:附加统计项2(第5项);size:32;not null;default:''"`
	ExtraStat         string `json:"extraStat" form:"extraStat" gorm:"column:extra_stat;comment:附加统计项(兼容旧版);size:32;not null;default:''"`
	SessionKey        string `json:"sessionKey" form:"sessionKey" gorm:"column:session_key;comment:会话标识;size:64;not null;uniqueIndex"`
}

func (PlayerSelection) TableName() string {
	return "player_selection"
}
