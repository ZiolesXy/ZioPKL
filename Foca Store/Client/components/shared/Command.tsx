"use client"

import {
  Command,
  CommandInput,
  CommandList,
  CommandItem,
  CommandEmpty
} from "@/components/ui/command"

import { useState, useRef, useEffect } from "react"
import { useRouter } from "next/navigation"
import { Product } from "@/types/product"
export function CommandDemo({ products }: { products: Product[] }) {
  const [query, setQuery] = useState("")
  const [open, setOpen] = useState(false)
  const router = useRouter()
  const wrapperRef = useRef<HTMLDivElement>(null)

  const filtered = products.filter((p) =>
    p.name.toLowerCase().includes(query.toLowerCase())
  )

  useEffect(() => {
    function handleClickOutside(e: MouseEvent) {
      if (wrapperRef.current && !wrapperRef.current.contains(e.target as Node)) {
        setOpen(false)
      }
    }
    document.addEventListener("mousedown", handleClickOutside)
    return () => document.removeEventListener("mousedown", handleClickOutside)
  }, [])

  return (
    <div ref={wrapperRef} className="relative w-75">
      <Command className="rounded border bg-whiteoverflow-visible">
        <CommandInput
          placeholder="Cari produk..."
          value={query}
          onValueChange={(val) => {
            setQuery(val)
            setOpen(val.length > 0)
          }}
          onFocus={() => {
            if (query.length > 0) setOpen(true)
          }}
        />

        {open && (
          <CommandList className="absolute top-full left-0 right-0 z-50 mt-1 max-h-[300px] overflow-y-auto rounded-md border bg-white shadow-lg">
            <CommandEmpty>Tidak ditemukan.</CommandEmpty>

            {filtered.map((item) => (
              <CommandItem
                key={item.id}
                onSelect={() => {
                  router.push(`/product/${item.slug}`)
                  setQuery("")
                  setOpen(false)
                }}
                className="cursor-pointer px-3 py-2 hover:bg-slate-100"
              >
                {item.name}
              </CommandItem>
            ))}
          </CommandList>
        )}
      </Command>
    </div>
  )
}
