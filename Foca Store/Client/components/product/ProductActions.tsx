"use client"

import { useState } from "react"
import { Button } from "@/components/ui/button"
import { ShoppingCart, Loader2 } from "lucide-react"
import { useCart } from "@/context/CartContext"
import { notifyError } from "@/lib/toast"
import Link from "next/link"

interface ProductActionsProps {
  productId: number;
  stock: number;
}

export function ProductActions({ productId, stock }: ProductActionsProps) {
  const { addToCart } = useCart()
  const [isLoading, setIsLoading] = useState(false)

  async function handleAdd() {
    if (stock <= 0) return

    setIsLoading(true)
    try {
      await addToCart(productId, 1)
    } catch (err) {
      console.error("Gagal menambahkan ke keranjang:", err)
      notifyError("Gagal menambahkan ke keranjang. Silakan coba lagi.")
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className="flex gap-4 pt-6">
      <Button
        size="lg"
        onClick={handleAdd}
        disabled={stock <= 0 || isLoading}
        className="flex-1 bg-teal-600 hover:bg-teal-700 h-14 text-lg gap-2"
      >
        {isLoading ? (
          <Loader2 className="size-5 animate-spin" />
        ) : (
          <ShoppingCart className="size-5" />
        )}
        {isLoading ? "Menambahkan..." : "Tambah ke Keranjang"}
      </Button>

      <Button
        size="lg"
        variant="outline"
        className="h-14 px-8"
        disabled={stock <= 0}
      >
        <Link href="/">Beli Sekarang</Link>
      </Button>
    </div>
  )
}