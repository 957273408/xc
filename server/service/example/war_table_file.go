package example

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/example"
	"github.com/xuri/excelize/v2"
)

// WarTableExcelService 战场表格Excel生成服务
type WarTableExcelService struct{}

var WarTableExcelServiceApp = new(WarTableExcelService)

// ============================================================
// 字段映射：代理接口字段名 → 中文表头
// ============================================================

var playerFieldMap = []struct {
	Key   string
	Label string
}{
	{"nickName", "玩家昵称"},
	{"uID", "玩家uid"},
	{"teamID", "队伍ID"},
	{"teamName", "队伍名"},
	{"totalKill", "淘汰数"},
	{"lostRoleRank", "队伍排名"},
	{"damage", "造成伤害"},
	{"inDamage", "受到伤害"},
	{"heal", "治疗量"},
	{"headShot", "爆头击杀数"},
	{"vehicleKill", "载具击杀数"},
	{"airDrop", "拾取空投数量"},
	{"maxKillDistance", "最远击杀距离"},
	{"maxHitDownDistance", "最远击倒距离"},
	{"survive", "生存时长(秒)"},
	{"helmetLevel", "头盔等级"},
	{"vestLevel", "防具等级"},
	{"bagLevel", "背包等级"},
	{"bulletTotalNum", "射击子弹总数"},
	{"hitBulletNum", "命中子弹总数"},
	{"hitHeadBulletNum", "命中头部子弹总数"},
	{"moveDistance", "行进距离"},
	{"driveDistance", "载具行驶距离(米)"},
	{"marchDistance", "非开车距离(米)"},
	{"assists", "助攻数"},
	{"rescue", "救援队友数"},
	{"hook", "拉钩钩使用次数"},
	{"big", "使用变大大次数"},
	{"fire", "火球使用次数"},
	{"fireDamage", "火球造成伤害量"},
	{"fireKill", "使用火球击倒数"},
	{"cannon", "传送大炮使用次数"},
	{"wormHole", "虫洞手雷使用次数"},
	{"healBomb", "治疗弹使用次数"},
	{"healBombHeal", "治疗弹治疗量"},
	{"shelter", "战术掩体使用次数"},
	{"cloudBomb", "云雾弹使用次数"},
	{"bomb", "手雷使用次数"},
	{"bombDamage", "手雷造成伤害量"},
	{"bombKill", "使用手雷击倒数"},
	{"bombMax", "手雷MAX使用次数"},
	{"bombMaxDamage", "手雷MAX造成伤害量"},
	{"bombMaxKill", "手雷MAX击倒数"},
	{"pickSwordNum", "圣剑拾取次数"},
	{"swordKill", "圣剑击杀"},
	{"swordDamage", "使用圣剑造成伤害"},
	{"mechaDamage", "使用机甲造成伤害"},
	{"mechaKill", "使用机甲淘汰数"},
	{"mechaDistance", "机甲移动距离"},
	{"dragonDamage", "使用呆呆龙造成伤害"},
	{"dragonKill", "呆呆龙淘汰数"},
	{"dragonDistance", "呆呆龙移动距离"},
	{"xcc", "变成小肠肠次数"},
	{"xccDistance", "变成小肠肠后移动距离"},
	{"hitWeak", "击倒敌人次数"},
	{"useResurrectionNum", "使用复活机次数"},
	{"allGunDamage", "枪械伤害"},
	{"trexKingDistance", "凶凶龙累计移动距离"},
	{"peterosaurDistance", "飘飘龙累计移动距离"},
	{"triceratopDistance", "憨憨龙累计移动距离"},
	{"raptorsDistance", "奔奔龙累计移动距离"},
	{"trexKingDamage", "凶凶龙累计造成伤害"},
	{"trexKingKill", "凶凶龙击杀数"},
	{"beUsedDragonNum", "呆呆龙登场次数"},
	{"beUsedTrexKingNum", "凶凶龙登场次数"},
	{"skillCard_1", "冥王身份卡使用次数"},
	{"skillCard_2", "海王身份卡使用次数"},
	{"skillCard_3", "神王身份卡使用次数"},
	{"skillCard_4", "军师身份卡使用次数"},
	{"skillCard_5", "武圣身份卡使用次数"},
	{"skillCard_6", "枭雄身份卡使用次数"},
	{"skillCard_7", "影武者身份卡使用次数"},
	{"skillCard_8", "甜心身份卡使用次数"},
	{"skillCard_9", "球球身份卡使用次数"},
	{"skillCard_10", "飞翼身份卡使用次数"},
	{"skillCard_11", "迪迦身份卡使用次数"},
	{"skillCard_12", "泽塔身份卡使用次数"},
	{"isFiring", "是否正在开火"},
	{"state", "玩家当前状态"},
	{"isCoinPicked", "复活币是否被拾取"},
	{"posX", "选手位置x坐标"},
	{"posY", "选手位置y坐标"},
	{"posZ", "选手位置z坐标"},
	{"health", "血量"},
	{"liveState", "存活状态"},
	{"curWeapon", "当前武器"},
	{"isOutCircle", "是否在毒圈外"},
	{"bountyCoin", "赏金"},
	{"warId", "战场ID"},
}

