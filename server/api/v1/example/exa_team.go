package example

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/example"
	exaReq "github.com/flipped-aurora/gin-vue-admin/server/model/example/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type TeamApi struct{}

func (t *TeamApi) CreateTeam(c *gin.Context) {
	var req exaReq.CreateTeamRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	team := example.ExaTeam{
		TeamName:    req.TeamName,
		TotalBounty: req.TotalBounty,
	}
	err = teamService.CreateTeam(team)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

func (t *TeamApi) DeleteTeam(c *gin.Context) {
	var req request.GetById
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(req, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = teamService.DeleteTeam(uint(req.ID))
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

func (t *TeamApi) UpdateTeam(c *gin.Context) {
	var req exaReq.UpdateTeamRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(req, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	team := example.ExaTeam{
		GVA_MODEL:   global.GVA_MODEL{ID: req.ID},
		TeamName:    req.TeamName,
		TotalBounty: req.TotalBounty,
	}
	err = teamService.UpdateTeam(&team)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

func (t *TeamApi) GetTeam(c *gin.Context) {
	var req request.GetById
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(req, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	team, err := teamService.GetTeam(uint(req.ID))
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(team, "获取成功", c)
}

func (t *TeamApi) GetTeamList(c *gin.Context) {
	var pageInfo request.PageInfo
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(pageInfo, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	teams, total, err := teamService.GetTeamList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     teams,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

func (t *TeamApi) SetTeamBounty(c *gin.Context) {
	var req exaReq.SetTeamBountyRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		global.GVA_LOG.Error("设置战队赏金请求参数解析失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	const minBounty = 0
	const maxBounty = 99999

	if req.Bounty < minBounty {
		global.GVA_LOG.Error("设置战队赏金失败: 赏金不能为负数",
			zap.Uint("teamID", req.TeamID),
			zap.Float64("bounty", req.Bounty))
		response.FailWithMessage("赏金不能为负数，最小允许值为0", c)
		return
	}

	if req.Bounty > maxBounty {
		global.GVA_LOG.Error("设置战队赏金失败: 赏金超过最大限制",
			zap.Uint("teamID", req.TeamID),
			zap.Float64("bounty", req.Bounty),
			zap.Float64("maxBounty", maxBounty))
		response.FailWithMessage("赏金超过最大限制，最大允许值为99999", c)
		return
	}

	if req.Bounty != float64(int64(req.Bounty)) {
		global.GVA_LOG.Error("设置战队赏金失败: 赏金必须为整数",
			zap.Uint("teamID", req.TeamID),
			zap.Float64("bounty", req.Bounty))
		response.FailWithMessage("赏金必须为整数", c)
		return
	}

	global.GVA_LOG.Info("设置战队赏金",
		zap.Uint("teamID", req.TeamID),
		zap.Float64("bounty", req.Bounty))

	err = teamService.SetTeamBounty(req.TeamID, req.Bounty)
	if err != nil {
		global.GVA_LOG.Error("设置失败!", zap.Uint("teamID", req.TeamID), zap.Error(err))
		response.FailWithMessage("设置失败", c)
		return
	}
	response.OkWithMessage("设置成功", c)
}
