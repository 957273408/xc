package response

type PlayerInfo struct {
	PlayerID     string  `json:"playerId"`
	NickName     string  `json:"nickName"`
	KillCount    int     `json:"killCount"`
	HeadshotRate float64 `json:"headshotRate"`
	AccuracyRate float64 `json:"accuracyRate"`
	DamageAmount int     `json:"damageAmount"`

	HealingAmount    *int     `json:"healingAmount,omitempty"`
	MovementDistance *float64 `json:"movementDistance,omitempty"`
	ThrowablesUsed   *int     `json:"throwablesUsed,omitempty"`
	IdentityCardUsed *int     `json:"identityCardUsed,omitempty"`
	LongestKillDist  *float64 `json:"longestKillDist,omitempty"`
}

type WarPlayersResponse struct {
	WarID   string       `json:"warId"`
	Players []PlayerInfo `json:"players"`
}

// MultiWarPlayersResponse 多场汇总响应
type MultiWarPlayersResponse struct {
	WarIDs     []string     `json:"warIds"`
	MatchCount int          `json:"matchCount"`
	Players    []PlayerInfo `json:"players"`
}

type PlayerSelectionResponse struct {
	ID                uint     `json:"id"`
	WarID             string   `json:"warId"`
	WarIDs            []string `json:"warIds,omitempty"`
	SelectedPlayerIDs []string `json:"selectedPlayerIds"`
	ExtraStat1        string   `json:"extraStat1"`
	ExtraStat2        string   `json:"extraStat2"`
	ExtraStat         string   `json:"extraStat"`
	SessionKey        string   `json:"sessionKey"`
	CreatedAt         string   `json:"createdAt"`
}

// StatItem 数据项（名称+数值）
type StatItem struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// PlayerDetailItem 选手详情（扁平化数据项字段）
type PlayerDetailItem struct {
	NickName   string `json:"nickName"`
	TeamName   string `json:"teamName"`
	TeamLogo   string `json:"teamLogo"`
	Data1Name  string `json:"data1name"`
	Data1Value string `json:"data1value"`
	Data2Name  string `json:"data2name"`
	Data2Value string `json:"data2value"`
	Data3Name  string `json:"data3name"`
	Data3Value string `json:"data3value"`
	Data4Name  string `json:"data4name"`
	Data4Value string `json:"data4value"`
	Data5Name  string `json:"data5name"`
	Data5Value string `json:"data5value"`
}

// LatestSelectionResponse 获取最新保存数据的响应
type LatestSelectionResponse struct {
	WarID             string             `json:"warId"`
	WarIDs            []string           `json:"warIds,omitempty"`
	TeamName          string             `json:"teamName"`
	TeamLogo          string             `json:"teamLogo"`
	Players           []PlayerDetailItem `json:"players"`
	SelectedPlayerIDs []string           `json:"selectedPlayerIds,omitempty"`
	ExtraStat1        string             `json:"extraStat1"`
	ExtraStat2        string             `json:"extraStat2"`
	ExtraStat         string             `json:"extraStat"`
	CreatedAt         string             `json:"createdAt"`
}

// MultiWarTop5Response 多场Top5统计响应
type MultiWarTop5Response struct {
	WarIDs        []string     `json:"warIds"`
	MatchCount    int          `json:"matchCount"`
	TopKills      []PlayerInfo `json:"topKills"`
	TopHeadshot   []PlayerInfo `json:"topHeadshot"`
	TopAccuracy   []PlayerInfo `json:"topAccuracy"`
	TopDamage     []PlayerInfo `json:"topDamage"`
}
