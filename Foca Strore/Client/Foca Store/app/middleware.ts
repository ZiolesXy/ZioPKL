import { NextResponse } from "next/server"
import type { NextRequest } from "next/server"

type JwtPayload = {
  role?: number
  exp?: number
}

const BACKEND_BASE_URL = "http://172.16.17.123:8080"

function decodeJwtPayload(token: string): JwtPayload | null {
  try {
    const parts = token.split(".")
    if (parts.length < 2) return null

    const base64 = parts[1].replace(/-/g, "+").replace(/_/g, "/")
    const padded = base64.padEnd(base64.length + ((4 - (base64.length % 4)) % 4), "=")
    const json = atob(padded)
    return JSON.parse(json)
  } catch {
    return null
  }
}

export async function middleware(req: NextRequest) {
  const refreshToken = req.cookies.get("refresh_token")?.value
  const accessToken = req.cookies.get("access_token")?.value

  const protectedPaths = [
    "/dashboard",
    "/delete",
    "/post",
    "/put",
  ]

  const adminOnlyPaths = ["/delete", "/post", "/put"]

  const isProtected = protectedPaths.some(path =>
    req.nextUrl.pathname.startsWith(path)
  )

  if (isProtected && !refreshToken) {
    return NextResponse.redirect(new URL("/", req.url))
  }

  let role: number | undefined

  const now = Math.floor(Date.now() / 1000)
  const payload = accessToken ? decodeJwtPayload(accessToken) : null
  const exp = typeof payload?.exp === "number" ? payload.exp : undefined
  const isAccessValid = !!accessToken && !!payload && !!exp && exp > now
  if (isAccessValid) {
    role = typeof payload?.role === "number" ? payload.role : undefined
  }

  // If route is protected and access token is missing/expired, refresh it using refresh_token
  let nextResponse: NextResponse | null = null
  let effectiveRole = role
  if (isProtected && !isAccessValid && refreshToken) {
    const refreshUpstream = await fetch(`${BACKEND_BASE_URL}/refresh-token`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ refresh_token: refreshToken }),
    })

    const refreshJson = await refreshUpstream.json().catch(() => null)

    if (!refreshUpstream.ok) {
      const res = NextResponse.redirect(new URL("/", req.url))
      res.cookies.delete("access_token")
      res.cookies.delete("refresh_token")
      return res
    }

    const newAccessToken = refreshJson?.data?.access_token
    if (!newAccessToken) {
      const res = NextResponse.redirect(new URL("/", req.url))
      res.cookies.delete("access_token")
      res.cookies.delete("refresh_token")
      return res
    }

    const newPayload = decodeJwtPayload(newAccessToken)
    effectiveRole = typeof newPayload?.role === "number" ? newPayload.role : undefined

    // Server Components in the *current* request can't see cookies that are only set on the response.
    // So after refreshing we redirect to the same URL, so the next request includes the new access_token.
    nextResponse = NextResponse.redirect(req.nextUrl)
    nextResponse.cookies.set("access_token", newAccessToken, {
      httpOnly: false,
      sameSite: "lax",
      path: "/",
    })
  }

  // Authorization: role 1 = Admin. Other roles only allowed to access /dashboard
  const isAdminOnly = adminOnlyPaths.some(path => req.nextUrl.pathname.startsWith(path))
  if (isProtected && isAdminOnly && effectiveRole !== 1) {
    return NextResponse.redirect(new URL("/dashboard", req.url))
  }

  return nextResponse ?? NextResponse.next()
}

export const config = {
  matcher: ["/dashboard/:path*", "/delete/:path*", "/post/:path*", "/put/:path*"],
}