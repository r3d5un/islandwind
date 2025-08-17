import { marked } from 'marked'

export interface IBlogpostResponse {
  data: IBlogpost
}

export interface IBlogpostListResponse {
  data: IBlogpost[]
  metadata: IMetadata
}

export class BlogpostPostBody {
  public data: BlogpostInput

  constructor(input: BlogpostInput) {
    this.data = input
  }
}

export class BlogpostInput {
  public title: string
  public content: string
  public published: boolean

  constructor(title: string, content: string, published: boolean) {
    this.title = title
    this.content = content
    this.published = published
  }
}

export class BlogpostListResponse {
  public data: Blogpost[]
  public metadata: Metadata

  constructor(data: IBlogpost[], metadata: IMetadata) {
    this.data = []
    for (const input of data) {
      this.data.push(new Blogpost(input))
    }
    this.metadata = new Metadata(metadata)
  }
}

export interface IBlogpost {
  id: string
  title: string
  content: string
  published: boolean
  createdAt: Date
  updatedAt: Date
  deleted: boolean
  deletedAt: Date
}

export class Blogpost {
  id: string
  title: string
  content: string
  published: boolean
  createdAt: Date
  updatedAt: Date
  deleted: boolean
  deletedAt: Date

  constructor(blogpost: IBlogpost) {
    this.id = blogpost.id
    this.title = blogpost.title
    this.content = blogpost.content
    this.published = blogpost.published
    this.createdAt = new Date(blogpost.createdAt)
    this.updatedAt = new Date(blogpost.updatedAt)
    this.deleted = blogpost.deleted
    this.deletedAt = new Date(blogpost.deletedAt)
  }

  public async markdownContent(): Promise<string> {
    return marked(this.content)
  }
}

export class BlogpostDeleteBody {
  public data: BlogpostDeleteOptions

  constructor(data: BlogpostDeleteOptions) {
    this.data = data
  }
}

export class BlogpostDeleteOptions {
  id: string
  purge: boolean

  constructor(id: string, purge: boolean) {
    this.id = id
    this.purge = purge
  }
}

export class BlogpostPatchBody {
  public data: BlogpostPatch

  constructor(data: BlogpostPatch) {
    this.data = data
  }
}

export class BlogpostPatch {
  public id: string
  public title?: string | null = null
  public content?: string | null = null
  public published?: boolean | null = null
  public deleted?: boolean | null = null

  constructor(input: IBlogpostPatchNamedParameters) {
    this.id = input.id
    this.title = input.title
    this.content = input.content
    this.published = input.published
    this.deleted = input.deleted
  }
}

export interface IBlogpostPatchNamedParameters {
  id: string
  title?: string
  content?: string
  published?: boolean
  deleted?: boolean
}

export interface IMetadata {
  lastSeen: string
  next: string
  responseLength: number
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
