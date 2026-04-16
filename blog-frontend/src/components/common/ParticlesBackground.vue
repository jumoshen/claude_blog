<template>
  <canvas ref="canvasRef" class="particles-bg"></canvas>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'

const canvasRef = ref(null)
let ctx = null
let particles = []
let animationId = null
let mouseX = 0
let mouseY = 0

const props = defineProps({
  particleCount: {
    type: Number,
    default: 50
  },
  particleColor: {
    type: String,
    default: 'var(--accent)'
  },
  connectDistance: {
    type: Number,
    default: 120
  },
  moveSpeed: {
    type: Number,
    default: 0.5
  }
})

class Particle {
  constructor(canvas) {
    this.x = Math.random() * canvas.width
    this.y = Math.random() * canvas.height
    this.vx = (Math.random() - 0.5) * props.moveSpeed
    this.vy = (Math.random() - 0.5) * props.moveSpeed
    this.radius = Math.random() * 2 + 1
    this.opacity = Math.random() * 0.5 + 0.2
  }

  update(canvas) {
    // Mouse repulsion
    const dx = this.x - mouseX
    const dy = this.y - mouseY
    const dist = Math.sqrt(dx * dx + dy * dy)
    if (dist < 100) {
      const force = (100 - dist) / 100
      this.vx += dx * force * 0.01
      this.vy += dy * force * 0.01
    }

    this.x += this.vx
    this.y += this.vy

    // Boundary check
    if (this.x < 0 || this.x > canvas.width) this.vx *= -1
    if (this.y < 0 || this.y > canvas.height) this.vy *= -1

    // Damping
    this.vx *= 0.99
    this.vy *= 0.99
  }

  draw(ctx) {
    ctx.beginPath()
    ctx.arc(this.x, this.y, this.radius, 0, Math.PI * 2)
    ctx.fillStyle = `rgba(179, 102, 255, ${this.opacity})`
    ctx.fill()
  }
}

const init = () => {
  const canvas = canvasRef.value
  if (!canvas) return

  ctx = canvas.getContext('2d')
  resizeCanvas()

  // Create particles
  particles = []
  for (let i = 0; i < props.particleCount; i++) {
    particles.push(new Particle(canvas))
  }

  animate()
}

const resizeCanvas = () => {
  const canvas = canvasRef.value
  if (!canvas) return
  canvas.width = window.innerWidth
  canvas.height = window.innerHeight
}

const animate = () => {
  const canvas = canvasRef.value
  if (!canvas || !ctx) return

  ctx.clearRect(0, 0, canvas.width, canvas.height)

  // Update and draw particles
  particles.forEach(p => {
    p.update(canvas)
    p.draw(ctx)
  })

  // Draw connections
  for (let i = 0; i < particles.length; i++) {
    for (let j = i + 1; j < particles.length; j++) {
      const dx = particles[i].x - particles[j].x
      const dy = particles[i].y - particles[j].y
      const dist = Math.sqrt(dx * dx + dy * dy)

      if (dist < props.connectDistance) {
        ctx.beginPath()
        ctx.moveTo(particles[i].x, particles[i].y)
        ctx.lineTo(particles[j].x, particles[j].y)
        ctx.strokeStyle = `rgba(179, 102, 255, ${0.15 * (1 - dist / props.connectDistance)})`
        ctx.lineWidth = 0.5
        ctx.stroke()
      }
    }
  }

  animationId = requestAnimationFrame(animate)
}

const handleMouseMove = (e) => {
  mouseX = e.clientX
  mouseY = e.clientY
}

onMounted(() => {
  init()
  window.addEventListener('resize', resizeCanvas)
  window.addEventListener('mousemove', handleMouseMove)
})

onUnmounted(() => {
  if (animationId) {
    cancelAnimationFrame(animationId)
  }
  window.removeEventListener('resize', resizeCanvas)
  window.removeEventListener('mousemove', handleMouseMove)
})
</script>

<style scoped>
.particles-bg {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
  z-index: 0;
  opacity: 0.6;
}
</style>
