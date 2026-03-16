"use client"

import { useEffect, useState, useCallback } from "react"
import Image from "next/image"
import { 
  ChevronLeft, 
  ChevronRight, 
  Edit2, 
  Trash2, 
  Plane, 
  MoreHorizontal, 
  Code,
  Building
} from "lucide-react"

import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { 
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"

import { getAirports, deleteAirport, createAirport, updateAirport } from "@/lib/api/AirportApi"
import { Airport } from "@/lib/type/airport"
import { GenericDeleteButton } from "./GenericDeleteButton"
import { UpsertForm, FormField } from "./UpsertForm"

export function AirportsTableData({ refreshKey }: { refreshKey?: number }) {
  const [airports, setAirports] = useState<Airport[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [currentPage, setCurrentPage] = useState(1)
  const [totalPages, setTotalPages] = useState(1)
  const limit = 10

  const airportFields: FormField[] = [
    { name: "code", label: "Kode Bandara", placeholder: "Contoh: CGK", required: true },
    { name: "name", label: "Nama Bandara", placeholder: "Contoh: Soekarno-Hatta", required: true },
    { name: "city", label: "Kota", placeholder: "Contoh: Jakarta", required: true },
  ]

  const fetchData = useCallback(async () => {
    setIsLoading(true)
    try {
      const response = await getAirports(currentPage, limit)
      setAirports(response.data)
      setTotalPages(Math.ceil(response.meta.total / response.meta.limit))
    } catch (error) {
      console.error("Gagal mengambil data bandara:", error)
    } finally {
      setIsLoading(false)
    }
  }, [currentPage, limit])

  useEffect(() => {
    fetchData()
  }, [fetchData, refreshKey])

  if (isLoading) {
    return (
      <div className="flex h-64 items-center justify-center rounded-xl border bg-card">
        <div className="flex flex-col items-center gap-2">
          <div className="h-8 w-8 animate-spin rounded-full border-4 border-primary border-t-transparent" />
          <p className="text-sm text-muted-foreground font-medium">Memuat data bandara...</p>
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-4">
      <div className="rounded-xl border bg-card shadow-sm overflow-hidden">
        <Table>
          <TableHeader className="bg-muted/50">
            <TableRow>
              <TableHead className="w-[80px] text-center">ID</TableHead>
              <TableHead>Kode Bandara</TableHead>
              <TableHead>Nama Bandara</TableHead>
              <TableHead>Kota</TableHead>
              <TableHead className="text-right">Aksi</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {airports.length > 0 ? (
              airports.map((airport) => (
                <TableRow key={airport.id} className="group hover:bg-muted/30 transition-colors">
                  <TableCell className="text-center font-mono text-xs text-muted-foreground">
                    {airport.id}
                  </TableCell>
                  <TableCell>
                    <Badge variant="secondary" className="font-mono tracking-wider">
                      <Code className="mr-1 h-3 w-3" />
                      {airport.code}
                    </Badge>
                  </TableCell>
                  <TableCell>
                    <Badge variant="secondary" className="font-mono tracking-wider">
                      <Plane className="mr-1 h-3 w-3" />
                      {airport.name}
                    </Badge>
                  </TableCell>
                  <TableCell>
                    <Badge variant="secondary" className="font-mono tracking-wider">
                      <Building className="mr-1 h-3 w-3" />
                      {airport.city}
                    </Badge>
                  </TableCell>  
                  <TableCell className="text-right">
                    <div className="hidden md:flex justify-end gap-2">
                      <UpsertForm
                        title="Bandara"
                        description="Pastikan data bandara sudah benar sebelum menyimpan perubahan."
                        fields={airportFields}
                        initialData={airport}
                        triggerVariant="icon"
                        onSubmit={(data) => updateAirport(airport.id, data)}
                        onSuccess={fetchData}
                      />
                     <GenericDeleteButton 
                        id={airport.id} 
                        name={airport.name} 
                        deleteApi={deleteAirport} 
                        onSuccess={fetchData}
                        />
                    </div>
                    {/* Mobile Menu */}
                    <div className="md:hidden">
                      <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                          <Button variant="ghost" size="icon">
                            <MoreHorizontal className="h-4 w-4" />
                          </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end">
                          <DropdownMenuLabel>Opsi</DropdownMenuLabel>
                          <DropdownMenuItem>Edit</DropdownMenuItem>
                          <DropdownMenuItem className="text-destructive">Hapus</DropdownMenuItem>
                        </DropdownMenuContent>
                      </DropdownMenu>
                    </div>
                  </TableCell>
                </TableRow>
              ))
            ) : (
              <TableRow>
                <TableCell colSpan={4} className="h-64 text-center">
                  <div className="flex flex-col items-center justify-center text-muted-foreground">
                    <Plane className="h-12 w-12 opacity-20 mb-4" />
                    <p className="text-lg font-medium">Belum ada maskapai</p>
                    <p className="text-sm">Silakan tambahkan maskapai baru melalui dashboard.</p>
                  </div>
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>

        {/* Pagination Section */}
        <div className="flex items-center justify-between border-t bg-muted/20 px-4 py-4">
          <p className="text-sm text-muted-foreground">
            Halaman <span className="font-medium text-foreground">{currentPage}</span> dari{" "}
            <span className="font-medium text-foreground">{totalPages}</span>
          </p>
          <div className="flex items-center gap-2">
            <Button
              variant="outline"
              size="sm"
              className="h-9"
              onClick={() => setCurrentPage((prev) => Math.max(prev - 1, 1))}
              disabled={currentPage === 1 || isLoading}
            >
              <ChevronLeft className="h-4 w-4 mr-1" />
              Kembali
            </Button>
            <Button
              variant="outline"
              size="sm"
              className="h-9"
              onClick={() => setCurrentPage((prev) => Math.min(prev + 1, totalPages))}
              disabled={currentPage === totalPages || isLoading}
            >
              Lanjut
              <ChevronRight className="h-4 w-4 ml-1" />
            </Button>
          </div>
        </div>
      </div>
    </div>
  )
}