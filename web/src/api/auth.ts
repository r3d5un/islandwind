import { handleRequestFailure, type RequestFailureError } from '@/api/errors.ts'
import axios, { type AxiosResponse } from 'axios'
import { defineStore } from 'pinia'
import { logger } from '@/ui/logging.ts'

export const useAuthStore = defineStore('tokens', {
  state: () => ({
    tokens: new Tokens({ accessToken: '', refreshToken: '' }),
    loggedIn: false,
    loading: false,
    logger: logger,
  }),
  actions: {
    async login(username: string, password: string): Promise<boolean> {
      this.loading = true
      try {
        logger.info('logging in')
        const result = await login(username, password)
        if (result instanceof Tokens) {
          this.tokens = result
          this.loggedIn = true
          return true
        }
        return false
      } catch (error) {
        logger.info('unable to login', { error: error })
        return false
      } finally {
        this.loading = false
      }
    },
  },
})

export async function login(
  username: string,
  password: string,
): Promise<Tokens | RequestFailureError> {
  try {
    const response: AxiosResponse<ITokens, number> = await axios({
      method: 'post',
      url: `${import.meta.env.VITE_API_URL}/api/v1/auth/login`,
      timeout: import.meta.env.VITE_API_TIMEOUT,
      auth: { username: username, password: password },
    })
    return new Tokens(response.data)
  } catch (error) {
    return handleRequestFailure(error)
  }
}

export async function refresh(refreshToken: string): Promise<Tokens | RequestFailureError> {
  try {
    const response: AxiosResponse<ITokens, number> = await axios.post(
      `${import.meta.env.VITE_API_URL}/api/v1/auth/refresh`,
      new RefreshRequestBody(refreshToken),
      { timeout: import.meta.env.VITE_API_TIMEOUT },
    )
    return new Tokens(response.data)
  } catch (error) {
    return handleRequestFailure(error)
  }
}

export async function invalidateRefreshToken(refreshToken: string): Promise<void> {
  try {
    await axios({
      method: 'post',
      url: `${import.meta.env.VITE_API_URL}/api/v1/auth/logout`,
      data: new RefreshRequestBody(refreshToken),
      timeout: import.meta.env.VITE_API_TIMEOUT,
    })
  } catch (error) {
    throw handleRequestFailure(error)
  }
}

export interface ITokens {
  accessToken: string
  refreshToken: string
}

export class Tokens {
  public accessToken: string
  public refreshToken: string

  constructor(input: ITokens) {
    this.accessToken = input.accessToken
    this.refreshToken = input.refreshToken
  }
}

class RefreshRequestBody {
  refreshToken: string

  constructor(refreshToken: string) {
    this.refreshToken = refreshToken
  }
}
