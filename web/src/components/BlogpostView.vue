<script lang="ts" setup>
import { useRoute } from 'vue-router'
import { onMounted, ref, watch } from 'vue'
import { Blogpost } from '@/api/blogposts.ts'
import { useLogger } from '@/ui/logging.ts'
import { useApiClient } from '@/ui/client.ts'
import { BlogpostApiClient } from '@/api/blog.ts'

const logger = useLogger()
const route = useRoute()
const apiClient = useApiClient()
const blogpostClient = new BlogpostApiClient(apiClient)
const blogpost = ref<Blogpost>()
const content = ref('')

function validateRouteId(input: string | string[]): string {
  if (!input) {
    logger.error('No ID provided in route')
    blogpost.value = undefined
    return ''
  }

  return Array.isArray(input) ? input[0] : input
}

async function retrieveBlogpost(id: string): Promise<void> {
  const response = await blogpostClient.get(id)
  if (response instanceof Blogpost) {
    logger.info('retrieved blogpost', { blogpost: blogpost })
    blogpost.value = response
    content.value = await response.markdownContent()
  } else {
    logger.error('Unable to retrieve blogpost', { id: id, error: blogpost })
    blogpost.value = undefined
  }
}

watch(
  () => route.params.id,
  async (newId) => {
    const id = validateRouteId(newId)
    await retrieveBlogpost(id)
  },
)

onMounted(async () => {
  const newId = route.params.id
  const id = validateRouteId(newId)
  await retrieveBlogpost(id)
})
</script>

<template>
  <div class="markdown-content" v-html="content"></div>
</template>

<style scoped>
.markdown-content {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
  font-family: Arial, sans-serif;
  line-height: 1.6;
  color: #333;
}

.markdown-content h1,
.markdown-content h2,
.markdown-content h3,
.markdown-content h4,
.markdown-content h5,
.markdown-content h6 {
  color: #2c3e50;
  margin-top: 1.5em;
  margin-bottom: 0.5em;
  line-height: 1.2;
}

.markdown-content h1 {
  font-size: 2em;
  border-bottom: 1px solid #eee;
  padding-bottom: 0.3em;
}

.markdown-content h2 {
  font-size: 1.5em;
  border-bottom: 1px solid #eee;
  padding-bottom: 0.3em;
}

.markdown-content h3 {
  font-size: 1.25em;
}

.markdown-content h4 {
  font-size: 1em;
}

.markdown-content p {
  margin-bottom: 1em;
}

.markdown-content a {
  color: #3498db;
  text-decoration: none;
}

.markdown-content a:hover {
  text-decoration: underline;
}

.markdown-content ul,
.markdown-content ol {
  margin-bottom: 1em;
  padding-left: 2em;
}

.markdown-content li {
  margin-bottom: 0.5em;
}

.markdown-content blockquote {
  border-left: 4px solid #ddd;
  padding-left: 1em;
  margin-left: 0;
  color: #7f8c8d;
  font-style: italic;
}

.markdown-content code {
  background-color: #f4f4f4;
  border-radius: 3px;
  padding: 0.2em 0.4em;
  font-family: monospace;
  font-size: 0.9em;
}

.markdown-content pre {
  background-color: #f8f8f8;
  border-radius: 3px;
  padding: 1em;
  overflow-x: auto;
  font-family: monospace;
  font-size: 0.9em;
}

.markdown-content table {
  border-collapse: collapse;
  width: 100%;
  margin-bottom: 1em;
}

.markdown-content th,
.markdown-content td {
  border: 1px solid #ddd;
  padding: 8px;
  text-align: left;
}

.markdown-content th {
  background-color: #f2f2f2;
}

.markdown-content img {
  max-width: 100%;
  height: auto;
}

.markdown-content hr {
  border: 0;
  height: 1px;
  background: #ddd;
  margin: 1em 0;
}
</style>
