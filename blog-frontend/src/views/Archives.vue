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
.archives h1 { margin-bottom: 30px; color: var(--text-h); font-size: 32px; }
.year-section { margin-bottom: 40px; }
.year {
  color: var(--accent);
  margin-bottom: 20px;
  font-size: 28px;
  font-weight: 600;
  padding-bottom: 10px;
  border-bottom: 2px solid var(--accent);
  display: inline-block;
}
.post-list { list-style: none; padding: 0; margin: 0; background: var(--bg); border-radius: 12px; overflow: hidden; }
.post-list li {
  display: flex;
  align-items: center;
  padding: 14px 20px;
  border-bottom: 1px solid var(--border);
  cursor: pointer;
  transition: all 0.2s;
}
.post-list li:last-child { border-bottom: none; }
.post-list li:hover { background: var(--accent-bg); padding-left: 28px; }
.post-list .date {
  width: 140px;
  color: var(--text);
  font-size: 13px;
  font-family: var(--mono);
  flex-shrink: 0;
}
.post-list .title {
  flex: 1;
  color: var(--text-h);
  font-weight: 500;
}
.post-list li:hover .title { color: var(--accent); }
</style>
