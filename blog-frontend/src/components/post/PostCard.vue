<template>
  <div class="post-card" :class="themeClass" @click="$router.push(`/post/${post.slug}`)">
    <div class="card-accent-line"></div>

    <div class="card-header">
      <h2 class="title">{{ post.title }}</h2>
      <div class="reading-time" v-if="readingTime">
        <span class="clock">⏱️</span> {{ readingTime }}
      </div>
    </div>

    <div class="meta">
      <span class="date">📅 {{ formatDate(post.date) }}</span>
      <span class="views">👁️ {{ post.views }}</span>
    </div>

    <div class="tags">
      <el-tag
        v-for="tag in post.tags"
        :key="tag"
        size="small"
        class="tag-item"
        @click.stop="$router.push({ path: '/', query: { tag } })"
      >{{ tag }}</el-tag>
    </div>

    <p class="summary">{{ post.summary }}</p>

    <div class="card-footer">
      <span class="read-more">阅读全文</span>
      <span class="arrow">→</span>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useStyleStore } from '../../store/style'

const props = defineProps({
  post: { type: Object, required: true },
})

const styleStore = useStyleStore()
const themeClass = computed(() => `theme-${styleStore.currentTheme}`)

// 估算阅读时间（按每分钟400字计算）
const readingTime = computed(() => {
  if (!props.post.content) return null
  const words = props.post.content.length
  const minutes = Math.ceil(words / 400)
  return minutes < 1 ? '< 1 分钟' : `${minutes} 分钟`
})

const formatDate = (date) => {
  const d = new Date(date)
  const pad = (n) => n.toString().padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}`
}
</script>

<style scoped>
.post-card {
  background: var(--card-bg);
  border-radius: 16px;
  padding: 24px;
  margin-bottom: 20px;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
  border: 1px solid var(--border);
  position: relative;
  overflow: hidden;
}

.card-accent-line {
  position: absolute;
  top: 0;
  left: 0;
  width: 0;
  height: 4px;
  background: linear-gradient(90deg, var(--accent), var(--accent-border));
  transition: width 0.4s ease;
}

.post-card:hover .card-accent-line {
  width: 100%;
}

.post-card:hover {
  border-color: var(--accent);
  box-shadow: 0 12px 40px var(--shadow);
}

/* 通用悬浮效果 */
.post-card:hover {
  transform: translateY(-4px);
}

/* 像素风：块状浮起 + 像素阴影 */
.theme-pixel.post-card:hover {
  transform: translateY(-4px);
  box-shadow: 8px 8px 0 var(--accent-border);
}

/* 可爱风：弹性放大 */
.theme-cute.post-card:hover {
  transform: translateY(-8px) scale(1.02);
}

/* Q版：玻璃效果 */
.theme-qver.post-card {
  backdrop-filter: blur(10px);
}
.theme-qver.post-card:hover {
  transform: translateY(-6px);
  box-shadow: 0 20px 50px var(--shadow);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 12px;
}

.title {
  margin: 0;
  color: var(--text-h);
  font-size: 20px;
  font-weight: 600;
  flex: 1;
  line-height: 1.4;
  transition: color 0.2s;
}

.post-card:hover .title {
  color: var(--accent);
}

.reading-time {
  font-size: 12px;
  color: var(--text);
  opacity: 0.7;
  white-space: nowrap;
  margin-left: 12px;
  display: flex;
  align-items: center;
  gap: 4px;
}

.meta {
  font-size: 13px;
  color: var(--text);
  margin-bottom: 12px;
  opacity: 0.8;
  display: flex;
  gap: 16px;
}

.tags {
  margin-bottom: 12px;
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.tag-item {
  cursor: pointer;
  transition: all 0.2s cubic-bezier(0.34, 1.56, 0.64, 1);
  border: none;
  background: var(--accent-bg);
  color: var(--accent);
  border-radius: 20px;
  padding: 0 12px;
}

.tag-item:hover {
  transform: scale(1.1) rotate(-2deg);
  background: var(--accent);
  color: white;
}

.summary {
  color: var(--text);
  line-height: 1.8;
  margin: 0 0 16px;
  opacity: 0.9;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 12px;
  border-top: 1px solid var(--border);
}

.read-more {
  font-size: 13px;
  color: var(--accent);
  font-weight: 500;
}

.arrow {
  font-size: 18px;
  color: var(--accent);
  opacity: 0;
  transform: translateX(-10px);
  transition: all 0.3s ease;
}

.post-card:hover .arrow {
  opacity: 1;
  transform: translateX(0);
}
</style>
