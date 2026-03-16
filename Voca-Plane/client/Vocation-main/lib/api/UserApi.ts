// lib/api-user.ts
import axiosInstance from "../axios";
import { getAuthHeaders } from "../getAuth";

export async function getProfile() {
    const isClient = typeof window !== "undefined";

    if (isClient) {
        const Cookies = (await import("js-cookie")).default;
        const token = Cookies.get("access_token") || Cookies.get("token");
        const headers = token ? { Authorization: `Bearer ${token}` } : {};
        const response = await axiosInstance.get("/user/profile", { headers });
        return response.data;
    }

    const headers = await getAuthHeaders();
    const response = await axiosInstance.get("/user/profile", { headers });
    return response.data; // Biasanya berisi { success: true, data: { name, email, ... } }
}

