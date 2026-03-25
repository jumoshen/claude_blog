<template>
  <div class="archives">
    <h1>Archives</h1>
    <div v-for="(posts, year) in archives" :key="year" class="year-section">
      <h2 class="year">{{ year }}</h2>
      <ul class="post-list">
        <li v-for="post in posts" :key="post.slug" @click="$router.push(`/post/${post.slug}`)">
          <span class="date">{{ formatDate(post.date) }}</span>
          <span class="title">{{ post.title }}</span>
        </li>
      </ul>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '../api'

const archives = ref({})

onMounted(async () => {
  const res = await api.getArchives()
  if (res.code === 0) archives.value = res.data
})

const formatDate = (date) => {
  const d = new Date(date)
  const pad = (n) => n.toString().padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}
</script>

<style scoped>
.archives { max-width: 800px; margin: 0 auto; padding: 40px 20px; }
.archives h1 { margin-bottom: 30px; }
.year-section { margin-bottom: 30px; }
.year { color: #409eff; margin-bottom: 15px; }
.post-list { list-style: none; padding: 0; margin: 0; }
.post-list li {
  display: flex;
  padding: 10px 0;
  border-bottom: 1px solid #eee;
  cursor: pointer;
}
.post-list li:hover { color: #409eff; }
.post-list .date { width: 80px; color: #999; font-size: 14px; }
.post-list .title { flex: 1; }
</style>
