import { NextRequest, NextResponse } from "next/server"

const API_URL = process.env.API_URL

export async function POST(request: NextRequest) {
    try {
        const body = await request.json()

        const backendResponse = await fetch(`${API_URL}/register`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(body),
        })

        const data = await backendResponse.json()

        if (!backendResponse.ok) {
            return NextResponse.json(
                { message: data.message || "Register gagal" },
                { status: backendResponse.status }
            )
        }

        return NextResponse.json({
            status: "success",
            message: "Register berhasil",
        })
    } catch (error) {
        console.error("Register API error:", error)
        return NextResponse.json(
            { message: "Terjadi kesalahan server" },
            { status: 500 }
        )
    }
}