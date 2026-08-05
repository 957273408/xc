package example

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/example"
	exaReq "github.com/flipped-aurora/gin-vue-admin/server/model/example/request"
	exaResp "github.com/flipped-aurora/gin-vue-admin/server/model/example/response"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type PlayerSelectionService struct{}

var PlayerSelectionServiceApp = new(PlayerSelectionService)

// 简易内存缓存（考虑到Redis不存在时的fallback
type cacheEntry struct {
	data      *exaResp.WarPlayersResponse
	expiresAt time.Time
}

var (
	playerCache   = make(map[string]cacheEntry)
	cacheMutex    sync.RWMutex
	cacheDuration = 5 * time.Minute
)

// GetWarPlayers 获取指定WarId下的玩家列表（带缓存
func (s *PlayerSelectionService) GetWarPlayers(warID string) (*exaResp.WarPlayersResponse, error) {
	// 1. 查缓存
	cacheMutex.RLock()
	if entry, ok := playerCache[warID]; ok && time.Now().Before(entry.expiresAt) {
		cacheMutex.RUnlock()
		return entry.data, nil
	}
	cacheMutex.RUnlock()

	// 2. 调用代理接口
	resp, err := ProxyServiceApp.Forward("GET", "/xdc/get_info", map[string]interface{}{
		"warId": warID,
	})
	if err != nil {
		global.GVA_LOG.Error("获取战场玩家信息失败", zap.String("warId", warID), zap.Error(err))
		return nil, fmt.Errorf("获取战场信息失败: %v", err)
	}

	warData, ok := resp.Body.(map[string]interface{})
	if !ok {
		// 尝试兼容包装格式
		if wrapped, ok2 := resp.Body.(map[string]interface{}); ok2 {
			if data, hasData := wrapped["data"].(map[string]interface{}); hasData {
				warData = data
			}
		}
	}
	if warData == nil {
		return nil, fmt.Errorf("代理响应格式异常")
	}

	playerList, ok := warData["playerInfoList"].([]interface{})
	if !ok {
		// playerInfoList不存在或者空数组都返回空列表不报错
		result := &exaResp.WarPlayersResponse{WarID: warID, Players: []exaResp.PlayerInfo{}}
		s.putCache(warID, result)
		return result, nil
	}

	players := make([]exaResp.PlayerInfo, 0, len(playerList))
	for _, p := range playerList {
		player, ok := p.(map[string]interface{})
		if !ok {
			continue
		}
		info := parsePlayerInfo(player)
		players = append(players, info)
	}

	result := &exaResp.WarPlayersResponse{
		WarID:   warID,
		Players: players,
	}

	s.putCache(warID, result)
	return result, nil
}

func (s *PlayerSelectionService) putCache(warID string, data *exaResp.WarPlayersResponse) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	playerCache[warID] = cacheEntry{
		data:      data,
		expiresAt: time.Now().Add(cacheDuration),
	}
}

