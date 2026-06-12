package example

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	exaReq "github.com/flipped-aurora/gin-vue-admin/server/model/example/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type BountyRecordApi struct{}

func (b *BountyRecordApi) GetRecordList(c *gin.Context) {
	var req exaReq.GetBountyRecordListRequest
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	pageInfo := request.PageInfo{
		Page:     req.Page,
		PageSize: req.PageSize,
	}
	err = utils.Verify(pageInfo, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	records, total, err := bountyRecordService.GetRecordList(pageInfo, uint(req.PlayerID), uint(req.TeamID))
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     records,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

func (b *BountyRecordApi) GetPoolInfo(c *gin.Context) {
	pool, err := bountyRecordService.GetPoolInfo()
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(pool, "获取成功", c)
}

// GetTeamBountyRanking 获取队伍赏金排行榜
func (b *BountyRecordApi) GetTeamBountyRanking(c *gin.Context) {
	var req exaReq.GetTeamBountyRankingRequest
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	// 默认每页10条
	if req.PageSize == 0 {
		req.PageSize = 10
	}
	if req.Page == 0 {
		req.Page = 1
	}
	result, err := bountyRecordService.GetTeamBountyRanking(req.Page, req.PageSize)
	if err != nil {
		global.GVA_LOG.Error("获取队伍排行榜失败!", zap.Error(err))
		response.FailWithMessage("获取队伍排行榜失败", c)
		return
	}
	response.OkWithDetailed(result, "获取成功", c)
}

// GetPlayerBountyRanking 获取选手赏金排行榜
func (b *BountyRecordApi) GetPlayerBountyRanking(c *gin.Context) {
	var req exaReq.GetPlayerBountyRankingRequest
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	// 默认每页10条
	if req.PageSize == 0 {
		req.PageSize = 10
	}
	if req.Page == 0 {
		req.Page = 1
	}
	result, err := bountyRecordService.GetPlayerBountyRanking(req.Page, req.PageSize)
	if err != nil {
		global.GVA_LOG.Error("获取选手排行榜失败!", zap.Error(err))
		response.FailWithMessage("获取选手排行榜失败", c)
		return
	}
	response.OkWithDetailed(result, "获取成功", c)
}
