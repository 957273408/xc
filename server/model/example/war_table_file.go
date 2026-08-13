package example

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// WarTableFile 场次表格文件映射，存储每个 warId 对应的表格下载链接
type WarTableFile struct {
	global.GVA_MODEL
	WarID       string `json:"warId" form:"warId" gorm:"column:war_id;comment:战场ID;size:64;not null;uniqueIndex"`
	FileName    string `json:"fileName" form:"fileName" gorm:"column:file_name;comment:表格文件名;size:255;not null"`
	DownloadURL string `json:"downloadUrl" form:"downloadUrl" gorm:"column:download_url;comment:表格下载地址;size:1000;not null"`
	Remark      string `json:"remark" form:"remark" gorm:"column:remark;comment:备注;size:500"`
}

func (WarTableFile) TableName() string {
	return "war_table_file"
}
