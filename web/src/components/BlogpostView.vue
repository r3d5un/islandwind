<script lang="ts" setup>
import { useRoute } from 'vue-router'
import { onMounted, ref, watch } from 'vue'
import { useClient } from '@/ui/client.ts'
import { Blogpost } from '@/api/blogposts.ts'
import { useLogger } from '@/ui/logging.ts'

const logger = useLogger()
const route = useRoute()
const client = useClient()
const blogpost = ref<Blogpost>()

function validateRouteId(input: string | string[]): string {
  if (!input) {
    logger.error('No ID provided in route')
    blogpost.value = undefined
    return ''
  }

  return Array.isArray(input) ? input[0] : input
}

async function retrieveBlogpost(id: string): Promise<void> {
  const response = await client.blogposts.get(id)
  if (response instanceof Blogpost) {
    logger.info('retrieved blogpost', { blogpost: blogpost })
    blogpost.value = response
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
  <h2 v-if="blogpost">{{ blogpost.title }}</h2>
  <p v-if="blogpost">{{ blogpost.createdAt }} - {{ blogpost.id }}</p>
</template>

<style scoped>
h2 {
  font-family: Monospaced, monospace;
}

p {
  font-family: Monospaced, monospace;
}
</style>
