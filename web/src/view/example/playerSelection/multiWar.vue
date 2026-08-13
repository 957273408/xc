<template>
  <div class="multi-war-page">
    <header class="page-header">
      <div class="header-content">
        <div class="brand-left">
          <div class="logo-badge">SMC</div>
          <span class="event-year">2026</span>
        </div>
        <h1 class="page-title">多场选手数据汇总</h1>
        <div class="brand-right">
          <span class="sub-title">SAUSAGE MAN</span>
          <span class="championship-label">CHAMPIONSHIP</span>
        </div>
      </div>
    </header>

    <main class="page-main">
      <!-- 控制区 -->
      <section class="control-panel">
        <div class="panel-card">
          <div class="input-row">
            <div class="warids-input-group">
              <label class="input-label">战场ID列表（每行一个，或用逗号/空格分隔）</label>
              <el-input
                v-model="warIdsText"
                type="textarea"
                :rows="4"
                placeholder="请输入多个 WarId，例如:&#10;123456789&#10;987654321&#10;或: 123, 456, 789"
                class="warids-textarea"
              />
            </div>
          </div>

          <div class="input-row">
            <div class="extra-stat-group">
              <label class="input-label">附加展示项 1（第4项）</label>
              <el-select v-model="selectedExtraStat1" size="large" class="extra-stat-select">
                <el-option label="伤害量" value="damage" />
                <el-option label="治疗量" value="healing" />
                <el-option label="移动距离" value="movement" />
                <el-option label="投掷物使用数" value="throwables" />
                <el-option label="身份卡使用数" value="identity_card" />
                <el-option label="最远击杀距离" value="longest_kill" />
              </el-select>
            </div>
            <div class="extra-stat-group">
              <label class="input-label">附加展示项 2（第5项）</label>
              <el-select v-model="selectedExtraStat2" size="large" class="extra-stat-select">
                <el-option label="伤害量" value="damage" />
                <el-option label="治疗量" value="healing" />
                <el-option label="移动距离" value="movement" />
                <el-option label="投掷物使用数" value="throwables" />
                <el-option label="身份卡使用数" value="identity_card" />
                <el-option label="最远击杀距离" value="longest_kill" />
              </el-select>
            </div>
            <div class="action-group">
              <el-button
                type="primary"
                size="large"
                :loading="fetching"
                :disabled="!warIdsText.trim()"
                class="fetch-btn"
                @click="handleFetch"
              >
                {{ fetching ? '汇总中...' : '获取汇总数据' }}
              </el-button>
              <el-button size="large" @click="handleClear">清空</el-button>
            </div>
          </div>

          <el-alert
            v-if="errorMessage"
            :title="errorMessage"
            type="error"
            show-icon
            :closable="true"
            @close="errorMessage = ''"
            class="error-alert"
          />
          <el-alert
            v-if="successMessage"
            :title="successMessage"
            type="success"
            show-icon
            :closable="true"
            @close="successMessage = ''"
            class="success-alert"
          />
        </div>
      </section>

      <!-- 汇总信息条 + 保存操作 -->
      <section v-if="playersList.length > 0" class="info-bar">
        <el-tag type="info" effect="dark">战场数: {{ matchCount }}</el-tag>
        <el-tag type="success" effect="dark">选手数: {{ playersList.length }}</el-tag>
        <el-tag type="warning" effect="dark">总淘汰: {{ totalKills }}</el-tag>
        <el-tag type="danger" effect="dark">总伤害: {{ formatNumber(totalDamage) }}</el-tag>
        <div class="info-bar-actions">
          <el-tag type="warning" effect="light">
            已选择 {{ selectedPlayers.length }} / 5
          </el-tag>
          <el-button
            type="success"
            size="default"
            :loading="saving"
            :disabled="selectedPlayers.length < 4 || selectedPlayers.length > 5"
            @click="handleSave"
          >
            {{ saving ? '保存中...' : `💾 保存选择 (${selectedPlayers.length}/5)` }}
          </el-button>
        </div>
      </section>

      <!-- 数据表格 -->
      <section v-if="playersList.length > 0" class="player-list-section">
        <div class="section-header">
          <h2 class="section-title">
            <span class="title-accent">📊</span>
            汇总选手数据（按淘汰数降序）
          </h2>
        </div>
        <div class="table-wrapper">
          <el-table
            ref="playerTableRef"
            :data="playersList"
            stripe
            height="560"
            class="player-table"
            @selection-change="handleSelectionChange"
          >
            <el-table-column type="selection" width="55" reserve-selection />
            <el-table-column type="index" label="#" width="60" align="center" />
            <el-table-column prop="nickName" label="玩家昵称" min-width="180">
              <template #default="{ row }">
                <div class="nickname-cell">
                  <div class="avatar-placeholder">{{ getInitial(row.nickName) }}</div>
                  <span class="nickname-text">{{ row.nickName }}</span>
                </div>
              </template>
            </el-table-column>
            <el-table-column prop="killCount" label="淘汰数" width="100" align="center" sortable>
              <template #default="{ row }"><span class="stat-kill">{{ row.killCount }}</span></template>
            </el-table-column>
            <el-table-column label="暴击率" width="110" align="center" sortable :sort-method="(a,b)=>a.headshotRate-b.headshotRate">
              <template #default="{ row }"><span class="stat-rate">{{ formatPercent(row.headshotRate) }}</span></template>
            </el-table-column>
            <el-table-column label="命中率" width="110" align="center" sortable :sort-method="(a,b)=>a.accuracyRate-b.accuracyRate">
              <template #default="{ row }"><span class="stat-rate">{{ formatPercent(row.accuracyRate) }}</span></template>
            </el-table-column>
            <el-table-column prop="damageAmount" label="伤害量" width="120" align="center" sortable>
              <template #default="{ row }"><span class="stat-damage">{{ formatNumber(row.damageAmount) }}</span></template>
            </el-table-column>
            <el-table-column label="治疗量" width="100" align="center" v-if="selectedExtraStat1 === 'healing' || selectedExtraStat2 === 'healing'">
              <template #default="{ row }"><span class="stat-extra">{{ row.healingAmount ?? '-' }}</span></template>
            </el-table-column>
            <el-table-column label="移动距离(m)" width="130" align="center" v-if="selectedExtraStat1 === 'movement' || selectedExtraStat2 === 'movement'">
              <template #default="{ row }"><span class="stat-extra">{{ row.movementDistance ?? '-' }}</span></template>
            </el-table-column>
            <el-table-column label="投掷物数" width="110" align="center" v-if="selectedExtraStat1 === 'throwables' || selectedExtraStat2 === 'throwables'">
              <template #default="{ row }"><span class="stat-extra">{{ row.throwablesUsed ?? '-' }}</span></template>
            </el-table-column>
            <el-table-column label="身份卡数" width="110" align="center" v-if="selectedExtraStat1 === 'identity_card' || selectedExtraStat2 === 'identity_card'">
              <template #default="{ row }"><span class="stat-extra">{{ row.identityCardUsed ?? '-' }}</span></template>
            </el-table-column>
            <el-table-column label="最远击杀(m)" width="130" align="center" v-if="selectedExtraStat1 === 'longest_kill' || selectedExtraStat2 === 'longest_kill'">
              <template #default="{ row }"><span class="stat-extra">{{ row.longestKillDist ?? '-' }}</span></template>
            </el-table-column>
          </el-table>
        </div>
      </section>

      <!-- Loading 骨架 -->
      <section v-else-if="fetching" class="player-list-section">
        <div class="skeleton-wrapper"><el-skeleton :rows="10" animated /></div>
      </section>

      <!-- 展示卡片区：选中的选手（或 Top 5） -->
      <section v-if="displayPlayers.length > 0" class="display-section">
        <div class="bg-decoration-text">SAUSAGE MAN</div>
        <div class="cards-container">
          <div
            v-for="(player, idx) in displayPlayers"
            :key="player.playerId + '_' + idx"
            class="player-card"
            :class="['card-' + (idx + 1)]"
          >
            <div class="card-frame">
              <div class="frame-corner tl"></div>
              <div class="frame-corner tr"></div>
              <div class="frame-corner bl"></div>
              <div class="frame-corner br"></div>
            </div>
            <div class="card-avatar-area">
              <div class="avatar-frame">
                <div class="avatar-placeholder large">{{ getInitial(player.nickName) }}</div>
              </div>
            </div>
            <div class="player-id-label" v-if="useLatestData">PLAYER</div>
            <div class="player-id-label" v-else>RANK #{{ idx + 1 }}</div>
            <div class="player-nickname">{{ truncate(player.nickName, 8) }}</div>
            <template v-if="!useLatestData">
              <div class="stat-row stat-row-option-1">
                <span class="stat-label">淘汰数</span>
                <span class="stat-value kill">{{ player.killCount }}</span>
              </div>
              <div class="stat-row stat-row-option-2">
                <span class="stat-label">暴击率</span>
                <span class="stat-value rate">{{ formatPercent(player.headshotRate) }}</span>
              </div>
              <div class="stat-row stat-row-option-3">
                <span class="stat-label">命中率</span>
                <span class="stat-value rate">{{ formatPercent(player.accuracyRate) }}</span>
              </div>
              <div class="stat-row stat-row-option-4">
                <span class="stat-label">{{ extraLabel(selectedExtraStat1) }}</span>
                <span class="stat-value damage">{{ extraValue(player, selectedExtraStat1) }}</span>
              </div>
              <div class="stat-row stat-row-extra">
                <span class="stat-label">{{ extraLabel(selectedExtraStat2) }}</span>
                <span class="stat-value extra">{{ extraValue(player, selectedExtraStat2) }}</span>
              </div>
            </template>
            <template v-else>
              <div v-for="(stat, sIdx) in player.stats" :key="sIdx" class="stat-row" :class="'stat-row-option-' + (sIdx + 1)">
                <span class="stat-label">{{ stat.name }}</span>
                <span class="stat-value" :class="statClass(stat.name)">{{ stat.value }}</span>
              </div>
            </template>
          </div>
        </div>
      </section>
    </main>

    <footer class="page-footer">
      <div class="footer-left">
        <span>SAUSAGE MAN</span>
        <span class="footer-year">· SMC 2026</span>
      </div>
      <div class="footer-right"><span>多场数据汇总</span></div>
    </footer>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import { getMultiWarPlayers, savePlayerSelection, getLatestSelection } from '@/api/playerSelection'

