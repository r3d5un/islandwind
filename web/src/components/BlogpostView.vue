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

watch(
  () => route.params.id,
  async (newId) => {
    if (!newId) {
      logger.error('No ID provided in route')
      blogpost.value = undefined
      return
    }

    const id = Array.isArray(newId) ? newId[0] : newId

    const response = await client.blogposts.get(id)
    if (response instanceof Blogpost) {
      logger.info('retrieved blogpost', { blogpost: blogpost })
      blogpost.value = response
    } else {
      logger.error('Unable to retrieve blogpost', { id: newId, error: blogpost })
    }
  },
)

onMounted(async () => {
  const newId = route.params.id
  if (!newId) {
    logger.error('No ID provided in route')
    blogpost.value = undefined
    return
  }
  const id = Array.isArray(newId) ? newId[0] : newId

  const response = await client.blogposts.get(id)
  if (response instanceof Blogpost) {
    logger.info('retrieved blogpost', { blogpost: blogpost })
    blogpost.value = response
  } else {
    logger.error('Unable to retrieve blogpost', { id: newId, error: blogpost })
  }
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
