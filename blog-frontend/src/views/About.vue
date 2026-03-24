<template>
  <div class="about">
    <div class="content" v-if="content" v-html="content"></div>
    <div v-else class="loading">Loading...</div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '../api'

const content = ref('')

onMounted(async () => {
  const res = await api.getAbout()
  if (res.code === 0) content.value = res.data.content
})
</script>

<style scoped>
.about { max-width: 800px; margin: 0 auto; padding: 40px 20px; }
.content { background: #fff; padding: 30px; border-radius: 8px; line-height: 1.8; }
.loading { text-align: center; padding: 40px; color: #999; }
.content :deep(h1), .content :deep(h2) { margin-top: 30px; margin-bottom: 15px; }
.content :deep(img) { max-width: 100%; }
</style>
