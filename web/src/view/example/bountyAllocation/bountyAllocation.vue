<template>
  <div class="team-bonus-container">
    <!-- 页面标题 -->
    <div class="page-header">
      <h1 class="page-title">
        <el-icon><Wallet /></el-icon>
        战队奖金分配
      </h1>
      <p class="page-subtitle">为战队选手分配奖金，支持最多4名选手同时分配</p>
      <div class="type-indicator">
        <el-tag type="warning" size="small">战队奖金</el-tag>
        <span class="type-note">* 此功能用于分配战队奖金，与赏金池奖金独立管理</span>
      </div>
    </div>

    <!-- 步骤指引 -->
    <div class="steps-guide">
      <el-steps :active="currentStep" align-center>
        <el-step title="选择战队" description="选择要分配奖金的战队" />
        <el-step title="选择选手" description="从战队中选择参赛选手" />
        <el-step title="分配奖金" description="输入每位选手的分配金额" />
        <el-step title="确认提交" description="确认并提交分配方案" />
      </el-steps>
    </div>

    <!-- 步骤1：选择战队 -->
    <div v-show="currentStep === 0" class="step-content team-selection">
      <div class="section-card">
        <div class="section-header">
          <h2 class="section-title">
            <el-icon><Trophy /></el-icon>
            选择战队
          </h2>
          <el-button type="text" @click="refreshTeams" :loading="teamsLoading">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </div>

        <!-- 加载状态 -->
        <div v-if="teamsLoading" class="loading-state">
          <el-icon class="is-loading"><Loading /></el-icon>
          <span>正在加载战队数据...</span>
        </div>

        <!-- 战队列表 -->
        <div v-else class="team-grid">
          <div
            v-for="team in teams"
            :key="team.ID"
            class="team-card"
            :class="{ active: selectedTeam?.ID === team.ID }"
            @click="selectTeam(team)"
          >
            <div class="team-avatar">
              <el-icon :size="32"><Trophy /></el-icon>
            </div>
            <div class="team-info">
              <div class="team-name">{{ team.teamName }}</div>
              <div class="team-bounty">
                <span class="bounty-label">战队奖金:</span>
                <span class="bounty-value">{{ formatAmount(team.totalBounty) }}</span>
              </div>
            </div>
            <div v-if="selectedTeam?.ID === team.ID" class="selected-indicator">
              <el-icon><Check /></el-icon>
            </div>
          </div>
        </div>

        <!-- 空状态 -->
        <el-empty v-if="!teamsLoading && teams.length === 0" description="暂无可用战队" />
      </div>
    </div>

    <!-- 步骤2：选择选手 -->
    <div v-show="currentStep === 1" class="step-content player-selection">
      <div class="section-card">
        <div class="section-header">
          <h2 class="section-title">
            <el-icon><User /></el-icon>
            选择选手 - {{ selectedTeam?.teamName }}
          </h2>
          <div class="player-search">
            <el-input
              v-model="playerSearch"
              placeholder="搜索选手姓名"
              clearable
              @input="filterPlayers"
            >
              <template #prefix>
                <el-icon><Search /></el-icon>
              </template>
            </el-input>
          </div>
        </div>

        <!-- 已选选手展示 -->
        <div v-if="selectedPlayers.length > 0" class="selected-players">
          <span class="selected-label">已选选手 ({{ selectedPlayers.length }}/4):</span>
          <el-tag
            v-for="player in selectedPlayers"
            :key="player.ID"
            closable
            @close="removePlayer(player)"
            type="primary"
            class="player-tag"
          >
            {{ player.playerName }}
          </el-tag>
        </div>

        <!-- 选手列表 -->
        <div v-loading="playersLoading" class="player-list">
          <el-empty v-if="!playersLoading && filteredPlayers.length === 0" description="暂无选手" />

          <div v-else class="player-grid">
            <div
              v-for="player in filteredPlayers"
              :key="player.ID"
              class="player-card"
              :class="{
                active: isPlayerSelected(player),
                disabled: !isPlayerSelected(player) && selectedPlayers.length >= 4
              }"
              @click="togglePlayer(player)"
            >
              <div class="player-avatar">
                <el-icon :size="28"><UserFilled /></el-icon>
              </div>
              <div class="player-info">
                <div class="player-name">{{ player.playerName }}</div>
                <div class="player-bounty">
                  <span>个人赏金: {{ formatAmount(player.bounty) }}</span>
                </div>
              </div>
              <div v-if="isPlayerSelected(player)" class="selected-badge">
                <el-icon><Check /></el-icon>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 操作按钮 -->
      <div class="step-actions">
        <div class="selection-hint">
          <span>已选 {{ selectedPlayers.length }}/4 名选手</span>
          <span v-if="selectedPlayers.length >= 4" class="hint-warning">已达最大选择数</span>
        </div>
        <div class="action-buttons">
          <el-button @click="prevStep">上一步</el-button>
          <el-button type="primary" :disabled="selectedPlayers.length === 0" @click="nextStep">
            下一步
          </el-button>
        </div>
      </div>
    </div>

    <!-- 步骤3：分配奖金 -->
    <div v-show="currentStep === 2" class="step-content bonus-allocation">
      <!-- 战队奖金信息 -->
      <div class="bonus-info-bar">
        <div class="info-item">
          <span class="info-label">战队:</span>
          <span class="info-value">{{ selectedTeam?.teamName }}</span>
        </div>
        <div class="info-item">
          <span class="info-label">可用战队奖金:</span>
          <span class="info-value highlight">{{ formatAmount(teamBonus) }}</span>
        </div>
        <div class="info-item">
          <span class="info-label">已分配:</span>
          <span class="info-value">{{ formatAmount(allocatedAmount) }}</span>
        </div>
        <div class="info-item">
          <span class="info-label">剩余:</span>
          <span class="info-value" :class="{ warning: remainingAmount < 0 }">
            {{ formatAmount(remainingAmount) }}
          </span>
        </div>
        <div class="info-item">
          <span class="info-label">分配进度:</span>
          <span class="info-value">{{ allocationPercentage }}%</span>
        </div>
      </div>

      <!-- 分配面板 -->
      <div class="section-card">
        <div class="section-header">
          <h2 class="section-title">
            <el-icon><Money /></el-icon>
            奖金分配
          </h2>
          <el-button type="text" @click="autoDistribute" :disabled="teamBonus <= 0">
            <el-icon><Sort /></el-icon>
            平均分配
          </el-button>
        </div>

        <div class="allocation-grid">
          <div
            v-for="(player, index) in allocationPlayers"
            :key="player.ID"
            class="allocation-card"
          >
            <div class="allocation-header">
              <span class="player-index">{{ index + 1 }}</span>
              <span class="player-name">{{ player.playerName }}</span>
              <span class="personal-bounty">个人赏金: {{ formatAmount(player.bounty) }}</span>
            </div>

            <div class="allocation-body">
              <!-- 金额输入 -->
              <div class="input-group">
                <span class="input-label">分配金额</span>
                <div class="amount-input-wrapper" :class="{ error: player.error }">
                  <span class="currency">¥</span>
                  <el-input
                    v-model="player.amount"
                    type="number"
                    placeholder="0.00"
                    step="0.01"
                    :min="0"
                    @input="handleAmountChange(index)"
                  />
                </div>
              </div>

              <!-- 滑块（独立控制，相对于战队总奖金） -->
              <div class="slider-group">
                <span class="slider-label">分配比例: {{ player.percentage.toFixed(1) }}%</span>
                <el-slider
                  v-model="player.percentage"
                  :min="0"
                  :max="100"
                  :step="1"
                  :disabled="teamBonus <= 0"
                  @input="handleSliderChange(index)"
                />
              </div>

              <!-- 快捷金额按钮 -->
              <div class="quick-buttons">
                <el-button
                  v-for="amount in quickAmounts"
                  :key="amount"
                  size="small"
                  :disabled="amount > remainingAmount && player.amount !== amount"
                  @click="setQuickAmount(index, amount)"
                >
                  {{ formatAmount(amount) }}
                </el-button>
              </div>
            </div>

            <div v-if="player.error" class="error-message">
              <el-icon><Warning /></el-icon>
              {{ player.error }}
            </div>
          </div>
        </div>
      </div>

      <!-- 操作按钮 -->
      <div class="step-actions">
        <el-button @click="prevStep">上一步</el-button>
        <el-button type="primary" :disabled="!canSubmit" @click="nextStep">
          下一步
        </el-button>
      </div>
    </div>

    <!-- 步骤4：确认提交 -->
    <div v-show="currentStep === 3" class="step-content confirm-submit">
      <div class="section-card">
        <div class="section-header">
          <h2 class="section-title">
            <el-icon><DocumentChecked /></el-icon>
            确认分配方案
          </h2>
        </div>

        <div class="confirm-summary">
          <div class="summary-header">
            <span class="team-name">{{ selectedTeam?.teamName }}</span>
            <span class="total-amount">{{ formatAmount(allocatedAmount) }}</span>
          </div>

          <div class="summary-table">
            <div class="table-header">
              <span>选手</span>
              <span>分配金额</span>
              <span>比例</span>
            </div>
            <div
              v-for="player in allocationPlayers"
              :key="player.ID"
              class="table-row"
            >
              <span>{{ player.playerName }}</span>
              <span class="amount">{{ formatAmount(player.amount) }}</span>
              <span class="percentage">{{ player.percentage.toFixed(1) }}%</span>
            </div>
            <div class="table-footer">
              <span>合计</span>
              <span class="amount">{{ formatAmount(allocatedAmount) }}</span>
              <span class="percentage">100%</span>
            </div>
          </div>
        </div>

        <el-alert
          title="提交后将从战队奖金扣除相应金额，此操作不可撤销"
          type="warning"
          show-icon
          :closable="false"
          class="confirm-warning"
        />
      </div>

      <!-- 操作按钮 -->
      <div class="step-actions">
        <el-button @click="prevStep">上一步</el-button>
        <el-button type="warning" @click="saveDraft" :disabled="!canSave">
          <el-icon><Upload /></el-icon>
          保存草稿
        </el-button>
        <el-button type="primary" @click="showConfirmDialog = true">
          <el-icon><Check /></el-icon>
          确认提交
        </el-button>
      </div>
    </div>

    <!-- 确认弹窗 -->
    <el-dialog
      v-model="showConfirmDialog"
      title="确认提交分配"
      width="500px"
      :close-on-click-modal="false"
    >
      <div class="confirm-content">
        <div class="confirm-info">
          <div class="confirm-team">{{ selectedTeam?.teamName }}</div>
          <div class="confirm-amount">分配总额: {{ formatAmount(allocatedAmount) }}</div>
          <div class="confirm-tip">提交将包含战队全部 {{ allPlayers.length }} 名选手</div>
        </div>
        
        <div class="confirm-detail">
          <div class="detail-header">
            <span>选手</span>
            <span>分配金额</span>
          </div>
          <div v-for="player in allPlayers" :key="player.ID" class="confirm-player" :class="{ unallocated: !isPlayerSelected(player) }">
            <span class="player-name">
              {{ player.playerName }}
              <span class="personal-bounty-small">(个人赏金: {{ formatAmount(player.bounty) }})</span>
            </span>
            <span class="amount">
              {{ isPlayerSelected(player) ? formatAmount(getPlayerAmount(player.ID)) : formatAmount(0) }}
            </span>
          </div>
        </div>
        
        <div class="confirm-summary">
          <div class="summary-item">
            <span>战队总选手数:</span>
            <span>{{ allPlayers.length }} 人</span>
          </div>
          <div class="summary-item">
            <span>已选中分配:</span>
            <span class="allocated">{{ selectedPlayers.length }} 人</span>
          </div>
          <div class="summary-item">
            <span>可用战队奖金:</span>
            <span>{{ formatAmount(teamBonus) }}</span>
          </div>
          <div class="summary-item">
            <span>已分配:</span>
            <span class="allocated">{{ formatAmount(allocatedAmount) }}</span>
          </div>
          <div class="summary-item">
            <span>剩余:</span>
            <span :class="{ warning: remainingAmount < 0 }">{{ formatAmount(remainingAmount) }}</span>
          </div>
        </div>
      </div>
      <template #footer>
        <el-button @click="showConfirmDialog = false">取消</el-button>
        <el-button type="primary" @click="confirmSubmit" :loading="submitting">
          确认提交
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { ElMessage } from 'element-plus'
import {
  Wallet,
  Trophy,
  User,
  UserFilled,
  Check,
  Refresh,
  Search,
  Loading,
  Money,
  Sort,
  Warning,
  DocumentChecked,
  Upload
} from '@element-plus/icons-vue'
import { getTeamList } from '@/api/team'
import { getPlayerList, allocateBounty } from '@/api/player'
import { saveDraft as saveDraftApi } from '@/api/bountyAllocation'

