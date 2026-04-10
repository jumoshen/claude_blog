<template>
  <!-- 滚动阅读进度条 -->
  <div class="reading-progress" :style="{ width: progress + '%' }">
    <span class="progress-text">{{ Math.round(progress) }}%</span>
  </div>

  <!-- 返回顶部按钮 -->
  <transition name="fade">
    <div v-if="showBackToTop" class="back-to-top" @click="scrollToTop">
      <span class="arrow">↑</span>
    </div>
  </transition>

  </template>

<script setup>
import { ref, watch, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const progress = ref(0)
const showBackToTop = ref(false)

const getStorageKey = () => `reading_progress_${route.params.slug}`

const saveProgress = () => {
  if (route.params.slug) {
    localStorage.setItem(getStorageKey(), progress.value.toString())
  }
}

const restoreProgress = () => {
  if (route.params.slug) {
    const saved = localStorage.getItem(getStorageKey())
    if (saved) {
      const savedProgress = parseFloat(saved)
      if (savedProgress > 0 && savedProgress < 100) {
        // 恢复阅读位置
        setTimeout(() => {
          const docHeight = document.documentElement.scrollHeight - window.innerHeight
          window.scrollTo({ top: (savedProgress / 100) * docHeight, behavior: 'instant' })
          progress.value = savedProgress
        }, 100)
      }
    }
  }
}

const updateProgress = () => {
  const scrollTop = window.scrollY
  const docHeight = document.documentElement.scrollHeight - window.innerHeight
  progress.value = docHeight > 0 ? (scrollTop / docHeight) * 100 : 0
  showBackToTop.value = scrollTop > 500

  // 保存进度
  saveProgress()
}

const scrollToTop = () => {
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

// 监听路由变化，重置进度条
watch(() => route.path, () => {
  progress.value = 0
  showBackToTop.value = false
  restoreProgress()
})

onMounted(() => {
  window.addEventListener('scroll', updateProgress)
  restoreProgress()
})

onUnmounted(() => {
  window.removeEventListener('scroll', updateProgress)
})
</script>

<style scoped>
/* 阅读进度条 */
.reading-progress {
  position: fixed;
  top: 0;
  left: 0;
  height: 4px;
  background: linear-gradient(90deg, var(--accent), var(--accent-border));
  z-index: 10000;
  transition: width 0.1s ease-out;
  display: flex;
  align-items: center;
  justify-content: flex-end;
}

.progress-text {
  position: absolute;
  right: 10px;
  top: 8px;
  font-size: 10px;
  color: var(--accent);
  font-weight: 600;
  opacity: 0;
  transition: opacity 0.3s;
}

.reading-progress:hover .progress-text {
  opacity: 1;
}

/* 返回顶部按钮 */
.back-to-top {
  position: fixed;
  bottom: 30px;
  right: 30px;
  width: 50px;
  height: 50px;
  background: var(--accent);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  z-index: 1000;
  box-shadow: 0 4px 20px var(--shadow);
  transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
}

.back-to-top:hover {
  transform: translateY(-5px) scale(1.1);
  box-shadow: 0 8px 30px var(--shadow);
}

.back-to-top .arrow {
  font-size: 24px;
  color: white;
  animation: bounce 1s infinite;
}

@keyframes bounce {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-3px); }
}

.fade-enter-active, .fade-leave-active {
  transition: opacity 0.3s;
}
.fade-enter-from, .fade-leave-to {
  opacity: 0;
}
</style>