const warIdsText = ref('')
const selectedExtraStat1 = ref('damage')
const selectedExtraStat2 = ref('healing')

const playersList = ref([])
const matchCount = ref(0)
const fetching = ref(false)
const saving = ref(false)
const errorMessage = ref('')
const successMessage = ref('')

const selectedPlayers = ref([])
const playerTableRef = ref(null)
const latestData = ref(null)

const EXTRA_STAT_OPTIONS = [
  { value: 'damage', label: '伤害量' },
  { value: 'healing', label: '治疗量' },
  { value: 'movement', label: '移动距离' },
  { value: 'throwables', label: '投掷物数' },
  { value: 'identity_card', label: '身份卡数' },
  { value: 'longest_kill', label: '最远击杀' }
]
function extraLabel(key) {
  const f = EXTRA_STAT_OPTIONS.find(o => o.value === key)
  return f ? f.label : '附加项'
}
function extraValue(p, key) {
  switch (key) {
    case 'damage': return formatNumber(p.damageAmount)
    case 'healing': return p.healingAmount ?? '-'
    case 'movement': return p.movementDistance ?? '-'
    case 'throwables': return p.throwablesUsed ?? '-'
    case 'identity_card': return p.identityCardUsed ?? '-'
    case 'longest_kill': return p.longestKillDist ?? '-'
    default: return '-'
  }
}
function statClass(name) {
  if (name?.includes('淘汰')) return 'kill'
  if (name?.includes('爆头') || name?.includes('命中')) return 'rate'
  if (name?.includes('伤害')) return 'damage'
  return 'extra'
}