defineOptions({
  name: 'TeamBonusAllocation'
})

// 快捷固定金额选项（支持2位小数）
const quickAmounts = ref([1000, 5000, 10000])

// 数据状态
const currentStep = ref(0)
const teams = ref([])
const teamsLoading = ref(false)
const selectedTeam = ref(null)
const playersLoading = ref(false)
const allPlayers = ref([])
const selectedPlayers = ref([])
const playerSearch = ref('')
const filteredPlayers = ref([])
const teamBonus = ref(0)
const allocationPlayers = ref([])
const showConfirmDialog = ref(false)
const submitting = ref(false)

// 计算属性
const allocatedAmount = computed(() => {
  return allocationPlayers.value.reduce((sum, p) => sum + Math.round(parseFloat(p.amount || 0) * 100) / 100, 0)
})

const remainingAmount = computed(() => {
  return Math.round((teamBonus.value - allocatedAmount.value) * 100) / 100
})

const allocationPercentage = computed(() => {
  if (teamBonus.value <= 0) return '0.00'
  return ((allocatedAmount.value / teamBonus.value) * 100).toFixed(2)
})

const canSave = computed(() => {
  return allocationPlayers.value.some(p => p.playerName && p.amount && parseFloat(p.amount) > 0)
})

const canSubmit = computed(() => {
  const hasValidPlayers = allocationPlayers.value.every(p => !p.error)
  const hasPositiveAmount = allocationPlayers.value.some(p => parseFloat(p.amount || 0) > 0)
  const isValidAmount = remainingAmount.value >= 0 && allocatedAmount.value > 0
  return hasValidPlayers && hasPositiveAmount && isValidAmount
})

