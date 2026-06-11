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

type PlayerApi struct{}

func (p *PlayerApi) CreatePlayer(c *gin.Context) {
	var req exaReq.CreatePlayerRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	player := example.ExaPlayer{
		PlayerName: req.PlayerName,
		UID:        req.UID,
		TeamID:     req.TeamID,
		Bounty:     req.Bounty,
	}
	err = playerService.CreatePlayer(player)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

func (p *PlayerApi) DeletePlayer(c *gin.Context) {
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
	err = playerService.DeletePlayer(uint(req.ID))
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

func (p *PlayerApi) UpdatePlayer(c *gin.Context) {
	var req exaReq.UpdatePlayerRequest
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
	player := example.ExaPlayer{
		GVA_MODEL:   global.GVA_MODEL{ID: req.ID},
		PlayerName: req.PlayerName,
		UID:        req.UID,
		TeamID:     req.TeamID,
		Bounty:     req.Bounty,
	}
	err = playerService.UpdatePlayer(&player)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

func (p *PlayerApi) GetPlayer(c *gin.Context) {
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
	player, err := playerService.GetPlayer(uint(req.ID))
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(player, "获取成功", c)
}

func (p *PlayerApi) GetPlayerList(c *gin.Context) {
	var pageInfo request.PageInfo
	var teamID request.GetById
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	c.ShouldBindQuery(&teamID)
	err = utils.Verify(pageInfo, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	players, total, err := playerService.GetPlayerList(pageInfo, uint(teamID.ID))
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     players,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

func (p *PlayerApi) AllocateBounty(c *gin.Context) {
	var req exaReq.AllocateBountyRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = playerService.AllocateBounty(req)
	if err != nil {
		global.GVA_LOG.Error("分配失败!", zap.Error(err))
		response.FailWithMessage("分配失败", c)
		return
	}
	response.OkWithMessage("分配成功", c)
}

func (p *PlayerApi) Kill(c *gin.Context) {
	var req exaReq.KillRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	amount, err := playerService.Kill(req.KillerID, req.VictimID)
	if err != nil {
		global.GVA_LOG.Error("击杀处理失败!", zap.Error(err))
		response.FailWithMessage("击杀处理失败", c)
		return
	}
	response.OkWithDetailed(gin.H{"amount": amount}, "击杀成功", c)
}

func (p *PlayerApi) Revive(c *gin.Context) {
	var req exaReq.ReviveRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	amount, err := playerService.Revive(req.PlayerID)
	if err != nil {
		global.GVA_LOG.Error("复活处理失败!", zap.Error(err))
		response.FailWithMessage("复活处理失败", c)
		return
	}
	response.OkWithDetailed(gin.H{"lostAmount": amount}, "复活成功", c)
}
