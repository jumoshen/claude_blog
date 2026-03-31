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
    <div class="dice-container" :class="{ 'is-rolling': isRolling }">
      <div class="dice-shadow"></div>
      <div class="dice" :style="diceStyle">
        <div class="dice-face face-1"><span class="dot"></span></div>
        <div class="dice-face face-2"><span class="dot"></span><span class="dot"></span></div>
        <div class="dice-face face-3"><span class="dot"></span><span class="dot"></span><span class="dot"></span></div>
        <div class="dice-face face-4"><span class="dot"></span><span class="dot"></span><span class="dot"></span><span class="dot"></span></div>
        <div class="dice-face face-5"><span class="dot"></span><span class="dot"></span><span class="dot"></span><span class="dot"></span><span class="dot"></span></div>
        <div class="dice-face face-6"><span class="dot"></span><span class="dot"></span><span class="dot"></span><span class="dot"></span><span class="dot"></span><span class="dot"></span></div>
      </div>
    </div>
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
const diceStyle = ref({})

// 面角度映射：每个点数对应的旋转角度
const FACE_ANGLES = {
  1: { x: 0, y: 0 },
  2: { x: 0, y: 180 },
  3: { x: 0, y: -90 },
  4: { x: 0, y: 90 },
  5: { x: -90, y: 0 },
  6: { x: 90, y: 0 }
}

// 掷骰子动画
const rollDice = (targetFace) => {
  const { x, y } = FACE_ANGLES[targetFace]
  // 额外整圈：720~1080度
  const extra = (Math.floor(Math.random() * 2) + 2) * 360
  const finalX = x + extra
  const finalY = y + extra

  // 先重置到初始位置（无过渡）
  diceStyle.value = {
    transform: 'rotateX(0deg) rotateY(0deg)',
    transition: 'transform 0.01s'
  }

  // 下一步设置最终位置触发动画
  setTimeout(() => {
    diceStyle.value = {
      transform: `rotateX(${finalX}deg) rotateY(${finalY}deg)`,
      transition: 'transform 1s cubic-bezier(0.2, 0.8, 0.3, 1)'
    }
  }, 50)
}

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
    if (res.code === 0 && res.data) {
      // 兼容新旧两种 API 格式
      const posts = res.data.list || (Array.isArray(res.data) ? res.data : [])
      if (posts.length > 0) {
        isRolling.value = true
        diceResult.value = null

        // 随机结果 1-6
        const finalResult = Math.floor(Math.random() * 6) + 1
        diceResult.value = finalResult

        // 启动掷骰子动画
        rollDice(finalResult)

        // 动画结束后跳转
        setTimeout(() => {
          const targetIndex = (finalResult - 1) % posts.length
          const randomPost = posts[targetIndex]
          isRolling.value = false
          window.scrollTo({ top: 0, behavior: 'instant' })
          if (randomPost && randomPost.slug) {
            router.push(`/post/${randomPost.slug}`)
          }
        }, 1300)
      }
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
  width: 56px;
  height: 56px;
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
  box-shadow: 0 4px 16px var(--shadow);
}

/* 像素风按钮更像素 */
.theme-pixel.random-btn {
  border-radius: 12px;
  border-width: 3px;
}

.theme-pixel.random-btn:hover {
  transform: scale(1.08);
}

.random-btn:hover {
  transform: scale(1.1);
  box-shadow: 0 6px 24px var(--shadow);
}

.random-btn.rolling {
  pointer-events: none;
}

/* 3D 骰子容器 */
.dice-container {
  width: 44px;
  height: 44px;
  perspective: 900px;
  perspective-origin: center center;
  position: relative;
}

/* 骰子阴影 */
.dice-shadow {
  position: absolute;
  width: 40px;
  height: 10px;
  background: radial-gradient(ellipse, rgba(0,0,0,0.25) 0%, transparent 70%);
  border-radius: 50%;
  bottom: -8px;
  left: 50%;
  transform: translateX(-50%);
  transition: all 0.3s ease;
}

.dice-container.is-rolling .dice-shadow {
  animation: shadowPulse 0.15s ease-in-out infinite alternate;
}

