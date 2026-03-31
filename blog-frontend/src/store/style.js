import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useStyleStore = defineStore('style', () => {
  const themes = {
    pixel: {
      name: '像素风',
      tagline: '8-BIT RETRO',
      logo: '/logo-pixel.png',
      icon: '🎮',
      preview: 'linear-gradient(135deg, #b366ff 0%, #7c3aed 100%)',
      colors: {
        accent: '#b366ff',
        accentRgb: '179, 102, 255',
        accentBg: 'rgba(179, 102, 255, 0.1)',
        accentBorder: 'rgba(179, 102, 255, 0.5)',
        text: '#6b6375',
        textHeading: '#2a1a4a',
        bg: '#fafafa',
        border: '#e5e4e7',
        codeBg: '#f3f0ff',
        headerBg: '#fff',
        cardBg: '#fff',
        shadow: 'rgba(179, 102, 255, 0.2)',
      },
      animation: {
        cardHover: 'pixelHover',
        floatElements: 'pixelFloat',
        loading: 'pixelBlink',
        cursor: 'pixelCursor',
      },
      cardStyle: 'pixelated',
      fontDisplay: 'Press Start 2P, monospace',
      fontBody: '"Noto Sans SC", "Microsoft YaHei", sans-serif',
    },
    cute: {
      name: '可爱风',
      tagline: 'SWEETY DAYS',
      logo: '/logo-cute.png',
      icon: '🌸',
      preview: 'linear-gradient(135deg, #ff6b9d 0%, #ffa8c5 100%)',
      colors: {
        accent: '#ff6b9d',
        accentRgb: '255, 107, 157',
        accentBg: 'rgba(255, 107, 157, 0.1)',
        accentBorder: 'rgba(255, 107, 157, 0.5)',
        text: '#8b7b8b',
        textHeading: '#4a2a3a',
        bg: '#fff5f8',
        border: '#ffd6e0',
        codeBg: '#fff0f5',
        headerBg: '#fff',
        cardBg: '#fff',
        shadow: 'rgba(255, 107, 157, 0.15)',
      },
      animation: {
        cardHover: 'cuteBounce',
        floatElements: 'cuteFloat',
        loading: 'cutePulse',
        cursor: 'cuteSparkle',
      },
      cardStyle: 'rounded',
      fontDisplay: 'Fredoka, cursive',
      fontBody: 'Nunito, sans-serif',
    },
    qver: {
      name: 'Q版',
      tagline: 'CHILL VIBES',
      logo: '/logo-qver.png',
      icon: '⚡',
      preview: 'linear-gradient(135deg, #7c6aff 0%, #a78bfa 100%)',
      colors: {
        accent: '#7c6aff',
        accentRgb: '124, 106, 255',
        accentBg: 'rgba(124, 106, 255, 0.1)',
        accentBorder: 'rgba(124, 106, 255, 0.5)',
        text: '#6b6375',
        textHeading: '#1a1a3a',
        bg: '#f8f9ff',
        border: '#e0dfff',
        codeBg: '#f4f3ff',
        headerBg: '#fff',
        cardBg: 'rgba(255, 255, 255, 0.8)',
        shadow: 'rgba(124, 106, 255, 0.1)',
      },
      animation: {
        cardHover: 'qverSlide',
        floatElements: 'qverFloat',
        loading: 'qverGlow',
        cursor: 'qverRipple',
      },
      cardStyle: 'soft',
      fontDisplay: 'Quicksand, sans-serif',
      fontBody: 'Quicksand, sans-serif',
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
    root.style.setProperty('--accent-rgb', t.colors.accentRgb)
    root.style.setProperty('--accent-bg', t.colors.accentBg)
    root.style.setProperty('--accent-border', t.colors.accentBorder)
    root.style.setProperty('--text', t.colors.text)
    root.style.setProperty('--text-h', t.colors.textHeading)
    root.style.setProperty('--bg', t.colors.bg)
    root.style.setProperty('--border', t.colors.border)
    root.style.setProperty('--code-bg', t.colors.codeBg)
    root.style.setProperty('--header-bg', t.colors.headerBg)
    root.style.setProperty('--card-bg', t.colors.cardBg)
    root.style.setProperty('--shadow', t.colors.shadow)
    root.style.setProperty('--font-display', t.fontDisplay)
    root.style.setProperty('--font-body', t.fontBody)
    document.body.style.background = t.colors.bg
  }

  return { themes, currentTheme, theme, setTheme, applyTheme }
})
