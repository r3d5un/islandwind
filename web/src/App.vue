<script lang="ts" setup>
import { useLogger } from '@/ui/logging.ts'
import NavComponent from '@/components/NavComponent.vue'

const logger = useLogger()

logger.info('Loading UI')
</script>

<template>
  <header>
    <h1>Islandwind</h1>
    <NavComponent />
  </header>
  <main>
    <RouterView />
  </main>
</template>

<style scoped></style>
