<template>
  <div class="danmu-layer" v-if="visible">
    <div
      v-for="item in activeDanmus"
      :key="item.id"
      class="danmu-item"
      :class="['theme-' + styleStore.currentTheme]"
      :style="getDanmuStyle(item)"
    >
      <span class="danmu-content" :style="{ color: item.textColor }">{{ item.content }}</span>
    </div>
  </div>
</template>

<script setup>
import { ref, onUnmounted } from 'vue'
import { useStyleStore } from '../../store/style'

const props = defineProps({
  visible: {
    type: Boolean,
    default: true
  },
  density: {
    type: Number,
    default: 6
  },
  fontSize: {
    type: Number,
    default: 16
  },
  fontColor: {
    type: String,
    default: '#ffffff'
  }
})

const emit = defineEmits(['danmu-end'])

const styleStore = useStyleStore()

const activeDanmus = ref([])
const danmuQueue = ref([])
const danmuIdCounter = ref(0)
let cleanupTimers = []

// Track occupied tracks
const occupiedTracks = ref(new Set())

// Danmu color palettes - bg color with corresponding text color
const colorPalettes = [
  { bg: 'rgba(100, 108, 255, 0.75)', text: '#ffffff' },   // 蓝紫
  { bg: 'rgba(255, 107, 157, 0.75)', text: '#ffffff' },  // 粉色
  { bg: 'rgba(124, 106, 255, 0.75)', text: '#ffffff' },  // 紫色
  { bg: 'rgba(0, 200, 150, 0.75)', text: '#ffffff' },    // 青绿
  { bg: 'rgba(255, 159, 64, 0.75)', text: '#ffffff' },   // 橙色
  { bg: 'rgba(255, 99, 71, 0.75)', text: '#ffffff' },    // 番茄红
  { bg: 'rgba(72, 209, 204, 0.75)', text: '#ffffff' },   // 湖蓝
  { bg: 'rgba(199, 125, 175, 0.75)', text: '#ffffff' }, // 玫瑰粉
  { bg: 'rgba(106, 168, 79, 0.75)', text: '#ffffff' },   // 草绿
  { bg: 'rgba(255, 217, 61, 0.75)', text: '#333333' },   // 黄色
]

const getRandomColor = () => {
  return colorPalettes[Math.floor(Math.random() * colorPalettes.length)]
}

const getTrackHeight = (track) => {
  return 60 + track * 50
}

const getDanmuStyle = (item) => {
  const track = item.track
  const height = getTrackHeight(track)

  return {
    top: `${height}px`,
    fontSize: `${props.fontSize}px`,
    animationDuration: '12s',
    animationTimingFunction: 'linear',
    animationDelay: '0s',
    opacity: 0.95,
    backgroundColor: item.bgColor,
    color: item.textColor
  }
}

const findAvailableTrack = () => {
  for (let i = 0; i < props.density; i++) {
    if (!occupiedTracks.value.has(i)) {
      return i
    }
  }
  return -1
}

const startDanmu = (danmu) => {
  const track = findAvailableTrack()
  if (track === -1) {
    // All tracks occupied, add to queue
    danmuQueue.value.push(danmu)
    return
  }

  occupiedTracks.value.add(track)
  danmu.track = track
  danmu.startTime = Date.now()

  activeDanmus.value.push(danmu)

  // Set timeout to remove danmu after animation
  const timer = setTimeout(() => {
    removeDanmu(danmu.id)
  }, 12000)
  cleanupTimers.push(timer)
}

const removeDanmu = (id) => {
  const index = activeDanmus.value.findIndex(d => d.id === id)
  if (index !== -1) {
    const danmu = activeDanmus.value[index]
    occupiedTracks.value.delete(danmu.track)
    activeDanmus.value.splice(index, 1)
    emit('danmu-end', danmu)

    // Process queue
    if (danmuQueue.value.length > 0) {
      const next = danmuQueue.value.shift()
      startDanmu(next)
    }
  }
}

const addDanmu = (data) => {
  const id = danmuIdCounter.value++
  const colors = getRandomColor()
  const danmu = {
    id,
    content: data.content || '',
    track: 0,
    bgColor: colors.bg,
    textColor: props.fontColor !== '#ffffff' ? props.fontColor : colors.text
  }
  startDanmu(danmu)
}

const addDanmus = (danmus) => {
  // Add historical comments with staggered delays
  danmus.forEach((danmu, index) => {
    setTimeout(() => {
      addDanmu(danmu)
    }, index * 500)
  })
}

const clearAllDanmus = () => {
  activeDanmus.value = []
  danmuQueue.value = []
  occupiedTracks.value.clear()
}

// Expose methods
defineExpose({
  addDanmu,
  addDanmus,
  clearAllDanmus
})

onUnmounted(() => {
  cleanupTimers.forEach(t => clearTimeout(t))
  cleanupTimers = []
})
</script>

<style scoped>
.danmu-layer {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100vh;
  pointer-events: none;
  z-index: 999;
  overflow: hidden;
}

.danmu-item {
  position: absolute;
  right: -300px;
  display: flex;
  align-items: center;
  padding: 8px 20px;
  border-radius: 24px;
  white-space: nowrap;
  animation: danmu-move 12s linear forwards;
  backdrop-filter: blur(4px);
}

@keyframes danmu-move {
  from {
    transform: translateX(0);
  }
  to {
    transform: translateX(calc(-100vw - 300px));
  }
}

.danmu-content {
  max-width: 500px;
  overflow: hidden;
  text-overflow: ellipsis;
  font-weight: 500;
  text-shadow: 1px 1px 2px rgba(0,0,0,0.5);
}

/* Theme variations */
.theme-pixel .danmu-item {
  border-radius: 4px;
  border: 2px solid var(--accent);
}

.theme-cute .danmu-item {
  border-radius: 16px;
}

.theme-qver .danmu-item {
  border-radius: 12px;
}
</style>
