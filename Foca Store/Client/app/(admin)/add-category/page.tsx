import { getCategories } from "@/lib/api/category"
import CategoryTable from "@/components/admin/CategoryTable"

export default async function AdminCategoryPage() {
  const categories = await getCategories()
  return (
    <CategoryTable categories={categories} />
  )
}
