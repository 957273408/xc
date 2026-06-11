package example

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type TeamRouter struct{}

func (t *TeamRouter) InitTeamRouter(Router *gin.RouterGroup) {
	teamRouter := Router.Group("team").Use(middleware.OperationRecord())
	teamRouterWithoutRecord := Router.Group("team")
	{
		teamRouter.POST("team", exaTeamApi.CreateTeam)           // 创建战队
		teamRouter.PUT("team", exaTeamApi.UpdateTeam)            // 更新战队
		teamRouter.DELETE("team", exaTeamApi.DeleteTeam)         // 删除战队
		teamRouter.POST("setBounty", exaTeamApi.SetTeamBounty)   // 设置战队初始赏金
	}
	{
		teamRouterWithoutRecord.GET("team", exaTeamApi.GetTeam)         // 获取单一战队信息
		teamRouterWithoutRecord.GET("teamList", exaTeamApi.GetTeamList) // 获取战队列表
	}
}