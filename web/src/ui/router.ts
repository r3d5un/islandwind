import { createMemoryHistory, createRouter } from 'vue-router'
import HomeView from '@/components/HomeView.vue'
import BlogpostView from '@/components/BlogpostView.vue'
import NotFoundView from '@/components/NotFoundView.vue'

const routes = [
  { path: '/', component: HomeView },
  { path: '/blog', component: HomeView },
  { path: '/blog/:id', component: BlogpostView },
  { path: '/:pathMatch(.*)*', name: 'NotFound', component: NotFoundView },
]

export const router = createRouter({
  history: createMemoryHistory(),
  routes,
})