// 格式化金额
const formatAmount = (amount) => {
  return `¥ ${parseFloat(amount || 0).toFixed(2)}`
}

// 格式化整数金额（用于个人奖金显示）
const formatInteger = (amount) => {
  const num = parseInt(amount || 0)
  return num.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ',')
}

// 获取选手分配金额
const getPlayerAmount = (playerId) => {
  const player = allocationPlayers.value.find(p => p.ID === playerId)
  return player ? player.amount : 0
}

// 加载战队列表
const loadTeams = async () => {
  teamsLoading.value = true
  try {
    const res = await getTeamList({ page: 1, pageSize: 100 })
    if (res.code === 0) {
      teams.value = res.data.list || []
    }
  } catch (error) {
    console.error('加载战队列表失败:', error)
    ElMessage.error('加载战队列表失败')
  } finally {
    teamsLoading.value = false
  }
}

// 刷新战队列表
const refreshTeams = () => {
  loadTeams()
}

// 选择战队
const selectTeam = (team) => {
  selectedTeam.value = team
  teamBonus.value = team.totalBounty || 0
  selectedPlayers.value = []
  allocationPlayers.value = []
  
  // 加载选手列表
  loadPlayers(team.ID)
  
  // 自动跳转到步骤2（选择选手）
  currentStep.value = 1
}

