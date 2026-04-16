<template>
  <div class="login">
    <!-- 动态渐变背景 -->
    <div class="gradient-bg">
      <div class="gradient-sphere sphere-1"></div>
      <div class="gradient-sphere sphere-2"></div>
      <div class="gradient-sphere sphere-3"></div>
      <div class="floating-shapes">
        <div class="shape shape-1"></div>
        <div class="shape shape-2"></div>
        <div class="shape shape-3"></div>
        <div class="shape shape-4"></div>
      </div>
    </div>

    <div class="login-box">
      <div class="login-icon">🔐</div>
      <h1>Welcome Back</h1>
      <p>Sign in with GitHub to access admin features</p>
      <button class="github-btn" @click="login" :disabled="loading">
        <span v-if="loading" class="loading-spinner"></span>
        <template v-else>
          <svg class="github-icon" width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
            <path d="M12 0C5.37 0 0 5.37 0 12c0 5.31 3.435 9.795 8.205 11.385.6.105.825-.255.825-.57 0-.285-.015-1.23-.015-2.235-3.015.555-3.795-.735-4.035-1.41-.135-.345-.72-1.41-1.23-1.695-.42-.225-1.02-.78-.015-.795.945-.015 1.62.87 1.845 1.23 1.08 1.815 2.805 1.305 3.495.99.105-.78.42-1.305.765-1.605-2.67-.3-5.46-1.335-5.46-5.925 0-1.305.465-2.385 1.23-3.225-.12-.3-.54-1.53.12-3.18 0 0 1.005-.315 3.3 1.23.96-.27 1.98-.405 3-.405s2.04.135 3 .405c2.295-1.56 3.3-1.23 3.3-1.23.66 1.65.24 2.88.12 3.18.765.84 1.23 1.905 1.23 3.225 0 4.605-2.805 5.625-5.475 5.925.435.375.81 1.095.81 2.22 0 1.605-.015 2.895-.015 3.3 0 .315.225.69.825.57A12.02 12.02 0 0024 12c0-6.63-5.37-12-12-12z"/>
          </svg>
          <span>Login with GitHub</span>
        </template>
      </button>
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
  position: relative;
  overflow: hidden;
}

/* 动态渐变背景 */
.gradient-bg {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: linear-gradient(-45deg, var(--bg), var(--card-bg), var(--header-bg), var(--bg));
  background-size: 400% 400%;
  animation: gradientShift 15s ease infinite;
  z-index: 0;
}

@keyframes gradientShift {
  0% { background-position: 0% 50%; }
  50% { background-position: 100% 50%; }
  100% { background-position: 0% 50%; }
}

/* 浮动光球 */
.gradient-sphere {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  opacity: 0.4;
  animation: floatSphere 10s ease-in-out infinite;
}

.sphere-1 {
  width: 400px;
  height: 400px;
  background: linear-gradient(135deg, var(--accent), var(--accent-border));
  top: -150px;
  right: -100px;
}

.sphere-2 {
  width: 300px;
  height: 300px;
  background: linear-gradient(135deg, var(--accent), var(--accent-border));
  bottom: -100px;
  left: -100px;
  animation-delay: -5s;
}

.sphere-3 {
  width: 200px;
  height: 200px;
  background: linear-gradient(135deg, var(--accent), var(--accent-border));
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  animation-delay: -2.5s;
}

@keyframes floatSphere {
  0%, 100% { transform: translate(0, 0) scale(1); }
  33% { transform: translate(30px, -30px) scale(1.05); }
  66% { transform: translate(-20px, 20px) scale(0.95); }
}

/* 浮动几何图形 */
.floating-shapes {
  position: absolute;
  width: 100%;
  height: 100%;
  overflow: hidden;
}

.shape {
  position: absolute;
  opacity: 0.15;
  animation: floatShape 20s ease-in-out infinite;
}

.shape-1 {
  width: 60px;
  height: 60px;
  background: linear-gradient(135deg, var(--accent), var(--accent-border));
  border-radius: 12px;
  top: 20%;
  left: 10%;
  animation-delay: 0s;
}

.shape-2 {
  width: 40px;
  height: 40px;
  background: linear-gradient(135deg, var(--accent), var(--accent-border));
  border-radius: 50%;
  top: 60%;
  right: 15%;
  animation-delay: -5s;
}

.shape-3 {
  width: 80px;
  height: 80px;
  background: linear-gradient(135deg, var(--accent), var(--accent-border));
  border-radius: 8px;
  bottom: 20%;
  left: 20%;
  animation-delay: -10s;
}

.shape-4 {
  width: 50px;
  height: 50px;
  background: linear-gradient(135deg, var(--accent), var(--accent-border));
  border-radius: 50%;
  top: 30%;
  right: 25%;
  animation-delay: -15s;
}

@keyframes floatShape {
  0%, 100% { transform: translate(0, 0) rotate(0deg); }
  25% { transform: translate(20px, -30px) rotate(90deg); }
  50% { transform: translate(-10px, 20px) rotate(180deg); }
  75% { transform: translate(30px, 10px) rotate(270deg); }
}

/* 登录卡片 */
.login-box {
  position: relative;
  z-index: 1;
  text-align: center;
  background: var(--card-bg);
  backdrop-filter: blur(20px);
  padding: 48px;
  border-radius: 24px;
  box-shadow: 0 25px 80px var(--shadow);
  animation: cardAppear 0.8s cubic-bezier(0.34, 1.56, 0.64, 1);
  border: 1px solid var(--border);
}

@keyframes cardAppear {
  from {
    opacity: 0;
    transform: translateY(40px) scale(0.95);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

.login-icon {
  font-size: 48px;
  margin-bottom: 16px;
  animation: iconBounce 2s ease-in-out infinite;
}

@keyframes iconBounce {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-10px); }
}

.login-box h1 {
  margin: 0 0 8px;
  font-size: 28px;
  font-weight: 700;
  color: var(--text-h);
}

.login-box p {
  color: var(--text);
  margin-bottom: 24px;
  font-size: 14px;
}

.login-box p {
  color: #666;
  margin-bottom: 24px;
  font-size: 14px;
}

/* GitHub按钮 */
.github-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  padding: 14px 28px;
  background: var(--accent);
  color: #fff;
  border: none;
  border-radius: 12px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
  box-shadow: 0 4px 20px var(--shadow);
  position: relative;
  overflow: hidden;
}

.github-btn::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.2), transparent);
  transition: left 0.5s;
}

.github-btn:hover::before {
  left: 100%;
}

.github-btn:hover {
  transform: translateY(-3px) scale(1.02);
  box-shadow: 0 8px 30px var(--shadow);
}

.github-btn:active {
  transform: translateY(-1px) scale(0.98);
}

.github-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
  transform: none;
}

.github-icon {
  width: 20px;
  height: 20px;
}

/* 加载动画 */
.loading-spinner {
  width: 20px;
  height: 20px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: #fff;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>
