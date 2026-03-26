<template>
  <div class="about">
    <div class="content" v-if="content" v-html="content"></div>
    <div v-else-if="!loading" class="empty">
      <p>About 页面暂无内容</p>
      <p class="tip">请在后台添加 About 页面内容</p>
    </div>
    <div v-else class="loading">Loading...</div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '../api'

const content = ref('')
const loading = ref(true)

onMounted(async () => {
  try {
    const res = await api.getAbout()
    if (res.code === 0 && res.data?.content) {
      content.value = res.data.content
    }
  } catch (e) {
    console.error('Failed to load about:', e)
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.about { max-width: 800px; margin: 0 auto; padding: 40px 20px; }
.content { background: var(--bg); padding: 30px; border-radius: 12px; line-height: 1.8; border: 1px solid var(--border); }
.loading, .empty { text-align: center; padding: 60px 40px; color: var(--text); background: var(--bg); border-radius: 12px; border: 1px solid var(--border); }
.empty p { margin: 10px 0; }
.empty .tip { color: var(--accent); font-size: 14px; }
.content :deep(h1), .content :deep(h2) { margin-top: 30px; margin-bottom: 15px; color: var(--text-h); }
.content :deep(img) { max-width: 100%; }
</style>
