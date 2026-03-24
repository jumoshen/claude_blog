<template>
  <footer class="footer">
    <div class="container">
      <p class="description">{{ siteInfo.description }}</p>
      <p class="beian">
        <a href="https://beian.miit.gov.cn/" target="_blank">{{ siteInfo.beian }}</a>
      </p>
      <p class="author">Powered by {{ siteInfo.author }}</p>
    </div>
  </footer>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '../../api'

const siteInfo = ref({
  author: '',
  beian: '',
  description: '',
})

onMounted(async () => {
  try {
    const res = await api.getSite()
    if (res.code === 0) {
      siteInfo.value = res.data
    }
  } catch (e) {
    // ignore
  }
})
</script>

<style scoped>
.footer {
  background: #fff;
  border-top: 1px solid #eee;
  padding: 20px 0;
  margin-top: 40px;
}
.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
  text-align: center;
}
.description {
  color: #666;
  font-size: 14px;
  margin: 0 0 10px;
}
.beian {
  margin: 0 0 5px;
}
.beian a {
  color: #999;
  font-size: 12px;
  text-decoration: none;
}
.author {
  color: #999;
  font-size: 12px;
  margin: 0;
}
</style>
