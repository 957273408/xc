<template>
  <div class="team-card" :class="{ 'is-expanded': expanded }" @click="$emit('click', team)">
    <div class="card-header">
      <div class="team-logo" :style="logoStyle">
        <span v-if="!team.teamLogo">{{ team.teamName?.charAt(0) || '?' }}</span>
        <img v-else :src="team.teamLogo" :alt="team.teamName" />
      </div>
      <div class="team-info">
        <h3 class="team-name">{{ team.teamName }}</h3>
        <span class="team-code">{{ team.teamCode }}</span>
      </div>
      <div class="team-score">
        <span class="score-value">{{ team.totalScore || 0 }}</span>
        <span class="score-label">总积分</span>
      </div>
    </div>
    <div class="card-actions" v-if="showActions">
      <el-button type="primary" size="small" @click.stop="$emit('edit', team)">
        编辑
      </el-button>
      <el-button type="danger" size="small" @click.stop="$emit('delete', team)">
        删除
      </el-button>
      <el-button type="success" size="small" @click.stop="$emit('manageWarId', team)">
        管理WarId
      </el-button>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  team: {
    type: Object,
    required: true
  },
  expanded: {
    type: Boolean,
    default: false
  },
  showActions: {
    type: Boolean,
    default: true
  }
})

defineEmits(['click', 'edit', 'delete', 'manageWarId'])

const logoStyle = computed(() => {
  if (!props.team.teamLogo) {
    return {}
  }
  return {
    backgroundImage: `url(${props.team.teamLogo})`,
    backgroundSize: 'cover',
    backgroundPosition: 'center'
  }
})
</script>

<style scoped>
.team-card {
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 100%);
  border-radius: 12px;
  padding: 20px;
  border: 1px solid rgba(255, 215, 0, 0.3);
  transition: all 0.3s ease;
  cursor: pointer;
  position: relative;
  overflow: hidden;
}

.team-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 3px;
  background: linear-gradient(90deg, #ffd700, #ff6b6b, #ffd700);
}

.team-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 40px rgba(255, 215, 0, 0.2);
  border-color: rgba(255, 215, 0, 0.6);
}

.card-header {
  display: flex;
  align-items: center;
  gap: 16px;
}

.team-logo {
  width: 64px;
  height: 64px;
  border-radius: 12px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  font-weight: bold;
  color: #fff;
  overflow: hidden;
  flex-shrink: 0;
}

.team-logo img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.team-info {
  flex: 1;
  min-width: 0;
}

.team-name {
  font-size: 18px;
  font-weight: bold;
  color: #fff;
  margin: 0 0 4px 0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.team-code {
  font-size: 12px;
  color: #ffd700;
  background: rgba(255, 215, 0, 0.1);
  padding: 2px 8px;
  border-radius: 4px;
  display: inline-block;
}

.team-score {
  text-align: center;
  padding: 0 12px;
}

.score-value {
  display: block;
  font-size: 28px;
  font-weight: bold;
  color: #ffd700;
  text-shadow: 0 0 20px rgba(255, 215, 0, 0.5);
}

.score-label {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.6);
}

.card-actions {
  display: flex;
  gap: 8px;
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.card-actions .el-button {
  flex: 1;
}

@media (max-width: 768px) {
  .team-card {
    padding: 16px;
  }

  .team-logo {
    width: 48px;
    height: 48px;
    font-size: 18px;
  }

  .team-name {
    font-size: 16px;
  }

  .score-value {
    font-size: 24px;
  }

  .card-actions {
    flex-wrap: wrap;
  }

  .card-actions .el-button {
    flex: 1 1 calc(50% - 4px);
  }
}
</style>
