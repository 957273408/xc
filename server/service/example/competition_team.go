package example

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"
	"unicode"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/example"
	exaResp "github.com/flipped-aurora/gin-vue-admin/server/model/example/response"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CompetitionTeamService struct{}

var CompetitionTeamServiceApp = new(CompetitionTeamService)

var rankScoreTable = map[int]int{
	1: 16, 2: 12, 3: 10, 4: 8, 5: 6, 6: 5, 7: 4, 8: 3, 9: 2, 10: 1,
	11: 0, 12: 0, 13: 0, 14: 0, 15: 0, 16: 0,
}

func (s *CompetitionTeamService) CalculateScore(rank, killCount int) (rankScore, killScore, totalScore int) {
	if score, ok := rankScoreTable[rank]; ok {
		rankScore = score
	} else if rank >= 11 && rank <= 16 {
		rankScore = 0
	} else {
		rankScore = 0
	}
	killScore = killCount
	totalScore = rankScore + killScore
	return
}

func (s *CompetitionTeamService) CreateCompetitionTeam(team example.CompetitionTeam) error {
	var existing example.CompetitionTeam
	err := global.GVA_DB.Where("team_code = ?", team.TeamCode).First(&existing).Error
	if err == nil {
		return errors.New("战队标识已存在")
	}
	now := time.Now()
	team.CreatedAt = now
	team.UpdatedAt = now
	return global.GVA_DB.Create(&team).Error
}

func (s *CompetitionTeamService) UpdateCompetitionTeam(team *example.CompetitionTeam) error {
	var existing example.CompetitionTeam
	err := global.GVA_DB.Where("team_code = ? AND id != ?", team.TeamCode, team.ID).First(&existing).Error
	if err == nil {
		return errors.New("战队标识已被其他战队使用")
	}
	return global.GVA_DB.Model(&example.CompetitionTeam{}).Where("id = ?", team.ID).Updates(map[string]interface{}{
		"team_code":   team.TeamCode,
		"team_name":   team.TeamName,
		"team_logo":   team.TeamLogo,
		"group_name":  team.GroupName,
		"total_score": team.TotalScore,
	}).Error
}

func (s *CompetitionTeamService) DeleteCompetitionTeam(id uint) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("team_id = ?", id).Delete(&example.TeamScore{}).Error; err != nil {
			return err
		}
		return tx.Delete(&example.CompetitionTeam{}, id).Error
	})
}

func (s *CompetitionTeamService) GetCompetitionTeam(id uint) (example.CompetitionTeam, error) {
	var team example.CompetitionTeam
	err := global.GVA_DB.Where("id = ?", id).First(&team).Error
	return team, err
}

func (s *CompetitionTeamService) GetCompetitionTeamList(info request.PageInfo) ([]example.CompetitionTeam, int64, error) {
	var teams []example.CompetitionTeam
	var total int64
	db := global.GVA_DB.Model(&example.CompetitionTeam{})

	if info.Keyword != "" {
		db = db.Where("team_name LIKE ? OR team_code LIKE ?", "%"+info.Keyword+"%", "%"+info.Keyword+"%")
	}

	err := db.Count(&total).Error
	if err != nil {
		return teams, total, err
	}

	err = db.Order("total_score DESC, id DESC").Limit(info.PageSize).Offset(info.PageSize * (info.Page - 1)).Find(&teams).Error
	return teams, total, err
}

func (s *CompetitionTeamService) AddWarID(teamID uint, warID string) (*example.TeamScore, error) {
	var team example.CompetitionTeam
	if err := global.GVA_DB.Where("id = ?", teamID).First(&team).Error; err != nil {
		return nil, errors.New("战队不存在")
	}

	var existingScore example.TeamScore
	err := global.GVA_DB.Where("team_id = ? AND war_id = ?", teamID, warID).First(&existingScore).Error
	if err == nil {
		return nil, errors.New("该WarId已添加")
	}

	playerCount, killCount, rank, bountyCoin, err := s.FetchWarInfo(warID, team.TeamCode)
	if err != nil {
		global.GVA_LOG.Error("获取战场信息失败", zap.Error(err))
		return nil, fmt.Errorf("获取战场信息失败: %v", err)
	}

	if rank == 0 {
		rank = playerCount
	}

	rankScore, killScore, totalScore := s.CalculateScore(rank, killCount)

	score := &example.TeamScore{
		TeamID:     teamID,
		WarID:      warID,
		Rank:       rank,
		KillCount:  killCount,
		RankScore:  rankScore,
		KillScore:  killScore,
		TotalScore: totalScore,
		BountyCoin: bountyCoin,
		SettleTime: time.Now(),
	}

	if err := global.GVA_DB.Create(score).Error; err != nil {
		return nil, err
	}

	if err := s.RecalculateTeamTotalScore(teamID); err != nil {
		return nil, err
	}

	return score, nil
}

