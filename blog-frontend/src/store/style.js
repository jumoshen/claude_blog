import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useStyleStore = defineStore('style', () => {
  // Helper to create dark mode variant
  const withDarkMode = (colors, isDark) => {
    if (!isDark) return colors
    return {
      accent: colors.accent,
      accentRgb: colors.accentRgb,
      accentBg: colors.accentBg ? colors.accentBg.replace(/[\d.]+\)$/, '0.2)') : 'rgba(179, 102, 255, 0.2)',
      accentBorder: colors.accentBorder ? colors.accentBorder.replace(/[\d.]+\)$/, '0.3)') : 'rgba(179, 102, 255, 0.3)',
      text: '#a8a8b3',
      textHeading: '#e8e8ed',
      bg: '#0f0f12',
      border: '#2a2a35',
      codeBg: '#1a1a24',
      headerBg: 'rgba(15, 15, 18, 0.9)',
      cardBg: 'rgba(30, 30, 40, 0.8)',
      shadow: 'rgba(0, 0, 0, 0.4)',
    }
  }

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
    pixelDark: {
      name: '像素风·暗',
      tagline: '8-BIT RETRO',
      logo: '/logo-pixel.png',
      icon: '🎮',
      preview: 'linear-gradient(135deg, #b366ff 0%, #7c3aed 100%)',
      isDark: true,
      colors: {
        accent: '#b366ff',
        accentRgb: '179, 102, 255',
        accentBg: 'rgba(179, 102, 255, 0.15)',
        accentBorder: 'rgba(179, 102, 255, 0.4)',
        text: '#a8a8b3',
        textHeading: '#e8e8ed',
        bg: '#0a0a0e',
        border: '#2a2a38',
        codeBg: '#1a1a28',
        headerBg: 'rgba(10, 10, 14, 0.95)',
        cardBg: 'rgba(20, 20, 30, 0.9)',
        shadow: 'rgba(179, 102, 255, 0.15)',
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
    cuteDark: {
      name: '可爱风·暗',
      tagline: 'SWEETY DAYS',
      logo: '/logo-cute.png',
      icon: '🌸',
      preview: 'linear-gradient(135deg, #ff6b9d 0%, #ffa8c5 100%)',
      isDark: true,
      colors: {
        accent: '#ff6b9d',
        accentRgb: '255, 107, 157',
        accentBg: 'rgba(255, 107, 157, 0.15)',
        accentBorder: 'rgba(255, 107, 157, 0.4)',
        text: '#b8a8b3',
        textHeading: '#f0e8ed',
        bg: '#12101a',
        border: '#2a2035',
        codeBg: '#1a1525',
        headerBg: 'rgba(18, 16, 26, 0.95)',
        cardBg: 'rgba(30, 25, 40, 0.9)',
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
    },
    qverDark: {
      name: 'Q版·暗',
      tagline: 'CHILL VIBES',
      logo: '/logo-qver.png',
      icon: '⚡',
      preview: 'linear-gradient(135deg, #7c6aff 0%, #a78bfa 100%)',
      isDark: true,
      colors: {
        accent: '#7c6aff',
        accentRgb: '124, 106, 255',
        accentBg: 'rgba(124, 106, 255, 0.15)',
        accentBorder: 'rgba(124, 106, 255, 0.4)',
        text: '#a8a8c3',
        textHeading: '#e0e0f0',
        bg: '#0d0d18',
        border: '#252540',
        codeBg: '#181828',
        headerBg: 'rgba(13, 13, 24, 0.95)',
        cardBg: 'rgba(25, 25, 45, 0.9)',
        shadow: 'rgba(124, 106, 255, 0.15)',
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
    },
  }

  const currentTheme = ref(localStorage.getItem('blog-theme') || 'pixel')
  const isDark = ref(localStorage.getItem('blog-dark-mode') === 'true')
  const theme = computed(() => {
    const baseTheme = themes[currentTheme.value]
    if (!baseTheme) return themes.pixel
    // If dark mode is on and we're not already on a dark variant, switch to dark variant
    if (isDark.value && !currentTheme.value.endsWith('Dark')) {
      const darkVariant = themes[currentTheme.value + 'Dark']
      if (darkVariant) return darkVariant
    }
    return baseTheme
  })

  function setTheme(name) {
    // Strip 'Dark' suffix if present when setting theme
    const baseName = name.replace('Dark', '')
    if (themes[name] || themes[baseName]) {
      currentTheme.value = themes[name] ? name : baseName
      localStorage.setItem('blog-theme', currentTheme.value)
      applyTheme(currentTheme.value)
    }
  }

  function toggleDarkMode() {
    isDark.value = !isDark.value
    localStorage.setItem('blog-dark-mode', isDark.value.toString())
    // Apply the appropriate variant
    const themeName = isDark.value ? (currentTheme.value + 'Dark') : currentTheme.value
    const targetTheme = themes[themeName] || themes[currentTheme.value]
    if (targetTheme) {
      applyTheme(themeName)
    }
  }

  function setDarkMode(dark) {
    isDark.value = dark
    localStorage.setItem('blog-dark-mode', isDark.value.toString())
    const themeName = isDark.value ? (currentTheme.value + 'Dark') : currentTheme.value
    const targetTheme = themes[themeName] || themes[currentTheme.value]
    if (targetTheme) {
      applyTheme(themeName)
    }
  }

  function applyTheme(name) {
    let t = themes[name]
    // If dark mode is on and we're applying a base theme, switch to dark variant
    if (!t) return
    if (isDark.value && !name.endsWith('Dark')) {
      const darkVariant = themes[name + 'Dark']
      if (darkVariant) t = darkVariant
    }

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

    // Store actual applied theme name (might be dark variant)
    const actualThemeName = t.name ? name : (isDark.value ? name + 'Dark' : name)
    // 更新Favicon
    const favicon = document.getElementById('favicon')
    if (favicon) {
      favicon.href = '/logo-' + name.replace('Dark', '') + '.ico'
    }
  }

  return { themes, currentTheme, isDark, theme, setTheme, applyTheme, toggleDarkMode, setDarkMode }
})
