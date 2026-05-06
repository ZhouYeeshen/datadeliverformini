import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useAdminStore = defineStore('admin', () => {
  const token = ref(localStorage.getItem('admin_token') || '')
  const username = ref(localStorage.getItem('admin_username') || '')

  function setAuth(t: string, name: string) {
    token.value = t
    username.value = name
    localStorage.setItem('admin_token', t)
    localStorage.setItem('admin_username', name)
  }

  function logout() {
    token.value = ''
    username.value = ''
    localStorage.removeItem('admin_token')
    localStorage.removeItem('admin_username')
  }

  return { token, username, setAuth, logout }
})
