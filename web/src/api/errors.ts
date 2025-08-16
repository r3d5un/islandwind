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
    super()
    this.name = 'BadRequestError'
  }
}

export class UnexpectedStatusCodeError extends Error {
  constructor(public message: string = 'Unexpected HTTP status code received') {
    super()
    this.name = 'UnexpectedStatusCodeError'
  }
}

export class UnknownRequestFailureError extends Error {
  constructor(public message: string = 'Unknown request failure') {
    super()
    this.name = 'UnknownRequestFailureError'
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
