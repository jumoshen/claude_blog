<template>
  <Teleport to="body">
    <Transition name="modal-fade">
      <div v-if="visible" class="share-modal" @click.self="close">
        <div class="share-container">
          <div class="share-header">
            <h3>分享文章</h3>
            <button class="close-btn" @click="close">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M18 6L6 18M6 6l12 12"/>
              </svg>
            </button>
          </div>

          <!-- Poster Preview -->
          <div class="poster-preview">
            <canvas ref="posterCanvas" class="poster-canvas"></canvas>
          </div>

          <!-- Share Actions -->
          <div class="share-actions">
            <button class="action-btn download" @click="downloadPoster">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4M7 10l5 5 5-5M12 15V3"/>
              </svg>
              <span>下载海报</span>
            </button>
            <button class="action-btn copy-link" @click="copyLink">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="9" y="9" width="13" height="13" rx="2" ry="2"/>
                <path d="M5 15H4a2 2 0 01-2-2V4a2 2 0 012-2h9a2 2 0 012 2v1"/>
              </svg>
              <span>{{ copied ? '已复制' : '复制链接' }}</span>
            </button>
          </div>

          <!-- Social Share -->
          <div class="social-share">
            <span class="share-label">分享到</span>
            <div class="social-icons">
              <button class="social-btn wechat" @click="shareToWechat" title="微信">
                <svg width="24" height="24" viewBox="0 0 24 24" fill="currentColor">
                  <path d="M8.5 13.5a1 1 0 100-2 1 1 0 000 2zm5 0a1 1 0 100-2 1 1 0 000 2z"/>
                  <path d="M12 2C6.477 2 2 6.145 2 11.243c0 2.936 1.526 5.55 3.926 7.227l-.928 2.503 3.377-1.44a11.9 11.9 0 002.625.387C11.075 19.92 11.515 20 12 20s.925-.08 1.375-.223c.886-.106 1.743-.274 2.625-.387l3.377 1.44-.928-2.503C20.474 16.793 22 14.18 22 11.243 22 6.145 17.523 2 12 2zm-5.5 5.586a.929.929 0 110 1.858.929.929 0 010-1.858zm5 0a.929.929 0 110 1.858.929.929 0 010-1.858z"/>
                </svg>
              </button>
              <button class="social-btn weibo" @click="shareToWeibo" title="微博">
                <svg width="24" height="24" viewBox="0 0 24 24" fill="currentColor">
                  <path d="M10.098 20c-4.612 0-8.348-2.727-8.348-6.086 0-1.77 1.135-3.466 3.195-5.06 2.252-1.74 4.762-2.456 6.022-2.456.29 0 .53.027.728.08.57-.36 1.082-.702 1.457-.976.44-.32.837-.607 1.238-.877C16.178 3.663 18.104 2 20.456 2c.46 0 .91.045 1.345.132-.874-.52-1.945-.82-3.13-.82-3.459 0-6.267 1.853-6.267 4.137 0 1.37.685 2.615 1.738 3.53-.044.002-.088.006-.133.006-.64 0-1.157-.32-1.157-.713 0-.55.762-.996 1.73-1.448 1.542-.72 2.873-1.778 2.873-3.39 0-1.93-1.57-3.495-3.502-3.495-2.023 0-3.71 1.32-3.71 2.947 0 1.027.533 1.945 1.12 2.58.07.076.09.2.05.304-.18.46-.48 1.1-.48 1.1s-.093.06-.09.148c.025.52.04 1.05.04 1.58 0 3.04-2.77 5.5-6.19 5.5-.47 0-.92-.04-1.35-.12-.05.11-.1.22-.16.32-.53.91-1.25 1.63-2.14 2.14.68.2 1.4.32 2.16.32 4.61 0 8.35-2.73 8.35-6.09 0 0 0 0 0 0 .53 2.34 2.44 4.09 4.93 4.09 1.8 0 3.42-.79 4.58-2.07.38.07.78.1 1.18.1 1.35 0 2.58-.34 3.64-.93a9.46 9.46 0 001.36-1.07c-.28.01-.56.02-.85.02-1.66 0-3.25-.36-4.7-1.04.25.02.5.03.76.03 1.35 0 2.68-.24 3.92-.69a14.12 14.12 0 01-3.34 2.17c-.4-.12-.82-.18-1.25-.18zm.38-6.86c-1.1 0-1.98-.7-1.98-1.57 0-.86.88-1.56 1.98-1.56 1.09 0 1.98.7 1.97 1.56 0 .87-.88 1.57-1.97 1.57z"/>
                </svg>
              </button>
              <button class="social-btn twitter" @click="shareToTwitter" title="Twitter/X">
                <svg width="24" height="24" viewBox="0 0 24 24" fill="currentColor">
                  <path d="M18.244 2.25h3.308l-7.227 8.26 8.502 11.24H16.17l-5.214-6.817L4.99 21.75H1.68l7.73-8.835L1.254 2.25H8.08l4.713 6.231zm-1.161 17.52h1.833L7.084 4.126H5.117z"/>
                </svg>
              </button>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useStyleStore } from '../../store/style'
