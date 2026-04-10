// v-lazy-load directive for image lazy loading
export const lazyLoad = {
  mounted(el, binding) {
    const options = {
      root: null,
      rootMargin: '50px 0px',
      threshold: 0.1
    }

    const observer = new IntersectionObserver((entries) => {
      entries.forEach(entry => {
        if (entry.isIntersecting) {
          const img = entry.target
          const src = img.getAttribute('data-src')
          const originalSrc = img.src

          if (src && src !== originalSrc) {
            // Add loading state
            img.classList.add('lazy-loading')

            // Create a new image to preload
            const preloadImg = new Image()
            preloadImg.onload = () => {
              img.src = src
              img.classList.remove('lazy-loading')
              img.classList.add('lazy-loaded')
            }
            preloadImg.onerror = () => {
              img.classList.remove('lazy-loading')
              img.classList.add('lazy-error')
            }
            preloadImg.src = src
          }

          // Stop observing once loaded
          observer.unobserve(img)
        }
      })
    }, options)

    // Store observer reference on element for cleanup
    el._lazyObserver = observer

    // Observe all images with data-src
    const images = el.querySelectorAll ? el.querySelectorAll('img[data-src]') : []
    images.forEach(img => {
      observer.observe(img)
    })

    // If the element itself is an image
    if (el.tagName === 'IMG' && el.getAttribute('data-src')) {
      observer.observe(el)
    }
  },

  unmounted(el) {
    if (el._lazyObserver) {
      el._lazyObserver.disconnect()
      delete el._lazyObserver
    }
  }
}

// Add CSS for lazy loading states
const style = document.createElement('style')
style.textContent = `
  .lazy-loading {
    opacity: 0.5;
    transition: opacity 0.3s ease-out;
  }
  .lazy-loaded {
    opacity: 1;
    transition: opacity 0.3s ease-out;
  }
  .lazy-error {
    opacity: 0.5;
    filter: grayscale(100%);
  }
`
document.head.appendChild(style)
