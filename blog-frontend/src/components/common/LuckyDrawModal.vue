<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="visible" class="modal-overlay" @click.self="handleClose">
        <div class="modal-container">
          <div class="modal-header">
            <h3>翻牌抽奖</h3>
            <button class="close-btn" @click="handleClose">&times;</button>
          </div>

          <div class="cards-container">
            <div
              v-for="(card, index) in cards"
              :key="index"
              class="card"
              :class="{ flipped: card.flipped }"
              @click="flipCard(index)"
            >
              <!-- 牌背面 -->
              <div class="card-face card-back">
                <div class="poker-pattern">
                  <span class="suit">&#9824;</span>
                </div>
              </div>

              <!-- 牌正面 -->
              <div class="card-face card-front">
                <div class="card-front-content">
                  <div class="card-title" @click.stop="goToPost(card)">
                    {{ card.title }}
                  </div>
                  <div class="card-date" v-if="card.date">{{ formatDate(card.date) }}</div>
                </div>
              </div>
            </div>
          </div>

          <p class="hint">点击任意一张牌翻开</p>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import api from '../../api'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['close'])

const router = useRouter()
const cards = ref([])

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN', { year: 'numeric', month: 'short', day: 'numeric' })
}

const fetchPosts = async () => {
  try {
    const res = await api.getPosts({ page: 1, page_size: 100 })
    if (res.code === 0 && res.data) {
      const posts = res.data.list || []
      // 随机选3篇
      const shuffled = [...posts].sort(() => Math.random() - 0.5)
      cards.value = shuffled.slice(0, 3).map(post => ({
        id: post.id,
        slug: post.slug,
        title: post.title || '无标题',
        date: post.date || '',
        flipped: false
      }))

      // 如果文章不足3篇，填充占位
      while (cards.value.length < 3) {
        cards.value.push({
          id: null,
          slug: null,
          title: '暂无文章',
          date: '',
          flipped: false
        })
      }
    }
  } catch (e) {
    console.error('Failed to fetch posts:', e)
    // 填充占位卡片
    cards.value = [
      { id: null, slug: null, title: '加载失败', date: '', flipped: false },
      { id: null, slug: null, title: '加载失败', date: '', flipped: false },
      { id: null, slug: null, title: '加载失败', date: '', flipped: false }
    ]
  }
}

const flipCard = (index) => {
  if (cards.value[index].flipped) return
  cards.value[index].flipped = true
}

const goToPost = (card) => {
  if (!card.slug) return
  emit('close')
  setTimeout(() => {
    window.scrollTo({ top: 0, behavior: 'instant' })
    const targetSlug = card.slug.startsWith('/') ? card.slug : `/post/${card.slug}`
    router.push(targetSlug)
  }, 300)
}

const handleClose = () => {
  emit('close')
}

watch(() => props.visible, (newVal) => {
  if (newVal) {
    // 重置所有卡片
    cards.value.forEach(card => card.flipped = false)
    fetchPosts()
  }
})
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.6);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
  backdrop-filter: blur(4px);
}

