import React from 'react'
import Link from 'next/link'
import { 
  Sprout, 
  Watch, 
  LayoutGrid, 
  Utensils, 
  Laptop, 
  CupSoda, 
  Smartphone,
  Gamepad2,
  ChevronRight
} from "lucide-react"
import { Category } from '@/types/category'
import { getCategories } from '@/lib/api/category'

export const dynamic = "force-dynamic"

const getCategoryIcon = (name: string) => {
  const n = name.toLowerCase();
  if (n.includes("pertanian")) return <Sprout className="w-8 h-8" />;
  if (n.includes("accessories")) return <Watch className="w-8 h-8" />;
  if (n.includes("makanan")) return <Utensils className="w-8 h-8" />;
  if (n.includes("laptop")) return <Laptop className="w-8 h-8" />;
  if (n.includes("minuman")) return <CupSoda className="w-8 h-8" />;
  if (n.includes("smartphone")) return <Smartphone className="w-8 h-8" />;
  if (n.includes("mainan")) return <Gamepad2 className="w-8 h-8" />;
  return <LayoutGrid className="w-8 h-8" />;
}

export default async function CategoryPage() {
  const categories: Category[] = await getCategories()

  return (
    <div className="min-h-screen bg-[#F8FAFC] p-6 md:p-12">
      <div className="max-w-6xl mx-auto">
        
        <div className="mb-10">
          <h1 className="text-3xl font-extrabold text-gray-900 tracking-tight">Jelajahi Kategori</h1>
          <p className="text-gray-500 mt-2">Temukan berbagai produk terbaik berdasarkan minat Anda.</p>
        </div>

        <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-6">
          {categories.map((cat) => (
            <Link 
              key={cat.id} 
              href={`/category/${cat.slug}`}
              className="group"
            >
              <div className="relative bg-white border border-gray-100 rounded-3xl p-6 flex flex-col items-center justify-center transition-all duration-300 shadow-sm hover:shadow-xl hover:-translate-y-2 aspect-square">
                
                <div className="absolute top-0 right-0 p-4 opacity-0 group-hover:opacity-100 transition-opacity">
                  <ChevronRight className="w-5 h-5 text-teal-500" />
                </div>

                <div className="mb-4 p-4 rounded-2xl bg-teal-50 text-teal-600 group-hover:bg-teal-600 group-hover:text-white transition-all duration-300 transform group-hover:rotate-6">
                  {getCategoryIcon(cat.name)}
                </div>

                <div className="text-center">
                  <h3 className="font-bold text-gray-800 text-sm md:text-base group-hover:text-teal-600 transition-colors">
                    {cat.name}
                  </h3>
                  <p className="text-xs text-gray-400 mt-1 font-medium">
                    {cat.product_count} Produk
                  </p>
                </div>
              </div>
            </Link>
          ))}
        </div>

        <div className="mt-16 bg-teal-600 rounded-3xl p-8 text-white flex flex-col md:flex-row items-center justify-between gap-6">
          <div>
            <h2 className="text-2xl font-bold">Butuh bantuan mencari barang?</h2>
            <p className="text-teal-100 opacity-90">Tim kami siap membantu mencarikan produk terbaik untuk Anda.</p>
          </div>
          <button className="bg-white text-teal-600 px-6 py-3 rounded-xl font-bold hover:bg-teal-50 transition-colors whitespace-nowrap">
            Hubungi CS
          </button>
        </div>
      </div>
    </div>
  )
}