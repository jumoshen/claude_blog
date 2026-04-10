import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: () => import('../views/Home.vue'),
  },
  {
    path: '/post/:slug',
    name: 'Post',
    component: () => import('../views/Post.vue'),
  },
  {
    path: '/archives',
    name: 'Archives',
    component: () => import('../views/Archives.vue'),
  },
  {
    path: '/about',
    name: 'About',
    component: () => import('../views/About.vue'),
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/Login.vue'),
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('../views/NotFound.vue'),
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    } else {
      return { top: 0, behavior: 'smooth' }
    }
  }
})

// Page transition class management
router.beforeEach((to, from) => {
  // Add transitioning class
  document.body.classList.add('page-transitioning')
})

router.afterEach(() => {
  // Remove transitioning class after transition
  setTimeout(() => {
    document.body.classList.remove('page-transitioning')
  }, 300)
})

export default router
