import Link from "next/link"
import { ArrowRight } from "lucide-react"
import type { LucideIcon } from "lucide-react"

interface SectionHeaderProps {
  icon: LucideIcon
  label: string
  title: string
  href: string
  linkText: string
}

export function SectionHeader({
  icon: Icon,
  label,
  title,
  href,
  linkText,
}: SectionHeaderProps) {
  return (
    <div className="flex items-end justify-between mb-6">
      <div className="space-y-1">
        <div className="flex items-center gap-2 text-teal-600 font-bold uppercase tracking-wider text-xs">
          <Icon size={14} />
          <span>{label}</span>
        </div>
        <h2 className="text-2xl font-bold text-slate-800">{title}</h2>
      </div>

      <Link
        href={href}
        className="group flex items-center gap-2 text-sm font-bold text-teal-600 hover:text-teal-700 transition-all"
      >
        {linkText}
        <ArrowRight
          size={16}
          className="group-hover:translate-x-1 transition-transform"
        />
      </Link>
    </div>
  )
}
