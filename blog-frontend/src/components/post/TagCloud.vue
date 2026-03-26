<template>
  <div class="tag-cloud" :class="themeClass">
    <h3>
      <span class="tag-icon" v-if="styleStore.currentTheme === 'cute'">✦</span>
      <span class="tag-icon" v-else-if="styleStore.currentTheme === 'pixel'">■</span>
      <span class="tag-icon" v-else>◇</span>
      Tags
    </h3>
    <div class="tags">
      <el-tag
        v-for="(count, tag) in tags"
        :key="tag"
        :type="getType(count)"
        class="tag-item"
        :class="{ active: activeTag === tag }"
        @click="handleClick(tag)"
      >
        {{ tag }} ({{ count }})
      </el-tag>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useStyleStore } from '../../store/style'

const props = defineProps({
  tags: { type: Object, default: () => ({}) },
})

const emit = defineEmits(['tag-click'])
const styleStore = useStyleStore()
const activeTag = ref('')

const themeClass = computed(() => `theme-${styleStore.currentTheme}`)

const handleClick = (tag) => {
  activeTag.value = activeTag.value === tag ? '' : tag
  emit('tag-click', activeTag.value)
}

const getType = (count) => {
  if (count >= 5) return 'danger'
  if (count >= 3) return 'warning'
  if (count >= 2) return 'success'
  return 'info'
}
</script>

<style scoped>
.tag-cloud {
  background: var(--card-bg);
  padding: 20px;
  border-radius: 12px;
  margin-bottom: 20px;
  border: 1px solid var(--border);
  transition: all 0.3s ease;
}

.tag-cloud h3 {
  margin: 0 0 15px;
  color: var(--text-h);
  font-size: 16px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.tag-icon {
  color: var(--accent);
  animation: iconPulse 2s infinite ease-in-out;
}

@keyframes iconPulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

.tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.tag-item {
  cursor: pointer;
  transition: all 0.2s cubic-bezier(0.34, 1.56, 0.64, 1);
  border: 1px solid var(--border);
  background: transparent;
  color: var(--text);
}

.tag-item:hover {
  transform: translateY(-2px) scale(1.05);
  border-color: var(--accent);
  color: var(--accent);
}

.tag-item.active {
  background: var(--accent);
  border-color: var(--accent);
  color: white;
}

/* 像素风：方块效果 */
.theme-pixel .tag-item:hover {
  transform: translateY(-2px) scale(1.1);
  box-shadow: 3px 3px 0 var(--accent);
}

/* 可爱风：弹性效果 */
.theme-cute .tag-item:hover {
  transform: translateY(-4px) scale(1.15) rotate(-3deg);
}

/* Q版：发光效果 */
.theme-qver .tag-item:hover {
  box-shadow: 0 0 15px var(--accent-border);
}
</style>
