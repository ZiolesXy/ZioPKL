"use server"

import { revalidatePath } from "next/cache"
import { updateProduct, deleteProduct } from "@/lib/api/product"
import { getErrorMessage } from "@/lib/utils"

export async function editProductAction(id: number, formData: FormData) {
  try {
    await updateProduct(id, formData)

    revalidatePath("/overview")
    return { success: true }
  } catch (error: unknown) {
    return {
      success: false,
      message: getErrorMessage(error, "Gagal update produk"),
    }
  }
}

export async function deleteProductAction(id: number) {
  try {
    await deleteProduct(id)

    revalidatePath("/overview")
    return { success: true }
  } catch (error: unknown) {
    return {
      success: false,
      message: getErrorMessage(error, "Gagal hapus produk"),
    }
  }
}