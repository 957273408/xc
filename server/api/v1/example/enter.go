package example

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	CustomerApi
	AttachmentCategoryApi
	FileUploadAndDownloadApi
	TeamApi
	PlayerApi
	BountyRecordApi
	ProxyApi
	CompetitionTeamApi
	PlayerSelectionApi
}

var (
	customerService              = service.ServiceGroupApp.ExampleServiceGroup.CustomerService
	attachmentCategoryService    = service.ServiceGroupApp.ExampleServiceGroup.AttachmentCategoryService
	fileUploadAndDownloadService = service.ServiceGroupApp.ExampleServiceGroup.FileUploadAndDownloadService
	teamService                  = service.ServiceGroupApp.ExampleServiceGroup.TeamService
	playerService                = service.ServiceGroupApp.ExampleServiceGroup.PlayerService
	bountyRecordService          = service.ServiceGroupApp.ExampleServiceGroup.BountyRecordService
	proxyService                 = service.ServiceGroupApp.ExampleServiceGroup.ProxyService
	competitionTeamService       = service.ServiceGroupApp.ExampleServiceGroup.CompetitionTeamService
	excelImportService           = service.ServiceGroupApp.ExampleServiceGroup.ExcelImportService
	playerSelectionService       = service.ServiceGroupApp.ExampleServiceGroup.PlayerSelectionService
	warTableExcelService         = service.ServiceGroupApp.ExampleServiceGroup.WarTableExcelService
)
