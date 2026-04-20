import { NextResponse } from "next/server"
import type { NextRequest } from "next/server"
import { jwtDecode } from "jwt-decode"
const API_URL = process.env.API_URL || "http://172.16.17.123:8080"

 type TokenPayload = {
     exp: number
     role: "Admin" | "User"
 }

export async function middleware(request: NextRequest) {
    const { pathname } = request.nextUrl

    // 1. Skip static assets and internal next requests FIRST
    if (
        pathname.startsWith("/_next") ||
        pathname.startsWith("/static") ||
        pathname.includes(".") ||
        pathname.startsWith("/api/auth") // Auth APIs manage their own cookies
    ) {
        return NextResponse.next()
    }

    let accessToken = request.cookies.get("access_token")?.value
    const refreshToken = request.cookies.get("refresh_token")?.value
    let nextResponse: NextResponse | null = null

    // Helper to check if token is expired
    const isExpired = (token: string): boolean => {
        try {
            const decoded = jwtDecode<TokenPayload>(token)
            return decoded.exp * 1000 < Date.now() + 5000 // 5s buffer
        } catch {
            return true
        }
    }

    // 2. TOKEN REFRESH LOGIC
    if ((!accessToken || isExpired(accessToken)) && refreshToken) {
        console.log(`Middleware: Token expires/missing, attempting refresh for ${pathname}`)
        try {
            const response = await fetch(`${API_URL}/refresh`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ refresh_token: refreshToken }),
            })

            if (response.ok) {
                const data = await response.json()
                const tokens = data.data || data

                if (tokens && tokens.access_token) {
                    console.log("Middleware: Token refreshed successfully")
                    accessToken = tokens.access_token

                    const requestHeaders = new Headers(request.headers)
                    requestHeaders.set("Authorization", `Bearer ${accessToken}`)

                    nextResponse = NextResponse.next({
                        request: { headers: requestHeaders },
                    })

                    if (accessToken) {
                        nextResponse.cookies.set("access_token", accessToken, {
                            httpOnly: true,
                            secure: process.env.NODE_ENV === "production",
                            sameSite: "lax",
                            path: "/",
                        })
                    }

                    if (tokens.refresh_token) {
                        nextResponse.cookies.set("refresh_token", tokens.refresh_token, {
                            httpOnly: true,
                            secure: process.env.NODE_ENV === "production",
                            sameSite: "lax",
                            path: "/",
                        })
                    }
                }
            }
        } catch (error) {
            console.error("Middleware: Refresh error:", error)
        }
    }

    // 3. ROUTE PROTECTION & REDIRECTION LOGIC (formerly proxy.ts)
    let role: "Admin" | "User" | null = null

    if (accessToken) {
        try {
            const decoded = jwtDecode<TokenPayload>(accessToken)
            if (decoded.exp * 1000 > Date.now()) {
                role = decoded.role
            }
        } catch (error) {
            console.error("Middleware: Token decode error:", error)
        }
    }

    // -- Redirection Rules --

    const isPublicRoute =
        pathname === "/" ||
        pathname.startsWith("/product") ||
        pathname.startsWith("/category") ||
        pathname.startsWith("/privacy") ||
        pathname.startsWith("/terms") ||
        pathname.startsWith("/refund")

    const isAuthRoute =
        pathname === "/login" ||
        pathname === "/register" ||
        pathname.startsWith("/forgot-password") ||
        pathname.startsWith("/reset-password")

    const isProtectedUserRoute =
        pathname.startsWith("/orders") ||
        pathname.startsWith("/coupons") ||
        pathname.startsWith("/notifications") ||
        pathname.startsWith("/update-profile") ||
        pathname.startsWith("/chat")

    const isProtectedAdminRoute =
        pathname.startsWith("/overview") ||
        pathname.startsWith("/add-products") ||
        pathname.startsWith("/check-products") ||
        pathname.startsWith("/customer-service") ||
        pathname.startsWith("/system") 

    // Logged in user hitting auth pages
    if (accessToken && role && isAuthRoute) {
        const redirectUrl = role === "Admin" ? "/overview" : "/"
        const response = NextResponse.redirect(new URL(redirectUrl, request.url))
        // If we recently refreshed, we MUST persist those cookies to the redirect response
        if (nextResponse) {
            nextResponse.cookies.getAll().forEach((c) => response.cookies.set(c.name, c.value, c))
        }
        return response
    }

    // Non-logged in user hitting protected routes
    if (!accessToken && (isProtectedUserRoute || isProtectedAdminRoute) && !pathname.startsWith("/api")) {
        const loginUrl = new URL("/login", request.url)
        loginUrl.searchParams.set("next", pathname)
        return NextResponse.redirect(loginUrl)
    }

    // Admin-only protection
    if (isProtectedAdminRoute && role !== "Admin") {
        return NextResponse.redirect(new URL("/", request.url))
    }

    // NOTE: public browsing should remain accessible even when not logged in.
    // If additional routes should be public, extend isPublicRoute above.
    if (!isPublicRoute && !isAuthRoute && !isProtectedUserRoute && !isProtectedAdminRoute) {
        // Default: allow access (no forced login) to avoid accidentally locking users out
        // if new routes are added.
    }

    return nextResponse || NextResponse.next()
}

// Optionally define which routes this middleware should run on
export const config = {
    matcher: [
        /*
         * Match all request paths except for the ones starting with:
         * - api/auth (auth endpoints)
         * - api/proxy/refresh (already handled locally)
         * - _next/static (static files)
         * - _next/image (image optimization files)
         * - favicon.ico (favicon file)
         */
        "/((?!api/auth|_next/static|_next/image|favicon.ico).*)",
    ],
}
