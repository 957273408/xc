package example

import (
	"encoding/json"
	"io"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/service/example"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ProxyApi struct{}

var knownEndpoints = []string{
	"/xdc/get_info",
	"/xdc/get_kill_info",
	"/xd/get_score_info",
	"/xd/get_template",
}

func (p *ProxyApi) GetProxyInfo(c *gin.Context) {
	proxyConfig := global.GVA_CONFIG.Proxy

	upstreamBaseURL := proxyConfig.UpstreamBaseURL
	if upstreamBaseURL == "" {
		upstreamBaseURL = "http://match.biubiubiu.io:8099"
	}

	response.OkWithDetailed(gin.H{
		"name":      "buibuibui-proxy",
		"upstream":  upstreamBaseURL,
		"endpoints": knownEndpoints,
		"usage": gin.H{
			"direct": "GET /proxy/xdc/get_info?warId=xxx",
			"post":   "POST /proxy/xd/get_score_info with JSON body",
		},
	}, "获取成功", c)
}

func (p *ProxyApi) Forward(c *gin.Context) {
	method := c.Request.Method
	if method != "GET" && method != "POST" {
		response.FailWithMessage("Only GET and POST are supported", c)
		return
	}

	targetPath := c.Param("path")
	if !strings.HasPrefix(targetPath, "/") {
		targetPath = "/" + targetPath
	}

	query, err := readRequestPayload(c)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	upstreamResp, err := proxyService.Forward(method, targetPath, query)
	if err != nil {
		global.GVA_LOG.Error("代理请求失败",
			zap.String("method", method),
			zap.String("path", targetPath),
			zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	writeUpstreamResponse(c, upstreamResp)
}

func readRequestPayload(c *gin.Context) (map[string]interface{}, error) {
	queryObject := map[string]interface{}{}
	for k, v := range c.Request.URL.Query() {
		if len(v) > 0 {
			queryObject[k] = example.CoerceScalar(v[0])
		}
	}

	if c.Request.Method == "GET" {
		return queryObject, nil
	}

	body, err := io.ReadAll(io.LimitReader(c.Request.Body, 1024*1024))
	if err != nil {
		return nil, err
	}
	defer c.Request.Body.Close()

	bodyText := strings.TrimSpace(string(body))
	if bodyText == "" {
		return queryObject, nil
	}

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(bodyText), &result); err != nil {
		return nil, err
	}

	if result == nil {
		return queryObject, nil
	}

	return result, nil
}

func writeUpstreamResponse(c *gin.Context, resp *example.UpstreamResponse) {
	isJSON := false
	var body []byte

	switch b := resp.Body.(type) {
	case string:
		body = []byte(b)
	default:
		isJSON = true
		body, _ = json.Marshal(resp.Body)
	}

	contentType := resp.ContentType
	if isJSON {
		contentType = "application/json; charset=utf-8"
	} else if contentType == "" {
		contentType = "text/plain; charset=utf-8"
	}

	c.Header("Content-Type", contentType)
	c.Status(resp.Status)
	c.Writer.Write(body)
}
