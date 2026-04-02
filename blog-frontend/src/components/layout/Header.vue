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
          <el-dropdown @command="handleCommand" trigger="click" :show-timeout="150" :hide-timeout="150">
            <span class="user">
              <div class="avatar-wrapper">
                <img src="https://avatars.githubusercontent.com/u/1?v=4" class="avatar" />
                <div class="avatar-ring"></div>
              </div>
              <span class="username">testuser</span>
              <span class="user-arrow">
                <svg width="10" height="10" viewBox="0 0 10 10" fill="currentColor">
                  <path d="M2 3.5L5 6.5L8 3.5" stroke="currentColor" stroke-width="1.5" fill="none" stroke-linecap="round"/>
                </svg>
              </span>
            </span>
            <template #dropdown>
              <div class="dropdown-wrapper">
                <el-dropdown-item command="logout" class="logout-item">
                  <span class="logout-icon">
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                      <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/>
                      <polyline points="16,17 21,12 16,7"/>
                      <line x1="21" y1="12" x2="9" y2="12"/>
                    </svg>
                  </span>
                  <span>退出登录</span>
                </el-dropdown-item>
              </div>
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
  padding: 6px 12px;
  border-radius: 24px;
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
}

.user:hover {
  background: var(--accent-bg);
}

.avatar-wrapper {
  position: relative;
  width: 32px;
  height: 32px;
}

.avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  z-index: 1;
}

.avatar-ring {
  position: absolute;
  inset: -3px;
  border-radius: 50%;
  border: 2px solid transparent;
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
}

.user:hover .avatar {
  transform: scale(1.05);
}

.user:hover .avatar-ring {
  border-color: var(--accent);
  box-shadow: 0 0 12px var(--accent);
}

.username {
  color: var(--text);
  font-size: 14px;
  font-weight: 500;
  transition: color 0.2s;
}

.user:hover .username {
  color: var(--accent);
}

.user-arrow {
  color: var(--text);
  opacity: 0.5;
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  align-items: center;
}

.user:hover .user-arrow {
  opacity: 1;
  transform: rotate(180deg);
}

.menu-icon {
  font-size: 16px;
  transition: all 0.2s;
}

/* 下拉菜单包装器 */
.dropdown-wrapper {
  background: var(--card-bg) !important;
  border: 1px solid var(--border) !important;
  border-radius: 20px !important;
  padding: 10px !important;
  min-width: 160px !important;
  box-shadow: 0 16px 48px var(--shadow) !important;
  animation: dropdownEnter 0.25s cubic-bezier(0.4, 0, 0.2, 1);
}

@keyframes dropdownEnter {
  from {
    opacity: 0;
    transform: translateY(-12px) scale(0.95);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

/* 退出按钮 */
:deep(.logout-item) {
  padding: 14px 18px !important;
  margin: 4px 0 !important;
  border-radius: 14px !important;
  font-size: 14px !important;
  color: var(--text) !important;
  background: transparent !important;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1) !important;
  display: flex !important;
  align-items: center !important;
  gap: 12px !important;
}

:deep(.logout-item:hover) {
  background: rgba(255, 80, 80, 0.08) !important;
  color: #ff5050 !important;
  transform: translateX(4px);
}

:deep(.logout-item:hover .logout-icon) {
  transform: scale(1.1) rotate(-5deg);
  color: #ff5050;
}

.logout-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  color: var(--text);
  opacity: 0.6;
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
