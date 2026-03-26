import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useStyleStore = defineStore('style', () => {
  const themes = {
    pixel: {
      name: '像素风',
      logo: '/logo-pixel.png',
      colors: {
        accent: '#b366ff',
        accentBg: 'rgba(179, 102, 255, 0.1)',
        accentBorder: 'rgba(179, 102, 255, 0.5)',
        text: '#6b6375',
        textHeading: '#2a1a4a',
        bg: '#fafafa',
        border: '#e5e4e7',
        codeBg: '#f3f0ff',
        headerBg: '#fff',
      }
    },
    cute: {
      name: '可爱风',
      logo: '/logo-cute.png',
      colors: {
        accent: '#ff6b9d',
        accentBg: 'rgba(255, 107, 157, 0.1)',
        accentBorder: 'rgba(255, 107, 157, 0.5)',
        text: '#8b7b8b',
        textHeading: '#4a2a3a',
        bg: '#fff5f8',
        border: '#ffd6e0',
        codeBg: '#fff0f5',
        headerBg: '#fff',
      }
    },
    qver: {
      name: 'Q版',
      logo: '/logo-qver.png',
      colors: {
        accent: '#7c6aff',
        accentBg: 'rgba(124, 106, 255, 0.1)',
        accentBorder: 'rgba(124, 106, 255, 0.5)',
        text: '#6b6375',
        textHeading: '#1a1a3a',
        bg: '#f8f9ff',
        border: '#e0dfff',
        codeBg: '#f4f3ff',
        headerBg: '#fff',
      }
    }
  }

  const currentTheme = ref(localStorage.getItem('blog-theme') || 'pixel')
  const theme = computed(() => themes[currentTheme.value])

  function setTheme(name) {
    if (themes[name]) {
      currentTheme.value = name
      localStorage.setItem('blog-theme', name)
      applyTheme(name)
    }
  }

  function applyTheme(name) {
    const t = themes[name]
    const root = document.documentElement
    root.style.setProperty('--accent', t.colors.accent)
    root.style.setProperty('--accent-bg', t.colors.accentBg)
    root.style.setProperty('--accent-border', t.colors.accentBorder)
    root.style.setProperty('--text', t.colors.text)
    root.style.setProperty('--text-h', t.colors.textHeading)
    root.style.setProperty('--bg', t.colors.bg)
    root.style.setProperty('--border', t.colors.border)
    root.style.setProperty('--code-bg', t.colors.codeBg)
    root.style.setProperty('--header-bg', t.colors.headerBg)
    document.body.style.background = t.colors.bg
  }

  return { themes, currentTheme, theme, setTheme, applyTheme }
})
