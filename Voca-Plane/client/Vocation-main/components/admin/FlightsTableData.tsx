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
import { ChevronLeft, ChevronRight, Plane, Edit2 } from "lucide-react"
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

  // Single Dialog State for performance optimization
  const [isUpsertOpen, setIsUpsertOpen] = useState(false)
  const [selectedFlight, setSelectedFlight] = useState<any>(null)

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

  const handleEdit = (flight: Flight) => {
    setSelectedFlight({
      ...flight,
      class_prices: flight.classes.map((c) => ({
        class_type: c.class_type,
        price: c.price
      }))
    })
    setIsUpsertOpen(true)
  }

        if (isLoading || isFieldsLoading) return <p className="p-4 text-center text-muted-foreground animate-pulse">Memuat data penerbangan...</p>

  return (
    <div className="space-y-4">
      {/* Single UpsertForm Instance for performance */}
      <UpsertForm
        open={isUpsertOpen}
        onOpenChange={setIsUpsertOpen}
        title="Flight"
        description="Securely update flight schedules and configurations."
        fields={flightFields}
        columns={2}
        maxWidth="sm:max-w-[1100px]"
        initialData={selectedFlight}
        triggerVariant="hidden"
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
          await updateFlight(selectedFlight.id, payload);
        }}
        onSuccess={fetchData}
      />

      <div className="rounded-[2.5rem] border border-slate-200/50 dark:border-white/5 bg-white/80 dark:bg-slate-900/80 backdrop-blur-md shadow-xl overflow-hidden animate-in fade-in slide-in-from-bottom-8 duration-500">
            <div className="overflow-x-auto">
              <Table>
                <TableHeader className="bg-slate-50/50 dark:bg-white/5">
                  <TableRow className="hover:bg-transparent border-b border-slate-100 dark:border-white/5">
                    <TableHead className="font-black uppercase tracking-[0.2em] text-[10px] text-slate-400 py-6 px-10">Carrier</TableHead>
                    <TableHead className="font-black uppercase tracking-[0.2em] text-[10px] text-slate-400 py-6">Code</TableHead>
                    <TableHead className="font-black uppercase tracking-[0.2em] text-[10px] text-slate-400 py-6">Route</TableHead>
                    <TableHead className="font-black uppercase tracking-[0.2em] text-[10px] text-slate-400 py-6">Schedule</TableHead>
                    <TableHead className="font-black uppercase tracking-[0.2em] text-[10px] text-slate-400 py-6">Inventory</TableHead>
                    <TableHead className="font-black uppercase tracking-[0.2em] text-[10px] text-slate-400 py-6">Base Fare</TableHead>
                    <TableHead className="text-right font-black uppercase tracking-[0.2em] text-[10px] text-slate-400 py-6 px-10">Control</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {flights && flights.length > 0 ? (
                    flights.map((flight) => {
                      const minPrice = flight.classes.length > 0 
                        ? Math.min(...flight.classes.map(c => c.price))
                        : 0
                      
                      return (
                        <TableRow key={flight.id} className="hover:bg-slate-50/50 dark:hover:bg-white/2 transition-all border-b border-slate-50 dark:border-white/2 group">
                          <TableCell className="px-10 py-6">
                            <div className="flex items-center gap-4">
                              <div className="h-12 w-12 relative shrink-0 bg-white rounded-xl border border-slate-100 dark:border-white/10 p-2 shadow-sm group-hover:scale-110 transition-transform">
                                <Image
                                  src={flight.airline.logo_url}
                                  alt={flight.airline.name}
                                  fill
                                  className="object-contain p-1"
                                />
                              </div>
                              <div className="flex flex-col gap-0.5">
                                <span className="font-black text-slate-900 dark:text-white text-sm tracking-tight">{flight.airline.name}</span>
                                <span className="text-[10px] text-primary font-black uppercase tracking-widest">{flight.airline.code}</span>
                              </div>
                            </div>
                          </TableCell>
                          <TableCell>
                            <Badge variant="outline" className="font-black text-[10px] px-3 py-1 rounded-full border-primary/20 text-primary bg-primary/5">
                              {flight.flight_number}
                            </Badge>
                          </TableCell>
                          <TableCell>
                            <div className="flex flex-col gap-1">
                              <div className="flex items-center gap-2">
                                <span className="font-black text-slate-900 dark:text-white text-sm">{flight.origin.code}</span>
                                <div className="h-[1px] w-4 bg-slate-200 dark:bg-white/10" />
                                <span className="font-black text-slate-900 dark:text-white text-sm">{flight.destination.code}</span>
                              </div>
                              <span className="text-[10px] text-slate-400 font-bold uppercase tracking-wider truncate max-w-[150px]">
                                {flight.origin.city.split(',')[0]} to {flight.destination.city.split(',')[0]}
                              </span>
                            </div>
                          </TableCell>
                          <TableCell>
                            <div className="flex flex-col gap-0.5">
                              <span className="text-[10px] font-black uppercase tracking-wider text-slate-400">{formatDate(flight.departure_time)}</span>
                              <span className="font-black text-slate-900 dark:text-white text-sm">
                                {formatTime(flight.departure_time)} <span className="text-slate-300 font-medium">→</span> {formatTime(flight.arrival_time)}
                              </span>
                            </div>
                          </TableCell>
                          <TableCell>
                            <div className="flex flex-col gap-2 w-36">
                              <div className="flex justify-between items-end">
                                <span className={`text-[10px] font-black uppercase tracking-wider ${flight.available_seats === 0 ? "text-red-500" : "text-emerald-500"}`}>
                                  {flight.available_seats} AVAIL
                                </span>
                                <span className="text-[9px] font-black text-slate-400">/ {flight.total_seats}</span>
                              </div>
                              <div className="w-full bg-slate-100 dark:bg-white/5 h-1.5 rounded-full overflow-hidden">
                                <div 
                                  className={`h-full transition-all duration-500 ease-out ${flight.available_seats === 0 ? "bg-red-500" : "bg-primary"}`} 
                                  style={{ width: `${(flight.available_seats / flight.total_seats) * 100}%` }}
                                />
                              </div>
                            </div>
                          </TableCell>
                          <TableCell>
                            <span className="font-black text-primary text-sm italic tracking-tight">
                              {minPrice > 0 ? formatCurrency(minPrice) : "-"}
                            </span>
                          </TableCell>
                          <TableCell className="text-right px-10">
                              <div className="flex justify-end gap-3 translate-x-4 opacity-0 group-hover:opacity-100 group-hover:translate-x-0 transition-all duration-300">
                                <Button 
                                  variant="outline" 
                                  size="sm" 
                                  className="rounded-xl h-9 px-4 font-black uppercase tracking-widest text-[10px] gap-2 transition-all hover:bg-slate-100 dark:hover:bg-white/5"
                                  onClick={() => handleEdit(flight)}
                                >
                                  <Edit2 className="size-3" />
                                  Edit
                                </Button>
  
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
                      <TableCell colSpan={7} className="text-center py-32">
                        <div className="flex flex-col items-center justify-center text-slate-400 gap-4">
                          <div className="size-16 rounded-3xl bg-slate-50 dark:bg-white/5 flex items-center justify-center border border-slate-100 dark:border-white/10 shadow-inner">
                            <Plane className="size-8 rotate-45 opacity-20" />
                          </div>
                          <div className="space-y-1">
                            <p className="text-lg font-black text-slate-900 dark:text-white tracking-tight">Empty Hangar</p>
                            <p className="text-sm font-medium">Initialize new flight schedules to populate this terminal.</p>
                          </div>
                        </div>
                      </TableCell>
                    </TableRow>
                  )}
                </TableBody>
              </Table>
            </div>
            <div className="flex items-center justify-between px-10 py-8 bg-slate-50/50 dark:bg-white/2 border-t border-slate-100 dark:border-white/5">
              <p className="text-[10px] font-black uppercase tracking-[0.25em] text-slate-400">
                Segment <span className="text-primary font-black">{currentPage}</span> of{" "}
                <span className="text-slate-600 dark:text-white">{totalPages}</span>
              </p>
              <div className="flex items-center gap-4">
                <Button
                  variant="ghost"
                  size="sm"
                  className="rounded-2xl h-11 px-6 font-black uppercase tracking-widest text-[10px] hover:bg-white dark:hover:bg-white/5 shadow-sm transition-all active:scale-95 disabled:opacity-30"
                  onClick={() => setCurrentPage((prev) => Math.max(prev - 1, 1))}
                  disabled={currentPage === 1 || isLoading}
                >
                  <ChevronLeft className="h-4 w-4 mr-2" />
                  Previous
                </Button>
                
                <div className="size-11 rounded-2xl bg-primary text-white flex items-center justify-center font-black text-xs shadow-lg shadow-primary/30">
                  {currentPage}
                </div>

                <Button
                  variant="ghost"
                  size="sm"
                  className="rounded-2xl h-11 px-6 font-black uppercase tracking-widest text-[10px] hover:bg-white dark:hover:bg-white/5 shadow-sm transition-all active:scale-95 disabled:opacity-30"
                  onClick={() => setCurrentPage((prev) => Math.min(prev + 1, totalPages))}
                  disabled={currentPage === totalPages || isLoading}
                >
                  Next
                  <ChevronRight className="h-4 w-4 ml-2" />
                </Button>
              </div>
            </div>
      </div>
    </div>
  )
}
