import { useLogger } from './logging.ts'
import { HttpClient } from '@/api/client.ts'

const baseUrl: string = 'http://localhost:4000'

const httpClient = new HttpClient(useLogger(), baseUrl)

export function useClient() {
  return httpClient
}
