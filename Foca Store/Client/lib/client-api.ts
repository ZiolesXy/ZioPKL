
const PROXY_BASE = "/api/proxy"

interface FetchOptions {
  method?: string
  body?: unknown
  headers?: Record<string, string>
}

export class ClientApiError extends Error {
  status: number
  data: Record<string, unknown>

  constructor(
    message: string,
    status: number,
    data: Record<string, unknown> = {}
  ) {
    super(message)
    this.name = "ClientApiError"
    this.status = status
    this.data = data
  }
}

async function clientFetch<T>(
  path: string,
  options: FetchOptions = {}
): Promise<T> {
  const { method = "GET", body, headers = {} } = options

  const fetchOptions: RequestInit = {
    method,
    headers: {
      ...headers,
    },
  }

  if (body) {
    if (body instanceof FormData) {
      // FormData — biarkan browser set Content-Type dengan boundary
      fetchOptions.body = body
    } else {
      fetchOptions.headers = {
        "Content-Type": "application/json",
        ...headers,
      }
      fetchOptions.body = JSON.stringify(body)
    }
  }

  const response = await fetch(`${PROXY_BASE}${path}`, fetchOptions)

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({}))
    throw new ClientApiError(
      errorData.message || `Request failed: ${response.status}`,
      response.status,
      errorData
    )
  }

  return response.json()
}

export const clientApi = {
  get: <T = unknown>(path: string, headers?: Record<string, string>) => clientFetch<T>(path, { headers }),

  post: <T = unknown>(path: string, body?: unknown, headers?: Record<string, string>) =>
    clientFetch<T>(path, { method: "POST", body, headers }),

  put: <T = unknown>(path: string, body?: unknown, headers?: Record<string, string>) =>
    clientFetch<T>(path, { method: "PUT", body, headers }),

  patch: <T = unknown>(path: string, body?: unknown, headers?: Record<string, string>) =>
    clientFetch<T>(path, { method: "PATCH", body, headers }),

  delete: <T = unknown>(path: string, body?: unknown, headers?: Record<string, string>) =>
    clientFetch<T>(path, { method: "DELETE", body }),

  request: <T = unknown>(options: { url: string; method: string; data?: unknown; headers?: Record<string, string> }) =>
    clientFetch<T>(options.url, { method: options.method, body: options.data, headers: options.headers }),
}
