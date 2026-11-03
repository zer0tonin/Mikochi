export const getApiUrl = (path) => {
  if (
    window.__MIKOCHI_CONFIG__ != null &&
    "pathPrefix" in window.__MIKOCHI_CONFIG__
  ) {
    return `${window.__MIKOCHI_CONFIG__["pathPrefix"]}api/${path}`;
  }
  return `/api/${path}`;
};

export const getAssetUrl = (path) => {
  if (
    window.__MIKOCHI_CONFIG__ != null &&
    "pathPrefix" in window.__MIKOCHI_CONFIG__
  ) {
    return `${window.__MIKOCHI_CONFIG__["pathPrefix"]}${path}`;
  }
  return `/${path}`;
};

export const getStreamUrl = (path, token) => {
  if (
    window.__MIKOCHI_CONFIG__ != null &&
    "pathPrefix" in window.__MIKOCHI_CONFIG__
  ) {
    return `${window.location.protocol}//${window.location.hostname}${window.location.port == "" ? "" : `:${window.location.port}`}${window.__MIKOCHI_CONFIG__["pathPrefix"]}api/stream${path}?${token.toString()}`;
  }
  return `${window.location.protocol}//${window.location.hostname}${window.location.port == "" ? "" : `:${window.location.port}`}/api/stream${path}?${token.toString()}`;
};

export const getLocation = (path) => {
  if (
    window.__MIKOCHI_CONFIG__ != null &&
    "pathPrefix" in window.__MIKOCHI_CONFIG__
  ) {
    return path.slice(window.__MIKOCHI_CONFIG__["pathPrefix"].length);
  }
  return path;
}
