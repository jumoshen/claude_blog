<template>
  <!-- 工具箱容器 -->
  <div class="toolbox" @mouseenter="isExpanded = true" @mouseleave="isExpanded = false">
    <!-- 翻牌弹窗 -->
    <LuckyDrawModal :visible="showCardModal" @close="showCardModal = false" />

    <!-- 展开内容 -->
    <transition name="expand">
      <div v-show="isExpanded" class="toolbox-content">
        <!-- 骰子工具 -->
        <div
          class="tool-btn dice-btn"
          :class="['theme-' + styleStore.currentTheme, { rolling: isRolling }]"
          @click="goRandom"
          title="随机文章"
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
          <span class="tool-label">{{ isRolling ? '摇啊摇...' : '随机文章' }}</span>
        </div>

        <!-- 翻牌工具 - 占位 -->
        <div
          class="tool-btn card-btn"
          :class="['theme-' + styleStore.currentTheme]"
          @click="flipCard"
          title="翻牌"
        >
          <div class="card-icon">
            <span class="card-back">?</span>
          </div>
          <span class="tool-label">翻牌</span>
        </div>
      </div>
    </transition>

    <!-- 收起时的图标按钮 -->
    <div class="toolbox-trigger" :class="{ active: isExpanded }">
      <span class="trigger-icon">🔧</span>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useStyleStore } from '../../store/style'
import api from '../../api'
import LuckyDrawModal from './LuckyDrawModal.vue'

const router = useRouter()
const styleStore = useStyleStore()
const isExpanded = ref(false)
const isRolling = ref(false)
const showCardModal = ref(false)
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
  const extra = (Math.floor(Math.random() * 2) + 2) * 360
  const finalX = x + extra
  const finalY = y + extra

  diceStyle.value = {
    transform: 'rotateX(0deg) rotateY(0deg)',
    transition: 'transform 0.01s'
  }

  setTimeout(() => {
    diceStyle.value = {
      transform: `rotateX(${finalX}deg) rotateY(${finalY}deg)`,
      transition: 'transform 1s cubic-bezier(0.2, 0.8, 0.3, 1)'
    }
  }, 50)
}

const goRandom = async () => {
  if (isRolling.value) return

  try {
    const res = await api.getPosts({ page: 1, page_size: 100 })
    if (res.code === 0 && res.data) {
      const posts = res.data.list || (Array.isArray(res.data) ? res.data : [])
      if (posts.length > 0) {
        isRolling.value = true
        diceResult.value = null

        const finalResult = Math.floor(Math.random() * 6) + 1
        diceResult.value = finalResult
        rollDice(finalResult)

        setTimeout(() => {
          const targetIndex = (finalResult - 1) % posts.length
          const randomPost = posts[targetIndex]
          isRolling.value = false
          window.scrollTo({ top: 0, behavior: 'instant' })
          if (randomPost && (randomPost.slug || randomPost.id)) {
            const targetSlug = randomPost.slug || `/post/${randomPost.id}`
            router.push(targetSlug.startsWith('/') ? targetSlug : `/post/${targetSlug}`)
          }
        }, 1300)
      }
    }
  } catch (e) {
    console.error('Failed to get random post:', e)
    isRolling.value = false
  }
}

// 翻牌功能
const flipCard = () => {
  showCardModal.value = true
}
</script>

<style scoped>
/* 工具箱容器 */
.toolbox {
  position: fixed;
  bottom: 30px;
  left: 30px;
  z-index: 1000;
  display: flex;
  align-items: flex-end;
  gap: 12px;
}

/* 展开内容 */
.toolbox-content {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 12px;
  background: var(--card-bg);
  border: 2px solid var(--accent);
  border-radius: 16px;
  box-shadow: 0 4px 20px var(--shadow);
}

/* 触发按钮 */
.toolbox-trigger {
  width: 44px;
  height: 44px;
  background: var(--card-bg);
  border: 2px solid var(--accent);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
  box-shadow: 0 4px 16px var(--shadow);
}

.toolbox-trigger:hover,
.toolbox-trigger.active {
  transform: scale(1.1);
  background: var(--accent);
}

.toolbox-trigger:hover .trigger-icon,
.toolbox-trigger.active .trigger-icon {
  color: white;
}

.trigger-icon {
  font-size: 18px;
  font-weight: bold;
  color: var(--accent);
  transition: color 0.3s;
}

/* 工具按钮 */
.tool-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  padding: 8px;
  background: transparent;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.tool-btn:hover {
  background: var(--accent-bg);
}

.tool-label {
  font-size: 10px;
  color: var(--text);
  white-space: nowrap;
}

