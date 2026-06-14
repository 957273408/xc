import service from '@/utils/request'

/**
 * 获取战队奖金信息
 * @description 获取当前战队的可用奖金总额
 * @returns {Promise}
 */
export const getTeamBonusInfo = () => {
  return service({
    url: '/bountyAllocation/teamBonus',
    method: 'get'
  })
}

/**
 * 提交战队奖金分配
 * @description 将战队奖金分配给选手
 * @param {Object} data - 分配数据
 * @param {string} data.teamId - 战队ID
 * @param {Array} data.players - 选手分配列表
 * @param {string} data.players[].playerId - 选手ID
 * @param {string} data.players[].name - 选手姓名
 * @param {number} data.players[].amount - 分配金额
 * @returns {Promise}
 */
export const submitAllocation = (data) => {
  return service({
    url: '/bountyAllocation/submit',
    method: 'post',
    data
  })
}

/**
 * 保存战队奖金分配草稿
 * @description 将当前分配方案保存为草稿
 * @param {Object} data - 分配数据
 * @param {string} data.teamId - 战队ID
 * @param {Array} data.players - 选手分配列表
 * @param {string} data.players[].playerId - 选手ID
 * @param {string} data.players[].name - 选手姓名
 * @param {number} data.players[].amount - 分配金额
 * @returns {Promise}
 */
export const saveDraft = (data) => {
  return service({
    url: '/bountyAllocation/saveDraft',
    method: 'post',
    data
  })
}