export const getApiUrl = (path) => {
  if (window.__MIKOCHI_CONFIG__ != null && 'pathPrefix' in window.__MIKOCHI_CONFIG__) {
    return `${window.__MIKOCHI_CONFIG__['pathPrefix']}api/${path}`
  }
  return `/api/${path}`
}
