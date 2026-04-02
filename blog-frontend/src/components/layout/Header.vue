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

        <el-dropdown @command="handleStyleChange" trigger="click" class="theme-dropdown" :popper-class="'theme-dropdown-panel theme-' + styleStore.currentTheme">
          <span class="style-switcher" :class="'theme-' + styleStore.currentTheme">
            <span class="switcher-icon">{{ styleStore.theme.icon }}</span>
            <span class="switcher-text">{{ styleStore.theme.name }}</span>
            <span class="arrow">▼</span>
          </span>
          <template #dropdown>
            <el-dropdown-menu class="theme-menu">
              <el-dropdown-item
                v-for="(t, key) in styleStore.themes"
                :key="key"
                :command="key"
                :class="{ active: key === styleStore.currentTheme, ['theme-' + key]: true }"
              >
                <div class="theme-card">
                  <div class="theme-preview" :style="{ background: t.preview }">
                    <span class="theme-icon">{{ t.icon }}</span>
                  </div>
                  <div class="theme-info">
                    <span class="theme-name">{{ t.name }}</span>
                    <span class="theme-tagline">{{ t.tagline }}</span>
                  </div>
                  <span v-if="key === styleStore.currentTheme" class="check">✓</span>
                </div>
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>

        <template v-if="true">
          <el-dropdown @command="handleCommand" trigger="click">
            <span class="user">
              <img src="https://avatars.githubusercontent.com/u/1?v=4" class="avatar" />
              <span class="username">testuser</span>
              <span class="user-arrow">▾</span>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="logout">
                  <span class="menu-icon">🚪</span>
                  退出登录
                </el-dropdown-item>
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
  border-radius: 12px;
  transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
  font-size: 14px;
  background: var(--accent-bg);
  border: 2px solid transparent;
}

.style-switcher:hover {
  transform: scale(1.05);
  border-color: var(--accent);
  box-shadow: 0 4px 20px var(--shadow);
}

.style-switcher .arrow {
  font-size: 10px;
  transition: transform 0.3s;
}

.style-switcher:hover .arrow {
  transform: rotate(180deg);
}

/* 主题菜单样式 */
:deep(.theme-menu) {
  padding: 8px !important;
  background: var(--card-bg) !important;
  border: 1px solid var(--border) !important;
  border-radius: 16px !important;
  box-shadow: 0 20px 60px var(--shadow) !important;
}

:deep(.el-dropdown-menu__item) {
  padding: 0 !important;
  margin: 4px 0;
  border-radius: 12px;
  overflow: hidden;
}

:deep(.el-dropdown-menu__item .theme-card) {
  display: flex;
  align-items: center;
  padding: 12px;
  gap: 12px;
  transition: all 0.3s;
}

:deep(.el-dropdown-menu__item:hover .theme-card) {
  transform: translateX(4px);
}

:deep(.el-dropdown-menu__item.active) {
  background: var(--accent-bg) !important;
}

.theme-preview {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  box-shadow: 0 4px 12px rgba(0,0,0,0.15);
  transition: transform 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
}

.theme-preview .theme-icon {
  font-size: 24px;
  filter: drop-shadow(0 2px 4px rgba(0,0,0,0.2));
}

:deep(.el-dropdown-menu__item:hover .theme-preview) {
  transform: scale(1.1) rotate(-5deg);
}

.theme-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.theme-name {
  font-weight: 600;
  color: var(--text-h);
  font-size: 14px;
}

.theme-tagline {
  font-size: 10px;
  color: var(--text);
  letter-spacing: 1px;
  text-transform: uppercase;
  opacity: 0.7;
}

.check {
  margin-left: auto;
  color: var(--accent);
  font-size: 16px;
  font-weight: bold;
}

/* 像素风主题卡片悬停 */
:deep(.theme-pixel:hover .theme-preview) {
  background: linear-gradient(135deg, #b366ff 0%, #7c3aed 100%) !important;
}

/* 可爱风主题卡片悬停 */
:deep(.theme-cute:hover .theme-preview) {
  background: linear-gradient(135deg, #ff6b9d 0%, #ffa8c5 100%) !important;
}

/* Q版主题卡片悬停 */
:deep(.theme-qver:hover .theme-preview) {
  background: linear-gradient(135deg, #7c6aff 0%, #a78bfa 100%) !important;
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

.user-arrow {
  font-size: 10px;
  color: var(--text);
  opacity: 0.6;
  margin-left: 2px;
  transition: transform 0.2s;
}

.user:hover .user-arrow {
  transform: rotate(180deg);
}

.menu-icon {
  font-size: 14px;
}

/* 用户下拉菜单 - Glassmorphism */
.user {
  position: relative;
}

:deep(.el-dropdown-menu) {
  padding: 8px !important;
  background: var(--card-bg) !important;
  border: 1px solid var(--border) !important;
  border-radius: 16px !important;
  box-shadow:
    0 8px 32px var(--shadow),
    0 0 0 1px rgba(255,255,255,0.1) inset !important;
  backdrop-filter: blur(20px) saturate(180%) !important;
  -webkit-backdrop-filter: blur(20px) saturate(180%) !important;
  min-width: 140px !important;
}

:deep(.el-dropdown-menu__item) {
  padding: 12px 16px !important;
  margin: 2px 0;
  border-radius: 10px !important;
  font-size: 14px !important;
  color: var(--text) !important;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1) !important;
  display: flex !important;
  align-items: center !important;
  gap: 10px !important;
}

:deep(.el-dropdown-menu__item::before) {
  content: '';
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--accent);
  opacity: 0;
  transform: scale(0);
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1) !important;
}

:deep(.el-dropdown-menu__item:hover) {
  background: var(--accent-bg) !important;
  color: var(--accent) !important;
  padding-left: 20px !important;
}

:deep(.el-dropdown-menu__item:hover::before) {
  opacity: 1;
  transform: scale(1);
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
