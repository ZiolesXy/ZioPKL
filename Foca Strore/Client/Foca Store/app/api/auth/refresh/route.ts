import { NextResponse } from "next/server"

const BACKEND_BASE_URL = "http://172.16.17.123:8080"

export async function POST(req: Request) {
  const refreshToken = req.headers
    .get("cookie")
    ?.split(";")
    .map((c) => c.trim())
    .find((c) => c.startsWith("refresh_token="))
    ?.split("=")
    ?.slice(1)
    .join("=")

  if (!refreshToken) {
    return NextResponse.json({ message: "refresh_token missing" }, { status: 401 })
  }

  const upstream = await fetch(`${BACKEND_BASE_URL}/refresh-token`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ refresh_token: refreshToken }),
  })

  const data = await upstream.json().catch(() => null)

  if (!upstream.ok) {
    const res = NextResponse.json(data ?? { message: "Refresh failed" }, { status: upstream.status })
    res.cookies.delete("access_token")
    res.cookies.delete("refresh_token")
    return res
  }

  const accessToken = data?.data?.access_token
  if (!accessToken) {
    return NextResponse.json({ message: "Missing access_token from backend" }, { status: 500 })
  }

  const res = NextResponse.json({ ok: true })
  res.cookies.set("access_token", accessToken, {
    httpOnly: false,
    sameSite: "lax",
    path: "/",
  })

  return res
}
