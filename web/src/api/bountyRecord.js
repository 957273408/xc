import service from '@/utils/request'

export const getRecordList = (params) => {
  return service({
    url: '/bountyRecord/recordList',
    method: 'get',
    params
  })
}