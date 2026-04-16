<template>
  <div class="home" :class="'theme-' + styleStore.currentTheme">
    <!-- 粒子背景 -->
    <ParticlesBackground :particle-count="40" />
    <div class="main-content">
      <div class="post-list">
        <transition name="fade" appear>
          <div v-if="activeTag" class="filter-bar">
            <span class="filter-icon">🏷️</span>
            <span>Tag: <strong>{{ activeTag }}</strong></span>
            <el-button text @click="clearTag">清除</el-button>
          </div>
        </transition>

        <transition name="fade" appear>
          <div v-if="activeCategory" class="filter-bar">
            <span class="filter-icon">📂</span>
            <span>分类: <strong>{{ activeCategory }}</strong></span>
            <el-button text @click="clearCategory">清除</el-button>
          </div>
        </transition>

        <!-- 搜索栏 -->
        <div class="search-bar">
          <input
            v-model="searchKeyword"
            type="text"
            placeholder="搜索文章标题或内容..."
            class="search-input"
            @keyup.enter="handleSearch"
          />
          <button class="search-btn" @click="handleSearch">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="11" cy="11" r="8"></circle>
              <path d="m21 21-4.35-4.35"></path>
            </svg>
          </button>
        </div>

        <!-- 置顶文章 -->
        <div v-if="featuredPosts.length > 0 && !searchKeyword && !activeTag && !activeCategory" class="featured-section">
          <h3 class="featured-title">推荐文章</h3>
          <div class="featured-list">
            <router-link
              v-for="post in featuredPosts"
              :key="post.slug"
              :to="'/post/' + post.slug"
              class="featured-item"
            >
              <span class="featured-item-title">{{ post.title }}</span>
              <span class="featured-item-views">{{ post.views }} 阅读</span>
            </router-link>
          </div>
        </div>

        <transition-group name="list" tag="div" class="posts-wrapper">
          <PostCard
            v-for="(post, index) in posts"
            :key="post.slug"
            :post="post"
            :style="{ '--delay': index * 0.05 + 's' }"
            class="post-item"
          />
        </transition-group>

        <!-- 加载更多触发器 -->
        <div ref="loadMoreTrigger" class="load-more-trigger">
          <div v-if="loading && posts.length > 0" class="loading-more">
            <div class="loading-dots">
              <span></span><span></span><span></span>
            </div>
            <p>加载更多文章...</p>
          </div>
        </div>

        <div v-if="loading && (!posts || posts.length === 0)" class="loading">
          <div class="loading-spinner"></div>
          <p>加载中...</p>
        </div>

        <div v-else-if="!posts || posts.length === 0" class="empty">
          <span class="empty-icon">📭</span>
          <p>暂无文章</p>
        </div>

        <div v-else-if="!hasMore" class="no-more">
          <span>— 已加载全部 —</span>
        </div>
      </div>
      <aside class="sidebar">
        <TagCloud :tags="tags" @tag-click="handleTagClick" />
        <CategoryList :categories="categories" :activeCategory="activeCategory" />
      </aside>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import api from '../api'
import PostCard from '../components/post/PostCard.vue'
import TagCloud from '../components/post/TagCloud.vue'
import CategoryList from '../components/post/CategoryList.vue'
import ParticlesBackground from '../components/common/ParticlesBackground.vue'
import { useStyleStore } from '../store/style'

const route = useRoute()
const router = useRouter()
const posts = ref([])
const tags = ref({})
const categories = ref({})
const styleStore = useStyleStore()
const loading = ref(false)
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const hasMore = computed(() => (posts.value?.length || 0) < total.value)
const loadMoreTrigger = ref(null)
let observer = null
const searchKeyword = ref('')
const featuredPosts = ref([])
const searchResults = ref([])

const activeTag = computed(() => route.query.tag || '')
const activeCategory = computed(() => route.query.category || '')

const fetchPosts = async (reset = false) => {
  if (loading.value) return
  if (!reset && !hasMore.value) return

  loading.value = true
  try {
    const res = await api.getPosts({
      page: reset ? 1 : page.value,
      page_size: pageSize.value,
      tag: activeTag.value || undefined,
      category: activeCategory.value || undefined
    })
    if (res.code === 0 && res.data) {
      // 兼容新旧两种 API 格式
      const list = res.data.list || (Array.isArray(res.data) ? res.data : [])
      const totalCount = res.data.total ?? (Array.isArray(res.data) ? res.data.length : 0)
      if (reset) {
        posts.value = list
        page.value = 2
      } else {
        posts.value = [...posts.value, ...list]
        page.value++
      }
      total.value = totalCount
    }
  } catch (e) {
    console.error('Failed to fetch posts:', e)
  } finally {
    loading.value = false
  }
}

const handleTagClick = (tag) => {
  router.push({ path: '/', query: { tag } })
}

const clearTag = () => {
  router.push({ path: '/', query: {} })
}

const clearCategory = () => {
  const query = { ...route.query }
  delete query.category
  router.push({ path: '/', query })
}

// 搜索
const handleSearch = async () => {
  if (!searchKeyword.value.trim()) return
  loading.value = true
  try {
    const res = await api.searchPosts(searchKeyword.value.trim())
    if (res.code === 0) {
      // 后端返回 { list, total } 格式
      const data = res.data || {}
      posts.value = data.list || data || []
      total.value = data.total || posts.value.length
    }
  } catch (e) {
    console.error('Failed to search posts:', e)
  } finally {
    loading.value = false
  }
}

const clearSearch = () => {
  searchKeyword.value = ''
  fetchPosts(true)
}

