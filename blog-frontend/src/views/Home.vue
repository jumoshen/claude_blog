<template>
  <div class="home" :class="'theme-' + styleStore.currentTheme">
    <div class="main-content">
      <div class="post-list">
        <transition name="fade" appear>
          <div v-if="activeTag" class="filter-bar">
            <span class="filter-icon">🏷️</span>
            <span>Tag: <strong>{{ activeTag }}</strong></span>
            <el-button text @click="clearTag">清除</el-button>
          </div>
        </transition>

        <transition-group name="list" tag="div" class="posts-wrapper">
          <PostCard
            v-for="(post, index) in filteredPosts"
            :key="post.slug"
            :post="post"
            :style="{ '--delay': index * 0.1 + 's' }"
            class="post-item"
          />
        </transition-group>

        <div v-if="filteredPosts.length === 0" class="empty">
          <span class="empty-icon">📭</span>
          <p>暂无文章</p>
        </div>
      </div>
      <aside class="sidebar">
        <TagCloud :tags="tags" @tag-click="handleTagClick" />
        <CategoryList :categories="categories" />
      </aside>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import api from '../api'
import PostCard from '../components/post/PostCard.vue'
import TagCloud from '../components/post/TagCloud.vue'
import CategoryList from '../components/post/CategoryList.vue'
import { useStyleStore } from '../store/style'

const route = useRoute()
const router = useRouter()
const posts = ref([])
const tags = ref({})
const categories = ref({})
const styleStore = useStyleStore()

const activeTag = computed(() => route.query.tag || '')

const filteredPosts = computed(() => {
  if (!activeTag.value) return posts.value
  return posts.value.filter(p => p.tags && p.tags.includes(activeTag.value))
})

const handleTagClick = (tag) => {
  router.push({ path: '/', query: { tag } })
}

const clearTag = () => {
  router.push({ path: '/', query: {} })
}

onMounted(async () => {
  const [postsRes, tagsRes, categoriesRes] = await Promise.all([
    api.getPosts(),
    api.getTags(),
    api.getCategories(),
  ])
  if (postsRes.code === 0) posts.value = postsRes.data
  if (tagsRes.code === 0) tags.value = tagsRes.data
  if (categoriesRes.code === 0) categories.value = categoriesRes.data
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

@media (max-width: 768px) {
  .main-content {
    flex-direction: column;
  }
  .sidebar {
    width: 100%;
  }
}
</style>