func (s *CompetitionTeamService) FetchWarInfo(warID string, teamCode string) (teamPlayerCount, killCount, rank int, bountyCoin int64, err error) {
	warData, err := s.fetchWarInfoRaw(warID)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	playerList, ok := warData["playerInfoList"].([]interface{})
	if !ok {
		return 0, 0, 0, 0, nil
	}

	// 统计该战队的选手数和淘汰数
	for _, p := range playerList {
		player, ok := p.(map[string]interface{})
		if !ok {
			continue
		}
		nickName, _ := player["nickName"].(string)

		if isTeamMatched(nickName, teamCode) {
			teamPlayerCount++
			if totalKill, ok := player["totalKill"].(float64); ok {
				killCount += int(totalKill)
			}
			if bc, ok := player["bountyCoin"].(float64); ok {
				bountyCoin += int64(bc)
			}

			// 直接从 lostRoleRank 获取排名
			if lr, exists := player["lostRoleRank"]; exists {
				if r, ok := lr.(float64); ok && r > 0 {
					if rank == 0 || int(r) < rank {
						rank = int(r)
					}
				}
			}
		}
	}

	return teamPlayerCount, killCount, rank, bountyCoin, nil
}

func (s *CompetitionTeamService) GetTeamScores(teamID uint, limit int) ([]exaResp.TeamScoreRecordResponse, error) {
	var scores []example.TeamScore
	db := global.GVA_DB.Model(&example.TeamScore{}).Where("team_id = ?", teamID).Order("settle_time DESC")

	if limit > 0 {
		db = db.Limit(limit)
	}

	if err := db.Find(&scores).Error; err != nil {
		return nil, err
	}

	var result []exaResp.TeamScoreRecordResponse
	for _, score := range scores {
		result = append(result, exaResp.TeamScoreRecordResponse{
			ID:         score.ID,
			WarID:      score.WarID,
			Rank:       score.Rank,
			KillCount:  score.KillCount,
			RankScore:  score.RankScore,
			KillScore:  score.KillScore,
			TotalScore: score.TotalScore,
			BountyCoin: score.BountyCoin,
			GameTime:   score.GameTime,
			SettleTime: score.SettleTime.Format("2006-01-02 15:04:05"),
		})
	}
	return result, nil
}

func (s *CompetitionTeamService) GetTeamRecentScores(teamID uint, limit int) ([]exaResp.TeamScoreRecordResponse, error) {
	if limit <= 0 {
		limit = 4
	}
	return s.GetTeamScores(teamID, limit)
}

func (s *CompetitionTeamService) RecalculateTeamTotalScore(teamID uint) error {
	var result struct {
		Total int `json:"total"`
	}
	err := global.GVA_DB.Model(&example.TeamScore{}).
		Select("COALESCE(SUM(total_score), 0) as total").
		Where("team_id = ?", teamID).
		Scan(&result).Error
	if err != nil {
		return err
	}

	return global.GVA_DB.Model(&example.CompetitionTeam{}).
		Where("id = ?", teamID).
		Update("total_score", result.Total).Error
}

func (s *CompetitionTeamService) GetTeamDetail(teamID uint) (*exaResp.TeamDetailResponse, error) {
	var team example.CompetitionTeam
	if err := global.GVA_DB.Where("id = ?", teamID).First(&team).Error; err != nil {
		return nil, err
	}

	scores, err := s.GetTeamScores(teamID, 0)
	if err != nil {
		return nil, err
	}

	return &exaResp.TeamDetailResponse{
		Team:       team,
		ScoreList:  scores,
		TotalScore: team.TotalScore,
	}, nil
}

func (s *CompetitionTeamService) GetAllTeamsScoreSummary() ([]exaResp.ScoreSummaryResponse, error) {
	type TeamWithScore struct {
		example.CompetitionTeam
		MatchCount int    `json:"matchCount"`
		LastRank   int    `json:"lastRank"`
	}

	var teams []TeamWithScore
	err := global.GVA_DB.Model(&example.CompetitionTeam{}).
		Select("competition_team.*, COUNT(team_score.id) as match_count, "+
			"(SELECT rank FROM team_score WHERE team_score.team_id = competition_team.id ORDER BY team_score.settle_time DESC LIMIT 1) as last_rank").
		Joins("LEFT JOIN team_score ON team_score.team_id = competition_team.id").
		Group("competition_team.id").
		Order("competition_team.total_score DESC").
		Scan(&teams).Error
	if err != nil {
		return nil, err
	}

	var result []exaResp.ScoreSummaryResponse
	for _, t := range teams {
		result = append(result, exaResp.ScoreSummaryResponse{
			TeamID:     t.ID,
			TeamName:   t.TeamName,
			TeamCode:   t.TeamCode,
			TotalScore: t.TotalScore,
			MatchCount: t.MatchCount,
			LastRank:   t.LastRank,
		})
	}
	return result, nil
}

