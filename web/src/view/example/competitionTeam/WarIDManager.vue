<template>
  <div class="warid-manager">
    <!-- 批量计算区域 -->
    <div class="batch-section">
      <h3>WarId 批量积分计算</h3>
      <p class="section-desc">输入 WarId，系统将自动为所有战队计算积分，验证后批量保存</p>
      <div class="input-row">
        <el-input
          v-model="warId"
          placeholder="请输入 WarId"
          clearable
          @keyup.enter="handleCalculate"
        />
        <el-button
          type="primary"
          :loading="calculating"
          :disabled="!warId.trim()"
          @click="handleCalculate"
        >
          {{ calculating ? '计算中...' : '计算积分' }}
        </el-button>
      </div>
    </div>

    <!-- 计算结果预览 -->
    <div class="result-section" v-if="calcResult">
      <div class="result-header">
        <h4>
          计算结果 - WarId: {{ calcResult.warId }}
          <el-tag size="small" type="success">匹配 {{ calcResult.matchedNum }}/{{ calcResult.totalTeams }}</el-tag>
        </h4>
        <div class="result-actions">
          <el-button
            type="success"
            :loading="confirming"
            :disabled="selectedTeamIds.length === 0"
            @click="handleConfirm"
          >
            {{ confirming ? '保存中...' : `确认保存 (${selectedTeamIds.length})` }}
          </el-button>
        </div>
      </div>

      <el-table
        :data="calcResult.items"
        style="width: 100%"
        @selection-change="handleSelectionChange"
        ref="tableRef"
      >
        <el-table-column type="selection" width="50" />
        <el-table-column prop="teamCode" label="战队标识" width="100" />
        <el-table-column prop="teamName" label="战队名称" width="150" />
        <el-table-column label="匹配状态" width="120">
          <template #default="{ row }">
            <el-tag :type="row.matched ? 'success' : 'danger'" size="small" effect="dark">
              {{ row.matched ? `已匹配 ${row.playerCount}人` : '未匹配' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="killCount" label="淘汰数" width="80" />
        <el-table-column prop="rank" label="排名" width="80">
          <template #default="{ row }">
            <span v-if="row.matched" :class="{ 'rank-good': row.rank <= 3 }">#{{ row.rank }}</span>
            <span v-else class="rank-na">-</span>
          </template>
        </el-table-column>
        <el-table-column prop="rankScore" label="排名分" width="80">
          <template #default="{ row }">{{ row.matched ? row.rankScore : '-' }}</template>
        </el-table-column>
        <el-table-column prop="killScore" label="淘汰分" width="80">
          <template #default="{ row }">{{ row.matched ? row.killScore : '-' }}</template>
        </el-table-column>
        <el-table-column prop="bountyCoin" label="赏金" width="100">
          <template #default="{ row }">
            <span v-if="row.matched && row.bountyCoin > 0" class="bounty-coin">{{ row.bountyCoin.toLocaleString() }}</span>
            <span v-else class="score-na">-</span>
          </template>
        </el-table-column>
        <el-table-column prop="totalScore" label="总积分" width="100">
          <template #default="{ row }">
            <span class="total-score" :class="{ 'score-na': !row.matched }">{{ row.matched ? row.totalScore : '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="message" label="备注" min-width="120">
          <template #default="{ row }">{{ row.message || (row.matched ? '正常' : '昵称前缀未匹配到战场数据') }}</template>
        </el-table-column>
      </el-table>
    </div>

    <!-- 确认结果反馈 -->
    <div class="confirm-result" v-if="confirmResult">
      <el-alert
        :title="`保存完成：成功 ${confirmResult.successCount} 条，失败 ${confirmResult.failCount} 条`"
        :type="confirmResult.failCount > 0 ? 'warning' : 'success'"
        show-icon
        :closable="false"
      >
        <template #default v-if="confirmResult.errors.length > 0">
          <div class="error-list">
            <div v-for="(err, i) in confirmResult.errors" :key="i">{{ err }}</div>
          </div>
        </template>
      </el-alert>
    </div>

    <!-- 积分规则说明 -->
    <div class="rules-section">
      <h4>积分规则</h4>
      <div class="rules-grid">
        <div class="rule-item">#1 = 16分</div>
        <div class="rule-item">#2 = 12分</div>
        <div class="rule-item">#3 = 10分</div>
        <div class="rule-item">#4 = 8分</div>
        <div class="rule-item">#5 = 6分</div>
        <div class="rule-item">#6 = 5分</div>
        <div class="rule-item">#7 = 4分</div>
        <div class="rule-item">#8 = 3分</div>
        <div class="rule-item">#9 = 2分</div>
        <div class="rule-item">#10 = 1分</div>
        <div class="rule-item">#11-16 = 0分</div>
        <div class="rule-item highlight">每淘汰1人 +1分</div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import { calculateWarIDForAllTeams, confirmWarIDScores } from '@/api/competitionTeam'

const warId = ref('')
const calculating = ref(false)
const confirming = ref(false)
const calcResult = ref(null)
const confirmResult = ref(null)
const selectedTeamIds = ref([])
const tableRef = ref(null)

const handleCalculate = async () => {
  if (!warId.value.trim()) {
    ElMessage.warning('请输入 WarId')
    return
  }

  calculating.value = true
  calcResult.value = null
  confirmResult.value = null
  selectedTeamIds.value = []

  try {
    const res = await calculateWarIDForAllTeams({ warId: warId.value.trim() })
    if (res.code === 0) {
      calcResult.value = res.data
      ElMessage.success(`计算完成，匹配 ${res.data.matchedNum}/${res.data.totalTeams} 个战队`)
      // 自动选中所有匹配的战队
      nextTick(() => {
        if (calcResult.value) {
          calcResult.value.items.forEach((item) => {
            if (item.matched && tableRef.value) {
              tableRef.value.toggleRowSelection(item, true)
            }
          })
        }
      })
    } else {
      ElMessage.error(res.msg || '计算失败')
    }
  } catch (error) {
    ElMessage.error('接口调用失败')
  } finally {
    calculating.value = false
  }
}

const handleSelectionChange = (selection) => {
  selectedTeamIds.value = selection.map((item) => item.teamId)
}

const handleConfirm = async () => {
  if (selectedTeamIds.value.length === 0) {
    ElMessage.warning('请至少选择一个战队')
    return
  }

  confirming.value = true
  confirmResult.value = null

  try {
    const res = await confirmWarIDScores({
      warId: warId.value.trim(),
      teamIds: selectedTeamIds.value
    })
    if (res.code === 0) {
      confirmResult.value = res.data
      const msg = `保存完成：成功 ${res.data.successCount} 条，失败 ${res.data.failCount} 条`
      if (res.data.failCount > 0) {
        ElMessage.warning(msg)
      } else {
        ElMessage.success(msg)
      }
      // 清除计算结果
      calcResult.value = null
      warId.value = ''
    } else {
      ElMessage.error(res.msg || '保存失败')
    }
  } catch (error) {
    ElMessage.error('接口调用失败')
  } finally {
    confirming.value = false
  }
}
</script>

<style scoped>
.warid-manager {
  padding: 20px;
}

.batch-section {
  background: rgba(255, 255, 255, 0.05);
  padding: 20px;
  border-radius: 12px;
  margin-bottom: 20px;
}

.batch-section h3 {
  margin: 0 0 8px 0;
  color: #fff;
}

.section-desc {
  color: rgba(255, 255, 255, 0.5);
  font-size: 13px;
  margin: 0 0 16px 0;
}

.input-row {
  display: flex;
  gap: 12px;
}

.input-row .el-input {
  flex: 1;
}

.result-section {
  margin-bottom: 20px;
}

.result-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.result-header h4 {
  margin: 0;
  color: #fff;
  display: flex;
  align-items: center;
  gap: 8px;
}

.result-actions {
  display: flex;
  gap: 8px;
}

.rank-good {
  color: #ffd700;
  font-weight: bold;
}

.rank-na,
.score-na {
  color: rgba(255, 255, 255, 0.3);
}

.bounty-coin {
  color: #ff6b6b;
  font-weight: bold;
}

.total-score {
  color: #ffd700;
  font-weight: bold;
  font-size: 16px;
}

.total-score.score-na {
  font-weight: normal;
  font-size: 14px;
}

.confirm-result {
  margin-bottom: 20px;
}

.error-list {
  margin-top: 8px;
  font-size: 12px;
}

.error-list div {
  margin-bottom: 4px;
}

.rules-section {
  background: rgba(255, 215, 0, 0.05);
  border: 1px solid rgba(255, 215, 0, 0.2);
  padding: 16px;
  border-radius: 8px;
}

.rules-section h4 {
  margin: 0 0 12px 0;
  color: #fff;
}

.rules-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(100px, 1fr));
  gap: 8px;
}

.rule-item {
  background: rgba(255, 255, 255, 0.05);
  padding: 6px 10px;
  border-radius: 4px;
  text-align: center;
  font-size: 13px;
  color: rgba(255, 255, 255, 0.7);
}

.rule-item.highlight {
  background: rgba(255, 215, 0, 0.2);
  color: #ffd700;
  font-weight: bold;
}

:deep(.el-table) {
  background: transparent;
  color: rgba(255, 255, 255, 0.8);
}

:deep(.el-table th) {
  background: rgba(255, 255, 255, 0.05) !important;
  color: rgba(255, 255, 255, 0.9);
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

:deep(.el-table tr) {
  background: transparent;
}

:deep(.el-table td) {
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

:deep(.el-table__row:hover > td) {
  background: rgba(255, 255, 255, 0.05) !important;
}

@media (max-width: 768px) {
  .input-row {
    flex-direction: column;
  }

  .rules-grid {
    grid-template-columns: repeat(auto-fill, minmax(80px, 1fr));
  }
}
</style>