const useLatestData = computed(() => !!latestData.value?.players?.length)

function flattenToStats(p) {
  const stats = []
  for (let i = 1; i <= 5; i++) {
    const name = p[`data${i}name`]
    const value = p[`data${i}value`]
    if (name != null && name !== '') {
      stats.push({ name, value: value ?? '-' })
    }
  }
  return { nickName: p.nickName, stats }
}

const displayPlayers = computed(() => {
  if (useLatestData.value) {
    return latestData.value.players.map(flattenToStats)
  }
  return selectedPlayers.value.length > 0
    ? selectedPlayers.value
    : playersList.value.slice(0, 5)
})

const totalKills = computed(() => playersList.value.reduce((s, p) => s + (p.killCount || 0), 0))
const totalDamage = computed(() => playersList.value.reduce((s, p) => s + (p.damageAmount || 0), 0))

function parseWarIds(text) {
  return text
    .split(/[\s,，、\n\r]+/)
    .map(s => s.trim())
    .filter(s => s.length > 0)
}

async function handleFetch() {
  errorMessage.value = ''
  successMessage.value = ''
  const ids = parseWarIds(warIdsText.value)
  if (ids.length === 0) {
    ElMessage.warning('请输入至少一个 WarId')
    return
  }
  fetching.value = true
  playersList.value = []
  selectedPlayers.value = []
  latestData.value = null
  try {
    const res = await getMultiWarPlayers({ warIds: ids })
    if (res.code === 0 && res.data) {
      playersList.value = res.data.players || []
      matchCount.value = res.data.matchCount || 0
      successMessage.value = `汇总完成：${matchCount.value} 场战场，${playersList.value.length} 名选手`
    } else {
      errorMessage.value = res.msg || '获取失败'
    }
  } catch (e) {
    errorMessage.value = '请求失败：' + (e.message || '网络错误')
  } finally {
    fetching.value = false
  }
}

