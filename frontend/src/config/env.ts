export const env = {
  // empty means relative URLs: Vite proxy in dev, nginx in production
  apiBaseUrl: import.meta.env.VITE_API_BASE_URL ?? '',
}
