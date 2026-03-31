import { NextResponse } from "next/server"

const BACKEND_BASE_URL = "http://172.16.17.123:8080"

export async function POST(req: Request) {
  const body = await req.json().catch(() => null)
  if (!body?.email || !body?.password) {
    return NextResponse.json({ message: "email & password required" }, { status: 400 })
  }

  const upstream = await fetch(`${BACKEND_BASE_URL}/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email: body.email, password: body.password }),
  })

  const data = await upstream.json().catch(() => null)

  if (!upstream.ok) {
    return NextResponse.json(data ?? { message: "Login failed" }, { status: upstream.status })
  }

  const accessToken = data?.data?.access_token
  const refreshToken = data?.data?.refresh_token
  const role = data?.data?.role

  if (!accessToken || !refreshToken) {
    return NextResponse.json({ message: "Missing tokens from backend" }, { status: 500 })
  }

  const res = NextResponse.json({ ok: true, role })

  res.cookies.set("access_token", accessToken, {
    httpOnly: false,
    sameSite: "lax",
    path: "/",
  })

  res.cookies.set("refresh_token", refreshToken, {
    httpOnly: true,
    sameSite: "lax",
    path: "/",
  })

  return res
}
