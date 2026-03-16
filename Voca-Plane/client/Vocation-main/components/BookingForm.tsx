"use client"

import { useState } from "react"
import { useRouter } from "next/navigation" // Untuk redirect setelah transaksi berhasil
import { SeatGrid } from "@/components/SeatGrid"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Card, CardContent, CardHeader } from "@/components/ui/card"
import { toast } from "sonner"
import { createTransaction } from "@/lib/api/TransactionApi" // Import api call kamu

export function BookingForm({ flight }: { flight: any }) {
  const router = useRouter()
  const [isLoading, setIsLoading] = useState(false)
  const [passengers, setPassengers] = useState([
    { full_name: "", nationality: "", passport_no: "", seat_number: "" }
  ])

   const addPassenger = () => {
    if (passengers.length < flight.available_seats) {
      setPassengers([...passengers, { full_name: "agus", nationality: "Indonesia", passport_no: "11111", seat_number: "" }])
      toast.success("Penumpang ditambahkan") // Feedback simpel
    } else {
      toast.error("Batas kursi tercapai")
    }
  }
  // Ambil class_id dari flight (misal user memilih kelas pertama yang tersedia)
  // Berdasarkan database kamu, transaksi butuh flight_id dan class_id [cite: 78, 80]
  const selectedClassId = flight.classes[0]?.id

  const handleCheckout = async () => {
    // 1. Validasi Kelengkapan Data
    const isDataIncomplete = passengers.some(
      p => !p.full_name || !p.nationality || !p.passport_no || !p.seat_number
    )

    if (isDataIncomplete) {
      toast.error("Data Belum Lengkap", {
        description: "Pastikan semua informasi penumpang dan kursi telah terisi."
      })
      return
    }

    setIsLoading(true)

    // 2. Susun Payload sesuai kebutuhan API
    const transactionPayload = {
      flight_id: flight.id,
      class_id: selectedClassId,
      promo_code: "", // B  isa dikembangkan dengan field input promo nanti 
      passengers: passengers
    }

    try {
      const response = await createTransaction(transactionPayload)
      
      toast.success("Booking Berhasil!", {
        description: "Pesanan Anda telah tercatat dalam sistem."
      })

      // 3. Redirect ke halaman status pembayaran menggunakan code transaksi 
      router.push(`/my-bookings/${response.data.code}`)
      
    } catch (error: any) {
      toast.error("Gagal Membuat Pesanan", {
        description: error.response?.data?.message || "Terjadi kesalahan pada server."
      })
    } finally {
      setIsLoading(false)
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


  return (
    <div className="grid grid-cols-1 lg:grid-cols-3 gap-8 mt-8">
      <div className="lg:col-span-2 space-y-6">
        {passengers.map((p, index) => (
          <Card key={index} className="overflow-hidden border-l-4 border-l-indigo-600">
            <CardHeader className="bg-slate-50/50">
               {/* UI Header Penumpang */}
            </CardHeader>
            <CardContent className="grid grid-cols-1 md:grid-cols-2 gap-4 pt-6">
              {/* Input Nama */}
              <div className="space-y-2">
                <label className="text-xs font-semibold text-slate-500">Nama Lengkap</label>
                <Input 
                  value={p.full_name} 
                  onChange={(e) => {
                    const newP = [...passengers]; newP[index].full_name = e.target.value; setPassengers(newP);
                  }} 
                />
              </div>

              {/* Input Kewarganegaraan (Dinamis, tidak hardcode lagi) */}
              <div className="space-y-2">
                <label className="text-xs font-semibold text-slate-500">Kewarganegaraan</label>
                <Input 
                  placeholder="Contoh: Indonesia"
                  value={p.nationality} 
                  onChange={(e) => {
                    const newP = [...passengers]; newP[index].nationality = e.target.value; setPassengers(newP);
                  }} 
                />
              </div>

              {/* Input Passport */}
              <div className="space-y-2">
                <label className="text-xs font-semibold text-slate-500">No. Passport / ID</label>
                <Input 
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
        <Card>
          <CardContent className="pt-6">
            <SeatGrid 
              flightData={flight} 
              selectedSeats={passengers.map(p => p.seat_number).filter(s => s !== "")} 
              onSeatChange={handleSeatSelect} 
            />
            <Button 
              className="w-full mt-8 bg-indigo-600 font-bold"
              disabled={isLoading}
              onClick={handleCheckout}
            >
              {isLoading ? "Memproses..." : "Lanjut ke Pembayaran"}
            </Button>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}