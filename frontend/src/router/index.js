import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  { path: '/', name: 'PredictsList', component: () => import('@/view/DigitsLists.vue') },
  { path: '/login', name: 'Login', component: () => import('@/view/Login.vue') },
  { path: '/smart-generation', name: 'SmartGeneration', component: () => import('@/view/SmartGeneration.vue') },
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
})

export default router
