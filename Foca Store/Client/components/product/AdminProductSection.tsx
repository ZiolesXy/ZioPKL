"use client"

import { useState, useMemo } from "react"
import ProductTable from "@/components/product/ProductTable"
import ProductSearch from "@/components/product/ProductSearch"
import PaginationControls from "@/components/admin/PaginationControl"
import { Product } from "@/types/product"

export default function AdminProductClient({ products }: { products: Product[] }) {
  const [query, setQuery] = useState("")
  const [page, setPage] = useState(1)

  const limit = 10

  const filtered = useMemo(() => {
    const q = query.toLowerCase()

    return products.filter((p) =>
      p.name.toLowerCase().includes(q) ||
      p.id.toString().includes(q)
    )
  }, [products, query])

  const totalPages = Math.ceil(filtered.length / limit)

  const paginated = useMemo(() => {
    const start = (page - 1) * limit
    return filtered.slice(start, start + limit)
  }, [filtered, page])

  function handleSearch(q: string) {
    setQuery(q)
    setPage(1)
  }

  return (
    <div className="space-y-5">
      <h1 className="text-2xl font-extrabold text-slate-900 tracking-tight">Product Management</h1>
      <p className="text-slate-500 mb-5">Utility untuk mengelola produk.</p>
      <ProductSearch onSearch={handleSearch} />

      <ProductTable products={paginated} />

      <PaginationControls
        page={page}
        totalPages={totalPages}
        onPageChange={setPage}
      />
    </div>
  )
}
