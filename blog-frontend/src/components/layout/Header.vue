<template>
  <header class="header">
    <div class="container">
      <div class="logo" @click="$router.push('/')">{{ siteTitle }}</div>
      <nav class="nav">
        <router-link to="/">Home</router-link>
        <router-link to="/archives">Archives</router-link>
        <router-link to="/about">About</router-link>
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
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../../store/user'
import api from '../../api'

const router = useRouter()
const userStore = useUserStore()
const siteTitle = ref('Blog')

onMounted(async () => {
  try {
    const res = await api.getSite()
    if (res.code === 0) {
      siteTitle.value = res.data.title
    }
  } catch (e) {
    // use default
  }
})

const handleCommand = (command) => {
  if (command === 'logout') {
    userStore.logout()
    router.push('/')
  }
}
</script>

<style scoped>
.header {
  background: #fff;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
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
}
.logo { font-size: 24px; font-weight: bold; cursor: pointer; color: #409eff; }
.nav { margin-left: auto; display: flex; gap: 20px; align-items: center; }
.nav a { color: #333; text-decoration: none; }
.nav a:hover { color: #409eff; }
.nav a.router-link-active { color: #409eff; }
.user { display: flex; align-items: center; gap: 8px; cursor: pointer; }
.avatar { width: 32px; height: 32px; border-radius: 50%; }
</style>