func buildPlayerFieldIndex() map[string]string {
	m := make(map[string]string, len(playerFieldMap))
	for _, f := range playerFieldMap {
		m[f.Key] = f.Label
	}
	return m
}

func getPlayerKeysOrdered() []string {
	keys := make([]string, len(playerFieldMap))
	for i, f := range playerFieldMap {
		keys[i] = f.Key
	}
	return keys
}

func getFieldLabel(key string) string {
	if label, ok := buildPlayerFieldIndex()[key]; ok {
		return label
	}
	return key
}

// ============================================================
// 聚合规则
// ============================================================

var sumKeys = map[string]bool{
	"totalKill": true, "damage": true, "inDamage": true, "heal": true,
	"headShot": true, "vehicleKill": true, "airDrop": true,
	"bulletTotalNum": true, "hitBulletNum": true, "hitHeadBulletNum": true,
	"moveDistance": true, "driveDistance": true, "marchDistance": true,
	"assists": true, "rescue": true,
	"hook": true, "big": true, "fire": true, "fireDamage": true, "fireKill": true,
	"cannon": true, "wormHole": true, "healBomb": true, "healBombHeal": true,
	"shelter": true, "cloudBomb": true, "bomb": true, "bombDamage": true,
	"bombKill": true, "bombMax": true, "bombMaxDamage": true, "bombMaxKill": true,
	"pickSwordNum": true, "swordKill": true, "swordDamage": true,
	"mechaDamage": true, "mechaKill": true, "mechaDistance": true,
	"dragonDamage": true, "dragonKill": true, "dragonDistance": true,
	"xcc": true, "xccDistance": true, "hitWeak": true, "useResurrectionNum": true,
	"allGunDamage": true,
	"trexKingDistance": true, "peterosaurDistance": true,
	"triceratopDistance": true, "raptorsDistance": true,
	"trexKingDamage": true, "trexKingKill": true,
	"beUsedDragonNum": true, "beUsedTrexKingNum": true,
	"skillCard_1": true, "skillCard_2": true, "skillCard_3": true, "skillCard_4": true,
	"skillCard_5": true, "skillCard_6": true, "skillCard_7": true, "skillCard_8": true,
	"skillCard_9": true, "skillCard_10": true, "skillCard_11": true, "skillCard_12": true,
	"bountyCoin": true,
}

var maxKeys = map[string]bool{
	"maxKillDistance": true, "maxHitDownDistance": true,
	"helmetLevel": true, "vestLevel": true, "bagLevel": true,
	"survive": true,
}

var minKeys = map[string]bool{
	"lostRoleRank": true,
}

// ============================================================
// 主线：生成 Excel 并返回 Buffer
// ============================================================

// GenerateWarTableExcel 调用代理获取战场数据并生成 Excel 文件
func (s *WarTableExcelService) GenerateWarTableExcel(warID string) (*bytes.Buffer, string, error) {
	// 1. 从代理获取原始数据
	rawData, err := competitionTeamServiceInstance.fetchWarInfoRaw(warID)
	if err != nil {
		return nil, "", fmt.Errorf("获取战场数据失败: %w", err)
	}

	// 2. 解析玩家列表
	playerListRaw, ok := rawData["playerInfoList"].([]interface{})
	if !ok || len(playerListRaw) == 0 {
		return nil, "", fmt.Errorf("战场ID [%s] 暂无选手数据", warID)
	}

	playerList := make([]map[string]interface{}, 0, len(playerListRaw))
	for _, p := range playerListRaw {
		if playerMap, ok := p.(map[string]interface{}); ok {
			playerList = append(playerList, playerMap)
		}
	}
	if len(playerList) == 0 {
		return nil, "", fmt.Errorf("战场ID [%s] 选手数据解析失败", warID)
	}

	// 3. 获取所有战队标识，用于匹配
	teamCodes := getAllTeamCodes()

	// 4. 按战队匹配分组（与 warScores 使用相同逻辑）
	teamGroups := matchAndGroupByTeam(playerList, teamCodes)

	// 5. 创建 Excel 工作簿
	f := excelize.NewFile()

	// Sheet 1: 战场基本信息
	writeBattlefieldInfoSheet(f, rawData, warID, len(teamGroups), len(playerList))

	// Sheet 2: 选手明细数据
	writePlayerDetailSheet(f, playerList)

	// Sheet 3: 队伍明细数据
	writeTeamDetailSheet(f, teamGroups)

	// 删除默认 Sheet1
	f.DeleteSheet("Sheet1")

	// 6. 写入 buffer
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, "", fmt.Errorf("生成Excel文件失败: %w", err)
	}

	filename := fmt.Sprintf("战场数据_%s_%s.xlsx", warID, time.Now().Format("20060102_150405"))
	return &buf, filename, nil
}

