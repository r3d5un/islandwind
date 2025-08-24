import { useLogger } from './logging.ts'
import { HttpClient } from '@/api/client.ts'
import axios, {
  AxiosError,
  type AxiosInstance,
  type AxiosResponse,
  type InternalAxiosRequestConfig,
} from 'axios'
import { refresh, Tokens, useAuthStore } from '@/api/auth.ts'
import { handleRequestFailure } from '@/api/errors.ts'

const baseUrl: string = 'http://localhost:4000'

const httpClient = new HttpClient(useLogger(), baseUrl)

const apiClient: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_URL,
  timeout: 5_000,
})

apiClient.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    let accessToken: string = ''
    const authStore = useAuthStore()
    if (authStore.loggedIn) {
      accessToken = authStore.tokens.accessToken
    }

    if (accessToken) {
      config.headers = config.headers || {}
      config.headers.Authorization = `Bearer ${accessToken}`
    }
    return config
  },
  (error: AxiosError) => {
    return Promise.reject(error)
  },
)

interface CustomAxiosRequestConfig extends InternalAxiosRequestConfig {
  _retry?: boolean
}

apiClient.interceptors.response.use(
  (response: AxiosResponse) => response,
  async (error: AxiosError) => {
    const originalRequest = error.config as CustomAxiosRequestConfig

    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true

      try {
        const authStore = useAuthStore()
        const newTokens = await refresh(authStore.tokens.refreshToken)

        if (newTokens instanceof Tokens) {
          authStore.tokens.accessToken = newTokens.accessToken
          authStore.tokens.refreshToken = newTokens.refreshToken
        }

        originalRequest.headers.Authorization = `Bearer ${authStore.tokens.accessToken}`

        return apiClient(originalRequest)
      } catch (refreshError) {
        return Promise.reject(handleRequestFailure(refreshError))
      }
    }

    return Promise.reject(handleRequestFailure(error))
  },
)

export function useApiClient() {
  return apiClient
}

export function useClient() {
  return httpClient
}
