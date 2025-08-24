import { type AxiosInstance, type AxiosResponse } from 'axios'
import {
  Blogpost,
  BlogpostDeleteBody,
  BlogpostDeleteOptions,
  BlogpostInput,
  BlogpostListResponse,
  BlogpostPatch,
  BlogpostPatchBody,
  BlogpostPostBody,
  type IBlogpostListResponse,
  type IBlogpostResponse,
} from '@/api/blogposts.ts'
import { handleRequestFailure, type RequestFailureError } from '@/api/errors.ts'

export class BlogpostApiClient {
  private client: AxiosInstance
  constructor(client: AxiosInstance) {
    this.client = client
  }

  public async get(id: string): Promise<Blogpost | RequestFailureError> {
    try {
      const response: AxiosResponse<IBlogpostResponse, number> = await this.client.get(
        `/api/v1/blog/post/${id}`,
      )
      return new Blogpost(response.data.data)
    } catch (error) {
      return handleRequestFailure(error)
    }
  }

  public async list(): Promise<BlogpostListResponse | RequestFailureError> {
    try {
      const response: AxiosResponse<IBlogpostListResponse, number> =
        await this.client.get('/api/v1/blog/post')
      return new BlogpostListResponse(response.data.data, response.data.metadata)
    } catch (error) {
      return handleRequestFailure(error)
    }
  }

  public async patch(blogpost: BlogpostPatch): Promise<Blogpost | RequestFailureError> {
    try {
      const response: AxiosResponse<IBlogpostResponse, number> = await this.client.patch(
        '/api/v1/blog/post',
        new BlogpostPatchBody(blogpost),
      )
      return new Blogpost(response.data.data)
    } catch (error) {
      return handleRequestFailure(error)
    }
  }

  public async post(input: BlogpostInput): Promise<Blogpost | RequestFailureError> {
    try {
      const response: AxiosResponse<IBlogpostResponse, number> = await this.client.post(
        '/api/v1/blog/post',
        new BlogpostPostBody(input),
      )
      return new Blogpost(response.data.data)
    } catch (error) {
      return handleRequestFailure(error)
    }
  }

  public async delete(id: string, purge: boolean): Promise<Blogpost | RequestFailureError> {
    try {
      const response: AxiosResponse<IBlogpostResponse, number> = await this.client.request({
        method: 'delete',
        url: '/api/v1/blog/post',
        data: new BlogpostDeleteBody(new BlogpostDeleteOptions(id, purge)),
      })
      return new Blogpost(response.data.data)
    } catch (error) {
      return handleRequestFailure(error)
    }
  }
}
