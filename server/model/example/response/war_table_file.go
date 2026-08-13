package response

// WarTableDownloadResponse 场次表格下载链接的响应结构
type WarTableDownloadResponse struct {
	WarID       string `json:"warId"`       // 战场ID
	FileName    string `json:"fileName"`    // 表格文件名
	DownloadURL string `json:"downloadUrl"` // 表格下载地址
}
