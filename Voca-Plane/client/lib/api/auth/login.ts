// lib/api-auth.ts
import axiosInstance from "@/lib/axios";

export async function Login(email: string, password: string) {
    try {
        const response = await axiosInstance.post('/auth/login', {
            email,
            password
        });
        return response.data; 
    } catch (error: any) {
        throw error.response?.data?.message || "Login failed";
    }
}