<template>
  <div class="competition-team-page">
    <div class="page-header">
      <div class="header-content">
        <h1 class="page-title">
          <span class="title-icon">🏆</span>
          战队信息管理
        </h1>
        <p class="page-subtitle">第八届锦标赛 · 战队数据管理系统</p>
      </div>
      <div class="header-actions">
        <el-button type="primary" icon="Plus" @click="openCreateDialog">
          新增战队
        </el-button>
        <el-button type="success" icon="Upload" @click="showImportDialog = true">
          Excel导入
        </el-button>
        <el-button icon="Refresh" @click="fetchAllData">
          刷新
        </el-button>
      </div>
    </div>

    <div class="stats-overview">
      <div class="stat-card">
        <div class="stat-icon teams-icon">👥</div>
        <div class="stat-content">
          <span class="stat-value">{{ totalTeams }}</span>
          <span class="stat-label">战队总数</span>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-card-content">
          <span class="stat-value">{{ totalMatches }}</span>
          <span class="stat-label">比赛场次</span>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-card-content">
          <span class="stat-value">{{ totalScore }}</span>
          <span class="stat-label">总积分</span>
        </div>
      </div>
    </div>

    <div class="main-content">
      <div class="teams-section">
        <div class="section-header">
          <h2>战队列表</h2>
          <div class="search-box">
            <el-input
              v-model="searchKeyword"
              placeholder="搜索战队名称或标识"
              clearable
              @keyup.enter="handleSearch"
            >
              <template #prefix>
                <el-icon><Search /></el-icon>
              </template>
            </el-input>
          </div>
        </div>

        <div v-loading="loading" class="teams-grid">
          <el-empty
            v-if="!loading && teamList.length === 0"
            description="暂无战队数据，请点击新增或Excel导入"
          />
          <TeamCard
            v-for="team in teamList"
            :key="team.ID"
            :team="team"
            @click="handleTeamClick"
            @edit="openEditDialog"
            @delete="handleDeleteTeam"
            @manageWarId="openWarIDManager"
          />
        </div>

        <div class="pagination-container" v-if="teamList.length > 0">
          <el-pagination
            v-model:current-page="currentPage"
            v-model:page-size="pageSize"
            :page-sizes="[12, 24, 48]"
            :total="total"
            layout="total, sizes, prev, pager, next"
            @current-change="handlePageChange"
            @size-change="handleSizeChange"
          />
        </div>
      </div>

      <div class="detail-section" v-if="selectedTeam">
        <div class="detail-header">
          <h2>战队详情</h2>
          <el-button size="small" @click="selectedTeam = null">关闭</el-button>
        </div>
        <ScoreChart
          ref="scoreChartRef"
          :scores="teamScores"
          :total-score="selectedTeam.totalScore"
          :match-count="teamScores.length"
          :last-rank="teamScores[0]?.rank || 0"
          :team-id="selectedTeam.ID"
          @edit-score="handleEditScore"
          @delete-score="handleDeleteScore"
          @refresh="fetchTeamDetail"
        />
      </div>
    </div>

    <el-dialog
      v-model="createDialogVisible"
      :title="editingTeam ? '编辑战队' : '新增战队'"
      width="500px"
    >
      <el-form :model="teamForm" :rules="formRules" ref="formRef" label-width="100px">
        <el-form-item label="战队标识" prop="teamCode">
          <el-input
            v-model="teamForm.teamCode"
            placeholder="3-8位字母或数字"
            maxlength="8"
          />
          <div class="form-tip">战队唯一标识，3-8位字母或数字组合</div>
        </el-form-item>
        <el-form-item label="战队名称" prop="teamName">
          <el-input
            v-model="teamForm.teamName"
            placeholder="请输入战队名称"
            maxlength="50"
          />
        </el-form-item>
        <el-form-item label="战队Logo">
          <el-input
            v-model="teamForm.teamLogo"
            placeholder="请输入Logo URL或上传图片"
          />
          <div class="upload-logo-btn">
            <el-upload
              :auto-upload="false"
              :show-file-list="false"
              :before-upload="handleLogoUpload"
              accept="image/jpeg,image/png"
            >
              <el-button size="small">上传Logo</el-button>
            </el-upload>
            <span class="form-tip">支持JPG/PNG格式，最大2MB</span>
          </div>
        </el-form-item>
        <el-form-item label="分组信息">
          <el-input
            v-model="teamForm.groupName"
            placeholder="请输入分组信息"
            maxlength="50"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmitTeam">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog
      v-model="warIDDialogVisible"
      title="WarId 管理"
      width="600px"
      :close-on-click-modal="false"
    >
      <WarIDManager
        v-if="warIDDialogVisible && selectedTeam"
        :team="selectedTeam"
        @close="warIDDialogVisible = false"
        @score-updated="fetchTeamDetail"
      />
    </el-dialog>

    <ImportExcelDialog
      v-model="showImportDialog"
      @success="fetchAllData"
    />
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Plus, Upload, Refresh } from '@element-plus/icons-vue'
import {
  createCompetitionTeam,
  updateCompetitionTeam,
  deleteCompetitionTeam,
  getCompetitionTeamList,
  getTeamRecentScores
} from '@/api/competitionTeam'
import TeamCard from './TeamCard.vue'
import WarIDManager from './WarIDManager.vue'
import ScoreChart from './ScoreChart.vue'
import ImportExcelDialog from './ImportExcelDialog.vue'

