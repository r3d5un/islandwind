<script lang="ts" setup>
import { onMounted, ref } from 'vue'
import { BlogpostListResponse } from '@/api/blogposts.ts'
import { useLogger } from '@/ui/logging.ts'
import { useClient } from '@/ui/client.ts'

const logger = useLogger()
const client = useClient()

logger.info('Retrieving blogposts')
const blogposts = ref<BlogpostListResponse['data']>([])
onMounted(async () => {
  const response = await client.blogposts.list()
  if (response instanceof BlogpostListResponse) {
    blogposts.value = response.data
  } else {
    logger.error('Unable to retrieve blogposts', { error: response })
  }
})
</script>

<template>
  <h2>HomeView</h2>
  <ul>
    <li v-for="post of blogposts" :key="post.id">
      <RouterLink :to="`/blog/${post.id}`"
        >{{ post.id }}: {{ post.title }} - {{ post.createdAt }}</RouterLink
      >
    </li>
  </ul>
</template>

<style scoped>
h2 {
  font-family: Monospaced, monospace;
}

a {
  font-family: Monospaced, monospace;
}
</style>
