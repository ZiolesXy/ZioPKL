import { Badge } from "@/components/ui/badge"
import { getProductsBySlug } from "@/lib/api/product"
import { getCategories } from "@/lib/api/category"
import { Star, ShieldCheck, Tag, ChevronRight } from "lucide-react"
import Image from "next/image"
import { ProductActions } from "@/components/product/ProductActions" 
import Link from "next/link"
import { formatRupiah } from "@/lib/utils"

function slugifyCategory(value: string): string {
  return value
    .trim()
    .toLowerCase()
    .replace(/[^a-z0-9\s-]/g, "")
    .replace(/\s+/g, "-")
    .replace(/-+/g, "-")
}

export default async function ProductDetailPage({
  params,
}: {
  params: Promise<{ slug: string }>
}) {
  const { slug } = await params


  const product = await getProductsBySlug(slug)

  const categoryName =
    product && typeof product.category === "object" && product.category
      ? product.category.name
      : product && typeof product.category === "string"
        ? product.category
        : null

  let categorySlug: string | null =
    product && typeof product.category === "object" && product.category
      ? product.category.slug
      : null

  if (!categorySlug && categoryName) {
    const categories = await getCategories()
    const normalized = categoryName.trim().toLowerCase()
    const matched = categories.find(
      (c) => c.slug.toLowerCase() === normalized || c.name.toLowerCase() === normalized
    )
    categorySlug = matched?.slug ?? slugifyCategory(categoryName)
  }

  if (!product) return <div>Produk tidak ditemukan</div>

  return (
    <div className="container mx-auto p-6 max-w-6xl">
      <nav className="flex items-center gap-2 text-sm text-muted-foreground mb-6">
        <Link href="/" className="hover:text-teal-600 transition">Home</Link>
        <ChevronRight className="size-4" />
        <Link 
          href="/category" 
          className="hover:text-teal-600 transition"
        >
          Category
        </Link>
        <ChevronRight className="size-4" />
        <Link 
          href={categorySlug ? `/category/${encodeURIComponent(categorySlug)}` : "/category"} 
          className="hover:text-teal-600 transition"
        >
          {typeof product.category === 'object' ? product.category.name : product.category}
        </Link>

        <ChevronRight className="size-4" />
        <span className="text-gray-900 font-medium truncate">{product.name}</span>
      </nav>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-10">
        {/* Gambar Produk */}
        <div className="space-y-4">
          <div className="aspect-square rounded-xl overflow-hidden border bg-white shadow-sm p-8">
            <Image
              src={product.image_url}
              alt={product.name}
              width={600}
              height={600}
              className="w-full h-full object-contain"
              priority
            />
          </div>
       
        </div>

        {/* Info Produk */}
        <div className="space-y-6">
          <div>
            {/* 2. CATEGORY LABEL - Di atas judul */}
            <div className="flex items-center gap-2 mb-2">
                <Badge variant="secondary" className="bg-teal-50 text-teal-700 hover:bg-teal-100 border-none">
                   <Tag className="size-3 mr-1" />
                   {typeof product.category === 'object' ? product.category.name : product.category}
                </Badge>
                <Badge variant="outline">Rekomendasi</Badge>
            </div>
            
            <h1 className="text-4xl font-bold text-gray-900 leading-tight">{product.name}</h1>
            
            <div className="flex items-center gap-2 mt-2 text-yellow-500">
               <Star className="fill-current size-4" />
               <span className="text-sm text-muted-foreground font-medium">4.8 (120 Ulasan)</span>
            </div>
          </div>
          <div className="border-y py-4">
            <p className="text-3xl font-bold text-teal-600">{formatRupiah(product.price)}</p>
            <p className={`text-sm mt-1 ${product.stock > 0 ? 'text-muted-foreground' : 'text-red-500 font-bold'}`}>
              {product.stock > 0 ? `${product.stock} Unit Tersedia` : 'Stok Habis'}
            </p>
          </div>

          <div className="space-y-4">
            <h3 className="font-semibold text-lg">Deskripsi Produk</h3>
            <p className="text-gray-600 leading-relaxed whitespace-pre-line">
              {product.description}
            </p>
          </div>

          <ProductActions productId={product.id} stock={product.stock} />

          <div className="pt-6 border-t grid grid-cols-2 gap-4 text-sm text-muted-foreground">
            <div className="flex items-center gap-2">
              <ShieldCheck className="text-teal-600 size-5" /> Garansi Resmi 2 Tahun
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
