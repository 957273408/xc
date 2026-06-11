package example

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/example"
	exaReq "github.com/flipped-aurora/gin-vue-admin/server/model/example/request"
)

type PlayerService struct{}

var PlayerServiceApp = new(PlayerService)

func (p *PlayerService) CreatePlayer(player example.ExaPlayer) error {
	return global.GVA_DB.Create(&player).Error
}

func (p *PlayerService) DeletePlayer(id uint) error {
	return global.GVA_DB.Delete(&example.ExaPlayer{}, id).Error
}

func (p *PlayerService) UpdatePlayer(player *example.ExaPlayer) error {
	return global.GVA_DB.Save(player).Error
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
		return err
	}

	var totalAllocated float64
	for _, pb := range req.PlayerBounties {
		totalAllocated += pb.Amount
	}

	if totalAllocated > team.TotalBounty {
		tx.Rollback()
		return nil
	}

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

	return tx.Commit().Error
}

func (p *PlayerService) Kill(killerID, victimID uint) (float64, error) {
	tx := global.GVA_DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var killer, victim example.ExaPlayer
	if err := tx.Where("id = ?", killerID).First(&killer).Error; err != nil {
		tx.Rollback()
		return 0, err
	}
	if err := tx.Where("id = ?", victimID).First(&victim).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	stealAmount := victim.Bounty * 0.5

	newKillerBounty := killer.Bounty + stealAmount
	newVictimBounty := victim.Bounty - stealAmount

	if err := tx.Model(&example.ExaPlayer{}).Where("id = ?", killerID).Update("bounty", newKillerBounty).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := tx.Model(&example.ExaPlayer{}).Where("id = ?", victimID).Update("bounty", newVictimBounty).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	killerRecord := example.ExaBountyRecord{
		PlayerID:    killerID,
		PlayerName:  killer.PlayerName,
		TeamID:      killer.TeamID,
		TeamName:    killer.Team.TeamName,
		ChangeType:  "kill",
		Amount:      stealAmount,
		Balance:     newKillerBounty,
		Reason:      "击杀获取赏金",
		RelatedID:   victimID,
		RelatedName: victim.PlayerName,
	}
	if err := tx.Create(&killerRecord).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	victimRecord := example.ExaBountyRecord{
		PlayerID:    victimID,
		PlayerName:  victim.PlayerName,
		TeamID:      victim.TeamID,
		TeamName:    victim.Team.TeamName,
		ChangeType:  "killed",
		Amount:      -stealAmount,
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

	return stealAmount, nil
}

func (p *PlayerService) Revive(playerID uint) (float64, error) {
	tx := global.GVA_DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var player example.ExaPlayer
	if err := tx.Where("id = ?", playerID).First(&player).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	lostAmount := player.Bounty * 0.5
	newBounty := player.Bounty - lostAmount

	if err := tx.Model(&example.ExaPlayer{}).Where("id = ?", playerID).Update("bounty", newBounty).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	record := example.ExaBountyRecord{
		PlayerID:    playerID,
		PlayerName:  player.PlayerName,
		TeamID:      player.TeamID,
		TeamName:    player.Team.TeamName,
		ChangeType:  "revive",
		Amount:      -lostAmount,
		Balance:     newBounty,
		Reason:      "复活损失赏金",
	}
	if err := tx.Create(&record).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := tx.Commit().Error; err != nil {
		return 0, err
	}

	return lostAmount, nil
}
