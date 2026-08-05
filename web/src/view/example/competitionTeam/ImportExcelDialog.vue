<template>
  <el-dialog
    v-model="visible"
    title="导入战队数据"
    width="500px"
    :close-on-click-modal="false"
  >
    <div class="import-dialog">
      <div class="mode-section">
        <span class="label">导入模式：</span>
        <el-radio-group v-model="mode">
          <el-radio value="incremental">增量导入</el-radio>
          <el-radio value="full">全量更新</el-radio>
        </el-radio-group>
        <el-tooltip
          :content="mode === 'incremental' ? '已存在的战队标识将被跳过' : '已存在的战队标识将被更新'"
          placement="top"
        >
          <el-icon class="info-icon"><QuestionFilled /></el-icon>
        </el-tooltip>
      </div>

      <div class="upload-section">
        <el-upload
          ref="uploadRef"
          :action="uploadUrl"
          :headers="uploadHeaders"
          :show-file-list="true"
          :limit="1"
          accept=".xlsx,.xls"
          :auto-upload="false"
          @change="handleFileChange"
          @exceed="handleExceed"
        >
          <el-button type="primary" icon="upload">选择Excel文件</el-button>
          <template #tip>
            <div class="upload-tip">
              支持 .xlsx 和 .xls 格式，文件大小不超过 10MB
            </div>
          </template>
        </el-upload>
      </div>

      <div class="preview-section" v-if="filePreview">
        <h4>文件预览</h4>
        <div class="file-info">
          <el-icon><Document /></el-icon>
          <span>{{ filePreview.name }}</span>
          <span class="file-size">{{ formatFileSize(filePreview.size) }}</span>
        </div>
      </div>

      <div class="result-section" v-if="result">
        <el-alert
          :title="resultMessage"
          :type="hasErrors ? 'warning' : 'success'"
          show-icon
          :closable="false"
        >
          <template #default>
            <div v-if="hasErrors" class="error-list">
              <div v-for="(error, index) in result.errors" :key="index">
                {{ error }}
              </div>
            </div>
          </template>
        </el-alert>
      </div>

      <div class="actions">
        <el-button @click="handleCancel">取消</el-button>
        <el-button
          type="primary"
          :loading="uploading"
          :disabled="!filePreview"
          @click="handleUpload"
        >
          {{ uploading ? '导入中...' : '开始导入' }}
        </el-button>
      </div>
    </div>
  </el-dialog>
</template>

<script setup>
import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { QuestionFilled, Document } from '@element-plus/icons-vue'
import { importExcel } from '@/api/competitionTeam'
import { useUserStore } from '@/pinia'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:modelValue', 'success'])

const userStore = useUserStore()

const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const mode = ref('incremental')
const filePreview = ref(null)
const uploading = ref(false)
const result = ref(null)

const uploadUrl = computed(() => {
  const baseUrl = import.meta.env.VITE_BASE_API
  const prefix = baseUrl === '/' ? '' : baseUrl
  return `${prefix}/competitionTeam/importExcel?mode=${mode.value}`
})

const uploadHeaders = computed(() => ({
  'x-token': userStore.token
}))

const hasErrors = computed(() => {
  return result.value && result.value.failCount > 0
})

const resultMessage = computed(() => {
  if (!result.value) return ''
  const { successCount, failCount } = result.value
  if (failCount === 0) {
    return `导入成功！共导入 ${successCount} 条数据`
  }
  return `导入完成：成功 ${successCount} 条，失败 ${failCount} 条`
})

const handleFileChange = (file) => {
  result.value = null
  
  if (!file.raw) {
    filePreview.value = null
    return
  }

  const validTypes = ['application/vnd.openxmlformats-officedocument.spreadsheetml.sheet', 'application/vnd.ms-excel']
  const validExts = ['.xlsx', '.xls']
  const ext = file.name.substring(file.name.lastIndexOf('.')).toLowerCase()

  if (!validTypes.includes(file.type) && !validExts.includes(ext)) {
    ElMessage.error('请上传Excel文件（.xlsx 或 .xls）')
    filePreview.value = null
    return
  }

  if (file.size > 10 * 1024 * 1024) {
    ElMessage.error('文件大小不能超过10MB')
    filePreview.value = null
    return
  }

  filePreview.value = file
}

const handleExceed = () => {
  ElMessage.warning('只能上传一个文件，请先移除已选文件')
}

const handleUpload = async () => {
  if (!filePreview.value) {
    ElMessage.warning('请先选择文件')
    return
  }

  uploading.value = true
  result.value = null

  try {
    const res = await importExcel(filePreview.value.raw, mode.value)
    if (res.code === 0) {
      result.value = res.data
      ElMessage.success(res.msg || '导入完成')
      emit('success', res.data)
    } else {
      result.value = { successCount: 0, failCount: 1, errors: [res.msg || '导入失败'] }
      ElMessage.error(res.msg || '导入失败')
    }
  } catch (error) {
    result.value = { successCount: 0, failCount: 1, errors: ['网络请求失败'] }
    ElMessage.error('网络请求失败')
  } finally {
    uploading.value = false
  }
}

const handleCancel = () => {
  visible.value = false
  resetState()
}

const resetState = () => {
  mode.value = 'incremental'
  filePreview.value = null
  result.value = null
}

const formatFileSize = (bytes) => {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / 1024 / 1024).toFixed(1) + ' MB'
}
</script>

<style scoped>
.import-dialog {
  padding: 10px 0;
}

.mode-section {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 20px;
}

.mode-section .label {
  font-weight: bold;
}

.info-icon {
  color: #909399;
  cursor: help;
}

.upload-section {
  margin-bottom: 20px;
}

.upload-tip {
  margin-top: 8px;
  font-size: 12px;
  color: #909399;
}

.preview-section {
  background: #f5f7fa;
  padding: 12px;
  border-radius: 8px;
  margin-bottom: 16px;
}

.preview-section h4 {
  margin: 0 0 8px 0;
  font-size: 14px;
}

.file-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.file-size {
  color: #909399;
  font-size: 12px;
}

.result-section {
  margin-bottom: 16px;
}

.error-list {
  margin-top: 8px;
  font-size: 12px;
}

.error-list div {
  margin-bottom: 4px;
  color: #e6a23c;
}

.actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style>
