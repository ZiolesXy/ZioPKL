import { NextRequest, NextResponse } from "next/server"
import { cookies } from "next/headers"

const API_URL = process.env.API_URL

async function refreshTokens() {
    const cookieStore = await cookies()
    const refreshToken = cookieStore.get("refresh_token")?.value

    if (!refreshToken) {
        console.log("Auto-refresh: No refresh token found in cookies")
        return null
    }

    try {
        console.log("Auto-refresh: Attempting to refresh token...")
        const response = await fetch(`${API_URL}/refresh`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ refresh_token: refreshToken }),
        })

        if (!response.ok) {
            console.log(`Auto-refresh: Refresh request failed with status ${response.status}`)
            return null
        }

        const data = await response.json()
        // Coba ambil dari data.data (seperti login) atau langsung dari data
        const tokens = data.data || data
        console.log("Auto-refresh: Token refreshed successfully")
        return tokens
    } catch (error) {
        console.error("Auto-refresh: Error during refresh process:", error)
        return null
    }
}

async function proxyRequest(request: NextRequest, method: string) {
    const cookieStore = await cookies()
    let token = cookieStore.get("access_token")?.value

    const url = new URL(request.url)
    const path = url.pathname.replace(/^\/api\/proxy/, "")
    const targetUrl = `${API_URL}${path}${url.search}`

    const buildHeaders = (t: string | undefined) => {
        const h: Record<string, string> = {}
        if (t) {
            h["Authorization"] = `Bearer ${t}`
        }
        const ct = request.headers.get("content-type")
        if (ct && !ct.includes("multipart/form-data")) {
            h["Content-Type"] = ct
        }
        return h
    }

    const getBody = async () => {
        if (method === "GET" || method === "HEAD") return undefined
        const ct = request.headers.get("content-type")
        if (ct?.includes("multipart/form-data")) {
            return await request.formData()
        }
        return await request.text()
    }

    const body = await getBody()

    try {
        const backendResponse = await fetch(targetUrl, {
            method,
            headers: buildHeaders(token),
            body,
        })

        // Jika unauthorized, coba refresh token
        if (backendResponse.status === 401) {
            console.log(`Auto-refresh: Detected 401 for ${path}, attempting refresh...`)
            const newData = await refreshTokens()

            if (newData && newData.access_token) {
                token = newData.access_token
                console.log(`Auto-refresh: Retrying original request ${path} with new token`)

                // Retry original request with new token
                const retryResponse = await fetch(targetUrl, {
                    method,
                    headers: buildHeaders(token),
                    body,
                })

                if (retryResponse.ok) {
                    console.log(`Auto-refresh: Retry successful for ${path}`)
                    const responseText = await retryResponse.text()

                    const responseHeaders = new Headers()
                    retryResponse.headers.forEach((value, key) => {
                        // Forward all headers except set-cookie and content-encoding
                        if (key.toLowerCase() !== "set-cookie" && key.toLowerCase() !== "content-encoding") {
                            responseHeaders.set(key, value)
                        }
                    })

                    const response = new NextResponse(responseText, {
                        status: retryResponse.status,
                        headers: responseHeaders,
                    })

                    // Update cookies with new tokens
                    response.cookies.set("access_token", newData.access_token, {
                        httpOnly: true,
                        secure: process.env.NODE_ENV === "production",
                        sameSite: "lax",
                        path: "/",
                    })

                    if (newData.refresh_token) {
                        response.cookies.set("refresh_token", newData.refresh_token, {
                            httpOnly: true,
                            secure: process.env.NODE_ENV === "production",
                            sameSite: "lax",
                            path: "/",
                        })
                    }

                    return response
                } else {
                    console.log(`Auto-refresh: Retry failed with status ${retryResponse.status}`)
                }
            }

            // Jika refresh gagal atau retry tetap gagal, logout paksa
            console.log("Auto-refresh: Refresh or retry failed, clearing sessions")
            const response = NextResponse.json(
                { message: "Sesi telah berakhir, silakan login kembali" },
                { status: 401 }
            )
            response.cookies.delete("access_token")
            response.cookies.delete("refresh_token")
            return response
        }

        const responseData = await backendResponse.text()
        return new NextResponse(responseData, {
            status: backendResponse.status,
            headers: {
                "Content-Type": backendResponse.headers.get("Content-Type") || "application/json",
            },
        })
    } catch (error) {
        console.error(`Proxy error [${method} ${path}]:`, error)
        return NextResponse.json(
            { message: "Gagal menghubungi server backend" },
            { status: 502 }
        )
    }
}

export async function GET(request: NextRequest) {
    return proxyRequest(request, "GET")
}

export async function POST(request: NextRequest) {
    return proxyRequest(request, "POST")
}

export async function PUT(request: NextRequest) {
    return proxyRequest(request, "PUT")
}

export async function PATCH(request: NextRequest) {
    return proxyRequest(request, "PATCH")
}

export async function DELETE(request: NextRequest) {
    return proxyRequest(request, "DELETE")
}
