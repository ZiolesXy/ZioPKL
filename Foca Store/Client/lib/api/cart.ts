import { clientApi } from "../client-api"
import { handleApiError } from "../utils"
import { Cart } from "@/types/cart"

export async function addToCart(productId: number, quantity: number) {
  try {
    const response = await clientApi.post("/api/cart/items", {
      product_id: productId,
      quantity: quantity
    })
    return response
  } catch (error: unknown) {
    handleApiError(error)
  }
}

export async function getCart(): Promise<Cart | undefined> {
  try {
    const response = await clientApi.get<{ data: Cart }>("/api/cart")
    return response.data
  } catch (error: unknown) {
    handleApiError(error)
  }
}

export async function deleteCartItem(itemId: number) {
  try {
    const response = await clientApi.delete(`/api/cart/items/${itemId}`)
    return response
  } catch (error: unknown) {
    handleApiError(error)
  }
}

export async function checkoutCart(itemIDs: number[], addressUid: string, couponCode?: string) {
  try {
    const response = await clientApi.post("/api/checkout", {
      cart_item_ids: itemIDs,
      address_uid: addressUid,
      coupon_code: couponCode
    })
    return response
  } catch (error: unknown) {
    handleApiError(error)
  }
}

export async function deleteCartItemsMany(itemIDs: number[]) {
  try {
    const response = await clientApi.delete("/api/cart/items", {
      cart_item_ids: itemIDs,
    })
    return response
  } catch (error: unknown) {
    handleApiError(error)
  }
}
