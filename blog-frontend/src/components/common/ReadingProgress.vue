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

  <!-- 随机文章入口 -->
  <div
    class="random-btn"
    :class="['theme-' + styleStore.currentTheme, { rolling: isRolling }]"
    @click="goRandom"
    title="随机一篇文章"
  >
    <span class="dice" :style="{ transform: `rotate(${diceRotation}deg)` }">🎲</span>
    <span class="dice-face" v-if="diceResult">{{ diceResult }}</span>
    <span class="tooltip">{{ isRolling ? '摇啊摇...' : '随机文章' }}</span>
  </div>
</template>

<script setup>
import { ref, watch, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useStyleStore } from '../../store/style'
import api from '../../api'

const router = useRouter()
const route = useRoute()
const styleStore = useStyleStore()
const progress = ref(0)
const showBackToTop = ref(false)
const isRolling = ref(false)
const diceResult = ref(null)
const diceRotation = ref(0)

const updateProgress = () => {
  const scrollTop = window.scrollY
  const docHeight = document.documentElement.scrollHeight - window.innerHeight
  progress.value = docHeight > 0 ? (scrollTop / docHeight) * 100 : 0
  showBackToTop.value = scrollTop > 500
}

const scrollToTop = () => {
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

const goRandom = async () => {
  if (isRolling.value) return

  try {
    const res = await api.getPosts({ page: 1, page_size: 100 })
    if (res.code === 0 && res.data && res.data.length > 0) {
      const posts = res.data

      // 开始滚动动画
      isRolling.value = true
      diceResult.value = null

      // 滚动动画：快速变换角度
      const rollDuration = 1500 // 1.5秒
      const rollInterval = 50 // 每50ms更新一次
      let elapsed = 0

      const rollAnim = setInterval(() => {
        elapsed += rollInterval
        diceRotation.value += 180 + Math.random() * 180

        if (elapsed >= rollDuration) {
          clearInterval(rollAnim)

          // 动画结束，确定最终点数 (1-6)
          const finalResult = Math.floor(Math.random() * 6) + 1
          diceResult.value = finalResult

          // 用骰子点数作为索引偏移
          const baseIndex = Math.floor(Math.random() * posts.length)
          const offsetIndex = (baseIndex + finalResult - 1) % posts.length
          const randomPost = posts[offsetIndex]

          // 延迟一下显示结果
          setTimeout(() => {
            isRolling.value = false
            window.scrollTo({ top: 0, behavior: 'instant' })
            router.push(`/post/${randomPost.slug}`)
          }, 500)
        }
      }, rollInterval)
    }
  } catch (e) {
    console.error('Failed to get random post:', e)
    isRolling.value = false
  }
}

// 监听路由变化，重置进度条
watch(() => route.path, () => {
  progress.value = 0
  showBackToTop.value = false
})

onMounted(() => {
  window.addEventListener('scroll', updateProgress)
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

/* 随机文章按钮 */
.random-btn {
  position: fixed;
  bottom: 30px;
  left: 30px;
  width: 50px;
  height: 50px;
  background: var(--card-bg);
  border: 2px solid var(--accent);
  border-radius: 50%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  z-index: 1000;
  transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
  overflow: visible;
}

.random-btn:hover {
  transform: scale(1.1);
}

.random-btn.rolling {
  animation: pulse 0.3s ease-in-out infinite;
  pointer-events: none;
}

@keyframes pulse {
  0%, 100% { transform: scale(1); }
  50% { transform: scale(1.15); }
}

.random-btn .dice {
  font-size: 24px;
  transition: transform 0.05s linear;
}

.random-btn.rolling .dice {
  font-size: 18px;
}

.random-btn .dice-face {
  position: absolute;
  bottom: -2px;
  right: -2px;
  background: var(--accent);
  color: white;
  font-size: 10px;
  font-weight: bold;
  width: 18px;
  height: 18px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 2px solid var(--card-bg);
}

.random-btn .tooltip {
  position: absolute;
  left: 60px;
  background: var(--accent);
  color: white;
  padding: 6px 12px;
  border-radius: 8px;
  font-size: 12px;
  white-space: nowrap;
  opacity: 0;
  pointer-events: none;
  transition: opacity 0.3s;
}

.random-btn:hover .tooltip {
  opacity: 1;
}

.fade-enter-active, .fade-leave-active {
  transition: opacity 0.3s;
}
.fade-enter-from, .fade-leave-to {
  opacity: 0;
}
</style>
