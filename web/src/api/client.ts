import { BlogpostClient } from '@/api/BlogpostsClient.ts'
import { type ILogObj, Logger } from 'tslog'

export class HttpClient {
  public baseUrl: string
  public blogposts: BlogpostClient

  constructor(logger: Logger<ILogObj>, baseUrl: string) {
    this.baseUrl = baseUrl
    this.blogposts = new BlogpostClient(baseUrl, logger)
  }
}
