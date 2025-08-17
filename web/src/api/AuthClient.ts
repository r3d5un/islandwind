import {
  BackendServerInternalError,
  BadRequestError,
  NetworkError,
  NotFoundError,
  type RequestFailureError,
  UnexpectedStatusCodeError,
} from '@/api/errors.ts'
import axios, { type AxiosResponse } from 'axios'
import { type ILogObj, Logger } from 'tslog'

export class AuthClient {
  readonly baseUrl: string
  private logger: Logger<ILogObj>
  private timeout: number = 5000

  constructor(baseUrl: string, logger: Logger<ILogObj>) {
    this.baseUrl = baseUrl
    this.logger = logger
  }

  public async get(username: string, password: string): Promise<boolean | RequestFailureError> {
    this.logger.info('Attempting basic auth login')
    try {
      const response: AxiosResponse<null, number> = await axios.get(
        `${this.baseUrl}/api/v1/auth/login`,
        { timeout: this.timeout, auth: { username: username, password: password } },
      )
      return response.status === 200
    } catch (error) {
      if (axios.isAxiosError(error)) {
        if (error.response) {
          switch (error.response.status) {
            case 400:
              return new BadRequestError()
            case 401:
              return false
            case 403:
              return false
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
}
