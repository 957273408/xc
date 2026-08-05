                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                              Q<template>
  <div class="player-selection-page">
    <!-- 顶部品牌栏 -->
    <header class="page-header">
      <div class="header-bg-overlay"></div>
      <div class="header-content">
        <div class="brand-left">
          <div class="logo-badge">SMC</div>
          <span class="event-year">2026</span>
        </div>
        <h1 class="page-title">选手数据选择展示系统</h1>
        <div class="brand-right">
          <span class="sub-title">SAUSAGE MAN</span>
          <span class="championship-label">CHAMPIONSHIP</span>
        </div>
      </div>
    </header>

    <main class="page-main">
      <!-- ============ 控制区 ============ -->
      <section class="control-panel">
        <div class="panel-card">
          <div class="input-row">
            <div class="warid-input-group">
              <label class="input-label">战场ID (WarId)</label>
              <div class="input-with-btn">
                <el-input
                  v-model="warIdInput"
                  placeholder="请输入战场 WarId，例如: 123456789"
                  clearable
                  size="large"
                  class="warid-input"
                  @keyup.enter="handleFetchPlayers"
                />
                <el-button
                  type="primary"
                  size="large"
                  :loading="fetchingPlayers"
                  :disabled="!warIdInput.trim()"
                  class="fetch-btn"
                  @click="handleFetchPlayers"
                >
                  {{ fetchingPlayers ? '加载中...' : '获取玩家数据' }}
                </el-button>
              </div>
            </div>

            <div class="team-select-group">
              <label class="input-label">选择战队（自动勾选该战队选手）</label>
              <el-select
                v-model="selectedTeamCode"
                placeholder="请先获取玩家数据后选择战队"
                size="large"
                class="team-select"
                filterable
                :disabled="playersList.length === 0"
                @change="handleTeamChange"
              >
                <el-option
                  v-for="t in teamList"
                  :key="t.teamCode"
                  :label="`${t.teamName} (${t.teamCode})`"
                  :value="t.teamCode"
                >
                  <div class="team-option">
                    <img v-if="t.teamLogo" :src="t.teamLogo" class="team-option-logo" alt="">
                    <span class="team-option-name">{{ t.teamName }}</span>
                    <span class="team-option-code">{{ t.teamCode }}</span>
                  </div>
                </el-option>
              </el-select>
            </div>
          </div>

          <div class="input-row">
            <div class="extra-stat-group">
              <label class="input-label">附加展示项 1（第4项）</label>
              <el-select
                v-model="selectedExtraStat1"
                placeholder="选择第4项数据"
                size="large"
                class="extra-stat-select"
              >
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
              <el-select
                v-model="selectedExtraStat2"
                placeholder="选择第5项数据"
                size="large"
                class="extra-stat-select"
              >
                <el-option label="伤害量" value="damage" />
                <el-option label="治疗量" value="healing" />
                <el-option label="移动距离" value="movement" />
                <el-option label="投掷物使用数" value="throwables" />
                <el-option label="身份卡使用数" value="identity_card" />
                <el-option label="最远击杀距离" value="longest_kill" />
              </el-select>
            </div>
          </div>

          <!-- 错误提示 -->
          <el-alert
            v-if="errorMessage"
            :title="errorMessage"
            type="error"
            show-icon
            :closable="true"
            @close="errorMessage = ''"
            class="error-alert"
          />
          <!-- 成功提示 -->
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

        <!-- 会话恢复提示 -->
        <div v-if="restoredFromSession" class="restore-tip">
          <el-tag type="info" effect="light">
            ✨ 已从历史会话恢复数据（SessionKey: {{ shortSessionKey }}）
          </el-tag>
        </div>
      </section>

      <!-- ============ 玩家列表 ============ -->
      <section v-if="playersList.length > 0" class="player-list-section">
        <div class="section-header">
          <h2 class="section-title">
            <span class="title-accent">👥</span>
            战场玩家列表
            <el-tag size="small" type="warning" class="select-count-tag">
              已选择 {{ selectedPlayers.length }} / 4
            </el-tag>
          </h2>
          <div class="section-actions">
            <el-button
              type="success"
              size="large"
              :loading="savingSelection"
              :disabled="selectedPlayers.length !== 4"
              class="save-btn"
              @click="handleSaveSelection"
            >
              {{ savingSelection ? '保存中...' : `💾 保存选择 (${selectedPlayers.length}/4)` }}
            </el-button>
            <el-button
              v-if="savedSelection && savedSelection.sessionKey"
              size="large"
              type="primary"
              plain
              @click="handleCopyShareLink"
            >
              🔗 复制分享链接
            </el-button>
          </div>
        </div>

        <div class="table-wrapper">
          <el-table
            ref="playerTableRef"
            :data="playersList"
            @selection-change="handleSelectionChange"
            stripe
            height="500"
            class="player-table"
          >
            <el-table-column type="selection" width="55" reserve-selection />
            <el-table-column prop="nickName" label="玩家昵称" min-width="180" fixed="left">
              <template #default="{ row }">
                <div class="nickname-cell">
                  <div class="avatar-placeholder">
                    {{ getInitial(row.nickName) }}
                  </div>
                  <span class="nickname-text">{{ row.nickName }}</span>
                </div>
              </template>
            </el-table-column>
            <el-table-column prop="killCount" label="淘汰数" width="100" align="center">
              <template #default="{ row }">
                <span class="stat-kill">{{ row.killCount }}</span>
              </template>
            </el-table-column>
            <el-table-column label="爆头率" width="110" align="center">
              <template #default="{ row }">
                <span class="stat-rate">{{ formatPercent(row.headshotRate) }}</span>
              </template>
            </el-table-column>
            <el-table-column label="命中率" width="110" align="center">
              <template #default="{ row }">
                <span class="stat-rate">{{ formatPercent(row.accuracyRate) }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="damageAmount" label="伤害量" width="120" align="center">
              <template #default="{ row }">
                <span class="stat-damage">{{ formatNumber(row.damageAmount) }}</span>
              </template>
            </el-table-column>
            <!-- 附加统计列 -->
            <el-table-column v-if="selectedExtraStat === 'healing'" label="治疗量" width="110" align="center">
              <template #default="{ row }">
                <span class="stat-extra">{{ row.healingAmount ?? '-' }}</span>
              </template>
            </el-table-column>
            <el-table-column v-if="selectedExtraStat === 'movement'" label="移动距离(m)" width="130" align="center">
              <template #default="{ row }">
                <span class="stat-extra">{{ row.movementDistance ?? '-' }}</span>
              </template>
            </el-table-column>
            <el-table-column v-if="selectedExtraStat === 'throwables'" label="投掷物数" width="110" align="center">
              <template #default="{ row }">
                <span class="stat-extra">{{ row.throwablesUsed ?? '-' }}</span>
              </template>
            </el-table-column>
            <el-table-column v-if="selectedExtraStat === 'identity_card'" label="身份卡数" width="110" align="center">
              <template #default="{ row }">
                <span class="stat-extra">{{ row.identityCardUsed ?? '-' }}</span>
              </template>
            </el-table-column>
            <el-table-column v-if="selectedExtraStat === 'longest_kill'" label="最远击杀(m)" width="130" align="center">
              <template #default="{ row }">
                <span class="stat-extra">{{ row.longestKillDist ?? '-' }}</span>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="100" align="center" fixed="right">
              <template #default="{ row }">
                <el-button
                  link
                  type="primary"
                  @click="toggleQuickSelect(row)"
                >
                  {{ isSelected(row) ? '取消' : '选中' }}
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </section>

      <!-- 玩家列表 Loading 骨架 -->
      <section v-else-if="fetchingPlayers" class="player-list-section">
        <div class="skeleton-wrapper">
          <el-skeleton :rows="8" animated />
        </div>
      </section>

      <!-- ============ 展示卡片区 ============ -->
      <section v-if="displayCards.length > 0" class="display-section">
        <!-- 装饰背景文字 -->
        <div class="bg-decoration-text">SAUSAGE MAN</div>

        <div class="cards-container">
          <div
            v-for="(player, idx) in displayCards"
            :key="player.nickName + '_' + idx"
            class="player-card"
            :class="['card-' + (idx + 1)]"
          >
            <!-- 卡片外框装饰 -->
            <div class="card-frame">
              <div class="frame-corner tl"></div>
              <div class="frame-corner tr"></div>
              <div class="frame-corner bl"></div>
              <div class="frame-corner br"></div>
            </div>

            <!-- 头像区 -->
            <div class="card-avatar-area">
              <div class="avatar-frame">
                <div class="avatar-placeholder large">
                  {{ getInitial(player.nickName) }}
                </div>
              </div>
            </div>

            <!-- PLAYER ID标签 -->
            <div class="player-id-label">PLAYER ID</div>

            <!-- 玩家昵称 -->
            <div class="player-nickname">{{ truncate(player.nickName, 8) }}</div>

            <!-- 数据行：直接渲染接口返回的 stats 数组 -->
            <div
              v-for="(stat, sIdx) in player.stats"
              :key="stat.name"
              class="stat-row"
              :class="'stat-row-option-' + (sIdx + 1)"
            >
              <span class="stat-label">{{ stat.name }}</span>
              <span class="stat-value">{{ stat.value }}</span>
            </div>
          </div>
        </div>

        <!-- 战队名称区 -->
        <div class="team-name-area">
          <div class="team-logo-placeholder">
            <img v-if="displayTeamLogo" :src="displayTeamLogo" class="team-logo-img" alt="logo">
            <div v-else class="wolf-icon">🐺</div>
          </div>
          <div class="team-text-area">
            <div class="team-label">TEAM</div>
            <div class="team-name-display">
              {{ displayTeamName }}
            </div>
          </div>
          <div class="championship-tag-area">
            <span class="tag-sm">SAUSAGE MAN</span>
            <span class="tag-lg">CHAMPIONSHIP</span>
          </div>
        </div>
      </section>
    </main>

    <!-- 页脚 -->
    <footer class="page-footer">
      <div class="footer-left">
        <span>SAUSAGE MAN</span>
        <span class="footer-year">· SMC 2026</span>
      </div>
      <div class="footer-right">
        <span>SAUSAGE MAN · CHAMPIONSHIP</span>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch, nextTick } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  getWarPlayers,
  savePlayerSelection,
  getPlayerSelection,
  getLatestSelection
} from '@/api/playerSelection'
import { getPublicTeamList } from '@/api/competitionTeam'