// parsePlayerInfo 根据接口文档字段精确解析玩家数据
// 文档字段: uID, nickName, totalKill, damage, hitHeadBulletNum, bulletTotalNum,
// hitBulletNum, heal, moveDistance, maxKillDistance, skillCard_1~12,
// hook, big, fire, cannon, wormHole, healBomb, shelter, cloudBomb, bomb, bombMax
func parsePlayerInfo(p map[string]interface{}) exaResp.PlayerInfo {
	info := exaResp.PlayerInfo{}

	// PlayerID: uID (int64)
	info.PlayerID = getStringOrFallback(p, []string{"uID", "uid"}, "")
	// NickName
	info.NickName = getStringOrFallback(p, []string{"nickName"}, "未知玩家")
	// KillCount: totalKill (int16)
	info.KillCount = int(getFloatOrFallback(p, []string{"totalKill"}, 0))

	// 伤害量: damage (float64)
	info.DamageAmount = int(getFloatOrFallback(p, []string{"damage"}, 0))

	// 爆头率: hitHeadBulletNum / bulletTotalNum * 100
	bulletTotal := getFloatOrFallback(p, []string{"bulletTotalNum"}, 0)
	hitHeadBullet := getFloatOrFallback(p, []string{"hitHeadBulletNum"}, 0)
	if bulletTotal > 0 {
		info.HeadshotRate = roundFloat((hitHeadBullet/bulletTotal)*100, 2)
	} else {
		info.HeadshotRate = 0
	}

	// 命中率: hitBulletNum / bulletTotalNum * 100
	hitBullet := getFloatOrFallback(p, []string{"hitBulletNum"}, 0)
	if bulletTotal > 0 {
		info.AccuracyRate = roundFloat((hitBullet/bulletTotal)*100, 2)
	} else {
		info.AccuracyRate = 0
	}

	// 治疗量: heal (float64)
	healing := int(getFloatOrFallback(p, []string{"heal"}, 0))
	info.HealingAmount = &healing

	// 移动距离: moveDistance (int32, 单位米)
	moveDist := getFloatOrFallback(p, []string{"moveDistance"}, 0)
	mv := roundFloat(moveDist, 2)
	info.MovementDistance = &mv

	// 最远击杀距离: maxKillDistance (int16, 单位米)
	longestKill := getFloatOrFallback(p, []string{"maxKillDistance"}, 0)
	lk := roundFloat(longestKill, 2)
	info.LongestKillDist = &lk

	// 身份卡使用次数: skillCard_1 ~ skillCard_12 总和
	idCardTotal := 0
	for i := 1; i <= 12; i++ {
		key := fmt.Sprintf("skillCard_%d", i)
		idCardTotal += int(getFloatOrFallback(p, []string{key}, 0))
	}
	info.IdentityCardUsed = &idCardTotal

	// 投掷物/道具使用次数: hook + big + fire + cannon + wormHole + healBomb + shelter + cloudBomb + bomb + bombMax
	throwableFields := []string{"hook", "big", "fire", "cannon", "wormHole", "healBomb", "shelter", "cloudBomb", "bomb", "bombMax"}
	throwableTotal := 0
	for _, field := range throwableFields {
		throwableTotal += int(getFloatOrFallback(p, []string{field}, 0))
	}
	info.ThrowablesUsed = &throwableTotal

	return info
}

func getStringOrFallback(m map[string]interface{}, keys []string, fallback string) string {
	for _, k := range keys {
		if v, exists := m[k]; exists {
			if str, ok := v.(string); ok && str != "" {
				return str
			}
			// 数字转字符串
			switch val := v.(type) {
			case float64:
				return fmt.Sprintf("%.0f", val)
			case int:
				return fmt.Sprintf("%d", val)
			}
		}
	}
	return fallback
}

func getFloatOrFallback(m map[string]interface{}, keys []string, fallback float64) float64 {
	for _, k := range keys {
		if v, exists := m[k]; exists {
			switch val := v.(type) {
			case float64:
				return val
			case int:
				return float64(val)
			case int64:
				return float64(val)
			case string:
				var f float64
				if _, err := fmt.Sscanf(val, "%f", &f); err == nil {
					return f
				}
			}
		}
	}
	return fallback
}

func roundFloat(val float64, precision int) float64 {
	format := fmt.Sprintf("%%.%df", precision)
	var result float64
	fmt.Sscanf(fmt.Sprintf(format, val), "%f", &result)
	return result
}

// generateSessionKey 生成随机会话标识
func generateSessionKey(warID string) string {
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		// fallback 时间戳
		return fmt.Sprintf("%s_%d", warID, time.Now().UnixNano())
	}
	return fmt.Sprintf("%s_%s", warID, hex.EncodeToString(b))
}

