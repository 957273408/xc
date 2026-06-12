package example

import (
	"errors"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/example"
)

type TeamService struct{}

var TeamServiceApp = new(TeamService)

func (t *TeamService) CreateTeam(team example.ExaTeam) error {
	return global.GVA_DB.Create(&team).Error
}

func (t *TeamService) DeleteTeam(id uint) error {
	return global.GVA_DB.Delete(&example.ExaTeam{}, id).Error
}

func (t *TeamService) UpdateTeam(team *example.ExaTeam) error {
	return global.GVA_DB.Save(team).Error
}

func (t *TeamService) GetTeam(id uint) (example.ExaTeam, error) {
	var team example.ExaTeam
	err := global.GVA_DB.Where("id = ?", id).First(&team).Error
	return team, err
}

func (t *TeamService) GetTeamList(info request.PageInfo) ([]example.ExaTeam, int64, error) {
	var teams []example.ExaTeam
	var total int64
	db := global.GVA_DB.Model(&example.ExaTeam{})
	err := db.Count(&total).Error
	if err != nil {
		return teams, total, err
	}
	err = db.Limit(info.PageSize).Offset(info.PageSize * (info.Page - 1)).Find(&teams).Error
	return teams, total, err
}

func (t *TeamService) SetTeamBounty(teamID uint, bounty float64) error {
	const minBounty = 0
	const maxBounty = 99999

	if bounty < minBounty || bounty > maxBounty {
		return errors.New("赏金必须在0到99999之间")
	}

	if bounty != float64(int64(bounty)) {
		return errors.New("赏金必须为整数")
	}

	return global.GVA_DB.Model(&example.ExaTeam{}).Where("id = ?", teamID).Update("total_bounty", bounty).Error
}

// RecalculateTeamBounty 重新计算并更新战队总赏金
// 总赏金 = 该战队所有选手当前赏金的总和
func (t *TeamService) RecalculateTeamBounty(teamID uint) error {
	var totalBounty float64
	err := global.GVA_DB.Model(&example.ExaPlayer{}).
		Select("COALESCE(SUM(bounty), 0)").
		Where("team_id = ?", teamID).
		Find(&totalBounty).Error
	if err != nil {
		return err
	}
	return t.SetTeamBounty(teamID, totalBounty)
}

// RecalculateAllTeamBounty 重新计算所有战队的总赏金
func (t *TeamService) RecalculateAllTeamBounty() error {
	var teams []example.ExaTeam
	err := global.GVA_DB.Find(&teams).Error
	if err != nil {
		return err
	}
	for _, team := range teams {
		if err := t.RecalculateTeamBounty(team.ID); err != nil {
			return err
		}
	}
	return nil
}
