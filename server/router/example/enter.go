package example

import (
	api "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
)

type RouterGroup struct {
	CustomerRouter
	AttachmentCategoryRouter
	FileUploadAndDownloadRouter
	TeamRouter
	PlayerRouter
	BountyRecordRouter
}

var (
	exaCustomerApi = api.ApiGroupApp.ExampleApiGroup.CustomerApi
	attachmentCategoryApi = api.ApiGroupApp.ExampleApiGroup.AttachmentCategoryApi
	exaFileUploadAndDownloadApi = api.ApiGroupApp.ExampleApiGroup.FileUploadAndDownloadApi
	exaTeamApi = api.ApiGroupApp.ExampleApiGroup.TeamApi
	exaPlayerApi = api.ApiGroupApp.ExampleApiGroup.PlayerApi
	exaBountyRecordApi = api.ApiGroupApp.ExampleApiGroup.BountyRecordApi
)