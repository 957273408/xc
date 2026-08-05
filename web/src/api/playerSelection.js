import service from '@/utils/request'

// 获取指定WarId下的所有玩家列表
export const getWarPlayers = (params) => {
  return service({
    url: '/playerSelection/warPlayers',
    method: 'get',
    params
  })
}

// 保存玩家选择（5人 + 附加统计项）
export const savePlayerSelection = (data) => {
  return service({
    url: '/playerSelection/save',
    method: 'post',
    data
  })
}

// 按SessionKey获取已保存的玩家选择
export const getPlayerSelection = (params) => {
  return service({
    url: '/playerSelection/get',
    method: 'get',
    params
  })
}

// 获取最新保存的玩家选择数据（无需sessionKey）
export const getLatestSelection = () => {
  return service({
    url: '/playerSelection/latest',
    method: 'get'
  })
}

// 多场汇总：传入多个WarId获取汇总后的选手数据
export const getMultiWarPlayers = (data) => {
  return service({
    url: '/playerSelection/multiWarPlayers',
    method: 'post',
    data
  })
}
