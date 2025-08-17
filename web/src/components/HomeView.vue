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
      <router-link :to="{ name: 'BlogpostView', params: { id: post.id } }"
        >{{ post.createdAt }} - {{ post.title }}</router-link
      >
    </li>
  </ul>
</template>

<style scoped></style>
