<script setup>
import { onMounted } from 'vue'
import Header from './components/layout/Header.vue'
import Footer from './components/layout/Footer.vue'
import Decorations from './components/common/Decorations.vue'
import ReadingProgress from './components/common/ReadingProgress.vue'
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
    <Header />
    <main class="main">
      <router-view :key="$route.fullPath" />
    </main>
    <Footer />
  </div>
</template>

<style>
* { box-sizing: border-box; margin: 0; padding: 0; }
body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, sans-serif;
  background: var(--bg);
  color: var(--text);
  line-height: 1.6;
  transition: background 0.3s ease, color 0.3s ease;
}
.main { min-height: calc(100vh - 60px); position: relative; z-index: 1; }
</style>
