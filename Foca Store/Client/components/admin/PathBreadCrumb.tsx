"use client"

import { usePathname } from "next/navigation"
import { ChevronRight } from "lucide-react"

// 1. Definisikan Dictionary di luar komponen
const routeLabels: Record<string, string> = {
  "/overview": "Overview",
  "/check-products": "Cek Produk",
  "/add-products": "Tambah Produk",
  "/add-category": "Tambah Kategori",
  "/system": "System",
  "/add-coupon": "Tambah Kupon",
  "/chat": "Customer Service"
}

export function DynamicBreadcrumb() {
  const pathname = usePathname()

  const getLabel = () => {
    if (routeLabels[pathname]) {
      return routeLabels[pathname]
    }

    const segments = pathname.split('/').filter(Boolean)
    const lastSegment = segments[segments.length - 1] || "Overview"
    
    return lastSegment
      .charAt(0).toUpperCase() + 
      lastSegment.slice(1).replace(/-/g, ' ')
  }

  return (
    <div className="flex items-center gap-2">
      <ChevronRight className="h-4 w-4 text-muted-foreground/50" />
      <span className="text-foreground transition-all">
        {getLabel()}
      </span>
    </div>
  )
}