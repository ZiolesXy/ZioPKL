export function extractArrayFromPayload<T>(
  payload: unknown,
  options: {
    keys?: ["entries"]
    label: string
  }
): T[] {
  if (Array.isArray(payload)) return payload as T[]

  if (payload === null || typeof payload !== "object") {
    throw new Error(
      `Unexpected ${options.label} response shape: ${typeof payload === "string" ? payload : JSON.stringify(payload)}`
    )
  }

  const obj = payload as Record<string, unknown>

  const keys = options.keys ?? []
  for (const key of ["data", ...keys, "result", "payload"]) {
    const value = obj[key]
    if (Array.isArray(value)) return value as T[]
  }

  throw new Error(
    `Unexpected ${options.label} response shape: ${typeof payload === "string" ? payload : JSON.stringify(payload)}`
  )
}
