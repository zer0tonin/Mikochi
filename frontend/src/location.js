export const removePrefix = (path) => {
  if (
    window.__MIKOCHI_CONFIG__ != null &&
    "pathPrefix" in window.__MIKOCHI_CONFIG__
  ) {
    return path.slice(window.__MIKOCHI_CONFIG__["pathPrefix"].length);
  }
  return path;
};

export const addPrefix = (path) => {
  if (
    window.__MIKOCHI_CONFIG__ != null &&
    "pathPrefix" in window.__MIKOCHI_CONFIG__
  ) {
    if (path.startsWith("/")) {
      return `${window.__MIKOCHI_CONFIG__["pathPrefix"]}${path.slice(1)}`;
    }
    return `${window.__MIKOCHI_CONFIG__["pathPrefix"]}${path}`;
  }
  return path;
};
