<template>
  <Teleport to="body">
    <Transition name="hint-fade">
      <div v-if="showHint" class="keyboard-hint">
        <div class="hint-content">
          <div class="hint-header">
            <kbd>J</kbd> / <kbd>K</kbd> 导航文章
          </div>
          <div class="hint-keys">
            <span class="hint-item"><kbd>J</kbd> 下一篇</span>
            <span class="hint-item"><kbd>K</kbd> 上一篇</span>
            <span class="hint-item"><kbd>T</kbd> 回到顶部</span>
            <span class="hint-item"><kbd>?</kbd> 帮助</span>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'

const showHint = ref(false)
let hintTimeout = null

const handleKeydown = (e) => {
  // Don't trigger if user is typing in an input
  if (e.target.tagName === 'INPUT' || e.target.tagName === 'TEXTAREA' || e.target.isContentEditable) {
    return
  }

  switch (e.key.toLowerCase()) {
    case 'j':
      // Navigate to next post
      navigatePost('next')
      break
    case 'k':
      // Navigate to previous post
      navigatePost('prev')
      break
    case 't':
      // Scroll to top
      window.scrollTo({ top: 0, behavior: 'smooth' })
      break
    case '?':
      showHint.value = !showHint.value
      if (showHint.value) {
        clearTimeout(hintTimeout)
        hintTimeout = setTimeout(() => {
          showHint.value = false
        }, 5000)
      }
      break
    case 'escape':
      showHint.value = false
      break
  }
}

const navigatePost = (direction) => {
  // Emit event for Post.vue to handle
  window.dispatchEvent(new CustomEvent('keyboard-nav', { detail: { direction } }))
}

onMounted(() => {
  window.addEventListener('keydown', handleKeydown)
  // Show hint briefly on mount
  setTimeout(() => {
    showHint.value = true
    hintTimeout = setTimeout(() => {
      showHint.value = false
    }, 4000)
  }, 2000)
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeydown)
  clearTimeout(hintTimeout)
})
</script>

<style scoped>
.keyboard-hint {
  position: fixed;
  bottom: 100px;
  left: 50%;
  transform: translateX(-50%);
  z-index: 999;
  pointer-events: none;
}

.hint-content {
  background: var(--card-bg, #fff);
  border: 1px solid var(--border, #e5e4e7);
  border-radius: 16px;
  padding: 16px 24px;
  box-shadow: 0 12px 40px var(--shadow, rgba(0,0,0,0.15));
  backdrop-filter: blur(10px);
}

.hint-header {
  font-size: 13px;
  color: var(--text, #6b6375);
  margin-bottom: 12px;
  text-align: center;
}

.hint-keys {
  display: flex;
  gap: 16px;
}

.hint-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: var(--text-h, #2a1a4a);
}

kbd {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 24px;
  height: 24px;
  padding: 0 6px;
  font-size: 11px;
  font-family: inherit;
  font-weight: 600;
  color: var(--accent, #b366ff);
  background: var(--accent-bg, rgba(179, 102, 255, 0.1));
  border: 1px solid var(--accent-border, rgba(179, 102, 255, 0.3));
  border-radius: 6px;
  box-shadow: 0 2px 0 var(--accent, #b366ff);
}

.hint-fade-enter-active,
.hint-fade-leave-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.hint-fade-enter-from,
.hint-fade-leave-to {
  opacity: 0;
  transform: translateX(-50%) translateY(20px);
}
</style>
