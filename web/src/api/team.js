import service from '@/utils/request'

export const createTeam = (data) => {
  return service({
    url: '/team/team',
    method: 'post',
    data
  })
}

export const updateTeamApi = (data) => {
  return service({
    url: '/team/team',
    method: 'put',
    data
  })
}

export const deleteTeam = (data) => {
  return service({
    url: '/team/team',
    method: 'delete',
    data
  })
}

export const getTeam = (params) => {
  return service({
    url: '/team/team',
    method: 'get',
    params
  })
}

export const getTeamList = (params) => {
  return service({
    url: '/team/teamList',
    method: 'get',
    params
  })
}

export const setTeamBounty = (data) => {
  return service({
    url: '/team/setBounty',
    method: 'post',
    data
  })
}