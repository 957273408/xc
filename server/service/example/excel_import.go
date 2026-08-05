package example

import (
	"fmt"
	"io"
	"mime/multipart"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/example"
	exaResp "github.com/flipped-aurora/gin-vue-admin/server/model/example/response"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
)

type ExcelImportService struct{}

var ExcelImportServiceApp = new(ExcelImportService)

type TeamExcelRow struct {
	TeamCode string
	TeamName string
	TeamLogo string
	RawRow   int
}

func (s *ExcelImportService) ImportTeamsFromExcel(file *multipart.FileHeader, mode string) (*exaResp.ImportResultResponse, error) {
	result := &exaResp.ImportResultResponse{
		SuccessCount: 0,
		FailCount:    0,
		Errors:       []string{},
	}

	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %v", err)
	}
	defer src.Close()

	rows, err := s.parseExcel(src)
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return result, nil
	}

	for _, row := range rows {
		if err := s.processTeamRow(row, mode); err != nil {
			result.FailCount++
			result.Errors = append(result.Errors, fmt.Sprintf("第%d行: %v", row.RawRow+1, err))
			global.GVA_LOG.Warn("导入战队失败",
				zap.Int("row", row.RawRow+1),
				zap.String("teamCode", row.TeamCode),
				zap.Error(err))
		} else {
			result.SuccessCount++
		}
	}

	return result, nil
}

func (s *ExcelImportService) parseExcel(src io.Reader) ([]TeamExcelRow, error) {
	f, err := excelize.OpenReader(src)
	if err != nil {
		return nil, fmt.Errorf("解析Excel失败: %v", err)
	}
	defer f.Close()

	sheetName := f.GetSheetName(0)
	if sheetName == "" {
		return nil, fmt.Errorf("Excel文件为空")
	}

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("读取Excel数据失败: %v", err)
	}

	if len(rows) < 2 {
		return nil, fmt.Errorf("Excel数据不足，需要标题行和至少一行数据")
	}

	headerRow := rows[0]
	columnMap := s.mapColumns(headerRow)

	var result []TeamExcelRow
	for i, row := range rows[1:] {
		if len(strings.TrimSpace(strings.Join(row, ""))) == 0 {
			continue
		}

		teamCode := s.getCellValue(row, columnMap["teamCode"])
		teamName := s.getCellValue(row, columnMap["teamName"])
		teamLogo := s.getCellValue(row, columnMap["teamLogo"])

		if teamCode == "" && teamName == "" {
			continue
		}

		result = append(result, TeamExcelRow{
			TeamCode: teamCode,
			TeamName: teamName,
			TeamLogo: teamLogo,
			RawRow:   i + 1,
		})
	}

	return result, nil
}

func (s *ExcelImportService) mapColumns(headerRow []string) map[string]int {
	columnMap := map[string]int{
		"teamCode": -1,
		"teamName": -1,
		"teamLogo": -1,
	}

	for i, header := range headerRow {
		header = strings.TrimSpace(header)
		switch {
		case containsAny(header, "战队标识", "队伍标识", "战队代码", "TeamCode", "teamCode", "code", "缩写"):
			columnMap["teamCode"] = i
		case containsAny(header, "战队名称", "队伍名称", "战队", "TeamName", "teamName", "name", "名称"):
			columnMap["teamName"] = i
		case containsAny(header, "战队Logo", "队伍Logo", "Logo", "logo", "队标"):
			columnMap["teamLogo"] = i
		}
	}

	return columnMap
}

func (s *ExcelImportService) getCellValue(row []string, index int) string {
	if index < 0 || index >= len(row) {
		return ""
	}
	return strings.TrimSpace(row[index])
}

func containsAny(s string, substrs ...string) bool {
	for _, sub := range substrs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}

func (s *ExcelImportService) processTeamRow(row TeamExcelRow, mode string) error {
	if row.TeamCode == "" {
		return fmt.Errorf("战队标识不能为空")
	}
	if row.TeamName == "" {
		return fmt.Errorf("战队名称不能为空")
	}

	var existing example.CompetitionTeam
	err := global.GVA_DB.Where("team_code = ?", row.TeamCode).First(&existing).Error

	if mode == "incremental" {
		if err == nil {
			return fmt.Errorf("战队标识 '%s' 已存在，跳过敏", row.TeamCode)
		}

		team := example.CompetitionTeam{
			TeamCode: row.TeamCode,
			TeamName: row.TeamName,
			TeamLogo: row.TeamLogo,
		}
		return global.GVA_DB.Create(&team).Error
	}

	if mode == "full" {
		if err == nil {
			existing.TeamName = row.TeamName
			if row.TeamLogo != "" {
				existing.TeamLogo = row.TeamLogo
			}
			return global.GVA_DB.Save(&existing).Error
		}

		team := example.CompetitionTeam{
			TeamCode: row.TeamCode,
			TeamName: row.TeamName,
			TeamLogo: row.TeamLogo,
		}
		return global.GVA_DB.Create(&team).Error
	}

	return fmt.Errorf("未知的导入模式: %s", mode)
}