// SavePlayerSelection 保存玩家选择（Upsert语义）
func (s *PlayerSelectionService) SavePlayerSelection(req *exaReq.SavePlayerSelectionRequest) (*exaResp.PlayerSelectionResponse, error) {
	// 校验长度（4或5人均可）
	if len(req.SelectedPlayerIDs) < 4 || len(req.SelectedPlayerIDs) > 5 {
		return nil, errors.New("必须选择正好4名或5名玩家")
	}

	sessionKey := req.SessionKey
	if sessionKey == "" {
		sessionKey = generateSessionKey(req.WarID)
	}

	idsJSON, err := json.Marshal(req.SelectedPlayerIDs)
	if err != nil {
		return nil, fmt.Errorf("序列化玩家ID失败: %v", err)
	}

	// 序列化 warIds
	warIDsJSON := ""
	if len(req.WarIDs) > 0 {
		if wj, err := json.Marshal(req.WarIDs); err == nil {
			warIDsJSON = string(wj)
		}
	}

	now := time.Now()
	var record example.PlayerSelection
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		// 查找是否存在
		err := tx.Where("session_key = ?", sessionKey).First(&record).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 创建新记录
			record = example.PlayerSelection{
				WarID:             req.WarID,
				WarIDs:            warIDsJSON,
				SelectedPlayerIDs: string(idsJSON),
				ExtraStat1:        req.ExtraStat1,
				ExtraStat2:        req.ExtraStat2,
				ExtraStat:         req.ExtraStat,
				SessionKey:        sessionKey,
			}
			record.CreatedAt = now
			record.UpdatedAt = now
			if err := tx.Create(&record).Error; err != nil {
				return err
			}
		} else {
			// 更新现有记录
			record.WarID = req.WarID
			record.WarIDs = warIDsJSON
			record.SelectedPlayerIDs = string(idsJSON)
			record.ExtraStat1 = req.ExtraStat1
			record.ExtraStat2 = req.ExtraStat2
			record.ExtraStat = req.ExtraStat
			record.UpdatedAt = now
			if err := tx.Save(&record).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		global.GVA_LOG.Error("保存玩家选择失败", zap.String("warId", req.WarID), zap.Error(err))
		return nil, fmt.Errorf("保存失败: %v", err)
	}

	// 解析返回
	var selectedIDs []string
	if err := json.Unmarshal([]byte(record.SelectedPlayerIDs), &selectedIDs); err != nil {
		selectedIDs = req.SelectedPlayerIDs
	}

	var warIDs []string
	if record.WarIDs != "" {
		_ = json.Unmarshal([]byte(record.WarIDs), &warIDs)
	}

	return &exaResp.PlayerSelectionResponse{
		ID:                record.ID,
		WarID:             record.WarID,
		WarIDs:            warIDs,
		SelectedPlayerIDs: selectedIDs,
		ExtraStat1:        record.ExtraStat1,
		ExtraStat2:        record.ExtraStat2,
		ExtraStat:         record.ExtraStat,
		SessionKey:        record.SessionKey,
		CreatedAt:         record.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// GetPlayerSelection 按SessionKey读取保存的选择
func (s *PlayerSelectionService) GetPlayerSelection(sessionKey string) (*exaResp.PlayerSelectionResponse, error) {
	var record example.PlayerSelection
	err := global.GVA_DB.Where("session_key = ?", sessionKey).First(&record).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("未找到对应的玩家选择数据")
		}
		return nil, err
	}

	var selectedIDs []string
	if err := json.Unmarshal([]byte(record.SelectedPlayerIDs), &selectedIDs); err != nil {
		selectedIDs = []string{}
	}

	var warIDs []string
	if record.WarIDs != "" {
		_ = json.Unmarshal([]byte(record.WarIDs), &warIDs)
	}

	return &exaResp.PlayerSelectionResponse{
		ID:                record.ID,
		WarID:             record.WarID,
		WarIDs:            warIDs,
		SelectedPlayerIDs: selectedIDs,
		ExtraStat1:        record.ExtraStat1,
		ExtraStat2:        record.ExtraStat2,
		ExtraStat:         record.ExtraStat,
		SessionKey:        record.SessionKey,
		CreatedAt:         record.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// GetLatestSelection 获取最新保存的选择数据（无需sessionKey）
func (s *PlayerSelectionService) GetLatestSelection() (*exaResp.LatestSelectionResponse, error) {
	var record example.PlayerSelection
	err := global.GVA_DB.Order("created_at DESC").First(&record).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("暂无保存的选择数据")
		}
		return nil, err
	}

	var selectedIDs []string
	if err := json.Unmarshal([]byte(record.SelectedPlayerIDs), &selectedIDs); err != nil {
		return nil, fmt.Errorf("解析玩家ID失败: %v", err)
	}

	// 解析 warIDs
	var warIDs []string
	if record.WarIDs != "" {
		_ = json.Unmarshal([]byte(record.WarIDs), &warIDs)
	}

	// 根据是否有 warIDs 决定使用单场还是多场接口
	var playerMap map[string]exaResp.PlayerInfo
	if len(warIDs) > 0 {
		multiResp, err := s.GetMultiWarPlayers(warIDs)
		if err != nil {
			return nil, fmt.Errorf("获取多场玩家信息失败: %v", err)
		}
		playerMap = make(map[string]exaResp.PlayerInfo)
		for _, p := range multiResp.Players {
			playerMap[p.PlayerID] = p
		}
	} else {
		warPlayers, err := s.GetWarPlayers(record.WarID)
		if err != nil {
			return nil, fmt.Errorf("获取战场玩家信息失败: %v", err)
		}
		playerMap = make(map[string]exaResp.PlayerInfo)
		for _, p := range warPlayers.Players {
			playerMap[p.PlayerID] = p
		}
	}

	players := make([]exaResp.PlayerDetailItem, 0, len(selectedIDs))
	teamName := ""
	for _, id := range selectedIDs {
		p, ok := playerMap[id]
		if !ok {
			players = append(players, exaResp.PlayerDetailItem{
				NickName: "未知选手",
			})
			continue
		}

		if teamName == "" {
			teamName = extractTeamPrefix(p.NickName)
		}

		// 兼容逻辑：extraStat1 对应 data4，extraStat2 对应 data5
		stat1Type := record.ExtraStat1
		stat2Type := record.ExtraStat2
		if stat1Type == "" && stat2Type == "" && record.ExtraStat != "" {
			stat1Type = "damage"
			stat2Type = record.ExtraStat
		}
		if stat1Type == "" {
			stat1Type = "damage"
		}
		extra1Name, extra1Value := getExtraStatDisplay(stat1Type, p)
		extra2Name, extra2Value := getExtraStatDisplay(stat2Type, p)

		players = append(players, exaResp.PlayerDetailItem{
			NickName:   p.NickName,
			Data1Name:  "淘汰数",
			Data1Value: fmt.Sprintf("%d", p.KillCount),
			Data2Name:  "爆头率",
			Data2Value: fmt.Sprintf("%.1f%%", p.HeadshotRate),
			Data3Name:  "命中率",
			Data3Value: fmt.Sprintf("%.1f%%", p.AccuracyRate),
			Data4Name:  extra1Name,
			Data4Value: extra1Value,
			Data5Name:  extra2Name,
			Data5Value: extra2Value,
		})
	}

	// 查找战队信息
	resolvedTeamName := teamName
	teamLogo := ""
	if teamName != "" {
		var team example.CompetitionTeam
		if err := global.GVA_DB.Where("team_code = ?", teamName).First(&team).Error; err == nil {
			teamLogo = team.TeamLogo
			if team.TeamName != "" {
				resolvedTeamName = team.TeamName
			}
		}
	}

	// 回填 teamName 和 teamLogo 到每个玩家
	for i := range players {
		players[i].TeamName = resolvedTeamName
		players[i].TeamLogo = teamLogo
	}

	// 兼容返回
	stat1 := record.ExtraStat1
	stat2 := record.ExtraStat2
	if stat1 == "" && stat2 == "" && record.ExtraStat != "" {
		stat1 = "damage"
		stat2 = record.ExtraStat
	}
	if stat1 == "" {
		stat1 = "damage"
	}

	return &exaResp.LatestSelectionResponse{
		WarID:             record.WarID,
		WarIDs:            warIDs,
		TeamName:          resolvedTeamName,
		TeamLogo:          teamLogo,
		Players:           players,
		SelectedPlayerIDs: selectedIDs,
		ExtraStat1:        stat1,
		ExtraStat2:        stat2,
		ExtraStat:         record.ExtraStat,
		CreatedAt:         record.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// extractTeamPrefix 从昵称中提取战队前缀
func extractTeamPrefix(nick string) string {
	if nick == "" {
		return ""
	}
	for i, r := range nick {
		if isSeparator(r) && i > 0 {
			return nick[:i]
		}
	}
	return ""
}

// getExtraStatDisplay 根据附加统计类型返回名称和数值
func getExtraStatDisplay(extraStat string, p exaResp.PlayerInfo) (name, value string) {
	switch extraStat {
	case "damage":
		name = "伤害量"
		value = fmt.Sprintf("%d", p.DamageAmount)
	case "healing":
		name = "治疗量"
		if p.HealingAmount != nil {
			value = fmt.Sprintf("%d", *p.HealingAmount)
		} else {
			value = "0"
		}
	case "movement":
		name = "移动距离"
		if p.MovementDistance != nil {
			value = fmt.Sprintf("%.1f", *p.MovementDistance)
		} else {
			value = "0"
		}
	case "throwables":
		name = "投掷物数"
		if p.ThrowablesUsed != nil {
			value = fmt.Sprintf("%d", *p.ThrowablesUsed)
		} else {
			value = "0"
		}
	case "identity_card":
		name = "身份卡数"
		if p.IdentityCardUsed != nil {
			value = fmt.Sprintf("%d", *p.IdentityCardUsed)
		} else {
			value = "0"
		}
	case "longest_kill":
		name = "最远击杀"
		if p.LongestKillDist != nil {
			value = fmt.Sprintf("%.1f", *p.LongestKillDist)
		} else {
			value = "0"
		}
	default:
		name = "附加项"
		value = "-"
	}
	return
}

// ===== 多场汇总 =====

// playerAgg 多场数据聚合中间结构
type playerAgg struct {
	playerID      string
	nickName      string
	matchCount    int
	totalKill     float64
	damage        float64
	bulletTotal   float64
	hitHeadBullet float64
	hitBullet     float64
	heal          float64
	moveDistance  float64
	maxKillDist   float64 // 多场取最大值
	skillCards    [12]float64
	throwables    map[string]float64
}

// GetMultiWarPlayers 获取多个 WarId 的玩家数据并按 playerId 聚合
// 聚合规则：数值类字段求和；maxKillDistance 取最大；爆头率/命中率用累计命中数/累计总弹数重算
func (s *PlayerSelectionService) GetMultiWarPlayers(warIDs []string) (*exaResp.MultiWarPlayersResponse, error) {
	if len(warIDs) == 0 {
		return nil, errors.New("warIds 不能为空")
	}

	throwableFields := []string{"hook", "big", "fire", "cannon", "wormHole", "healBomb", "shelter", "cloudBomb", "bomb", "bombMax"}
	aggMap := make(map[string]*playerAgg)
	matchCount := 0

	for _, warID := range warIDs {
		warID = strTrimSpace(warID)
		if warID == "" {
			continue
		}
		matchCount++

		// 直接调用代理获取原始数据（绕过缓存以拿到原始字段）
		resp, err := ProxyServiceApp.Forward("GET", "/xdc/get_info", map[string]interface{}{
			"warId": warID,
		})
		if err != nil {
			global.GVA_LOG.Error("多场汇总-获取战场信息失败", zap.String("warId", warID), zap.Error(err))
			continue
		}

		warData, ok := resp.Body.(map[string]interface{})
		if !ok {
			if wrapped, ok2 := resp.Body.(map[string]interface{}); ok2 {
				if data, hasData := wrapped["data"].(map[string]interface{}); hasData {
					warData = data
				}
			}
		}
		if warData == nil {
			continue
		}

		playerList, ok := warData["playerInfoList"].([]interface{})
		if !ok {
			continue
		}

		for _, p := range playerList {
			player, ok := p.(map[string]interface{})
			if !ok {
				continue
			}
			pid := getStringOrFallback(player, []string{"uID", "uid"}, "")
			if pid == "" {
				continue
			}

			agg, exists := aggMap[pid]
			if !exists {
				agg = &playerAgg{
					playerID:   pid,
					nickName:   getStringOrFallback(player, []string{"nickName"}, "未知玩家"),
					throwables: make(map[string]float64),
				}
				aggMap[pid] = agg
			}
			agg.matchCount++
			agg.totalKill += getFloatOrFallback(player, []string{"totalKill"}, 0)
			agg.damage += getFloatOrFallback(player, []string{"damage"}, 0)
			agg.bulletTotal += getFloatOrFallback(player, []string{"bulletTotalNum"}, 0)
			agg.hitHeadBullet += getFloatOrFallback(player, []string{"hitHeadBulletNum"}, 0)
			agg.hitBullet += getFloatOrFallback(player, []string{"hitBulletNum"}, 0)
			agg.heal += getFloatOrFallback(player, []string{"heal"}, 0)
			agg.moveDistance += getFloatOrFallback(player, []string{"moveDistance"}, 0)

			if mk := getFloatOrFallback(player, []string{"maxKillDistance"}, 0); mk > agg.maxKillDist {
				agg.maxKillDist = mk
			}

			for i := 1; i <= 12; i++ {
				agg.skillCards[i-1] += getFloatOrFallback(player, []string{fmt.Sprintf("skillCard_%d", i)}, 0)
			}
			for _, f := range throwableFields {
				agg.throwables[f] += getFloatOrFallback(player, []string{f}, 0)
			}
		}
	}

	// 转换为 PlayerInfo 切片
	players := make([]exaResp.PlayerInfo, 0, len(aggMap))
	for _, a := range aggMap {
		info := exaResp.PlayerInfo{
			PlayerID:     a.playerID,
			NickName:     a.nickName,
			KillCount:    int(a.totalKill),
			DamageAmount: int(a.damage),
		}
		if a.bulletTotal > 0 {
			info.HeadshotRate = roundFloat((a.hitHeadBullet/a.bulletTotal)*100, 2)
			info.AccuracyRate = roundFloat((a.hitBullet/a.bulletTotal)*100, 2)
		}

		healing := int(a.heal)
		info.HealingAmount = &healing
		mv := roundFloat(a.moveDistance, 2)
		info.MovementDistance = &mv
		lk := roundFloat(a.maxKillDist, 2)
		info.LongestKillDist = &lk

		idCardTotal := 0
		for _, v := range a.skillCards {
			idCardTotal += int(v)
		}
		info.IdentityCardUsed = &idCardTotal

		throwableTotal := 0
		for _, v := range a.throwables {
			throwableTotal += int(v)
		}
		info.ThrowablesUsed = &throwableTotal

		players = append(players, info)
	}

	// 按淘汰数降序
	sortPlayersByKillDesc(players)

	return &exaResp.MultiWarPlayersResponse{
		WarIDs:     warIDs,
		MatchCount: matchCount,
		Players:    players,
	}, nil
}

func strTrimSpace(s string) string {
	for len(s) > 0 && (s[0] == ' ' || s[0] == '\t' || s[0] == '\n' || s[0] == '\r') {
		s = s[1:]
	}
	for len(s) > 0 && (s[len(s)-1] == ' ' || s[len(s)-1] == '\t' || s[len(s)-1] == '\n' || s[len(s)-1] == '\r') {
		s = s[:len(s)-1]
	}
	return s
}

// GetMultiWarTop5 多场数据Top5统计：淘汰数、爆头率、命中率、伤害量各取前5
func (s *PlayerSelectionService) GetMultiWarTop5(warIDs []string) (*exaResp.MultiWarTop5Response, error) {
	// 复用现有多场聚合逻辑
	multiResp, err := s.GetMultiWarPlayers(warIDs)
	if err != nil {
		return nil, err
	}

	players := multiResp.Players
	topKills := top5(players, func(p exaResp.PlayerInfo) float64 { return float64(p.KillCount) })
	topHeadshot := top5(players, func(p exaResp.PlayerInfo) float64 { return p.HeadshotRate })
	topAccuracy := top5(players, func(p exaResp.PlayerInfo) float64 { return p.AccuracyRate })
	topDamage := top5(players, func(p exaResp.PlayerInfo) float64 { return float64(p.DamageAmount) })

	return &exaResp.MultiWarTop5Response{
		WarIDs:      warIDs,
		MatchCount:  multiResp.MatchCount,
		TopKills:    topKills,
		TopHeadshot: topHeadshot,
		TopAccuracy: topAccuracy,
		TopDamage:   topDamage,
	}, nil
}

func top5(players []exaResp.PlayerInfo, score func(exaResp.PlayerInfo) float64) []exaResp.PlayerInfo {
	copied := make([]exaResp.PlayerInfo, len(players))
	copy(copied, players)
	sort.Slice(copied, func(i, j int) bool {
		return score(copied[i]) > score(copied[j])
	})
	n := 5
	if len(copied) < n {
		n = len(copied)
	}
	return copied[:n]
}

func sortPlayersByKillDesc(players []exaResp.PlayerInfo) {
	for i := 1; i < len(players); i++ {
		for j := i; j > 0 && players[j].KillCount > players[j-1].KillCount; j-- {
			players[j], players[j-1] = players[j-1], players[j]
		}
	}
}
