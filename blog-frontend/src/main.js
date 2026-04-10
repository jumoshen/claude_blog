import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'

import App from './App.vue'
import router from './router'
import { codeCopy } from './directives/codeCopy'
import { lazyLoad } from './directives/lazyLoad'

const app = createApp(App)
const pinia = createPinia()

// Register Element Plus icons
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

// Register global directives
app.directive('code-copy', codeCopy)
app.directive('lazy-load', lazyLoad)

app.use(pinia)
app.use(router)
app.use(ElementPlus)

app.mount('#app')
