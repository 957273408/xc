package example

import (
	"github.com/gin-gonic/gin"
)

type BountyRecordRouter struct{}

func (b *BountyRecordRouter) InitBountyRecordRouter(Router *gin.RouterGroup) {
	bountyRecordRouter := Router.Group("bountyRecord")
	{
		bountyRecordRouter.GET("recordList", exaBountyRecordApi.GetRecordList) // 获取赏金记录列表
	}
}