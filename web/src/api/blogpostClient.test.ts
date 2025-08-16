import { afterAll, beforeAll, describe, expect, it } from 'vitest'
import { type ILogObj, Logger } from 'tslog'
import { BlogpostClient } from '@/api/blogpostsClient.ts'
import { Blogpost, BlogpostInput, BlogpostListResponse, BlogpostPatch } from '@/api/blogposts.ts'
import type { RequestFailureError } from '@/api/errors.ts'
import { PostgreSqlContainer, StartedPostgreSqlContainer } from '@testcontainers/postgresql'
import { Client } from 'pg'
import { migrate, MigrationDirection } from '@/testsuite/db.ts'

describe('BlogpostClient', () => {
  let postgresContainer: StartedPostgreSqlContainer
  let postgresClient: Client

  const baseUrl: string = 'http://localhost:4000'
  const logger: Logger<ILogObj> = new Logger({
    hideLogPositionForProduction: false,
    type: 'pretty',
  })
  const blogpostClient: BlogpostClient = new BlogpostClient(baseUrl, logger)
  blogpostClient.username = 'islandwind'
  blogpostClient.password = 'islandwind'
  const blogpostIds: string[] = []

  beforeAll(async () => {
    logger.info('Setting up Docker network')

    logger.info('Starting Postgres container')
    postgresContainer = await new PostgreSqlContainer('postgres:17.6').start()
    postgresClient = new Client({ connectionString: postgresContainer.getConnectionUri() })
    await postgresClient.connect()
    logger.info('Postgres container started')

    logger.info('performing up migrations')
    try {
      await migrate(postgresClient, MigrationDirection.up)
    } catch (error) {
      logger.error('unable to run migrations', { error: error })
    }
  }, 60_000)

  afterAll(async () => {
    logger.info('performing down migrations')
    try {
      await migrate(postgresClient, MigrationDirection.down)
    } catch (error) {
      logger.error('unable to run migrations', { error: error })
    }

    await postgresClient.end()
    await postgresContainer.stop()
  })

  it('should SELECT anything', async () => {
    const result = await postgresClient.query("SELECT 'Hello, World!' AS col1")
    expect(result.rows[0].col1).toEqual('Hello, World!')
  })

  it('should SELECT from blog.post', async () => {
    try {
      await postgresClient.query('SELECT * FROM blog.post;')
    } catch (error) {
      logger.error('something went wrong', { error: error })
    }
  })

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
