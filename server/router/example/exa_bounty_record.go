package example

import (
	"github.com/gin-gonic/gin"
)

type BountyRecordRouter struct{}

func (b *BountyRecordRouter) InitBountyRecordRouter(Router *gin.RouterGroup) {
	bountyRecordRouter := Router.Group("bountyRecord")
	{
		bountyRecordRouter.GET("recordList", exaBountyRecordApi.GetRecordList)             // 获取赏金记录列表
		bountyRecordRouter.GET("poolInfo", exaBountyRecordApi.GetPoolInfo)                 // 获取公共赏金池信息
		bountyRecordRouter.GET("teamRanking", exaBountyRecordApi.GetTeamBountyRanking)     // 获取队伍赏金排行榜
		bountyRecordRouter.GET("playerRanking", exaBountyRecordApi.GetPlayerBountyRanking) // 获取选手赏金排行榜
	}
}
