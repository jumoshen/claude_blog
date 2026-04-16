import { defineStore } from 'pinia'
import api from '../api'

export const useUserStore = defineStore('user', {
  state: () => ({
    user: JSON.parse(localStorage.getItem('user') || 'null'),
    token: localStorage.getItem('token') || null,
    isLoggedIn: !!localStorage.getItem('token'),
  }),

  actions: {
    setToken(token) {
      this.token = token
      this.isLoggedIn = true
      localStorage.setItem('token', token)
      // Fetch user info with the token
      this.fetchMe()
    },

    setUser(user, token) {
      this.user = user
      this.token = token
      this.isLoggedIn = true
      localStorage.setItem('user', JSON.stringify(user))
      localStorage.setItem('token', token)
    },

    async login(code) {
      const res = await api.loginCallback(code)
      if (res.code === 0) {
        this.setUser(res.data.user, res.data.token)
        return true
      }
      return false
    },

    async logout() {
      try {
        await api.logout()
      } catch (e) {
        // ignore
      }
      this.user = null
      this.token = null
      this.isLoggedIn = false
      localStorage.removeItem('user')
      localStorage.removeItem('token')
    },

    async fetchMe() {
      if (!this.token) return null
      try {
        const res = await api.getMe()
        if (res.code === 0) {
          this.user = res.data
          localStorage.setItem('user', JSON.stringify(this.user))
          return this.user
        }
      } catch (e) {
        this.logout()
      }
      return null
    },
  },
})
