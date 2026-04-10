<template>
  <Teleport to="body">
    <Transition name="modal-pop">
      <div v-if="visible" class="donation-modal" @click.self="close">
        <div class="donation-container">
          <!-- Close button -->
          <button class="close-btn" @click="close">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M18 6L6 18M6 6l12 12"/>
            </svg>
          </button>

          <!-- Header -->
          <div class="donation-header">
            <div class="coffee-icon">
              <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path d="M18 8h1a4 4 0 010 8h-1M2 8h16v9a4 4 0 01-4 4H6a4 4 0 01-4-4V8z"/>
                <path d="M6 1v3M10 1v3M14 1v3"/>
              </svg>
            </div>
            <h2>请作者喝杯咖啡</h2>
            <p class="subtitle">您的支持是我创作的最大动力</p>
          </div>

          <!-- Donation Methods -->
          <div class="donation-methods">
            <!-- WeChat Pay -->
            <div
              class="method-card"
              :class="{ active: selectedMethod === 'wechat' }"
              @click="selectedMethod = 'wechat'"
            >
              <div class="method-icon wechat">
                <svg width="28" height="28" viewBox="0 0 24 24" fill="currentColor">
                  <path d="M8.5 13.5a1 1 0 100-2 1 1 0 000 2zm5 0a1 1 0 100-2 1 1 0 000 2z"/>
                  <path d="M12 2C6.477 2 2 6.145 2 11.243c0 2.936 1.526 5.55 3.926 7.227l-.928 2.503 3.377-1.44a11.9 11.9 0 002.625.387C11.075 19.92 11.515 20 12 20s.925-.08 1.375-.223c.886-.106 1.743-.274 2.625-.387l3.377 1.44-.928-2.503C20.474 16.793 22 14.18 22 11.243 22 6.145 17.523 2 12 2z"/>
                </svg>
              </div>
              <span>微信支付</span>
              <div class="check-icon" v-if="selectedMethod === 'wechat'">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
                  <path d="M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41z"/>
                </svg>
              </div>
            </div>

            <!-- Alipay -->
            <div
              class="method-card"
              :class="{ active: selectedMethod === 'alipay' }"
              @click="selectedMethod = 'alipay'"
            >
              <div class="method-icon alipay">
                <svg width="28" height="28" viewBox="0 0 24 24" fill="currentColor">
                  <path d="M21.93 8.44c.27.58.37 1.18.37 1.8 0 2.84-2.75 5.14-6.14 5.14H8.07c-.87 0-1.58-.71-1.58-1.58 0-.71.47-1.31 1.12-1.5l.97-.29.37-.11c1.22-.36 2.23-.9 3.02-1.62.84-.76 1.31-1.72 1.31-2.82 0-.41-.06-.81-.17-1.2-.12-.41-.3-.81-.54-1.18-.24-.36-.54-.7-.89-.98-.36-.27-.78-.49-1.22-.63-.45-.14-.93-.21-1.42-.21-.69 0-1.35.12-1.97.36-.61.24-1.15.58-1.61 1.01-.45.43-.81.94-1.06 1.53-.25.59-.37 1.24-.37 1.93H4c0-.73.12-1.44.36-2.12.24-.68.58-1.28 1.01-1.79.43-.51.95-.91 1.54-1.2.59-.29 1.24-.44 1.94-.44.88 0 1.69.21 2.42.63.74.42 1.34.99 1.79 1.7.45.72.72 1.54.8 2.48.06.6-.01 1.16-.22 1.68-.2.53-.52.99-.95 1.38.59.19 1.1.51 1.54.97.44.46.74.99.89 1.6.15.61.13 1.25-.06 1.93h2.17c.73 0 1.33-.6 1.33-1.33 0-.16-.03-.32-.08-.47z"/>
                </svg>
              </div>
              <span>支付宝</span>
              <div class="check-icon" v-if="selectedMethod === 'alipay'">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
                  <path d="M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41z"/>
                </svg>
              </div>
            </div>
          </div>

          <!-- QR Code Display -->
          <div class="qr-display">
            <Transition name="fade" mode="out-in">
              <div v-if="selectedMethod === 'wechat'" key="wechat" class="qr-card">
                <div class="qr-label">微信扫码支付</div>
                <div class="qr-frame wechat">
                  <div class="qr-inner">
                    <!-- WeChat QR placeholder -->
                    <svg width="180" height="180" viewBox="0 0 180 180">
                      <rect width="180" height="180" fill="#f5f5f5" rx="8"/>
                      <g fill="#07c160">
                        <rect x="20" y="20" width="50" height="50" rx="4"/>
                        <rect x="110" y="20" width="50" height="50" rx="4"/>
                        <rect x="20" y="110" width="50" height="50" rx="4"/>
                        <rect x="30" y="30" width="30" height="30" fill="#fff"/>
                        <rect x="120" y="30" width="30" height="30" fill="#fff"/>
                        <rect x="30" y="120" width="30" height="30" fill="#fff"/>
                        <rect x="80" y="80" width="20" height="20" rx="2"/>
                        <rect x="40" y="40" width="10" height="10" fill="#07c160"/>
                        <rect x="130" y="40" width="10" height="10" fill="#07c160"/>
                        <rect x="40" y="130" width="10" height="10" fill="#07c160"/>
                      </g>
                    </svg>
                  </div>
                </div>
                <p class="qr-hint">打开微信扫一扫</p>
              </div>

              <div v-else key="alipay" class="qr-card">
                <div class="qr-label">支付宝扫码支付</div>
                <div class="qr-frame alipay">
                  <div class="qr-inner">
                    <!-- Alipay QR placeholder -->
                    <svg width="180" height="180" viewBox="0 0 180 180">
                      <rect width="180" height="180" fill="#f5f5f5" rx="8"/>
                      <g fill="#1677ff">
                        <rect x="20" y="20" width="50" height="50" rx="4"/>
                        <rect x="110" y="20" width="50" height="50" rx="4"/>
                        <rect x="20" y="110" width="50" height="50" rx="4"/>
                        <circle cx="55" cy="55" r="15" fill="#fff"/>
                        <circle cx="135" cy="55" r="15" fill="#fff"/>
                        <circle cx="55" cy="145" r="15" fill="#fff"/>
                        <rect x="85" y="75" width="10" height="30" rx="2"/>
                        <rect x="75" y="85" width="30" height="10" rx="2"/>
                      </g>
                    </svg>
                  </div>
                </div>
                <p class="qr-hint">打开支付宝扫一扫</p>
              </div>
            </Transition>
          </div>

          <!-- Amount Selection -->
          <div class="amount-section">
            <p class="section-label">自定义金额</p>
            <div class="amount-options">
              <button
                v-for="amt in [5, 10, 20, 50, 100]"
                :key="amt"
                class="amount-btn"
                :class="{ active: customAmount === amt }"
                @click="customAmount = amt"
              >
                ¥{{ amt }}
              </button>
            </div>
          </div>

          <!-- Motivation Message -->
          <div class="motivation">
            <div class="motivation-icon">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
                <path d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z"/>
              </svg>
            </div>
            <p>{{ currentMessage }}</p>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup>
