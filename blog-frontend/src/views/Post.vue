<template>
  <div class="post-view" v-if="post">
    <!-- Password Protection -->
    <div v-if="needsPassword" class="password-overlay">
      <div class="password-card">
        <h2>这篇文章受密码保护</h2>
        <p>请输入密码来查看文章内容</p>
        <div class="password-input-wrapper">
          <input
            v-model="passwordInput"
            type="password"
            placeholder="输入密码..."
            class="password-input"
            @keydown.enter="verifyPassword"
          />
          <button class="password-btn" @click="verifyPassword" :disabled="verifying">
            {{ verifying ? '验证中...' : '确认' }}
          </button>
        </div>
        <p v-if="passwordError" class="password-error">{{ passwordError }}</p>
      </div>
    </div>

    <!-- Main Content -->
    <div class="main-wrapper" v-else>
      <article class="post-content">
        <h1 class="title">{{ post.title }}</h1>
        <div class="meta">
          <span>{{ formatDate(post.date) }}</span>
          <span>{{ post.views }} views</span>
          <span v-if="readingTime">约 {{ readingTime }} 分钟读完</span>
          <span v-if="post.is_pinned" class="pin-badge">置顶</span>
          <span v-if="post.is_featured" class="feature-badge">推荐</span>
          <div class="post-actions">
            <button class="action-btn like-btn" :class="{ active: isLiked }" @click="toggleLike" :title="isLiked ? '取消点赞' : '点赞'">
              <svg width="18" height="18" viewBox="0 0 24 24" :fill="isLiked ? 'currentColor' : 'none'" stroke="currentColor" stroke-width="2">
                <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"/>
              </svg>
              <span class="like-count">{{ likeCount }}</span>
            </button>
            <button class="action-btn favorite-btn" :class="{ active: isFavorited }" @click="toggleFavorite" :title="isFavorited ? '取消收藏' : '收藏'">
              <svg width="18" height="18" viewBox="0 0 24 24" :fill="isFavorited ? 'currentColor' : 'none'" stroke="currentColor" stroke-width="2">
                <polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"/>
              </svg>
            </button>
            <button class="action-btn share-btn" @click="showSharePoster = true" title="分享海报">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M4 12v8a2 2 0 002 2h12a2 2 0 002-2v-8M16 6l-4-4-4 4M12 2v13"/>
              </svg>
            </button>
            <button class="action-btn donate-btn" @click="showDonation = true" title="打赏作者">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor">
                <path d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z"/>
              </svg>
            </button>
          </div>
        </div>
        <div class="tags">
          <el-tag v-for="tag in post.tags" :key="tag" size="small">{{ tag }}</el-tag>
        </div>
        <div class="content" v-html="renderedContent" v-code-copy></div>
      </article>

      <!-- 评论区域 -->
      <div class="comment-section">
        <div class="comment-header" @click="showComments = !showComments">
          <h3 class="comment-title">评论 {{ comments.length > 0 ? '(' + comments.length + ')' : '' }}</h3>
          <span class="comment-toggle">{{ showComments ? '收起' : '展开' }}</span>
        </div>

        <transition name="slide">
          <div v-if="showComments">
            <!-- 评论输入框 -->
            <div class="comment-form">
              <div class="form-row">
                <input
                  v-model="commentForm.nickname"
                  type="text"
                  placeholder="昵称（必填）"
                  class="nickname-input"
                  maxlength="50"
                />
              </div>
              <div class="form-row">
                <textarea
                  v-model="commentForm.content"
                  placeholder="说点什么...（支持弹幕显示）"
                  class="comment-input"
                  rows="3"
                  maxlength="500"
                  @keydown.enter.exact.prevent="submitComment"
                ></textarea>
              </div>
              <div class="form-row form-actions">
                <span class="char-count">{{ commentForm.content.length }}/500</span>
                <button
                  class="submit-btn"
                  :disabled="submitting || !canSubmit"
                  @click="submitComment"
                >
                  {{ submitting ? '发送中...' : '发送评论' }}
                </button>
              </div>
            </div>

            <!-- 评论列表 -->
            <div class="comment-list" v-if="comments.length > 0">
              <div
                v-for="comment in comments"
                :key="comment.id"
                class="comment-item"
              >
                <div class="comment-header-row">
                  <span class="comment-nickname">{{ comment.nickname }}</span>
                  <span class="comment-time">{{ formatCommentTime(comment.created_at) }}</span>
                </div>
                <div class="comment-content">{{ comment.content }}</div>
              </div>
            </div>
            <div v-else class="no-comments">
              暂无评论，来抢沙发吧~
            </div>
          </div>
        </transition>
      </div>
    </div>

    <!-- Right Sidebar -->
    <aside class="post-sidebar">
      <!-- TOC -->
      <div v-if="toc.length > 0" class="sidebar-block">
        <div class="block-title">文章目录</div>
        <nav class="toc-nav">
          <a
            v-for="item in toc"
            :key="item.id"
            :href="'#' + item.id"
            class="toc-item"
            :class="{ 'toc-h2': item.level === 2, 'toc-h3': item.level === 3 }"
            @click.prevent="scrollToHeading(item.id)"
          >
            {{ item.text }}
          </a>
        </nav>
      </div>

      <!-- Navigation: Prev/Next -->
      <div v-if="navigation" class="sidebar-block">
        <div class="block-title">导航</div>
        <div class="nav-list">
          <router-link
            v-if="navigation.prev"
            :to="'/post/' + navigation.prev.slug"
            class="nav-item nav-prev"
          >
            <span class="nav-arrow">←</span>
            <span class="nav-text">{{ navigation.prev.title }}</span>
          </router-link>
          <router-link
            v-if="navigation.next"
            :to="'/post/' + navigation.next.slug"
            class="nav-item nav-next"
          >
            <span class="nav-text">{{ navigation.next.title }}</span>
            <span class="nav-arrow">→</span>
          </router-link>
        </div>
      </div>

      <!-- Related Posts -->
      <div v-if="relatedPosts.length > 0" class="sidebar-block">
        <div class="block-title">相关推荐</div>
        <div class="related-list">
          <router-link
            v-for="related in relatedPosts"
            :key="related.slug"
            :to="'/post/' + related.slug"
            class="related-item"
          >
            <span class="related-title">{{ related.title }}</span>
            <span class="related-views">{{ related.views }} 阅读</span>
          </router-link>
        </div>
      </div>
    </aside>

    <!-- Share Poster Modal -->
    <SharePoster
      :visible="showSharePoster"
      :post="post"
      @close="showSharePoster = false"
    />

    <!-- Donation Modal -->
    <DonationModal
      :visible="showDonation"
      :post="post"
      @close="showDonation = false"
    />
  </div>
