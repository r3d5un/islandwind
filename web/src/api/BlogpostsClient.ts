import {
  Blogpost,
  BlogpostDeleteBody,
  BlogpostDeleteOptions,
  BlogpostInput,
  BlogpostListResponse,
  BlogpostPatch,
  BlogpostPatchBody,
  BlogpostPostBody,
} from './blogposts'
import { type IBlogpostListResponse, type IBlogpostResponse } from './blogposts.ts'
import { type ILogObj, Logger } from 'tslog'
import axios, { type AxiosResponse } from 'axios'
import { type RequestFailureError, UnauthorizedError, handleRequestFailure } from '@/api/errors.ts'

export class BlogpostClient {
  readonly baseUrl: string
  private logger: Logger<ILogObj>
  private timeout: number = 5000

  constructor(baseUrl: string, logger: Logger<ILogObj>) {
    this.baseUrl = baseUrl
    this.logger = logger
    this._username = null
    this._password = null
  }

  private _username: string | null

  set username(value: string) {
    this._username = value
  }

  private _password: string | null

  set password(value: string) {
    this._password = value
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
      return handleRequestFailure(error)
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
      return handleRequestFailure(error)
    }
  }

  public async patch(
    accessToken: string,
    blogpost: BlogpostPatch,
  ): Promise<Blogpost | RequestFailureError> {
    this.logger.info('updating blogpost', { blogpost: blogpost })

    if (!this._username || !this._password) {
      return new UnauthorizedError('missing basic authentication credentials')
    }

    try {
      const response: AxiosResponse<IBlogpostResponse, number> = await axios.patch(
        `${this.baseUrl}/api/v1/blog/post`,
        new BlogpostPatchBody(blogpost),
        {
          timeout: this.timeout,
          headers: {
            Authorization: `Bearer ${accessToken}`,
          },
        },
      )
      this.logger.info('blogpost updated')
      return new Blogpost(response.data.data)
    } catch (error) {
      this.logger.error('Error listing blogposts', { error: error })
      return handleRequestFailure(error)
    }
  }

  public async post(
    accessToken: string,
    blogpost: BlogpostInput,
  ): Promise<Blogpost | RequestFailureError> {
    this.logger.info('creating blogpost', { blogpost: blogpost })

    if (!this._username || !this._password) {
      return new UnauthorizedError('missing basic authentication credentials')
    }

    try {
      const response: AxiosResponse<IBlogpostResponse, number> = await axios.post(
        `${this.baseUrl}/api/v1/blog/post`,
        new BlogpostPostBody(blogpost),
        {
          timeout: this.timeout,
          headers: {
            Authorization: `Bearer ${accessToken}`,
          },
        },
      )
      this.logger.info('blogpost created')
      return new Blogpost(response.data.data)
    } catch (error) {
      this.logger.error('Error listing blogposts', { error: error })
      return handleRequestFailure(error)
    }
  }

  public async delete(
    accessToken: string,
    id: string,
    purge: boolean,
  ): Promise<Blogpost | RequestFailureError> {
    this.logger.info('deleting blogpost', { id: id })

    if (!this._username || !this._password) {
      return new UnauthorizedError('missing basic authentication credentials')
    }

    try {
      const response: AxiosResponse<IBlogpostResponse, number> = await axios({
        method: 'delete',
        url: `${this.baseUrl}/api/v1/blog/post`,
        data: new BlogpostDeleteBody(new BlogpostDeleteOptions(id, purge)),
        timeout: this.timeout,
        auth: { username: this._username, password: this._password },
        headers: { Authorization: `Bearer ${accessToken}` },
      })
      this.logger.info('blogpost created')
      return new Blogpost(response.data.data)
    } catch (error) {
      this.logger.error('Error listing blogposts', { error: error })
      return handleRequestFailure(error)
    }
  }
}
