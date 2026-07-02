package example

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProxyRouter struct{}

// proxyCors 与原 server.js 一致，使用通配符 * 放行所有来源
func proxyCors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "content-type,t,s")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func (p *ProxyRouter) InitProxyRouter(Router *gin.RouterGroup) {
	// 同时注册 /proxy 和 /api/proxy 两种前缀，兼容直连和 nginx 转发
	for _, prefix := range []string{"proxy", "api/proxy"} {
		proxyRouter := Router.Group(prefix).Use(proxyCors())
		{
			proxyRouter.GET("", exaProxyApi.GetProxyInfo)
			proxyRouter.Any("/*path", exaProxyApi.Forward)
		}
	}
}
