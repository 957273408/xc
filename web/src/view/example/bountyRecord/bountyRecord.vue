<template>
  <div>
    <div class="gva-table-box">
      <el-form :inline="true" :model="searchForm" class="gva-search-form">
        <el-form-item label="选手ID">
          <el-input v-model.number="searchForm.playerId" type="number" placeholder="请输入选手ID" />
        </el-form-item>
        <el-form-item label="战队ID">
          <el-input v-model.number="searchForm.teamId" type="number" placeholder="请输入战队ID" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="getTableData">查询</el-button>
        </el-form-item>
      </el-form>
      <el-table
        ref="multipleTable"
        :data="tableData"
        style="width: 100%"
        tooltip-effect="dark"
        row-key="ID"
      >
        <el-table-column align="left" label="记录时间" width="180">
          <template #default="scope">
            <span>{{ formatDate(scope.row.CreatedAt) }}</span>
          </template>
        </el-table-column>
        <el-table-column
          align="left"
          label="选手姓名"
          prop="playerName"
          width="120"
        />
        <el-table-column
          align="left"
          label="所属战队"
          prop="teamName"
          width="150"
        />
        <el-table-column
          align="left"
          label="变动类型"
          width="120"
        >
          <template #default="scope">
            <el-tag :type="getChangeTypeTag(scope.row.changeType)">
              {{ getChangeTypeLabel(scope.row.changeType) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          align="left"
          label="变动金额"
          width="120"
        >
          <template #default="scope">
            <span :class="scope.row.amount >= 0 ? 'text-green' : 'text-red'">
              {{ scope.row.amount >= 0 ? '+' : '' }}{{ scope.row.amount }}
            </span>
          </template>
        </el-table-column>
        <el-table-column
          align="left"
          label="变动后余额"
          prop="balance"
          width="120"
        />
        <el-table-column
          align="left"
          label="变动原因"
          prop="reason"
          width="180"
        />
        <el-table-column
          align="left"
          label="关联对象"
          width="120"
        >
          <template #default="scope">
            <span>{{ scope.row.relatedName || '-' }}</span>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination
          :current-page="page"
          :page-size="pageSize"
          :page-sizes="[10, 30, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @current-change="handleCurrentChange"
          @size-change="handleSizeChange"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { getRecordList } from '@/api/bountyRecord'
import { ref } from 'vue'
import { formatDate } from '@/utils/format'

defineOptions({
  name: 'BountyRecord'
})

const searchForm = ref({
  playerId: 0,
  teamId: 0
})

const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])

const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

const getTableData = async () => {
  const table = await getRecordList({
    page: page.value,
    pageSize: pageSize.value,
    playerId: searchForm.value.playerId || undefined,
    teamId: searchForm.value.teamId || undefined
  })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

getTableData()

const getChangeTypeLabel = (type) => {
  const labels = {
    'allocate': '赏金分配',
    'kill': '击杀获取',
    'killed': '被击杀',
    'revive': '复活损失'
  }
  return labels[type] || type
}

const getChangeTypeTag = (type) => {
  const tags = {
    'allocate': 'success',
    'kill': 'success',
    'killed': 'danger',
    'revive': 'warning'
  }
  return tags[type] || 'info'
}
</script>

<style scoped>
.text-green {
  color: #67c23a;
}
.text-red {
  color: #f56c6c;
}
</style>