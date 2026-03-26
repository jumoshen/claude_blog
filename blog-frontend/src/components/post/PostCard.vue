<template>
  <div class="post-card" @click="$router.push(`/post/${post.slug}`)">
    <h2 class="title">{{ post.title }}</h2>
    <div class="meta">
      <span class="date">{{ formatDate(post.date) }}</span>
      <span class="views">{{ post.views }} views</span>
    </div>
    <div class="tags">
      <el-tag
        v-for="tag in post.tags"
        :key="tag"
        size="small"
        class="tag-item"
        @click.stop="$router.push({ path: '/', query: { tag } })"
      >{{ tag }}</el-tag>
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
  background: var(--bg);
  border-radius: 12px;
  padding: 24px;
  margin-bottom: 20px;
  cursor: pointer;
  transition: all 0.2s;
  border: 1px solid var(--border);
}
.post-card:hover {
  border-color: var(--accent);
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(0,0,0,0.08);
}
.title { margin: 0 0 12px; color: var(--text-h); font-size: 20px; font-weight: 600; }
.meta { font-size: 13px; color: var(--text); margin-bottom: 12px; }
.meta span { margin-right: 16px; }
.tags { margin-bottom: 12px; display: flex; flex-wrap: wrap; gap: 6px; }
.tag-item { cursor: pointer; transition: all 0.2s; }
.tag-item:hover { transform: scale(1.05); }
.summary { color: var(--text); line-height: 1.6; margin: 0; }
</style>
