<script lang="ts" setup>
import { useLogger } from '@/ui/logging.ts'
import { useAuthStore } from '@/api/auth.ts'

const logger = useLogger()
const authStore = useAuthStore()

logger.info('Loading UI')
</script>

<template>
  <header>
    <h1>Islandwind</h1>
    <p>Path: {{ $route.fullPath }}</p>
    <nav>
      <router-link :to="{ name: 'Home' }">Home</router-link>
      <router-link v-if="authStore.loggedIn === false" :to="{ name: 'Login' }">Login</router-link>
    </nav>
  </header>
  <main>
    <RouterView />
  </main>
</template>

<style scoped>
nav > a + a {
  margin-left: 10px;
}
</style>
