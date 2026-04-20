"use client"

import { useState} from "react"
import { Search } from "lucide-react"
import { Input } from "@/components/ui/input"

export default function ProductSearch({onSearch} : {onSearch: (query: string) => void}) {
    const [query, setQuery] = useState("")
    function handleChange(e: React.ChangeEvent<HTMLInputElement>) {
    const value = e.target.value
    setQuery(value)
    onSearch(value)
  }

  return (
    <div className="relative flex-1 max-w-sm">
      <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
      <Input
        placeholder="Cari nama atau ID produk..."
        className="pl-8"
        value={query}
        onChange={handleChange}
      />
    </div>
  )
}
