import api from "../axios"
import { handleApiError } from "../utils"
import { Category, CategoryWithProducts } from "@/types/category"
import { clientApi } from "../client-api"

export async function getCategories(): Promise<Category[]> {
  try {
    const response = await api.get<{ data: { entries: Category[] } }>("/category")
    return response.data?.data?.entries || []
  } catch (error: unknown) {
    handleApiError(error)
    return[]
  }
}

export async function getCategoryBySlug(slug: string): Promise<CategoryWithProducts | undefined> {
  try {
    const response = await api.get<{ data: CategoryWithProducts }>(`/category/${slug}`)
    return response.data?.data
  } catch (error: unknown) {
    handleApiError(error)
  }
}

export async function addCategory(data: { name: string; icon: File }) {
  const formData = new FormData()
  formData.append("name", data.name)
  formData.append("icon", data.icon)

  try {
    const response = await clientApi.post<{ data: Category }>(
      "/api/admin/category",
      formData
    )

    return response.data
  } catch (error: unknown) {
    handleApiError(error)
  }
}

export async function deleteCategory(id: number) {
  try {
    const response = await clientApi.delete<{ data: Category }>(`/api/admin/category/${id}`)
    return response.data
  } catch (error: unknown) {
    handleApiError(error)
  }
}

export async function editCategory(id: number, formData: FormData) {
  try {
    const response = await clientApi.put<{ data: Category }>(`/api/admin/category/${id}`, formData)
    return response.data
  } catch (error: unknown) {
    handleApiError(error)
  }
}