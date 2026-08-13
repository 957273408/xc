package example

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

// GetWarTableDownload 获取指定场次的表格下载
// @Summary      获取指定场次Excel表格
// @Description  根据warId从代理获取战场数据，生成Excel文件(含3个Sheet:战场基本信息/选手明细数据/队伍明细数据)并自动下载
// @Tags         WarTable
// @Accept       json
// @Produce      application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Param        warId query string true "战场ID"
// @Success      200  {file}   binary  "Excel文件流"
// @Failure      400  {object} response.Response  "warId参数错误或不存在"
// @Router       /competitionTeam/public/warTable [get]
func (a *CompetitionTeamApi) GetWarTableDownload(c *gin.Context) {
	warID := c.Query("warId")
	if warID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 7,
			"msg":  "warId参数不能为空",
		})
		return
	}

	buf, filename, err := warTableExcelService.GenerateWarTableExcel(warID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 7,
			"msg":  err.Error(),
		})
		return
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename*=UTF-8''"+url.QueryEscape(filename))
	c.Header("Content-Transfer-Encoding", "binary")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf.Bytes())
}
