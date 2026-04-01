# Lucky Draw Tool Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 在C端新增一个可收起的工具箱，包含骰子和翻牌两个工具，点击翻牌显示3张扑克牌，点击后3D翻转到正面显示文章标题。

**Architecture:** 将现有的独立骰子按钮封装到Toolbox组件中，新增LuckyDrawModal翻牌弹窗组件，主题样式跟随全局主题配置。

**Tech Stack:** Vue 3 Composition API, CSS 3D Transforms, CSS Variables for theming

---

## File Structure

```
blog-frontend/src/components/common/
├── Toolbox.vue                    # 新增 - 工具箱容器
├── LuckyDrawModal.vue              # 新增 - 翻牌弹窗组件
├── ReadingProgress.vue             # 修改 - 移除独立骰子按钮，保留进度条和返回顶部
└── (existing files...)
```

---

## Task 1: 创建Toolbox.vue工具箱组件

**Files:**
- Create: `blog-frontend/src/components/common/Toolbox.vue`
- Modify: `blog-frontend/src/components/common/ReadingProgress.vue` (移除骰子部分)
- Test: `http://localhost:5173` (手动验证)

- [ ] **Step 1: 创建Toolbox.vue基础结构**

```vue
<template>
  <div class="toolbox" :class="{ expanded: isExpanded }" @mouseenter="expand" @mouseleave="collapse">
    <!-- 收起状态的图标 -->
    <div class="toolbox-trigger">
      <span class="icon">🎰</span>
    </div>

    <!-- 展开的工具列表 -->
    <div class="toolbox-content">
      <div class="tool-item dice-tool" @click="handleDiceClick">
        <div class="dice-mini"></div>
        <span>骰子</span>
      </div>
      <div class="tool-item card-tool" @click="handleCardClick">
        <div class="card-mini"></div>
        <span>翻牌</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const isExpanded = ref(false)

const expand = () => { isExpanded.value = true }
const collapse = () => { isExpanded.value = false }

const handleDiceClick = () => {
  // 触发骰子功能
}

const handleCardClick = () => {
  // 触发翻牌功能
}
</script>

<style scoped>
.toolbox {
  position: fixed;
  bottom: 30px;
  left: 30px;
  z-index: 1000;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
}

.toolbox-trigger {
  width: 48px;
  height: 48px;
  background: var(--card-bg);
  border: 2px solid var(--accent);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.3s;
}

.toolbox-trigger:hover {
  transform: scale(1.1);
}

.toolbox-content {
  display: flex;
  flex-direction: column;
  gap: 12px;
  overflow: hidden;
  max-height: 0;
  transition: max-height 0.3s ease, opacity 0.3s ease;
  opacity: 0;
}

.toolbox.expanded .toolbox-content {
  max-height: 200px;
  opacity: 1;
  margin-top: 12px;
}

.tool-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  background: var(--card-bg);
  border: 2px solid var(--accent);
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s;
}

.tool-item:hover {
  transform: translateX(4px);
  background: var(--accent);
  color: white;
}
</style>
```

- [ ] **Step 2: 在App.vue或Layout中引入Toolbox**

在主布局组件中引入Toolbox组件，替换原有的ReadingProgress中的骰子按钮

- [ ] **Step 3: 测试hover展开**

Run: `npm run dev` → 访问页面 → hover工具箱看是否展开
Expected: hover显示骰子和翻牌选项

- [ ] **Step 4: Commit**

```bash
git add blog-frontend/src/components/common/Toolbox.vue
git commit -m "feat: add Toolbox component with hover expand"
```

---

## Task 2: 创建LuckyDrawModal.vue翻牌弹窗

**Files:**
- Create: `blog-frontend/src/components/common/LuckyDrawModal.vue`
- Modify: `blog-frontend/src/components/common/Toolbox.vue` (导入并使用)

- [ ] **Step 1: 创建LuckyDrawModal基础结构**

