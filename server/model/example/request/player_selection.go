package request

type SavePlayerSelectionRequest struct {
	WarID             string   `json:"warId" form:"warId" binding:"required" label:"战场ID"`
	WarIDs            []string `json:"warIds" form:"warIds" binding:"omitempty" label:"多战场ID列表"`
	SelectedPlayerIDs []string `json:"selectedPlayerIds" form:"selectedPlayerIds" binding:"required,min=4,max=5" label:"选中玩家ID"`
	ExtraStat1        string   `json:"extraStat1" form:"extraStat1" binding:"omitempty,oneof=damage healing movement throwables identity_card longest_kill" label:"附加统计项1"`
	ExtraStat2        string   `json:"extraStat2" form:"extraStat2" binding:"omitempty,oneof=damage healing movement throwables identity_card longest_kill" label:"附加统计项2"`
	ExtraStat         string   `json:"extraStat" form:"extraStat" binding:"omitempty,oneof=healing movement throwables identity_card longest_kill" label:"附加统计项(兼容旧版)"`
	SessionKey        string   `json:"sessionKey" form:"sessionKey" label:"会话标识"`
}

type GetPlayerSelectionRequest struct {
	SessionKey string `form:"sessionKey" binding:"required" label:"会话标识"`
}

type GetWarPlayersRequest struct {
	WarID string `form:"warId" binding:"required" label:"战场ID"`
}

// GetMultiWarPlayersRequest 多场汇总请求
type GetMultiWarPlayersRequest struct {
	WarIDs []string `json:"warIds" form:"warIds" binding:"required,min=1" label:"战场ID列表"`
}
