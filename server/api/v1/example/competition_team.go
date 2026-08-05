package example

import (
	"fmt"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/example"
	exaReq "github.com/flipped-aurora/gin-vue-admin/server/model/example/request"
	exaResp "github.com/flipped-aurora/gin-vue-admin/server/model/example/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CompetitionTeamApi struct{}

func (a *CompetitionTeamApi) CreateCompetitionTeam(c *gin.Context) {
	var req exaReq.CreateCompetitionTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	team := example.CompetitionTeam{
		TeamCode:  req.TeamCode,
		TeamName:  req.TeamName,
		TeamLogo:  req.TeamLogo,
		GroupName: req.GroupName,
	}

	if err := competitionTeamService.CreateCompetitionTeam(team); err != nil {
		global.GVA_LOG.Error("创建战队失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithMessage("创建成功", c)
}

func (a *CompetitionTeamApi) UpdateCompetitionTeam(c *gin.Context) {
	var req exaReq.UpdateCompetitionTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	team := &example.CompetitionTeam{
		GVA_MODEL: global.GVA_MODEL{ID: req.ID},
		TeamCode:  req.TeamCode,
		TeamName:  req.TeamName,
		TeamLogo:  req.TeamLogo,
		GroupName: req.GroupName,
	}

	if err := competitionTeamService.UpdateCompetitionTeam(team); err != nil {
		global.GVA_LOG.Error("更新战队失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithMessage("更新成功", c)
}

func (a *CompetitionTeamApi) DeleteCompetitionTeam(c *gin.Context) {
	var req exaReq.DeleteCompetitionTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := utils.Verify(req, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := competitionTeamService.DeleteCompetitionTeam(req.ID); err != nil {
		global.GVA_LOG.Error("删除战队失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithMessage("删除成功", c)
}

func (a *CompetitionTeamApi) GetCompetitionTeam(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := utils.Verify(req, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	team, err := competitionTeamService.GetCompetitionTeam(uint(req.ID))
	if err != nil {
		global.GVA_LOG.Error("获取战队失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}

	response.OkWithDetailed(team, "获取成功", c)
}

func (a *CompetitionTeamApi) GetCompetitionTeamList(c *gin.Context) {
	var pageInfo request.PageInfo
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := utils.Verify(pageInfo, utils.PageInfoVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	teams, total, err := competitionTeamService.GetCompetitionTeamList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取战队列表失败!", zap.Error(err))
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

func (a *CompetitionTeamApi) ImportExcel(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.FailWithMessage("请上传Excel文件", c)
		return
	}
	defer file.Close()

	mode := c.DefaultQuery("mode", "incremental")

	if mode != "incremental" && mode != "full" {
		response.FailWithMessage("无效的导入模式", c)
		return
	}

	result, err := excelImportService.ImportTeamsFromExcel(header, mode)
	if err != nil {
		global.GVA_LOG.Error("导入Excel失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithDetailed(result, fmt.Sprintf("导入完成: 成功%d条, 失败%d条", result.SuccessCount, result.FailCount), c)
}

func (a *CompetitionTeamApi) AddWarID(c *gin.Context) {
	var req exaReq.AddWarIDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	score, err := competitionTeamService.AddWarID(req.TeamID, req.WarID)
	if err != nil {
		global.GVA_LOG.Error("添加WarId失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithDetailed(score, "添加成功", c)
}

func (a *CompetitionTeamApi) GetTeamScores(c *gin.Context) {
	teamIDStr := c.Query("teamId")
	limitStr := c.DefaultQuery("limit", "0")

	var teamID uint
	fmt.Sscanf(teamIDStr, "%d", &teamID)

	var limit int
	fmt.Sscanf(limitStr, "%d", &limit)

	if teamID == 0 {
		response.FailWithMessage("战队ID不能为空", c)
		return
	}

	scores, err := competitionTeamService.GetTeamScores(teamID, limit)
	if err != nil {
		global.GVA_LOG.Error("获取积分记录失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}

	response.OkWithDetailed(scores, "获取成功", c)
}

func (a *CompetitionTeamApi) GetTeamRecentScores(c *gin.Context) {
	teamIDStr := c.Query("teamId")
	limitStr := c.DefaultQuery("limit", "4")

	var teamID uint
	fmt.Sscanf(teamIDStr, "%d", &teamID)

	var limit int
	fmt.Sscanf(limitStr, "%d", &limit)

	if teamID == 0 {
		response.FailWithMessage("战队ID不能为空", c)
		return
	}

	scores, err := competitionTeamService.GetTeamRecentScores(teamID, limit)
	if err != nil {
		global.GVA_LOG.Error("获取最近积分记录失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}

	response.OkWithDetailed(scores, "获取成功", c)
}

func (a *CompetitionTeamApi) GetTeamDetail(c *gin.Context) {
	teamIDStr := c.Query("teamId")

	var teamID uint
	fmt.Sscanf(teamIDStr, "%d", &teamID)

	if teamID == 0 {
		response.FailWithMessage("战队ID不能为空", c)
		return
	}

	detail, err := competitionTeamService.GetTeamDetail(teamID)
	if err != nil {
		global.GVA_LOG.Error("获取战队详情失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}

	response.OkWithDetailed(detail, "获取成功", c)
}

func (a *CompetitionTeamApi) GetAllTeamsScoreSummary(c *gin.Context) {
	summary, err := competitionTeamService.GetAllTeamsScoreSummary()
	if err != nil {
		global.GVA_LOG.Error("获取积分汇总失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}

	response.OkWithDetailed(summary, "获取成功", c)
}

// GetTeamScoreRanking 战队积分排名（包含最近4次积分变动记录）
func (a *CompetitionTeamApi) GetTeamScoreRanking(c *gin.Context) {
	groupName := c.Query("groupName")
	result, err := competitionTeamService.GetTeamScoreRanking(groupName)
	if err != nil {
		global.GVA_LOG.Error("获取积分排名失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}

	response.OkWithDetailed(result, "获取成功", c)
}

func (a *CompetitionTeamApi) DeleteTeamScore(c *gin.Context) {
	teamIDStr := c.Query("teamId")
	warID := c.Query("warId")

	var teamID uint
	fmt.Sscanf(teamIDStr, "%d", &teamID)

	if teamID == 0 || warID == "" {
		response.FailWithMessage("参数不完整", c)
		return
	}

	if err := competitionTeamService.DeleteTeamScore(teamID, warID); err != nil {
		global.GVA_LOG.Error("删除积分记录失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithMessage("删除成功", c)
}

func (a *CompetitionTeamApi) UpdateTeamScore(c *gin.Context) {
	var req exaReq.UpdateScoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := competitionTeamService.UpdateTeamScore(
		req.ID, req.TeamID,
		req.Rank, req.KillCount, req.RankScore, req.KillScore, req.TotalScore,
		req.BountyCoin,
		req.SettleTime,
	); err != nil {
		global.GVA_LOG.Error("修改积分记录失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithMessage("修改成功", c)
}

// GetPublicWarScores 公开接口：获取指定WarId下各战队当场积分（无需鉴权）
func (a *CompetitionTeamApi) GetPublicWarScores(c *gin.Context) {
	warID := c.Query("warId")
	if warID == "" {
		response.FailWithMessage("缺少warId参数", c)
		return
	}

	result, err := competitionTeamService.GetPublicWarScores(warID)
	if err != nil {
		global.GVA_LOG.Error("获取WarId积分失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithDetailed(result, "获取成功", c)
}

// GetPublicWarBounty 公开接口：获取指定WarId下各战队选手赏金分配（无需鉴权）
func (a *CompetitionTeamApi) GetPublicWarBounty(c *gin.Context) {
	warID := c.Query("warId")
	if warID == "" {
		response.FailWithMessage("缺少warId参数", c)
		return
	}
	result, err := competitionTeamService.GetPublicWarBounty(warID)
	if err != nil {
		global.GVA_LOG.Error("获取赏金分配失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(result, "获取成功", c)
}

// GetPublicTeamList 公开接口：获取所有战队的简化列表（无需鉴权）
// @Summary 获取所有战队列表
// @Tags CompetitionTeam
// @Produce json
// @Success 200 {object} response.Response{data=response.PublicTeamListResponse}
// @Router /competitionTeam/public/teamList [get]
func (a *CompetitionTeamApi) GetPublicTeamList(c *gin.Context) {
	result, err := competitionTeamService.GetPublicTeamList()
	if err != nil {
		global.GVA_LOG.Error("获取战队列表失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(result, "获取成功", c)
}

// CalculateWarIDForAllTeams 批量计算指定WarId下所有战队积分（预览，不保存）
func (a *CompetitionTeamApi) CalculateWarIDForAllTeams(c *gin.Context) {
	var req exaReq.WarIDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	result, err := competitionTeamService.CalculateWarIDForAllTeams(req.WarID)
	if err != nil {
		global.GVA_LOG.Error("批量计算WarId积分失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithDetailed(result, "计算完成", c)
}

// ConfirmWarIDScores 确认并保存批量积分
func (a *CompetitionTeamApi) ConfirmWarIDScores(c *gin.Context) {
	var req exaResp.ConfirmWarIDScoresRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	result, err := competitionTeamService.ConfirmWarIDScores(req.WarID, req.TeamIDs)
	if err != nil {
		global.GVA_LOG.Error("确认批量积分失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	msg := fmt.Sprintf("确认完成: 成功%d条, 失败%d条", result.SuccessCount, result.FailCount)
	response.OkWithDetailed(result, msg, c)
}
