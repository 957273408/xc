package example

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"go.uber.org/zap"
)

type ProxyService struct{}

var ProxyServiceApp = new(ProxyService)

type UpstreamResponse struct {
	Status      int
	ContentType string
	Body        interface{}
}

func (s *ProxyService) Forward(method, path string, query map[string]interface{}) (*UpstreamResponse, error) {
	proxyConfig := global.GVA_CONFIG.Proxy

	upstreamBaseURL := strings.TrimRight(proxyConfig.UpstreamBaseURL, "/")
	if upstreamBaseURL == "" {
		upstreamBaseURL = "http://match.biubiubiu.io:8099"
	}

	signKey := proxyConfig.SignKey
	if signKey == "" {
		signKey = "nCRmkbuTucUv"
	}

	signAlgorithm := proxyConfig.SignAlgorithm
	if signAlgorithm == "" {
		signAlgorithm = "sha256"
	}

	upstreamQueryMode := proxyConfig.UpstreamQueryMode
	if upstreamQueryMode == "" {
		upstreamQueryMode = "json-query"
	}

	requestTimeout := proxyConfig.RequestTimeout
	if requestTimeout <= 0 {
		requestTimeout = 30
	}

	jsonQuery := stableJSONStringify(query)
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	signature := createSignature(signKey, signAlgorithm, timestamp, path, jsonQuery)

	headers := map[string]string{
		"Accept": "application/json",
		"t":      timestamp,
		"s":      signature,
	}

	reqURL, body := buildUpstreamRequest(method, upstreamBaseURL, path, jsonQuery, upstreamQueryMode)

	var reqBody io.Reader
	if body != "" {
		reqBody = bytes.NewBufferString(body)
		headers["Content-Type"] = "application/json"
	}

	req, err := http.NewRequest(method, reqURL.String(), reqBody)
	if err != nil {
		global.GVA_LOG.Error("创建上游请求失败", zap.Error(err))
		return nil, fmt.Errorf("创建上游请求失败: %w", err)
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{Timeout: time.Duration(requestTimeout) * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		global.GVA_LOG.Error("上游请求失败",
			zap.String("method", method),
			zap.String("path", path),
			zap.Error(err))
		return nil, fmt.Errorf("上游请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		global.GVA_LOG.Error("读取上游响应失败", zap.Error(err))
		return nil, fmt.Errorf("读取上游响应失败: %w", err)
	}

	contentType := resp.Header.Get("Content-Type")

	return &UpstreamResponse{
		Status:      resp.StatusCode,
		ContentType: contentType,
		Body:        parseJSONIfPossible(string(respBody)),
	}, nil
}

func NormalizeProxyPath(pathname string) (string, error) {
	path := pathname
	if strings.HasPrefix(pathname, "/api/") {
		path = pathname[4:]
	}

	if !strings.HasPrefix(path, "/") {
		return "", fmt.Errorf("invalid upstream path")
	}

	return path, nil
}

func CoerceScalar(value string) interface{} {
	if value == "true" {
		return true
	}
	if value == "false" {
		return false
	}
	if value != "" {
		if n, err := strconv.ParseFloat(value, 64); err == nil {
			if n == float64(int64(n)) {
				return int64(n)
			}
			return n
		}
	}
	return value
}

func buildUpstreamRequest(method, baseURL, path, jsonQuery, queryMode string) (*url.URL, string) {
	u, _ := url.Parse(fmt.Sprintf("%s%s", baseURL, path))

	if method == http.MethodPost {
		return u, jsonQuery
	}

	// 与 Node.js server.js 一致：仅当 jsonQuery != "{}" 时才添加查询字符串
	if jsonQuery != "{}" {
		switch queryMode {
		case "json-param":
			q := u.Query()
			q.Set("q", jsonQuery)
			u.RawQuery = q.Encode()
		case "normal-query":
			var queryMap map[string]interface{}
			if err := json.Unmarshal([]byte(jsonQuery), &queryMap); err == nil {
				q := u.Query()
				for k, v := range queryMap {
					q.Set(k, fmt.Sprintf("%v", v))
				}
				u.RawQuery = q.Encode()
			}
		default:
			u.RawQuery = jsonQuery
		}
	}

	return u, ""
}

func createSignature(signKey, algorithm, timestamp, path, jsonQuery string) string {
	data := fmt.Sprintf("%s%s%s%s", signKey, timestamp, path, jsonQuery)

	switch algorithm {
	case "sha256":
		h := sha256.Sum256([]byte(data))
		return base64.StdEncoding.EncodeToString(h[:])
	default:
		h := sha256.Sum256([]byte(data))
		return base64.StdEncoding.EncodeToString(h[:])
	}
}

func stableJSONStringify(value map[string]interface{}) string {
	sorted := sortObjectKeys(value)
	b, _ := json.Marshal(sorted)
	return string(b)
}

func sortObjectKeys(value interface{}) interface{} {
	switch v := value.(type) {
	case []interface{}:
		result := make([]interface{}, len(v))
		for i, item := range v {
			result[i] = sortObjectKeys(item)
		}
		return result
	case map[string]interface{}:
		keys := make([]string, 0, len(v))
		for k := range v {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		result := make(map[string]interface{}, len(keys))
		for _, k := range keys {
			result[k] = sortObjectKeys(v[k])
		}
		return result
	default:
		return value
	}
}

func parseJSONIfPossible(text string) interface{} {
	if text == "" {
		return nil
	}

	var result interface{}
	if err := json.Unmarshal([]byte(text), &result); err != nil {
		return text
	}
	return result
}
