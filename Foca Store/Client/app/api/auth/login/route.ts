import { NextRequest, NextResponse } from "next/server"

const API_URL = process.env.API_URL

export async function POST(request: NextRequest) {
    try {
        const body = await request.json()

        const backendResponse = await fetch(`${API_URL}/login`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(body),
        })

        const data = await backendResponse.json()

        if (!backendResponse.ok) {
            return NextResponse.json(
                { message: data.message || "Login gagal" },
                { status: backendResponse.status }
            )
        }

        const { access_token, refresh_token, user } = data.data

        const response = NextResponse.json({
            status: "success",
            data: { user },
        })

        // Set HttpOnly cookies — tidak bisa diakses JavaScript
        response.cookies.set("access_token", access_token, {
            httpOnly: true,
            secure: process.env.NODE_ENV === "production",
            sameSite: "lax",
            path: "/",
        })

        response.cookies.set("refresh_token", refresh_token, {
            httpOnly: true,
            secure: process.env.NODE_ENV === "production",
            sameSite: "lax",
            path: "/",
        })

        return response
    } catch (error) {
        console.error("Login API error:", error)
        return NextResponse.json(
            { message: "Terjadi kesalahan server" },
            { status: 500 }
        )
    }
}
