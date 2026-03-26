<template>
  <div class="category-list">
    <h3>Categories</h3>
    <ul>
      <li v-for="(count, cat) in categories" :key="cat" @click="handleClick(cat)">
        <span class="cat-name">{{ cat }}</span>
        <span class="count">({{ count }})</span>
      </li>
    </ul>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'

defineProps({
  categories: { type: Object, default: () => ({}) },
})

const router = useRouter()
const handleClick = (cat) => {
  router.push({ path: '/', query: { category: cat } })
}
</script>

<style scoped>
.category-list { background: var(--bg); padding: 20px; border-radius: 12px; margin-bottom: 20px; border: 1px solid var(--border); }
.category-list h3 { margin: 0 0 15px; color: var(--text-h); font-size: 16px; }
.category-list ul { list-style: none; padding: 0; margin: 0; }
.category-list li {
  padding: 8px 12px;
  margin: 4px 0;
  border-radius: 6px;
  cursor: pointer;
  display: flex;
  justify-content: space-between;
  transition: all 0.2s;
}
.category-list li:hover { background: var(--accent-bg); }
.cat-name { color: var(--text-h); }
.count { color: var(--accent); font-size: 13px; }
</style>
