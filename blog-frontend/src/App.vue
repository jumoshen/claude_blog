<script setup>
import { onMounted, ref, provide, reactive } from 'vue'
import Header from './components/layout/Header.vue'
import Footer from './components/layout/Footer.vue'
import Decorations from './components/common/Decorations.vue'
import ReadingProgress from './components/common/ReadingProgress.vue'
import Toolbox from './components/common/Toolbox.vue'
import DanmuLayer from './components/common/DanmuLayer.vue'
import KeyboardNav from './components/common/KeyboardNav.vue'
import { useStyleStore } from './store/style'
import { useUserStore } from './store/user'

const styleStore = useStyleStore()
const userStore = useUserStore()

// DEBUG: Auto login as test user for development
if (import.meta.env.DEV) {
  userStore.debugAutoLogin()
}

// Danmu state - shared across app
const danmuLayer = ref(null)
const danmuSettings = reactive({
  enabled: true,
  density: 6,
  fontSize: 16,
  fontColor: '#ffffff'
})

// Load settings from localStorage
const saved = localStorage.getItem('danmu-settings')
if (saved) {
  try {
    Object.assign(danmuSettings, JSON.parse(saved))
  } catch (e) {}
}

const handleDanmuSettingsChange = (settings) => {
  Object.assign(danmuSettings, settings)
  localStorage.setItem('danmu-settings', JSON.stringify(settings))
}

// Provide danmu functions to children
provide('danmuLayer', danmuLayer)
provide('danmuSettings', danmuSettings)
provide('addDanmu', (comment) => {
  if (danmuLayer.value && danmuSettings.enabled) {
    danmuLayer.value.addDanmu(comment)
  }
})
provide('clearDanmus', () => {
  if (danmuLayer.value) {
    danmuLayer.value.clearAllDanmus()
  }
})

onMounted(() => {
  styleStore.applyTheme(styleStore.currentTheme)
})
</script>

<template>
  <div id="app">
    <Decorations />
    <ReadingProgress />
    <Toolbox @danmu-settings-change="handleDanmuSettingsChange" />
    <KeyboardNav />
    <DanmuLayer
      ref="danmuLayer"
      :visible="danmuSettings.enabled"
      :density="danmuSettings.density"
      :fontSize="danmuSettings.fontSize"
      :fontColor="danmuSettings.fontColor"
    />
    <Header />
    <main class="main">
      <router-view v-slot="{ Component }">
        <transition name="router-view" mode="out-in">
          <component :is="Component" :key="$route.fullPath" />
        </transition>
      </router-view>
    </main>
    <Footer />
  </div>
</template>

<style>
@import url('https://fonts.googleapis.com/css2?family=Noto+Sans+SC:wght@400;500;600;700&display=swap');

* { box-sizing: border-box; margin: 0; padding: 0; }
body {
  font-family: var(--font-body, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Noto Sans SC', sans-serif);
  background: var(--bg);
  color: var(--text);
  line-height: 1.7;
  transition: background 0.3s ease, color 0.3s ease;
}
h1, h2, h3, h4, h5, h6 {
  font-family: var(--font-display, var(--font-body)), sans-serif;
  color: var(--text-h);
  line-height: 1.3;
}
/* 像素风标题特殊处理 - 只对h1应用像素字体 */
.theme-pixel h1 {
  font-family: 'Press Start 2P', monospace;
  font-size: 18px;
  letter-spacing: 1px;
}
.main { min-height: calc(100vh - 60px); position: relative; z-index: 1; }

/* Page transitions */
.page-transitioning .main {
  animation: pageTransition 0.3s ease-out;
}

@keyframes pageTransition {
  0% {
    opacity: 0;
    transform: translateY(10px);
  }
  100% {
    opacity: 1;
    transform: translateY(0);
  }
}

/* Route view transitions */
.router-view-enter-active,
.router-view-leave-active {
  transition: opacity 0.25s ease, transform 0.25s ease;
}

.router-view-enter-from {
  opacity: 0;
  transform: translateY(15px);
}

.router-view-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}
</style>
