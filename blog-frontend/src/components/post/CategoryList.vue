<template>
  <div class="category-list">
    <h3>Categories</h3>
    <ul>
      <li v-for="(count, cat) in categories" :key="cat">
        <el-link :href="`/?category=${cat}`">{{ cat }} ({{ count }})</el-link>
      </li>
    </ul>
  </div>
</template>

<script setup>
defineProps({
  categories: { type: Object, default: () => ({}) },
})
</script>

<style scoped>
.category-list { background: #fff; padding: 20px; border-radius: 8px; }
.category-list h3 { margin: 0 0 15px; }
.category-list ul { list-style: none; padding: 0; margin: 0; }
.category-list li { padding: 5px 0; }
</style>
