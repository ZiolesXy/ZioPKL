import UserLayout from "./(user)/layout"
import ProductCard from "@/components/product/ProductCard"
import { getProducts } from "@/lib/api/product"
import FeaturedCarousel from "@/components/product/FeaturedCarousel"
import { getCategories } from "@/lib/api/category"
import { CategoryCard } from "@/components/shared/CardCategory"
import Image from "next/image"
import { Sparkles, LayoutGrid, ShoppingBag } from "lucide-react"
import { SectionHeader } from "@/components/shared/SectionHeader"
import { Badge } from "@/components/ui/badge"
import Link from "next/link"

const MAX_DISPLAY_ITEMS = 8

export default async function DashboardPage() {
  const [categories, products] = await Promise.all([
    getCategories(),
    getProducts(),
  ])

  return (
    <UserLayout>
      <div className="min-h-screen bg-slate-50/50">
        <section className="w-full bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900 overflow-hidden relative border-b border-slate-800 shadow-2xl">
          <div className="absolute top-0 left-1/4 w-96 h-96 bg-teal-500/10 rounded-full blur-3xl pointer-events-none"></div>
          <div className="max-w-7xl mx-auto px-4 md:px-8 py-20 md:py-32 relative z-10">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-12 items-center">

              <div className="space-y-8 text-white order-2 md:order-1 text-center md:text-left animate-in fade-in slide-in-from-left-8 duration-1000">
                <div className="space-y-4">
                  <Badge variant="outline" className="text-teal-400 border-teal-500/30 bg-teal-500/10 mb-4 px-3 py-1 font-medium tracking-wide">PILIHAN EDITOR</Badge>
                  <h1 className="text-5xl md:text-7xl font-black tracking-tight leading-tight">
                    Acer Nitro <span className="text-transparent bg-clip-text bg-linear-to-r from-teal-400 to-teal-200">Lite 16</span>
                  </h1>
                  <p className="text-slate-300 text-lg md:text-xl md:leading-relaxed max-w-lg mx-auto md:mx-0 font-medium">
                    Performance Gaming yang Terjangkau. Rasakan kecepatan tanpa batas dengan teknologi pendingin terbaru.
                  </p>
                </div>

                <div className="flex flex-wrap gap-4 justify-center md:justify-start">
                  <button className="px-8 py-4 bg-teal-500 text-white font-bold rounded-xl hover:bg-teal-400 transition-all shadow-lg hover:shadow-teal-500/25 hover:-translate-y-1">
                    <Link href="/">Beli Sekarang</Link>
                  </button>
                  <button className="px-8 py-4 bg-slate-800/50 text-slate-200 font-bold rounded-xl border border-slate-700 hover:bg-slate-800 transition-all hover:-translate-y-1">
                    Selengkapnya
                  </button>
                </div>
              </div>

              <div className="order-1 md:order-2 flex justify-center items-center animate-in fade-in zoom-in-95 duration-1000 md:delay-150">
                <div className="relative group animate-bounce-slow">
                  <div className="absolute inset-0 bg-teal-400 blur-[100px] opacity-20 rounded-full"></div>

                  <Image
                    src="/acer.png"
                    alt="Acer Nitro Lite 16"
                    width={800}
                    height={800}
                    className="relative z-10 drop-shadow-2xl object-contain transform group-hover:scale-105 transition-transform duration-500"
                    priority
                  />
                </div>
              </div>

            </div>
          </div>
        </section>


        <div className="max-w-7xl mx-auto p-4 md:p-8 space-y-12">
          <section className="rounded-3xl bg-teal-50/50 border border-teal-100 shadow-sm p-8" >
            <div className="flex items-center gap-2 mb-6">
              <div className="p-2 bg-amber-100 rounded-lg text-amber-600">
                <Sparkles size={20} />
              </div>
              <h2 className="text-2xl md:text-3xl font-extrabold text-slate-800">
                Hot Produk
              </h2>
            </div>
            <div className="">
              <FeaturedCarousel products={products.slice(0, MAX_DISPLAY_ITEMS)} />
            </div>
          </section>

          <section>
            <SectionHeader
              icon={LayoutGrid}
              label="Kategori"
              title="Jelajahi Kategori"
              href="/category"
              linkText="Lihat Semua"
            />

            <div className="flex flex-nowrap overflow-x-auto gap-4 md:gap-6 pb-4 no-scrollbar -mx-4 px-4 md:mx-0 md:px-0">
              {categories.slice(0, MAX_DISPLAY_ITEMS).map((category) => (
                <CategoryCard
                  key={category.id}
                  icon_url={category.icon_url}
                  name={category.name}
                />
              ))}
            </div>
          </section>

          <section>
            <SectionHeader
              icon={ShoppingBag}
              label="Untuk Anda"
              title="Rekomendasi Produk"
              href="/product"
              linkText="Semua Produk"
            />

            <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4 md:gap-8">
              {products.slice(0, MAX_DISPLAY_ITEMS).map((item) => (
                <div key={item.id} className="hover:-translate-y-1 transition-transform duration-300">
                  <ProductCard
                    id={item.id}
                    slug={item.slug}
                    name={item.name}
                    price={item.price}
                    category={item.category}
                    stock={item.stock}
                    image={item.image_url}
                  />
                </div>
              ))}
            </div>
          </section>

        </div>
      </div>
    </UserLayout>
  )
}
