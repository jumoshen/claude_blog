import axios from 'axios'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1',
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
  getArchives: () => api.get('/archives'),
  getAbout: () => api.get('/about'),
  getTags: () => api.get('/tags'),
  getCategories: () => api.get('/categories'),

  // Auth
  getLoginInfo: () => api.get('/auth/login'),
  loginCallback: (code) => api.post('/auth/callback', null, { params: { code } }),
  logout: () => api.post('/auth/logout'),
  getMe: () => api.get('/auth/me'),

  // Admin
  refresh: () => api.post('/admin/refresh'),
}
