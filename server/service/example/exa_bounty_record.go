package example

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/example"
	exaRes "github.com/flipped-aurora/gin-vue-admin/server/model/example/response"
)

type BountyRecordService struct{}

var BountyRecordServiceApp = new(BountyRecordService)

// GetRecordList 获取赏金记录列表
func (b *BountyRecordService) GetRecordList(info request.PageInfo, playerID, teamID uint) ([]example.ExaBountyRecord, int64, error) {
	var records []example.ExaBountyRecord
	var total int64
	db := global.GVA_DB.Model(&example.ExaBountyRecord{})
	if playerID > 0 {
		db = db.Where("player_id = ?", playerID)
	}
	if teamID > 0 {
		db = db.Where("team_id = ?", teamID)
	}
	err := db.Count(&total).Error
	if err != nil {
		return records, total, err
	}
	err = db.Order("created_at DESC").Limit(info.PageSize).Offset(info.PageSize * (info.Page - 1)).Find(&records).Error
	return records, total, err
}

// GetPoolInfo 获取公共赏金池信息
func (b *BountyRecordService) GetPoolInfo() (example.ExaPublicBountyPool, error) {
	var pool example.ExaPublicBountyPool
	err := global.GVA_DB.FirstOrCreate(&pool, example.ExaPublicBountyPool{PoolName: "公共赏金池"}).Error
	return pool, err
}

// GetTeamBountyRanking 获取队伍赏金排行榜
func (b *BountyRecordService) GetTeamBountyRanking(page, pageSize int) (exaRes.TeamBountyRankingResponse, error) {
	var teams []example.ExaTeam
	var total int64
	var result exaRes.TeamBountyRankingResponse

	db := global.GVA_DB.Model(&example.ExaTeam{})
	err := db.Count(&total).Error
	if err != nil {
		return result, err
	}

	// 按总赏金降序排序
	err = db.Order("total_bounty DESC").
		Limit(pageSize).
		Offset(pageSize * (page - 1)).
		Find(&teams).Error
	if err != nil {
		return result, err
	}

	// 构建排行榜项
	rankItems := make([]exaRes.TeamBountyRankingItem, len(teams))
	for i, team := range teams {
		rankItems[i] = exaRes.TeamBountyRankingItem{
			ID:          team.ID,
			TeamName:    team.TeamName,
			TotalBounty: team.TotalBounty,
			Rank:        pageSize*(page-1) + i + 1, // 计算实际排名
		}
	}

	result.List = rankItems
	result.Total = total
	result.Page = page
	return result, nil
}

// GetPlayerBountyRanking 获取选手赏金排行榜
func (b *BountyRecordService) GetPlayerBountyRanking(page, pageSize int) (exaRes.PlayerBountyRankingResponse, error) {
	var players []example.ExaPlayer
	var total int64
	var result exaRes.PlayerBountyRankingResponse

	db := global.GVA_DB.Model(&example.ExaPlayer{}).Preload("Team")
	err := db.Count(&total).Error
	if err != nil {
		return result, err
	}

	// 按赏金降序排序
	err = db.Order("bounty DESC").
		Limit(pageSize).
		Offset(pageSize * (page - 1)).
		Find(&players).Error
	if err != nil {
		return result, err
	}

	// 构建排行榜项
	rankItems := make([]exaRes.PlayerBountyRankingItem, len(players))
	for i, player := range players {
		teamName := ""
		if player.Team.ID > 0 {
			teamName = player.Team.TeamName
		}
		rankItems[i] = exaRes.PlayerBountyRankingItem{
			ID:         player.ID,
			PlayerName: player.PlayerName,
			Avatar:     "", // 暂时为空，后续可根据实际需求添加头像字段
			TeamName:   teamName,
			Bounty:     player.Bounty,
			Rank:       pageSize*(page-1) + i + 1, // 计算实际排名
		}
	}

	result.List = rankItems
	result.Total = total
	result.Page = page
	return result, nil
}

// AddToPool 向公共赏金池添加金额（内部使用）
func (b *BountyRecordService) AddToPool(amount float64, reason string) error {
	if amount <= 0 {
		return nil
	}
	tx := global.GVA_DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var pool example.ExaPublicBountyPool
	if err := tx.FirstOrCreate(&pool, example.ExaPublicBountyPool{PoolName: "公共赏金池"}).Error; err != nil {
		tx.Rollback()
		return err
	}

	newPoolAmount := pool.TotalAmount + amount
	if err := tx.Model(&pool).Update("total_amount", newPoolAmount).Error; err != nil {
		tx.Rollback()
		return err
	}

	record := example.ExaBountyRecord{
		PlayerID:   0,
		PlayerName: "系统",
		TeamID:     0,
		TeamName:   "公共赏金池",
		ChangeType: "pool_add",
		Amount:     amount,
		Balance:    newPoolAmount,
		Reason:     reason,
	}
	if err := tx.Create(&record).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// CreateRecord 创建赏金记录并同步更新战队总赏金
func (b *BountyRecordService) CreateRecord(record example.ExaBountyRecord) error {
	if err := global.GVA_DB.Create(&record).Error; err != nil {
		return err
	}
	return b.onRecordChanged(record.TeamID)
}

// UpdateRecord 更新赏金记录并同步更新战队总赏金
func (b *BountyRecordService) UpdateRecord(record example.ExaBountyRecord) error {
	var oldRecord example.ExaBountyRecord
	if err := global.GVA_DB.Where("id = ?", record.ID).First(&oldRecord).Error; err != nil {
		return err
	}
	if err := global.GVA_DB.Save(&record).Error; err != nil {
		return err
	}
	// 同步更新新旧战队的总赏金
	if err := b.onRecordChanged(oldRecord.TeamID); err != nil {
		return err
	}
	return b.onRecordChanged(record.TeamID)
}

// DeleteRecord 删除赏金记录并同步更新战队总赏金
func (b *BountyRecordService) DeleteRecord(recordID uint) error {
	var record example.ExaBountyRecord
	if err := global.GVA_DB.Where("id = ?", recordID).First(&record).Error; err != nil {
		return err
	}
	if err := global.GVA_DB.Delete(&record).Error; err != nil {
		return err
	}
	return b.onRecordChanged(record.TeamID)
}

// onRecordChanged 赏金记录变更后的处理钩子
// 自动重新计算并更新战队总赏金
func (b *BountyRecordService) onRecordChanged(teamID uint) error {
	if teamID == 0 {
		return nil
	}
	return TeamServiceApp.RecalculateTeamBounty(teamID)
}