import QRCode from 'qrcode'

const props = defineProps({
  visible: Boolean,
  post: Object
})

const emit = defineEmits(['close'])

const styleStore = useStyleStore()
const posterCanvas = ref(null)
const qrCanvas = ref(null)
const copied = ref(false)

const close = () => {
  emit('close')
}

const generatePoster = async () => {
  if (!posterCanvas.value || !props.post) return

  const canvas = posterCanvas.value
  const ctx = canvas.getContext('2d')

  // Poster dimensions - 2:3 aspect ratio for social media
  const width = 540
  const height = 810
  canvas.width = width
  canvas.height = height

  // Get theme colors
  const theme = styleStore.theme
  const accent = theme.colors.accent
  const bgColor = theme.colors.bg
  const textColor = theme.colors.text
  const headingColor = theme.colors.textHeading

  // Background gradient
  const gradient = ctx.createLinearGradient(0, 0, width, height)
  gradient.addColorStop(0, bgColor)
  gradient.addColorStop(1, '#f8f9ff')
  ctx.fillStyle = gradient
  ctx.fillRect(0, 0, width, height)

  // Decorative shapes
  ctx.save()
  ctx.globalAlpha = 0.1
  ctx.fillStyle = accent
  ctx.beginPath()
  ctx.arc(width - 50, 150, 200, 0, Math.PI * 2)
  ctx.fill()
  ctx.beginPath()
  ctx.arc(-100, height - 100, 250, 0, Math.PI * 2)
  ctx.fill()
  ctx.restore()

  // Header decoration line
  ctx.strokeStyle = accent
  ctx.lineWidth = 3
  ctx.beginPath()
  ctx.moveTo(40, 80)
  ctx.lineTo(width - 40, 80)
  ctx.stroke()

  // Blog title
  ctx.fillStyle = accent
  ctx.font = 'bold 24px "Noto Sans SC", sans-serif'
  ctx.textAlign = 'left'
  ctx.fillText('Jumoshen', 40, 50)

  // Article title
  ctx.fillStyle = headingColor
  ctx.font = 'bold 36px "Noto Sans SC", sans-serif'

  // Wrap long title
  const title = props.post.title || '无标题'
  const maxWidth = width - 80
  const titleLines = wrapText(ctx, title, maxWidth)
  let y = 160
  titleLines.forEach((line, i) => {
    ctx.fillText(line, 40, y + i * 48)
  })

  // Divider
  y += titleLines.length * 48 + 30
  ctx.strokeStyle = accent
  ctx.lineWidth = 2
  ctx.beginPath()
  ctx.moveTo(40, y)
  ctx.lineTo(width - 40, y)
  ctx.stroke()

  // Summary
  y += 40
  ctx.fillStyle = textColor
  ctx.font = '20px "Noto Sans SC", sans-serif'
  const summary = props.post.summary || ''
  const summaryLines = wrapText(ctx, summary, maxWidth)
  summaryLines.slice(0, 3).forEach((line, i) => {
    ctx.fillText(line, 40, y + i * 32)
  })

  // Tags
  y += 130
  if (props.post.tags && props.post.tags.length > 0) {
    ctx.font = '18px "Noto Sans SC", sans-serif'
    let x = 40
    props.post.tags.slice(0, 4).forEach(tag => {
      const tagWidth = ctx.measureText('#' + tag).width + 24
      if (x + tagWidth > width - 40) return

      ctx.fillStyle = accent
      ctx.globalAlpha = 0.15
      roundRect(ctx, x, y - 22, tagWidth, 36, 18)
      ctx.fill()
      ctx.globalAlpha = 1

      ctx.fillStyle = accent
      ctx.fillText('#' + tag, x + 12, y)
      x += tagWidth + 12
    })
  }

  // QR code placeholder area
  y = height - 220
  ctx.fillStyle = '#fff'
  roundRect(ctx, width/2 - 60, y, 120, 120, 12)
  ctx.fill()
  ctx.strokeStyle = accent
  ctx.lineWidth = 2
  ctx.stroke()

  // Generate real QR code
  const articleUrl = window.location.href
  try {
    await QRCode.toCanvas(document.createElement('canvas'), articleUrl, {
      width: 80,
      margin: 0,
      color: { dark: accent, light: '#ffffff' }
    }).then(qr => {
      ctx.drawImage(qr, width/2 - 40, y + 20, 80, 80)
    })
  } catch (e) {
    console.error('QR generation failed:', e)
  }

  // URL text
  ctx.fillStyle = textColor
  ctx.font = '16px "Noto Sans SC", sans-serif'
  ctx.textAlign = 'center'
  ctx.fillText('jumoshen.cn', width / 2, height - 60)
  ctx.fillText('扫码阅读完整文章', width / 2, height - 35)

  // Footer tagline
  ctx.fillStyle = accent
  ctx.font = 'bold 14px "Noto Sans SC", sans-serif'
  ctx.fillText('一加一永远大于二', width / 2, height - 15)
}

