import axios from 'axios'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api/v1',
  timeout: 10000,
})

// 请求拦截器
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => Promise.reject(error)
)

// 响应拦截器
api.interceptors.response.use(
  (response) => response.data,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

export default {
  // Site
  getSite: () => api.get('/site'),

  // Posts
  getPosts: (params) => api.get('/posts', { params }),
  getPost: (slug) => api.get(`/posts/${slug}`),
  getNavigation: (slug) => api.get(`/posts/${slug}/navigation`),
  getTOC: (slug) => api.get(`/posts/${slug}/toc`),
  searchPosts: (keyword) => api.get('/posts/search', { params: { q: keyword } }),
  listPopularPosts: (limit = 10) => api.get('/posts/popular', { params: { limit } }),
  listRelatedPosts: (slug, limit = 5) => api.get(`/posts/${slug}/related`, { params: { limit } }),
  listFeaturedPosts: () => api.get('/posts/featured'),
  getArchives: () => api.get('/archives'),
  getAbout: () => api.get('/about'),
  getTags: () => api.get('/tags'),
  getCategories: () => api.get('/categories'),

  // Auth
  getLoginInfo: () => api.get('/auth/login'),
  loginCallback: (code) => api.get('/auth/callback', { params: { code } }),
  logout: () => api.post('/auth/logout'),
  getMe: () => api.get('/auth/me'),

  // Captcha
  getCaptcha: () => api.get('/captcha'),

  // Password Login
  loginWithPassword: (data) => api.post('/auth/login', data),

  // Comments
  getComments: (postSlug) => api.get(`/comments/${postSlug}`),
  createComment: (data) => api.post('/comments', data),

  // Like & Favorite
  likePost: (slug) => api.post(`/posts/${slug}/like`),
  getPostLikes: (slug) => api.get(`/posts/${slug}/likes`),
  getPostFavorite: (slug) => api.get(`/posts/${slug}/favorite`),
  favoritePost: (slug) => api.post(`/posts/${slug}/favorite`),
  getMyFavorites: () => api.get('/users/me/favorites'),

  // Post Password
  checkPostPassword: (slug) => api.get(`/posts/${slug}/check`),
  verifyPostPassword: (slug, password) => api.post(`/posts/${slug}/verify`, { password }),

  // Admin
  refresh: () => api.post('/admin/refresh'),
}
