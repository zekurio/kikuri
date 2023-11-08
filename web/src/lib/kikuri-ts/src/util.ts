export function replaceDoublePath(url: string): string {
  const split = url.split("://");
  split[split.length - 1] = split[split.length - 1].replace(/\/\//g, "/");
  return split.join("://");
}