/* 骰子按钮 */
.dice-btn.rolling {
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

/* 像素风主题 */
.theme-pixel .dice-face {
  border-radius: 4px;
  border: 2px solid var(--accent);
  background: #fff;
  box-shadow:
    inset -2px -2px 0 rgba(0,0,0,0.08),
    inset 2px 2px 0 rgba(255,255,255,0.9),
    0 0 0 1px rgba(0,0,0,0.05);
}

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

/* 可爱风主题 */
.theme-cute .dice-face {
  border-radius: 12px;
  border: none;
  background: linear-gradient(145deg, #ffffff 0%, #fff5f7 100%);
  box-shadow:
    0 4px 12px rgba(255, 107, 157, 0.2),
    inset 0 2px 4px rgba(255,255,255,0.9);
}

.theme-cute .dot {
  border-radius: 50%;
  width: 9px;
  height: 9px;
  background: linear-gradient(145deg, #ff8ab5 0%, #ff6b9d 100%);
  box-shadow:
    inset 0 2px 3px rgba(255,255,255,0.6),
    inset 0 -2px 3px rgba(0,0,0,0.1),
    0 2px 4px rgba(255, 107, 157, 0.3);
}

/* Q版主题 */
.theme-qver .dice-face {
  border-radius: 14px;
  border: 2px solid rgba(124, 106, 255, 0.3);
  background: linear-gradient(145deg, #ffffff 0%, #f0f0ff 100%);
  box-shadow:
    0 6px 20px rgba(124, 106, 255, 0.25),
    inset 0 2px 6px rgba(255,255,255,1),
    inset 0 -1px 3px rgba(124, 106, 255, 0.1);
}

.theme-qver .dot {
  border-radius: 50%;
  width: 10px;
  height: 10px;
  background: linear-gradient(145deg, #a78bfa 0%, #7c6aff 100%);
  box-shadow:
    inset 0 2px 4px rgba(255,255,255,0.5),
    inset 0 -2px 3px rgba(0,0,0,0.15),
    0 3px 6px rgba(124, 106, 255, 0.4);
}

/* 骰子六个面的定位 */
.face-1 { transform: rotateY(0deg) translateZ(22px); }
.face-2 { transform: rotateY(180deg) translateZ(22px); }
.face-3 { transform: rotateY(90deg) translateZ(22px); }
.face-4 { transform: rotateY(-90deg) translateZ(22px); }
.face-5 { transform: rotateX(90deg) translateZ(22px); }
.face-6 { transform: rotateX(-90deg) translateZ(22px); }

/* 各面的点数布局 */
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
  grid-template-columns: 1fr 1fr 1fr;
  grid-template-rows: 1fr 1fr 1fr;
  padding: 8px;
}
.face-5 .dot:nth-child(1) { grid-area: 1 / 1; justify-self: start; align-self: start; }
.face-5 .dot:nth-child(2) { grid-area: 1 / 3; justify-self: end; align-self: start; }
.face-5 .dot:nth-child(3) { grid-area: 2 / 2; justify-self: center; align-self: center; }
.face-5 .dot:nth-child(4) { grid-area: 3 / 1; justify-self: start; align-self: end; }
.face-5 .dot:nth-child(5) { grid-area: 3 / 3; justify-self: end; align-self: end; }

.face-6 { place-items: center; grid-template-columns: 1fr 1fr; grid-template-rows: 1fr 1fr 1fr; }

/* 翻牌图标 */
.card-icon {
  width: 44px;
  height: 44px;
  background: linear-gradient(145deg, var(--accent-bg) 0%, var(--card-bg) 100%);
  border: 2px solid var(--accent);
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  font-weight: bold;
  color: var(--accent);
}

.theme-pixel .card-icon {
  border-radius: 4px;
}

.theme-cute .card-icon {
  border-radius: 12px;
  border: none;
  background: linear-gradient(145deg, #ffffff 0%, #fff5f7 100%);
  box-shadow: 0 4px 12px rgba(255, 107, 157, 0.2);
}

.theme-qver .card-icon {
  border-radius: 14px;
  border: 2px solid rgba(124, 106, 255, 0.3);
  background: linear-gradient(145deg, #ffffff 0%, #f0f0ff 100%);
  box-shadow: 0 6px 20px rgba(124, 106, 255, 0.25);
}

.card-back {
  font-family: var(--font-display);
}

/* 展开动画 */
.expand-enter-active,
.expand-leave-active {
  transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
}

.expand-enter-from,
.expand-leave-to {
  opacity: 0;
  transform: translateX(-20px);
}
</style>
