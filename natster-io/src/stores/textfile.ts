import { defineStore } from 'pinia'

export const textFileStore = defineStore('textfile', {
  state: () => ({
    body: '',
    title: '',
    show: false
  }),
  actions: {
    showTextFile(title, body) {
      this.title = title
      this.body = body
      this.show = true
    }
  }
})
