import service from '@/utils/request'

export const createCompetitionTeam = (data) => {
  return service({
    url: '/competitionTeam/competitionTeam',
    method: 'post',
    data
  })
}

export const updateCompetitionTeam = (data) => {
  return service({
    url: '/competitionTeam/competitionTeam',
    method: 'put',
    data
  })
}

export const deleteCompetitionTeam = (data) => {
  return service({
    url: '/competitionTeam/competitionTeam',
    method: 'delete',
    data
  })
}

export const getCompetitionTeam = (params) => {
  return service({
    url: '/competitionTeam/competitionTeam',
    method: 'get',
    params
  })
}

export const getCompetitionTeamList = (params) => {
  return service({
    url: '/competitionTeam/competitionTeamList',
    method: 'get',
    params
  })
}

export const importExcel = (file, mode) => {
  const formData = new FormData()
  formData.append('file', file)
  return service({
    url: `/competitionTeam/importExcel?mode=${mode}`,
    method: 'post',
    data: formData,
    headers: { 'Content-Type': 'multipart/form-data' }
  })
}

export const addWarID = (data) => {
  return service({
    url: '/competitionTeam/addWarID',
    method: 'post',
    data
  })
}

export const getTeamScores = (params) => {
  return service({
    url: '/competitionTeam/scores',
    method: 'get',
    params
  })
}

export const getTeamRecentScores = (params) => {
  return service({
    url: '/competitionTeam/recentScores',
    method: 'get',
    params
  })
}

export const getTeamDetail = (params) => {
  return service({
    url: '/competitionTeam/detail',
    method: 'get',
    params
  })
}

export const getAllTeamsScoreSummary = () => {
  return service({
    url: '/competitionTeam/allScores',
    method: 'get'
  })
}

export const deleteTeamScore = (params) => {
  return service({
    url: '/competitionTeam/deleteScore',
    method: 'delete',
    params
  })
}

export const updateTeamScore = (data) => {
  return service({
    url: '/competitionTeam/updateScore',
    method: 'put',
    data
  })
}

export const calculateWarIDForAllTeams = (data) => {
  return service({
    url: '/competitionTeam/calculateWarID',
    method: 'post',
    data
  })
}

export const confirmWarIDScores = (data) => {
  return service({
    url: '/competitionTeam/confirmWarID',
    method: 'post',
    data
  })
}

// 公开接口（无需鉴权）
export const getPublicTeamList = () => {
  return service({
    url: '/competitionTeam/public/teamList',
    method: 'get'
  })
}
