"use server"
import { cookies, headers } from "next/headers"

/**
 * Helper to get authentication headers for server-side API calls.
 * Works in Server Components, Server Actions, and API Routes.
 */
export async function getAuthHeaders() {
    // 1. Cek header Authorization (mungkin sudah diset oleh middleware setelah refresh)
    const headerList = await headers()
    const authHeader = headerList.get("Authorization")

    if (authHeader) {
        return {
            Authorization: authHeader,
        }
    }

    // 2. Fallback ke cookies
    const cookieStore = await cookies()
    const token = cookieStore.get("access_token")?.value

    if (!token) {
        return {}
    }

    return {
        Authorization: `Bearer ${token}`,
    }
}
