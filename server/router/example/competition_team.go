package example

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type CompetitionTeamRouter struct{}

func (c *CompetitionTeamRouter) InitCompetitionTeamRouter(Router, PublicRouter *gin.RouterGroup) {
	router := Router.Group("competitionTeam").Use(middleware.OperationRecord())
	routerWithoutRecord := Router.Group("competitionTeam")
	{
		router.POST("competitionTeam", competitionTeamApi.CreateCompetitionTeam)
		router.PUT("competitionTeam", competitionTeamApi.UpdateCompetitionTeam)
		router.DELETE("competitionTeam", competitionTeamApi.DeleteCompetitionTeam)
		router.POST("importExcel", competitionTeamApi.ImportExcel)
		router.POST("addWarID", competitionTeamApi.AddWarID)
		router.POST("calculateWarID", competitionTeamApi.CalculateWarIDForAllTeams)
		router.POST("confirmWarID", competitionTeamApi.ConfirmWarIDScores)
		router.PUT("updateScore", competitionTeamApi.UpdateTeamScore)
		router.DELETE("deleteScore", competitionTeamApi.DeleteTeamScore)
	}
	{
		routerWithoutRecord.GET("competitionTeam", competitionTeamApi.GetCompetitionTeam)
		routerWithoutRecord.GET("competitionTeamList", competitionTeamApi.GetCompetitionTeamList)
		routerWithoutRecord.GET("scores", competitionTeamApi.GetTeamScores)
		routerWithoutRecord.GET("recentScores", competitionTeamApi.GetTeamRecentScores)
		routerWithoutRecord.GET("detail", competitionTeamApi.GetTeamDetail)
		routerWithoutRecord.GET("allScores", competitionTeamApi.GetAllTeamsScoreSummary)
		routerWithoutRecord.GET("warTable", competitionTeamApi.GetWarTableDownload)
	}
	// 公开路由（无鉴权）
	publicGroup := PublicRouter.Group("competitionTeam")
	{
		publicGroup.GET("public/warScores", competitionTeamApi.GetPublicWarScores)
		publicGroup.GET("public/warBounty", competitionTeamApi.GetPublicWarBounty)
		publicGroup.GET("public/teamList", competitionTeamApi.GetPublicTeamList)
		publicGroup.GET("public/warTable", competitionTeamApi.GetWarTableDownload)
		publicGroup.GET("ranking", competitionTeamApi.GetTeamScoreRanking)
		publicGroup.GET("bountyRanking", competitionTeamApi.GetTeamBountyRanking)
	}
}