.modal-container {
  background: var(--card-bg);
  border: 2px solid var(--accent);
  border-radius: 20px;
  padding: 24px;
  max-width: 600px;
  width: 90%;
  box-shadow: 0 8px 32px var(--shadow);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.modal-header h3 {
  margin: 0;
  font-size: 20px;
  color: var(--text);
  font-family: var(--font-display);
}

.close-btn {
  background: none;
  border: none;
  font-size: 28px;
  color: var(--text);
  cursor: pointer;
  padding: 0;
  line-height: 1;
  opacity: 0.6;
  transition: opacity 0.2s;
}

.close-btn:hover {
  opacity: 1;
}

.cards-container {
  display: flex;
  justify-content: center;
  gap: 20px;
  perspective: 1000px;
}

.card {
  width: 140px;
  height: 200px;
  position: relative;
  transform-style: preserve-3d;
  transition: transform 0.8s cubic-bezier(0.4, 0.0, 0.2, 1);
  cursor: pointer;
}

.card.flipped {
  transform: rotateY(180deg);
}

.card-face {
  position: absolute;
  width: 100%;
  height: 100%;
  backface-visibility: hidden;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 2px solid var(--accent);
}

.card-back {
  background: linear-gradient(145deg, var(--card-bg) 0%, var(--accent-bg) 100%);
  transform: rotateY(0deg);
}

.card-front {
  background: linear-gradient(145deg, #ffffff 0%, var(--card-bg) 100%);
  transform: rotateY(180deg);
  padding: 16px;
}

.card-front-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  gap: 8px;
}

.poker-pattern {
  width: 80px;
  height: 112px;
  background: var(--accent);
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: inset 0 2px 8px rgba(0,0,0,0.1);
}

.suit {
  font-size: 48px;
  color: white;
  text-shadow: 0 2px 4px rgba(0,0,0,0.2);
}

.card-title {
  font-size: 13px;
  color: var(--text);
  text-align: center;
  word-break: break-word;
  line-height: 1.4;
  cursor: pointer;
  transition: color 0.2s;
  max-height: 120px;
  overflow: hidden;
  display: -webkit-box;
  -webkit-line-clamp: 5;
  -webkit-box-orient: vertical;
}

.card-title:hover {
  color: var(--accent);
}

.card-date {
  font-size: 10px;
  color: var(--text);
  opacity: 0.6;
}

.hint {
  text-align: center;
  margin-top: 20px;
  font-size: 12px;
  color: var(--text);
  opacity: 0.6;
}

/* 主题适配 */
.theme-pixel .card-face {
  border-radius: 4px;
  border-width: 2px;
}

.theme-cute .card-face {
  border-radius: 16px;
  border: none;
  box-shadow: 0 4px 16px rgba(255, 107, 157, 0.2);
}

.theme-qver .card-face {
  border-radius: 16px;
  border: 2px solid rgba(124, 106, 255, 0.3);
  box-shadow: 0 6px 20px rgba(124, 106, 255, 0.2);
}

/* 主题适配 - 扑克牌背面样式 */
.theme-pixel .card-back {
  background:
    repeating-linear-gradient(
      0deg,
      #ff6b9d 0px,
      #ff6b9d 10px,
      #c44569 10px,
      #c44569 20px
    ),
    repeating-linear-gradient(
      90deg,
      #ff6b9d 0px,
      #ff6b9d 10px,
      #c44569 10px,
      #c44569 20px
    );
  border: 4px solid #2d2d2d;
  box-shadow: 4px 4px 0 #2d2d2d;
}

.theme-pixel .card-back .poker-pattern {
  width: 60px;
  height: 84px;
  background: #2d2d2d;
  border: 3px solid #fff;
  box-shadow: inset 0 0 0 2px #2d2d2d;
}

.theme-pixel .card-back .suit {
  font-size: 36px;
  color: #ff6b9d;
  text-shadow: none;
}

.theme-cute .card-back {
  background: linear-gradient(135deg, #ffecd2 0%, #fcb69f 100%);
  border: none;
  border-radius: 20px;
  box-shadow: 0 8px 32px rgba(255, 107, 157, 0.3);
}

.theme-cute .card-back .poker-pattern {
  width: 60px;
  height: 84px;
  background: linear-gradient(135deg, #ff8ab5 0%, #ff6b9d 100%);
  border-radius: 12px;
  box-shadow: 0 4px 12px rgba(255, 107, 157, 0.4);
}

.theme-cute .card-back .suit {
  font-size: 40px;
  background: linear-gradient(135deg, #fff 0%, #fff5f7 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.theme-qver .card-back {
  background:
    linear-gradient(135deg, rgba(124, 106, 255, 0.1) 0%, transparent 50%),
    linear-gradient(225deg, rgba(124, 106, 255, 0.1) 0%, transparent 50%),
    linear-gradient(135deg, #1a1a2e 0%, #16213e 100%);
  border: 2px solid rgba(124, 106, 255, 0.6);
  border-radius: 16px;
  box-shadow:
    0 0 20px rgba(124, 106, 255, 0.4),
    inset 0 0 30px rgba(124, 106, 255, 0.1);
}

.theme-qver .card-back .poker-pattern {
  width: 60px;
  height: 84px;
  background: rgba(124, 106, 255, 0.2);
  border: 2px solid rgba(124, 106, 255, 0.6);
  border-radius: 8px;
  box-shadow:
    0 0 15px rgba(124, 106, 255, 0.3),
    inset 0 0 15px rgba(124, 106, 255, 0.2);
}

.theme-qver .card-back .suit {
  font-size: 38px;
  color: #a78bfa;
  text-shadow: 0 0 10px #a78bfa, 0 0 20px #7c6aff, 0 0 30px #7c6aff;
}

/* Modal 动画 */
.modal-enter-active,
.modal-leave-active {
  transition: all 0.3s cubic-bezier(0.4, 0.0, 0.2, 1);
}

.modal-enter-active .modal-container,
.modal-leave-active .modal-container {
  transition: all 0.3s cubic-bezier(0.4, 0.0, 0.2, 1);
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-from .modal-container,
.modal-leave-to .modal-container {
  transform: translateY(100px);
}
</style>
