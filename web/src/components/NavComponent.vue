<script setup lang="ts">
import { useAuthStore } from '@/api/auth.ts'

const authStore = useAuthStore()
</script>

<template>
  <nav>
    <div class="nav-container">
      <div class="nav-left">
        <router-link :to="{ name: 'Home' }">Islandwind</router-link>
      </div>
      <div class="nav-right">
        <router-link v-if="!authStore.loggedIn" :to="{ name: 'Login' }">Login</router-link>
      </div>
    </div>
  </nav>
</template>

<style scoped>
nav {
  font-family: Arial, sans-serif;
  background-color: #f8f8f8;
  padding: 10px 20px;
  position: fixed;
  width: 100%;
  top: 0;
  z-index: 1000;
}
.nav-container {
  display: flex;
  justify-content: space-between;
  align-items: center;
  max-width: 800px;
  margin: 0 auto;
  width: 100%;
}
.nav-left,
.nav-right {
  display: flex;
  align-items: center;
}
.nav-right {
  gap: 10px;
}
nav h1 {
  margin: 0;
  font-size: 1.5rem;
}
nav a {
  text-decoration: none;
  color: #333;
  font-weight: bold;
  padding: 5px 10px;
  border-radius: 4px;
  transition: background-color 0.3s;
}
nav a:hover {
  background-color: #ddd;
}

/* Media query for smaller screens */
@media (max-width: 850px) {
  .nav-container {
    flex-direction: column;
    align-items: flex-start;
  }
  .nav-right {
    margin-top: 10px;
    width: 100%;
  }
  .nav-left,
  .nav-right {
    width: 100%;
  }
}
</style>