// 获取置顶文章
const fetchFeaturedPosts = async () => {
  try {
    const res = await api.listFeaturedPosts()
    if (res.code === 0) {
      featuredPosts.value = res.data || []
    }
  } catch (e) {
    console.error('Failed to fetch featured posts:', e)
  }
}

// 监听 tag 和 category 变化，重新加载
let initialized = false
watch(activeTag, () => {
  if (initialized) {
    fetchPosts(true)
  }
})

watch(activeCategory, () => {
  if (initialized) {
    fetchPosts(true)
  }
})

// 设置滚动加载
onMounted(async () => {
  const [tagsRes, categoriesRes] = await Promise.all([
    api.getTags(),
    api.getCategories(),
  ])
  if (tagsRes.code === 0) tags.value = tagsRes.data
  if (categoriesRes.code === 0) categories.value = categoriesRes.data

  await fetchPosts(true)
  await fetchFeaturedPosts()
  initialized = true

  // 滚动加载 - 提前 300px 触发
  observer = new IntersectionObserver((entries) => {
    if (entries[0].isIntersecting && !loading.value && hasMore.value) {
      fetchPosts(false)
    }
  }, { threshold: 0, rootMargin: '300px 0px 0px 0px' })

  if (loadMoreTrigger.value) {
    observer.observe(loadMoreTrigger.value)
  }
})

onUnmounted(() => {
  if (observer) {
    observer.disconnect()
  }
})
</script>

<style scoped>
.home {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
  position: relative;
  z-index: 1;
}

.main-content {
  display: flex;
  gap: 24px;
}

.post-list {
  flex: 1;
  min-width: 0;
}

.sidebar {
  width: 280px;
  flex-shrink: 0;
}

.posts-wrapper {
  position: relative;
}

.post-item {
  animation: slideUp 0.6s ease-out forwards;
  animation-delay: var(--delay);
  opacity: 0;
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.list-enter-active {
  transition: all 0.5s ease-out;
}

.list-leave-active {
  transition: all 0.3s ease-in;
}

.list-enter-from {
  opacity: 0;
  transform: translateX(-30px);
}

.list-leave-to {
  opacity: 0;
  transform: translateX(30px);
}

.filter-bar {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  background: var(--accent-bg);
  border-radius: 12px;
  margin-bottom: 20px;
  color: var(--text);
  border: 1px solid var(--accent-border);
}

.filter-icon {
  font-size: 16px;
}

/* 搜索栏 */
.search-bar {
  display: flex;
  gap: 8px;
  margin-bottom: 20px;
}
.search-input {
  flex: 1;
  padding: 10px 14px;
  border: 1px solid var(--border);
  border-radius: 8px;
  font-size: 14px;
  background: var(--card-bg);
  color: var(--text);
}
.search-input:focus {
  outline: none;
  border-color: var(--accent);
}
.search-btn {
  padding: 10px 14px;
  background: var(--accent);
  color: #fff;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}
.search-btn:hover {
  opacity: 0.9;
  transform: translateY(-1px);
}

/* 置顶推荐 */
.featured-section {
  margin-bottom: 20px;
  padding: 16px;
  background: var(--accent-bg);
  border-radius: 12px;
  border: 1px solid var(--accent-border);
}
.featured-title {
  margin: 0 0 12px;
  font-size: 14px;
  font-weight: 600;
  color: var(--text-h);
}
.featured-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.featured-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 12px;
  background: var(--card-bg);
  border-radius: 6px;
  text-decoration: none;
  transition: all 0.2s;
}
.featured-item:hover {
  background: var(--accent);
}
.featured-item:hover .featured-item-title,
.featured-item:hover .featured-item-views {
  color: #fff;
}
.featured-item-title {
  font-size: 14px;
  color: var(--text-h);
}
.featured-item-views {
  font-size: 12px;
  color: var(--text);
  opacity: 0.7;
}

.empty {
  text-align: center;
  padding: 60px 20px;
  color: var(--text);
}

.empty-icon {
  font-size: 48px;
  display: block;
  margin-bottom: 16px;
  animation: float 3s ease-in-out infinite;
}

@keyframes float {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-10px); }
}

.fade-enter-active, .fade-leave-active {
  transition: opacity 0.3s ease;
}
.fade-enter-from, .fade-leave-to {
  opacity: 0;
}

.load-more-trigger {
  height: 60px;
  margin: 20px 0;
}

.loading {
  text-align: center;
  padding: 60px 20px;
  color: var(--text);
}

.loading-spinner {
  width: 40px;
  height: 40px;
  margin: 0 auto 16px;
  border: 3px solid var(--border);
  border-top-color: var(--accent);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

.loading-more {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  padding: 20px;
  animation: fadeIn 0.3s ease;
}

.loading-dots {
  display: flex;
  gap: 6px;
}

.loading-dots span {
  width: 8px;
  height: 8px;
  background: var(--accent);
  border-radius: 50%;
  animation: bounce 1.4s ease-in-out infinite both;
}

.loading-dots span:nth-child(1) { animation-delay: -0.32s; }
.loading-dots span:nth-child(2) { animation-delay: -0.16s; }

@keyframes spin {
  to { transform: rotate(360deg); }
}

@keyframes bounce {
  0%, 80%, 100% {
    transform: scale(0);
  }
  40% {
    transform: scale(1);
  }
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

.no-more {
  text-align: center;
  padding: 20px;
  color: var(--text);
  opacity: 0.6;
  font-size: 14px;
}

@media (max-width: 768px) {
  .main-content {
    flex-direction: column;
  }
  .sidebar {
    width: 100%;
  }
}
</style>
