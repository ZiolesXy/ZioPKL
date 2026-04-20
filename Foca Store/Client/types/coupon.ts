export interface AvailableCoupon {
  id: number
  code: string
  type: "percentage" | "fixed"
  value: number
  quota: number
  is_active: boolean
  minimum_purchase: number
  used_count: number
  remaining: number
  expires_at: string
}

export interface MyCoupon {
  id: number
  coupon_code: string
  coupon_type: "percentage" | "fixed"
  value: number
  used_at: string | null
  claimed_at: string
}