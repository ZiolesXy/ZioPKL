"use server"

import { cookies } from "next/headers"

export async function getWsConfig() {
    const cookieStore = await cookies()
    const token = cookieStore.get("access_token")?.value || null
    
    // Fallback to localhost if env var is missing during dev
    const apiUrl = process.env.API_URL || "http://172.16.17.11:8080"
    
    // Replace http(s) with ws(s)
    let wsUrl = apiUrl.replace(/^http/, "ws")
    
    // Append /ws per API specs (the base path for WebSocket is /ws)
    if (!wsUrl.endsWith("/ws")) {
        // Only append if it's not already there by some configured path
        // Wait, the API spec says ws://172.16.17.11:8080/ws/chat/{session_uid}
        // So the base WS URL we'll return is just the host + /ws
        wsUrl = `${wsUrl}/ws`
    }
    
    return {
        token,
        wsUrl
    }
}
