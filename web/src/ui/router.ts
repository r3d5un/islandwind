import { createMemoryHistory, createRouter } from 'vue-router'
import HomeView from '@/components/HomeView.vue'
import BlogpostView from '@/components/BlogpostView.vue'

const routes = [
  { path: '/', component: HomeView },
  { path: '/blogpost', component: BlogpostView },
]

export const router = createRouter({
  history: createMemoryHistory(),
  routes,
})
