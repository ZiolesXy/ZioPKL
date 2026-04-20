  import { clientApi } from "../client-api"
import api from "../axios"
import { handleApiError } from "../utils"
import { AvailableCoupon, MyCoupon } from "@/types/coupon"

type CreateCouponRequest = {
  code: string
  type: "percentage" | "fixed" | string
  value: number
  quota: number
  is_active: boolean
  expire_at: string
  minimum_purchase?: number
}

type UpdateCouponRequest = Partial<CreateCouponRequest>

export async function getCoupon(): Promise<AvailableCoupon[]> {
  try {
    const response = await api.get<{ data?: { entries?: AvailableCoupon[] }; entries?: AvailableCoupon[] }>("/coupons")
    return response.data?.data?.entries ?? response.data?.entries ?? []
  } catch (error: unknown) {
    handleApiError(error)
    return []
  }
}

export async function claimCoupon(id: number) {
  try {
    const response = await clientApi.post<{ data?: unknown }>(`/api/coupons/${id}/claim`)
    return response?.data ?? response
  } catch (error: unknown) {
    handleApiError(error)
  }
}

export async function getMyCoupon(): Promise<MyCoupon[]> {
  try {
    const response = await clientApi.get<{ data?: { entries?: MyCoupon[] }; entries?: MyCoupon[] }>("/api/coupons/me")
    return response?.data?.entries ?? response?.entries ?? []
  } catch (error: unknown) {
    handleApiError(error)
    return []
  }
}

export async function deleteCoupon(id: number) {
  try {
    const response = await clientApi.delete<{ data?: unknown }>(`/api/admin/coupon/${id}`)
    return response?.data ?? response
  } catch (error: unknown) {
    handleApiError(error)
  }
}

export async function editCoupon(id: number, payload: Pick<AvailableCoupon, "code" | "type" | "value" | "quota" | "minimum_purchase" | "expires_at" | "is_active">) {
  try {
    const requestPayload: UpdateCouponRequest = {
      code: payload.code,
      type: payload.type,
      value: payload.value,
      quota: payload.quota,
      minimum_purchase: payload.minimum_purchase,
      expire_at: payload.expires_at,
      is_active: payload.is_active,
    }
    const response = await clientApi.put<{ data?: unknown }>(`/api/admin/coupon/${id}`, requestPayload)
    return response?.data ?? response
  } catch (error: unknown) {
    handleApiError(error)
  } 
}

export async function addCoupon(
  payload: Pick<AvailableCoupon, "code" | "type" | "value" | "quota" | "minimum_purchase" | "expires_at" | "is_active">
) {
  try {
    const requestPayload: CreateCouponRequest = {
      code: payload.code,
      type: payload.type,
      value: payload.value,
      quota: payload.quota,
      is_active: payload.is_active,
      minimum_purchase: payload.minimum_purchase,
      expire_at: payload.expires_at,
    }
    const response = await clientApi.post<{ data?: unknown }>(
      "/api/admin/coupons",
      requestPayload
    )
    return response?.data ?? response
  } catch (error: unknown) {
    handleApiError(error)
  }
}
