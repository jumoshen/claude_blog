<template>
  <div class="archives" :class="themeClass">
    <h1>
      <span class="title-icon" v-if="styleStore.currentTheme === 'pixel'">▣</span>
      <span class="title-icon" v-else-if="styleStore.currentTheme === 'cute'">❀</span>
      <span class="title-icon" v-else>◇</span>
      Archives
    </h1>
    <div v-for="(posts, year) in archives" :key="year" class="year-section">
      <h2 class="year">
        <span class="year-num">{{ year }}</span>
      </h2>
      <ul class="post-list">
        <li v-for="post in posts" :key="post.slug" @click="$router.push(`/post/${post.slug}`)">
          <span class="date">{{ formatDate(post.date) }}</span>
          <span class="title">{{ post.title }}</span>
          <span class="arrow">→</span>
        </li>
      </ul>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import api from '../api'
import { useStyleStore } from '../store/style'

const archives = ref({})
const styleStore = useStyleStore()
const themeClass = computed(() => `theme-${styleStore.currentTheme}`)

onMounted(async () => {
  const res = await api.getArchives()
  if (res.code === 0) archives.value = res.data
})

const formatDate = (date) => {
  const d = new Date(date)
  const pad = (n) => n.toString().padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}`
}
</script>

<style scoped>
.archives {
  max-width: 800px;
  margin: 0 auto;
  padding: 40px 20px;
  position: relative;
  z-index: 1;
}

.archives h1 {
  margin-bottom: 30px;
  color: var(--text-h);
  font-size: 32px;
  display: flex;
  align-items: center;
  gap: 12px;
}

.title-icon {
  color: var(--accent);
  animation: iconSpin 4s infinite linear;
}

@keyframes iconSpin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.year-section { margin-bottom: 40px; }

.year {
  color: var(--accent);
  margin-bottom: 20px;
  font-size: 28px;
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 12px;
}

.year-num {
  position: relative;
}

.year-num::after {
  content: '';
  position: absolute;
  bottom: -4px;
  left: 0;
  width: 100%;
  height: 3px;
  background: var(--accent);
  transform: scaleX(0);
  transition: transform 0.3s ease;
}

.year-section:hover .year-num::after {
  transform: scaleX(1);
}

.post-list {
  list-style: none;
  padding: 0;
  margin: 0;
  background: var(--card-bg);
  border-radius: 12px;
  overflow: hidden;
  border: 1px solid var(--border);
}

.post-list li {
  display: flex;
  align-items: center;
  padding: 14px 20px;
  border-bottom: 1px solid var(--border);
  cursor: pointer;
  transition: all 0.2s;
  position: relative;
}

.post-list li:last-child { border-bottom: none; }

.post-list li::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  width: 3px;
  height: 100%;
  background: var(--accent);
  transform: scaleY(0);
  transition: transform 0.2s ease;
}

.post-list li:hover {
  background: var(--accent-bg);
  padding-left: 28px;
}

.post-list li:hover::before {
  transform: scaleY(1);
}

.post-list .date {
  width: 120px;
  color: var(--text);
  font-size: 13px;
  flex-shrink: 0;
  opacity: 0.7;
}

.post-list .title {
  flex: 1;
  color: var(--text-h);
  font-weight: 500;
  transition: color 0.2s;
}

.post-list li:hover .title {
  color: var(--accent);
}

.post-list .arrow {
  opacity: 0;
  transform: translateX(-10px);
  transition: all 0.2s;
  color: var(--accent);
}

.post-list li:hover .arrow {
  opacity: 1;
  transform: translateX(0);
}

/* 像素风：块状效果 */
.theme-pixel .post-list li:hover {
  background: var(--accent-bg);
}

/* 可爱风：轻微摇晃 */
.theme-cute .post-list li:hover {
  background: var(--accent-bg);
  animation: cuteShake 0.3s ease;
}

@keyframes cuteShake {
  0%, 100% { transform: translateX(0); }
  25% { transform: translateX(-3px); }
  75% { transform: translateX(3px); }
}

/* Q版：玻璃效果 */
.theme-qver .post-list {
  backdrop-filter: blur(10px);
  background: rgba(255, 255, 255, 0.7);
}
</style>
