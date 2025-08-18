import { BlogpostClient } from '@/api/BlogpostsClient.ts'
import { type ILogObj, Logger } from 'tslog'
import { AuthClient } from '@/api/AuthClient.ts'
import { LoginError } from '@/api/errors.ts'

export class HttpClient {
  public baseUrl: string
  public blogposts: BlogpostClient
  public auth: AuthClient
  private logger: Logger<ILogObj>

  constructor(logger: Logger<ILogObj>, baseUrl: string) {
    this.baseUrl = baseUrl
    this.blogposts = new BlogpostClient(baseUrl, logger)
    this.auth = new AuthClient(baseUrl, logger)
    this.logger = logger
  }

  public async login(username: string, password: string): Promise<void | LoginError> {
    const result = await this.auth.get(username, password)
    if (typeof result === 'boolean') {
      if (result) {
        this.logger.info('login successful')
        this.blogposts.username = username
        this.blogposts.password = password
        return
      }
    }
    this.logger.info('unable to login')
    return new LoginError(`Unable to login: ${result}`)
  }
}
