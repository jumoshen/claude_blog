<template>
  <div class="login">
    <div class="login-box">
      <h1>Login</h1>
      <p>Sign in with GitHub to access admin features</p>
      <el-button type="primary" @click="login" :loading="loading">
        Login with GitHub
      </el-button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import api from '../api'
import { useUserStore } from '../store/user'

const router = useRouter()
const userStore = useUserStore()
const loading = ref(false)

onMounted(async () => {
  // Check if already logged in
  if (userStore.isLoggedIn) {
    router.push('/')
  }

  // Handle OAuth callback with token (redirected from API)
  const token = new URLSearchParams(window.location.search).get('token')
  if (token) {
    userStore.setToken(token)
    ElMessage.success('Login successful!')
    router.push('/')
    return
  }

  // Handle OAuth callback with code (direct GitHub redirect)
  const code = new URLSearchParams(window.location.search).get('code')
  if (code) {
    loading.value = true
    const success = await userStore.login(code)
    loading.value = false
    if (success) {
      ElMessage.success('Login successful!')
      router.push('/')
    } else {
      ElMessage.error('Login failed')
    }
  }
})

const login = async () => {
  loading.value = true
  const res = await api.getLoginInfo()
  if (res.code === 0) {
    // Redirect to GitHub OAuth
    const state = res.data.state
    window.location.href = `https://github.com/login/oauth/authorize?client_id=${res.data.client_id}&redirect_uri=${encodeURIComponent(res.data.callback_url)}&state=${state}`
  }
}
</script>

<style scoped>
.login {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 60vh;
}
.login-box {
  text-align: center;
  background: #fff;
  padding: 40px;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0,0,0,0.1);
}
.login-box h1 { margin: 0 0 10px; }
.login-box p { color: #666; margin-bottom: 20px; }
</style>
