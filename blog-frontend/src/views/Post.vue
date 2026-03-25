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
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { marked } from 'marked'
import api from '../api'

const route = useRoute()
const post = ref(null)
const content = ref('')

onMounted(async () => {
  const res = await api.getPost(route.params.slug)
  if (res.code === 0) {
    post.value = res.data.post
    content.value = res.data.content
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
</style>
