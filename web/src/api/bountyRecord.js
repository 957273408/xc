import service from '@/utils/request'

export const getRecordList = (params) => {
  return service({
    url: '/bountyRecord/recordList',
    method: 'get',
    params
  })
}

export const getPoolInfo = () => {
  return service({
    url: '/bountyRecord/poolInfo',
    method: 'get'
  })
}

/**
 * 获取队伍赏金排行榜
 * @param {Object} params - 查询参数
 * @param {number} params.page - 页码
 * @param {number} params.pageSize - 每页数量
 * @returns {Promise}
 */
export const getTeamBountyRanking = (params) => {
  return service({
    url: '/bountyRecord/teamRanking',
    method: 'get',
    params
  })
}

/**
 * 获取选手赏金排行榜
 * @param {Object} params - 查询参数
 * @param {number} params.page - 页码
 * @param {number} params.pageSize - 每页数量
 * @returns {Promise}
 */
export const getPlayerBountyRanking = (params) => {
  return service({
    url: '/bountyRecord/playerRanking',
    method: 'get',
    params
  })
}