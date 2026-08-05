<template>
  <div class="score-chart">
    <div class="chart-header">
      <h3>积分统计</h3>
    </div>

    <div class="total-score-section">
      <div class="score-circle">
        <span class="score-number">{{ totalScore }}</span>
        <span class="score-unit">分</span>
      </div>
      <div class="score-details">
        <div class="detail-item">
          <span class="label">比赛场次</span>
          <span class="value">{{ matchCount }}</span>
        </div>
        <div class="detail-item">
          <span class="label">最近排名</span>
          <span class="value">{{ lastRank > 0 ? '#' + lastRank : '-' }}</span>
        </div>
        <div class="detail-item">
          <span class="label">累计赏金</span>
          <span class="value bounty-value">{{ totalBountyCoin.toLocaleString() }}</span>
        </div>
      </div>
    </div>

    <div class="recent-records" v-if="scores.length > 0">
      <h4>最近 {{ scores.length }} 次记录</h4>
      <div class="records-list">
        <div
          v-for="(score, index) in scores"
          :key="score.id"
          class="record-item"
          :class="{ 'rank-high': score.rank <= 3 }"
        >
          <div class="record-rank">
            <span class="rank-num">{{ score.rank || '-' }}</span>
          </div>
          <div class="record-info">
            <div class="record-header">
              <span class="war-id">{{ score.warId }}</span>
              <span class="record-time">{{ score.settleTime }}</span>
            </div>
            <div class="record-stats">
              <span>淘汰 {{ score.killCount }} 人</span>
              <span>排名分 {{ score.rankScore }}</span>
              <span>淘汰分 {{ score.killScore }}</span>
              <span v-if="score.bountyCoin > 0" class="record-bounty">赏金 {{ score.bountyCoin.toLocaleString() }}</span>
            </div>
          </div>
          <div class="record-score">
            <span class="score-total">+{{ score.totalScore }}</span>
          </div>
          <div class="record-actions">
            <el-button type="primary" size="small" @click="$emit('editScore', score)">
              <el-icon><Edit /></el-icon>
            </el-button>
            <el-button type="danger" size="small" @click="$emit('deleteScore', score)">
              <el-icon><Delete /></el-icon>
            </el-button>
          </div>
          <div class="rank-badge-indicator" v-if="index === 0">
            <span class="latest-tag">最新</span>
          </div>
        </div>
      </div>
    </div>

    <div class="empty-state" v-else>
      <el-empty description="暂无积分记录，请添加WarId获取积分" />
    </div>

    <div class="score-breakdown" v-if="scores.length > 0">
      <h4>积分构成</h4>
      <div class="breakdown-bar">
        <div class="bar-segment rank-segment" :style="{ width: rankScorePercent + '%' }">
          <span>排名分: {{ totalRankScore }}</span>
        </div>
        <div class="bar-segment kill-segment" :style="{ width: killScorePercent + '%' }">
          <span>淘汰分: {{ totalKillScore }}</span>
        </div>
      </div>
    </div>

    <!-- 编辑积分对话框 -->
    <el-dialog v-model="editDialogVisible" title="编辑积分" width="500px" @close="handleEditClose">
      <el-form :model="editForm" label-width="80px">
        <el-form-item label="排名">
          <el-input-number v-model="editForm.rank" :min="0" :max="16" />
        </el-form-item>
        <el-form-item label="淘汰数">
          <el-input-number v-model="editForm.killCount" :min="0" />
        </el-form-item>
        <el-form-item label="排名分">
          <el-input-number v-model="editForm.rankScore" :min="0" />
        </el-form-item>
        <el-form-item label="淘汰分">
          <el-input-number v-model="editForm.killScore" :min="0" />
        </el-form-item>
        <el-form-item label="总积分">
          <el-input-number v-model="editForm.totalScore" :min="0" />
        </el-form-item>
        <el-form-item label="赏金">
          <el-input-number v-model="editForm.bountyCoin" :min="0" />
        </el-form-item>
        <el-form-item label="结算时间">
          <el-date-picker
            v-model="editForm.settleTime"
            type="datetime"
            placeholder="选择结算时间"
            value-format="YYYY-MM-DD HH:mm:ss"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="handleSaveEdit">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Edit, Delete } from '@element-plus/icons-vue'
import { updateTeamScore, deleteTeamScore } from '@/api/competitionTeam'

const props = defineProps({
  scores: {
    type: Array,
    default: () => []
  },
  totalScore: {
    type: Number,
    default: 0
  },
  matchCount: {
    type: Number,
    default: 0
  },
  lastRank: {
    type: Number,
    default: 0
  },
  teamId: {
    type: Number,
    default: 0
  }
})

const emit = defineEmits(['editScore', 'deleteScore', 'refresh'])

const totalRankScore = computed(() => {
  return props.scores.reduce((sum, s) => sum + (s.rankScore || 0), 0)
})

const totalKillScore = computed(() => {
  return props.scores.reduce((sum, s) => sum + (s.killScore || 0), 0)
})

const totalBountyCoin = computed(() => {
  return props.scores.reduce((sum, s) => sum + (s.bountyCoin || 0), 0)
})

const rankScorePercent = computed(() => {
  const total = totalRankScore.value + totalKillScore.value
  if (total === 0) return 50
  return (totalRankScore.value / total) * 100
})

const killScorePercent = computed(() => {
  return 100 - rankScorePercent.value
})

// 编辑积分
const editDialogVisible = ref(false)
const saving = ref(false)
const editForm = ref({
  id: null,
  teamId: null,
  rank: 0,
  killCount: 0,
  rankScore: 0,
  killScore: 0,
  totalScore: 0,
  bountyCoin: 0,
  settleTime: ''
})

