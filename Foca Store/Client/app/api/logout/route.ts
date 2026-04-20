import { NextResponse } from "next/server"
import { cookies } from "next/headers"

const API_URL = process.env.API_URL

export async function POST() {
  const cookieStore = await cookies()
  const token = cookieStore.get("access_token")?.value

  // Memberitahu backend untuk blacklist token jika token ada
  if (token) {
    try {
      console.log("Logging out from backend...")
      const res = await fetch(`${API_URL}/api/logout`, {
        method: "POST",
        headers: {
          "Authorization": `Bearer ${token}`,
          "Content-Type": "application/json",
        },
      })
      console.log(`Backend logout status: ${res.status}`)
    } catch (error) {
      console.error("Backend logout error:", error)
      // Tetap lanjut untuk menghapus cookie di sisi client
    }
  }

  const response = NextResponse.json({ ok: true })

  response.cookies.set("access_token", "", {
    httpOnly: true,
    secure: process.env.NODE_ENV === "production",
    sameSite: "lax",
    path: "/",
    maxAge: 0,
  })

  response.cookies.set("refresh_token", "", {
    httpOnly: true,
    secure: process.env.NODE_ENV === "production",
    sameSite: "lax",
    path: "/",
    maxAge: 0,
  })

  return response
}
