<script lang="ts" setup>
import { type ILogObj, Logger } from 'tslog'

const logger: Logger<ILogObj> = new Logger({
  hideLogPositionForProduction: true,
  type: 'json',
})

logger.info('Loading UI')
</script>

<template>
  <header>
    <h1>Islandwind</h1>
    <p>Path: {{ $route.fullPath }}</p>
    <nav>
      <RouterLink to="/">Home</RouterLink>
      <RouterLink to="/blogpost">Blogpost</RouterLink>
    </nav>
    <p>
      <strong></strong>
    </p>
  </header>
  <main>
    <RouterView />
  </main>
</template>

<style scoped>
h1 {
  font-family: Monospaced, monospace;
}

h2 {
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
