import { MyCoupon } from "./coupon"

export type CheckoutStatus = "pending" | "approved" | "rejected"

export interface CheckoutUser {
  id: number
  name: string
  email: string
}

export interface CheckoutProduct {
  id: number
  name: string
}

export interface CheckoutItem {
  id: number
  quantity: number
  price: number
  product: CheckoutProduct
}

export interface CheckoutNotification {
  id: number
  user: CheckoutUser
  uid: string
  total_price: number
  status: CheckoutStatus
  items: CheckoutItem[]
  created_at: string
  updated_at: string
}

export interface GetMyCheckoutResponse {
  status: string
  message: string
  data: {
    entries: CheckoutNotification[]
  }
}

export interface CheckoutAddress {
  address_uid: string
  recipient_name: string
  phone: string
  full_address: string
}

export interface CheckoutDetailUser {
  id: number
  name: string
  email: string
  telephone_number: string
}

export interface CheckoutDetail {
  id: number
  uid: string
  user: CheckoutDetailUser
  items: CheckoutItem[]
  coupon:   MyCoupon | null
  subtotal: number
  discount_amount: number
  snap_token: string | null
  whatsapp_url: string | null
  total_price: number
  status: CheckoutStatus
  address: CheckoutAddress
  created_at: string
  updated_at: string
}

export interface GetCheckoutDetailResponse {
  status: string
  message: string
  data: CheckoutDetail
}