defineOptions({
  name: 'CompetitionTeam'
})

const loading = ref(false)
const teamList = ref([])
const currentPage = ref(1)
const pageSize = ref(12)
const total = ref(0)
const searchKeyword = ref('')

const selectedTeam = ref(null)
const teamScores = ref([])
const scoreChartRef = ref(null)

const createDialogVisible = ref(false)
const editingTeam = ref(null)
const teamForm = reactive({
  teamCode: '',
  teamName: '',
  teamLogo: '',
  groupName: ''
})

const formRef = ref(null)
const formRules = {
  teamCode: [
    { required: true, message: '请输入战队标识', trigger: 'blur' }
  ],
  teamName: [
    { required: true, message: '请输入战队名称', trigger: 'blur' },
    { max: 50, message: '最长50字符', trigger: 'blur' }
  ]
}

const showImportDialog = ref(false)
const warIDDialogVisible = ref(false)

const totalTeams = computed(() => total.value)
const totalMatches = computed(() => {
  return teamList.value.reduce((sum, t) => sum + (t.matchCount || 0), 0)
})
const totalScore = computed(() => {
  return teamList.value.reduce((sum, t) => sum + (t.totalScore || 0), 0)
})

const fetchTeamList = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      pageSize: pageSize.value
    }
    if (searchKeyword.value) {
      params.keyword = searchKeyword.value
    }
    const res = await getCompetitionTeamList(params)
    if (res.code === 0) {
      teamList.value = res.data?.list || []
      total.value = res.data?.total || 0
    }
  } catch (error) {
    console.error('获取战队列表失败:', error)
    ElMessage.error('获取战队列表失败')
  } finally {
    loading.value = false
  }
}

const fetchTeamDetail = async () => {
  if (!selectedTeam.value) return
  try {
    const res = await getTeamRecentScores({
      teamId: selectedTeam.value.ID,
      limit: 50
    })
    if (res.code === 0) {
      teamScores.value = res.data || []
      const updatedTeam = teamList.value.find(t => t.ID === selectedTeam.value.ID)
      if (updatedTeam) {
        selectedTeam.value = { ...updatedTeam }
      }
    }
  } catch (error) {
    console.error('获取战队详情失败:', error)
  }
}

const fetchAllData = async () => {
  await fetchTeamList()
  if (selectedTeam.value) {
    await fetchTeamDetail()
  }
}

const handleEditScore = (score) => {
  scoreChartRef.value?.handleEditScore(score)
}

const handleDeleteScore = (score) => {
  scoreChartRef.value?.handleDeleteScore(score)
}

const handleSearch = () => {
  currentPage.value = 1
  fetchTeamList()
}

const handlePageChange = (page) => {
  currentPage.value = page
  fetchTeamList()
}

const handleSizeChange = (size) => {
  pageSize.value = size
  currentPage.value = 1
  fetchTeamList()
}

const handleTeamClick = (team) => {
  selectedTeam.value = team
  fetchTeamDetail()
}

const openCreateDialog = () => {
  editingTeam.value = null
  teamForm.teamCode = ''
  teamForm.teamName = ''
  teamForm.teamLogo = ''
  teamForm.groupName = ''
  createDialogVisible.value = true
}

const openEditDialog = (team) => {
  editingTeam.value = team
  teamForm.teamCode = team.teamCode
  teamForm.teamName = team.teamName
  teamForm.teamLogo = team.teamLogo
  teamForm.groupName = team.groupName || ''
  createDialogVisible.value = true
}

