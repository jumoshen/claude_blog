<template>
  <div class="home">
    <div class="main-content">
      <div class="post-list">
        <PostCard v-for="post in posts" :key="post.slug" :post="post" />
        <div v-if="posts.length === 0" class="empty">No posts yet.</div>
      </div>
      <aside class="sidebar">
        <TagCloud :tags="tags" />
        <CategoryList :categories="categories" />
      </aside>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '../api'
import PostCard from '../components/post/PostCard.vue'
import TagCloud from '../components/post/TagCloud.vue'
import CategoryList from '../components/post/CategoryList.vue'

const posts = ref([])
const tags = ref({})
const categories = ref({})

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
@media (max-width: 768px) {
  .main-content { flex-direction: column; }
  .sidebar { width: 100%; }
}
</style>
