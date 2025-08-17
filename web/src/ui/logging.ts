import { type ILogObj, Logger } from 'tslog'

export const logger: Logger<ILogObj> = new Logger({
  hideLogPositionForProduction: true,
  type: 'json',
})

export function useLogger() {
  return logger
}
