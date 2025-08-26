<script setup lang="ts">
import { useAuthStore } from '@/api/auth.ts'
import { logout } from '@/ui/auth.ts'

const authStore = useAuthStore()
</script>

<template>
  <nav>
    <div class="nav-container">
      <router-link :to="{ name: 'Home' }">Islandwind</router-link>
      <router-link v-if="!authStore.loggedIn" :to="{ name: 'Login' }">Login</router-link>
      <a v-else href="#" @click.prevent="logout">Logout</a>
    </div>
  </nav>
</template>

<style scoped>
nav {
  font-family: Arial, sans-serif;
  background-color: #f8f8f8;
  padding: 10px 0;
  position: fixed;
  width: 100%;
  z-index: 1000;
}
.nav-container {
  display: flex;
  align-items: center;
  max-width: 800px;
  margin: auto;
  width: 100%;
}
nav h1 {
  margin: 0;
  font-size: 1.5rem;
}
nav a {
  text-decoration: none;
  color: #333;
  font-weight: lighter;
  padding: 5px 10px;
  border-radius: 10px;
  transition: background-color 0.3s;
}
nav a:hover {
  background-color: #ddd;
}
</style>
