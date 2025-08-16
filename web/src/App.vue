<script setup lang="ts">
import { type ILogObj, Logger } from 'tslog'

const logger: Logger<ILogObj> = new Logger({
  hideLogPositionForProduction: true,
  type: 'json',
})

logger.info('Loading UI')
</script>

<template>
  <main>
    <h1>Hello, World!</h1>
  </main>
</template>

<style scoped>
h1 {
  font-family: Monospaced, monospace;
}
</style>
