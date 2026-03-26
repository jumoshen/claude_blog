<template>
  <header class="header">
    <div class="container">
      <div class="logo" @click="$router.push('/')">
        <img :src="styleStore.theme.logo" :alt="styleStore.theme.name" class="logo-img" />
      </div>
      <nav class="nav">
        <router-link to="/">Home</router-link>
        <router-link to="/archives">Archives</router-link>
        <router-link to="/about">About</router-link>
        <el-dropdown @command="handleStyleChange" trigger="click">
          <span class="style-switcher">
            {{ styleStore.theme.name }}
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item v-for="(t, key) in styleStore.themes" :key="key" :command="key" :class="{ active: key === styleStore.currentTheme }">
                <img :src="t.logo" class="theme-logo" :alt="t.name" />
                {{ t.name }}
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
        <template v-if="userStore.isLoggedIn">
          <el-dropdown @command="handleCommand">
            <span class="user">
              <img :src="userStore.user?.avatar_url" class="avatar" />
              {{ userStore.user?.login }}
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="logout">Logout</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </template>
        <router-link v-else to="/login">Login</router-link>
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
  background: var(--header-bg, #fff);
  box-shadow: 0 2px 8px rgba(0,0,0,0.08);
  position: sticky;
  top: 0;
  z-index: 100;
}
.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
  display: flex;
  align-items: center;
  height: 60px;
  gap: 20px;
}
.logo { cursor: pointer; display: flex; align-items: center; }
.logo-img { height: 44px; width: auto; }
.nav { margin-left: auto; display: flex; gap: 20px; align-items: center; }
.nav a { color: var(--text); text-decoration: none; transition: color 0.2s; }
.nav a:hover { color: var(--accent); }
.nav a.router-link-active { color: var(--accent); font-weight: 500; }
.user { display: flex; align-items: center; gap: 8px; cursor: pointer; }
.avatar { width: 32px; height: 32px; border-radius: 50%; }
.style-switcher {
  color: var(--accent);
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 14px;
  transition: all 0.2s;
}
.style-switcher:hover {
  background: var(--accent-bg);
}
.theme-logo {
  width: 24px;
  height: 24px;
  margin-right: 8px;
  vertical-align: middle;
}
:deep(.el-dropdown-menu__item) {
  display: flex;
  align-items: center;
}
:deep(.el-dropdown-menu__item.active) {
  color: var(--accent);
  background: var(--accent-bg);
}
</style>