const wrapText = (ctx, text, maxWidth) => {
  const words = text.split('')
  const lines = []
  let currentLine = ''

  for (const char of words) {
    const testLine = currentLine + char
    const metrics = ctx.measureText(testLine)
    if (metrics.width > maxWidth && currentLine) {
      lines.push(currentLine)
      currentLine = char
    } else {
      currentLine = testLine
    }
  }
  if (currentLine) lines.push(currentLine)
  return lines
}

const roundRect = (ctx, x, y, w, h, r) => {
  ctx.beginPath()
  ctx.moveTo(x + r, y)
  ctx.lineTo(x + w - r, y)
  ctx.quadraticCurveTo(x + w, y, x + w, y + r)
  ctx.lineTo(x + w, y + h - r)
  ctx.quadraticCurveTo(x + w, y + h, x + w - r, y + h)
  ctx.lineTo(x + r, y + h)
  ctx.quadraticCurveTo(x, y + h, x, y + h - r)
  ctx.lineTo(x, y + r)
  ctx.quadraticCurveTo(x, y, x + r, y)
  ctx.closePath()
}

const downloadPoster = () => {
  if (!posterCanvas.value) return
  const link = document.createElement('a')
  link.download = `poster-${props.post?.slug || 'article'}.png`
  link.href = posterCanvas.value.toDataURL('image/png')
  link.click()
}

const copyLink = async () => {
  const url = window.location.href
  try {
    await navigator.clipboard.writeText(url)
    copied.value = true
    setTimeout(() => { copied.value = false }, 2000)
  } catch (e) {
    console.error('Copy failed:', e)
  }
}

const shareToWechat = () => {
  copyLink()
  alert('链接已复制，请粘贴到微信中分享')
}

const shareToWeibo = () => {
  const title = encodeURIComponent(props.post?.title || '')
  const url = encodeURIComponent(window.location.href)
  window.open(`https://service.weibo.com/share/share.php?title=${title}&url=${url}`, '_blank')
}

const shareToTwitter = () => {
  const title = encodeURIComponent(props.post?.title || '')
  const url = encodeURIComponent(window.location.href)
  window.open(`https://twitter.com/intent/tweet?text=${title}&url=${url}`, '_blank')
}

watch(() => props.visible, (val) => {
  if (val) {
    setTimeout(generatePoster, 100)
  }
})
</script>

<style scoped>
.share-modal {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.7);
  backdrop-filter: blur(8px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.share-container {
  background: var(--bg, #fff);
  border-radius: 20px;
  padding: 24px;
  width: 90%;
  max-width: 600px;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 25px 80px rgba(0, 0, 0, 0.3);
}

.share-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.share-header h3 {
  font-size: 20px;
  color: var(--text-h, #2a1a4a);
  margin: 0;
}

.close-btn {
  width: 36px;
  height: 36px;
  border: none;
  background: var(--accent-bg, rgba(179, 102, 255, 0.1));
  border-radius: 50%;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--accent, #b366ff);
  transition: all 0.2s;
}

.close-btn:hover {
  background: var(--accent, #b366ff);
  color: #fff;
  transform: rotate(90deg);
}

.poster-preview {
  background: #f0f0f0;
  border-radius: 12px;
  padding: 20px;
  display: flex;
  justify-content: center;
  margin-bottom: 20px;
}

.poster-canvas {
  max-width: 100%;
  height: auto;
  border-radius: 8px;
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.15);
}

.share-actions {
  display: flex;
  gap: 12px;
  margin-bottom: 20px;
}

.action-btn {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 14px 20px;
  border: none;
  border-radius: 12px;
  font-size: 15px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.action-btn.download {
  background: var(--accent, #b366ff);
  color: #fff;
}

.action-btn.download:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px var(--shadow, rgba(179, 102, 255, 0.4));
}

.action-btn.copy-link {
  background: var(--accent-bg, rgba(179, 102, 255, 0.1));
  color: var(--accent, #b366ff);
}

.action-btn.copy-link:hover {
  background: var(--accent, #b366ff);
  color: #fff;
}

.social-share {
  display: flex;
  align-items: center;
  gap: 16px;
  padding-top: 16px;
  border-top: 1px solid var(--border, #e5e4e7);
}

.share-label {
  font-size: 14px;
  color: var(--text, #6b6375);
}

.social-icons {
  display: flex;
  gap: 12px;
}

.social-btn {
  width: 44px;
  height: 44px;
  border: none;
  border-radius: 50%;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.social-btn.wechat {
  background: #07c160;
  color: #fff;
}

.social-btn.weibo {
  background: #e6162d;
  color: #fff;
}

.social-btn.twitter {
  background: #000;
  color: #fff;
}

.social-btn:hover {
  transform: scale(1.1);
}

.modal-fade-enter-active,
.modal-fade-leave-active {
  transition: all 0.3s ease;
}

.modal-fade-enter-from,
.modal-fade-leave-to {
  opacity: 0;
}

.modal-fade-enter-from .share-container,
.modal-fade-leave-to .share-container {
  transform: scale(0.9) translateY(20px);
}
</style>
