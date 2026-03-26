<template>
  <div class="tag-cloud">
    <h3>Tags</h3>
    <div class="tags">
      <el-tag
        v-for="(count, tag) in tags"
        :key="tag"
        :type="getType(count)"
        class="tag-item"
        @click="emit('tag-click', tag)"
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

const emit = defineEmits(['tag-click'])

const getType = (count) => {
  if (count >= 5) return 'danger'
  if (count >= 3) return 'warning'
  if (count >= 2) return 'success'
  return 'info'
}
</script>

<style scoped>
.tag-cloud { background: var(--bg); padding: 20px; border-radius: 12px; margin-bottom: 20px; border: 1px solid var(--border); }
.tag-cloud h3 { margin: 0 0 15px; color: var(--text-h); font-size: 16px; }
.tags { display: flex; flex-wrap: wrap; gap: 8px; }
.tag-item { cursor: pointer; transition: all 0.2s; }
.tag-item:hover { transform: translateY(-2px); }
</style>
