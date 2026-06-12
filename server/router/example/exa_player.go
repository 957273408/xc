package example

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type PlayerRouter struct{}

func (p *PlayerRouter) InitPlayerRouter(Router *gin.RouterGroup) {
	playerRouter := Router.Group("player").Use(middleware.OperationRecord())
	playerRouterWithoutRecord := Router.Group("player")
	{
		playerRouter.POST("player", exaPlayerApi.CreatePlayer)           // 创建选手
		playerRouter.PUT("player", exaPlayerApi.UpdatePlayer)            // 更新选手
		playerRouter.DELETE("player", exaPlayerApi.DeletePlayer)         // 删除选手
		playerRouter.POST("allocateBounty", exaPlayerApi.AllocateBounty) // 战队内部赏金分配
		playerRouter.POST("kill", exaPlayerApi.Kill)                     // 击杀处理
		playerRouter.POST("revive", exaPlayerApi.Revive)                 // 复活处理
		playerRouter.POST("claimFromPool", exaPlayerApi.ClaimFromPool)   // 从赏金池领取
	}
	{
		playerRouterWithoutRecord.GET("player", exaPlayerApi.GetPlayer)         // 获取单一选手信息
		playerRouterWithoutRecord.GET("playerList", exaPlayerApi.GetPlayerList) // 获取选手列表
	}
}
