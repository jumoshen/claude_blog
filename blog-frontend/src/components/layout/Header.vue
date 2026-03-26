<template>
  <header class="header" :class="'theme-' + styleStore.currentTheme">
    <div class="container">
      <div class="logo" @click="$router.push('/')">
        <img :src="styleStore.theme.logo" :alt="styleStore.theme.name" class="logo-img" />
        <span class="logo-text">{{ styleStore.theme.name }}</span>
      </div>
      <nav class="nav">
        <router-link to="/" class="nav-link">
          <span class="link-icon">🏠</span>
          <span class="link-text">Home</span>
        </router-link>
        <router-link to="/archives" class="nav-link">
          <span class="link-icon">📁</span>
          <span class="link-text">Archives</span>
        </router-link>
        <router-link to="/about" class="nav-link">
          <span class="link-icon">💡</span>
          <span class="link-text">About</span>
        </router-link>

        <el-dropdown @command="handleStyleChange" trigger="click" class="theme-dropdown">
          <span class="style-switcher">
            <span class="switcher-icon">🎨</span>
            <span class="switcher-text">{{ styleStore.theme.name }}</span>
            <span class="arrow">▼</span>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item
                v-for="(t, key) in styleStore.themes"
                :key="key"
                :command="key"
                :class="{ active: key === styleStore.currentTheme }"
              >
                <img :src="t.logo" class="theme-logo" :alt="t.name" />
                {{ t.name }}
                <span v-if="key === styleStore.currentTheme" class="check">✓</span>
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>

        <template v-if="userStore.isLoggedIn">
          <el-dropdown @command="handleCommand">
            <span class="user">
              <img :src="userStore.user?.avatar_url" class="avatar" />
              <span class="username">{{ userStore.user?.login }}</span>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="logout">登出</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </template>
        <router-link v-else to="/login" class="nav-link login-link">
          <span class="link-icon">🔑</span>
          <span class="link-text">登录</span>
        </router-link>
      </nav>
    </div>
  </header>
</template>

<script setup>
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../../store/user'
import { useStyleStore } from '../../store/style'
import api from '../../api'

const router = useRouter()
const userStore = useUserStore()
const styleStore = useStyleStore()

onMounted(async () => {
  try {
    const res = await api.getSite()
    if (res.code === 0) {
      // siteTitle.value = res.data.title
    }
  } catch (e) {
    // use default
  }
})

const handleStyleChange = (theme) => {
  styleStore.setTheme(theme)
}

const handleCommand = (command) => {
  if (command === 'logout') {
    userStore.logout()
    router.push('/')
  }
}
</script>

<style scoped>
.header {
  background: var(--header-bg);
  backdrop-filter: blur(10px);
  box-shadow: 0 2px 8px rgba(0,0,0,0.06);
  position: sticky;
  top: 0;
  z-index: 100;
  border-bottom: 1px solid var(--border);
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
  display: flex;
  align-items: center;
  height: 64px;
  gap: 24px;
}

.logo {
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 10px;
  transition: transform 0.3s;
}

.logo:hover {
  transform: scale(1.05);
}

.logo-img {
  height: 40px;
  width: auto;
}

.logo-text {
  font-weight: 600;
  color: var(--text-h);
  font-size: 16px;
}

.nav {
  margin-left: auto;
  display: flex;
  gap: 8px;
  align-items: center;
}

.nav-link {
  display: flex;
  align-items: center;
  gap: 6px;
  color: var(--text);
  text-decoration: none;
  padding: 8px 12px;
  border-radius: 8px;
  transition: all 0.2s;
  font-size: 14px;
}

.nav-link:hover {
  color: var(--accent);
  background: var(--accent-bg);
}

.nav-link.router-link-active {
  color: var(--accent);
  background: var(--accent-bg);
}

.link-icon {
  font-size: 14px;
}

.login-link {
  background: var(--accent-bg);
  color: var(--accent);
}

.login-link:hover {
  background: var(--accent);
  color: white;
}

.theme-dropdown {
  margin-left: 8px;
}

.style-switcher {
  display: flex;
  align-items: center;
  gap: 6px;
  color: var(--accent);
  cursor: pointer;
  padding: 8px 12px;
  border-radius: 8px;
  transition: all 0.2s;
  font-size: 14px;
  background: var(--accent-bg);
}

.style-switcher:hover {
  background: var(--accent);
  color: white;
}

.switcher-icon {
  font-size: 14px;
}

.arrow {
  font-size: 10px;
  transition: transform 0.2s;
}

.style-switcher:hover .arrow {
  transform: rotate(180deg);
}

.theme-logo {
  width: 20px;
  height: 20px;
  margin-right: 8px;
  vertical-align: middle;
}

.check {
  margin-left: auto;
  color: var(--accent);
}

:deep(.el-dropdown-menu__item) {
  display: flex;
  align-items: center;
}

:deep(.el-dropdown-menu__item.active) {
  color: var(--accent);
  background: var(--accent-bg);
}

.user {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 8px;
  transition: background 0.2s;
}

.user:hover {
  background: var(--accent-bg);
}

.avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  border: 2px solid var(--border);
}

.username {
  color: var(--text);
  font-size: 14px;
}

/* 像素风特殊效果 */
.theme-pixel .logo-img {
  image-rendering: pixelated;
}

/* 可爱风特殊效果 */
.theme-cute .logo:hover {
  animation: cuteLogo 0.5s ease;
}

@keyframes cuteLogo {
  0%, 100% { transform: scale(1); }
  50% { transform: scale(1.1) rotate(5deg); }
}

/* Q版特殊效果 */
.theme-qver .header {
  background: rgba(255, 255, 255, 0.8);
}

@media (max-width: 768px) {
  .logo-text, .link-text, .switcher-text, .username {
    display: none;
  }
  .container {
    gap: 12px;
  }
  .nav-link {
    padding: 8px;
  }
}
</style>
