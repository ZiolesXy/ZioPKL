"use server"
import { cookies, headers } from "next/headers"


export async function getAuthHeaders() {
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
