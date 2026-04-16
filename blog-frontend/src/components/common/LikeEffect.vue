<template>
  <canvas ref="canvasRef" class="like-effect-canvas"></canvas>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'

const props = defineProps({
  x: { type: Number, required: true },
  y: { type: Number, required: true },
  color: { type: String, default: '#ff6b6b' }
})

const canvasRef = ref(null)
let ctx = null
let particles = []
let animationId = null

class Particle {
  constructor(x, y, color) {
    this.x = x
    this.y = y
    this.color = color
    this.vx = (Math.random() - 0.5) * 8
    this.vy = (Math.random() - 0.5) * 8 - 3
    this.radius = Math.random() * 4 + 2
    this.life = 1
    this.decay = Math.random() * 0.02 + 0.015
    this.gravity = 0.15
  }

  update() {
    this.vy += this.gravity
    this.x += this.vx
    this.y += this.vy
    this.life -= this.decay
    this.radius *= 0.97
  }

  draw(ctx) {
    ctx.beginPath()
    ctx.arc(this.x, this.y, this.radius, 0, Math.PI * 2)
    ctx.fillStyle = this.color
    ctx.globalAlpha = this.life
    ctx.fill()
    ctx.globalAlpha = 1
  }
}

const emit = defineEmits(['complete'])

const init = () => {
  const canvas = canvasRef.value
  if (!canvas) return

  canvas.width = window.innerWidth
  canvas.height = window.innerHeight
  ctx = canvas.getContext('2d')

  // Create particles
  const rect = canvas.getBoundingClientRect()
  const x = props.x - rect.left
  const y = props.y - rect.top

  for (let i = 0; i < 20; i++) {
    particles.push(new Particle(x, y, props.color))
  }

  animate()
}

const animate = () => {
  const canvas = canvasRef.value
  if (!canvas || !ctx) return

  ctx.clearRect(0, 0, canvas.width, canvas.height)

  particles = particles.filter(p => p.life > 0)
  particles.forEach(p => {
    p.update()
    p.draw(ctx)
  })

  if (particles.length > 0) {
    animationId = requestAnimationFrame(animate)
  } else {
    emit('complete')
  }
}

onMounted(() => {
  init()
})

onUnmounted(() => {
  if (animationId) {
    cancelAnimationFrame(animationId)
  }
})
</script>

<style scoped>
.like-effect-canvas {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
  z-index: 9999;
}
</style>