// 加载选手列表
const loadPlayers = async (teamId) => {
  playersLoading.value = true
  try {
    // 使用 id 参数筛选战队的选手
    const res = await getPlayerList({ id: teamId, page: 1, pageSize: 100 })
    if (res.code === 0) {
      allPlayers.value = res.data.list || []
      filterPlayers()
    }
  } catch (error) {
    console.error('加载选手列表失败:', error)
    ElMessage.error('加载选手列表失败')
  } finally {
    playersLoading.value = false
  }
}

// 搜索过滤选手
const filterPlayers = () => {
  if (!playerSearch.value) {
    filteredPlayers.value = allPlayers.value
  } else {
    filteredPlayers.value = allPlayers.value.filter(p =>
      p.playerName?.toLowerCase().includes(playerSearch.value.toLowerCase())
    )
  }
}

// 检查选手是否已选
const isPlayerSelected = (player) => {
  return selectedPlayers.value.some(p => p.ID === player.ID)
}

// 切换选手选择（最多选择4名选手）
const togglePlayer = (player) => {
  if (isPlayerSelected(player)) {
    removePlayer(player)
  } else {
    if (selectedPlayers.value.length < 4) {
      selectedPlayers.value.push(player)
    }
  }
}

// 移除选手
const removePlayer = (player) => {
  const index = selectedPlayers.value.findIndex(p => p.ID === player.ID)
  if (index > -1) {
    selectedPlayers.value.splice(index, 1)
    // 同时清零该选手的分配金额
    const allocationIndex = allocationPlayers.value.findIndex(p => p.ID === player.ID)
    if (allocationIndex > -1) {
      allocationPlayers.value[allocationIndex].amount = 0
    }
  }
}

