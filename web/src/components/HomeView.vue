<script lang="ts" setup>
import { onMounted, ref } from 'vue'
import { BlogpostListResponse } from '@/api/blogposts.ts'
import { useLogger } from '@/ui/logging.ts'
import { useApiClient } from '@/ui/client.ts'
import { BlogpostApiClient } from '@/api/blog.ts'
import { formatDate } from '@/ui/time.ts'

const logger = useLogger()
const apiClient = useApiClient()
const blogpostClient = new BlogpostApiClient(apiClient)

logger.info('Retrieving blogposts')
const blogposts = ref<BlogpostListResponse['data']>([])
onMounted(async () => {
  const response = await blogpostClient.list()
  if (response instanceof BlogpostListResponse) {
    blogposts.value = response.data.reverse()
  } else {
    logger.error('Unable to retrieve blogposts', { error: response })
  }
})
</script>

<template>
  <ul class="blogpost-list">
    <li v-for="post of blogposts" :key="post.id" class="blogpost-item">
      <router-link :to="{ name: 'BlogpostView', params: { id: post.id } }" class="blogpost-link">
        {{ formatDate(post.createdAt) }} - {{ post.title }}
      </router-link>
    </li>
  </ul>
</template>

<style scoped>
h2 {
  font-family: Arial, sans-serif;
  color: #333;
  text-align: center;
  margin-bottom: 20px;
}

.blogpost-list {
  list-style: none;
  padding: 0;
  max-width: 800px;
  margin: 50px auto 0;
  font-family: Arial, sans-serif;
}

.blogpost-item {
  margin-bottom: 15px;
  padding: 15px;
  background-color: #f9f9f9;
  border-radius: 5px;
  transition: background-color 0.3s;
}

.blogpost-item:hover {
  background-color: #e9e9e9;
}

.blogpost-link {
  text-decoration: none;
  color: #333;
  font-weight: bold;
  display: block;
}

.blogpost-link:hover {
  color: #007bff;
}
</style>
