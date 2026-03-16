import Image from "next/image"
import { Plane } from "lucide-react"
import { Card, CardContent } from "@/components/ui/card"
import Link from "next/link"
import { Flight } from "@/lib/type/flight"

interface FlightCardProps {
  flight: Flight
}

export function CardTicket({ flight }: FlightCardProps) {
  const minPrice = Math.min(...flight.classes.map((c) => c.price))
  
  return (
    <div className="relative w-full max-w-95 group">
      <Card className="relative h-72 rounded-[2.5rem] overflow-hidden border-none shadow-2xl transition-all duration-500 group-hover:shadow-primary/20 group-hover:-translate-y-2">
        <Image
          src="/destination.jpg"
          alt="Destination"
          fill
          className="object-cover transition-transform duration-1000 group-hover:scale-110"
        />

        <div className="absolute inset-0 bg-gradient-to-t from-slate-950/90 via-slate-950/20 to-transparent" />

        <div className="absolute inset-0 p-8 flex flex-col justify-between text-white">
          <div className="flex justify-between items-start">
            <span className="bg-white/10 backdrop-blur-xl px-4 py-2 rounded-full text-[10px] font-bold border border-white/20 uppercase tracking-widest leading-none">
              {flight.flight_number}
            </span>

            <span className="text-[11px] font-black tracking-[0.3em] uppercase opacity-80">
              {flight.origin.code} <Plane className="inline size-3 mx-1 opacity-50 rotate-45" /> {flight.destination.code}
            </span>
          </div>

          <div className="space-y-1">
            <h3 className="text-4xl font-black tracking-tighter leading-none">
              {flight.destination.city}
            </h3>
            <p className="text-[11px] font-bold text-white/70 uppercase tracking-widest">
              {new Date(flight.departure_time).toLocaleDateString("id-ID", {
                weekday: "short",
                day: "numeric",
                month: "short"
              })}
            </p>
          </div>
        </div>
      </Card>

      <div className="absolute -bottom-12 left-1/2 -translate-x-1/2 w-[82%] group-hover:w-[88%] transition-all duration-500">
        <Card className="rounded-[2rem] shadow-2xl border border-white/50 bg-white/95 backdrop-blur-2xl">
          <Link href={`/flight/${flight.id}`}>
            <CardContent className="p-6 flex justify-between items-center">
              <div className="space-y-0.5">
                <p className="text-[10px] text-slate-500 font-extrabold uppercase tracking-widest">
                  Best Fare
                </p>
                <div className="flex items-baseline gap-0.5">
                  <span className="text-[10px] font-bold text-slate-400">IDR</span>
                  <p className="text-2xl font-black text-primary tracking-tighter">
                    {minPrice.toLocaleString("id-ID")}
                  </p>
                </div>
              </div>

              <div className="h-14 w-14 rounded-3xl bg-primary text-white flex items-center justify-center shadow-xl shadow-primary/20 transition-all duration-300 hover:scale-110 hover:rotate-12 active:scale-95">
                <Plane className="size-6 rotate-45" />
              </div>
            </CardContent>
          </Link>
        </Card>
      </div>
    </div>
  )
}