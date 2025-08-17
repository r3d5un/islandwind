import { BlogpostClient } from '@/api/BlogpostsClient.ts'
import { type ILogObj, Logger } from 'tslog'
import { AuthClient } from '@/api/AuthClient.ts'

export class HttpClient {
  public baseUrl: string
  public blogposts: BlogpostClient
  public auth: AuthClient

  constructor(logger: Logger<ILogObj>, baseUrl: string) {
    this.baseUrl = baseUrl
    this.blogposts = new BlogpostClient(baseUrl, logger)
    this.auth = new AuthClient(baseUrl, logger)
  }
}
