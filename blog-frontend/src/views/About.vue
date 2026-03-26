<template>
  <div class="about" :class="'theme-' + styleStore.currentTheme">
    <!-- 背景装饰 -->
    <div class="bg-orbs">
      <div class="orb orb-1"></div>
      <div class="orb orb-2"></div>
      <div class="orb orb-3"></div>
    </div>

    <div class="about-card" v-if="!loading">
      <!-- 头像区域 -->
      <div class="hero-section">
        <div class="avatar-wrapper">
          <div class="avatar-ring"></div>
          <div class="avatar">
            <span class="avatar-icon">👨‍💻</span>
          </div>
        </div>
        <div class="title-area">
          <h1 class="title">About Me</h1>
          <p class="subtitle">
            <span class="typing-text">{{ typingText }}</span>
            <span class="cursor">|</span>
          </p>
        </div>
      </div>

      <div class="divider">
        <span class="divider-icon">✦</span>
      </div>

      <!-- 内容区域 -->
      <div class="content-section" v-if="renderedContent" v-html="renderedContent"></div>

      <div class="content-section" v-else>
        <div class="info-grid">
          <div
            v-for="(item, index) in defaultContent"
            :key="index"
            class="info-card"
            :style="{ '--delay': index * 0.15 + 's' }"
          >
            <span class="info-icon">{{ item.icon }}</span>
            <span class="info-text">{{ item.text }}</span>
          </div>
        </div>
      </div>

      <!-- 技能标签 -->
      <div class="skills-section">
        <h3 class="skills-title">技术栈</h3>
        <div class="skills-tags">
          <span class="skill-tag" v-for="skill in skills" :key="skill">{{ skill }}</span>
        </div>
      </div>

      <!-- 社交链接 -->
      <div class="social-section">
        <a href="mailto:job.lianghui@gmail.com" class="social-btn email">
          <span class="social-icon">📧</span>
          <span>发邮件</span>
        </a>
        <a href="https://github.com/jumoshen" target="_blank" class="social-btn github">
          <span class="social-icon">🐙</span>
          <span>GitHub</span>
        </a>
      </div>
    </div>

    <!-- 加载状态 -->
    <div v-else class="loading">
      <div class="loader"></div>
      <p>加载中...</p>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { marked } from 'marked'
import api from '../api'
import { useStyleStore } from '../store/style'

const styleStore = useStyleStore()
const content = ref('')
const loading = ref(true)
const typingText = ref('')

const defaultContent = [
  { icon: '👋', text: 'Hi, I\'m liang hui' },
  { icon: '👀', text: 'I\'m a php/golang programmer' },
  { icon: '📫', text: 'Connect me: job.lianghui@gmail.com' },
]

const skills = ['PHP', 'Go', 'Golang', 'Vue.js', 'React', 'Docker', 'MySQL', 'Redis', 'Linux']

// 打字机效果
const fullText = 'A programmer who loves coding'

onMounted(async () => {
  try {
    const res = await api.getAbout()
    if (res.code === 0 && res.data?.content && !res.data.content.includes('not found')) {
      content.value = res.data.content
    }
  } catch (e) {
    console.error('Failed to load about:', e)
  } finally {
    loading.value = false
  }

  // 启动打字机效果
  let index = 0
  const typeInterval = setInterval(() => {
    if (index <= fullText.length) {
      typingText.value = fullText.slice(0, index)
      index++
    } else {
      clearInterval(typeInterval)
    }
  }, 80)
})

const renderedContent = computed(() => {
  if (!content.value) return ''
  return marked.parse(content.value)
})
</script>

<style scoped>
.about {
  min-height: calc(100vh - 120px);
  padding: 60px 20px;
  position: relative;
  overflow: hidden;
}

/* 背景装饰 */
.bg-orbs {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
  z-index: 0;
}

.orb {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  opacity: 0.4;
}

.orb-1 {
  width: 400px;
  height: 400px;
  background: var(--accent);
  top: -150px;
  right: -100px;
  animation: floatOrb 10s ease-in-out infinite;
}

.orb-2 {
  width: 300px;
  height: 300px;
  background: #f472b6;
  bottom: -100px;
  left: -100px;
  animation: floatOrb 12s ease-in-out infinite reverse;
}

.orb-3 {
  width: 200px;
  height: 200px;
  background: #a78bfa;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  animation: pulseOrb 6s ease-in-out infinite;
}

@keyframes floatOrb {
  0%, 100% { transform: translate(0, 0); }
  50% { transform: translate(30px, -30px); }
}

@keyframes pulseOrb {
  0%, 100% { transform: translate(-50%, -50%) scale(1); opacity: 0.4; }
  50% { transform: translate(-50%, -50%) scale(1.2); opacity: 0.2; }
}