// GetTeamScoreRanking 获取战队积分排名（包含最近4次积分变动记录）
func (s *CompetitionTeamService) GetTeamScoreRanking(groupName string) (*exaResp.TeamScoreRankingResponse, error) {
	// 1. 查询战队，按总积分降序
	var teams []example.CompetitionTeam
	query := global.GVA_DB.Order("total_score DESC")
	if groupName != "" {
		query = query.Where("group_name = ?", groupName)
	}
	if err := query.Find(&teams).Error; err != nil {
		return nil, err
	}

	// 2. 获取所有战队 ID
	teamIDs := make([]uint, len(teams))
	for i, t := range teams {
		teamIDs[i] = t.ID
	}

	// 3. 批量获取各战队最近4场积分（只需 total_score）
	type ScoreRecord struct {
		TeamID     uint `gorm:"column:team_id"`
		TotalScore int  `gorm:"column:total_score"`
	}
	var allRecords []ScoreRecord
	subQuery := global.GVA_DB.Model(&example.TeamScore{}).
	Select("team_id, total_score, ROW_NUMBER() OVER (PARTITION BY team_id ORDER BY created_at DESC) AS rn").
	Where("team_id IN (?)", teamIDs)
	err := global.GVA_DB.Table("(?) AS ranked", subQuery).
		Where("rn <= ?", 4).
		Order("team_id ASC, rn ASC").
		Find(&allRecords).Error
	if err != nil {
		return nil, err
	}

	// 4. 按 teamID 分组，只收集 totalScore
	teamHistoryMap := make(map[uint][]int)
	for _, r := range allRecords {
		teamHistoryMap[r.TeamID] = append(teamHistoryMap[r.TeamID], r.TotalScore)
	}

	// 5. 批量查询各战队赏金总数
	type BountySum struct {
		TeamID      uint  `gorm:"column:team_id"`
		TotalBounty int64 `gorm:"column:total_bounty"`
	}
	var bountySums []BountySum
	global.GVA_DB.Model(&example.TeamScore{}).
		Select("team_id, SUM(bounty_coin) AS total_bounty").
		Where("team_id IN (?)", teamIDs).
		Group("team_id").
		Find(&bountySums)
	teamBountyMap := make(map[uint]int64)
	for _, b := range bountySums {
		teamBountyMap[b.TeamID] = b.TotalBounty
	}

	// 6. 批量查询各战队总淘汰数
	type KillSum struct {
		TeamID    uint `gorm:"column:team_id"`
		TotalKills int  `gorm:"column:total_kills"`
	}
	var killSums []KillSum
	global.GVA_DB.Model(&example.TeamScore{}).
		Select("team_id, SUM(kill_count) AS total_kills").
		Where("team_id IN (?)", teamIDs).
		Group("team_id").
		Find(&killSums)
	teamKillsMap := make(map[uint]int)
	for _, k := range killSums {
		teamKillsMap[k.TeamID] = k.TotalKills
	}

	// 6.5 同分按总淘汰数排序
	sort.SliceStable(teams, func(i, j int) bool {
		if teams[i].TotalScore != teams[j].TotalScore {
			return teams[i].TotalScore > teams[j].TotalScore
		}
		return teamKillsMap[teams[i].ID] > teamKillsMap[teams[j].ID]
	})

	// 7. 组装排名结果
	items := make([]exaResp.TeamRankingItem, 0, len(teams))
	for i, t := range teams {
		history := teamHistoryMap[t.ID]
		if history == nil {
			history = []int{}
		}
		hasPlayed := len(history) > 0
		item := exaResp.TeamRankingItem{
			Rank:        i + 1,
			TeamID:      t.ID,
			TeamName:    t.TeamCode,
			TeamCode:    t.TeamCode,
			TeamLogo:    t.TeamLogo,
			GroupName:   t.GroupName,
			TotalScore:  exaResp.FlexInt{Valid: hasPlayed, Value: t.TotalScore},
			TotalKills:  exaResp.FlexInt{Valid: hasPlayed, Value: teamKillsMap[t.ID]},
			TotalBounty: exaResp.FlexInt64{Valid: hasPlayed, Value: teamBountyMap[t.ID]},
			MatchCount:  exaResp.FlexInt{Valid: hasPlayed, Value: len(history)},
		}
		// 反转 history 使其从最旧到最新排列，然后从 scoreHistory4 开始填充
		rev := make([]int, len(history))
		for j := 0; j < len(history); j++ {
			rev[j] = history[len(history)-1-j]
		}
		if len(rev) > 0 {
			item.ScoreHistory4 = exaResp.FlexInt{Valid: true, Value: rev[0]} // 第1场（最早）
		}
		if len(rev) > 1 {
			item.ScoreHistory3 = exaResp.FlexInt{Valid: true, Value: rev[1]} // 第2场
		}
		if len(rev) > 2 {
			item.ScoreHistory2 = exaResp.FlexInt{Valid: true, Value: rev[2]} // 第3场
		}
		if len(rev) > 3 {
			item.ScoreHistory1 = exaResp.FlexInt{Valid: true, Value: rev[3]} // 第4场（最晚）
		}
		items = append(items, item)
	}

	return &exaResp.TeamScoreRankingResponse{Items: items}, nil
}

