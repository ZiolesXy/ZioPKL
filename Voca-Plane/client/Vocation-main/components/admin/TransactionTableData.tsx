"use client"

import { useEffect, useState, useCallback } from "react"
import Image from "next/image"
import { 
  ChevronLeft, 
  ChevronRight, 
  Edit2, 
  Trash2, 
  Plane, 
  MoreHorizontal 
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

import { getAllTransaction } from "@/lib/api/TransactionApi"
import { Transaction } from "@/lib/type/transaction"

export function TransactionTableData() {
  const [transactions, setTransactions] = useState<Transaction[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [currentPage, setCurrentPage] = useState(1)
  const [totalPages, setTotalPages] = useState(1)
  const limit = 10

  const fetchData = useCallback(async () => {
    setIsLoading(true)
    try {
      const response = await getAllTransaction(currentPage, limit)
      setTransactions(response.data)
      setTotalPages(Math.ceil(response.meta.total / response.meta.limit))
    } catch (error) {
      console.error("Gagal mengambil data transaksi:", error)
    } finally {
      setIsLoading(false)
    }
  }, [currentPage, limit])

  useEffect(() => {
    fetchData()
  }, [fetchData])

  if (isLoading) {
    return (
      <div className="flex h-64 items-center justify-center rounded-xl border bg-card">
        <div className="flex flex-col items-center gap-2">
          <div className="h-8 w-8 animate-spin rounded-full border-4 border-primary border-t-transparent" />
          <p className="text-sm text-muted-foreground font-medium">Memuat data Transaksi...</p>
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
              <TableHead className="w-[80px]">ID</TableHead>
              <TableHead>Penerbangan</TableHead>
              <TableHead>Rute</TableHead>
              <TableHead>Penumpang</TableHead>
              <TableHead>Total Bayar</TableHead>
              <TableHead>Status</TableHead>
              <TableHead className="text-right">Aksi</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {transactions.length > 0 ? (
              transactions.map((tx) => (
                <TableRow key={tx.id} className="group hover:bg-muted/30 transition-colors">
                  <TableCell className="font-mono text-xs text-muted-foreground">
                    #{tx.id}
                  </TableCell>
                  <TableCell>
                    <div className="flex items-center gap-3">
                      <div className="relative h-8 w-8 shrink-0 overflow-hidden rounded-md border bg-white p-1">
                        <Image
                          src={tx.flight.airline.logo_url}
                          alt={tx.flight.airline.name}
                          fill
                          className="object-contain"
                        />
                      </div>
                      <div className="flex flex-col">
                        <span className="font-bold text-sm">{tx.flight.airline.name}</span>
                        <span className="text-xs text-muted-foreground">{tx.flight.flight_number}</span>
                      </div>
                    </div>
                  </TableCell>
                  <TableCell>
                    <div className="flex flex-col">
                      <span className="font-medium text-sm">{tx.flight.origin.code} → {tx.flight.destination.code}</span>
                      <span className="text-xs text-muted-foreground">{tx.flight.origin.city}</span>
                    </div>
                  </TableCell>
                  <TableCell>
                    <div className="flex flex-col">
                      <span className="font-bold text-sm">{tx.transactions_passangers.length} Penumpang</span>
                      <span className="text-[10px] text-muted-foreground">Detail dipesan</span>
                    </div>
                  </TableCell>
                  <TableCell>
                    <div className="flex flex-col">
                      <span className="font-bold text-sm text-primary">
                        {new Intl.NumberFormat("id-ID", { style: "currency", currency: "IDR", maximumFractionDigits: 0 }).format(tx.total_price)}
                      </span>
                      {tx.discount > 0 && <span className="text-[10px] text-destructive">Hemat Rp {tx.discount.toLocaleString()}</span>}
                    </div>
                  </TableCell>
                  <TableCell>
                    <Badge 
                      variant={tx.payment_status === "PAID" ? "default" : "destructive"}
                      className={tx.payment_status === "PAID" ? "bg-green-500 hover:bg-green-600" : ""}
                    >
                      {tx.payment_status}
                    </Badge>
                  </TableCell>
                  <TableCell className="text-right">
                    <Button variant="outline" size="sm" className="h-8">Detail</Button>
                  </TableCell>
                </TableRow>
              ))
            ) : (
              <TableRow>
                <TableCell colSpan={4} className="h-64 text-center">
                  <div className="flex flex-col items-center justify-center text-muted-foreground">
                    <Plane className="h-12 w-12 opacity-20 mb-4" />
                    <p className="text-lg font-medium">Belum ada Transaksi</p>
                    <p className="text-sm">Transaksi akan muncul setelah user melakukan transaksi.</p>
                  </div>
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>

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