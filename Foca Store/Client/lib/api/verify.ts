import { clientApi } from "../client-api"
import { handleApiError } from "../utils"

// Note: getCheckout returns admin order entries (ApiOrder shape defined in check-products/page.tsx)
export async function getCheckout() {
    try {
        const response = await clientApi.get<{ data: { entries: unknown[] } }>("/api/admin/checkout")
        return response.data.entries
    } catch (error) {
        handleApiError(error)
    }
}

export async function verifyCartApproved(id: number) {
    try {
        const response = await clientApi.patch(`/api/admin/checkout/${id}/approve`)
        return response
    } catch (error) {
        handleApiError(error)
    }
}

export async function verifyCartRejected(id: number) {
    try {
        const response = await clientApi.patch(`/api/admin/checkout/${id}/reject`)
        return response
    } catch (error) {
        handleApiError(error)
    }
}