// ========== 响应式状态 ==========
const route = useRoute()

const warIdInput = ref('')
const selectedExtraStat1 = ref('damage')
const selectedExtraStat2 = ref('healing')
// 兼容旧变量
const selectedExtraStat = computed(() => selectedExtraStat2.value)

const teamList = ref([])
const selectedTeamCode = ref('')

const playersList = ref([])
const selectedPlayers = ref([])
const playerTableRef = ref(null)

const fetchingPlayers = ref(false)
const savingSelection = ref(false)

const errorMessage = ref('')
const successMessage = ref('')

const savedSelection = ref(null)
const restoredFromSession = ref(false)
const teamName = ref('')
const latestData = ref(null)

// ========== 工具：附加统计映射 ==========
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

// ========== 计算属性 ==========
// 把 latestData 返回的扁平化 dataXname/dataXvalue 统一转为 stats 数组渲染
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
// 统一卡片数据源
const displayCards = computed(() => {
  if (latestData.value?.players?.length > 0) {
    return latestData.value.players.map(flattenToStats)
  }
  return selectedPlayers.value.slice(0, 4).map(p => ({
    nickName: p.nickName,
    stats: [
      { name: '淘汰数', value: String(p.killCount ?? 0) },
      { name: '爆头率', value: formatPercent(p.headshotRate) },
      { name: '命中率', value: formatPercent(p.accuracyRate) },
      { name: extraLabel(selectedExtraStat1.value), value: String(extraValue(p, selectedExtraStat1.value)) },
      { name: extraLabel(selectedExtraStat2.value), value: String(extraValue(p, selectedExtraStat2.value)) }
    ]
  }))
})