const handleSubmitTeam = async () => {
  try {
    await formRef.value.validate()
    
    if (editingTeam.value) {
      const res = await updateCompetitionTeam({
        id: editingTeam.value.ID,
        teamCode: teamForm.teamCode,
        teamName: teamForm.teamName,
        teamLogo: teamForm.teamLogo,
        groupName: teamForm.groupName
      })
      if (res.code === 0) {
        ElMessage.success('更新成功')
        createDialogVisible.value = false
        await fetchTeamList()
      } else {
        ElMessage.error(res.msg || '更新失败')
      }
    } else {
      const res = await createCompetitionTeam({
        teamCode: teamForm.teamCode,
        teamName: teamForm.teamName,
        teamLogo: teamForm.teamLogo,
        groupName: teamForm.groupName
      })
      if (res.code === 0) {
        ElMessage.success('创建成功')
        createDialogVisible.value = false
        await fetchTeamList()
      } else {
        ElMessage.error(res.msg || '创建失败')
      }
    }
  } catch (error) {
    // 表单验证失败
  }
}

const handleDeleteTeam = async (team) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除战队 "${team.teamName}" 吗？删除后该战队的所有积分记录也将被删除。`,
      '确认删除',
      { type: 'warning' }
    )
    const res = await deleteCompetitionTeam({ id: team.ID })
    if (res.code === 0) {
      ElMessage.success('删除成功')
      if (selectedTeam.value?.ID === team.ID) {
        selectedTeam.value = null
      }
      await fetchTeamList()
    } else {
      ElMessage.error(res.msg || '删除失败')
    }
  } catch (e) {
    // 用户取消删除
  }
}

const openWarIDManager = (team) => {
  selectedTeam.value = team
  warIDDialogVisible.value = true
}

const handleLogoUpload = (file) => {
  const isImage = file.type === 'image/jpeg' || file.type === 'image/png'
  const isLt2M = file.size / 1024 / 1024 < 2

  if (!isImage) {
    ElMessage.error('只支持JPG/PNG格式')
    return false
  }
  if (!isLt2M) {
    ElMessage.error('Logo大小不能超过2MB')
    return false
  }

  const reader = new FileReader()
  reader.onload = (e) => {
    teamForm.teamLogo = e.target.result
  }
  reader.readAsDataURL(file)
  return false
}

onMounted(() => {
  fetchTeamList()
})
</script>

<style scoped>
.competition-team-page {
  padding: 24px;
  min-height: 100vh;
  background: linear-gradient(135deg, #0f0f23 0%, #1a1a2e 50%, #16213e 100%);
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 30px;
  flex-wrap: wrap;
  gap: 20px;
}

.header-content {
  color: #fff;
}

.page-title {
  font-size: 28px;
  font-weight: bold;
  margin: 0 0 8px 0;
  display: flex;
  align-items: center;
  gap: 12px;
}

.title-icon {
  font-size: 32px;
}

.page-subtitle {
  color: rgba(255, 255, 255, 0.5);
  margin: 0;
}

.header-actions {
  display: flex;
  gap: 12px;
}

.stats-overview {
  display: flex;
  gap: 20px;
  margin-bottom: 30px;
}

.stat-card {
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.1) 0%, rgba(255, 255, 255, 0.05) 100%);
  border: 1px solid rgba(255, 215, 0, 0.2);
  border-radius: 12px;
  padding: 20px 24px;
  flex: 1;
  display: flex;
  align-items: center;
  gap: 16px;
  transition: all 0.3s;
}

.stat-card:hover {
  border-color: rgba(255, 215, 0, 0.5);
  transform: translateY(-2px);
}

.stat-icon {
  font-size: 36px;
}

.stat-content,
.stat-card-content {
  display: flex;
  flex-direction: column;
}

.stat-value {
  font-size: 28px;
  font-weight: bold;
  color: #ffd700;
}

.stat-label {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.6);
}

.main-content {
  display: grid;
  grid-template-columns: 1fr;
  gap: 24px;
}

.teams-section,
.detail-section {
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 16px;
  padding: 24px;
}

.section-header,
.detail-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.section-header h2,
.detail-header h2 {
  color: #fff;
  margin: 0;
}

.search-box {
  width: 280px;
}

.teams-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 16px;
}

.pagination-container {
  display: flex;
  justify-content: center;
  margin-top: 24px;
}

.form-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}

.upload-logo-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 8px;
}

@media (max-width: 768px) {
  .competition-team-page {
    padding: 16px;
  }

  .page-header {
    flex-direction: column;
  }

  .page-title {
    font-size: 22px;
  }

  .stats-overview {
    flex-wrap: wrap;
  }

  .stat-card {
    flex: 1 1 calc(50% - 10px);
  }

  .search-box {
    width: 100%;
  }

  .teams-grid {
    grid-template-columns: 1fr;
  }
}
</style>
