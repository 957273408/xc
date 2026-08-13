package request

// GetWarTableDownloadRequest 获取场次表格下载链接的请求参数
type GetWarTableDownloadRequest struct {
	WarID string `json:"warId" form:"warId" binding:"required"` // 战场ID
}
