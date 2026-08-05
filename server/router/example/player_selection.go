package example

import (
	"github.com/gin-gonic/gin"
)

type PlayerSelectionRouter struct{}

func (p *PlayerSelectionRouter) InitPlayerSelectionRouter(Router, PublicRouter *gin.RouterGroup) {
	publicGroup := PublicRouter.Group("playerSelection")
	{
		publicGroup.GET("warPlayers", playerSelectionApi.GetWarPlayers)
		publicGroup.POST("save", playerSelectionApi.SavePlayerSelection)
		publicGroup.GET("get", playerSelectionApi.GetPlayerSelection)
		publicGroup.POST("multiWarPlayers", playerSelectionApi.GetMultiWarPlayers)
		publicGroup.GET("latest", playerSelectionApi.GetLatestSelection)
	}
}