// 步骤导航
const nextStep = () => {
  if (currentStep.value === 1) {
    // 从选择选手到分配奖金，初始化分配数据
    initAllocationPlayers()
  }
  if (currentStep.value < 3) {
    currentStep.value++
  }
}

const prevStep = () => {
  if (currentStep.value > 0) {
    currentStep.value--
  }
}

// 初始化分配选手
const initAllocationPlayers = () => {
  allocationPlayers.value = selectedPlayers.value.map(p => ({
    ID: p.ID,
    playerName: p.playerName,
    bounty: p.bounty,
    amount: '',
    percentage: 0,
    error: ''
  }))
}

// 金额变化处理（独立模块）
const handleAmountChange = (index) => {
  const player = allocationPlayers.value[index]
  let amount = parseFloat(player.amount)

  player.error = ''

  // 支持2位小数
  if (isNaN(amount) || amount < 0) {
    player.amount = 0
    player.percentage = 0
    return
  }
  
  // 保留2位小数
  amount = Math.round(amount * 100) / 100
  player.amount = amount

  // 根据金额计算百分比（相对于战队总奖金）
  player.percentage = teamBonus.value > 0 
    ? Math.min(100, Math.round((amount / teamBonus.value) * 100))
    : 0

  // 检查是否超出总额
  if (allocatedAmount.value > teamBonus.value) {
    player.error = `超出可用奖金，最多 ¥${teamBonus.value.toFixed(2)}`
  }
}

// 滑块变化处理（独立模块，直接设置金额）
const handleSliderChange = (index) => {
  const player = allocationPlayers.value[index]
  
  // 清除错误
  player.error = ''
  
  // 根据滑块百分比计算金额
  const amount = Math.round(teamBonus.value * player.percentage / 100 * 100) / 100
  player.amount = amount
  
  // 验证总分配是否超出
  if (allocatedAmount.value > teamBonus.value) {
    player.error = `超出可用奖金，最多 ¥${teamBonus.value.toFixed(2)}`
  }
}

// 设置快捷金额
const setQuickAmount = (index, amount) => {
  const player = allocationPlayers.value[index]
  player.amount = amount
  player.error = ''
  handleAmountChange(index)
}

// 平均分配
const autoDistribute = () => {
  const count = allocationPlayers.value.length
  if (count === 0) return

  const avgAmount = Math.round(teamBonus.value / count * 100) / 100
  allocationPlayers.value.forEach(p => {
    p.amount = avgAmount
    p.error = ''
  })
  // 处理最后一个选手分配剩余金额
  const total = avgAmount * count
  const remainder = Math.round((teamBonus.value - total) * 100) / 100
  if (Math.abs(remainder) > 0.01) {
    allocationPlayers.value[0].amount = Math.round((allocationPlayers.value[0].amount + remainder) * 100) / 100
  }
  allocationPlayers.value.forEach((p, index) => {
    handleAmountChange(index)
  })
}