// GetTeamBountyRanking 按队伍总赏金排名，同赏金按总积分降序
func (s *CompetitionTeamService) GetTeamBountyRanking(groupName string) (*exaResp.TeamBountyRankingResponse, error) {
	// 1. 查询所有战队
	var teams []example.CompetitionTeam
	query := global.GVA_DB
	if groupName != "" {
		query = query.Where("group_name = ?", groupName)
	}
	if err := query.Find(&teams).Error; err != nil {
		return nil, err
	}

	teamIDs := make([]uint, len(teams))
	for i, t := range teams {
		teamIDs[i] = t.ID
	}

	// 2. 批量查询各战队总赏金（team_score.bounty_coin）
	type BountySum struct {
		TeamID      uint  `gorm:"column:team_id"`
		TotalBounty int64 `gorm:"column:total_bounty"`
	}
	var bountySums []BountySum
	global.GVA_DB.Model(&example.TeamScore{}).
		Select("team_id, SUM(bounty_coin) AS total_bounty").
		Where("team_id IN (?)", teamIDs).
		Group("team_id").
		Find(&bountySums)
	teamBountyMap := make(map[uint]int64)
	for _, b := range bountySums {
		teamBountyMap[b.TeamID] = b.TotalBounty
	}

	// 3. 批量查询各战队总淘汰数
	type KillSum struct {
		TeamID    uint `gorm:"column:team_id"`
		TotalKills int  `gorm:"column:total_kills"`
	}
	var killSums []KillSum
	global.GVA_DB.Model(&example.TeamScore{}).
		Select("team_id, SUM(kill_count) AS total_kills").
		Where("team_id IN (?)", teamIDs).
		Group("team_id").
		Find(&killSums)
	teamKillsMap := make(map[uint]int)
	for _, k := range killSums {
		teamKillsMap[k.TeamID] = k.TotalKills
	}

	// 4. 批量获取各战队最近4场积分
	type ScoreRecord struct {
		TeamID     uint `gorm:"column:team_id"`
		TotalScore int  `gorm:"column:total_score"`
	}
	var allRecords []ScoreRecord
	subQuery := global.GVA_DB.Model(&example.TeamScore{}).
		Select("team_id, total_score, ROW_NUMBER() OVER (PARTITION BY team_id ORDER BY created_at DESC) AS rn").
		Where("team_id IN (?)", teamIDs)
	global.GVA_DB.Table("(?) AS ranked", subQuery).
		Where("rn <= ?", 4).
		Order("team_id ASC, rn ASC").
		Find(&allRecords)
	teamHistoryMap := make(map[uint][]int)
	for _, r := range allRecords {
		teamHistoryMap[r.TeamID] = append(teamHistoryMap[r.TeamID], r.TotalScore)
	}

	// 5. 按总赏金降序排序，同赏金按总积分降序
	sort.SliceStable(teams, func(i, j int) bool {
		bi, bj := teamBountyMap[teams[i].ID], teamBountyMap[teams[j].ID]
		if bi != bj {
			return bi > bj
		}
		return teams[i].TotalScore > teams[j].TotalScore
	})

	// 6. 组装结果
	items := make([]exaResp.TeamBountyRankingItem, 0, len(teams))
	for i, t := range teams {
		history := teamHistoryMap[t.ID]
		if history == nil {
			history = []int{}
		}
		hasPlayed := len(history) > 0
		item := exaResp.TeamBountyRankingItem{
			Rank:        i + 1,
			TeamID:      t.ID,
			TeamName:    t.TeamCode,
			TeamCode:    t.TeamCode,
			TeamLogo:    t.TeamLogo,
			GroupName:   t.GroupName,
			TotalBounty: exaResp.FlexInt64{Valid: hasPlayed, Value: teamBountyMap[t.ID]},
			TotalScore:  exaResp.FlexInt{Valid: hasPlayed, Value: t.TotalScore},
			TotalKills:  exaResp.FlexInt{Valid: hasPlayed, Value: teamKillsMap[t.ID]},
			MatchCount:  exaResp.FlexInt{Valid: hasPlayed, Value: len(history)},
		}
		// 反转 history 从最旧到最新，填充 scoreHistory4→1
		rev := make([]int, len(history))
		for j := 0; j < len(history); j++ {
			rev[j] = history[len(history)-1-j]
		}
		if len(rev) > 0 {
			item.ScoreHistory4 = exaResp.FlexInt{Valid: true, Value: rev[0]}
		}
		if len(rev) > 1 {
			item.ScoreHistory3 = exaResp.FlexInt{Valid: true, Value: rev[1]}
		}
		if len(rev) > 2 {
			item.ScoreHistory2 = exaResp.FlexInt{Valid: true, Value: rev[2]}
		}
		if len(rev) > 3 {
			item.ScoreHistory1 = exaResp.FlexInt{Valid: true, Value: rev[3]}
		}
		items = append(items, item)
	}

	return &exaResp.TeamBountyRankingResponse{Items: items}, nil
}

func (s *CompetitionTeamService) DeleteTeamScore(teamID uint, warID string) error {
	var score example.TeamScore
	// 使用 Unscoped 查找（包括已软删除的记录）
	if err := global.GVA_DB.Unscoped().Where("team_id = ? AND war_id = ?", teamID, warID).First(&score).Error; err != nil {
		return errors.New("积分记录不存在")
	}

	// 使用 Unscoped 硬删除
	if err := global.GVA_DB.Unscoped().Delete(&score).Error; err != nil {
		return err
	}

	return s.RecalculateTeamTotalScore(teamID)
}

func (s *CompetitionTeamService) UpdateTeamScore(id, teamID uint, rank, killCount, rankScore, killScore, totalScore int, bountyCoin int64, settleTime string) error {
	var score example.TeamScore
	if err := global.GVA_DB.Unscoped().Where("id = ?", id).First(&score).Error; err != nil {
		return errors.New("积分记录不存在")
	}

	score.Rank = rank
	score.KillCount = killCount
	score.RankScore = rankScore
	score.KillScore = killScore
	score.TotalScore = totalScore
	score.BountyCoin = bountyCoin
	// 修改时不更新 settle_time，保持为创建时间

	if err := global.GVA_DB.Save(&score).Error; err != nil {
		return err
	}

	return s.RecalculateTeamTotalScore(teamID)
}

