import { defineStore } from 'pinia'

export const fileStore = defineStore('file', {
  state: () => ({
    body: null,
    buffer: [],
    loading: true,
    title: '',
    show: false,
    mimeType: null,
    codec: null,

    mediaSource: null,
    mediaUrl: null,

    audioSourceBuffer: null,
    videoSourceBuffer: null,

    appendCount: 0,
    appendInterval: null
  }),
  actions: {
    endStream() {
      if (this.mediaSource) {
        this.mediaSource.endOfStream()
        console.log('stream ended')
      }
    },
    render(title, mimeType, data) {
      this.title = title
      this.show = true
      this.mimeType = mimeType

      if (mimeType.toLowerCase().indexOf('video/') === 0) {
        if (!this.mediaSource && !this.mediaUrl && !this.videoSourceBuffer) {
          this.mediaSource = new MediaSource()
          this.mediaUrl = URL.createObjectURL(this.mediaSource)

          this.codec = 'avc1.640028,mp4a.40.2' //'avc1.42C028,mp4a.40.2' // FIXME-- read this from headers and pass it in to render()
          console.log(MediaSource.isTypeSupported(`video/mp4; codecs="${this.codec}"`))

          const _data = data
          this.mediaSource.addEventListener('sourceopen', () => {
            this.videoSourceBuffer = this.mediaSource.addSourceBuffer(
              `video/mp4; codecs="${this.codec}"`
            )

            this.videoSourceBuffer.addEventListener('updatestart', (e) => {
              // console.log(e)
            })

            this.videoSourceBuffer.addEventListener('update', (e) => {
              // console.log(e)
            })

            this.videoSourceBuffer.addEventListener('updateend', (e) => {
              // console.log(e)
            })

            this.videoSourceBuffer.addEventListener('error', (e) => {
              console.log(e)
            })

            this.videoSourceBuffer.addEventListener('abort', (e) => {
              // console.log(e)
            })
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

          this.mediaSource.addEventListener('sourceclose', (e) => {
            // console.log(e)
          })

          this.mediaSource.addEventListener('error', (e) => {
            // console.log(e)
          })
        }

        this.appendInterval = setInterval(() => {
          if (
            this.videoSourceBuffer &&
            !this.videoSourceBuffer.updating &&
            this.buffer.length > 0
          ) {
            this.videoSourceBuffer.appendBuffer(this.buffer.shift())

            this.appendCount++
            // console.log(`append count: ${this.appendCount}`)
          }
        }, 10)

        this.buffer.push(data)
        // console.log(`buffer length: ${this.buffer.length}`)
      } else {
        this.body = data
      }
    },
    reset() {
      if (this.appendInterval) {
        clearInterval(this.appendInterval)
        this.appendInterval = null
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

      console.log('reset!')
    }
  }
})
