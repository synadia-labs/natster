interface Catalog {
  selected: bool
  online: bool
  to: string
  from: string
  name: string
  files: File[]
}

interface File {
  byte_size: int
  description: string
  hash: string
  mime_type: string
  path: string
}

export { Catalog }