// CalculateWarIDForAllTeams 批量计算指定WarId下所有战队的积分（仅计算，不保存）
func (s *CompetitionTeamService) CalculateWarIDForAllTeams(warID string) (*exaResp.BatchWarIDCalcResponse, error) {
	// 获取所有战队
	var teams []example.CompetitionTeam
	if err := global.GVA_DB.Find(&teams).Error; err != nil {
		return nil, fmt.Errorf("获取战队列表失败: %v", err)
	}

	if len(teams) == 0 {
		return nil, errors.New("系统中暂无战队数据")
	}

	// 调用接口获取战场数据
	warData, err := s.fetchWarInfoRaw(warID)
	if err != nil {
		return nil, fmt.Errorf("获取战场信息失败: %v", err)
	}

	// 收集所有战队标识
	teamCodes := make([]string, 0, len(teams))
	for _, team := range teams {
		teamCodes = append(teamCodes, team.TeamCode)
	}

	// 按昵称前缀分组，计算各战队排名
	teamRankMap := s.calculateTeamRanks(warData, teamCodes)

	// 一次性匹配所有选手到战队，避免重复计数
	teamStatsMap := s.matchAllTeamPlayers(warData, teamCodes)

	resp := &exaResp.BatchWarIDCalcResponse{
		WarID:      warID,
		TotalTeams: len(teams),
		Items:      make([]exaResp.BatchWarIDScoreItem, 0, len(teams)),
	}

	for _, team := range teams {
		item := exaResp.BatchWarIDScoreItem{
			TeamID:   team.ID,
			TeamCode: team.TeamCode,
			TeamName: team.TeamName,
		}

		// 从预匹配结果中获取
		stats := teamStatsMap[team.TeamCode]
		if stats.playerCount == 0 {
			item.Matched = false
			item.Message = "未匹配到选手"
			resp.Items = append(resp.Items, item)
			continue
		}

		item.Matched = true
		item.PlayerCount = stats.playerCount
		item.KillCount = stats.killCount
		item.BountyCoin = stats.bountyCoin

		// 直接从 rankMap 获取排名（key 是 teamCode）
		rank := teamRankMap[team.TeamCode]
		item.Rank = rank

		rankScore, killScore, totalScore := s.CalculateScore(rank, stats.killCount)
		item.RankScore = rankScore
		item.KillScore = killScore
		item.TotalScore = totalScore

		resp.Items = append(resp.Items, item)
		resp.MatchedNum++
	}

	return resp, nil
}

// GetPublicWarScores 公开接口：获取指定WarId下所有战队的当场积分（按积分降序）
// extraTeamCodes 指定的队伍code会追加到列表末尾（即使本场未参赛）
func (s *CompetitionTeamService) GetPublicWarScores(warID string, extraTeamCodes []string) (*exaResp.PublicWarScoreResponse, error) {
	// 复用批量计算逻辑
	calcResp, err := s.CalculateWarIDForAllTeams(warID)
	if err != nil {
		return nil, err
	}

	type teamInfo struct {
		TeamName   string
		TeamCode   string
		TeamLogo   string
		GroupName  string
		TotalScore int
		RankScore  int
		KillCount  int
		BountyCoin int64
	}
	all := make([]teamInfo, 0, len(calcResp.Items))
	maxKills := 0
	for _, item := range calcResp.Items {
		if !item.Matched {
			continue
		}
		all = append(all, teamInfo{
			TeamName:   item.TeamName,
			TeamCode:   item.TeamCode,
			TeamLogo:   "",
			GroupName:  "",
			TotalScore: item.TotalScore,
			RankScore:  item.RankScore,
			KillCount:  item.KillCount,
			BountyCoin: item.BountyCoin,
		})
		if item.KillCount > maxKills {
			maxKills = item.KillCount
		}
	}

	// 查询战队Logo和GroupName
	var teams []example.CompetitionTeam
	global.GVA_DB.Find(&teams)
	logoMap := make(map[string]string)
	groupNameMap := make(map[string]string)
	for _, t := range teams {
		logoMap[t.TeamName] = t.TeamLogo
		groupNameMap[t.TeamName] = t.GroupName
	}

	// 按积分降序排序
	sort.Slice(all, func(i, j int) bool {
		return all[i].TotalScore > all[j].TotalScore
	})

	// 构建响应
	rankOneImg := "https://asset.fangguo.com/front-end/upload/3071503/2026-08/img/43d79914-9722-4920.png"
	resp := &exaResp.PublicWarScoreResponse{
		WarID: warID,
		Items: make([]exaResp.PublicWarScoreItem, 0, len(all)),
	}
	for i, m := range all {
		rank := i + 1
		rankOne := ""
		if m.RankScore == 16 {
			rankOne = rankOneImg
		}
		resp.Items = append(resp.Items, exaResp.PublicWarScoreItem{
			Rank:        rank,
			TeamName:    m.TeamCode,
			TeamLogo:    logoMap[m.TeamName],
			GroupName:   groupNameMap[m.TeamName],
			TotalScore:  exaResp.FlexInt{Valid: true, Value: m.TotalScore},
			RankScore:   exaResp.FlexInt{Valid: true, Value: m.RankScore},
			KillCount:   exaResp.FlexInt{Valid: true, Value: m.KillCount},
			BountyCoin:  exaResp.FlexInt64{Valid: true, Value: m.BountyCoin},
			IsTopKiller: m.KillCount == maxKills && maxKills > 0,
			RankOne:     rankOne,
		})
	}

	// 追加指定队伍到末尾（未在本场出现的队伍）
	if len(extraTeamCodes) > 0 {
		// 已出现在列表中的队伍code
		exists := make(map[string]struct{}, len(resp.Items))
		for _, item := range resp.Items {
			exists[item.TeamName] = struct{}{}
		}
		teamByCode := make(map[string]example.CompetitionTeam, len(teams))
		for _, t := range teams {
			teamByCode[t.TeamCode] = t
		}
		rank := len(resp.Items) + 1
		for _, code := range extraTeamCodes {
			if _, ok := exists[code]; ok {
				continue
			}
			t, ok := teamByCode[code]
			if !ok {
				continue // 系统中不存在的队伍code跳过
			}
			exists[code] = struct{}{}
			resp.Items = append(resp.Items, exaResp.PublicWarScoreItem{
				Rank:        rank,
				TeamName:    t.TeamCode,
				TeamLogo:    t.TeamLogo,
				GroupName:   t.GroupName,
				TotalScore:  exaResp.FlexInt{Valid: false},
				RankScore:   exaResp.FlexInt{Valid: false},
				KillCount:   exaResp.FlexInt{Valid: false},
				BountyCoin:  exaResp.FlexInt64{Valid: false},
				IsTopKiller: false,
			})
			rank++
		}
	}

	return resp, nil
}

