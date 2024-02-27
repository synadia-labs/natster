import { defineStore } from 'pinia'

export const fileStore = defineStore('file', {
  state: () => ({
    body: null,
    buffer: [],
    title: '',
    show: false,
    mimeType: null,
    codec: null,

    mediaSource: null,
    mediaUrl: null,
    updatingBuffer: false,

    audioSourceBuffer: null,
    videoSourceBuffer: null,

    appendCount: 0
  }),
  actions: {
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

          const _data = data;
          this.mediaSource.addEventListener("sourceopen", () => {
            this.videoSourceBuffer = this.mediaSource.addSourceBuffer(`video/mp4; codecs="${this.codec}"`)

            this.videoSourceBuffer.addEventListener('updatestart', (e) => {
              // this.updatingBuffer = true;
              console.log(e)
            })

            this.videoSourceBuffer.addEventListener('update', (e) => {
              console.log(e)
            })

            this.videoSourceBuffer.addEventListener('updateend', (e) => {
              // this.updatingBuffer = false;
              console.log(e)
            })

            this.videoSourceBuffer.addEventListener('error', (e) => {
              console.log(e)
            })

            this.videoSourceBuffer.addEventListener('abort', (e) => {
              console.log(e)
            })
          })

          this.mediaSource.addEventListener("sourceended", (e) => {
            console.log(e)
          })

          this.mediaSource.addEventListener("sourceclose", (e) => {
            console.log(e)
          })

          this.mediaSource.addEventListener("error", (e) => {
            console.log(e)
          })
        }

        let interval = null
        interval = setInterval(() => {
          if (this.videoSourceBuffer && !this.videoSourceBuffer.updating && this.buffer.length > 0) {
            this.videoSourceBuffer.appendBuffer(this.buffer.shift());

            this.appendCount++
            console.log(`append count: ${this.appendCount}`)
          }
        }, 10)

        this.buffer.push(data)
        console.log(`buffer length: ${this.buffer.length}`)
      } else {
        this.body = data
      }
    }
  }
})
