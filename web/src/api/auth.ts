import { handleRequestFailure, type RequestFailureError } from '@/api/errors.ts'
import axios, { type AxiosResponse } from 'axios'
import { defineStore, type StoreDefinition } from 'pinia'

export const useAuthStore: StoreDefinition<'tokens', { tokens: Tokens }> = defineStore('tokens', {
  state: () => ({ tokens: new Tokens({ accessToken: '', refreshToken: '' }) }),
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
