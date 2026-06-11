<template>
  <div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="openDrawer">新增战队</el-button>
      </div>
      <el-table
        ref="multipleTable"
        :data="tableData"
        style="width: 100%"
        tooltip-effect="dark"
        row-key="ID"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column align="left" label="创建时间" width="180">
          <template #default="scope">
            <span>{{ formatDate(scope.row.CreatedAt) }}</span>
          </template>
        </el-table-column>
        <el-table-column
          align="left"
          label="战队名称"
          prop="teamName"
          width="180"
        />
        <el-table-column
          align="left"
          label="总赏金"
          prop="totalBounty"
          width="120"
        />
        <el-table-column align="left" label="操作" min-width="200">
          <template #default="scope">
            <el-button
              type="primary"
              link
              icon="edit"
              @click="updateTeam(scope.row)"
              >编辑</el-button
            >
            <el-button
              type="primary"
              link
              icon="delete"
              @click="deleteTeamHandler(scope.row)"
              >删除</el-button
            >
            <el-button
              type="primary"
              link
              icon="wallet"
              @click="setBounty(scope.row)"
              >设置赏金</el-button
            >
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
    <el-drawer
      v-model="drawerFormVisible"
      :before-close="closeDrawer"
      :show-close="false"
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ type === 'create' ? '新增战队' : '编辑战队' }}</span>
          <div>
            <el-button @click="closeDrawer">取 消</el-button>
            <el-button type="primary" @click="enterDrawer">确 定</el-button>
          </div>
        </div>
      </template>
      <el-form :inline="false" :model="form" label-width="100px">
        <el-form-item label="战队名称">
          <el-input v-model="form.teamName" autocomplete="off" />
        </el-form-item>
        <el-form-item label="初始赏金">
          <el-input v-model.number="form.totalBounty" type="number" autocomplete="off" />
        </el-form-item>
      </el-form>
    </el-drawer>
    <el-drawer
      v-model="bountyDrawerVisible"
      :before-close="closeBountyDrawer"
      :show-close="false"
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">设置战队赏金</span>
          <div>
            <el-button @click="closeBountyDrawer">取 消</el-button>
            <el-button type="primary" @click="saveBounty">确 定</el-button>
          </div>
        </div>
      </template>
      <el-form :inline="false" :model="bountyForm" label-width="100px">
        <el-form-item label="战队名称">
          <el-input :value="bountyForm.teamName" disabled />
        </el-form-item>
        <el-form-item label="当前赏金">
          <el-input :value="bountyForm.currentBounty" disabled />
        </el-form-item>
        <el-form-item label="新赏金金额">
          <el-input v-model.number="bountyForm.newBounty" type="number" autocomplete="off" />
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<script setup>
import {
  createTeam,
  updateTeamApi,
  deleteTeam,
  getTeam,
  getTeamList,
  setTeamBounty
} from '@/api/team'
import { ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { formatDate } from '@/utils/format'

defineOptions({
  name: 'Team'
})

const form = ref({
  teamName: '',
  totalBounty: 0
})

const bountyForm = ref({
  teamId: 0,
  teamName: '',
  currentBounty: 0,
  newBounty: 0
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
  const table = await getTeamList({
    page: page.value,
    pageSize: pageSize.value
  })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

getTableData()

const drawerFormVisible = ref(false)
const bountyDrawerVisible = ref(false)
const type = ref('')

const updateTeam = async (row) => {
  const res = await getTeam({ ID: row.ID })
  type.value = 'update'
  if (res.code === 0) {
    form.value = {
      ID: res.data.ID,
      teamName: res.data.teamName,
      totalBounty: res.data.totalBounty
    }
    drawerFormVisible.value = true
  }
}

const closeDrawer = () => {
  drawerFormVisible.value = false
  form.value = {
    teamName: '',
    totalBounty: 0
  }
}

const closeBountyDrawer = () => {
  bountyDrawerVisible.value = false
  bountyForm.value = {
    teamId: 0,
    teamName: '',
    currentBounty: 0,
    newBounty: 0
  }
}

const deleteTeamHandler = async (row) => {
  ElMessageBox.confirm('确定要删除吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    const res = await deleteTeam({ ID: row.ID })
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: '删除成功'
      })
      if (tableData.value.length === 1 && page.value > 1) {
        page.value--
      }
      getTableData()
    }
  })
}

const enterDrawer = async () => {
  let res
  switch (type.value) {
    case 'create':
      res = await createTeam(form.value)
      break
    case 'update':
      res = await updateTeamApi(form.value)
      break
    default:
      res = await createTeam(form.value)
      break
  }

  if (res.code === 0) {
    closeDrawer()
    getTableData()
  }
}

const openDrawer = () => {
  type.value = 'create'
  drawerFormVisible.value = true
}

const setBounty = (row) => {
  bountyForm.value = {
    teamId: row.ID,
    teamName: row.teamName,
    currentBounty: row.totalBounty,
    newBounty: row.totalBounty
  }
  bountyDrawerVisible.value = true
}

const saveBounty = async () => {
  const res = await setTeamBounty({
    teamId: bountyForm.value.teamId,
    bounty: bountyForm.value.newBounty
  })
  if (res.code === 0) {
    closeBountyDrawer()
    getTableData()
    ElMessage({
      type: 'success',
      message: '设置成功'
    })
  }
}
</script>

<style></style>