```vue
<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="visible" class="modal-overlay" @click.self="close">
        <div class="modal-content">
          <h3>选择一张卡牌</h3>
          <div class="cards-container">
            <div
              v-for="(card, index) in cards"
              :key="card.slug"
              class="card-wrapper"
              :style="{ animationDelay: index * 0.1 + 's' }"
            >
              <div
                class="card"
                :class="{ flipped: card.flipped }"
                @click="flipCard(card)"
              >
                <div class="card-front">
                  <div class="card-back-pattern"></div>
                </div>
                <div class="card-back">
                  <span class="card-title">{{ card.title }}</span>
                </div>
              </div>
            </div>
          </div>
          <button class="close-btn" @click="close">关闭</button>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import api from '../api'

const props = defineProps({
  visible: Boolean
})

const emit = defineEmits(['close'])
const router = useRouter()
const cards = ref([])

const fetchCards = async () => {
  const res = await api.getPosts({ page: 1, page_size: 3 })
  const posts = res.data.list || []
  cards.value = posts.map(post => ({
    ...post,
    flipped: false
  }))
}

const flipCard = (card) => {
  if (!card.flipped) {
    card.flipped = true
    // 翻开后延迟跳转
    setTimeout(() => {
      router.push(`/post/${card.slug}`)
      close()
    }, 800)
  }
}

const close = () => {
  emit('close')
}

watch(() => props.visible, (val) => {
  if (val) {
    fetchCards()
  }
})
</script>

<style scoped>
/* 基础样式 */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
}

.modal-content {
  background: var(--card-bg);
  padding: 32px;
  border-radius: 16px;
  text-align: center;
}

.cards-container {
  display: flex;
  gap: 24px;
  justify-content: center;
  margin: 24px 0;
}

/* 3D翻牌效果 */
.card-wrapper {
  perspective: 1000px;
}

.card {
  width: 140px;
  height: 200px;
  position: relative;
  transform-style: preserve-3d;
  transition: transform 0.8s cubic-bezier(0.4, 0.0, 0.2, 1);
  cursor: pointer;
}

.card.flipped {
  transform: rotateY(180deg);
}

.card-front, .card-back {
  position: absolute;
  width: 100%;
  height: 100%;
  backface-visibility: hidden;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.card-front {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: 3px solid var(--accent);
}

.card-back {
  background: var(--card-bg);
  transform: rotateY(180deg);
  padding: 12px;
}

.card-title {
  font-size: 14px;
  text-align: center;
  overflow: hidden;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
}
</style>
```

- [ ] **Step 2: 在Toolbox中集成LuckyDrawModal**

更新Toolbox.vue，点击翻牌按钮时显示弹窗

- [ ] **Step 3: 测试翻牌动画**

Run: `npm run dev` → hover工具箱 → 点击翻牌 → 应该显示3张牌 → 点击任意一张牌应该翻转
Expected: 牌3D翻转，显示标题后跳转

- [ ] **Step 4: Commit**

```bash
git add blog-frontend/src/components/common/LuckyDrawModal.vue blog-frontend/src/components/common/Toolbox.vue
git commit -m "feat: add LuckyDrawModal with 3D card flip animation"
```

---

## Task 3: 主题适配 - 三种风格扑克牌背面

**Files:**
- Modify: `blog-frontend/src/components/common/LuckyDrawModal.vue`

- [ ] **Step 1: 添加像素风主题扑克牌背面**

在LuckyDrawModal.vue的`<style>`中添加像素风样式：

```css
/* 像素风主题 - 扑克牌背面 */
.theme-pixel .card-front {
  background:
    repeating-linear-gradient(
      0deg,
      transparent,
      transparent 8px,
      rgba(0,0,0,0.1) 8px,
      rgba(0,0,0,0.1) 16px
    ),
    repeating-linear-gradient(
      90deg,
      transparent,
      transparent 8px,
      rgba(0,0,0,0.1) 8px,
      rgba(0,0,0,0.1) 16px
    ),
    linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: 4px solid #000;
  box-shadow:
    4px 4px 0 #000,
    inset 0 0 0 2px var(--accent);
}

.theme-pixel .card-front::before {
  content: '♠';
  font-size: 48px;
  color: var(--accent);
  text-shadow: 2px 2px 0 #000;
}
```

