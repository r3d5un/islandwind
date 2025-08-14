import { describe, it } from 'vitest'
import { type ILogObj, Logger } from 'tslog'
import { BlogpostClient } from '@/api/blogpostsClient.ts'
import { Blogpost, BlogpostListResponse } from '@/api/blogposts.ts'
import type { RequestFailureError } from '@/api/errors.ts'

describe('BlogpostClient', () => {
  const baseUrl: string = 'http://localhost:4000'
  const logger: Logger<ILogObj> = new Logger({
    hideLogPositionForProduction: false,
    type: 'pretty',
  })
  const blogpostClient: BlogpostClient = new BlogpostClient(baseUrl, logger)
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
})
