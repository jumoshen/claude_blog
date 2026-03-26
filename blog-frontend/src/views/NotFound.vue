<template>
  <div class="not-found" :class="'theme-' + styleStore.currentTheme">
    <div class="troll-face">{{ trollFace }}</div>
    <h1 class="code">404</h1>
    <h2 class="title">页面走丢了...</h2>
    <p class="desc">{{ description }}</p>

    <div class="actions">
      <button class="btn-primary" @click="$router.push('/')">
        <span>🏠</span> 返回首页
      </button>
      <button class="btn-secondary" @click="$router.back()">
        <span>←</span> 返回上页
      </button>
    </div>

    <div class="floating-elements">
      <span v-for="i in 5" :key="i" class="floating-item" :style="getStyle(i)">
        {{ getEmoji(i) }}
      </span>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useStyleStore } from '../store/style'

const styleStore = useStyleStore()

const trollFaces = ['ಠ_ಠ', '(╯°□°)╯︵ ┻━┻', '┬─┬ノ( º _ ºノ)', '(ノಠ益ಠ)ノ', 'ಠ╭╮ಠ']
const descriptions = [
  '这只巨魔正在寻找你想要的页面...',
  '糟糕！页面被外星人绑架了！',
  '这个页面好像逃跑了...',
  '呜~页面不见了...',
  '糟糕，迷路了！'
]

const trollFace = computed(() => trollFaces[Math.floor(Math.random() * trollFaces.length)])
const description = computed(() => descriptions[Math.floor(Math.random() * descriptions.length)])

const getStyle = (i) => ({
  left: `${10 + i * 18}%`,
  animationDelay: `${i * 0.3}s`,
  fontSize: `${16 + (i % 3) * 8}px`
})

const getEmoji = (i) => {
  const emojis = ['✨', '⭐', '🌟', '💫', '🎈']
  return emojis[i % emojis.length]
}
</script>

<style scoped>
.not-found {
  min-height: calc(100vh - 120px);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
  padding: 40px 20px;
  position: relative;
  overflow: hidden;
}

.troll-face {
  font-size: 64px;
  margin-bottom: 20px;
  animation: bounce 1s infinite;
}

@keyframes bounce {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-20px); }
}

.code {
  font-size: 120px;
  font-weight: 800;
  color: var(--accent);
  opacity: 0.3;
  margin: 0;
  line-height: 1;
}

.title {
  font-size: 28px;
  color: var(--text-h);
  margin: 16px 0;
}

.desc {
  font-size: 16px;
  color: var(--text);
  margin-bottom: 32px;
  max-width: 400px;
}

.actions {
  display: flex;
  gap: 16px;
  flex-wrap: wrap;
  justify-content: center;
}

.btn-primary, .btn-secondary {
  padding: 12px 28px;
  border-radius: 50px;
  border: none;
  font-size: 15px;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 8px;
  transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
}

.btn-primary {
  background: var(--accent);
  color: white;
}

.btn-primary:hover {
  transform: scale(1.05) translateY(-2px);
  box-shadow: 0 8px 25px var(--shadow);
}

.btn-secondary {
  background: var(--card-bg);
  color: var(--text);
  border: 2px solid var(--border);
}

.btn-secondary:hover {
  border-color: var(--accent);
  color: var(--accent);
  transform: translateY(-2px);
}

.floating-elements {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
  z-index: -1;
}

.floating-item {
  position: absolute;
  animation: float 3s ease-in-out infinite;
  opacity: 0.5;
}

@keyframes float {
  0%, 100% { transform: translateY(0) rotate(0deg); }
  50% { transform: translateY(-20px) rotate(10deg); }
}

/* 像素风特殊效果 */
.theme-pixel .code {
  font-family: monospace;
  letter-spacing: -5px;
}

.theme-pixel .troll-face {
  image-rendering: pixelated;
}

/* 可爱风特殊效果 */
.theme-cute .troll-face {
  animation: wiggle 0.5s ease-in-out infinite;
}

@keyframes wiggle {
  0%, 100% { transform: rotate(-5deg); }
  50% { transform: rotate(5deg); }
}

/* Q版特殊效果 */
.theme-qver .code {
  background: linear-gradient(135deg, var(--accent), #a78bfa);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}
</style>