/* 卡片主体 */
.about-card {
  max-width: 700px;
  margin: 0 auto;
  background: var(--card-bg);
  border-radius: 32px;
  padding: 48px;
  border: 1px solid var(--border);
  box-shadow: 0 20px 60px var(--shadow);
  position: relative;
  z-index: 1;
  animation: slideUp 0.8s ease-out;
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(40px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* Hero 区域 */
.hero-section {
  text-align: center;
  margin-bottom: 32px;
}

.avatar-wrapper {
  position: relative;
  display: inline-block;
  margin-bottom: 24px;
}

.avatar-ring {
  position: absolute;
  top: -8px;
  left: -8px;
  right: -8px;
  bottom: -8px;
  border: 3px solid var(--accent);
  border-radius: 50%;
  animation: rotateRing 8s linear infinite;
}

@keyframes rotateRing {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.avatar {
  width: 120px;
  height: 120px;
  background: linear-gradient(135deg, var(--accent), var(--accent-border));
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 8px 32px var(--shadow);
}

.avatar-icon {
  font-size: 56px;
}

.title-area {
  margin-top: 16px;
}

.title {
  font-size: 36px;
  font-weight: 800;
  color: var(--text-h);
  margin: 0 0 12px;
  background: linear-gradient(135deg, var(--text-h), var(--accent));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.subtitle {
  font-size: 16px;
  color: var(--text);
  margin: 0;
  min-height: 24px;
}

.typing-text {
  color: var(--accent);
  font-weight: 500;
}

.cursor {
  animation: blink 0.8s infinite;
  color: var(--accent);
}

@keyframes blink {
  0%, 100% { opacity: 1; }
  50% { opacity: 0; }
}

/* 分隔线 */
.divider {
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 32px 0;
  gap: 16px;
}

.divider::before,
.divider::after {
  content: '';
  flex: 1;
  height: 1px;
  background: linear-gradient(90deg, transparent, var(--border));
}

.divider::after {
  background: linear-gradient(90deg, var(--border), transparent);
}

.divider-icon {
  color: var(--accent);
  font-size: 12px;
  animation: pulse 2s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; transform: scale(1); }
  50% { opacity: 0.5; transform: scale(1.2); }
}

/* 内容区域 */
.content-section {
  margin-bottom: 32px;
}

.info-grid {
  display: grid;
  gap: 16px;
}

.info-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px 24px;
  background: var(--accent-bg);
  border-radius: 16px;
  border: 1px solid var(--accent-border);
  animation: fadeInUp 0.6s ease-out forwards;
  animation-delay: var(--delay);
  opacity: 0;
  transition: all 0.3s ease;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.info-card:hover {
  transform: translateX(8px);
  border-color: var(--accent);
  box-shadow: 0 4px 20px var(--shadow);
}

.info-icon {
  font-size: 28px;
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--card-bg);
  border-radius: 12px;
}

.info-text {
  font-size: 16px;
  color: var(--text);
  line-height: 1.5;
}

/* 技能标签 */
.skills-section {
  margin-bottom: 32px;
}

.skills-title {
  font-size: 14px;
  color: var(--text);
  text-transform: uppercase;
  letter-spacing: 2px;
  margin: 0 0 16px;
  text-align: center;
}

.skills-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  justify-content: center;
}

.skill-tag {
  padding: 8px 16px;
  background: var(--card-bg);
  border: 1px solid var(--border);
  border-radius: 20px;
  font-size: 13px;
  color: var(--text);
  transition: all 0.3s ease;
}

.skill-tag:hover {
  background: var(--accent);
  color: white;
  border-color: var(--accent);
  transform: scale(1.05);
}

/* 社交链接 */
.social-section {
  display: flex;
  gap: 16px;
  justify-content: center;
}

.social-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 24px;
  border-radius: 50px;
  text-decoration: none;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
}

.social-btn.email {
  background: var(--accent);
  color: white;
}

.social-btn.github {
  background: var(--card-bg);
  color: var(--text);
  border: 1px solid var(--border);
}

.social-btn:hover {
  transform: translateY(-4px) scale(1.05);
  box-shadow: 0 8px 24px var(--shadow);
}

.social-icon {
  font-size: 18px;
}

/* 加载状态 */
.loading {
  text-align: center;
  padding: 100px 40px;
}

.loader {
  width: 48px;
  height: 48px;
  border: 4px solid var(--accent-bg);
  border-top-color: var(--accent);
  border-radius: 50%;
  margin: 0 auto 20px;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.loading p {
  color: var(--text);
  font-size: 14px;
}

/* 内容样式 */
.content-section :deep(h1),
.content-section :deep(h2),
.content-section :deep(h3) {
  color: var(--text-h);
  margin: 20px 0 12px;
}

.content-section :deep(p) {
  color: var(--text);
  line-height: 1.8;
  margin: 12px 0;
}

.content-section :deep(ul) {
  padding-left: 20px;
  color: var(--text);
}

.content-section :deep(li) {
  margin: 8px 0;
  line-height: 1.6;
}

.content-section :deep(a) {
  color: var(--accent);
  text-decoration: none;
}

.content-section :deep(a:hover) {
  text-decoration: underline;
}

/* 响应式 */
@media (max-width: 768px) {
  .about-card {
    padding: 32px 24px;
    border-radius: 24px;
  }

  .title {
    font-size: 28px;
  }

  .avatar {
    width: 100px;
    height: 100px;
  }

  .avatar-icon {
    font-size: 48px;
  }

  .social-section {
    flex-direction: column;
  }
}
</style>