// GetPublicWarBounty 公开接口：获取指定WarId下各战队的选手赏金分配（按bountyCoin从大到小）
func (s *CompetitionTeamService) GetPublicWarBounty(warID string) (*exaResp.PublicWarBountyResponse, error) {
	warData, err := s.fetchWarInfoRaw(warID)
	if err != nil {
		return nil, fmt.Errorf("获取战场信息失败: %v", err)
	}

	var teams []example.CompetitionTeam
	global.GVA_DB.Find(&teams)
	if len(teams) == 0 {
		return &exaResp.PublicWarBountyResponse{WarID: warID, Items: []exaResp.TeamBountyItem{}}, nil
	}

	teamCodes := make([]string, 0, len(teams))
	teamLogoMap := make(map[string]string)
	teamGroupMap := make(map[string]string)
	for _, t := range teams {
		teamCodes = append(teamCodes, t.TeamCode)
		teamLogoMap[t.TeamCode] = t.TeamLogo
		teamGroupMap[t.TeamCode] = t.GroupName
	}

	type teamDetail struct {
		teamCode    string
		teamLogo    string
		groupName   string
		totalBounty int64
		players     []int64 // 按赏金降序
	}
	teamMap := make(map[string]*teamDetail)

	playerList, _ := warData["playerInfoList"].([]interface{})
	for _, p := range playerList {
		player, ok := p.(map[string]interface{})
		if !ok {
			continue
		}
		nickName, _ := player["nickName"].(string)
		code := matchPlayerToTeam(nickName, teamCodes)
		if code == "" {
			continue
		}
		if _, exists := teamMap[code]; !exists {
			teamMap[code] = &teamDetail{
				teamCode:  code,
				teamLogo:  teamLogoMap[code],
				groupName: teamGroupMap[code],
			}
		}
		td := teamMap[code]
		var bounty int64
		if b, ok := player["bountyCoin"]; ok {
			switch v := b.(type) {
			case float64:
				bounty = int64(v)
			case string:
				fmt.Sscanf(v, "%d", &bounty)
			}
		}
		td.totalBounty += bounty
		td.players = append(td.players, bounty)
	}

	// 构建结果
	type resultItem struct {
		code string
		teamDetail
	}
	results := make([]resultItem, 0, len(teamMap))
	for code, td := range teamMap {
		// 选手赏金降序排序
		sort.Slice(td.players, func(i, j int) bool {
			return td.players[i] > td.players[j]
		})
		results = append(results, resultItem{code: code, teamDetail: *td})
	}

	// 按组别排序，同组按总赏金降序
	sort.Slice(results, func(i, j int) bool {
		if results[i].groupName != results[j].groupName {
			return results[i].groupName < results[j].groupName
		}
		return results[i].totalBounty > results[j].totalBounty
	})

	// 计算各组内最大赏金
	groupMaxBounty := make(map[string]int64)
	for _, r := range results {
		if r.totalBounty > groupMaxBounty[r.groupName] {
			groupMaxBounty[r.groupName] = r.totalBounty
		}
	}

	resp := &exaResp.PublicWarBountyResponse{
		WarID: warID,
		Items: make([]exaResp.TeamBountyItem, 0, len(results)),
	}
	// 组内排名
	groupRanks := make(map[string]int)
	for _, r := range results {
		groupRanks[r.groupName]++
		rank := groupRanks[r.groupName]
		item := exaResp.TeamBountyItem{
			Rank:        rank,
			TeamName:    r.teamCode,
			TeamLogo:    r.teamLogo,
			GroupName:   r.groupName,
			TotalBounty: r.totalBounty,
			IsTopBounty: r.totalBounty == groupMaxBounty[r.groupName] && r.totalBounty > 0,
		}
		// 填充最多4个选手
		for p := 0; p < len(r.players) && p < 4; p++ {
			switch p {
			case 0:
				item.Player1 = r.players[p]
			case 1:
				item.Player2 = r.players[p]
			case 2:
				item.Player3 = r.players[p]
			case 3:
				item.Player4 = r.players[p]
			}
		}
		resp.Items = append(resp.Items, item)
	}

	return resp, nil
}