</template>

<script setup>
import { ref, watch, onMounted, onUnmounted, computed, inject } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { marked } from 'marked'
import api from '../api'
import SharePoster from '../components/post/SharePoster.vue'
import DonationModal from '../components/post/DonationModal.vue'

const route = useRoute()
const router = useRouter()

// Get danmu functions from App.vue via inject
const danmuSettings = inject('danmuSettings')
const addDanmu = inject('addDanmu')
const clearDanmus = inject('clearDanmus')
const danmuLayer = inject('danmuLayer')

const post = ref(null)
const content = ref('')
const comments = ref([])
const submitting = ref(false)
const showSharePoster = ref(false)
const showDonation = ref(false)
const showComments = ref(false)
const navigation = ref(null)
const toc = ref([])
const relatedPosts = ref([])
const isLiked = ref(false)
const isFavorited = ref(false)
const likeCount = ref(0)

// Password protection
const needsPassword = ref(false)
const passwordInput = ref('')
const verifying = ref(false)
const passwordError = ref('')

// Check if post needs password and if already verified
const checkPostPassword = async (slug) => {
  try {
    const res = await api.checkPostPassword(slug)
    if (res.code === 0 && res.data.protected) {
      // Check localStorage for verified status
      const verifiedKey = `post_verified_${slug}`
      const verifiedTime = localStorage.getItem(verifiedKey)
      if (verifiedTime) {
        const elapsed = Date.now() - parseInt(verifiedTime)
        if (elapsed < 24 * 60 * 60 * 1000) {
          // Already verified within 24 hours
          needsPassword.value = false
          return
        }
      }
      needsPassword.value = true
    } else {
      needsPassword.value = false
    }
  } catch (e) {
    console.error('Failed to check password:', e)
    needsPassword.value = false
  }
}

