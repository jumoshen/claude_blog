<template>
  <div class="danmu-control">
    <div
      class="tool-btn danmu-btn"
      :class="['theme-' + styleStore.currentTheme, { active: danmuEnabled }]"
      @click="toggleDanmu"
      title="弹幕"
    >
      <span class="danmu-icon">💬</span>
      <span class="tool-label">{{ danmuEnabled ? '弹幕ON' : '弹幕OFF' }}</span>
    </div>

    <!-- Danmu settings panel -->
    <transition name="slide">
      <div v-if="showSettings" class="danmu-settings">
        <div class="setting-row">
          <label>密度</label>
          <input
            type="range"
            v-model.number="localDensity"
            min="1"
            max="10"
            step="1"
            @change="updateDensity"
          />
          <span class="setting-value">{{ localDensity }}</span>
        </div>
        <div class="setting-row">
          <label>字体</label>
          <input
            type="range"
            v-model.number="localFontSize"
            min="12"
            max="24"
            step="2"
            @change="updateFontSize"
          />
          <span class="setting-value">{{ localFontSize }}px</span>
        </div>
        <div class="setting-row">
          <label>颜色</label>
          <input
            type="color"
            v-model="localFontColor"
            @change="updateFontColor"
            class="color-picker"
          />
          <span class="setting-value">{{ localFontColor }}</span>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue'
import { useStyleStore } from '../../store/style'

const styleStore = useStyleStore()

const emit = defineEmits(['update:enabled', 'update:density', 'update:fontSize', 'update:fontColor'])

// Danmu state
const danmuEnabled = ref(true)
const showSettings = ref(false)
const localDensity = ref(6)
const localFontSize = ref(16)
const localFontColor = ref('#ffffff')

// Load settings from localStorage
onMounted(() => {
  const saved = localStorage.getItem('danmu-settings')
  if (saved) {
    try {
      const settings = JSON.parse(saved)
      danmuEnabled.value = settings.enabled ?? true
      localDensity.value = settings.density ?? 6
      localFontSize.value = settings.fontSize ?? 16
      localFontColor.value = settings.fontColor ?? '#ffffff'
    } catch (e) {
      // ignore
    }
  }
})

const saveSettings = () => {
  localStorage.setItem('danmu-settings', JSON.stringify({
    enabled: danmuEnabled.value,
    density: localDensity.value,
    fontSize: localFontSize.value,
    fontColor: localFontColor.value
  }))
}

const toggleDanmu = () => {
  danmuEnabled.value = !danmuEnabled.value
  showSettings.value = !showSettings.value || !danmuEnabled.value
  emit('update:enabled', danmuEnabled.value)
  saveSettings()
}

const updateDensity = () => {
  emit('update:density', localDensity.value)
  saveSettings()
}

const updateFontSize = () => {
  emit('update:fontSize', localFontSize.value)
  saveSettings()
}

const updateFontColor = () => {
  emit('update:fontColor', localFontColor.value)
  saveSettings()
}

// Watch for external changes
watch(danmuEnabled, (val) => {
  emit('update:enabled', val)
})

watch(localDensity, (val) => {
  emit('update:density', val)
})

watch(localFontSize, (val) => {
  emit('update:fontSize', val)
})

watch(localFontColor, (val) => {
  emit('update:fontColor', val)
})
</script>

<style scoped>
.danmu-control {
  position: relative;
}

.tool-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  padding: 8px;
  background: transparent;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.tool-btn:hover {
  background: var(--accent-bg);
}

.tool-btn.active {
  background: var(--accent-bg);
}

.danmu-icon {
  font-size: 24px;
}

.tool-label {
  font-size: 10px;
  color: var(--text);
  white-space: nowrap;
}

.danmu-settings {
  position: absolute;
  bottom: 100%;
  left: 0;
  margin-bottom: 10px;
  padding: 12px;
  background: var(--card-bg);
  border: 2px solid var(--accent);
  border-radius: 12px;
  box-shadow: 0 4px 20px var(--shadow);
  min-width: 180px;
  z-index: 1001;
}

.setting-row {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.setting-row:last-child {
  margin-bottom: 0;
}

.setting-row label {
  font-size: 12px;
  color: var(--text);
  min-width: 32px;
}

.setting-row input[type="range"] {
  flex: 1;
  height: 4px;
  -webkit-appearance: none;
  background: var(--accent-bg);
  border-radius: 2px;
  outline: none;
}

.setting-row input[type="range"]::-webkit-slider-thumb {
  -webkit-appearance: none;
  width: 14px;
  height: 14px;
  background: var(--accent);
  border-radius: 50%;
  cursor: pointer;
}

.setting-row input[type="color"] {
  width: 24px;
  height: 24px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  padding: 0;
  background: transparent;
}

.setting-value {
  font-size: 11px;
  color: var(--text-secondary);
  min-width: 50px;
  text-align: right;
}

.color-picker {
  border: 1px solid var(--accent);
}

/* Slide transition */
.slide-enter-active,
.slide-leave-active {
  transition: all 0.2s ease;
}

.slide-enter-from,
.slide-leave-to {
  opacity: 0;
  transform: translateY(10px);
}
</style>
