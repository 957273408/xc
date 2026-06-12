<template>
  <div class="ranking-container">
    <!-- 排行榜头部 -->
    <div class="ranking-header">
      <div class="header-title">
        <el-icon><Trophy /></el-icon>
        <span>队伍赏金排行榜</span>
      </div>
      <div class="header-stats">
        共 {{ total }} 支队伍
      </div>
    </div>

    <!-- 排行榜内容 -->
    <div class="ranking-content" v-loading="loading">
      <el-empty v-if="!loading && teamList.length === 0" description="暂无队伍数据" />

      <div class="ranking-list" v-else>
        <div
          v-for="team in teamList"
          :key="team.id"
          class="ranking-item"
          :class="getRankClass(team.rank)"
        >
          <!-- 排名标识 -->
          <div class="rank-badge">
            <span v-if="team.rank <= 3" class="special-rank">
              <el-icon v-if="team.rank === 1"><Trophy /></el-icon>
              <el-icon v-else-if="team.rank === 2"><Medal /></el-icon>
              <el-icon v-else><Star /></el-icon>
            </span>
            <span v-else class="normal-rank">{{ team.rank }}</span>
          </div>

          <!-- 队伍信息 -->
          <div class="team-info">
            <div class="team-name">{{ team.teamName }}</div>
            <div class="team-id">ID: {{ team.id }}</div>
          </div>

          <!-- 赏金金额 -->
          <div class="bounty-amount">
            <span class="amount-value">¥ {{ team.totalBounty?.toFixed(2) || '0.00' }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 分页 -->
    <div class="ranking-pagination" v-if="teamList.length > 0">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[10, 20, 50]"
        :total="total"
        layout="total, prev, pager, next"
        @current-change="handlePageChange"
        @size-change="handleSizeChange"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getTeamBountyRanking } from '@/api/bountyRecord'
import { Trophy, Medal, Star } from '@element-plus/icons-vue'

defineOptions({
  name: 'TeamRanking'
})

const loading = ref(false)
const teamList = ref([])
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

// 获取排名样式类名
const getRankClass = (rank) => {
  if (rank === 1) return 'rank-first'
  if (rank === 2) return 'rank-second'
  if (rank === 3) return 'rank-third'
  return ''
}

// 获取队伍排行榜数据
const fetchTeamRanking = async () => {
  loading.value = true
  try {
    const res = await getTeamBountyRanking({
      page: currentPage.value,
      pageSize: pageSize.value
    })
    if (res.code === 0) {
      teamList.value = res.data.list || []
      total.value = res.data.total || 0
    }
  } catch (error) {
    console.error('获取队伍排行榜失败:', error)
  } finally {
    loading.value = false
  }
}

// 分页变化
const handlePageChange = () => {
  fetchTeamRanking()
}

const handleSizeChange = () => {
  currentPage.value = 1
  fetchTeamRanking()
}

onMounted(() => {
  fetchTeamRanking()
})
</script>

<style scoped>
.ranking-container {
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.ranking-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
  padding: 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-title {
  display: flex;
  align-items: center;
  font-size: 18px;
  font-weight: bold;
}

.header-title .el-icon {
  margin-right: 10px;
  font-size: 24px;
}

.header-stats {
  font-size: 14px;
}

.ranking-content {
  max-height: 500px;
  overflow-y: auto;
}

.ranking-list {
  padding: 10px 0;
}

.ranking-item {
  display: flex;
  align-items: center;
  padding: 15px 20px;
  border-bottom: 1px solid #f0f0f0;
  transition: all 0.3s;
}

.ranking-item:hover {
  background: #f5f7fa;
}

.ranking-item:last-child {
  border-bottom: none;
}

/* 前三名特殊样式 */
.rank-first {
  background: linear-gradient(90deg, #fff5f5 0%, #fff 100%);
}

.rank-second {
  background: linear-gradient(90deg, #f5f5ff 0%, #fff 100%);
}

.rank-third {
  background: linear-gradient(90deg, #fffaf0 0%, #fff 100%);
}

.rank-badge {
  width: 50px;
  height: 50px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 20px;
}

.special-rank {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 28px;
}

.rank-first .special-rank {
  background: linear-gradient(135deg, #ffd700 0%, #ffed4e 100%);
  color: #fff;
  box-shadow: 0 4px 10px rgba(255, 215, 0, 0.3);
}

.rank-second .special-rank {
  background: linear-gradient(135deg, #c0c0c0 0%, #e8e8e8 100%);
  color: #fff;
  box-shadow: 0 4px 10px rgba(192, 192, 192, 0.3);
}

.rank-third .special-rank {
  background: linear-gradient(135deg, #cd7f32 0%, #daa06d 100%);
  color: #fff;
  box-shadow: 0 4px 10px rgba(205, 127, 50, 0.3);
}

.normal-rank {
  width: 40px;
  height: 40px;
  border-radius: 8px;
  background: #f0f0f0;
  color: #606266;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  font-weight: bold;
}

.team-info {
  flex: 1;
}

.team-name {
  font-size: 16px;
  font-weight: bold;
  color: #303133;
  margin-bottom: 5px;
}

.team-id {
  font-size: 12px;
  color: #909399;
}

.bounty-amount {
  text-align: right;
}

.amount-value {
  font-size: 18px;
  font-weight: bold;
  color: #67c23a;
}

.ranking-pagination {
  padding: 20px;
  display: flex;
  justify-content: center;
  background: #f5f7fa;
}

/* 响应式适配 */
@media (max-width: 768px) {
  .ranking-header {
    padding: 15px;
  }

  .header-title {
    font-size: 16px;
  }

  .ranking-item {
    padding: 12px 15px;
  }

  .rank-badge {
    width: 40px;
    height: 40px;
    margin-right: 15px;
  }

  .special-rank {
    width: 40px;
    height: 40px;
    font-size: 24px;
  }

  .team-name {
    font-size: 14px;
  }

  .amount-value {
    font-size: 16px;
  }
}
</style>