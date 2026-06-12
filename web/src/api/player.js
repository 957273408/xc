import service from '@/utils/request'

export const createPlayer = (data) => {
  return service({
    url: '/player/player',
    method: 'post',
    data
  })
}

export const updatePlayer = (data) => {
  return service({
    url: '/player/player',
    method: 'put',
    data
  })
}

export const deletePlayer = (data) => {
  return service({
    url: '/player/player',
    method: 'delete',
    data
  })
}

export const getPlayer = (params) => {
  return service({
    url: '/player/player',
    method: 'get',
    params
  })
}

export const getPlayerList = (params) => {
  return service({
    url: '/player/playerList',
    method: 'get',
    params
  })
}

export const allocateBounty = (data) => {
  return service({
    url: '/player/allocateBounty',
    method: 'post',
    data
  })
}

export const kill = (data) => {
  return service({
    url: '/player/kill',
    method: 'post',
    data
  })
}

export const revive = (data) => {
  return service({
    url: '/player/revive',
    method: 'post',
    data
  })
}

export const claimFromPool = (data) => {
  return service({
    url: '/player/claimFromPool',
    method: 'post',
    data
  })
}