const handleEditScore = (score) => {
  editForm.value = {
    id: score.id,
    teamId: props.teamId,
    rank: score.rank,
    killCount: score.killCount,
    rankScore: score.rankScore,
    killScore: score.killScore,
    totalScore: score.totalScore,
    bountyCoin: score.bountyCoin || 0,
    settleTime: score.settleTime
  }
  editDialogVisible.value = true
}

const handleSaveEdit = async () => {
  saving.value = true
  try {
    const res = await updateTeamScore(editForm.value)
    if (res.code === 0) {
      ElMessage.success('修改成功')
      editDialogVisible.value = false
      emit('refresh')
    } else {
      ElMessage.error(res.msg || '修改失败')
    }
  } catch (error) {
    ElMessage.error('接口调用失败')
  } finally {
    saving.value = false
  }
}

const handleEditClose = () => {
  editDialogVisible.value = false
}

// 删除积分
const handleDeleteScore = async (score) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除 WarId ${score.warId} 的积分记录吗？删除后不可恢复。`,
      '删除确认',
      {
        confirmButtonText: '确定删除',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    const res = await deleteTeamScore({
      teamId: props.teamId,
      warId: score.warId
    })
    if (res.code === 0) {
      ElMessage.success('删除成功')
      emit('refresh')
    } else {
      ElMessage.error(res.msg || '删除失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

// 暴露给父组件调用
defineExpose({
  handleEditScore,
  handleDeleteScore
})
</script>

<style scoped>
.score-chart {
  padding: 20px;
}

.chart-header h3 {
  margin: 0 0 20px 0;
  color: #fff;
}

.total-score-section {
  display: flex;
  align-items: center;
  gap: 30px;
  margin-bottom: 30px;
  padding: 20px;
  background: linear-gradient(135deg, rgba(255, 215, 0, 0.1) 0%, rgba(255, 107, 107, 0.1) 100%);
  border-radius: 16px;
  border: 1px solid rgba(255, 215, 0, 0.3);
}

.score-circle {
  width: 100px;
  height: 100px;
  border-radius: 50%;
  background: linear-gradient(135deg, #ffd700 0%, #ff6b6b 100%);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  box-shadow: 0 8px 30px rgba(255, 215, 0, 0.4);
}

.score-number {
  font-size: 32px;
  font-weight: bold;
  color: #fff;
  line-height: 1;
}

.score-unit {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.8);
}

.score-details {
  flex: 1;
  display: flex;
  gap: 30px;
}

.detail-item {
  display: flex;
  flex-direction: column;
}

.detail-item .label {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.5);
  margin-bottom: 4px;
}

.detail-item .value {
  font-size: 24px;
  font-weight: bold;
  color: #fff;
}

.bounty-value {
  color: #ff6b6b !important;
}

.recent-records h4,
.score-breakdown h4 {
  margin: 0 0 12px 0;
  color: #fff;
}

.records-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.record-item {
  display: flex;
  align-items: center;
  gap: 16px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 10px;
  padding: 12px 16px;
  transition: all 0.3s;
}

.record-item:hover {
  background: rgba(255, 255, 255, 0.08);
}

.record-item.rank-high {
  border-color: rgba(255, 215, 0, 0.4);
  background: rgba(255, 215, 0, 0.05);
}

.record-rank {
  width: 48px;
  height: 48px;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.1);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.rank-num {
  font-size: 20px;
  font-weight: bold;
  color: #fff;
}

.record-info {
  flex: 1;
  min-width: 0;
}

.record-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 4px;
}

.war-id {
  font-family: monospace;
  font-size: 13px;
  color: #ffd700;
}

.record-time {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.4);
}

.record-stats {
  display: flex;
  gap: 12px;
  font-size: 12px;
  color: rgba(255, 255, 255, 0.6);
}

.record-bounty {
  color: #ff6b6b;
  font-weight: bold;
}

.record-score {
  flex-shrink: 0;
}

.score-total {
  font-size: 20px;
  font-weight: bold;
  color: #ffd700;
}

.record-actions {
  display: flex;
  gap: 6px;
  flex-shrink: 0;
}

.record-actions .el-button {
  padding: 6px;
}

.record-actions .el-icon {
  font-size: 14px;
}

.rank-badge-indicator {
  flex-shrink: 0;
}

.latest-tag {
  font-size: 11px;
  background: linear-gradient(135deg, #ff6b6b, #ffd700);
  color: #fff;
  padding: 2px 8px;
  border-radius: 10px;
}

.empty-state {
  padding: 40px;
}

.score-breakdown {
  margin-top: 20px;
}

.breakdown-bar {
  display: flex;
  height: 32px;
  border-radius: 16px;
  overflow: hidden;
  background: rgba(255, 255, 255, 0.1);
}

.bar-segment {
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 12px;
  font-weight: bold;
  transition: width 0.5s ease;
}

.bar-segment.rank-segment {
  background: linear-gradient(90deg, #667eea, #764ba2);
}

.bar-segment.kill-segment {
  background: linear-gradient(90deg, #ff6b6b, #ffd700);
}

@media (max-width: 768px) {
  .total-score-section {
    flex-direction: column;
    text-align: center;
  }

  .score-details {
    justify-content: center;
  }

  .record-item {
    flex-wrap: wrap;
    gap: 8px;
  }

  .record-info {
    flex-basis: 100%;
    order: 3;
  }

  .bar-segment {
    font-size: 10px;
  }
}
</style>
