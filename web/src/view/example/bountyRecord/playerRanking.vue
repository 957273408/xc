<template>
  <div class="ranking-container">
    <!-- 排行榜头部 -->
    <div class="ranking-header">
      <div class="header-title">
        <el-icon><User /></el-icon>
        <span>选手赏金排行榜</span>
      </div>
      <div class="header-stats">
        共 {{ total }} 位选手
      </div>
    </div>

    <!-- 排行榜内容 -->
    <div class="ranking-content" v-loading="loading">
      <el-empty v-if="!loading && playerList.length === 0" description="暂无选手数据" />

      <div class="ranking-list" v-else>
        <div
          v-for="player in playerList"
          :key="player.id"
          class="ranking-item"
          :class="getRankClass(player.rank)"
        >
          <!-- 排名标识 -->
          <div class="rank-badge">
            <span v-if="player.rank <= 3" class="special-rank">
              <el-icon v-if="player.rank === 1"><Trophy /></el-icon>
              <el-icon v-else-if="player.rank === 2"><Medal /></el-icon>
              <el-icon v-else><Star /></el-icon>
            </span>
            <span v-else class="normal-rank">{{ player.rank }}</span>
          </div>

          <!-- 选手头像 -->
          <div class="player-avatar">
            <el-avatar
              :size="50"
              :src="player.avatar || 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'"
            />
          </div>

          <!-- 选手信息 -->
          <div class="player-info">
            <div class="player-name">{{ player.playerName }}</div>
            <div class="player-team">{{ player.teamName || '无战队' }}</div>
          </div>

          <!-- 赏金金额 -->
          <div class="bounty-amount">
            <span class="amount-value">¥ {{ player.bounty?.toFixed(2) || '0.00' }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 分页 -->
    <div class="ranking-pagination" v-if="playerList.length > 0">
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
import { getPlayerBountyRanking } from '@/api/bountyRecord'
import { User, Trophy, Medal, Star } from '@element-plus/icons-vue'

defineOptions({
  name: 'PlayerRanking'
})

const loading = ref(false)
const playerList = ref([])
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

// 获取选手排行榜数据
const fetchPlayerRanking = async () => {
  loading.value = true
  try {
    const res = await getPlayerBountyRanking({
      page: currentPage.value,
      pageSize: pageSize.value
    })
    if (res.code === 0) {
      playerList.value = res.data.list || []
      total.value = res.data.total || 0
    }
  } catch (error) {
    console.error('获取选手排行榜失败:', error)
  } finally {
    loading.value = false
  }
}

// 分页变化
const handlePageChange = () => {
  fetchPlayerRanking()
}

const handleSizeChange = () => {
  currentPage.value = 1
  fetchPlayerRanking()
}

onMounted(() => {
  fetchPlayerRanking()
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
  background: linear-gradient(135deg, #11998e 0%, #38ef7d 100%);
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
  margin-right: 15px;
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

.player-avatar {
  margin-right: 15px;
}

.player-info {
  flex: 1;
}

.player-name {
  font-size: 16px;
  font-weight: bold;
  color: #303133;
  margin-bottom: 5px;
}

.player-team {
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
    margin-right: 10px;
  }

  .special-rank {
    width: 40px;
    height: 40px;
    font-size: 24px;
  }

  .player-avatar {
    margin-right: 10px;
  }

  .player-name {
    font-size: 14px;
  }

  .amount-value {
    font-size: 16px;
  }
}
</style>