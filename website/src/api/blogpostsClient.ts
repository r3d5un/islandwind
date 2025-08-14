import { BlogpostListResponse, Blogpost } from './blogposts'
import { type IBlogpostListResponse, type IBlogpostResponse } from './blogposts.ts'
import { type ILogObj, Logger } from 'tslog'
import axios, { type AxiosResponse } from 'axios'
import {
  BackendServerInternalError,
  NotFoundError,
  NetworkError,
  UnauthorizedError,
  UnexpectedStatusCodeError,
  ForbiddenError,
  BadRequestError,
  type RequestFailureError,
} from '@/api/errors.ts'

export class BlogpostClient {
  readonly baseUrl: string
  private logger: Logger<ILogObj>
  private timeout: number = 5000

  constructor(baseUrl: string, logger: Logger<ILogObj>) {
    this.baseUrl = baseUrl
    this.logger = logger
  }

  public async get(id: string): Promise<Blogpost | RequestFailureError> {
    this.logger.info('retrieving blogpost', { id: id })
    try {
      const response: AxiosResponse<IBlogpostResponse, number> = await axios.get(
        `${this.baseUrl}/api/v1/blog/post/${id}`,
        { timeout: this.timeout },
      )
      return new Blogpost(response.data.data)
    } catch (error) {
      this.logger.error('unable to retrieve blogpost', { error: error })
      return this.handleRequestFailure(error)
    }
  }

  public async list(): Promise<BlogpostListResponse | RequestFailureError> {
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
      return this.handleRequestFailure(error)
    }
  }

  private handleRequestFailure(error: unknown): RequestFailureError {
    if (axios.isAxiosError(error)) {
      if (error.response) {
        switch (error.response.status) {
          case 400:
            return new BadRequestError()
          case 401:
            return new UnauthorizedError()
          case 403:
            return new ForbiddenError()
          case 404:
            return new NotFoundError()
          case 500:
            return new BackendServerInternalError()
          default:
            return new UnexpectedStatusCodeError()
        }
      } else if (error.request) {
        return new NetworkError()
      } else {
        return new Error('Unknown request failure')
      }
    }

    if (error instanceof Error) {
      return new NetworkError(error.message)
    }

    return new Error('Unknown request failure')
  }
}
