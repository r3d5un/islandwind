import { handleRequestFailure, type RequestFailureError } from '@/api/errors.ts'
import { useAuthStore, invalidateRefreshToken } from '@/api/auth.ts'
import { useLogger } from '@/ui/logging.ts'

export async function logout(): Promise<void | RequestFailureError> {
  const authStore = useAuthStore()
  const logger = useLogger()

  try {
    logger.info('invalidating token')
    await invalidateRefreshToken(authStore.tokens.refreshToken)

    logger.info('purging local tokens')
    authStore.loggedIn = false
    authStore.tokens = { accessToken: '', refreshToken: '' }
  } catch (error) {
    logger.error('unable to logout', { error: error })
    return handleRequestFailure(error)
  }
}
