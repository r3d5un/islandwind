export interface IBlogpostResponse {
  data: IBlogpost
}

export interface IBlogpostListResponse {
  data: IBlogpost[]
  metadata: IMetadata
}

export interface IBlogpost {
  id: string
  content: string
  published: boolean
  createdAt: Date
  deleted: boolean
  deletedAt: Date
}

export interface IMetadata {
  lastSeen: string
  next: string
  responseLength: number
}

export class BlogpostResponse {
  public data: Blogpost

  constructor(data: IBlogpost) {
    this.data = data
  }
}

export class BlogpostListResponse {
  public data: Blogpost[]
  public metadata: Metadata

  constructor(data: IBlogpost[], metadata: IMetadata) {
    this.data = data
    this.metadata = metadata
  }
}

export class Blogpost {
  id: string
  content: string
  published: boolean
  createdAt: Date
  deleted: boolean
  deletedAt: Date

  constructor(blogpost: IBlogpost) {
    this.id = blogpost.id
    this.content = blogpost.content
    this.published = blogpost.published
    this.createdAt = blogpost.createdAt
    this.deleted = blogpost.deleted
    this.deletedAt = blogpost.deletedAt
  }
}

export class Metadata {
  lastSeen: string
  next: string
  responseLength: number

  constructor(metadata: IMetadata) {
    this.lastSeen = metadata.lastSeen
    this.next = metadata.next
    this.responseLength = metadata.responseLength
  }
}
