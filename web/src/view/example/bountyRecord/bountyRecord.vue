<template>
  <div class="bounty-pool-container">
    <!-- 赏金池概览 -->
    <div class="pool-overview">
      <el-row :gutter="20">
        <el-col :xs="24" :sm="24">
          <div class="stat-card pool-total">
            <div class="stat-icon">
              <el-icon><Wallet /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-label">赏金池总额</div>
              <div class="stat-value">¥ {{ poolInfo.totalAmount?.toFixed(2) || '0.00' }}</div>
            </div>
          </div>
        </el-col>
      </el-row>
    </div>

      <div class="pool-overview">
      <el-row :gutter="20">
        <el-col :xs="24" :sm="12">
          <TeamRanking />
        </el-col>
        <el-col :xs="24" :sm="12">
        <PlayerRanking />
        </el-col>
    
      </el-row>
    </div>

    <!-- 赏金变更记录 -->
    <div class="section-wrapper">
      <!-- 筛选区域 -->
      <div class="filter-section">
        <el-form :inline="true" :model="searchForm" class="search-form">
          <el-form-item label="变动类型">
            <el-select v-model="searchForm.changeType" placeholder="请选择" clearable>
              <el-option label="全部" value="" />
              <el-option label="击杀获取" value="kill" />
              <el-option label="被击杀损失" value="killed" />
              <el-option label="复活损失" value="revive" />
              <el-option label="赏金分配" value="allocate" />
              <el-option label="赏金池增加" value="pool_add" />
              <el-option label="赏金池减少" value="pool_reduce" />
              <el-option label="从赏金池领取" value="pool_claim" />
              <el-option label="初始赏金" value="init" />
            </el-select>
          </el-form-item>
          <el-form-item label="选手筛选">
            <el-input v-model="searchForm.playerName" placeholder="选手姓名" clearable />
          </el-form-item>
          <el-form-item label="战队筛选">
            <el-input v-model="searchForm.teamName" placeholder="战队名称" clearable />
          </el-form-item>
          <el-form-item label="时间范围">
            <el-date-picker
              v-model="dateRange"
              type="daterange"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              value-format="YYYY-MM-DD"
              @change="handleDateChange"
            />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="handleSearch">搜索</el-button>
            <el-button @click="handleReset">重置</el-button>
          </el-form-item>
        </el-form>
      </div>

      <!-- 赏金变更记录列表 -->
      <div class="record-section">
        <div class="section-header">
          <span class="section-title">赏金变更记录</span>
          <span class="record-count">共 {{ total }} 条记录</span>
        </div>

        <el-table
          v-loading="loading"
          :data="tableData"
          style="width: 100%"
          row-key="ID"
          :header-cell-style="{ background: '#f5f7fa', color: '#606266' }"
        >
          <el-table-column prop="CreatedAt" label="操作时间" width="180" align="center">
            <template #default="scope">
              {{ formatDate(scope.row.CreatedAt) }}
            </template>
          </el-table-column>
          <el-table-column prop="changeType" label="变动类型" width="120" align="center">
            <template #default="scope">
              <el-tag :type="getChangeTypeTag(scope.row.changeType)">
                {{ getChangeTypeText(scope.row.changeType) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="playerName" label="操作人" width="120" align="center" />
          <el-table-column prop="teamName" label="所属战队" width="150" align="center" />
          <el-table-column prop="amount" label="变动金额" width="120" align="center">
            <template #default="scope">
              <span :class="getAmountClass(scope.row.amount)">
                {{ scope.row.amount > 0 ? '+' : '' }}{{ scope.row.amount?.toFixed(2) }}
              </span>
            </template>
          </el-table-column>
          <el-table-column prop="balance" label="变动后余额" width="120" align="center">
            <template #default="scope">
              {{ scope.row.balance?.toFixed(2) }}
            </template>
          </el-table-column>
          <el-table-column prop="reason" label="变动原因" min-width="200" />
          <el-table-column prop="relatedName" label="关联对象" width="120" align="center">
            <template #default="scope">
              {{ scope.row.relatedName || '-' }}
            </template>
          </el-table-column>
        </el-table>

        <!-- 空数据提示 -->
        <el-empty v-if="!loading && tableData.length === 0" description="暂无操作记录" />

        <!-- 分页 -->
        <div class="pagination" v-if="tableData.length > 0">
          <el-pagination
            v-model:current-page="page"
            v-model:page-size="pageSize"
            :page-sizes="[10, 20, 50, 100]"
            :total="total"
            layout="total, sizes, prev, pager, next, jumper"
            @current-change="handlePageChange"
            @size-change="handleSizeChange"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { getPoolInfo, getRecordList } from '@/api/bountyRecord'
import { formatDate } from '@/utils/format'
import TeamRanking from './teamRanking.vue'
import PlayerRanking from './playerRanking.vue'

defineOptions({
  name: 'BountyPool'
})

// 数据状态
const loading = ref(false)
const tableData = ref([])
const poolInfo = ref({})
const poolStats = ref({ totalIn: 0, totalOut: 0 })
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const dateRange = ref([])

// 筛选表单
const searchForm = reactive({
  changeType: '',
  playerName: '',
  teamName: ''
})

// 变动类型映射（完整）
const changeTypeMap = {
  kill: { text: '击杀获取', type: 'success' },
  killed: { text: '被击杀损失', type: 'danger' },
  revive: { text: '复活损失', type: 'warning' },
  allocate: { text: '赏金分配', type: 'primary' },
  pool_add: { text: '赏金池增加', type: 'success' },
  pool_reduce: { text: '赏金池减少', type: 'danger' },
  pool_claim: { text: '领取赏金', type: 'warning' },
  init: { text: '初始赏金', type: 'info' }
}

const getChangeTypeText = (type) => {
  return changeTypeMap[type]?.text || type
}

const getChangeTypeTag = (type) => {
  return changeTypeMap[type]?.type || 'info'
}

const getAmountClass = (amount) => {
  return amount > 0 ? 'amount-positive' : 'amount-negative'
}

// 获取赏金池信息
const fetchPoolInfo = async () => {
  try {
    const res = await getPoolInfo()
    if (res.code === 0) {
      poolInfo.value = res.data || {}
    }
  } catch (error) {
    console.error('获取赏金池信息失败:', error)
  }
}

// 计算赏金池统计数据
const calculateStats = (records) => {
  let totalIn = 0
  let totalOut = 0
  records.forEach(record => {
    if (record.changeType === 'pool_add') {
      totalIn += record.amount
    } else if (record.changeType === 'pool_reduce') {
      totalOut += Math.abs(record.amount)
    }
  })
  poolStats.value = { totalIn, totalOut }
}

// 获取记录列表
const fetchRecordList = async () => {
  loading.value = true
  try {
    const params = {
      page: page.value,
      pageSize: pageSize.value
    }
    const res = await getRecordList(params)
    if (res.code === 0) {
      // 显示所有赏金变更记录
      let allRecords = res.data.list || []
      
      // 按变动类型筛选
      if (searchForm.changeType) {
        allRecords = allRecords.filter(r => r.changeType === searchForm.changeType)
      }
      
      // 按选手姓名筛选
      if (searchForm.playerName) {
        allRecords = allRecords.filter(r => 
          r.playerName?.includes(searchForm.playerName)
        )
      }
      
      // 按战队名称筛选
      if (searchForm.teamName) {
        allRecords = allRecords.filter(r => 
          r.teamName?.includes(searchForm.teamName)
        )
      }
      
      // 按时间范围筛选
      if (dateRange.value && dateRange.value.length === 2) {
        const [startDate, endDate] = dateRange.value
        allRecords = allRecords.filter(r => {
          const recordDate = r.createdAt?.split('T')[0]
          return recordDate >= startDate && recordDate <= endDate
        })
      }
      
      tableData.value = allRecords
      
      // 计算赏金池统计
      if (!searchForm.changeType && !searchForm.playerName && !searchForm.teamName && dateRange.value.length === 0) {
        calculateStats(res.data.list || [])
      }
      
      total.value = allRecords.length
    }
  } catch (error) {
    console.error('获取记录列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  page.value = 1
  fetchRecordList()
}

// 重置
const handleReset = () => {
  searchForm.changeType = ''
  searchForm.playerName = ''
  searchForm.teamName = ''
  dateRange.value = []
  page.value = 1
  fetchRecordList()
}

// 日期变化
const handleDateChange = () => {
  handleSearch()
}

// 分页变化
const handlePageChange = () => {
  fetchRecordList()
}

const handleSizeChange = () => {
  page.value = 1
  fetchRecordList()
}

onMounted(() => {
  fetchPoolInfo()
  fetchRecordList()
})
</script>

<style scoped>
.bounty-pool-container {
  padding: 20px;
}

.pool-overview {
  margin-bottom: 20px;
}

.section-wrapper {
  margin-bottom: 20px;
}

.stat-card {
  display: flex;
  align-items: center;
  padding: 20px;
  border-radius: 8px;
  background: #fff;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  margin-bottom: 10px;
}

.stat-icon {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 20px;
  font-size: 28px;
}

.pool-total .stat-icon {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
}

.pool-in .stat-icon {
  background: linear-gradient(135deg, #11998e 0%, #38ef7d 100%);
  color: #fff;
}

.pool-out .stat-icon {
  background: linear-gradient(135deg, #eb3349 0%, #f45c43 100%);
  color: #fff;
}

.pool-total {
  border-left: 4px solid #667eea;
}

.pool-in {
  border-left: 4px solid #11998e;
}

.pool-out {
  border-left: 4px solid #eb3349;
}

.stat-content {
  flex: 1;
}

.stat-label {
  font-size: 14px;
  color: #909399;
  margin-bottom: 8px;
}

.stat-value {
  font-size: 28px;
  font-weight: bold;
  color: #303133;
}

.filter-section {
  background: #fff;
  padding: 20px;
  border-radius: 8px;
  margin-bottom: 20px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.record-section {
  background: #fff;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
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
}

.record-count {
  font-size: 14px;
  color: #909399;
}

.amount-positive {
  color: #67c23a;
  font-weight: bold;
}

.amount-negative {
  color: #f56c6c;
  font-weight: bold;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

/* 响应式适配 */
@media (max-width: 768px) {
  .bounty-pool-container {
    padding: 10px;
  }

  .stat-card {
    padding: 15px;
  }

  .stat-icon {
    width: 50px;
    height: 50px;
    font-size: 24px;
  }

  .stat-value {
    font-size: 22px;
  }

  .filter-section :deep(.el-form) {
    display: block;
  }

  .filter-section :deep(.el-form-item) {
    display: block;
    margin-bottom: 10px;
  }
}
</style>