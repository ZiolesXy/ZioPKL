import api from "../axios"
import { clientApi } from "../client-api"
import { getAuthHeaders } from "../auth-server"
import { handleApiError } from "../utils"
import { User, ChangePasswordData, AddressData, CreateAddressData } from "@/types/auth"
import type { GetMyCheckoutResponse, GetCheckoutDetailResponse } from "@/types/checkout"

// Server-side: get user profile (untuk server components / server actions)
export async function getUser(): Promise<User | undefined> {
  try {
    const headers = await getAuthHeaders()
    const response = await api.get("/api/profile", { headers })
    return response.data.data
  } catch {
    return undefined
  }
}

// Client-side: get user profile (untuk hooks/client components)
export async function getUserClient(): Promise<User | undefined> {
  try {
    const response = await clientApi.get<{ data: User }>("/api/profile")
    return response.data
  } catch (error: unknown) {
    handleApiError(error)
  }
}

export async function forgotPassword(email: string) {
  try {
    const response = await clientApi.post("/forgot-password", { email })
    return response
  } catch (error: unknown) {
    handleApiError(error)
  }
}

export async function resetPassword(data: { email: string; otp: string; new_password: string }) {
  try {
    const response = await clientApi.post("/verify-otp", data)
    return response
  } catch (error: unknown) {
    handleApiError(error)
  }
}

export async function updateUser(data: FormData) {
  try {
    const response = await clientApi.put("/api/profile", data)
    return response
  } catch (error: unknown) {
    handleApiError(error)
  }
}



export async function getMyCheckouts() {
  try {
    const response =
      await clientApi.get<GetMyCheckoutResponse>("/api/checkout/me")

    return response.data.entries
  } catch (error: unknown) {
    handleApiError(error)
  }
}

export async function getCheckoutDetail(uid: string) {
  try {
    const response =
      await clientApi.get<GetCheckoutDetailResponse>(`/api/checkout/${uid}`)

    return response.data
  } catch (error: unknown) {
    handleApiError(error)
  }
}
export async function changePassword(data: ChangePasswordData) {
  try {
    const response = await clientApi.put("/api/change-password", data)
    return response
  } catch (error: unknown) {
    handleApiError(error)
  }
}

export async function getAddresses(): Promise<AddressData[]> {
  try {
    const response = await clientApi.get<{ data: { entries: AddressData[] } }>("/api/addresses")
    return response.data.entries
  } catch (error: unknown) {
    handleApiError(error)
    return []
  }
}

export async function addAddresses(data: CreateAddressData) {
  try {
    const response = await clientApi.post("/api/addresses", data)
    return response
  } catch (error: unknown) {
    handleApiError(error)
  }
}

export async function updateAddress(uid: string, data: Partial<AddressData>) {
  try {
    const response = await clientApi.put(`/api/addresses/${uid}`, data)
    return response
  } catch (error: unknown) {
    handleApiError(error)
  }
}

export async function deleteAddress(uid: string) {
  try {
    const response = await clientApi.delete(`/api/addresses/${uid}`)
    return response
  } catch (error: unknown) {
    handleApiError(error)
  }
}