// 保存草稿
const saveDraft = async () => {
  try {
    const data = {
      teamId: selectedTeam.value.ID,
      players: allocationPlayers.value.map(p => ({
        playerId: p.ID,
        name: p.playerName,
        amount: parseFloat(p.amount || 0)
      }))
    }
    const res = await saveDraftApi(data)
    if (res.code === 0) {
      ElMessage.success('草稿保存成功')
    }
  } catch (error) {
    console.error('保存草稿失败:', error)
    ElMessage.error('保存草稿失败')
  }
}

// 确认提交
const confirmSubmit = async () => {
  showConfirmDialog.value = false
  submitting.value = true

  try {
    // 验证：确保已选择战队
    if (!selectedTeam.value) {
      ElMessage.error('请先选择战队')
      return
    }

    // 验证：确保提交包含战队所有选手
    if (!allPlayers.value || allPlayers.value.length === 0) {
      ElMessage.error('该战队暂无选手信息')
      return
    }

    // 构建包含所有战队选手的提交数据
    // 已选中的选手使用分配的金额，未选中的选手金额为0
    const playerBounties = allPlayers.value.map(player => {
      // 检查该选手是否在已分配列表中
      const allocatedPlayer = allocationPlayers.value.find(p => p.ID === player.ID)
      
      return {
        playerId: player.ID,
        amount: allocatedPlayer ? parseFloat(allocatedPlayer.amount || 0) : 0
      }
    })

    // 验证数据完整性
    console.log('提交数据:', {
      teamId: selectedTeam.value.ID,
      teamName: selectedTeam.value.teamName,
      totalPlayers: allPlayers.value.length,
      selectedPlayers: selectedPlayers.value.length,
      playerBounties
    })

    const res = await allocateBounty({
      teamId: selectedTeam.value.ID,
      playerBounties
    })
    
    if (res.code === 0) {
      ElMessage.success('分配提交成功')
      // 重置状态
      currentStep.value = 0
      selectedTeam.value = null
      selectedPlayers.value = []
      allocationPlayers.value = []
      allPlayers.value = []
      loadTeams()
    } else {
      ElMessage.error(res.message || '提交失败')
    }
  } catch (error) {
    console.error('提交分配失败:', error)
    ElMessage.error('提交分配失败，请重试')
  } finally {
    submitting.value = false
  }
}

// 初始化
onMounted(() => {
  loadTeams()
})
</script>

<style scoped>
.team-bonus-container {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: 24px;
}

.page-title {
  font-size: 24px;
  font-weight: bold;
  color: #303133;
  display: flex;
  align-items: center;
  gap: 10px;
  margin: 0 0 8px 0;
}

.page-title :deep(.el-icon) {
  color: #409eff;
  font-size: 28px;
}

.page-subtitle {
  font-size: 14px;
  color: #909399;
  margin: 0 0 10px 0;
}

.type-indicator {
  display: flex;
  align-items: center;
  gap: 10px;
}

.type-note {
  font-size: 12px;
  color: #e6a23c;
}

/* 步骤指引 */
.steps-guide {
  background: #fff;
  padding: 20px;
  border-radius: 12px;
  margin-bottom: 24px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}

/* 步骤内容 */
.step-content {
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

/* 卡片样式 */
.section-card {
  background: #fff;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
  margin-bottom: 20px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.section-title {
  font-size: 18px;
  font-weight: bold;
  color: #303133;
  margin: 0;
  display: flex;
  align-items: center;
  gap: 8px;
}

.section-title :deep(.el-icon) {
  color: #409eff;
}

/* 加载状态 */
.loading-state {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  padding: 40px;
  color: #909399;
}

/* 战队选择 */
.team-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 16px;
}

.team-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px;
  border: 2px solid #e4e7ed;
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.3s ease;
  position: relative;
}

.team-card:hover {
  border-color: #409eff;
  box-shadow: 0 0 10px rgba(64, 158, 255, 0.1);
}

.team-card.active {
  border-color: #409eff;
  background: linear-gradient(135deg, rgba(64, 158, 255, 0.1) 0%, rgba(64, 158, 255, 0.05) 100%);
}

.team-avatar {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
}

.team-info {
  flex: 1;
}

.team-name {
  font-size: 16px;
  font-weight: bold;
  color: #303133;
  margin-bottom: 4px;
}