// GetPublicTeamList 公开接口：获取所有战队的简化列表（TeamCode/TeamName/TeamLogo）
func (s *CompetitionTeamService) GetPublicTeamList() (*exaResp.PublicTeamListResponse, error) {
	var teams []example.CompetitionTeam
	err := global.GVA_DB.Order("total_score DESC, id DESC").Find(&teams).Error
	if err != nil {
		return nil, err
	}
	items := make([]exaResp.PublicTeamItem, 0, len(teams))
	for _, t := range teams {
		items = append(items, exaResp.PublicTeamItem{
			TeamCode: t.TeamCode,
			TeamName: t.TeamName,
			TeamLogo: t.TeamLogo,
		})
	}
	return &exaResp.PublicTeamListResponse{Teams: items}, nil
}

// fetchWarInfoRaw 通过代理服务获取战场原始数据
func (s *CompetitionTeamService) fetchWarInfoRaw(warID string) (map[string]interface{}, error) {
	resp, err := ProxyServiceApp.Forward("GET", "/xdc/get_info", map[string]interface{}{
		"warId": warID,
	})
	if err != nil {
		return nil, err
	}

	// Body 已由代理服务解析为 interface{}
	result, ok := resp.Body.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("代理响应格式异常: %T", resp.Body)
	}

	// 兼容包装格式：如果响应包含 data 字段且为对象，则解包
	if data, ok := result["data"].(map[string]interface{}); ok {
		if _, hasPlayerList := data["playerInfoList"]; hasPlayerList {
			return data, nil
		}
	}

	return result, nil
}

// isSeparator 判断字符是否为战队前缀分隔符
// 对应 JS: [.\uFF0E\u00B7\u2022\u30FB_\-#\s]
func isSeparator(r rune) bool {
	switch r {
	case '.', '\uFF0E', '\u00B7', '\u2022', '\u30FB', '_', '-', '#':
		return true
	}
	return unicode.IsSpace(r)
}

// getTeamFromNick 从昵称中提取队伍前缀
// "XHZ.辣子鸡" → "XHZ", "XYT·繁仙" → "XYT", "MS•兮辞" → "MS"
func getTeamFromNick(nick string) string {
	if nick == "" {
		return ""
	}
	for i, r := range nick {
		if isSeparator(r) {
			if i > 0 {
				return nick[:i]
			}
			return ""
		}
	}
	return nick
}

// isTeamMatched 判断昵称是否匹配战队标识
// 策略1: 分隔符前缀精确匹配（如 "XHZ.辣子鸡" → "XHZ"）
// 策略2: 对长标识(>=3字符)回退到包含匹配（区分大小写，避免短标识误匹配）
func isTeamMatched(nickName, teamCode string) bool {
	if nickName == "" || teamCode == "" {
		return false
	}
	prefix := getTeamFromNick(nickName)
	// 前缀精确匹配（不区分大小写）
	if strings.EqualFold(prefix, teamCode) {
		return true
	}
	// 仅对长标识(>=3字符)使用包含匹配，避免短标识如 "Y" 误匹配
	if len(teamCode) >= 3 {
		return strings.Contains(nickName, teamCode)
	}
	return false
}

// matchPlayerToTeam 从所有战队标识中找到该选手最佳匹配的战队
// 返回匹配的 teamCode，未匹配返回空字符串
// 优先匹配最长的战队标识，避免 "XYT" 的选手误匹配到 "Y"
func matchPlayerToTeam(nickName string, teamCodes []string) string {
	if nickName == "" || len(teamCodes) == 0 {
		return ""
	}
	prefix := getTeamFromNick(nickName)

	// 第一轮：前缀精确匹配，选最长的
	bestMatch := ""
	for _, code := range teamCodes {
		if strings.EqualFold(prefix, code) {
			if len(code) > len(bestMatch) {
				bestMatch = code
			}
		}
	}
	if bestMatch != "" {
		return bestMatch
	}

	// 第二轮：对长标识(>=3字符)使用包含匹配，选最长的
	for _, code := range teamCodes {
		if len(code) < 3 {
			continue
		}
		if strings.Contains(nickName, code) {
			if len(code) > len(bestMatch) {
				bestMatch = code
			}
		}
	}
	return bestMatch
}

// calculateTeamRanks 按昵称前缀分组玩家，使用 lostRoleRank 字段计算各战队排名
func (s *CompetitionTeamService) calculateTeamRanks(warData map[string]interface{}, teamCodes []string) map[string]int {
	playerList, ok := warData["playerInfoList"].([]interface{})
	if !ok {
		return map[string]int{}
	}

	// 用 teamCode 分组，收集每个战队的排名
	rankMap := make(map[string]int)
	for _, p := range playerList {
		player, ok := p.(map[string]interface{})
		if !ok {
			continue
		}
		nickName, _ := player["nickName"].(string)

		// 每个选手只匹配一个战队（最长匹配优先）
		code := matchPlayerToTeam(nickName, teamCodes)
		if code == "" {
			continue
		}
		// 使用 lostRoleRank 作为队伍排名
		if lr, exists := player["lostRoleRank"]; exists {
			if r, ok := lr.(float64); ok && r > 0 {
				// 取该战队最小的排名值（排名数字越小名次越高）
				if _, has := rankMap[code]; !has || int(r) < rankMap[code] {
					rankMap[code] = int(r)
				}
			}
		}
	}

	return rankMap
}