// Verify post password
const verifyPassword = async () => {
  if (!passwordInput.value || verifying.value) return

  verifying.value = true
  passwordError.value = ''

  try {
    const res = await api.verifyPostPassword(route.params.slug, passwordInput.value)
    if (res.code === 0) {
      // Store verification in localStorage
      const verifiedKey = `post_verified_${route.params.slug}`
      localStorage.setItem(verifiedKey, Date.now().toString())
      needsPassword.value = false
      passwordInput.value = ''
    } else {
      passwordError.value = res.message || '密码错误'
    }
  } catch (e) {
    console.error('Failed to verify password:', e)
    passwordError.value = '验证失败，请稍后重试'
  } finally {
    verifying.value = false
  }
}

// 阅读时间计算 (字数/200字每分钟)
const readingTime = computed(() => {
  if (!content.value) return 0
  const words = content.value.replace(/[^\u4e00-\u9fa5a-zA-Z0-9]/g, '').length
  return Math.max(1, Math.ceil(words / 200))
})

// Comment form
const commentForm = ref({
  nickname: '',
  content: '',
  device_id: ''
})

// Generate device ID
const generateDeviceId = () => {
  let deviceId = localStorage.getItem('device_id')
  if (!deviceId) {
    deviceId = 'dev_' + Math.random().toString(36).substr(2, 9) + Date.now().toString(36)
    localStorage.setItem('device_id', deviceId)
  }
  return deviceId
}

// WebSocket connection
let ws = null
let wsReconnectTimer = null

const connectWebSocket = (postSlug) => {
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  const host = window.location.host
  const wsUrl = `${protocol}//${host}/ws/comments/${postSlug}`

  try {
    ws = new WebSocket(wsUrl)

    ws.onopen = () => {
      console.log('WebSocket connected')
      ws.send(JSON.stringify({
        type: 'subscribe',
        post_slug: postSlug
      }))
    }

    ws.onmessage = (event) => {
      try {
        const msg = JSON.parse(event.data)
        if (msg.type === 'new_comment' && msg.data) {
          comments.value.unshift(msg.data)
          // Show as danmu if enabled
          if (danmuSettings?.enabled) {
            addDanmu(msg.data)
          }
        }
      } catch (e) {
        console.error('Failed to parse WebSocket message:', e)
      }
    }

    ws.onerror = (error) => {
      console.error('WebSocket error:', error)
    }

    ws.onclose = () => {
      console.log('WebSocket disconnected')
      wsReconnectTimer = setTimeout(() => {
        if (route.params.slug) {
          connectWebSocket(route.params.slug)
        }
      }, 3000)
    }
  } catch (e) {
    console.error('Failed to create WebSocket:', e)
  }
}

const disconnectWebSocket = () => {
  if (wsReconnectTimer) {
    clearTimeout(wsReconnectTimer)
    wsReconnectTimer = null
  }
  if (ws) {
    ws.close()
    ws = null
  }
}

