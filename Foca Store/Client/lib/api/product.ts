import api from "../axios"
import { getAuthHeaders } from "../auth-server"
import { handleApiError } from "../utils"
import { Product } from "@/types/product"

export async function getProducts(): Promise<Product[]> {
    try {
        const headers = await getAuthHeaders()
        const response = await api.get<{ data: { entries: Product[] } }>("/products", { headers })
        return response.data.data.entries ?? []
    } catch (error: unknown) {
        handleApiError(error)
        return []
    }
}

export async function getProductsBySlug(slug: string): Promise<Product | undefined> {
    try {
        const headers = await getAuthHeaders()
        const response = await api.get<{ data: Product }>(`/product/${slug}`, { headers })
        return response.data.data
    } catch (error: unknown) {
        handleApiError(error)
    }
}

export async function updateProduct(id: number, formData: FormData) {
  const headers = await getAuthHeaders()

  const response = await api.put(
    `/api/admin/products/${id}`,
    formData,
    {
      headers: {
        ...headers,
        "Content-Type": "multipart/form-data",
      },
    }
  )

  return response.data
}

export async function deleteProduct(id: number) {
  const headers = await getAuthHeaders()

  const response = await api.delete(
    `/api/admin/products/${id}`,
    { headers }
  )

  return response.data
}

export async function addProduct(formData: FormData) {
    const headers = await getAuthHeaders();
    const response = await api.post("/api/admin/products", formData, {
        headers: {
            ...headers,
        },
    });
    return response.data;
}