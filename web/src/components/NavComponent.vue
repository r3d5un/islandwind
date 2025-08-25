<script setup lang="ts">
import { useAuthStore } from '@/api/auth.ts'

const authStore = useAuthStore()
</script>

<template>
  <nav>
    <router-link :to="{ name: 'Home' }">Home</router-link>
    <router-link v-if="authStore.loggedIn === false" :to="{ name: 'Login' }">Login</router-link>
  </nav>
</template>

<style scoped>
nav > a + a {
  margin-left: 10px;
}
</style>
