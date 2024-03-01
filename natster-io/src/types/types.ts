interface Catalog {
  selected: bool
  description: string
  image: string
  online: bool
  lastSeen: Date
  to: string
  from: string
  name: string
  pending_invite: bool
  files: File[]
  status: Date
}

interface File {
  byte_size: int
  description: string
  hash: string
  mime_type: string
  path: string
}

export { Catalog, File }