function handleClear() {
  warIdsText.value = ''
  playersList.value = []
  matchCount.value = 0
  selectedPlayers.value = []
  latestData.value = null
  errorMessage.value = ''
  successMessage.value = ''
}

function handleSelectionChange(selection) {
  if (selection.length > 5) {
    const diff = selection.slice(5)
    diff.forEach(row => playerTableRef.value?.toggleRowSelection(row, false))
    ElMessage.warning('最多只能选择 5 名玩家')
    selectedPlayers.value = selection.slice(0, 5)
    return
  }
  selectedPlayers.value = selection
}

async function handleSave() {
  if (selectedPlayers.value.length < 4 || selectedPlayers.value.length > 5) {
    ElMessage.error(`必须选择 4 或 5 名玩家，当前 ${selectedPlayers.value.length} 人`)
    return
  }
  saving.value = true
  try {
    const warIds = parseWarIds(warIdsText.value)
    const payload = {
      warId: warIds[0] || '',
      warIds,
      selectedPlayerIds: selectedPlayers.value.map(p => p.playerId),
      extraStat1: selectedExtraStat1.value,
      extraStat2: selectedExtraStat2.value,
      extraStat: selectedExtraStat2.value
    }
    const res = await savePlayerSelection(payload)
    if (res.code === 0) {
      successMessage.value = `保存成功！SessionKey: ${res.data.sessionKey?.slice(0, 12)}...`
      ElMessage.success('保存成功')
    } else {
      errorMessage.value = res.msg || '保存失败'
    }
  } catch (e) {
    errorMessage.value = '保存失败：' + (e.message || '网络错误')
  } finally {
    saving.value = false
  }
}

async function restoreFromLatest() {
  try {
    const res = await getLatestSelection()
    if (res.code === 0 && res.data && res.data.players && res.data.players.length > 0) {
      const data = res.data
      latestData.value = data
      selectedExtraStat1.value = data.extraStat1 || 'damage'
      selectedExtraStat2.value = data.extraStat2 || data.extraStat || 'healing'
      // 如果有 warIds，回填并重新获取汇总数据供表格展示
      if (data.warIds && data.warIds.length > 0) {
        warIdsText.value = data.warIds.join('\n')
        await handleFetch()
        // 恢复选中的选手
        await nextTick()
        if (data.selectedPlayerIds) {
          const ids = data.selectedPlayerIds
          playersList.value.forEach(p => {
            if (ids.includes(p.playerId)) {
              playerTableRef.value?.toggleRowSelection(p, true)
            }
          })
        }
        latestData.value = null // 清除 latestData，改用表格选中数据
      } else if (data.warId) {
        warIdsText.value = data.warId
        await handleFetch()
        await nextTick()
        if (data.selectedPlayerIds) {
          const ids = data.selectedPlayerIds
          playersList.value.forEach(p => {
            if (ids.includes(p.playerId)) {
              playerTableRef.value?.toggleRowSelection(p, true)
            }
          })
        }
        latestData.value = null
      }
    }
  } catch (e) {
    // 静默失败
  }
}

