interface Catalog {
  selected: boolean
  description: string
  image: string
  online: boolean
  lastSeen: Date
  to: string
  from: string
  name: string
  pending_invite: boolean
  files: File[]
  status: Date
}

interface File {
  byte_size: number
  description: string
  hash: string
  mime_type: string
  path: string
}

export type { Catalog, File }
