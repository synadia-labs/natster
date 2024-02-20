import { defineStore } from 'pinia'

export const notificationStore = defineStore('notification', {
  state: () => ({
    message: '',
    title: '',
    show: false
  }),
  actions: {
    setNotification(title, message) {
      console.log('In setNotification', title, message)
      this.title = title
      this.message = message
      this.show = true
      setTimeout(() => {
        this.show = false
        this.message = ''
        this.title = ''
      }, 3000)
    }
  }
})
