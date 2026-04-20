import { getProducts } from '@/lib/api/product'
import AdminProductSection from '@/components/product/AdminProductSection'

async function AdminProductPage() {
  const products = await getProducts()
  

  return (
    <div className="p-8 space-y-6">
      <AdminProductSection products={products} />
    </div>
  )
}
export default AdminProductPage
export const dynamic = "force-dynamic"

