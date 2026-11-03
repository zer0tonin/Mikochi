export const removePrefix = (path) => {
  if (
    window.__MIKOCHI_CONFIG__ != null &&
    "pathPrefix" in window.__MIKOCHI_CONFIG__
  ) {
    return path.slice(window.__MIKOCHI_CONFIG__["pathPrefix"].length);
  }
  return path;
}

export const addPrefix = (path) => {
  if (
    window.__MIKOCHI_CONFIG__ != null &&
    "pathPrefix" in window.__MIKOCHI_CONFIG__
  ) {
    return `${window.__MIKOCHI_CONFIG__["pathPrefix"]}${path}` // FIXME: frequenty introduces double slash
  }
  return path;
}
