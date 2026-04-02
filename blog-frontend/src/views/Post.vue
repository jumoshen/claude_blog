<template>
  <div class="post-view" v-if="post">
    <article class="post-content">
      <h1 class="title">{{ post.title }}</h1>
      <div class="meta">
        <span>{{ formatDate(post.date) }}</span>
        <span>{{ post.views }} views</span>
      </div>
      <div class="tags">
        <el-tag v-for="tag in post.tags" :key="tag" size="small">{{ tag }}</el-tag>
      </div>
      <div class="content" v-html="renderedContent"></div>
    </article>

    <!-- 评论区域 -->
    <div class="comment-section">
      <h3 class="comment-title">评论</h3>

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
          <div class="comment-header">
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
  </div>
</template>

<script setup>
import { ref, watch, onMounted, onUnmounted, computed, inject } from 'vue'
import { useRoute } from 'vue-router'
import { marked } from 'marked'
import api from '../api'

const route = useRoute()

// Get danmu functions from App.vue via inject
const danmuSettings = inject('danmuSettings')
const addDanmu = inject('addDanmu')
const clearDanmus = inject('clearDanmus')
const danmuLayer = inject('danmuLayer')

const post = ref(null)
const content = ref('')
const comments = ref([])
const submitting = ref(false)

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

onMounted(async () => {
  await fetchPost()
  const slug = route.params.slug
  if (slug) {
    await fetchComments(slug)
    connectWebSocket(slug)
  }
})

onUnmounted(() => {
  disconnectWebSocket()
})

watch(() => route.params.slug, async (newSlug) => {
  disconnectWebSocket()
  // Clear danmus when switching posts
  clearDanmus()
  if (newSlug) {
    await fetchPost()
    await fetchComments(newSlug)
    connectWebSocket(newSlug)
  }
})

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
.post-view { max-width: 800px; margin: 0 auto; padding: 40px 20px; }
.post-content { background: #fff; padding: 30px; border-radius: 8px; }
.title { margin: 0 0 15px; }
.meta { color: #999; font-size: 14px; margin-bottom: 15px; }
.meta span { margin-right: 15px; }
.tags { margin-bottom: 20px; }
.tags .el-tag { margin-right: 5px; }
.content { line-height: 1.8; color: #333; }
.content :deep(h1), .content :deep(h2), .content :deep(h3) { margin-top: 30px; margin-bottom: 15px; }
.content :deep(pre) { background: #f6f8fa; padding: 15px; border-radius: 4px; overflow-x: auto; }
.content :deep(code) { background: #f6f8fa; padding: 2px 5px; border-radius: 3px; font-family: monospace; }
.content :deep(img) { max-width: 100%; }

/* 评论区域 */
.comment-section {
  margin-top: 30px;
  background: #fff;
  padding: 20px;
  border-radius: 8px;
}

.comment-title {
  margin: 0 0 20px;
  font-size: 18px;
  color: #333;
}

.comment-form {
  margin-bottom: 20px;
}

.form-row {
  margin-bottom: 10px;
}

.nickname-input {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 14px;
  box-sizing: border-box;
}

.nickname-input:focus {
  outline: none;
  border-color: var(--accent, #646cff);
}

.comment-input {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 14px;
  resize: vertical;
  font-family: inherit;
  box-sizing: border-box;
}

.comment-input:focus {
  outline: none;
  border-color: var(--accent, #646cff);
}

.form-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.char-count {
  font-size: 12px;
  color: #999;
}

.submit-btn {
  padding: 8px 20px;
  background: var(--accent, #646cff);
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
}

.submit-btn:hover:not(:disabled) {
  opacity: 0.9;
  transform: translateY(-1px);
}

.submit-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.comment-list {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.comment-item {
  padding: 12px;
  background: #f9f9f9;
  border-radius: 6px;
}

.comment-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 6px;
}

.comment-nickname {
  font-weight: bold;
  color: var(--accent, #646cff);
  font-size: 14px;
}

.comment-time {
  font-size: 12px;
  color: #999;
}

.comment-content {
  font-size: 14px;
  color: #333;
  line-height: 1.6;
  word-break: break-all;
}

.no-comments {
  text-align: center;
  color: #999;
  padding: 20px;
}
</style>
