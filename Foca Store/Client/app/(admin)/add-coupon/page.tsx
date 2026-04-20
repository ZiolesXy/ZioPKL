
import { getCoupon} from "@/lib/api/coupon"
import CouponTable from "@/components/admin/CouponTable"

export default async function AdminCouponPage() {
  const coupons = await getCoupon()
  return (
    <CouponTable coupons={coupons} />
  )
}
