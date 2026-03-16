      "use client"

      import { useState, useEffect, useMemo } from "react"
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
      import { getFlights, deleteFlight, updateFlight } from "@/lib/api/FlightApi"
      import { Flight } from "@/lib/type/flight"
      import Image from "next/image"
      import { ChevronLeft, ChevronRight } from "lucide-react"
      import { GenericDeleteButton } from "./GenericDeleteButton"
      import { UpsertForm, FormField } from "./UpsertForm"
      import { createFlight } from "@/lib/api/FlightApi"
      import { getAirlines } from "@/lib/api/AirlineApi"
      import { getAirports } from "@/lib/api/AirportApi"
import { useFlightFields } from "@/hooks/useFlightFields"

      export function FlightsTableData({ refreshKey }: { refreshKey?: number }) {
        const { flightFields, isDataLoading: isFieldsLoading } = useFlightFields()
        
        const [flights, setFlights] = useState<Flight[]>([])
        const [isLoading, setIsLoading] = useState(true)

        const [currentPage, setCurrentPage] = useState(1)
        const [totalPages, setTotalPages] = useState(1)
        const limit = 10

        const fetchData = async () => {
          setIsLoading(true)
          try {
            const response = await getFlights(currentPage, limit) 
            setFlights(response.data)
            setTotalPages(Math.ceil(response.meta.total / response.meta.limit))
          } catch (error) {
            console.error("Gagal mengambil data:", error)
          } finally {
            setIsLoading(false)
          }
        }

        useEffect(() => {
          fetchData()
        }, [currentPage, refreshKey])

        const formatCurrency = (amount: number) => {
          return new Intl.NumberFormat('id-ID', {
            style: 'currency',
            currency: 'IDR',
            maximumFractionDigits: 0
          }).format(amount)
        }

        const formatTime = (dateString: string) => {
          return new Date(dateString).toLocaleTimeString('id-ID', {
            hour: '2-digit',
            minute: '2-digit'
          })
        }

        const formatDate = (dateString: string) => {
          return new Date(dateString).toLocaleDateString('id-ID', {
            day: '2-digit',
            month: 'short',
            year: 'numeric'
          })
        }

        if (isLoading || isFieldsLoading) return <p className="p-4 text-center text-muted-foreground animate-pulse">Memuat data penerbangan...</p>

        return (
          <div className="rounded-xl border bg-card shadow-sm overflow-hidden">
            <Table>
              <TableHeader className="bg-muted/50">
                <TableRow>
                  <TableHead className="font-bold">Maskapai</TableHead>
                  <TableHead className="font-bold">No. Penerbangan</TableHead>
                  <TableHead className="font-bold">Rute</TableHead>
                  <TableHead className="font-bold">Waktu</TableHead>
                  <TableHead className="font-bold">Kursi</TableHead>
                  <TableHead className="font-bold">Harga Terendah</TableHead>
                  <TableHead className="text-right font-bold">Aksi</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {flights && flights.length > 0 ? (
                  flights.map((flight) => {
                    const minPrice = flight.classes.length > 0 
                      ? Math.min(...flight.classes.map(c => c.price))
                      : 0
                    
                    return (
                      <TableRow key={flight.id} className="hover:bg-muted/30 transition-colors">
                        <TableCell>
                          <div className="flex items-center gap-3">
                            <div className="h-10 w-10 relative shrink-0 bg-white rounded-lg border p-1 shadow-sm">
                              <Image
                                src={flight.airline.logo_url}
                                alt={flight.airline.name}
                                fill
                                className="object-contain p-1"
                              />
                            </div>
                            <div className="flex flex-col">
                              <span className="font-semibold text-sm">{flight.airline.name}</span>
                              <span className="text-[10px] text-muted-foreground uppercase tracking-wider">{flight.airline.code}</span>
                            </div>
                          </div>
                        </TableCell>
                        <TableCell>
                          <Badge variant="secondary" className="font-mono font-medium">
                            {flight.flight_number}
                          </Badge>
                        </TableCell>
                        <TableCell>
                          <div className="flex flex-col">
                            <div className="flex items-center gap-2">
                              <span className="font-bold text-base">{flight.origin.code}</span>
                              <span className="text-muted-foreground/50 text-xs">→</span>
                              <span className="font-bold text-base">{flight.destination.code}</span>
                            </div>
                            <span className="text-[10px] text-muted-foreground truncate max-w-[150px]">
                              {flight.origin.city} ke {flight.destination.city}
                            </span>
                          </div>
                        </TableCell>
                        <TableCell>
                          <div className="flex flex-col">
                            <span className="text-xs font-medium text-muted-foreground">{formatDate(flight.departure_time)}</span>
                            <span className="font-bold text-sm">
                              {formatTime(flight.departure_time)} - {formatTime(flight.arrival_time)}
                            </span>
                          </div>
                        </TableCell>
                        <TableCell>
                          <div className="flex flex-col gap-1.5 w-32">
                            <div className="flex justify-between items-end">
                              <span className={`text-xs font-bold ${flight.available_seats === 0 ? "text-destructive" : "text-emerald-600"}`}>
                                {flight.available_seats} Tersisa
                              </span>
                              <span className="text-[10px] text-muted-foreground">/ {flight.total_seats}</span>
                            </div>
                            <div className="w-full bg-secondary h-2 rounded-full overflow-hidden">
                              <div 
                                className={`h-full transition-all duration-500 ${flight.available_seats === 0 ? "bg-destructive" : "bg-primary"}`} 
                                style={{ width: `${(flight.available_seats / flight.total_seats) * 100}%` }}
                              />
                            </div>
                          </div>
                        </TableCell>
                        <TableCell>
                          <span className="font-bold text-primary">
                            {minPrice > 0 ? formatCurrency(minPrice) : "-"}
                          </span>
                        </TableCell>
                        <TableCell className="text-right">
                          <div className="flex justify-end gap-2">
                            <UpsertForm
                              title="Flight"
                              description="Update data penerbangan"
                              fields={flightFields}
                              columns={2}
                              maxWidth="sm:max-w-[700px]"
                              initialData={{
                                ...flight,
                                class_prices: flight.classes.map((c) => ({
                                  class_type: c.class_type,
                                  price: c.price
                                }))
                              }}
                              triggerVariant="outline"
                              triggerLabel="Edit"
                              onSubmit={async (data) => {
                                const payload = {
                                  ...data,
                                  airline_id: Number(data.airline_id),
                                  origin_id: Number(data.origin_id),
                                  destination_id: Number(data.destination_id),
                                  total_rows: Number(data.total_rows),
                                  total_columns: Number(data.total_columns),
                                  total_seats: Number(data.total_rows) * Number(data.total_columns),
                                  class_count: data.class_prices?.length || 0,
                                  class_prices: data.class_prices?.map((cp: any) => ({
                                    ...cp,
                                    price: Number(cp.price)
                                  }))
                                };
                                await updateFlight(flight.id, payload);
                              }}
                              onSuccess={fetchData}
                            />

                            <GenericDeleteButton 
                              id={flight.id} 
                              name={flight.flight_number} 
                              deleteApi={deleteFlight} 
                              onSuccess={fetchData}
                            />
                          </div>
                        </TableCell>
                      </TableRow>
                    )
                  })
                ) : (
                  <TableRow>
                    <TableCell colSpan={7} className="text-center py-20">
                      <div className="flex flex-col items-center justify-center text-muted-foreground">
                        <p className="text-lg font-medium">Belum ada jadwal penerbangan</p>
                        <p className="text-sm">Silakan tambahkan jadwal penerbangan baru untuk memulai.</p>
                      </div>
                    </TableCell>
                  </TableRow>
                )}
              </TableBody>
            </Table>
            <div className="flex items-center justify-between px-2 py-4">
              <p className="text-sm text-muted-foreground">
                Menampilkan halaman <span className="font-medium">{currentPage}</span> dari{" "}
                <span className="font-medium">{totalPages}</span>
              </p>
              <div className="flex items-center space-x-2">
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => setCurrentPage((prev) => Math.max(prev - 1, 1))}
                  disabled={currentPage === 1 || isLoading}
                >
                  <ChevronLeft className="h-4 w-4 mr-2" />
                  Previous
                </Button>
                <div className="flex items-center gap-1">
                  <span className="text-sm font-medium px-2">{currentPage}</span>
                </div>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => setCurrentPage((prev) => Math.min(prev + 1, totalPages))}
                  disabled={currentPage === totalPages || isLoading}
                >
                  Next
                  <ChevronRight className="h-4 w-4 ml-2" />
                </Button>
              </div>
            </div>
          </div>
        )
      }
