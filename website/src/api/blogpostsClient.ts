import { BlogpostListResponse, Blogpost } from './blogposts'
import { type IBlogpostListResponse, type IBlogpostResponse} from './blogposts.ts'
import { type ILogObj, Logger } from 'tslog'
import axios, { type AxiosResponse } from 'axios'

export class BlogpostClient {
  private baseUrl: string
  private logger: Logger<ILogObj>
  private timeout: number = 5000

  constructor(baseUrl: string, logger: Logger<ILogObj>) {
    this.baseUrl = baseUrl
    this.logger = logger
  }

  public async get(id: string): Promise<Blogpost | null> {
    this.logger.info('retrieving blogpost', { id: id })
    try {
      const response: AxiosResponse<IBlogpostResponse, number> = await axios.get(
        `${this.baseUrl}/api/v1/blog/post/${id}`,
        { timeout: this.timeout },
      )
      return new Blogpost(response.data.data)
    } catch (error) {
      this.logger.error('unable to retrieve blogpost', { error: error })
    }

    return null
  }

  public async list(): Promise<BlogpostListResponse | null> {
    this.logger.info('listing blogposts')
    try {
      const response: AxiosResponse<IBlogpostListResponse, number> = await axios.get(
        `${this.baseUrl}/api/v1/blog/post`,
        { timeout: this.timeout },
      )
      this.logger.info('blogposts listed')
      return new BlogpostListResponse(response.data.data, response.data.metadata)
    } catch (error) {
      this.logger.error('Error listing blogposts', { error: error })
      return null
    }
  }
}
