export function stripSuffix(v: string, suffix: string): string {
  if (!v.endsWith(suffix)) return v;
  return v.substring(0, v.length - suffix.length);
}