const displayTeamName = computed(() => {
  if (latestData.value?.teamName) return latestData.value.teamName
  return teamName.value || 'NAMEA'
})

const displayTeamLogo = computed(() => {
  return latestData.value?.teamLogo || ''
})

const shortSessionKey = computed(() => {
  if (!savedSelection.value?.sessionKey) return ''
  const sk = savedSelection.value.sessionKey
  return sk.length > 16 ? sk.slice(0, 8) + '...' + sk.slice(-6) : sk
})

// ========== 生命周期 ==========
onMounted(async () => {
  // 先加载战队列表
  await loadTeamList()

  // 优先处理 URL warId 参数
  const urlWarId = route.query.warId
  if (urlWarId) {
    warIdInput.value = urlWarId
    handleFetchPlayers()
  } else {
    // 直接获取最新保存的数据
    loadLatestSelection()
  }
})

// ========== 方法 ==========
async function loadTeamList() {
  try {
    const res = await getPublicTeamList()
    if (res.code === 0 && res.data?.teams) {
      teamList.value = res.data.teams
    }
  } catch (e) { /* 静默失败 */ }
}

// 分隔符（与后端 isSeparator 对齐）
const TEAM_SEP_RE = /[.．·•・_\-#\s]/
function isTeamMatched(nick, teamCode) {
  if (!nick || !teamCode) return false
  const code = teamCode.toLowerCase()
  // 情况1：战队前缀 + 分隔符 开头精确匹配
  if (nick.toLowerCase().startsWith(code)) {
    const after = nick.charAt(code.length)
    if (!after || TEAM_SEP_RE.test(after)) return true
  }
  // 情况2：分隔符包围的包含匹配，例如 xxx·TEAMCODE·xxx
  const parts = nick.split(TEAM_SEP_RE)
  return parts.some(part => part.toLowerCase() === code)
}

function handleTeamChange(teamCode) {
  if (!teamCode) return
  if (playersList.value.length === 0) {
    ElMessage.warning('请先获取战场玩家数据')
    return
  }
  const matched = playersList.value.filter(p => isTeamMatched(p.nickName || '', teamCode))
  if (matched.length === 0) {
    ElMessage.warning(`当前战场未找到战队 ${teamCode} 的选手`)
    return
  }
  if (matched.length < 4) {
    ElMessage.warning(`战队 ${teamCode} 在该战场仅找到 ${matched.length} 名选手，不足 4 人`)
  }
  // 先清空所有选择
  selectedPlayers.value.forEach(p => playerTableRef.value?.toggleRowSelection(p, false))
  selectedPlayers.value = []
  // 勾选匹配的选手，最多4个
  const toSelect = matched.slice(0, 4)
  toSelect.forEach(p => {
    playerTableRef.value?.toggleRowSelection(p, true)
  })
  selectedPlayers.value = toSelect
  // 设置战队名（用于展示卡片区）
  const t = teamList.value.find(x => x.teamCode === teamCode)
  teamName.value = t?.teamName || teamCode
  if (matched.length >= 4) {
    successMessage.value = `已自动勾选战队 ${teamCode} 的 ${toSelect.length} 名选手`
    setTimeout(() => successMessage.value = '', 3000)
  }
}

async function loadLatestSelection() {
  try {
    const res = await getLatestSelection()
    if (res.code === 0 && res.data && res.data.players && res.data.players.length > 0) {
      const data = res.data
      // 直接用返回的完整数据渲染展示卡片
      latestData.value = data
      warIdInput.value = data.warId
      selectedExtraStat1.value = data.extraStat1 || 'damage'
      selectedExtraStat2.value = data.extraStat2 || data.extraStat || 'healing'
      restoredFromSession.value = true
    }
  } catch (e) {
    // 静默失败，用户可手动输入WarId
  }
}

async function restoreBySessionKey(sessionKey) {
  try {
    const res = await getPlayerSelection({ sessionKey })
    if (res.code === 0 && res.data) {
      const data = res.data
      savedSelection.value = data
      warIdInput.value = data.warId
      selectedExtraStat1.value = data.extraStat1 || 'damage'
      selectedExtraStat2.value = data.extraStat2 || data.extraStat || 'healing'
      restoredFromSession.value = true

      await handleFetchPlayers(true)

      await nextTick()
      await selectPlayersByIds(data.selectedPlayerIds)
      teamName.value = extractTeamName()
    }
  } catch (e) {
    errorMessage.value = '恢复会话失败：' + (e.message || '网络错误')
  }
}

async function handleFetchPlayers(silent = false) {
  const warId = warIdInput.value.trim()
  if (!warId) {
    if (!silent) {
      ElMessage.warning('请输入 WarId')
    }
    return
  }

  fetchingPlayers.value = true
  errorMessage.value = ''
  if (!silent) successMessage.value = ''
  // 用户主动拉取数据时清除 latestData，改由 selectedPlayers 驱动展示
  latestData.value = null

  try {
    const res = await getWarPlayers({ warId })
    if (res.code === 0) {
      playersList.value = res.data.players || []
      // 保存lastWarId
      try {
        localStorage.setItem('playerSelectionLastWarId', warId)
      } catch (e) {}

      if (!silent) {
        if (playersList.value.length === 0) {
          successMessage.value = '获取成功，但该战场暂无玩家数据'
        } else {
          successMessage.value = `成功加载 ${playersList.value.length} 名玩家数据`
        }
      }
    } else {
      errorMessage.value = res.msg || '获取玩家数据失败'
    }
  } catch (e) {
    errorMessage.value = '获取失败：' + (e.message || '网络错误，请稍后重试')
  } finally {
    fetchingPlayers.value = false
  }
}

function handleSelectionChange(selection) {
  // 限制最多选4人
  if (selection.length > 4) {
    // 回滚：只保留前4个 + 提示
    const diff = selection.slice(4)
    diff.forEach((row) => playerTableRef.value?.toggleRowSelection(row, false))
    ElMessage.warning('最多只能选择 4 名玩家')
    selectedPlayers.value = selection.slice(0, 4)
    return
  }
  selectedPlayers.value = selection
  teamName.value = extractTeamName()
}

function toggleQuickSelect(row) {
  const isSel = isSelected(row)
  if (!isSel && selectedPlayers.value.length >= 4) {
    ElMessage.warning('最多只能选择 4 名玩家')
    return
  }
  playerTableRef.value?.toggleRowSelection(row, !isSel)
}

function isSelected(row) {
  return selectedPlayers.value.some((p) => p.playerId === row.playerId)
}

async function selectPlayersByIds(idList) {
  if (!playerTableRef.value) return
  const map = new Map(playersList.value.map((p) => [p.playerId, p]))
  idList.forEach((id) => {
    const player = map.get(id)
    if (player) {
      playerTableRef.value.toggleRowSelection(player, true)
    }
  })
}

async function handleSaveSelection() {
  if (selectedPlayers.value.length !== 4) {
    ElMessage.error(`必须选择恰好 4 名玩家，当前 ${selectedPlayers.value.length} 人`)
    return
  }

  savingSelection.value = true
  try {
    const playerIds = selectedPlayers.value.map((p) => p.playerId)
    const payload = {
      warId: warIdInput.value.trim(),
      selectedPlayerIds: playerIds,
      extraStat1: selectedExtraStat1.value,
      extraStat2: selectedExtraStat2.value,
      extraStat: selectedExtraStat2.value
    }
    // 如果已有sessionKey，带上以便更新
    if (savedSelection.value?.sessionKey) {
      payload.sessionKey = savedSelection.value.sessionKey
    }

    const res = await savePlayerSelection(payload)
    if (res.code === 0) {
      savedSelection.value = res.data
      successMessage.value = `保存成功！SessionKey: ${res.data.sessionKey.slice(0, 12)}...`

      // 存到 localStorage
      try {
        const warId = warIdInput.value.trim()
        localStorage.setItem(
          `playerSelectionSession_${warId}`,
          JSON.stringify({
            sessionKey: res.data.sessionKey,
            savedAt: Date.now()
          })
        )
        localStorage.setItem('playerSelectionLastWarId', warId)
      } catch (e) {}
    } else {
      errorMessage.value = res.msg || '保存失败'
    }
  } catch (e) {
    errorMessage.value = '保存失败：' + (e.message || '网络错误')
  } finally {
    savingSelection.value = false
  }
}

function handleCopyShareLink() {
  if (!savedSelection.value?.sessionKey) return
  const url = `${window.location.origin}${window.location.pathname}#/playerSelection?sessionKey=${savedSelection.value.sessionKey}`
  // 复制
  if (navigator.clipboard?.writeText) {
    navigator.clipboard.writeText(url).then(() => {
      ElMessage.success('分享链接已复制到剪贴板')
    }).catch(() => fallbackCopy(url))
  } else {
    fallbackCopy(url)
  }
}
function fallbackCopy(text) {
  const ta = document.createElement('textarea')
  ta.value = text
  document.body.appendChild(ta)
  ta.select()
  try {
    document.execCommand('copy')
    ElMessage.success('分享链接已复制')
  } catch (e) {
    ElMessage.info('请手动复制: ' + text)
  }
  document.body.removeChild(ta)
}

// ========== 格式化工具 ==========
function formatPercent(v) {
  if (v == null || isNaN(v)) return '0%'
  return Number(v).toFixed(1) + '%'
}
function formatNumber(v) {
  if (v == null || isNaN(v)) return '0'
  return Number(v).toLocaleString()
}
function getInitial(name) {
  if (!name) return '?'
  return name.trim().charAt(0).toUpperCase()
}
function truncate(s, n) {
  if (!s) return ''
  return s.length > n ? s.slice(0, n - 1) + '…' : s
}
function extractTeamName() {
  // 从昵称前缀推断战队名
  if (selectedPlayers.value.length === 0) return ''
  const firstNick = selectedPlayers.value[0].nickName || ''
  // 用分隔符切前缀
  const m = firstNick.match(/^([^.．·•・_\-#\s]+)[.．·•・_\-#\s]/)
  if (m) return m[1]
  return ''
}
</script>

<style lang="scss" scoped>
/* ============= 全局变量 & 基础 ============= */
$theme-red: #d92338;
$theme-red-dark: #8b0e1d;
$theme-gold: #ffd700;
$theme-bg: #0a0a0a;
$theme-bg-2: #121214;
$theme-card-bg: linear-gradient(180deg, #1a1a1d 0%, #0f0f10 100%);

.player-selection-page {
  min-height: 100vh;
  background: #000;
  color: #fff;
  font-family: 'Segoe UI', 'PingFang SC', 'Microsoft YaHei', sans-serif;
  background-image:
    radial-gradient(ellipse at 20% 0%, rgba(217, 35, 56, 0.12) 0%, transparent 50%),
    radial-gradient(ellipse at 80% 100%, rgba(217, 35, 56, 0.08) 0%, transparent 50%);
  overflow-x: hidden;
}

/* ============= 顶部品牌栏 ============= */
.page-header {
  position: relative;
  height: 80px;
  border-bottom: 3px solid $theme-red;
  overflow: hidden;

  .header-bg-overlay {
    position: absolute; inset: 0;
    background:
      repeating-linear-gradient(90deg, transparent, transparent 60px, rgba(255,255,255,0.02) 60px, rgba(255,255,255,0.02) 120px),
      linear-gradient(180deg, #1a0508 0%, #000 100%);
  }

  .header-content {
    position: relative;
    height: 100%;
    max-width: 1800px;
    margin: 0 auto;
    padding: 0 32px;
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .brand-left, .brand-right {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .logo-badge {
    width: 56px; height: 56px;
    background: $theme-red;
    border-radius: 8px;
    display: grid; place-items: center;
    font-weight: 900;
    font-size: 18px;
    letter-spacing: 1px;
    color: #fff;
    box-shadow: 0 0 20px rgba(217, 35, 56, 0.5);
    font-style: italic;
  }

  .event-year { font-size: 14px; color: rgba(255,255,255,0.5); letter-spacing: 2px; }

  .page-title {
    margin: 0;
    font-size: 24px;
    font-weight: 800;
    background: linear-gradient(90deg, #fff 0%, $theme-gold 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    letter-spacing: 4px;
  }

  .sub-title { font-size: 12px; color: rgba(255,255,255,0.5); letter-spacing: 2px; }
  .championship-label {
    font-size: 13px; color: $theme-red;
    font-weight: 700; letter-spacing: 3px;
    padding: 4px 10px;
    border: 1px solid $theme-red;
    border-radius: 4px;
  }
}

/* ============= 主内容区 ============= */
.page-main {
  max-width: 1800px;
  margin: 0 auto;
  padding: 32px;
}

/* ============= 控制面板 ============= */
.control-panel { margin-bottom: 32px; }

.panel-card {
  background: $theme-card-bg;
  border: 1px solid rgba(217, 35, 56, 0.2);
  border-radius: 16px;
  padding: 28px;
  backdrop-filter: blur(10px);
  box-shadow: 0 8px 32px rgba(0,0,0,0.6);
}

.input-row {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 28px;
  align-items: end;
}

.input-label {
  display: block;
  font-size: 13px;
  font-weight: 600;
  color: rgba(255,255,255,0.7);
  margin-bottom: 10px;
  letter-spacing: 1px;
}

.input-with-btn {
  display: flex;
  gap: 12px;

  :deep(.el-input__wrapper) {
    background: rgba(255,255,255,0.05) !important;
    border: 1px solid rgba(255,255,255,0.1) !important;
    box-shadow: none !important;
    border-radius: 10px;
    padding: 6px 16px;
  }
  :deep(.el-input__inner) { color: #fff !important; }
}

.fetch-btn {
  background: linear-gradient(135deg, $theme-red 0%, $theme-red-dark 100%);
  border: none;
  border-radius: 10px;
  padding: 0 32px;
  font-weight: 700;
  letter-spacing: 2px;
  box-shadow: 0 4px 16px rgba(217,35,56,0.4);
  transition: all .3s;
  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 6px 24px rgba(217,35,56,0.6);
  }
  &:disabled { opacity: .5; }
}

.extra-stat-select,
.team-select {
  width: 100%;
  :deep(.el-select__wrapper) {
    background: rgba(255,255,255,0.05) !important;
    border: 1px solid rgba(255,255,255,0.1) !important;
    box-shadow: none !important;
    border-radius: 10px;
  }
  :deep(.el-select__placeholder),
  :deep(.el-select__selected-item),
  :deep(.el-input__inner) { color: #fff !important; }
}

.team-option {
  display: flex; align-items: center; gap: 10px;
  .team-option-logo {
    width: 24px; height: 24px; border-radius: 6px; object-fit: cover;
    background: #fff2;
  }
  .team-option-name { font-weight: 600; color: #111; }
  .team-option-code { color: #888; font-size: 12px; margin-left: auto; }
}

.error-alert, .success-alert {
  margin-top: 20px;
  border-radius: 10px;
  border: none;
}

.restore-tip {
  margin-top: 16px;
  text-align: center;
}

/* ============= 玩家列表 ============= */
.player-list-section { margin-bottom: 48px; }

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20px;
  flex-wrap: wrap;
  gap: 16px;
}

.section-title {
  margin: 0;
  font-size: 22px;
  font-weight: 800;
  display: flex;
  align-items: center;
  gap: 12px;
  letter-spacing: 2px;

  .title-accent { font-size: 24px; }
}
.select-count-tag {
  margin-left: 8px;
  --el-tag-bg-color: rgba(255, 215, 0, 0.1);
  --el-tag-border-color: rgba(255, 215, 0, 0.5);
  --el-tag-text-color: $theme-gold;
  font-weight: 700;
}

.section-actions {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.save-btn {
  background: linear-gradient(135deg, #0ea5e9 0%, #0369a1 100%);
  border: none;
  font-weight: 700;
  border-radius: 10px;
  padding: 0 28px;
  box-shadow: 0 4px 16px rgba(14,165,233,0.3);
  transition: all .3s;
  &:hover:not(:disabled) { transform: translateY(-2px); }
  &:disabled { opacity: .45; }
}

.table-wrapper {
  background: $theme-card-bg;
  border-radius: 16px;
  border: 1px solid rgba(255,255,255,0.08);
  overflow: hidden;
}

.player-table {
  :deep(.el-table) {
    --el-table-bg-color: transparent;
    --el-table-tr-bg-color: transparent;
    --el-table-header-bg-color: rgba(255,255,255,0.03);
    --el-table-border-color: rgba(255,255,255,0.06);
    --el-table-text-color: rgba(255,255,255,0.85);
    --el-table-header-text-color: $theme-gold;
  }
  :deep(.el-table th.el-table__cell) {
    background: rgba(217,35,56,0.08) !important;
    font-weight: 700;
    letter-spacing: 1px;
  }
  :deep(.el-table tr:hover > td) {
    background: rgba(217,35,56,0.08) !important;
  }
  :deep(.el-table__row--striped td) {
    background: rgba(255,255,255,0.015);
  }
}

.nickname-cell {
  display: flex;
  align-items: center;
  gap: 10px;
}

.avatar-placeholder {
  width: 32px; height: 32px;
  background: linear-gradient(135deg, $theme-red, $theme-red-dark);
  border-radius: 8px;
  display: grid; place-items: center;
  font-weight: 800;
  font-size: 13px;
  color: #fff;
  &.large {
    width: 72px; height: 72px;
    font-size: 28px;
    border-radius: 50%;
  }
}

.stat-kill { color: #ff6b6b; font-weight: 800; font-size: 15px; }
.stat-rate { color: #38d9a9; font-weight: 700; }
.stat-damage { color: #ffd43b; font-weight: 800; }
.stat-extra { color: #74c0fc; font-weight: 700; }

.skeleton-wrapper {
  background: $theme-card-bg;
  border-radius: 16px;
  padding: 24px;
  :deep(.el-skeleton__item) {
    --el-skeleton-color: rgba(255,255,255,0.08);
    --el-skeleton-to-color: rgba(255,255,255,0.15);
  }
}

/* ============= 展示卡片区 ============= */
.display-section {
  position: relative;
  background: linear-gradient(180deg, #0d0204 0%, #000 50%, #0d0204 100%);
  border: 2px solid rgba(217,35,56,0.4);
  border-radius: 24px;
  padding: 60px 48px 48px;
  overflow: hidden;
  box-shadow:
    0 0 60px rgba(217,35,56,0.15),
    inset 0 0 120px rgba(0,0,0,0.9);
}

.bg-decoration-text {
  position: absolute;
  top: 10%;
  left: 50%;
  transform: translateX(-50%);
  font-size: clamp(100px, 16vw, 220px);
  font-weight: 900;
  letter-spacing: 20px;
  color: transparent;
  -webkit-text-stroke: 2px rgba(217,35,56,0.1);
  pointer-events: none;
  user-select: none;
  white-space: nowrap;
}

.cards-container {
  position: relative;
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 24px;
  max-width: 1600px;
  margin: 0 auto;
  margin-bottom: 56px;
}

/* -------- 单张卡片 -------- */
.player-card {
  position: relative;
  background: linear-gradient(180deg, #1a0a0c 0%, #0f0f11 60%, #000 100%);
  border: 2px solid transparent;
  border-image: linear-gradient(180deg, $theme-red 0%, rgba(255,215,0,0.3) 50%, $theme-red 100%) 1;
  padding: 18px 14px 14px;
  display: flex;
  flex-direction: column;
  align-items: center;
  transition: transform .4s cubic-bezier(.34,1.56,.64,1), box-shadow .4s;

  &:hover {
    transform: translateY(-8px) scale(1.02);
    box-shadow:
      0 20px 40px rgba(217,35,56,0.3),
      0 0 40px rgba(255,215,0,0.1);
    z-index: 2;
  }

  /* 每个卡片轻微错开高度 */
  &.card-1, &.card-5 { transform: translateY(10px); &:hover { transform: translateY(2px) scale(1.02); } }
  &.card-2, &.card-4 { transform: translateY(5px); &:hover { transform: translateY(-3px) scale(1.02); } }
  &.card-3 { transform: translateY(0); &:hover { transform: translateY(-8px) scale(1.03); } }
}

/* 四角装饰 */
.card-frame {
  position: absolute; inset: 0;
  pointer-events: none;
  .frame-corner {
    position: absolute;
    width: 18px; height: 18px;
    border: 2px solid $theme-gold;
    &.tl { top: 4px; left: 4px; border-right: none; border-bottom: none; }
    &.tr { top: 4px; right: 4px; border-left: none; border-bottom: none; }
    &.bl { bottom: 4px; left: 4px; border-right: none; border-top: none; }
    &.br { bottom: 4px; right: 4px; border-left: none; border-top: none; }
  }
}

.card-avatar-area {
  width: 100%;
  display: flex;
  justify-content: center;
  margin: 4px 0 12px;
}
.avatar-frame {
  position: relative;
  padding: 4px;
  background: linear-gradient(135deg, $theme-red, $theme-gold, $theme-red);
  border-radius: 50%;
  animation: rotate-border 4s linear infinite;
  .avatar-placeholder.large {
    position: relative;
    z-index: 1;
    background: #111;
    border: 2px solid #000;
  }
}
@keyframes rotate-border {
  0% { filter: hue-rotate(0deg); }
  100% { filter: hue-rotate(360deg); }
}

.player-id-label {
  font-size: 11px;
  letter-spacing: 3px;
  color: $theme-gold;
  font-weight: 800;
  border-top: 1px solid rgba(255,215,0,0.3);
  border-bottom: 1px solid rgba(255,215,0,0.3);
  padding: 4px 0;
  width: 100%;
  text-align: center;
  margin-bottom: 8px;
}

.player-nickname {
  font-size: 18px;
  font-weight: 900;
  color: #fff;
  text-shadow: 0 0 10px rgba(255,255,255,0.3);
  margin-bottom: 16px;
  text-align: center;
  letter-spacing: 1px;
  max-width: 100%;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* 数据行 */
.stat-row {
  width: 100%;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 6px 10px;
  border-radius: 4px;
  margin-bottom: 4px;

  .stat-label {
    font-size: 11px;
    color: rgba(255,255,255,0.5);
    letter-spacing: 1px;
  }
  .stat-value {
    font-weight: 900;
    font-size: 19px;
    font-family: 'Courier New', monospace;
    text-shadow: 0 0 8px currentColor;

    &.kill { color: #ff4757; }
    &.rate { color: #2ed573; }
    &.damage { color: #ffa502; }
    &.extra { color: #1e90ff; }
  }

  &.stat-row-option-1 { background: rgba(255,71,87,0.12); }
  &.stat-row-option-2 { background: rgba(46,213,115,0.10); }
  &.stat-row-option-3 { background: rgba(46,213,115,0.06); }
  &.stat-row-option-4 { background: rgba(255,165,2,0.10); }
  &.stat-row-extra { background: rgba(30,144,255,0.10); }
}

/* -------- 战队名称区 -------- */
.team-name-area {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 40px;
  padding: 28px 40px;
  border-top: 3px solid rgba(217,35,56,0.5);
  max-width: 1200px;
  margin: 0 auto;
}

.team-logo-placeholder {
  width: 90px; height: 90px;
  background: $theme-gold;
  border-radius: 20px;
  display: grid; place-items: center;
  box-shadow: 0 0 40px rgba(255,215,0,0.4);
  overflow: hidden;
  .wolf-icon { font-size: 48px; }
  .team-logo-img {
    width: 100%; height: 100%;
    object-fit: cover;
    display: block;
  }
}

.team-text-area { display: flex; align-items: baseline; gap: 20px; }
.team-label {
  font-size: 24px;
  font-weight: 900;
  letter-spacing: 4px;
  color: rgba(255,255,255,0.7);
}
.team-name-display {
  font-size: 88px;
  font-weight: 900;
  letter-spacing: 12px;
  color: $theme-red;
  line-height: 1;
  text-shadow:
    0 0 40px rgba(217,35,56,0.5),
    4px 4px 0 $theme-red-dark;
  font-style: italic;
}

.championship-tag-area {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 6px;
  .tag-sm { font-size: 11px; letter-spacing: 2px; color: rgba(255,255,255,0.4); }
  .tag-lg {
    font-size: 14px;
    font-weight: 800;
    letter-spacing: 3px;
    color: $theme-gold;
    padding: 4px 12px;
    border: 1px solid $theme-gold;
  }
}

/* ============= 页脚 ============= */
.page-footer {
  border-top: 1px solid rgba(255,255,255,0.08);
  padding: 20px 48px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 12px;
  letter-spacing: 2px;
  color: rgba(255,255,255,0.35);
  background: rgba(0,0,0,0.5);

  .footer-year { color: rgba(217,35,56,0.6); }
}

/* ============= 响应式 ============= */
@media (max-width: 1400px) {
  .cards-container { grid-template-columns: repeat(4, minmax(0, 1fr)); gap: 16px; }
  .stat-row .stat-value { font-size: 17px; }
  .team-name-display { font-size: 64px; letter-spacing: 8px; }
}

@media (max-width: 1200px) {
  .page-title { font-size: 18px; letter-spacing: 2px; }
  .input-row { grid-template-columns: 1fr; gap: 20px; }

  .cards-container {
    grid-template-columns: repeat(3, 1fr);
    row-gap: 32px;
  }
  .player-card { &.card-1,&.card-2,&.card-3,&.card-4,&.card-5 { transform: none; } }

  .team-name-area {
    flex-direction: column;
    gap: 16px;
    text-align: center;
    .championship-tag-area { align-items: center; }
  }
}

@media (max-width: 768px) {
  .page-main { padding: 16px; }
  .panel-card { padding: 20px; }
  .page-header { height: auto; padding: 16px 0;
    .header-content { flex-direction: column; gap: 12px; padding: 12px 16px; }
  }
  .section-actions { width: 100%;
    .el-button { flex: 1; }
  }

  .display-section { padding: 32px 16px 24px; }

  .cards-container {
    grid-template-columns: repeat(2, 1fr);
    gap: 14px;
    margin-bottom: 32px;
  }
  .player-nickname { font-size: 15px; }
  .stat-row { padding: 4px 8px;
    .stat-label { font-size: 10px; }
    .stat-value { font-size: 15px; }
  }
  .player-id-label { font-size: 9px; letter-spacing: 2px; }

  .team-logo-placeholder { width: 64px; height: 64px; border-radius: 14px;
    .wolf-icon { font-size: 32px; }
  }
  .team-name-area { gap: 12px; }
  .team-label { font-size: 16px; }
  .team-name-display { font-size: 42px; letter-spacing: 4px; }

  .page-footer {
    flex-direction: column;
    gap: 8px;
    padding: 16px;
    text-align: center;
  }
}

@media (max-width: 480px) {
  .cards-container {
    grid-template-columns: 1fr;
    max-width: 320px;
    margin-left: auto;
    margin-right: auto;
  }
  .team-name-display { font-size: 36px; letter-spacing: 3px; }
}
</style>
