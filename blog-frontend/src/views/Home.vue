<template>
  <div class="home">
    <div class="main-content">
      <div class="post-list">
        <div v-if="activeTag" class="filter-bar">
          <span>Tag: <strong>{{ activeTag }}</strong></span>
          <el-button text @click="clearTag">清除</el-button>
        </div>
        <PostCard v-for="post in filteredPosts" :key="post.slug" :post="post" />
        <div v-if="filteredPosts.length === 0" class="empty">No posts yet.</div>
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

const route = useRoute()
const router = useRouter()
const posts = ref([])
const tags = ref({})
const categories = ref({})

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
.home { max-width: 1200px; margin: 0 auto; padding: 20px; }
.main-content { display: flex; gap: 20px; }
.post-list { flex: 1; }
.sidebar { width: 280px; }
.empty { text-align: center; padding: 40px; color: #999; }
.filter-bar {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  background: var(--accent-bg);
  border-radius: 8px;
  margin-bottom: 20px;
  color: var(--text);
}
@media (max-width: 768px) {
  .main-content { flex-direction: column; }
  .sidebar { width: 100%; }
}
</style>
