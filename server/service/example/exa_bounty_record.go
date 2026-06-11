package example

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/example"
)

type BountyRecordService struct{}

var BountyRecordServiceApp = new(BountyRecordService)

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
