<template>
  <div class="tag-cloud">
    <h3>Tags</h3>
    <div class="tags">
      <el-tag
        v-for="(count, tag) in tags"
        :key="tag"
        :type="getType(count)"
        class="tag-item"
      >
        {{ tag }} ({{ count }})
      </el-tag>
    </div>
  </div>
</template>

<script setup>
defineProps({
  tags: { type: Object, default: () => ({}) },
})

const getType = (count) => {
  if (count >= 5) return 'danger'
  if (count >= 3) return 'warning'
  if (count >= 2) return 'success'
  return 'info'
}
</script>

<style scoped>
.tag-cloud { background: #fff; padding: 20px; border-radius: 8px; margin-bottom: 20px; }
.tag-cloud h3 { margin: 0 0 15px; }
.tags { display: flex; flex-wrap: wrap; gap: 8px; }
.tag-item { cursor: pointer; }
</style>
