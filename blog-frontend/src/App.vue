<script setup>
import { onMounted } from 'vue'
import Header from './components/layout/Header.vue'
import Footer from './components/layout/Footer.vue'
import Decorations from './components/common/Decorations.vue'
import ReadingProgress from './components/common/ReadingProgress.vue'
import Toolbox from './components/common/Toolbox.vue'
import { useStyleStore } from './store/style'

const styleStore = useStyleStore()

onMounted(() => {
  styleStore.applyTheme(styleStore.currentTheme)
})
</script>

<template>
  <div id="app">
    <Decorations />
    <ReadingProgress />
    <Toolbox />
    <Header />
    <main class="main">
      <router-view :key="$route.fullPath" />
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
</style>
