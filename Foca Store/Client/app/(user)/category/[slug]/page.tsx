import { getCategoryBySlug } from "@/lib/api/category"
import ProductCard from "@/components/product/ProductCard"
import { notFound } from "next/navigation"
import { Product } from "@/types/product"

export default async function CategoryPage({ params }: { params: Promise<{ slug: string }> }) {
  const { slug } = await params

  const category = await getCategoryBySlug(slug)

  if (!category) {
    notFound()
  }

  return (
    <main className="min-h-screen bg-[#F8FAFC]">
      <div className="max-w-350 mx-auto p-6">
        <div className="mb-8 border-b pb-4">
          <h1 className="text-3xl font-extrabold text-gray-900 capitalize">
            Kategori: {category.name}
          </h1>
          <p className="text-gray-500 mt-2">
            Menampilkan {category.products?.length || 0} produk pilihan
          </p>
        </div>

        <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
          {category.products?.map((item: Product) => (
            <ProductCard
              key={item.id}
              {...item}
              image={item.image_url}
            />
          ))}
        </div>

        {(!category.products || category.products.length === 0) && (
          <div className="text-center py-20 bg-white rounded-xl border-dashed border-2">
            <p className="text-gray-400 font-medium">Belum ada produk di kategori ini.</p>
          </div>
        )}
      </div>
    </main>
  )
}
