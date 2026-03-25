<template>
  <div class="post-card" @click="$router.push(`/post/${post.slug}`)">
    <h2 class="title">{{ post.title }}</h2>
    <div class="meta">
      <span class="date">{{ formatDate(post.date) }}</span>
      <span class="views">{{ post.views }} views</span>
    </div>
    <div class="tags">
      <el-tag v-for="tag in post.tags" :key="tag" size="small">{{ tag }}</el-tag>
    </div>
    <p class="summary">{{ post.summary }}</p>
  </div>
</template>

<script setup>
defineProps({
  post: { type: Object, required: true },
})

const formatDate = (date) => {
  const d = new Date(date)
  const pad = (n) => n.toString().padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}
</script>

<style scoped>
.post-card {
  background: #fff;
  border-radius: 8px;
  padding: 20px;
  margin-bottom: 20px;
  cursor: pointer;
  transition: box-shadow 0.2s;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}
.post-card:hover { box-shadow: 0 4px 16px rgba(0,0,0,0.15); }
.title { margin: 0 0 10px; color: #333; }
.meta { font-size: 14px; color: #999; margin-bottom: 10px; }
.meta span { margin-right: 15px; }
.tags { margin-bottom: 10px; }
.tags .el-tag { margin-right: 5px; }
.summary { color: #666; line-height: 1.6; }
</style>
