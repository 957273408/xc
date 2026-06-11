package example

import (
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
	return global.GVA_DB.Model(&example.ExaTeam{}).Where("id = ?", teamID).Update("total_bounty", bounty).Error
}