const fetchPost = async () => {
  const res = await api.getPost(route.params.slug)
  if (res.code === 0) {
    post.value = res.data.post
    content.value = res.data.content
  }
  // Fetch navigation
  try {
    const navRes = await api.getNavigation(route.params.slug)
    if (navRes.code === 0) {
      navigation.value = navRes.data
    }
  } catch (e) {
    console.error('Failed to fetch navigation:', e)
  }
  // Fetch TOC
  try {
    const tocRes = await api.getTOC(route.params.slug)
    if (tocRes.code === 0) {
      toc.value = tocRes.data || []
    }
  } catch (e) {
    console.error('Failed to fetch TOC:', e)
  }
  // Fetch related posts
  try {
    const relatedRes = await api.listRelatedPosts(route.params.slug, 5)
    if (relatedRes.code === 0) {
      relatedPosts.value = relatedRes.data || []
    }
  } catch (e) {
    console.error('Failed to fetch related posts:', e)
  }
  // Fetch like count
  try {
    const likeRes = await api.getPostLikes(route.params.slug)
    if (likeRes.code === 0) {
      likeCount.value = likeRes.data.count
    }
  } catch (e) {
    console.error('Failed to fetch like count:', e)
  }
}

// 滚动到标题
const scrollToHeading = (id) => {
  const el = document.getElementById(id)
  if (el) {
    el.scrollIntoView({ behavior: 'smooth', block: 'start' })
  }
}

const fetchComments = async (postSlug) => {
  try {
    const res = await api.getComments(postSlug)
    if (res.code === 0) {
      comments.value = res.data || []
      // Show historical comments as danmu
      if (danmuSettings?.enabled && comments.value.length > 0) {
        comments.value.forEach((comment, index) => {
          setTimeout(() => {
            addDanmu(comment)
          }, index * 500)
        })
      }
    }
  } catch (e) {
    console.error('Failed to fetch comments:', e)
  }
}

const submitComment = async () => {
  if (!canSubmit.value || submitting.value) return

  submitting.value = true
  try {
    const res = await api.createComment({
      post_slug: route.params.slug,
      nickname: commentForm.value.nickname.trim(),
      content: commentForm.value.content.trim(),
      device_id: generateDeviceId()
    })

    if (res.code === 0) {
      commentForm.value.content = ''
      const newComment = {
        id: res.data.id,
        nickname: res.data.nickname,
        content: res.data.content,
        created_at: res.data.created_at
      }
      comments.value.unshift(newComment)
      // Show as danmu immediately
      if (danmuSettings?.enabled) {
        addDanmu(newComment)
      }
    } else {
      alert(res.message || '评论失败')
    }
  } catch (e) {
    console.error('Failed to submit comment:', e)
    alert('评论失败，请稍后重试')
  } finally {
    submitting.value = false
  }
}

const canSubmit = computed(() => {
  return commentForm.value.nickname.trim().length > 0 &&
         commentForm.value.content.trim().length > 0
})

const toggleLike = async () => {
  try {
    const res = await api.likePost(route.params.slug)
    if (res.code === 0) {
      isLiked.value = res.data.liked
      likeCount.value = res.data.count
    }
  } catch (e) {
    console.error('Failed to toggle like:', e)
    alert('请先登录')
  }
}

const toggleFavorite = async () => {
  try {
    const res = await api.favoritePost(route.params.slug)
    if (res.code === 0) {
      isFavorited.value = res.data.favorited
    }
  } catch (e) {
    console.error('Failed to toggle favorite:', e)
    alert('请先登录')
  }
}

onMounted(async () => {
  await fetchPost()
  const slug = route.params.slug
  if (slug) {
    await checkPostPassword(slug)
    if (!needsPassword.value) {
      await fetchComments(slug)
      connectWebSocket(slug)
    }
  }
  // Listen for keyboard navigation events
  window.addEventListener('keyboard-nav', handleKeyboardNav)
})

onUnmounted(() => {
  disconnectWebSocket()
  window.removeEventListener('keyboard-nav', handleKeyboardNav)
})

watch(() => route.params.slug, async (newSlug) => {
  disconnectWebSocket()
  // Clear danmus when switching posts
  clearDanmus()
  if (newSlug) {
    await fetchPost()
    await checkPostPassword(newSlug)
    if (!needsPassword.value) {
      await fetchComments(newSlug)
      connectWebSocket(newSlug)
    }
  }
})