// getAllTeamCodes 从数据库中获取所有战队标识（team_code）
func getAllTeamCodes() []string {
	var teams []example.CompetitionTeam
	if err := global.GVA_DB.Find(&teams).Error; err != nil {
		return nil
	}
	codes := make([]string, len(teams))
	for i, t := range teams {
		codes[i] = t.TeamCode
	}
	return codes
}

// ============================================================
// Sheet 1: 战场基本信息
// ============================================================

func writeBattlefieldInfoSheet(f *excelize.File, rawData map[string]interface{}, warID string, teamCount, playerCount int) {
	f.SetSheetName("Sheet1", "战场基本信息")

	battleInfoFields := []struct {
		Key   string
		Label string
	}{
		{"warId", "战场ID"},
		{"customRoomName", "自定义房间名"},
		{"startTime", "比赛开始时间"},
		{"playerCount", "玩家数量"},
		{"teamCount", "队伍数量"},
		{"circleLevel", "当前圈层"},
		{"circlePosX", "当前安全区x坐标"},
		{"circlePosY", "当前安全区y坐标"},
		{"circleRadius", "当前安全区半径"},
		{"airlineStartPosX", "航线起点X"},
		{"airlineStartPosY", "航线起点Y"},
		{"airlineEndPosX", "航线终点X"},
		{"airlineEndPosY", "航线终点Y"},
	}

	f.SetCellValue("战场基本信息", "A1", "键")
	f.SetCellValue("战场基本信息", "B1", "值")

	row := 2
	for _, field := range battleInfoFields {
		f.SetCellValue("战场基本信息", fmt.Sprintf("A%d", row), field.Label)
		var value string
		if field.Key == "warId" {
			value = warID
		} else if field.Key == "playerCount" {
			value = strconv.Itoa(playerCount)
		} else if field.Key == "teamCount" {
			value = strconv.Itoa(teamCount)
		} else if v, ok := rawData[field.Key]; ok {
			value = fmt.Sprintf("%v", v)
		}
		f.SetCellValue("战场基本信息", fmt.Sprintf("B%d", row), value)
		row++
	}

	f.SetCellValue("战场基本信息", fmt.Sprintf("A%d", row), "数据导出时间")
	f.SetCellValue("战场基本信息", fmt.Sprintf("B%d", row), time.Now().Format("2006-01-02 15:04:05"))

	f.SetColWidth("战场基本信息", "A", "A", 20)
	f.SetColWidth("战场基本信息", "B", "B", 40)
}

// ============================================================
// Sheet 2: 选手明细数据
// ============================================================

func writePlayerDetailSheet(f *excelize.File, playerList []map[string]interface{}) {
	index, err := f.NewSheet("选手明细数据")
	if err != nil {
		return
	}
	f.SetActiveSheet(index)

	orderedKeys := getPlayerKeysOrdered()
	allKeys := collectAllKeys(playerList, orderedKeys)

	// 标题行
	for colIdx, key := range allKeys {
		cell, _ := excelize.CoordinatesToCellName(colIdx+1, 1)
		f.SetCellValue("选手明细数据", cell, getFieldLabel(key))
	}

	// 数据行
	for rowIdx, player := range playerList {
		for colIdx, key := range allKeys {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, rowIdx+2)
			if val := player[key]; val != nil {
				f.SetCellValue("选手明细数据", cell, val)
			}
		}
	}
}

// ============================================================
// Sheet 3: 队伍明细数据
// ============================================================

type teamAggData struct {
	TeamCode    string
	PlayerCount int
	Stats       map[string]float64
}

