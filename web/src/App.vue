<script lang="ts" setup>
import { useLogger } from '@/ui/logging.ts'
import NavComponent from '@/components/NavComponent.vue'
import { onMounted, ref } from 'vue'

const logger = useLogger()

const navbarHeight = ref(0)

logger.info('Loading UI')

onMounted(() => {
  const navbar = document.querySelector('nav')
  if (navbar) {
    navbarHeight.value = navbar.clientHeight
  }
})
</script>

<template>
  <nav>
    <NavComponent />
  </nav>
  <main>
    <RouterView />
  </main>
</template>

<style scoped>
main {
  padding-top: 40px;
}
</style>
