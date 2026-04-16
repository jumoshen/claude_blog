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
      <p>登录后即可点赞、收藏和评论文章</p>

      <!-- 账号密码登录表单 -->
      <form class="login-form" @submit.prevent="handleLogin">
        <div class="form-group">
          <input
            v-model="loginForm.username"
            type="text"
            placeholder="用户名"
            class="form-input"
            autocomplete="username"
          />
        </div>
        <div class="form-group">
          <input
            v-model="loginForm.password"
            type="password"
            placeholder="密码"
            class="form-input"
            autocomplete="current-password"
          />
        </div>
        <div class="form-group captcha-group">
          <input
            v-model="loginForm.captcha"
            type="text"
            placeholder="验证码"
            class="form-input captcha-input"
            maxlength="4"
            autocomplete="off"
          />
          <div class="captcha-img" @click="refreshCaptcha" title="点击刷新">
            <img v-if="captchaImg" :src="captchaImg" alt="验证码" />
            <span v-else class="captcha-loading">加载中...</span>
          </div>
        </div>
        <p v-if="errorMsg" class="error-msg">{{ errorMsg }}</p>
        <button type="submit" class="login-btn" :disabled="loading">
          <span v-if="loading" class="loading-spinner"></span>
          <span v-else>登录</span>
        </button>
      </form>

      <div class="divider">
        <span>或</span>
      </div>

      <!-- GitHub登录 -->
      <button class="github-btn" @click="loginWithGithub" :disabled="loading">
        <svg class="github-icon" width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
          <path d="M12 0C5.37 0 0 5.37 0 12c0 5.31 3.435 9.795 8.205 11.385.6.105.825-.255.825-.57 0-.285-.015-1.23-.015-2.235-3.015.555-3.795-.735-4.035-1.41-.135-.345-.72-1.41-1.23-1.695-.42-.225-1.02-.78-.015-.795.945-.015 1.62.87 1.845 1.23 1.08 1.815 2.805 1.305 3.495.99.105-.78.42-1.305.765-1.605-2.67-.3-5.46-1.335-5.46-5.925 0-1.305.465-2.385 1.23-3.225-.12-.3-.54-1.53.12-3.18 0 0 1.005-.315 3.3 1.23.96-.27 1.98-.405 3-.405s2.04.135 3 .405c2.295-1.56 3.3-1.23 3.3-1.23.66 1.65.24 2.88.12 3.18.765.84 1.23 1.905 1.23 3.225 0 4.605-2.805 5.625-5.475 5.925.435.375.81 1.095.81 2.22 0 1.605-.015 2.895-.015 3.3 0 .315.225.69.825.57A12.02 12.02 0 0024 12c0-6.63-5.37-12-12-12z"/>
        </svg>
        <span>使用 GitHub 登录</span>
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
const errorMsg = ref('')
const captchaImg = ref('')
const captchaId = ref('')

const loginForm = ref({
  username: '',
  password: '',
  captcha: ''
})

onMounted(async () => {
  // Check if already logged in
  if (userStore.isLoggedIn) {
    router.push('/')
    return
  }

  // Handle OAuth callback with token (redirected from API)
  const token = new URLSearchParams(window.location.search).get('token')
  if (token) {
    userStore.setToken(token)
    ElMessage.success('登录成功!')
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
      ElMessage.success('登录成功!')
      router.push('/')
    } else {
      ElMessage.error('登录失败')
    }
    return
  }

  // 获取验证码
  await refreshCaptcha()
})

// 获取验证码
const refreshCaptcha = async () => {
  try {
    const res = await api.getCaptcha()
    if (res.code === 0) {
      captchaImg.value = res.data.image
      captchaId.value = res.data.id
    }
  } catch (e) {
    console.error('Failed to get captcha:', e)
  }
}

// 账号密码登录
const handleLogin = async () => {
  if (!loginForm.value.username || !loginForm.value.password || !loginForm.value.captcha) {
    errorMsg.value = '请填写完整信息'
    return
  }

  loading.value = true
  errorMsg.value = ''

  try {
    const res = await api.loginWithPassword({
      username: loginForm.value.username,
      password: loginForm.value.password,
      captcha: loginForm.value.captcha,
      captcha_id: captchaId.value
    })

    if (res.code === 0) {
      userStore.setUser(res.data.user, res.data.token)
      ElMessage.success('登录成功!')
      router.push('/')
    } else {
      errorMsg.value = res.message || '登录失败'
      await refreshCaptcha()
    }
  } catch (e) {
    errorMsg.value = '登录失败，请稍后重试'
    await refreshCaptcha()
  } finally {
    loading.value = false
  }
}

// GitHub登录
const loginWithGithub = async () => {
  loading.value = true
  const res = await api.getLoginInfo()
  if (res.code === 0) {
    const state = res.data.state
    window.location.href = `https://github.com/login/oauth/authorize?client_id=${res.data.client_id}&redirect_uri=${encodeURIComponent(res.data.callback_url)}&state=${state}`
  } else {
    loading.value = false
    ElMessage.error('获取登录信息失败')
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
  width: 100%;
  max-width: 400px;
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

.login-box > p {
  color: var(--text);
  margin-bottom: 24px;
  font-size: 14px;
}

/* 登录表单 */
.login-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-bottom: 20px;
}

.form-group {
  position: relative;
}

.form-input {
  width: 100%;
  padding: 14px 16px;
  border: 2px solid var(--border);
  border-radius: 12px;
  font-size: 15px;
  background: var(--bg);
  color: var(--text);
  transition: all 0.2s;
  box-sizing: border-box;
}

.form-input:focus {
  outline: none;
  border-color: var(--accent);
  box-shadow: 0 0 0 4px var(--accent-bg);
}

.form-input::placeholder {
  color: var(--text);
  opacity: 0.5;
}

/* 验证码 */
.captcha-group {
  display: flex;
  gap: 12px;
}

.captcha-input {
  flex: 1;
}

.captcha-img {
  width: 120px;
  height: 48px;
  border-radius: 12px;
  overflow: hidden;
  cursor: pointer;
  border: 2px solid var(--border);
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg);
}

.captcha-img:hover {
  border-color: var(--accent);
  transform: scale(1.02);
}

.captcha-img img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.captcha-loading {
  font-size: 12px;
  color: var(--text);
  opacity: 0.6;
}

.error-msg {
  color: #ff6b6b;
  font-size: 13px;
  margin: 0;
  text-align: left;
}

/* 登录按钮 */
.login-btn {
  width: 100%;
  padding: 14px;
  background: var(--accent);
  color: #fff;
  border: none;
  border-radius: 12px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.login-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px var(--shadow);
}

.login-btn:active:not(:disabled) {
  transform: translateY(0);
}

.login-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

/* 分隔线 */
.divider {
  display: flex;
  align-items: center;
  gap: 16px;
  margin: 20px 0;
  color: var(--text);
  opacity: 0.6;
  font-size: 13px;
}

.divider::before,
.divider::after {
  content: '';
  flex: 1;
  height: 1px;
  background: var(--border);
}

/* GitHub按钮 */
.github-btn {
  width: 100%;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  padding: 14px;
  background: var(--card-bg);
  color: var(--text);
  border: 2px solid var(--border);
  border-radius: 12px;
  font-size: 15px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
}

.github-btn:hover:not(:disabled) {
  border-color: var(--accent);
  background: var(--accent-bg);
  transform: translateY(-2px);
}

.github-btn:active:not(:disabled) {
  transform: translateY(0);
}

.github-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
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