import { ref, watch } from 'vue'

const props = defineProps({
  visible: Boolean,
  post: Object
})

const emit = defineEmits(['close'])

const selectedMethod = ref('wechat')
const customAmount = ref(10)

const messages = [
  '您的支持是我创作的最大动力 💪',
  '感谢您的认可，一起让知识更有价值 ✨',
  '每一份支持都让我更有动力分享 🚀',
  '谢谢您的慷慨，世界因您更美好 🌟',
  '您的打赏是对作者最大的鼓励 🎉'
]

const currentMessage = ref(messages[Math.floor(Math.random() * messages.length)])

watch(() => props.visible, (val) => {
  if (val) {
    currentMessage.value = messages[Math.floor(Math.random() * messages.length)]
  }
})

const close = () => {
  emit('close')
}
</script>

<style scoped>
.donation-modal {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.65);
  backdrop-filter: blur(12px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 20px;
}

.donation-container {
  background: var(--bg, #fff);
  border-radius: 24px;
  padding: 32px;
  width: 100%;
  max-width: 420px;
  position: relative;
  box-shadow:
    0 25px 80px rgba(0, 0, 0, 0.3),
    0 0 0 1px rgba(255, 255, 255, 0.1);
}

.close-btn {
  position: absolute;
  top: 16px;
  right: 16px;
  width: 40px;
  height: 40px;
  border: none;
  background: var(--accent-bg, rgba(179, 102, 255, 0.1));
  border-radius: 50%;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text, #6b6375);
  transition: all 0.25s;
}

.close-btn:hover {
  background: var(--accent, #b366ff);
  color: #fff;
  transform: rotate(90deg);
}

.donation-header {
  text-align: center;
  margin-bottom: 28px;
}

.coffee-icon {
  width: 80px;
  height: 80px;
  margin: 0 auto 16px;
  background: linear-gradient(135deg, #ff6b9d 0%, #ffa8c5 100%);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  box-shadow: 0 12px 40px rgba(255, 107, 157, 0.4);
  animation: float 3s ease-in-out infinite;
}

@keyframes float {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-8px); }
}

.donation-header h2 {
  font-size: 24px;
  font-weight: 700;
  color: var(--text-h, #2a1a4a);
  margin: 0 0 8px;
}

.subtitle {
  font-size: 14px;
  color: var(--text, #6b6375);
  margin: 0;
}

.donation-methods {
  display: flex;
  gap: 12px;
  margin-bottom: 24px;
}

.method-card {
  flex: 1;
  padding: 16px;
  border: 2px solid var(--border, #e5e4e7);
  border-radius: 16px;
  display: flex;
  align-items: center;
  gap: 12px;
  cursor: pointer;
  transition: all 0.25s;
  position: relative;
}

.method-card:hover {
  border-color: var(--accent, #b366ff);
  background: var(--accent-bg, rgba(179, 102, 255, 0.05));
}

.method-card.active {
  border-color: var(--accent, #b366ff);
  background: var(--accent-bg, rgba(179, 102, 255, 0.1));
}

.method-icon {
  width: 44px;
  height: 44px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
}

.method-icon.wechat {
  background: linear-gradient(135deg, #07c160 0%, #00b55a 100%);
}

.method-icon.alipay {
  background: linear-gradient(135deg, #1677ff 0%, #0958d9 100%);
}

.method-card span {
  font-size: 15px;
  font-weight: 500;
  color: var(--text-h, #2a1a4a);
}

.check-icon {
  position: absolute;
  top: -8px;
  right: -8px;
  width: 24px;
  height: 24px;
  background: var(--accent, #b366ff);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
}

.qr-display {
  background: var(--accent-bg, rgba(179, 102, 255, 0.05));
  border-radius: 20px;
  padding: 24px;
  margin-bottom: 24px;
}

.qr-card {
  text-align: center;
}

.qr-label {
  font-size: 14px;
  color: var(--text, #6b6375);
  margin-bottom: 16px;
}

.qr-frame {
  display: inline-block;
  padding: 12px;
  border-radius: 16px;
}

.qr-frame.wechat {
  background: linear-gradient(135deg, rgba(7, 193, 96, 0.1) 0%, rgba(0, 181, 90, 0.1) 100%);
  border: 2px solid rgba(7, 193, 96, 0.3);
}

.qr-frame.alipay {
  background: linear-gradient(135deg, rgba(22, 119, 255, 0.1) 0%, rgba(9, 88, 217, 0.1) 100%);
  border: 2px solid rgba(22, 119, 255, 0.3);
}

.qr-inner {
  border-radius: 8px;
  overflow: hidden;
}

.qr-hint {
  margin: 12px 0 0;
  font-size: 13px;
  color: var(--text, #6b6375);
}

.amount-section {
  margin-bottom: 20px;
}

.section-label {
  font-size: 13px;
  color: var(--text, #6b6375);
  margin: 0 0 12px;
}

.amount-options {
  display: flex;
  gap: 8px;
}

.amount-btn {
  flex: 1;
  padding: 12px 8px;
  border: 2px solid var(--border, #e5e4e7);
  background: transparent;
  border-radius: 12px;
  font-size: 15px;
  font-weight: 600;
  color: var(--text-h, #2a1a4a);
  cursor: pointer;
  transition: all 0.2s;
}

.amount-btn:hover {
  border-color: var(--accent, #b366ff);
  color: var(--accent, #b366ff);
}

.amount-btn.active {
  border-color: var(--accent, #b366ff);
  background: var(--accent, #b366ff);
  color: #fff;
}

.motivation {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  background: linear-gradient(135deg, rgba(255, 107, 157, 0.08) 0%, rgba(179, 102, 255, 0.08) 100%);
  border-radius: 12px;
}

.motivation-icon {
  width: 36px;
  height: 36px;
  background: linear-gradient(135deg, #ff6b9d 0%, #b366ff 100%);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  flex-shrink: 0;
  animation: heartbeat 1.5s ease-in-out infinite;
}

@keyframes heartbeat {
  0%, 100% { transform: scale(1); }
  50% { transform: scale(1.1); }
}

.motivation p {
  font-size: 14px;
  color: var(--text, #6b6375);
  margin: 0;
  line-height: 1.5;
}

/* Transitions */
.modal-pop-enter-active,
.modal-pop-leave-active {
  transition: all 0.35s cubic-bezier(0.4, 0, 0.2, 1);
}

.modal-pop-enter-from,
.modal-pop-leave-to {
  opacity: 0;
}

.modal-pop-enter-from .donation-container,
.modal-pop-leave-to .donation-container {
  transform: scale(0.85) translateY(30px);
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s, transform 0.2s;
}

.fade-enter-from {
  opacity: 0;
  transform: translateX(10px);
}

.fade-leave-to {
  opacity: 0;
  transform: translateX(-10px);
}
</style>