.team-bounty {
  font-size: 14px;
  color: #606266;
}

.bounty-value {
  color: #67c23a;
  font-weight: bold;
}

.selected-indicator {
  position: absolute;
  top: 8px;
  right: 8px;
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background: #409eff;
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
}

/* 选手选择 */
.player-search {
  width: 250px;
}

.selected-players {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
  margin-bottom: 16px;
  padding: 12px;
  background: #f5f7fa;
  border-radius: 8px;
}

.selected-label {
  font-size: 14px;
  color: #606266;
}

.player-tag {
  margin: 4px;
}

.player-list {
  min-height: 200px;
}

.player-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
  gap: 16px;
}

.player-card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  border: 2px solid #e4e7ed;
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.3s ease;
  position: relative;
}

.player-card:hover:not(.disabled) {
  border-color: #409eff;
}

.player-card.active {
  border-color: #409eff;
  background: linear-gradient(135deg, rgba(64, 158, 255, 0.1) 0%, rgba(64, 158, 255, 0.05) 100%);
}

.player-card.disabled {
  opacity: 0.5;
  cursor: not-allowed;
  background: #f5f7fa;
}

.player-card.disabled:hover {
  border-color: #e4e7ed;
}

.player-avatar {
  width: 44px;
  height: 44px;
  border-radius: 50%;
  background: #409eff;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
}

.player-info {
  flex: 1;
}

.player-name {
  font-size: 15px;
  font-weight: bold;
  color: #303133;
  margin-bottom: 4px;
}

.player-bounty {
  font-size: 13px;
  color: #909399;
}

.selected-badge {
  position: absolute;
  top: 8px;
  right: 8px;
  width: 22px;
  height: 22px;
  border-radius: 50%;
  background: #409eff;
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
}

/* 奖金分配信息条 */
.bonus-info-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 16px 24px;
  border-radius: 12px;
  margin-bottom: 20px;
  color: #fff;
}

.info-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.info-label {
  font-size: 12px;
  opacity: 0.8;
}

.info-value {
  font-size: 16px;
  font-weight: bold;
}

.info-value.highlight {
  color: #a8e063;
}

.info-value.warning {
  color: #ff6b6b;
}

/* 分配面板 */
.allocation-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 20px;
}

.allocation-card {
  background: #fafafa;
  border: 1px solid #e4e7ed;
  border-radius: 10px;
  padding: 20px;
  transition: all 0.3s ease;
}

.allocation-card:hover {
  border-color: #409eff;
  box-shadow: 0 0 10px rgba(64, 158, 255, 0.1);
}

.allocation-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 16px;
}

.allocation-header .personal-bounty {
  margin-left: auto;
  font-size: 13px;
  color: #909399;
  background: #f5f7fa;
  padding: 4px 10px;
  border-radius: 12px;
}

.player-index {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: bold;
}

.allocation-body {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.input-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.input-label {
  font-size: 14px;
  color: #606266;
}

.amount-input-wrapper {
  display: flex;
  align-items: center;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  padding: 0 10px;
  background: #fff;
}

.amount-input-wrapper.error {
  border-color: #f56c6c;
}

.currency {
  color: #909399;
  font-size: 14px;
}

.amount-input-wrapper :deep(.el-input) {
  width: 120px;
}

.amount-input-wrapper :deep(.el-input__wrapper) {
  box-shadow: none !important;
}

.amount-input-wrapper :deep(.el-input__inner) {
  border: none;
  padding-left: 5px;
}

.slider-group {
  padding: 0 4px;
}

.slider-label {
  display: block;
  font-size: 14px;
  color: #606266;
  margin-bottom: 8px;
}

.quick-buttons {
  display: flex;
  gap: 8px;
}

.quick-buttons :deep(.el-button) {
  flex: 1;
}

.error-message {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #f56c6c;
  font-size: 13px;
  margin-top: 8px;
}

/* 确认提交 */
.confirm-summary {
  background: #f5f7fa;
  border-radius: 10px;
  padding: 20px;
  margin-bottom: 20px;
}

.summary-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-bottom: 16px;
  border-bottom: 1px dashed #dcdfe6;
  margin-bottom: 16px;
}

