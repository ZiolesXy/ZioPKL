import axios from "axios"

function toBearer(token: string): string {
  const t = token.trim()
  return /^Bearer\s+/i.test(t) ? t : `Bearer ${t}`
}

let refreshPromise: Promise<boolean> | null = null

async function refreshAccessToken(): Promise<boolean> {
  if (typeof window === "undefined") return false

  try {
    const res = await axios.post(
      "/api/auth/refresh",
      {},
      {
        withCredentials: true,
        headers: {
          "Content-Type": "application/json",
        },
        validateStatus: () => true,
      }
    )
    return res.status >= 200 && res.status < 300
  } catch {
    return false
  }
}

const api = axios.create({
  baseURL: "http://172.16.17.123:8080", 
  withCredentials: true,
  headers: {
    "Content-Type": "application/json",
  },
})

api.interceptors.request.use((config) => {
  if (typeof document !== "undefined") {
    const token = document.cookie
      .split(";")
      .map((c) => c.trim())
      .find((c) => c.startsWith("access_token="))
      ?.split("=")
      ?.slice(1)
      .join("=")

    if (token) {
      config.headers = config.headers ?? {}
      config.headers.Authorization = toBearer(decodeURIComponent(token))
    }
  }

  return config
})

api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error?.config as
      | (typeof error.config & { _retry?: boolean })
      | undefined

    const status = error?.response?.status

    if (!originalRequest || status !== 401 || originalRequest._retry) {
      return Promise.reject(error)
    }

    originalRequest._retry = true

    refreshPromise = refreshPromise ?? refreshAccessToken()
    const ok = await refreshPromise.finally(() => {
      refreshPromise = null
    })

    if (!ok) return Promise.reject(error)

    return api(originalRequest)
  }
)

export default api