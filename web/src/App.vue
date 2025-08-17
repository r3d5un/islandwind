<script lang="ts" setup>
import { useLogger } from '@/ui/logging.ts'

const logger = useLogger()
logger.info('Loading UI')
</script>

<template>
  <header>
    <h1>Islandwind</h1>
    <p>Path: {{ $route.fullPath }}</p>
    <nav>
      <router-link :to="{ name: 'Home' }">Home</router-link>
    </nav>
  </header>
  <main>
    <RouterView />
  </main>
</template>

<style scoped>
h1 {
  font-family: Monospaced, monospace;
}

p {
  font-family: Monospaced, monospace;
}

a {
  font-family: Monospaced, monospace;
}
nav > a + a {
  margin-left: 10px;
}
</style>
