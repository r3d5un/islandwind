import axios from 'axios'

export class NetworkError extends Error {
  constructor(public message: string = 'Unhandled network error occurred') {
    super(message)
    this.name = 'NetworkError'
  }
}

export class NotFoundError extends Error {
  constructor(public message: string = 'Resource not found') {
    super(message)
    this.name = 'NotFoundError'
  }
}

export class BackendServerInternalError extends Error {
  constructor(public message: string = 'Backend API had an server internal error') {
    super(message)
    this.name = 'BackendServerInternalError'
  }
}

export class UnauthorizedError extends Error {
  constructor(public message: string = 'Request unauthorized') {
    super(message)
    this.name = 'UnauthorizedError'
  }
}

export class ForbiddenError extends Error {
  constructor(public message: string = 'Request forbidden') {
    super(message)
    this.name = 'ForbiddenError'
  }
}

export class BadRequestError extends Error {
  constructor(public message: string = 'The request was not accepted') {
    super(message)
    this.name = 'BadRequestError'
  }
}

export class UnexpectedStatusCodeError extends Error {
  constructor(public message: string = 'Unexpected HTTP status code received') {
    super(message)
    this.name = 'UnexpectedStatusCodeError'
  }
}

export class UnknownRequestFailureError extends Error {
  constructor(public message: string = 'Unknown request failure') {
    super(message)
    this.name = 'UnknownRequestFailureError'
  }
}

export class LoginError extends Error {
  constructor(public message: string = 'Unable to login') {
    super(message)
    this.name = 'LoginError'
  }
}

export type RequestFailureError =
  | BadRequestError
  | UnauthorizedError
  | ForbiddenError
  | NotFoundError
  | BackendServerInternalError
  | UnexpectedStatusCodeError
  | NetworkError
  | UnknownRequestFailureError

export function handleRequestFailure(error: unknown): RequestFailureError {
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
