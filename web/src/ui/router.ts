import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '@/components/HomeView.vue'
import BlogpostView from '@/components/BlogpostView.vue'
import NotFoundView from '@/components/NotFoundView.vue'

const routes = [
  { path: '/', name: 'Home', component: HomeView },
  { path: '/blog', redirect: '/' },
  { path: '/blog/:id', name: 'BlogpostView', component: BlogpostView },
  { path: '/:pathMatch(.*)*', name: 'NotFound', component: NotFoundView },
]

export const router = createRouter({
  history: createWebHistory(),
  routes,
})
