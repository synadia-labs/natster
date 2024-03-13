import { defineStore } from 'pinia'

export const fileStore = defineStore('file', {
  state: () => ({
    blob: null,
    body: null,
    buffer: [],
    loading: true,
    catalog: null,
    title: '',
    description: '',
    show: false,
    mimeType: null,
    codec: null,

    mediaSource: null,
    mediaUrl: null,

    audioSourceBuffer: null,
    videoSourceBuffer: null,

    appendCount: 0,
    appendInterval: null,
    streamEndInterval: null,

    onReset: null
  }),
  actions: {
    endStream() {
      if (this.mediaSource) {
        this.streamEndInterval = setInterval(() => {
          if (this.buffer.length === 0) {
            clearInterval(this.streamEndInterval)
            this.streamEndInterval = null

            this.mediaSource?.endOfStream()
            console.log('stream ended')
          }
        }, 100)
      }
    },
    initMediaSource() {
      const re = /ipad|iphone/i
      if (navigator.userAgent.match(re) && typeof ManagedMediaSource !== 'undefined') {
        return new ManagedMediaSource()
      }
      return new MediaSource()
    },
    load(title, description, mimeType, catalog, onReset) {
      this.title = title
      this.description = description
      this.mimeType = mimeType
      this.catalog = catalog
      this.loading = true
      this.show = true
      this.onReset = onReset
    },
    render(title, description, mimeType, data, catalog, chunkIdx, totalChunks) {
      this.title = title
      this.description = description
      this.mimeType = mimeType
      this.catalog = catalog
      this.show = true

      if (mimeType.toLowerCase().indexOf('video/') === 0) {
        if (!this.mediaSource && !this.mediaUrl && !this.videoSourceBuffer) {
          this.mediaSource = this.initMediaSource()
          this.mediaUrl = URL.createObjectURL(this.mediaSource)

          this.codec = 'avc1.640028,mp4a.40.2' //'avc1.42C028,mp4a.40.2' // FIXME-- read this from headers and pass it in to render()
          console.log(MediaSource.isTypeSupported(`video/mp4; codecs="${this.codec}"`))

          this.mediaSource.addEventListener('sourceopen', () => {
            this.videoSourceBuffer = this.mediaSource.addSourceBuffer(
              `video/mp4; codecs="${this.codec}"`
            )
            this.videoSourceBuffer.addEventListener('error', (e) => {
              console.log(e)
            })

            this.videoSourceBuffer.addEventListener('abort', (e) => {})
            this.videoSourceBuffer.addEventListener('updatestart', (e) => {})
            this.videoSourceBuffer.addEventListener('update', (e) => {})
            this.videoSourceBuffer.addEventListener('updateend', (e) => {})
          })

          this.mediaSource.addEventListener('sourceended', (e) => {
            this.mediaSource = null
            this.audioSourceBuffer = null
            this.videoSourceBuffer = null

            if (this.appendInterval) {
              clearInterval(this.appendInterval)
              this.appendInterval = null
            }

            this.buffer = []
          })

          this.mediaSource.addEventListener('sourceclose', (e) => {})
          this.mediaSource.addEventListener('error', (e) => {})

          this.appendInterval = setInterval(() => {
            if (
              this.videoSourceBuffer &&
              !this.videoSourceBuffer.updating &&
              this.buffer.length > 0
            ) {
              this.videoSourceBuffer.appendBuffer(this.buffer.shift())
              this.appendCount++
            }
          }, 10)
        }

        this.buffer.push(data)
      } else if (mimeType.toLowerCase() === 'audio/mpeg') {
        if (!this.mediaSource && !this.mediaUrl && !this.audioSourceBuffer) {
          this.mediaSource = this.initMediaSource()
          this.mediaUrl = URL.createObjectURL(this.mediaSource)
          console.log(MediaSource.isTypeSupported(`audio/mpeg`))

          this.mediaSource.addEventListener('sourceopen', () => {
            this.audioSourceBuffer = this.mediaSource.addSourceBuffer(`audio/mpeg`)
            this.audioSourceBuffer.addEventListener('error', (e) => {})
            this.audioSourceBuffer.addEventListener('abort', (e) => {})
            this.audioSourceBuffer.addEventListener('updatestart', (e) => {})
            this.audioSourceBuffer.addEventListener('update', (e) => {})
            this.audioSourceBuffer.addEventListener('updateend', (e) => {})
          })

          this.mediaSource.addEventListener('sourceended', (e) => {
            this.mediaSource = null
            this.audioSourceBuffer = null
            this.videoSourceBuffer = null

            if (this.appendInterval) {
              clearInterval(this.appendInterval)
              this.appendInterval = null
            }

            this.buffer = []
          })

          const re = /sourcebuffer is full/i // FIXME? how does this work across browsers?

          this.appendInterval = setInterval(() => {
            if (
              this.audioSourceBuffer &&
              !this.audioSourceBuffer.updating &&
              this.buffer.length > 0
            ) {
              const _data = this.buffer.shift()

              try {
                this.audioSourceBuffer.appendBuffer(_data)
                this.appendCount++
              } catch (e) {
                if (e.toString().match(re)) {
                  this.buffer.unshift(_data)
                }
              }
            }
          }, 10)
        }

        this.buffer.push(data)
      } else {
        this.buffer.push(data)

        if (chunkIdx == totalChunks - 1) {
          this.blob = new Blob(this.buffer, { type: mimeType })
          this.loading = false
        }
      }
    },
    reset() {
      if (this.appendInterval) {
        clearInterval(this.appendInterval)
        this.appendInterval = null
      }

      if (this.streamEndInterval) {
        clearInterval(this.streamEndInterval)
        this.streamEndInterval = null
      }

      if (this.onReset && typeof this.onReset === 'function') {
        this.onReset()
        this.onReset = null
      }

      this.body = null
      this.buffer = []
      this.codec = null
      this.loading = true
      this.mediaSource = null
      this.mediaUrl = null
      this.mimeType = null

      this.audioSourceBuffer = null
      this.videoSourceBuffer = null

      this.show = false
      this.title = ''
      this.description = ''

      console.log('reset!')
    }
  }
})
