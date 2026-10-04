// ApiError messages carry the raw response body, which is {"error": "..."} for backend errors.
export function errorMessage(err: unknown): string {
  const raw = err instanceof Error ? err.message : String(err)
  try {
    const parsed = JSON.parse(raw) as { error?: string }
    return parsed.error ?? raw
  } catch {
    return raw
  }
}