@keyframes shadowPulse {
  from { transform: translateX(-50%) scale(0.9); opacity: 0.8; }
  to { transform: translateX(-50%) scale(1.1); opacity: 0.4; }
}

/* 骰子本体 */
.dice {
  width: 44px;
  height: 44px;
  position: relative;
  transform-style: preserve-3d;
  will-change: transform;
  transform: rotateX(-25deg) rotateY(-40deg);
}

.dice-container.is-rolling .dice {
  transition: none;
}

/* 骰子每个面 - 清晰的白色+紫色边框 */
.dice-face {
  position: absolute;
  width: 44px;
  height: 44px;
  background: linear-gradient(145deg, #ffffff 0%, #f5f5f5 100%);
  border: 2px solid var(--accent);
  border-radius: 10px;
  box-sizing: border-box;
  backface-visibility: hidden;
}

/* 骰子点数 - 使用 grid 精确定位 */
.dice-face {
  display: grid;
  padding: 6px;
}

.dot {
  width: 8px;
  height: 8px;
  background: linear-gradient(145deg, var(--accent) 0%, color-mix(in srgb, var(--accent) 85%, black) 100%);
  border-radius: 50%;
  box-shadow: inset 0 1px 2px rgba(255,255,255,0.4), inset 0 -1px 1px rgba(0,0,0,0.15);
  justify-self: center;
  align-self: center;
}

/* 像素风主题 - 骰子像素化但保持圆润 */
.theme-pixel .dice-face {
  border-radius: 4px;
  border: 2px solid var(--accent);
  background: #fff;
  box-shadow:
    inset -2px -2px 0 rgba(0,0,0,0.08),
    inset 2px 2px 0 rgba(255,255,255,0.9),
    0 0 0 1px rgba(0,0,0,0.05);
}

/* 像素风点数 - 小方块带像素风阴影 */
.theme-pixel .dot {
  border-radius: 0;
  width: 7px;
  height: 7px;
  background: var(--accent);
  box-shadow:
    inset -1px -1px 0 rgba(0,0,0,0.25),
    inset 1px 1px 0 rgba(255,255,255,0.35),
    1px 1px 0 rgba(0,0,0,0.15);
}

/* 骰子六个面的定位 */
.face-1 { transform: rotateY(0deg) translateZ(22px); }
.face-2 { transform: rotateY(180deg) translateZ(22px); }
.face-3 { transform: rotateY(90deg) translateZ(22px); }
.face-4 { transform: rotateY(-90deg) translateZ(22px); }
.face-5 { transform: rotateX(90deg) translateZ(22px); }
.face-6 { transform: rotateX(-90deg) translateZ(22px); }

/* 各面的点数布局 - 使用 grid 精确定位 */
.face-1 { place-items: center; }

.face-2 { place-items: center; grid-template-columns: 1fr 1fr; grid-template-rows: 1fr 1fr; }
.face-2 .dot:first-child { grid-area: 1 / 1; justify-self: start; align-self: start; }
.face-2 .dot:last-child { grid-area: 2 / 2; justify-self: end; align-self: end; }

.face-3 { place-items: center; grid-template-columns: 1fr 1fr 1fr; grid-template-rows: 1fr 1fr 1fr; }
.face-3 .dot:nth-child(1) { grid-area: 1 / 1; justify-self: start; align-self: start; }
.face-3 .dot:nth-child(2) { grid-area: 2 / 2; }
.face-3 .dot:nth-child(3) { grid-area: 3 / 3; justify-self: end; align-self: end; }

.face-4 { place-items: center; grid-template-columns: 1fr 1fr; grid-template-rows: 1fr 1fr; }

.face-5 {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  grid-template-rows: repeat(3, 1fr);
  place-items: center;
  padding: 8px;
}
.face-5 .dot:nth-child(1) { grid-area: 1 / 1; }
.face-5 .dot:nth-child(2) { grid-area: 1 / 3; }
.face-5 .dot:nth-child(3) { grid-area: 2 / 2; }
.face-5 .dot:nth-child(4) { grid-area: 3 / 1; }
.face-5 .dot:nth-child(5) { grid-area: 3 / 3; }

.face-6 { place-items: center; grid-template-columns: 1fr 1fr; grid-template-rows: 1fr 1fr 1fr; }

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
