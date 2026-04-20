"use client"

import React, { useRef, useEffect, useState, useMemo } from "react"
import { ChevronLeft, ChevronRight } from "lucide-react"
import ProductCard from "./ProductCard"
import { Product } from "@/types/product"

const CARD_WIDTH = 320
const CARD_GAP = 20
const SCROLL_AMOUNT = CARD_WIDTH + CARD_GAP
const AUTO_SCROLL_INTERVAL_MS = 2500

export default function FeaturedCarousel({
  products,
}: {
  products: Product[]
}) {
  const scrollRef = useRef<HTMLDivElement>(null)
  const [isHovered, setIsHovered] = useState(false)

  // Duplikasi data untuk efek infinite
  const displayProducts = useMemo(
    () => [...products, ...products, ...products],
    [products]
  )

  const handleInfiniteScroll = () => {
    const el = scrollRef.current
    if (!el) return

    // Hitung lebar satu set produk asli
    const singleSetWidth = el.scrollWidth / 3

    // Jika user scroll terlalu ke kanan, kembalikan ke tengah
    if (el.scrollLeft >= singleSetWidth * 2) {
      el.style.scrollBehavior = "auto" // Matikan smooth sementara
      el.scrollLeft -= singleSetWidth
    }
    // Jika user scroll terlalu ke kiri, kembalikan ke tengah
    else if (el.scrollLeft <= 0) {
      el.style.scrollBehavior = "auto"
      el.scrollLeft += singleSetWidth
    }
  }

  const scroll = (dir: "left" | "right") => {
    const el = scrollRef.current
    if (!el) return

    el.style.scrollBehavior = "smooth" // Pastikan smooth aktif saat tombol ditekan
    el.scrollLeft += dir === "left" ? -SCROLL_AMOUNT : SCROLL_AMOUNT
  }

  useEffect(() => {
    const el = scrollRef.current
    if (!el) return

    // Set posisi awal ke set produk kedua (tengah) agar bisa scroll ke kiri/kanan
    const singleSetWidth = el.scrollWidth / 3
    el.scrollLeft = singleSetWidth
  }, [products])

  useEffect(() => {
    if (isHovered) return

    const interval = setInterval(() => {
      const el = scrollRef.current
      if (el) {
        el.style.scrollBehavior = "smooth"
        el.scrollLeft += SCROLL_AMOUNT
      }
    }, AUTO_SCROLL_INTERVAL_MS)

    return () => clearInterval(interval)
  }, [isHovered])

  return (
    <div
      className="relative group w-full px-4"
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
    >
      {/* Tombol Navigasi */}
      <button
        onClick={() => scroll("left")}
        className="absolute left-0 top-1/2 -translate-y-1/2 z-10 bg-white shadow-xl rounded-full p-2 border opacity-0 group-hover:opacity-100 transition hover:bg-gray-50"
      >
        <ChevronLeft className="w-6 h-6" />
      </button>

      <button
        onClick={() => scroll("right")}
        className="absolute right-0 top-1/2 -translate-y-1/2 z-10 bg-white shadow-xl rounded-full p-2 border opacity-0 group-hover:opacity-100 transition hover:bg-gray-50"
      >
        <ChevronRight className="w-6 h-6" />
      </button>

      {/* Container Carousel */}
      <div
        ref={scrollRef}
        onScroll={handleInfiniteScroll}
        className="flex overflow-x-auto gap-5 pb-4 no-scrollbar"
        style={{ scrollSnapType: "x mandatory" }}
      >
        {displayProducts.map((item, index) => (
          <div
            key={`${item.id}-${index}`}
            className="flex-none w-[280px] sm:w-[320px]"
            style={{ scrollSnapAlign: "start" }}
          >
            <ProductCard {...item} image={item.image_url} />
          </div>
        ))}
      </div>
    </div>
  )
}