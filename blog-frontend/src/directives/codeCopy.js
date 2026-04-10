// Directive v-code-copy for adding copy buttons to code blocks
export const codeCopy = {
  mounted(el, binding) {
    // Find all pre/code blocks in the element
    const codeBlocks = el.querySelectorAll('pre')

    codeBlocks.forEach(pre => {
      // Create copy button
      const button = document.createElement('button')
      button.className = 'code-copy-btn'
      button.innerHTML = `
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
          <path d="M5 15H4a2 2 0 01-2-2V4a2 2 0 012-2h9a2 2 0 012 2v1"></path>
        </svg>
        <span class="copy-text">复制</span>
      `
      button.style.cssText = `
        position: absolute;
        top: 8px;
        right: 8px;
        padding: 6px 10px;
        background: var(--accent-bg, rgba(179, 102, 255, 0.1));
        border: 1px solid var(--accent-border, rgba(179, 102, 255, 0.3));
        border-radius: 8px;
        color: var(--accent, #b366ff);
        font-size: 12px;
        cursor: pointer;
        display: flex;
        align-items: center;
        gap: 4px;
        opacity: 0;
        transition: all 0.2s;
        z-index: 10;
      `

      // Add hover styles via CSS
      const style = document.createElement('style')
      style.textContent = `
        .code-copy-btn:hover {
          background: var(--accent, #b366ff) !important;
          color: #fff !important;
          transform: scale(1.05);
        }
        .code-copy-btn.copied {
          background: #10b981 !important;
          border-color: #10b981 !important;
          color: #fff !important;
        }
        pre { position: relative; }
        pre:hover .code-copy-btn { opacity: 1; }
      `
      document.head.appendChild(style)

      pre.style.position = 'relative'
      pre.appendChild(button)

      button.addEventListener('click', async () => {
        const code = pre.querySelector('code')
        const text = code ? code.textContent : pre.textContent

        try {
          await navigator.clipboard.writeText(text)
          button.classList.add('copied')
          button.innerHTML = `
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="20,6 9,17 4,12"></polyline>
            </svg>
            <span class="copy-text">已复制</span>
          `
          setTimeout(() => {
            button.classList.remove('copied')
            button.innerHTML = `
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
                <path d="M5 15H4a2 2 0 01-2-2V4a2 2 0 012-2h9a2 2 0 012 2v1"></path>
              </svg>
              <span class="copy-text">复制</span>
            `
          }, 2000)
        } catch (err) {
          console.error('Failed to copy:', err)
        }
      })
    })
  }
}
