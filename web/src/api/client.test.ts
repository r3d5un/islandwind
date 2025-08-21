import { afterAll, beforeAll, describe, expect, it } from 'vitest'
import { type ILogObj, Logger } from 'tslog'
import { Blogpost, BlogpostInput, BlogpostListResponse, BlogpostPatch } from '@/api/blogposts.ts'
import type { RequestFailureError } from '@/api/errors.ts'
import { Client, type QueryResult } from 'pg'
import { DockerComposeEnvironment, StartedDockerComposeEnvironment, Wait } from 'testcontainers'
import { HttpClient } from '@/api/client.ts'
import { login, Tokens } from '@/api/auth.ts'

interface IBlogpostID {
  id: string
}

describe('BlogpostClient', () => {
  const logger: Logger<ILogObj> = new Logger({
    hideLogPositionForProduction: false,
    type: 'pretty',
  })

  const databaseConnectionString: string =
    'postgres://postgres:postgres@localhost:15432/islandwind?sslmode=disable'
  let databaseClient: Client

  const baseUrl: string = 'http://localhost:14000'
  const client: HttpClient = new HttpClient(logger, baseUrl)
  client.blogposts.username = 'islandwind'
  client.blogposts.password = 'islandwind'
  let tokens: Tokens | RequestFailureError

  let environment: StartedDockerComposeEnvironment

  beforeAll(
    async () => {
      logger.info('Setting up Docker Compose testing environment')
      try {
        environment = await new DockerComposeEnvironment(
          './../.', // Project root
          './deployments/docker-compose.testing.yaml',
        )
          // One shot startup strategy is for containers than run briefly then exit on their own with exit code 0.
          // The migrate container executes all up migrations after the database is, then exits. This ensures that
          // the container is ready before proceeding.
          .withWaitStrategy('migrate-1', Wait.forOneShotStartup())
          .up()
      } catch (error) {
        logger.error('unable to start Docker Compose', { error: error })
        throw error
      }

      logger.info('Connecting database client')
      databaseClient = new Client({ connectionString: databaseConnectionString })
      await databaseClient.connect()

      logger.info('Inserting test data')
      await databaseClient.query(`
          INSERT INTO blog.post (title, content, published)
          VALUES ('Read Me', 'Read Me', false),
                 ('Update Me', 'Update Me', false),
                 ('Delete Me', 'Delete Me', false);
      `)

      logger.info('logging in')
      tokens = await login('islandwind', 'islandwind')
      if (!(tokens instanceof Tokens)) {
        throw tokens
      }
      logger.info('logged in', { tokens: tokens })
    },
    // Timeout set to two minutes because the container environment can take some time to be ready
    120_000,
  )

  afterAll(async () => {
    logger.info('Cleaning up database')
    await databaseClient.query('DROP SCHEMA IF EXISTS blog CASCADE;')
    logger.info('Cleanup complete')

    logger.info('Closing database client')
    await databaseClient.end()

    logger.info('Shutting down Docker Compose testing environment')
    await environment.down()
  })

  it('should create a blogpost', async () => {
    const input = new BlogpostInput('Created Blogpost', 'Content', false)
    if (tokens instanceof Tokens) {
      const result: Blogpost | RequestFailureError = await client.blogposts.post(
        tokens.accessToken,
        input,
      )

      expect(result).toBeInstanceOf(Blogpost)
      if (result instanceof Blogpost) {
        expect(result.id.length).toBeGreaterThan(0)
        expect(result.title).toBe(input.title)
        expect(result.content).toBe(input.content)
        expect(result.published).toBe(input.published)
        expect(result.createdAt).toBeInstanceOf(Date)
      }
    } else {
      throw tokens
    }
  })

  it('should read a blogpost', async () => {
    const queryResult: QueryResult<IBlogpostID> = await databaseClient.query(
      "SELECT id FROM blog.post WHERE title = 'Read Me';",
    )
    const result: Blogpost | RequestFailureError = await client.blogposts.get(
      queryResult.rows[0].id,
    )

    expect(result).toBeInstanceOf(Blogpost)
    if (result instanceof Blogpost) {
      expect(result.title).toBe('Read Me')
      expect(result.content).toBe('Read Me')
      expect(result.createdAt).toBeInstanceOf(Date)
    }
  })

  it('should list blogposts', async () => {
    const result: BlogpostListResponse | RequestFailureError = await client.blogposts.list()
    expect(result).toBeInstanceOf(BlogpostListResponse)

    if (result instanceof BlogpostListResponse) {
      expect(result.data.length).toBeGreaterThan(0)
      expect(result.metadata.next)
      expect(result.metadata.responseLength).toBe(result.data.length)
    }
  })

  it('should update a blogpost', async () => {
    const queryResult: QueryResult<IBlogpostID> = await databaseClient.query(
      "SELECT id FROM blog.post WHERE title = 'Update Me';",
    )
    if (tokens instanceof Tokens) {
      const patch: BlogpostPatch = new BlogpostPatch({
        id: queryResult.rows[0].id,
        published: true,
      })
      const result: Blogpost | RequestFailureError = await client.blogposts.patch(
        tokens.accessToken,
        patch,
      )

      expect(result).toBeInstanceOf(Blogpost)
      if (result instanceof Blogpost) {
        expect(result.id).toBe(queryResult.rows[0].id)
        expect(result.published).toBe(true)
        expect(result.createdAt).toBeInstanceOf(Date)
      }
    } else {
      throw tokens
    }
  })

  it('should delete a blogpost', async () => {
    const queryResult: QueryResult<IBlogpostID> = await databaseClient.query(
      "SELECT id FROM blog.post WHERE title = 'Delete Me';",
    )
    if (tokens instanceof Tokens) {
      const result: Blogpost | RequestFailureError = await client.blogposts.delete(
        tokens.accessToken,
        queryResult.rows[0].id,
        true,
      )
      expect(result).toBeInstanceOf(Blogpost)
      if (result instanceof Blogpost) {
        expect(result.id).toBe(queryResult.rows[0].id)
        expect(result.deleted).toBe(true)
        expect(result.createdAt).toBeInstanceOf(Date)
        expect(result.deletedAt).toBeInstanceOf(Date)
      }
    } else {
      throw tokens
    }
  })
})