.summary-header .team-name {
  font-size: 18px;
  font-weight: bold;
  color: #303133;
}

.summary-header .total-amount {
  font-size: 24px;
  font-weight: bold;
  color: #409eff;
}

.summary-table {
  font-size: 14px;
}

.table-header,
.table-row,
.table-footer {
  display: grid;
  grid-template-columns: 1fr 1fr 80px;
  padding: 10px 0;
  gap: 16px;
}

.table-header {
  font-weight: bold;
  color: #606266;
  border-bottom: 1px solid #e4e7ed;
}

.table-row {
  color: #303133;
}

.table-row .amount {
  font-weight: bold;
  color: #303133;
}

.table-row .percentage {
  color: #909399;
}

.table-footer {
  font-weight: bold;
  color: #303133;
  border-top: 1px solid #e4e7ed;
  margin-top: 10px;
  padding-top: 10px;
}

.table-footer .amount {
  color: #409eff;
}

.confirm-warning {
  margin-top: 20px;
}

/* 确认弹窗 */
.confirm-content {
  padding: 10px 0;
}

.confirm-info {
  text-align: center;
  padding-bottom: 20px;
  border-bottom: 1px solid #e4e7ed;
  margin-bottom: 20px;
}

.confirm-team {
  font-size: 18px;
  font-weight: bold;
  color: #303133;
  margin-bottom: 8px;
}

.confirm-amount {
  font-size: 28px;
  font-weight: bold;
  color: #409eff;
}

.confirm-players {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.confirm-player {
  display: flex;
  justify-content: space-between;
  padding: 8px 0;
  border-bottom: 1px dashed #e4e7ed;
}

.confirm-player:last-child {
  border-bottom: none;
}

.confirm-player span:first-child {
  color: #606266;
}

.confirm-player span:last-child {
  font-weight: bold;
  color: #303133;
}

/* 未分配选手样式 */
.confirm-player.unallocated {
  opacity: 0.6;
  background: #fafafa;
}

.confirm-player.unallocated span:first-child {
  color: #909399;
}

.confirm-player.unallocated .amount {
  color: #c0c4cc;
}

.confirm-player.unallocated .personal-bounty-small {
  color: #c0c4cc;
}

/* 提示样式 */
.confirm-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 8px;
}

/* 新的确认弹窗样式 */
.confirm-detail {
  background: #f5f7fa;
  border-radius: 10px;
  padding: 16px;
  margin-bottom: 20px;
}

.detail-header {
  display: flex;
  justify-content: space-between;
  padding: 0 0 12px 0;
  border-bottom: 1px solid #e4e7ed;
  margin-bottom: 12px;
  font-weight: bold;
  color: #606266;
  font-size: 14px;
}

.confirm-player .player-name {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.personal-bounty-small {
  font-size: 12px;
  color: #909399;
  font-weight: normal;
}

.confirm-player .amount {
  font-weight: bold;
  color: #409eff;
  font-size: 16px;
}

.confirm-summary .summary-item {
  display: flex;
  justify-content: space-between;
  padding: 8px 0;
  font-size: 14px;
}

.confirm-summary .summary-item .allocated {
  color: #409eff;
  font-weight: bold;
}

.confirm-summary .summary-item .warning {
  color: #f56c6c;
  font-weight: bold;
}

/* 操作按钮 */
.step-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
  margin-top: 20px;
}

.selection-hint {
  display: flex;
  flex-direction: column;
  gap: 4px;
  font-size: 14px;
  color: #606266;
}

.selection-hint .hint-warning {
  color: #e6a23c;
  font-size: 12px;
}

.action-buttons {
  display: flex;
  gap: 12px;
}

/* 响应式 */
@media (max-width: 768px) {
  .team-bonus-container {
    padding: 10px;
  }

  .bonus-info-bar {
    flex-wrap: wrap;
    gap: 12px;
    justify-content: center;
  }

  .info-item {
    min-width: 80px;
  }

  .team-grid,
  .player-grid,
  .allocation-grid {
    grid-template-columns: 1fr;
  }

  .step-actions {
    flex-direction: column;
  }

  .step-actions .el-button {
    width: 100%;
  }
}
</style>