- [ ] **Step 2: 添加可爱风主题扑克牌背面**

```css
/* 可爱风主题 - 扑克牌背面 */
.theme-cute .card-front {
  background: linear-gradient(135deg, #ffecd2 0%, #fcb69f 100%);
  border: none;
  border-radius: 20px;
  box-shadow:
    0 8px 32px rgba(255, 107, 157, 0.3),
    inset 0 2px 4px rgba(255,255,255,0.8);
}

.theme-cute .card-front::before {
  content: '♥';
  font-size: 56px;
  background: linear-gradient(135deg, #ff8ab5 0%, #ff6b9d 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  filter: drop-shadow(2px 2px 0 rgba(255,107,157,0.3));
}
```

- [ ] **Step 3: 添加Q版主题扑克牌背面**

```css
/* Q版主题 - 霓虹扑克牌背面 */
.theme-qver .card-front {
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 100%);
  border: 2px solid rgba(124, 106, 255, 0.5);
  border-radius: 16px;
  box-shadow:
    0 0 20px rgba(124, 106, 255, 0.3),
    inset 0 0 30px rgba(124, 106, 255, 0.1);
  animation: neonPulse 2s ease-in-out infinite;
}

.theme-qver .card-front::before {
  content: '♦';
  font-size: 52px;
  color: #a78bfa;
  text-shadow:
    0 0 10px #a78bfa,
    0 0 20px #a78bfa,
    0 0 40px #7c6aff;
}

@keyframes neonPulse {
  0%, 100% { box-shadow: 0 0 20px rgba(124, 106, 255, 0.3); }
  50% { box-shadow: 0 0 30px rgba(124, 106, 255, 0.5); }
}
```

- [ ] **Step 4: 添加翻牌弹窗主题适配**

更新LuckyDrawModal弹窗整体样式，跟随主题变化

- [ ] **Step 5: 测试三种主题**

Run: 切换三种主题 → 打开翻牌工具 → 验证牌背面风格是否跟随主题变化

- [ ] **Step 6: Commit**

```bash
git add blog-frontend/src/components/common/LuckyDrawModal.vue
git commit -m "feat: add theme-specific card back designs for pixel/cute/qver"
```

---

## Task 4: 修改ReadingProgress.vue移除独立骰子按钮

**Files:**
- Modify: `blog-frontend/src/components/common/ReadingProgress.vue`

- [ ] **Step 1: 移除骰子相关代码**

删除骰子按钮的HTML、JS逻辑和CSS样式，只保留：
- 阅读进度条
- 返回顶部按钮

- [ ] **Step 2: 验证功能完整性**

Run: `npm run dev` → 滚动页面看进度条 → 点击返回顶部 → 工具箱骰子功能正常

- [ ] **Step 3: Commit**

```bash
git commit -m "refactor: remove dice from ReadingProgress, move to Toolbox"
```

---

## Task 5: 整体测试与微调

- [ ] **Step 1: 完整流程测试**

1. hover左下角工具箱 → 应该展开
2. 点击骰子工具 → 骰子动画 → 跳转文章
3. 点击翻牌工具 → 弹出3张牌
4. 点击任意牌 → 3D翻转 → 显示标题 → 跳转

- [ ] **Step 2: 主题切换测试**

切换像素/可爱/Q版三种主题，验证翻牌工具和骰子都正常

- [ ] **Step 3: 移动端测试**

验证移动端hover不可用时的交互（考虑改为点击展开）

---

## 验证清单

- [ ] 工具箱hover展开/收起正常
- [ ] 骰子功能在工具箱中正常工作
- [ ] 翻牌弹窗显示3张牌
- [ ] 牌背面风格跟随主题
- [ ] 3D翻牌动画流畅
- [ ] 翻牌后跳转正确文章
- [ ] ReadingProgress中骰子已移除
- [ ] 三种主题样式完整
