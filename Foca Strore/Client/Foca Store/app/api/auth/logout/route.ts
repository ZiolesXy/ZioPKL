import { NextResponse } from "next/server"
import { cookies } from "next/headers"

const BACKEND_BASE_URL = "http://172.16.17.123:8080"

export async function POST() {
  const cookieStore = await cookies()
  const refreshToken = cookieStore.get("refresh_token")?.value

  let upstreamOk = true
  let upstreamStatus = 200
  let upstreamBody: unknown = null

  if (refreshToken) {
    const upstream = await fetch(`${BACKEND_BASE_URL}/logout`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "X-Refresh-Token": decodeURIComponent(refreshToken),
      },
    })

    upstreamOk = upstream.ok
    upstreamStatus = upstream.status
    upstreamBody = await upstream.json().catch(() => null)
  }

  const res = NextResponse.json(
    refreshToken
      ? upstreamOk
        ? upstreamBody ?? { ok: true }
        : upstreamBody ?? { ok: false, message: "Logout failed" }
      : { ok: true, message: "refresh_token missing; cleared cookies" },
    { status: refreshToken ? upstreamStatus : 200 }
  )

  res.cookies.delete("access_token")
  res.cookies.delete("refresh_token")
  return res
}