const handleKeyboardNav = (e) => {
  const { direction } = e.detail
  if (direction === 'prev' && navigation.value?.prev) {
    router.push(`/post/${navigation.value.prev.slug}`)
  } else if (direction === 'next' && navigation.value?.next) {
    router.push(`/post/${navigation.value.next.slug}`)
  }
}

const renderedContent = computed(() => {
  if (!content.value) return ''
  return marked.parse(content.value)
})

const formatDate = (date) => {
  const d = new Date(date)
  const pad = (n) => n.toString().padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}

const formatCommentTime = (date) => {
  const d = new Date(date)
  const now = new Date()
  const diff = now - d

  if (diff < 60000) return '刚刚'
  if (diff < 3600000) return `${Math.floor(diff / 60000)}分钟前`
  if (diff < 86400000) return `${Math.floor(diff / 3600000)}小时前`
  if (diff < 604800000) return `${Math.floor(diff / 86400000)}天前`

  return formatDate(date)
}
</script>

<style scoped>
/* 布局 */
.post-view { display: flex; gap: 30px; max-width: 1200px; margin: 0 auto; padding: 40px 20px; }
.main-wrapper { flex: 1; min-width: 0; }

/* 文章内容 */
.post-content { background: var(--card-bg); padding: 30px; border-radius: 8px; box-shadow: 0 4px 20px var(--shadow); }
.title { margin: 0 0 15px; color: var(--text-h); }
.meta { color: var(--text); opacity: 0.7; font-size: 14px; margin-bottom: 15px; display: flex; align-items: center; gap: 15px; flex-wrap: wrap; }
.meta span { margin-right: 0; }
.post-actions { display: flex; gap: 8px; margin-left: auto; }
.action-btn {
  width: 44px; height: 44px; border: 2px solid var(--border); border-radius: 50%; cursor: pointer;
  display: flex; align-items: center; justify-content: center;
  transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1); background: var(--card-bg);
  position: relative; overflow: hidden;
}
.action-btn::before { content: ''; position: absolute; inset: 0; background: linear-gradient(135deg, transparent 0%, rgba(255,255,255,0.1) 50%, transparent 100%); opacity: 0; transition: opacity 0.3s; }
.action-btn:hover::before { opacity: 1; }
.action-btn:hover { transform: translateY(-4px) scale(1.1); box-shadow: 0 12px 35px var(--shadow); }
.action-btn:active { transform: translateY(-2px) scale(1.05); }
.share-btn { color: var(--accent); border-color: var(--accent-border); }
.share-btn:hover { background: var(--accent); color: #fff; border-color: var(--accent); }
.donate-btn { color: #ff6b9d; border-color: rgba(255, 107, 157, 0.3); }
.donate-btn:hover { background: linear-gradient(135deg, #ff6b9d, #ffa8c5); color: #fff; border-color: #ff6b9d; box-shadow: 0 12px 35px rgba(255, 107, 157, 0.4); }
.like-btn { color: #ff6b6b; border-color: rgba(255, 107, 107, 0.3); display: flex; align-items: center; gap: 4px; }
.like-btn:hover { background: linear-gradient(135deg, #ff6b6b, #ff8e8e); color: #fff; border-color: #ff6b6b; box-shadow: 0 12px 35px rgba(255, 107, 107, 0.4); }
.like-btn.active { background: #ff6b6b; color: #fff; border-color: #ff6b6b; animation: likePulse 0.4s ease; }
@keyframes likePulse { 0% { transform: scale(1); } 50% { transform: scale(1.3); } 100% { transform: scale(1); } }
.like-count { font-size: 12px; font-weight: 600; }
.favorite-btn { color: #ffd700; border-color: rgba(255, 215, 0, 0.3); }
.favorite-btn:hover { background: linear-gradient(135deg, #ffd700, #ffed4a); color: #fff; border-color: #ffd700; box-shadow: 0 12px 35px rgba(255, 215, 0, 0.4); }
.favorite-btn.active { background: #ffd700; color: #fff; border-color: #ffd700; animation: starPulse 0.5s ease; }
@keyframes starPulse { 0% { transform: scale(1) rotate(0deg); } 50% { transform: scale(1.3) rotate(10deg); } 100% { transform: scale(1) rotate(0deg); } }
.tags { margin-bottom: 20px; }
.tags .el-tag { margin-right: 5px; background: var(--accent-bg); color: var(--accent); border-color: var(--accent-border); }
.content { line-height: 1.8; color: var(--text); }
.content :deep(h1), .content :deep(h2), .content :deep(h3) { margin-top: 30px; margin-bottom: 15px; color: var(--text-h); }
.content :deep(h2), .content :deep(h3) { padding-top: 60px; margin-top: -45px; }
.content :deep(h2):first-child, .content :deep(h3):first-child { padding-top: 0; margin-top: 0; }
.content :deep(pre) { background: var(--code-bg); padding: 15px; border-radius: 4px; overflow-x: auto; }
.content :deep(code) { background: var(--code-bg); padding: 2px 5px; border-radius: 3px; font-family: monospace; }
.content :deep(img) { max-width: 100%; }
.content :deep(a) { color: var(--accent); }
.content :deep(blockquote) { border-left: 4px solid var(--accent); margin: 20px 0; padding: 10px 20px; background: var(--accent-bg); border-radius: 0 8px 8px 0; }

/* 置顶/推荐标签 */
.pin-badge, .feature-badge { background: var(--accent); color: #fff; padding: 2px 8px; border-radius: 4px; font-size: 12px; font-weight: 500; }
.feature-badge { background: #10b981; }

/* 右侧边栏 */
.post-sidebar { width: 260px; flex-shrink: 0; display: flex; flex-direction: column; gap: 20px; }
.sidebar-block { background: var(--card-bg); border-radius: 8px; padding: 16px; box-shadow: 0 4px 20px var(--shadow); }
.block-title { font-size: 14px; font-weight: 600; color: var(--text-h); margin-bottom: 12px; padding-bottom: 8px; border-bottom: 1px solid var(--border); }

/* TOC */
.toc-nav { display: flex; flex-direction: column; gap: 4px; max-height: 300px; overflow-y: auto; }
.toc-item { font-size: 13px; color: var(--text); text-decoration: none; padding: 4px 8px; border-radius: 4px; transition: all 0.2s; }
.toc-item:hover { background: var(--accent-bg); color: var(--accent); }
.toc-h2 { padding-left: 8px; }
.toc-h3 { padding-left: 16px; font-size: 12px; }

/* 导航 */
.nav-list { display: flex; flex-direction: column; gap: 8px; }
.nav-item { display: flex; align-items: center; gap: 8px; padding: 10px 12px; background: var(--accent-bg); border-radius: 6px; text-decoration: none; transition: all 0.2s; }
.nav-item:hover { background: var(--accent); }
.nav-item:hover .nav-text, .nav-item:hover .nav-arrow { color: #fff; }
.nav-prev { flex-direction: row; }
.nav-next { flex-direction: row-reverse; }
.nav-text { font-size: 13px; color: var(--text-h); flex: 1; }
.nav-arrow { font-size: 14px; color: var(--accent); }

/* 相关推荐 */
.related-list { display: flex; flex-direction: column; gap: 8px; }
.related-item { display: flex; flex-direction: column; gap: 4px; padding: 10px 12px; background: var(--accent-bg); border-radius: 6px; text-decoration: none; transition: all 0.2s; }
.related-item:hover { background: var(--accent); }
.related-item:hover .related-title, .related-item:hover .related-views { color: #fff; }
.related-title { font-size: 13px; color: var(--text-h); line-height: 1.4; }
.related-views { font-size: 11px; color: var(--text); opacity: 0.7; }

/* 评论区域 */
.comment-section { margin-top: 30px; background: var(--card-bg); padding: 20px; border-radius: 8px; box-shadow: 0 4px 20px var(--shadow); }
.comment-header { display: flex; justify-content: space-between; align-items: center; cursor: pointer; user-select: none; }
.comment-title { margin: 0; font-size: 18px; color: var(--text-h); }
.comment-toggle { font-size: 13px; color: var(--accent); padding: 4px 12px; background: var(--accent-bg); border-radius: 12px; }
.comment-toggle:hover { background: var(--accent); color: #fff; }

/* 折叠动画 */
.slide-enter-active, .slide-leave-active { transition: all 0.3s ease; overflow: hidden; }
.slide-enter-from, .slide-leave-to { opacity: 0; max-height: 0; }
.slide-enter-to, .slide-leave-from { opacity: 1; max-height: 1000px; }

.comment-form { margin-top: 16px; }
.form-row { margin-bottom: 10px; }
.nickname-input { width: 100%; padding: 10px 12px; border: 1px solid var(--border); border-radius: 6px; font-size: 14px; box-sizing: border-box; background: var(--bg); color: var(--text); }
.nickname-input:focus { outline: none; border-color: var(--accent); }
.comment-input { width: 100%; padding: 10px 12px; border: 1px solid var(--border); border-radius: 6px; font-size: 14px; resize: vertical; font-family: inherit; box-sizing: border-box; background: var(--bg); color: var(--text); }
.comment-input:focus { outline: none; border-color: var(--accent); }
.form-actions { display: flex; justify-content: space-between; align-items: center; }
.char-count { font-size: 12px; color: #999; }
.submit-btn { padding: 8px 20px; background: var(--accent); color: white; border: none; border-radius: 6px; font-size: 14px; cursor: pointer; transition: all 0.2s; }
.submit-btn:hover:not(:disabled) { opacity: 0.9; transform: translateY(-1px); }
.submit-btn:disabled { opacity: 0.6; cursor: not-allowed; }
.comment-list { display: flex; flex-direction: column; gap: 15px; margin-top: 16px; }
.comment-item { padding: 12px; background: var(--accent-bg); border-radius: 6px; }
.comment-header-row { display: flex; justify-content: space-between; margin-bottom: 6px; }
.comment-nickname { font-weight: bold; color: var(--accent); font-size: 14px; }
.comment-time { font-size: 12px; color: var(--text); opacity: 0.7; }
.comment-content { font-size: 14px; color: var(--text); line-height: 1.6; word-break: break-all; }
.no-comments { text-align: center; color: var(--text); opacity: 0.6; padding: 20px; }

/* 响应式 */
@media (max-width: 1024px) {
  .post-sidebar { width: 220px; }
}
@media (max-width: 768px) {
  .post-view { flex-direction: column; }
  .post-sidebar { width: 100%; }
}

/* Password Protection */
.password-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: var(--bg);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}
.password-card {
  background: var(--card-bg);
  padding: 40px;
  border-radius: 12px;
  box-shadow: 0 8px 40px var(--shadow);
  text-align: center;
  max-width: 400px;
}
.password-card h2 {
  margin: 0 0 12px;
  color: var(--text-h);
}
.password-card p {
  color: var(--text);
  opacity: 0.7;
  margin: 0 0 20px;
}
.password-input-wrapper {
  display: flex;
  gap: 10px;
}
.password-input {
  flex: 1;
  padding: 12px 16px;
  border: 2px solid var(--border);
  border-radius: 8px;
  font-size: 14px;
  background: var(--bg);
  color: var(--text);
  transition: border-color 0.2s;
}
.password-input:focus {
  outline: none;
  border-color: var(--accent);
}
.password-btn {
  padding: 12px 24px;
  background: var(--accent);
  color: #fff;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
}
.password-btn:hover:not(:disabled) {
  opacity: 0.9;
}
.password-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
.password-error {
  color: #ff6b6b;
  margin-top: 12px;
  font-size: 14px;
}
</style>
