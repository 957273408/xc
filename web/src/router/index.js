import { createRouter, createWebHashHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    redirect: '/login'
  },
  {
    path: '/init',
    name: 'Init',
    component: () => import('@/view/init/index.vue')
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/view/login/index.vue')
  },
  {
    path: '/scanUpload',
    name: 'ScanUpload',
    meta: {
      title: '扫码上传',
      client: true
    },
    component: () => import('@/view/example/upload/scanUpload.vue')
  },
  {
    path: '/competitionTeam',
    name: 'CompetitionTeam',
    meta: {
      title: '战队信息管理'
    },
    component: () => import('@/view/example/competitionTeam/index.vue')
  },
  {
    path: '/playerSelection',
    name: 'PlayerSelection',
    meta: {
      title: '选手数据选择展示',
      client: true
    },
    component: () => import('@/view/example/playerSelection/index.vue')
  },
  {
    path: '/multiWar',
    name: 'MultiWar',
    meta: {
      title: '多场选手数据汇总',
      client: true
    },
    component: () => import('@/view/example/playerSelection/multiWar.vue')
  },
  {
    path: '/:catchAll(.*)',
    meta: {
      closeTab: true
    },
    component: () => import('@/view/error/index.vue')
  },
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

export default router
