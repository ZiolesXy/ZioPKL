"use server"
import { cookies } from "next/headers"

const BACKEND_BASE_URL = "http://172.16.17.123:8080"

export async function serverFetchWithAuth(path: string, init?: RequestInit) {
  const cookieStore = await cookies()
  const accessToken = cookieStore.get("access_token")?.value
  const refreshToken = cookieStore.get("refresh_token")?.value

  const headers = new Headers(init?.headers)
  headers.set("Content-Type", "application/json")

  if (accessToken) {
    headers.set("Authorization", /^Bearer\s+/i.test(accessToken) ? accessToken : `Bearer ${accessToken}`)
  }

  const doFetch = (h: Headers) =>
    fetch(`${BACKEND_BASE_URL}${path}`, {
      ...init,
      headers: h,
      cache: "no-store",
    })

  const res = await doFetch(headers)

  if (res.status !== 401 || !refreshToken) {
    return res
  }

  const refreshUpstream = await fetch(`${BACKEND_BASE_URL}/refresh-token`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ refresh_token: refreshToken }),
    cache: "no-store",
  })

  if (!refreshUpstream.ok) return res

  const refreshJson = await refreshUpstream.json().catch(() => null)
  const newAccessToken = refreshJson?.data?.access_token
  if (!newAccessToken || typeof newAccessToken !== "string") return res

  const retryHeaders = new Headers(init?.headers)
  retryHeaders.set("Content-Type", "application/json")
  retryHeaders.set(
    "Authorization",
    /^Bearer\s+/i.test(newAccessToken) ? newAccessToken : `Bearer ${newAccessToken}`
  )

  return doFetch(retryHeaders)

  
}
