"use client"
import Link from "next/link"
import { LayoutGrid } from "lucide-react"
import Image from "next/image"

export function CategoryCard({
  name,
  icon_url,
}: {
  name: string
  icon_url?: string
}) {
  return (
    <Link
      href={`/category/${name.toLowerCase()}`}
      className="group flex flex-col items-center justify-center 
                 w-32 h-32 md:w-40 md:h-40 p-4
                 bg-white border border-slate-200 rounded-3xl shadow-sm 
                 hover:shadow-xl hover:border-teal-500 hover:-translate-y-2 
                 transition-all duration-300 ease-out"
    >
      <div className="mb-3 p-3 rounded-2xl bg-slate-50 group-hover:bg-teal-50 transition-colors duration-300">
        {icon_url ? (
          <Image
            src={icon_url}
            alt={name}
            width={40}
            height={40}
            className="object-contain"
          />
        ) : (
          <LayoutGrid
            className="w-8 h-8 md:w-10 md:h-10 text-slate-600 group-hover:text-teal-600 transition-colors"
            strokeWidth={1.5}
          />
        )}
      </div>

      <span className="text-xs md:text-sm font-bold text-slate-700 group-hover:text-teal-600 text-center line-clamp-1 transition-colors">
        {name}
      </span>
    </Link>
  )
}