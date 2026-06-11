<template>
  <div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="openDrawer">新增选手</el-button>
        <el-button type="success" icon="share" @click="openAllocateDrawer">赏金分配</el-button>
      </div>
      <el-form :inline="true" :model="searchForm" class="gva-search-form">
        <el-form-item label="战队筛选">
          <el-select v-model="searchForm.teamId" placeholder="请选择战队">
            <el-option :value="0" label="全部战队" />
            <el-option
              v-for="team in teamList"
              :key="team.ID"
              :label="team.teamName"
              :value="team.ID"
            />
          </el-select>
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
        <el-table-column type="selection" width="55" />
        <el-table-column align="left" label="创建时间" width="180">
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
          label="UID"
          prop="uid"
          width="150"
        />
        <el-table-column
          align="left"
          label="所属战队"
          width="150"
        >
          <template #default="scope">
            <span>{{ scope.row.Team?.teamName || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column
          align="left"
          label="当前赏金"
          prop="bounty"
          width="120"
        />
        <el-table-column align="left" label="操作" min-width="200">
          <template #default="scope">
            <el-button
              type="primary"
              link
              icon="edit"
              @click="updatePlayer(scope.row)"
              >编辑</el-button
            >
            <el-button
              type="primary"
              link
              icon="delete"
              @click="deletePlayer(scope.row)"
              >删除</el-button
            >
            <el-button
              type="danger"
              link
              icon="swords"
              @click="handleKill(scope.row)"
              >击杀</el-button
            >
            <el-button
              type="success"
              link
              icon="refresh"
              @click="handleRevive(scope.row)"
              >复活</el-button
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
          <span class="text-lg">{{ type === 'create' ? '新增选手' : '编辑选手' }}</span>
          <div>
            <el-button @click="closeDrawer">取 消</el-button>
            <el-button type="primary" @click="enterDrawer">确 定</el-button>
          </div>
        </div>
      </template>
      <el-form :inline="false" :model="form" label-width="100px">
        <el-form-item label="选手姓名">
          <el-input v-model="form.playerName" autocomplete="off" />
        </el-form-item>
        <el-form-item label="UID">
          <el-input v-model="form.uid" autocomplete="off" />
        </el-form-item>
        <el-form-item label="所属战队">
          <el-select v-model="form.teamId" placeholder="请选择战队">
            <el-option
              v-for="team in teamList"
              :key="team.ID"
              :label="team.teamName"
              :value="team.ID"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="初始赏金">
          <el-input v-model.number="form.bounty" type="number" autocomplete="off" />
        </el-form-item>
      </el-form>
    </el-drawer>
    <el-drawer
      v-model="allocateDrawerVisible"
      :before-close="closeAllocateDrawer"
      :show-close="false"
      size="800px"
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">战队赏金分配</span>
          <div>
            <el-button @click="closeAllocateDrawer">取 消</el-button>
            <el-button type="primary" @click="saveAllocate">确 定</el-button>
          </div>
        </div>
      </template>
      <el-form :inline="false" :model="allocateForm" label-width="100px">
        <el-form-item label="选择战队">
          <el-select v-model="allocateForm.teamId" placeholder="请选择战队" @change="onTeamChange">
            <el-option
              v-for="team in teamList"
              :key="team.ID"
              :label="team.teamName"
              :value="team.ID"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="战队总赏金">
          <el-input :value="allocateForm.totalBounty" disabled />
        </el-form-item>
        <el-form-item label="已分配金额">
          <el-input :value="allocateForm.allocatedAmount" disabled />
        </el-form-item>
        <el-form-item label="剩余金额">
          <el-input :value="allocateForm.remainingAmount" disabled />
        </el-form-item>
      </el-form>
      <el-table
        v-if="allocateForm.teamId"
        :data="allocatePlayers"
        style="width: 100%"
        row-key="ID"
      >
        <el-table-column align="left" label="选手姓名" prop="playerName" />
        <el-table-column align="left" label="当前赏金" prop="bounty" />
        <el-table-column align="left" label="分配金额">
          <template #default="scope">
            <el-input
              v-model.number="scope.row.allocateAmount"
              type="number"
              @change="calculateAllocated"
            />
          </template>
        </el-table-column>
      </el-table>
    </el-drawer>
    <el-drawer
      v-model="killDrawerVisible"
      :before-close="closeKillDrawer"
      :show-close="false"
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">击杀操作</span>
          <div>
            <el-button @click="closeKillDrawer">取 消</el-button>
            <el-button type="danger" @click="saveKill">确 认击杀</el-button>
          </div>
        </div>
      </template>
      <el-form :inline="false" :model="killForm" label-width="100px">
        <el-form-item label="击杀者">
          <el-select v-model="killForm.killerId" placeholder="请选择击杀者">
            <el-option
              v-for="player in otherPlayers"
              :key="player.ID"
              :label="player.playerName"
              :value="player.ID"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="被击杀者">
          <el-input :value="killForm.victimName" disabled />
        </el-form-item>
        <el-form-item label="被击杀者赏金">
          <el-input :value="killForm.victimBounty" disabled />
        </el-form-item>
        <el-form-item label="夺取赏金(50%)">
          <el-input :value="killForm.stealAmount" disabled />
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<script setup>
import {
  createPlayer,
  updatePlayer as updatePlayerApi,
  deletePlayer as deletePlayerApi,
  getPlayer,
  getPlayerList,
  allocateBounty,
  kill,
  revive
} from '@/api/player'
import { getTeamList } from '@/api/team'
import { ref, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { formatDate } from '@/utils/format'

defineOptions({
  name: 'Player'
})

const form = ref({
  playerName: '',
  uid: '',
  teamId: 0,
  bounty: 0
})

const searchForm = ref({
  teamId: 0
})

const allocateForm = ref({
  teamId: 0,
  totalBounty: 0,
  allocatedAmount: 0,
  remainingAmount: 0
})

const killForm = ref({
  killerId: 0,
  victimId: 0,
  victimName: '',
  victimBounty: 0,
  stealAmount: 0
})

const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const teamList = ref([])
const allocatePlayers = ref([])
const currentVictim = ref(null)

const otherPlayers = computed(() => {
  if (!currentVictim.value) return []
  return tableData.value.filter(p => p.ID !== currentVictim.value.ID)
})

const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

const getTableData = async () => {
  const table = await getPlayerList({
    page: page.value,
    pageSize: pageSize.value,
    teamId: searchForm.value.teamId
  })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

const loadTeams = async () => {
  const res = await getTeamList({ page: 1, pageSize: 100 })
  if (res.code === 0) {
    teamList.value = res.data.list
  }
}

getTableData()
loadTeams()

const drawerFormVisible = ref(false)
const allocateDrawerVisible = ref(false)
const killDrawerVisible = ref(false)
const type = ref('')

const updatePlayer = async (row) => {
  const res = await getPlayer({ ID: row.ID })
  type.value = 'update'
  if (res.code === 0) {
    form.value = {
      ID: res.data.ID,
      playerName: res.data.playerName,
      uid: res.data.uid,
      teamId: res.data.TeamID,
      bounty: res.data.bounty
    }
    drawerFormVisible.value = true
  }
}

const closeDrawer = () => {
  drawerFormVisible.value = false
  form.value = {
    playerName: '',
    uid: '',
    teamId: 0,
    bounty: 0
  }
}

const closeAllocateDrawer = () => {
  allocateDrawerVisible.value = false
  allocateForm.value = {
    teamId: 0,
    totalBounty: 0,
    allocatedAmount: 0,
    remainingAmount: 0
  }
  allocatePlayers.value = []
}

const closeKillDrawer = () => {
  killDrawerVisible.value = false
  killForm.value = {
    killerId: 0,
    victimId: 0,
    victimName: '',
    victimBounty: 0,
    stealAmount: 0
  }
  currentVictim.value = null
}

const deletePlayer = async (row) => {
  ElMessageBox.confirm('确定要删除吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    const res = await deletePlayerApi({ ID: row.ID })
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
      res = await createPlayer(form.value)
      break
    case 'update':
      res = await updatePlayerApi(form.value)
      break
    default:
      res = await createPlayer(form.value)
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

const openAllocateDrawer = () => {
  allocateDrawerVisible.value = true
}

const onTeamChange = async () => {
  if (!allocateForm.value.teamId) {
    allocatePlayers.value = []
    allocateForm.value.totalBounty = 0
    return
  }
  const team = teamList.value.find(t => t.ID === allocateForm.value.teamId)
  if (team) {
    allocateForm.value.totalBounty = team.totalBounty
  }
  const res = await getPlayerList({ page: 1, pageSize: 100, teamId: allocateForm.value.teamId })
  if (res.code === 0) {
    allocatePlayers.value = res.data.list.map(p => ({
      ...p,
      allocateAmount: p.bounty
    }))
    calculateAllocated()
  }
}

const calculateAllocated = () => {
  allocateForm.value.allocatedAmount = allocatePlayers.value.reduce((sum, p) => sum + (p.allocateAmount || 0), 0)
  allocateForm.value.remainingAmount = allocateForm.value.totalBounty - allocateForm.value.allocatedAmount
}

const saveAllocate = async () => {
  const playerBounties = allocatePlayers.value.map(p => ({
    playerId: p.ID,
    amount: p.allocateAmount || 0
  }))
  const res = await allocateBounty({
    teamId: allocateForm.value.teamId,
    playerBounties
  })
  if (res.code === 0) {
    closeAllocateDrawer()
    getTableData()
    ElMessage({
      type: 'success',
      message: '分配成功'
    })
  }
}

const handleKill = (row) => {
  currentVictim.value = row
  killForm.value = {
    killerId: 0,
    victimId: row.ID,
    victimName: row.playerName,
    victimBounty: row.bounty,
    stealAmount: row.bounty * 0.5
  }
  killDrawerVisible.value = true
}

const saveKill = async () => {
  if (!killForm.value.killerId) {
    ElMessage({
      type: 'warning',
      message: '请选择击杀者'
    })
    return
  }
  const res = await kill({
    killerId: killForm.value.killerId,
    victimId: killForm.value.victimId
  })
  if (res.code === 0) {
    closeKillDrawer()
    getTableData()
    ElMessage({
      type: 'success',
      message: `击杀成功，夺取赏金: ${res.data.amount}`
    })
  }
}

const handleRevive = async (row) => {
  ElMessageBox.confirm(`确定让 ${row.playerName} 复活吗？复活将损失50%赏金`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    const res = await revive({ playerId: row.ID })
    if (res.code === 0) {
      getTableData()
      ElMessage({
        type: 'success',
        message: `复活成功，损失赏金: ${res.data.lostAmount}`
      })
    }
  })
}
</script>

<style></style>