// matchAllTeamPlayers 一次性将所有选手匹配到最佳战队，返回每个战队的统计
type teamStat struct {
	playerCount int
	killCount   int
	bountyCoin  int64
}

func (s *CompetitionTeamService) matchAllTeamPlayers(warData map[string]interface{}, teamCodes []string) map[string]teamStat {
	result := make(map[string]teamStat)
	playerList, ok := warData["playerInfoList"].([]interface{})
	if !ok {
		return result
	}

	for _, p := range playerList {
		player, ok := p.(map[string]interface{})
		if !ok {
			continue
		}
		nickName, _ := player["nickName"].(string)
		// 每个选手只匹配一个战队
		code := matchPlayerToTeam(nickName, teamCodes)
		if code == "" {
			continue
		}
		stat := result[code]
		stat.playerCount++
		if totalKill, ok := player["totalKill"].(float64); ok {
			stat.killCount += int(totalKill)
		}
		if bountyCoin, ok := player["bountyCoin"].(float64); ok {
			stat.bountyCoin += int64(bountyCoin)
		}
		result[code] = stat
	}

	return result
}

// matchTeamPlayers 用昵称匹配单个战队标识（保留兼容，仅用于单战队场景）
func (s *CompetitionTeamService) matchTeamPlayers(warData map[string]interface{}, teamCode string) (playerCount, killCount int) {
	playerList, ok := warData["playerInfoList"].([]interface{})
	if !ok {
		return 0, 0
	}

	for _, p := range playerList {
		player, ok := p.(map[string]interface{})
		if !ok {
			continue
		}
		nickName, _ := player["nickName"].(string)
		if isTeamMatched(nickName, teamCode) {
			playerCount++
			if totalKill, ok := player["totalKill"].(float64); ok {
				killCount += int(totalKill)
			}
		}
	}

	return playerCount, killCount
}

// ConfirmWarIDScores 确认并保存批量积分
func (s *CompetitionTeamService) ConfirmWarIDScores(warID string, teamIDs []uint) (*exaResp.ConfirmWarIDResultResponse, error) {
	result := &exaResp.ConfirmWarIDResultResponse{
		WarID:  warID,
		Errors: []string{},
	}

	if len(teamIDs) == 0 {
		return result, errors.New("未选择任何战队")
	}

	// 重新计算（确保数据一致性）
	calcResp, err := s.CalculateWarIDForAllTeams(warID)
	if err != nil {
		return nil, fmt.Errorf("重新计算失败: %v", err)
	}

	// 构建 teamID -> 计算结果 的映射
	calcMap := make(map[uint]*exaResp.BatchWarIDScoreItem)
	for i := range calcResp.Items {
		calcMap[calcResp.Items[i].TeamID] = &calcResp.Items[i]
	}

	// 在事务中批量保存
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		for _, teamID := range teamIDs {
			item, exists := calcMap[teamID]
			if !exists {
				result.FailCount++
				result.Errors = append(result.Errors, fmt.Sprintf("战队ID %d: 未找到战队数据", teamID))
				continue
			}

			// 检查是否已存在（包括软删除的）
			var existing example.TeamScore
			hasExisting := false
			if err := tx.Unscoped().Where("team_id = ? AND war_id = ?", teamID, warID).First(&existing).Error; err == nil {
				// 硬删除旧记录，以便重新创建
				tx.Unscoped().Delete(&existing)
				hasExisting = true
			}

			// 保存积分记录
			score := &example.TeamScore{
				TeamID:     teamID,
				WarID:      warID,
				Rank:       item.Rank,
				KillCount:  item.KillCount,
				RankScore:  item.RankScore,
				KillScore:  item.KillScore,
				TotalScore: item.TotalScore,
				BountyCoin: item.BountyCoin,
				SettleTime: time.Now(),
			}
			if err := tx.Create(score).Error; err != nil {
				result.FailCount++
				result.Errors = append(result.Errors, fmt.Sprintf("战队 %s: 保存失败 - %v", item.TeamName, err))
				continue
			}

			// 更新战队总积分
			if err := s.recalculateTeamTotalScoreTx(tx, teamID); err != nil {
				result.FailCount++
				result.Errors = append(result.Errors, fmt.Sprintf("战队 %s: 更新总积分失败 - %v", item.TeamName, err))
				continue
			}

			if hasExisting {
				result.Errors = append(result.Errors, fmt.Sprintf("战队 %s: 已更新该WarId积分", item.TeamName))
			}

			if !item.Matched {
				result.Errors = append(result.Errors, fmt.Sprintf("战队 %s: 未匹配选手，已保存0分", item.TeamName))
			}

			result.SuccessCount++
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

// recalculateTeamTotalScoreTx 在事务中重算战队总积分
func (s *CompetitionTeamService) recalculateTeamTotalScoreTx(tx *gorm.DB, teamID uint) error {
	var result struct {
		Total int `json:"total"`
	}
	err := tx.Model(&example.TeamScore{}).
		Select("COALESCE(SUM(total_score), 0) as total").
		Where("team_id = ?", teamID).
		Scan(&result).Error
	if err != nil {
		return err
	}

	return tx.Model(&example.CompetitionTeam{}).
		Where("id = ?", teamID).
		Update("total_score", result.Total).Error
}
