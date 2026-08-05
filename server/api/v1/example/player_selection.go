package example

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	exaReq "github.com/flipped-aurora/gin-vue-admin/server/model/example/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PlayerSelectionApi struct{}

// GetWarPlayers 公开接口：获取指定WarId下的所有玩家列表
// @Summary 获取战场玩家列表
// @Tags PlayerSelection
// @Accept json
// @Produce json
// @Param warId query string true "战场ID"
// @Success 200 {object} response.Response{data=response.WarPlayersResponse}
// @Router /playerSelection/warPlayers [get]
func (p *PlayerSelectionApi) GetWarPlayers(c *gin.Context) {
	var req exaReq.GetWarPlayersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	result, err := playerSelectionService.GetWarPlayers(req.WarID)
	if err != nil {
		global.GVA_LOG.Error("获取战场玩家列表失败", zap.String("warId", req.WarID), zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithDetailed(result, "获取成功", c)
}

// SavePlayerSelection 公开接口：保存玩家选择（5人+附加统计）
// @Summary 保存玩家选择
// @Tags PlayerSelection
// @Accept json
// @Produce json
// @Param data body exaReq.SavePlayerSelectionRequest true "保存请求"
// @Success 200 {object} response.Response{data=response.PlayerSelectionResponse}
// @Router /playerSelection/save [post]
func (p *PlayerSelectionApi) SavePlayerSelection(c *gin.Context) {
	var req exaReq.SavePlayerSelectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	result, err := playerSelectionService.SavePlayerSelection(&req)
	if err != nil {
		global.GVA_LOG.Error("保存玩家选择失败", zap.String("warId", req.WarID), zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithDetailed(result, "保存成功", c)
}

// GetPlayerSelection 公开接口：按SessionKey读取已保存的选择
// @Summary 获取已保存的玩家选择
// @Tags PlayerSelection
// @Accept json
// @Produce json
// @Param sessionKey query string true "会话标识"
// @Success 200 {object} response.Response{data=response.PlayerSelectionResponse}
// @Router /playerSelection/get [get]
func (p *PlayerSelectionApi) GetPlayerSelection(c *gin.Context) {
	var req exaReq.GetPlayerSelectionRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	result, err := playerSelectionService.GetPlayerSelection(req.SessionKey)
	if err != nil {
		global.GVA_LOG.Error("获取玩家选择失败", zap.String("sessionKey", req.SessionKey), zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithDetailed(result, "获取成功", c)
}

// GetMultiWarPlayers 公开接口：传入多个WarId获取汇总后的选手数据
// @Summary 多场选手数据汇总
// @Tags PlayerSelection
// @Accept json
// @Produce json
// @Param data body request.GetMultiWarPlayersRequest true "多场请求"
// @Success 200 {object} response.Response{data=response.MultiWarPlayersResponse}
// @Router /playerSelection/multiWarPlayers [post]
func (p *PlayerSelectionApi) GetMultiWarPlayers(c *gin.Context) {
	var req exaReq.GetMultiWarPlayersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	result, err := playerSelectionService.GetMultiWarPlayers(req.WarIDs)
	if err != nil {
		global.GVA_LOG.Error("多场汇总失败", zap.Any("warIds", req.WarIDs), zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithDetailed(result, "获取成功", c)
}

// GetMultiWarTop5 多场Top5统计：淘汰数、爆头率、命中率、伤害量各前5
// @Summary 多场Top5统计
// @Tags PlayerSelection
// @Accept json
// @Produce json
// @Param data body request.GetMultiWarPlayersRequest true "多场请求"
// @Success 200 {object} response.Response{data=response.MultiWarTop5Response}
// @Router /playerSelection/multiWarTop5 [post]
func (p *PlayerSelectionApi) GetMultiWarTop5(c *gin.Context) {
	var req exaReq.GetMultiWarPlayersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	result, err := playerSelectionService.GetMultiWarTop5(req.WarIDs)
	if err != nil {
		global.GVA_LOG.Error("多场Top5统计失败", zap.Any("warIds", req.WarIDs), zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithDetailed(result, "获取成功", c)
}

// GetLatestSelection 公开接口：获取最新保存的选择数据（无需sessionKey）
// @Summary 获取最新保存的玩家选择
// @Tags PlayerSelection
// @Produce json
// @Success 200 {object} response.Response{data=response.LatestSelectionResponse}
// @Router /playerSelection/latest [get]
func (p *PlayerSelectionApi) GetLatestSelection(c *gin.Context) {
	result, err := playerSelectionService.GetLatestSelection()
	if err != nil {
		global.GVA_LOG.Error("获取最新玩家选择失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithDetailed(result, "获取成功", c)
}
