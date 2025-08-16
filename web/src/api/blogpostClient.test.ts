import { describe, it } from 'vitest'
import { type ILogObj, Logger } from 'tslog'
import { BlogpostClient } from '@/api/blogpostsClient.ts'
import { Blogpost, BlogpostInput, BlogpostListResponse, BlogpostPatch } from '@/api/blogposts.ts'
import type { RequestFailureError } from '@/api/errors.ts'

describe('BlogpostClient', () => {
  const baseUrl: string = 'http://localhost:4000'
  const logger: Logger<ILogObj> = new Logger({
    hideLogPositionForProduction: false,
    type: 'pretty',
  })
  const blogpostClient: BlogpostClient = new BlogpostClient(baseUrl, logger)
  blogpostClient.username = 'islandwind'
  blogpostClient.password = 'islandwind'
  const blogpostIds: string[] = []

  it('should list blogposts', async () => {
    const blogposts: BlogpostListResponse | RequestFailureError = await blogpostClient.list()
    if (blogposts instanceof BlogpostListResponse) {
      blogposts?.data.forEach((value: Blogpost) => {
        logger.info('blogpost', { blogpost: value })
        blogpostIds.push(value.id)
      })
    }

    for (const id of blogpostIds) {
      logger.info('test GET request', { id: id })
      const blogpost: Blogpost | RequestFailureError = await blogpostClient.get(id)
      logger.info('blogpost retrieved', { blogpost: blogpost })
    }
  })

  it('should create then delete a blogpost', async () => {
    const blogpost: Blogpost | RequestFailureError = await blogpostClient.post(
      new BlogpostInput('Example title', 'This is some sample content', false),
    )
    logger.info('post performed', { blogpost: blogpost })

    if (blogpost instanceof Blogpost && blogpost) {
      await blogpostClient.delete(blogpost.id, true)
    }
  })

  it('should create, update, then delete a blogpost', async () => {
    const posted: Blogpost | RequestFailureError = await blogpostClient.post(
      new BlogpostInput('Update me', 'This content should be updated', false),
    )
    logger.info('post performed', { blogpost: posted })

    if (posted instanceof Blogpost && posted) {
      const updated: Blogpost | RequestFailureError = await blogpostClient.patch(
        new BlogpostPatch({
          id: posted.id,
          title: 'New Title',
          content: 'This content is updated',
          published: true,
        }),
      )
      logger.info('blogpost updated', { blogpost: updated })

      if (updated instanceof Blogpost && updated) {
        await blogpostClient.delete(updated.id, true)
      }
    }
  })
})