func writeTeamDetailSheet(f *excelize.File, teamGroups []teamAggData) {
	index, err := f.NewSheet("队伍明细数据")
	if err != nil {
		return
	}
	f.SetActiveSheet(index)

	orderedKeys := getPlayerKeysOrdered()
	allKeys := make([]string, 0)
	for _, k := range orderedKeys {
		if k == "nickName" || k == "uID" || k == "teamID" || k == "teamName" {
			continue // 跳过个人标识字段
		}
		allKeys = append(allKeys, k)
	}

	// 按队伍排名排序
	sort.Slice(teamGroups, func(i, j int) bool {
		ri := teamGroups[i].Stats["lostRoleRank"]
		rj := teamGroups[j].Stats["lostRoleRank"]
		if ri != rj {
			return ri < rj
		}
		return teamGroups[i].PlayerCount > teamGroups[j].PlayerCount
	})

	// 标题行
	f.SetCellValue("队伍明细数据", "A1", "队伍名")
	f.SetCellValue("队伍明细数据", "B1", "队伍ID")
	f.SetCellValue("队伍明细数据", "C1", "选手数量")
	for colIdx, key := range allKeys {
		cell, _ := excelize.CoordinatesToCellName(colIdx+4, 1)
		f.SetCellValue("队伍明细数据", cell, getFieldLabel(key))
	}

	// 数据行
	for rowIdx, team := range teamGroups {
		rowNum := rowIdx + 2
		f.SetCellValue("队伍明细数据", fmt.Sprintf("A%d", rowNum), team.TeamCode)
		f.SetCellValue("队伍明细数据", fmt.Sprintf("B%d", rowNum), team.TeamCode)
		f.SetCellValue("队伍明细数据", fmt.Sprintf("C%d", rowNum), team.PlayerCount)

		for colIdx, key := range allKeys {
			cell, _ := excelize.CoordinatesToCellName(colIdx+4, rowNum)
			if val, ok := team.Stats[key]; ok {
				f.SetCellValue("队伍明细数据", cell, val)
			}
		}
	}
}

// ============================================================
// 战队匹配与聚合（复用 warScores 的 matchPlayerToTeam 逻辑）
// ============================================================

// matchAndGroupByTeam 根据昵称前缀匹配战队，然后按战队聚合选手数据
func matchAndGroupByTeam(playerList []map[string]interface{}, teamCodes []string) []teamAggData {
	teamMap := make(map[string]*teamAggData)
	teamOrder := make([]string, 0)

	for _, player := range playerList {
		nickName, _ := player["nickName"].(string)
		if nickName == "" {
			continue
		}

		// 使用与 warScores 完全相同的匹配逻辑
		matchedCode := matchPlayerToTeam(nickName, teamCodes)
		if matchedCode == "" {
			continue
		}

		agg, ok := teamMap[matchedCode]
		if !ok {
			agg = &teamAggData{
				TeamCode: matchedCode,
				Stats:    make(map[string]float64),
			}
			teamMap[matchedCode] = agg
			teamOrder = append(teamOrder, matchedCode)
		}

		agg.PlayerCount++

		// 聚合字段
		for key, val := range player {
			numVal := toFloat64(val)
			if numVal == 0 && !isNumericKey(key) {
				continue
			}

			if maxKeys[key] {
				if numVal > agg.Stats[key] {
					agg.Stats[key] = numVal
				}
			} else if minKeys[key] {
				if curr, exists := agg.Stats[key]; !exists || numVal < curr {
					agg.Stats[key] = numVal
				}
			} else if sumKeys[key] {
				agg.Stats[key] += numVal
			}
			// 不在预设聚合规则中的字段跳过（非数值字段）
		}
	}

	result := make([]teamAggData, 0, len(teamOrder))
	for _, tid := range teamOrder {
		result = append(result, *teamMap[tid])
	}
	return result
}

// ============================================================
// 辅助函数
// ============================================================

func collectAllKeys(playerList []map[string]interface{}, orderedKeys []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0)

	// 先按预设顺序
	for _, key := range orderedKeys {
		for _, player := range playerList {
			if _, ok := player[key]; ok {
				result = append(result, key)
				seen[key] = true
				break
			}
		}
	}

	// 额外字段
	for _, player := range playerList {
		for key := range player {
			if !seen[key] {
				seen[key] = true
				result = append(result, key)
			}
		}
	}

	return result
}

func isNumericKey(key string) bool {
	if sumKeys[key] || maxKeys[key] || minKeys[key] {
		return true
	}
	if len(key) > 10 && key[:10] == "skillCard_" {
		return true
	}
	return false
}

func toFloat64(val interface{}) float64 {
	switch v := val.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	case int64:
		return float64(v)
	case int32:
		return float64(v)
	case string:
		f, _ := strconv.ParseFloat(v, 64)
		return f
	default:
		return 0
	}
}

// competitionTeamServiceInstance 用于调用 fetchWarInfoRaw
var competitionTeamServiceInstance = &CompetitionTeamService{}