onMounted(() => {
  restoreFromLatest()
})

function getInitial(name) {
  if (!name) return '?'
  return name.charAt(0).toUpperCase()
}
function truncate(str, n) {
  if (!str) return ''
  return str.length > n ? str.slice(0, n) + '…' : str
}
function formatPercent(v) {
  if (v == null) return '-'
  return Number(v).toFixed(1) + '%'
}
function formatNumber(v) {
  if (v == null) return '-'
  return Number(v).toLocaleString()
}
</script>

<style lang="scss" scoped>
.multi-war-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #0a0e1a 0%, #1a1f35 50%, #0d1117 100%);
  color: #fff;
  font-family: 'Microsoft YaHei', sans-serif;
}

.page-header {
  position: relative;
  padding: 24px 40px;
  background: linear-gradient(90deg, rgba(217,35,56,0.15), rgba(10,14,26,0.6));
  border-bottom: 2px solid rgba(255,215,0,0.3);
  .header-content {
    display: flex; align-items: center; justify-content: space-between;
    max-width: 1400px; margin: 0 auto;
  }
  .brand-left, .brand-right { display: flex; align-items: center; gap: 12px; }
  .logo-badge {
    width: 44px; height: 44px; border-radius: 10px;
    background: linear-gradient(135deg, #d92338, #8b0000);
    display: flex; align-items: center; justify-content: center;
    font-weight: 900; font-size: 16px; color: #ffd700;
    box-shadow: 0 4px 12px rgba(217,35,56,0.5);
  }
  .event-year { font-size: 20px; font-weight: 700; color: #ffd700; }
  .page-title { font-size: 26px; font-weight: 800; letter-spacing: 4px; color: #fff; margin: 0; }
  .sub-title { font-size: 14px; color: #aaa; letter-spacing: 2px; }
  .championship-label { font-size: 12px; color: #ffd700; letter-spacing: 3px; font-weight: 700; }
}

.page-main {
  max-width: 1400px; margin: 0 auto; padding: 24px;
}

.control-panel { margin-bottom: 20px; }
.panel-card {
  background: rgba(255,255,255,0.04);
  border: 1px solid rgba(255,255,255,0.1);
  border-radius: 16px; padding: 24px;
}
.input-row {
  display: flex; gap: 20px; margin-bottom: 16px; flex-wrap: wrap;
  &:last-child { margin-bottom: 0; }
}
.warids-input-group { flex: 1 1 100%; }
.extra-stat-group, .action-group { flex: 1 1 240px; display: flex; flex-direction: column; gap: 8px; }
.action-group { flex-direction: row; align-items: flex-end; }
.input-label {
  display: block; font-size: 13px; color: #bbb; margin-bottom: 8px; letter-spacing: 1px;
}
.warids-textarea {
  :deep(.el-textarea__inner) {
    background: rgba(255,255,255,0.05) !important;
    border: 1px solid rgba(255,255,255,0.1) !important;
    color: #fff !important; border-radius: 10px; font-family: monospace;
  }
}
.extra-stat-select {
  width: 100%;
  :deep(.el-select__wrapper) {
    background: rgba(255,255,255,0.05) !important;
    border: 1px solid rgba(255,255,255,0.1) !important;
    box-shadow: none !important; border-radius: 10px;
  }
  :deep(.el-select__placeholder), :deep(.el-select__selected-item) { color: #fff !important; }
}
.fetch-btn {
  background: linear-gradient(135deg, #d92338, #8b0000) !important;
  border: none !important; font-weight: 700;
}
.error-alert, .success-alert { margin-top: 12px; }

.info-bar {
  display: flex; gap: 12px; margin-bottom: 16px; flex-wrap: wrap; align-items: center;
  .info-bar-actions { margin-left: auto; display: flex; gap: 12px; align-items: center; }
}

.player-list-section {
  background: rgba(255,255,255,0.03);
  border: 1px solid rgba(255,255,255,0.08);
  border-radius: 16px; padding: 20px; margin-bottom: 24px;
}
.section-header { margin-bottom: 16px; }
.section-title { font-size: 18px; font-weight: 700; color: #fff; display: flex; align-items: center; gap: 8px; }
.title-accent { font-size: 20px; }

.table-wrapper {
  :deep(.el-table) {
    background: transparent;
    --el-table-bg-color: transparent;
    --el-table-tr-bg-color: transparent;
    --el-table-header-bg-color: rgba(217,35,56,0.2);
    --el-table-border-color: rgba(255,255,255,0.08);
    --el-table-text-color: #fff;
    --el-table-header-text-color: #ffd700;
  }
  :deep(.el-table th.el-table__cell) { background: rgba(217,35,56,0.2) !important; }
  :deep(.el-table--striped .el-table__body tr.el-table__row--striped td.el-table__cell) {
    background: rgba(255,255,255,0.02) !important;
  }
}

.nickname-cell { display: flex; align-items: center; gap: 10px; }
.avatar-placeholder {
  width: 32px; height: 32px; border-radius: 50%;
  background: linear-gradient(135deg, #d92338, #8b0000);
  display: flex; align-items: center; justify-content: center;
  color: #ffd700; font-weight: 700; font-size: 14px;
  &.large { width: 64px; height: 64px; font-size: 24px; }
}
.stat-kill { color: #ff6b6b; font-weight: 700; }
.stat-rate { color: #ffd700; }
.stat-damage { color: #ff9f43; font-weight: 700; }
.stat-extra { color: #4ecdc4; }

.skeleton-wrapper { padding: 20px 0; }

.display-section {
  position: relative; padding: 32px 0;
  border-top: 1px solid rgba(255,215,0,0.2);
}
.bg-decoration-text {
  position: absolute; top: 20px; left: 50%; transform: translateX(-50%);
  font-size: 80px; font-weight: 900; color: rgba(255,215,0,0.04);
  letter-spacing: 8px; pointer-events: none; white-space: nowrap;
}
.cards-container {
  display: grid;
  grid-template-columns: repeat(5, minmax(0, 1fr));
  gap: 20px; position: relative; z-index: 1;
}
.player-card {
  position: relative;
  background: linear-gradient(160deg, rgba(217,35,56,0.15), rgba(10,14,26,0.8));
  border: 1px solid rgba(255,215,0,0.2);
  border-radius: 16px; padding: 20px 16px; text-align: center;
  overflow: hidden;
  &.card-1 { border-color: rgba(255,215,0,0.5); box-shadow: 0 0 30px rgba(255,215,0,0.2); }
}
.card-frame { position: absolute; inset: 0; pointer-events: none; }
.frame-corner {
  position: absolute; width: 18px; height: 18px;
  border: 2px solid #ffd700;
  &.tl { top: 6px; left: 6px; border-right: none; border-bottom: none; }
  &.tr { top: 6px; right: 6px; border-left: none; border-bottom: none; }
  &.bl { bottom: 6px; left: 6px; border-right: none; border-top: none; }
  &.br { bottom: 6px; right: 6px; border-left: none; border-top: none; }
}
.card-avatar-area { margin-bottom: 12px; }
.avatar-frame {
  display: inline-block; padding: 3px; border-radius: 50%;
  background: linear-gradient(135deg, #ffd700, #d92338);
}
.player-id-label {
  font-size: 11px; color: #ffd700; letter-spacing: 2px; margin-bottom: 4px;
}
.player-nickname {
  font-size: 16px; font-weight: 700; color: #fff; margin-bottom: 12px;
  letter-spacing: 1px;
}
.stat-row {
  display: flex; justify-content: space-between; align-items: center;
  padding: 5px 10px; margin-bottom: 5px; border-radius: 8px;
  background: rgba(255,255,255,0.04);
  .stat-label { font-size: 11px; color: #aaa; }
  .stat-value { font-size: 14px; font-weight: 700; }
  .stat-value.kill { color: #ff6b6b; }
  .stat-value.rate { color: #ffd700; }
  .stat-value.damage { color: #ff9f43; }
  .stat-value.extra { color: #4ecdc4; }
}

.page-footer {
  display: flex; justify-content: space-between; align-items: center;
  padding: 20px 40px; border-top: 1px solid rgba(255,255,255,0.08);
  font-size: 12px; color: #888; max-width: 1400px; margin: 0 auto;
  .footer-year { color: #ffd700; margin-left: 8px; }
}

@media (max-width: 1024px) {
  .cards-container { grid-template-columns: repeat(2, minmax(0, 1fr)); }
  .input-row { flex-direction: column; }
}
</style>
