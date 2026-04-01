<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="visible" class="modal-overlay" @click.self="handleClose">
        <div class="modal-container">
          <div class="modal-header">
            <h3>翻牌抽奖</h3>
            <button class="close-btn" @click="handleClose">&times;</button>
          </div>

          <div class="cards-container" :class="currentTheme">
            <div
              v-for="(card, index) in cards"
              :key="index"
              class="card"
              :class="{ flipped: card.flipped }"
              @click="flipCard(index)"
            >
              <!-- 牌背面 -->
              <div class="card-face card-back" :style="{ backgroundImage: `url('/card-backs/${card.backImage}')` }">
                <div class="poker-pattern">
                  <span class="suit">&#9824;</span>
                </div>
              </div>

              <!-- 牌正面 -->
              <div class="card-face card-front" :style="{ backgroundImage: `url('/card-backs/${card.frontImage}')` }">
                <div class="card-front-content">
                  <div class="card-title" @click.stop="goToPost(card)">
                    {{ card.title }}
                  </div>
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
import { ref, watch, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useStyleStore } from '../../store/style'
import api from '../../api'

const styleStore = useStyleStore()
const currentTheme = computed(() => 'theme-' + styleStore.currentTheme)

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

// 根据主题获取背面图片列表
const getBackImages = (theme) => {
  const images = {
    'theme-pixel': ['back_1.png', 'back_4.png'],
    'theme-cute': ['back_2.jpg', 'back_5.png'],
    'theme-qver': ['back_3.png', 'back_6.png']
  }
  return images[theme] || images['theme-pixel']
}

// 正面图片列表
const frontImages = ['fron_1.png', 'fron_2.png', 'fron_3.png', 'fron_4.png']

// 随机从列表中选择一个
const randomPick = (list) => list[Math.floor(Math.random() * list.length)]

const fetchPosts = async () => {
  try {
    const res = await api.getPosts({ page: 1, page_size: 100 })
    if (res.code === 0 && res.data) {
      const posts = res.data.list || []
      // 随机选3篇
      const shuffled = [...posts].sort(() => Math.random() - 0.5)
      const backImages = getBackImages(currentTheme.value)

      cards.value = shuffled.slice(0, 3).map((post, index) => ({
        id: post.id,
        slug: post.slug,
        title: post.title || '无标题',
        date: post.date || '',
        flipped: false,
        backImage: randomPick(backImages), // 随机背面图片
        frontImage: randomPick(frontImages) // 随机正面图片
      }))

      // 如果文章不足3篇，填充占位
      while (cards.value.length < 3) {
        cards.value.push({
          id: null,
          slug: null,
          title: '暂无文章',
          date: '',
          flipped: false,
          backImage: randomPick(backImages),
          frontImage: randomPick(frontImages)
        })
      }
    }
  } catch (e) {
    console.error('Failed to fetch posts:', e)
    const backImages = getBackImages(currentTheme.value)
    // 填充占位卡片
    cards.value = [
      { id: null, slug: null, title: '加载失败', date: '', flipped: false, backImage: randomPick(backImages), frontImage: randomPick(frontImages) },
      { id: null, slug: null, title: '加载失败', date: '', flipped: false, backImage: randomPick(backImages), frontImage: randomPick(frontImages) },
      { id: null, slug: null, title: '加载失败', date: '', flipped: false, backImage: randomPick(backImages), frontImage: randomPick(frontImages) }
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
  background-size: cover;
  background-position: center;
  transform: rotateY(0deg);
}

.card-front {
  background: linear-gradient(145deg, #ffffff 0%, var(--card-bg) 100%);
  background-size: cover;
  background-position: center;
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
  padding: 12px;
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
  color: #fff;
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
  text-shadow: 0 1px 3px rgba(0,0,0,0.8), 0 0 8px rgba(0,0,0,0.5);
}

.card-title:hover {
  color: var(--accent);
}

.card-date {
  font-size: 10px;
  color: #fff;
  opacity: 0.9;
  text-shadow: 0 1px 3px rgba(0,0,0,0.8);
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

/* 主题适配 - 扑克牌背面样式（使用独立图片） */

/* 像素风: 2张不同的像素卡 */
.theme-pixel .card:nth-child(1) .card-back {
  background-image: url('/card-backs/back_1.png');
}
.theme-pixel .card:nth-child(2) .card-back {
  background-image: url('/card-backs/back_4.png');
}
.theme-pixel .card:nth-child(3) .card-back {
  background-image: url('/card-backs/back_1.png');
}

/* 可爱风: 2张可爱卡 */
.theme-cute .card:nth-child(1) .card-back {
  background-image: url('/card-backs/back_2.jpg');
}
.theme-cute .card:nth-child(2) .card-back {
  background-image: url('/card-backs/back_5.png');
}
.theme-cute .card:nth-child(3) .card-back {
  background-image: url('/card-backs/back_2.jpg');
}

/* Q版: 2张Q版卡 */
.theme-qver .card:nth-child(1) .card-back {
  background-image: url('/card-backs/back_3.png');
}
.theme-qver .card:nth-child(2) .card-back {
  background-image: url('/card-backs/back_6.png');
}
.theme-qver .card:nth-child(3) .card-back {
  background-image: url('/card-backs/back_3.png');
}

/* 隐藏扑克牌花色符号（改用背景图） */
.poker-pattern {
  display: none;
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
