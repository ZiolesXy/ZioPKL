"use client"

import { useState } from "react"
import { SeatGrid } from "@/components/SeatGrid"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { toast } from "sonner" // Import langsung dari sonner
import { cn } from "@/lib/utils"

export function BookingForm({ flight }: { flight: any }) {
  const [passengers, setPassengers] = useState([
    { full_name: "", nationality: "Indonesia", passport_no: "", seat_number: "" }
  ])

  const addPassenger = () => {
    if (passengers.length < flight.available_seats) {
      setPassengers([...passengers, { full_name: "", nationality: "Indonesia", passport_no: "", seat_number: "" }])
      toast.success("Penumpang ditambahkan") // Feedback simpel
    } else {
      toast.error("Batas kursi tercapai")
    }
  }

  const handleSeatSelect = (seatCode: string) => {
    const isAlreadyChosen = passengers.some(p => p.seat_number === seatCode)
    
    if (isAlreadyChosen) {
      setPassengers(passengers.map(p => 
        p.seat_number === seatCode ? { ...p, seat_number: "" } : p
      ))
      toast.info(`Kursi ${seatCode} dilepas`)
      return
    }

    const emptySeatIndex = passengers.findIndex(p => p.seat_number === "")
    
    if (emptySeatIndex !== -1) {
      const newPassengers = [...passengers]
      newPassengers[emptySeatIndex].seat_number = seatCode
      setPassengers(newPassengers)
      toast.success(`Kursi ${seatCode} dipilih untuk Penumpang ${emptySeatIndex + 1}`)
    } else {
      toast.warning("Slot Penuh", {
        description: "Semua penumpang sudah punya kursi. Tambah penumpang baru jika perlu.",
      })
    }
  }

  const selectedSeatCodes = passengers.map(p => p.seat_number).filter(s => s !== "")

  return (
    <div className="grid grid-cols-1 lg:grid-cols-3 gap-8 mt-8">
      <div className="lg:col-span-2 space-y-6">
        {passengers.map((p, index) => (
          <Card key={index} className="overflow-hidden border-l-4 border-l-indigo-600">
            <CardHeader className="bg-slate-50/50">
              <CardTitle className="text-sm font-bold flex justify-between">
                <span>PENUMPANG {index + 1}</span>
                {p.seat_number && <span className="text-indigo-600">Kursi: {p.seat_number}</span>}
              </CardTitle>
            </CardHeader>
            <CardContent className="grid grid-cols-1 md:grid-cols-2 gap-4 pt-6">
              <div className="space-y-2">
                <label className="text-xs font-semibold text-slate-500">Nama Lengkap</label>
                <Input 
                  placeholder="Contoh: Budi Santoso" 
                  value={p.full_name} 
                  onChange={(e) => {
                    const newP = [...passengers]; newP[index].full_name = e.target.value; setPassengers(newP);
                  }} 
                />
              </div>
              <div className="space-y-2">
                <label className="text-xs font-semibold text-slate-500">No. Passport / ID</label>
                <Input 
                  placeholder="A1234567" 
                  value={p.passport_no} 
                  onChange={(e) => {
                    const newP = [...passengers]; newP[index].passport_no = e.target.value; setPassengers(newP);
                  }} 
                />
              </div>
            </CardContent>
          </Card>
        ))}
        <Button 
          variant="outline" 
          onClick={addPassenger}
          className="w-full border-dashed border-2 py-8 text-slate-500 hover:text-indigo-600 hover:border-indigo-600 transition-all"
        >
          + Tambah Penumpang
        </Button>
      </div>

      <div className="lg:col-span-1">
        <div className="sticky top-6 space-y-6">
          <Card>
            <CardHeader>
              <CardTitle className="text-center text-lg">Konfigurasi Kursi</CardTitle>
            </CardHeader>
            <CardContent>
              <SeatGrid 
                flightData={flight} 
                selectedSeats={selectedSeatCodes} 
                onSeatChange={handleSeatSelect} 
              />
              <Button 
                className="w-full mt-8 bg-indigo-600 hover:bg-indigo-700 h-12 font-bold"
                onClick={() => {
                   if (selectedSeatCodes.length < passengers.length) {
                     toast.error("Lengkapi Kursi", { description: "Beberapa penumpang belum memilih kursi!" });
                   } else {
                     toast.promise(new Promise((res) => setTimeout(res, 2000)), {
                       loading: 'Menyiapkan transaksi...',
                       success: 'Booking berhasil dibuat!',
                       error: 'Gagal membuat booking',
                     });
                   }
                }}
              >
                Lanjut ke Pembayaran
              </Button>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  )
}