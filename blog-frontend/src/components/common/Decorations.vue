<template>
  <div class="bg-decorations" v-if="styleStore.currentTheme !== 'qver'">
    <!-- 像素风：星星和方块 -->
    <template v-if="styleStore.currentTheme === 'pixel'">
      <div class="pixel-star" v-for="i in 8" :key="'star-'+i" :style="getStarStyle(i)">
        <svg viewBox="0 0 20 20" width="20" height="20">
          <rect x="6" y="6" width="8" height="8" :fill="styleStore.theme.colors.accent"/>
        </svg>
      </div>
      <div class="pixel-block" v-for="i in 5" :key="'block-'+i" :style="getBlockStyle(i)"></div>
    </template>

    <!-- 可爱风：爱心和泡泡 -->
    <template v-if="styleStore.currentTheme === 'cute'">
      <div class="cute-heart" v-for="i in 6" :key="'heart-'+i" :style="getHeartStyle(i)">♥</div>
      <div class="cute-bubble" v-for="i in 8" :key="'bubble-'+i" :style="getBubbleStyle(i)"></div>
    </template>
  </div>

  <!-- Q版：渐变光晕 -->
  <div class="bg-decorations qver-glow" v-else>
    <div class="glow-orb glow-orb-1"></div>
    <div class="glow-orb glow-orb-2"></div>
    <div class="glow-orb glow-orb-3"></div>
  </div>

  <!-- 滚动进度条 -->
  <div class="scroll-progress" :style="{ width: scrollProgress + '%' }"></div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useStyleStore } from '../../store/style'

const styleStore = useStyleStore()
const scrollProgress = ref(0)

const getStarStyle = (i) => ({
  left: `${(i * 13) % 100}%`,
  top: `${(i * 17) % 100}%`,
  animationDelay: `${i * 0.5}s`,
  animationDuration: `${2 + (i % 3)}s`,
})

const getBlockStyle = (i) => ({
  left: `${(i * 19) % 100}%`,
  top: `${(i * 23) % 100}%`,
  animationDelay: `${i * 0.3}s`,
  width: `${12 + (i % 3) * 8}px`,
  height: `${12 + (i % 3) * 8}px`,
})

const getHeartStyle = (i) => ({
  left: `${(i * 15 + 5) % 100}%`,
  top: `${(i * 19) % 100}%`,
  animationDelay: `${i * 0.4}s`,
  fontSize: `${12 + (i % 3) * 6}px`,
  color: i % 2 === 0 ? styleStore.theme.colors.accent : '#ff9a9a',
})

const getBubbleStyle = (i) => ({
  left: `${(i * 12 + 8) % 100}%`,
  top: `${(i * 21) % 100}%`,
  animationDelay: `${i * 0.3}s`,
  width: `${8 + (i % 4) * 4}px`,
  height: `${8 + (i % 4) * 4}px`,
})

const handleScroll = () => {
  const scrollTop = window.scrollY
  const docHeight = document.documentElement.scrollHeight - window.innerHeight
  scrollProgress.value = docHeight > 0 ? (scrollTop / docHeight) * 100 : 0
}

onMounted(() => {
  window.addEventListener('scroll', handleScroll)
})

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll)
})
</script>

<style scoped>
.bg-decorations {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
  overflow: hidden;
  z-index: 0;
}

/* 滚动进度条 */
.scroll-progress {
  position: fixed;
  top: 0;
  left: 0;
  height: 3px;
  background: linear-gradient(90deg, var(--accent), var(--accent-border));
  z-index: 9999;
  transition: width 0.1s ease-out;
}

/* ========== 像素风样式 ========== */
.pixel-star {
  position: absolute;
  animation: pixelBlink 1.5s infinite ease-in-out;
}

.pixel-block {
  position: absolute;
  background: var(--accent);
  opacity: 0.1;
  animation: pixelMove 4s infinite ease-in-out;
}

@keyframes pixelBlink {
  0%, 100% { opacity: 1; transform: scale(1); }
  50% { opacity: 0.3; transform: scale(0.8); }
}

@keyframes pixelMove {
  0%, 100% { transform: translate(0, 0); }
  25% { transform: translate(10px, -15px); }
  50% { transform: translate(-5px, 10px); }
  75% { transform: translate(8px, 5px); }
}

/* ========== 可爱风样式 ========== */
.cute-heart {
  position: absolute;
  animation: cuteFloat 3s infinite ease-in-out;
}

.cute-bubble {
  position: absolute;
  border-radius: 50%;
  background: linear-gradient(135deg, var(--accent-bg), transparent);
  border: 1px solid var(--accent-border);
  animation: cuteBubble 4s infinite ease-in-out;
}

@keyframes cuteFloat {
  0%, 100% { transform: translateY(0) rotate(0deg); }
  25% { transform: translateY(-15px) rotate(5deg); }
  50% { transform: translateY(-5px) rotate(-3deg); }
  75% { transform: translateY(-20px) rotate(3deg); }
}

@keyframes cuteBubble {
  0%, 100% { transform: translateY(0) scale(1); opacity: 0.6; }
  50% { transform: translateY(-30px) scale(1.1); opacity: 0.3; }
}

/* ========== Q版样式 ========== */
.qver-glow {
  background: transparent;
}

.glow-orb {
  position: absolute;
  border-radius: 50%;
  filter: blur(60px);
  opacity: 0.3;
  animation: qverGlowMove 8s infinite ease-in-out;
}

.glow-orb-1 {
  width: 400px;
  height: 400px;
  background: var(--accent);
  top: -100px;
  left: -100px;
  animation-delay: 0s;
}

.glow-orb-2 {
  width: 300px;
  height: 300px;
  background: #a78bfa;
  bottom: -50px;
  right: -50px;
  animation-delay: -3s;
}

.glow-orb-3 {
  width: 200px;
  height: 200px;
  background: #f472b6;
  top: 50%;
  left: 50%;
  animation-delay: -5s;
}

@keyframes qverGlowMove {
  0%, 100% { transform: translate(0, 0) scale(1); }
  33% { transform: translate(30px, -30px) scale(1.1); }
  66% { transform: translate(-20px, 20px) scale(0.95); }
}
</style>
