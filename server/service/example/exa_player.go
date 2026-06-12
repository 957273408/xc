package example

import (
	"errors"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/example"
	exaReq "github.com/flipped-aurora/gin-vue-admin/server/model/example/request"
	"go.uber.org/zap"
)

type PlayerService struct{}

var PlayerServiceApp = new(PlayerService)

func (p *PlayerService) CreatePlayer(player example.ExaPlayer) error {
	if err := global.GVA_DB.Create(&player).Error; err != nil {
		return err
	}
	if player.TeamID > 0 {
		go TeamServiceApp.RecalculateTeamBounty(player.TeamID)
	}
	return nil
}

func (p *PlayerService) DeletePlayer(id uint) error {
	var player example.ExaPlayer
	if err := global.GVA_DB.Where("id = ?", id).First(&player).Error; err != nil {
		return err
	}
	if err := global.GVA_DB.Delete(&example.ExaPlayer{}, id).Error; err != nil {
		return err
	}
	if player.TeamID > 0 {
		go TeamServiceApp.RecalculateTeamBounty(player.TeamID)
	}
	return nil
}

func (p *PlayerService) UpdatePlayer(player *example.ExaPlayer) error {
	var oldPlayer example.ExaPlayer
	if err := global.GVA_DB.Where("id = ?", player.ID).First(&oldPlayer).Error; err != nil {
		return err
	}
	if err := global.GVA_DB.Save(player).Error; err != nil {
		return err
	}
	if oldPlayer.TeamID > 0 {
		go TeamServiceApp.RecalculateTeamBounty(oldPlayer.TeamID)
	}
	if player.TeamID > 0 && player.TeamID != oldPlayer.TeamID {
		go TeamServiceApp.RecalculateTeamBounty(player.TeamID)
	}
	return nil
}

func (p *PlayerService) GetPlayer(id uint) (example.ExaPlayer, error) {
	var player example.ExaPlayer
	err := global.GVA_DB.Where("id = ?", id).Preload("Team").First(&player).Error
	return player, err
}

func (p *PlayerService) GetPlayerList(info request.PageInfo, teamID uint) ([]example.ExaPlayer, int64, error) {
	var players []example.ExaPlayer
	var total int64
	db := global.GVA_DB.Model(&example.ExaPlayer{}).Preload("Team")
	if teamID > 0 {
		db = db.Where("team_id = ?", teamID)
	}
	err := db.Count(&total).Error
	if err != nil {
		return players, total, err
	}
	err = db.Limit(info.PageSize).Offset(info.PageSize * (info.Page - 1)).Find(&players).Error
	return players, total, err
}

func (p *PlayerService) GetPlayersByTeam(teamID uint) ([]example.ExaPlayer, error) {
	var players []example.ExaPlayer
	err := global.GVA_DB.Where("team_id = ?", teamID).Find(&players).Error
	return players, err
}

func (p *PlayerService) UpdatePlayerBounty(playerID uint, bounty float64) error {
	return global.GVA_DB.Model(&example.ExaPlayer{}).Where("id = ?", playerID).Update("bounty", bounty).Error
}

