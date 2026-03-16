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
    <div className="grid grid-cols-1 lg:grid-cols-3 gap-12 mt-16 animate-in fade-in slide-in-from-bottom-8 duration-700">
      <div className="lg:col-span-2 space-y-10">
        <div className="flex items-center justify-between">
          <h2 className="text-3xl font-black tracking-tighter text-slate-900 dark:text-white">Passenger Details</h2>
          <span className="text-[10px] font-black uppercase tracking-[0.2em] text-primary bg-primary/10 px-4 py-1.5 rounded-full">Secure Checkout</span>
        </div>

        {passengers.map((p, index) => (
          <Card key={index} className="overflow-hidden border-none shadow-[0_20px_50px_rgba(0,0,0,0.08)] bg-white/60 dark:bg-slate-900/40 backdrop-blur-3xl rounded-[2.5rem] transition-all hover:shadow-primary/5 border border-white/20 dark:border-white/5">
            <CardHeader className="bg-slate-50/50 dark:bg-white/5 px-10 py-6 border-b border-white/20 dark:border-white/5 flex flex-row items-center justify-between">
               <div className="flex items-center gap-4">
                 <div className="size-10 rounded-xl bg-primary/20 flex items-center justify-center text-primary font-black text-xs">
                   0{index + 1}
                 </div>
                 <h3 className="font-black text-slate-900 dark:text-white uppercase tracking-widest text-xs">Passenger Information</h3>
               </div>
               {index > 0 && (
                 <button 
                   onClick={() => setPassengers(passengers.filter((_, i) => i !== index))}
                   className="text-[10px] font-black uppercase tracking-widest text-red-400 hover:text-red-500 transition-colors"
                 >
                   Remove
                 </button>
               )}
            </CardHeader>
            <CardContent className="grid grid-cols-1 md:grid-cols-2 gap-8 p-10">
              {/* Input Nama */}
              <div className="space-y-3">
                <label className="text-[11px] font-black uppercase tracking-[0.2em] text-slate-400 ml-2">Full Name (as passport)</label>
                <Input 
                  placeholder="e.g. ALEXANDER SMITH"
                  className="h-14 px-6 bg-slate-100/50 dark:bg-white/5 border-none rounded-2xl focus:ring-4 focus:ring-primary/20 text-slate-900 dark:text-white font-bold"
                  value={p.full_name} 
                  onChange={(e) => {
                    const newP = [...passengers]; newP[index].full_name = e.target.value; setPassengers(newP);
                  }} 
                />
              </div>

              {/* Input Kewarganegaraan */}
              <div className="space-y-3">
                <label className="text-[11px] font-black uppercase tracking-[0.2em] text-slate-400 ml-2">Nationality</label>
                <Input 
                  placeholder="e.g. INDONESIA"
                  className="h-14 px-6 bg-slate-100/50 dark:bg-white/5 border-none rounded-2xl focus:ring-4 focus:ring-primary/20 text-slate-900 dark:text-white font-bold"
                  value={p.nationality} 
                  onChange={(e) => {
                    const newP = [...passengers]; newP[index].nationality = e.target.value; setPassengers(newP);
                  }} 
                />
              </div>

              {/* Input Passport */}
              <div className="md:col-span-2 space-y-3">
                <label className="text-[11px] font-black uppercase tracking-[0.2em] text-slate-400 ml-2">Passport / National ID Number</label>
                <Input 
                  placeholder="e.g. A12345678"
                  className="h-14 px-6 bg-slate-100/50 dark:bg-white/5 border-none rounded-2xl focus:ring-4 focus:ring-primary/20 text-slate-900 dark:text-white font-bold"
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
          className="w-full border-dashed border-2 py-10 rounded-[2.5rem] text-slate-400 hover:text-primary hover:border-primary hover:bg-primary/5 transition-all font-black uppercase tracking-[0.3em] text-xs"
        >
          + Add another passenger
        </Button>
      </div>

      <div className="lg:col-span-1">
        <div className="sticky top-32 space-y-8">
          <Card className="rounded-[3rem] shadow-[0_32px_80px_-16px_rgba(0,0,0,0.1)] border-none bg-white/80 dark:bg-slate-900/80 backdrop-blur-3xl overflow-hidden border border-white/20 dark:border-white/5">
            <CardHeader className="p-10 pb-0">
               <h3 className="font-black text-2xl tracking-tighter text-slate-900 dark:text-white">Seat Selection</h3>
            </CardHeader>
            <CardContent className="p-10">
              <SeatGrid 
                flightData={flight} 
                selectedSeats={passengers.map(p => p.seat_number).filter(s => s !== "")} 
                onSeatChange={handleSeatSelect} 
              />
              
              <div className="mt-12 pt-10 border-t border-slate-100 dark:border-white/5 space-y-4">
                <div className="flex justify-between items-center text-xs font-black uppercase tracking-widest">
                  <span className="text-slate-400">Total Passengers</span>
                  <span className="text-slate-900 dark:text-white">{passengers.length}</span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-xs font-black uppercase tracking-widest text-slate-400">Grand Total</span>
                  <span className="text-2xl font-black text-primary tracking-tighter italic">IDR {(passengers.length * (flight.classes[0]?.price || 0)).toLocaleString("id-ID")}</span>
                </div>
              </div>

              <Button 
                className="w-full mt-10 h-16 bg-primary text-white font-black text-lg rounded-2xl shadow-2xl shadow-primary/30 transition-all hover:scale-[1.02] active:scale-95 disabled:opacity-50"
                disabled={isLoading}
                onClick={handleCheckout}
              >
                {isLoading ? (
                  <div className="flex items-center gap-2">
                    <span className="w-2 h-2 rounded-full bg-white animate-bounce" />
                    <span className="w-2 h-2 rounded-full bg-white animate-bounce [animation-delay:0.2s]" />
                    <span className="w-2 h-2 rounded-full bg-white animate-bounce [animation-delay:0.4s]" />
                  </div>
                ) : "Confirm & Pay"}
              </Button>
            </CardContent>
          </Card>

          <p className="text-center text-[10px] font-black uppercase tracking-[0.2em] text-slate-400">
            Encrypted with 256-bit SSL
          </p>
        </div>
      </div>
    </div>
  )
}