func (p *PlayerService) AllocateBounty(req exaReq.AllocateBountyRequest) error {
	tx := global.GVA_DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var team example.ExaTeam
	if err := tx.Where("id = ?", req.TeamID).First(&team).Error; err != nil {
		tx.Rollback()
		global.GVA_LOG.Error("赏金分配失败: 战队不存在", zap.Uint("teamID", req.TeamID), zap.Error(err))
		return err
	}

	var totalAllocated float64
	for _, pb := range req.PlayerBounties {
		totalAllocated += pb.Amount
	}

	remainingAmount := team.TotalBounty - totalAllocated

	global.GVA_LOG.Info("赏金分配校验开始",
		zap.Uint("teamID", req.TeamID),
		zap.String("teamName", team.TeamName),
		zap.Float64("totalBounty", team.TotalBounty),
		zap.Float64("totalAllocated", totalAllocated),
		zap.Float64("remainingAmount", remainingAmount))

	if remainingAmount != 0 {
		tx.Rollback()
		global.GVA_LOG.Error("赏金分配失败: 剩余金额不为0",
			zap.Uint("teamID", req.TeamID),
			zap.String("teamName", team.TeamName),
			zap.Float64("remainingAmount", remainingAmount))
		return errors.New("分配失败：剩余金额必须为0才能执行分配操作")
	}

	global.GVA_LOG.Info("赏金分配校验通过",
		zap.Uint("teamID", req.TeamID),
		zap.String("teamName", team.TeamName))

	for _, pb := range req.PlayerBounties {
		if err := tx.Model(&example.ExaPlayer{}).Where("id = ?", pb.PlayerID).Update("bounty", pb.Amount).Error; err != nil {
			tx.Rollback()
			return err
		}

		var player example.ExaPlayer
		if err := tx.Where("id = ?", pb.PlayerID).First(&player).Error; err != nil {
			tx.Rollback()
			return err
		}

		record := example.ExaBountyRecord{
			PlayerID:    pb.PlayerID,
			PlayerName:  player.PlayerName,
			TeamID:      req.TeamID,
			TeamName:    team.TeamName,
			ChangeType:  "allocate",
			Amount:      pb.Amount,
			Balance:     pb.Amount,
			Reason:      "战队赏金分配",
			RelatedID:   req.TeamID,
			RelatedName: team.TeamName,
		}
		if err := tx.Create(&record).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	if req.TeamID > 0 {
		go TeamServiceApp.RecalculateTeamBounty(req.TeamID)
	}
	return nil
}

func (p *PlayerService) Kill(killerID, victimID uint) (float64, error) {
	tx := global.GVA_DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 获取伤害者和受害者信息（预加载战队数据）
	var killer, victim example.ExaPlayer
	if err := tx.Preload("Team").Where("id = ?", killerID).First(&killer).Error; err != nil {
		tx.Rollback()
		return 0, err
	}
	if err := tx.Preload("Team").Where("id = ?", victimID).First(&victim).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	// 计算受害者实际可扣除的赏金（确保不为负）
	// 规则：受害者损失当前赏金的50%，全部转移给伤害者
	actualLoss := victim.Bounty * 0.5
	if actualLoss < 0 {
		actualLoss = 0
	}
	if victim.Bounty-actualLoss < 0 {
		actualLoss = victim.Bounty
	}

	// 伤害者获得受害者损失的全部赏金
	killerGain := actualLoss

	newKillerBounty := killer.Bounty + killerGain
	newVictimBounty := victim.Bounty - actualLoss

	// 更新伤害者赏金
	if err := tx.Model(&example.ExaPlayer{}).Where("id = ?", killerID).Update("bounty", newKillerBounty).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	// 更新受害者赏金
	if err := tx.Model(&example.ExaPlayer{}).Where("id = ?", victimID).Update("bounty", newVictimBounty).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	// 伤害者获得赏金记录
	killerRecord := example.ExaBountyRecord{
		PlayerID:    killerID,
		PlayerName:  killer.PlayerName,
		TeamID:      killer.TeamID,
		TeamName:    killer.Team.TeamName,
		ChangeType:  "kill",
		Amount:      killerGain,
		Balance:     newKillerBounty,
		Reason:      "击杀获取赏金",
		RelatedID:   victimID,
		RelatedName: victim.PlayerName,
	}
	if err := tx.Create(&killerRecord).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	// 受害者损失赏金记录
	victimRecord := example.ExaBountyRecord{
		PlayerID:    victimID,
		PlayerName:  victim.PlayerName,
		TeamID:      victim.TeamID,
		TeamName:    victim.Team.TeamName,
		ChangeType:  "killed",
		Amount:      -actualLoss,
		Balance:     newVictimBounty,
		Reason:      "被击杀损失赏金",
		RelatedID:   killerID,
		RelatedName: killer.PlayerName,
	}
	if err := tx.Create(&victimRecord).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := tx.Commit().Error; err != nil {
		return 0, err
	}

	// 同步更新双方战队的总赏金
	if killer.TeamID > 0 {
		go TeamServiceApp.RecalculateTeamBounty(killer.TeamID)
	}
	if victim.TeamID > 0 {
		go TeamServiceApp.RecalculateTeamBounty(victim.TeamID)
	}

	return killerGain, nil
}

func (p *PlayerService) Revive(playerID uint) (float64, error) {
	tx := global.GVA_DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var player example.ExaPlayer
	if err := tx.Preload("Team").Where("id = ?", playerID).First(&player).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	// 复活损失50%赏金，这部分进入公共赏金池
	// 确保赏金不会变成负数
	lostAmount := player.Bounty * 0.5
	if lostAmount < 0 {
		lostAmount = 0
	}
	if player.Bounty-lostAmount < 0 {
		lostAmount = player.Bounty
	}
	newBounty := player.Bounty - lostAmount

	if err := tx.Model(&example.ExaPlayer{}).Where("id = ?", playerID).Update("bounty", newBounty).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	record := example.ExaBountyRecord{
		PlayerID:   playerID,
		PlayerName: player.PlayerName,
		TeamID:     player.TeamID,
		TeamName:   player.Team.TeamName,
		ChangeType: "revive",
		Amount:     -lostAmount,
		Balance:    newBounty,
		Reason:     "复活损失赏金",
	}
	if err := tx.Create(&record).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	// 复活损失的赏金进入公共赏金池
	var pool example.ExaPublicBountyPool
	if err := tx.FirstOrCreate(&pool, example.ExaPublicBountyPool{PoolName: "公共赏金池"}).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	newPoolAmount := pool.TotalAmount + lostAmount
	if err := tx.Model(&pool).Update("total_amount", newPoolAmount).Error; err != nil {
		tx.Rollback()
		return 0, err
	}
	poolRecord := example.ExaBountyRecord{
		PlayerID:    0,
		PlayerName:  "系统",
		TeamID:      0,
		TeamName:    "公共赏金池",
		ChangeType:  "pool_add",
		Amount:      lostAmount,
		Balance:     newPoolAmount,
		Reason:      "复活进入赏金池",
		RelatedID:   playerID,
		RelatedName: player.PlayerName,
	}
	if err := tx.Create(&poolRecord).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := tx.Commit().Error; err != nil {
		return 0, err
	}

	// 同步更新玩家战队的总赏金
	if player.TeamID > 0 {
		go TeamServiceApp.RecalculateTeamBounty(player.TeamID)
	}

	return lostAmount, nil
}

// ClaimFromPool 从公共赏金池领取赏金
// 规则：玩家可以领取赏金池中50%的金额上限，领多少扣多少
func (p *PlayerService) ClaimFromPool(playerID uint, claimAmount float64) (float64, error) {
	tx := global.GVA_DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var player example.ExaPlayer
	if err := tx.Preload("Team").Where("id = ?", playerID).First(&player).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	var pool example.ExaPublicBountyPool
	if err := tx.FirstOrCreate(&pool, example.ExaPublicBountyPool{PoolName: "公共赏金池"}).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	// 检查池子余额是否足够
	if pool.TotalAmount <= 0 {
		tx.Rollback()
		return 0, nil
	}

	// 实际领取金额不能超过池子余额，且最多只能领取池子的50%
	maxClaim := pool.TotalAmount * 0.5
	if claimAmount > maxClaim {
		claimAmount = maxClaim
	}
	if claimAmount <= 0 {
		tx.Rollback()
		return 0, nil
	}

	// 玩家获得claimAmount，池子减少claimAmount
	newPlayerBounty := player.Bounty + claimAmount
	newPoolAmount := pool.TotalAmount - claimAmount

	if err := tx.Model(&example.ExaPlayer{}).Where("id = ?", playerID).Update("bounty", newPlayerBounty).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := tx.Model(&pool).Update("total_amount", newPoolAmount).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	// 玩家领取赏金记录
	record := example.ExaBountyRecord{
		PlayerID:    playerID,
		PlayerName:  player.PlayerName,
		TeamID:      player.TeamID,
		TeamName:    player.Team.TeamName,
		ChangeType:  "pool_claim",
		Amount:      claimAmount,
		Balance:     newPlayerBounty,
		Reason:      "从公共赏金池领取",
		RelatedID:   pool.ID,
		RelatedName: pool.PoolName,
	}
	if err := tx.Create(&record).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	// 赏金池减少记录
	poolRecord := example.ExaBountyRecord{
		PlayerID:    0,
		PlayerName:  "系统",
		TeamID:      0,
		TeamName:    "公共赏金池",
		ChangeType:  "pool_reduce",
		Amount:      -claimAmount,
		Balance:     newPoolAmount,
		Reason:      "玩家领取赏金",
		RelatedID:   playerID,
		RelatedName: player.PlayerName,
	}
	if err := tx.Create(&poolRecord).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := tx.Commit().Error; err != nil {
		return 0, err
	}

	if player.TeamID > 0 {
		go TeamServiceApp.RecalculateTeamBounty(player.TeamID)
	}
	return claimAmount, nil
}
