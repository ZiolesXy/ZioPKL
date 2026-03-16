import Image from "next/image"
import { Plane, ArrowRight } from "lucide-react"
import { Card, CardContent } from "@/components/ui/card"
import Link from "next/link"
import { Flight } from "@/lib/type/flight"

interface FlightCardProps {
  flight: Flight
}

export function CardTicket({ flight }: FlightCardProps) {
  const minPrice = Math.min(...flight.classes.map((c) => c.price))
  
  return (
    <div className="relative group">
      <Card className="relative h-72 rounded-[3rem] overflow-hidden border-none shadow-xl transition-transform duration-500 group-hover:-translate-y-2">
        <Image
          src="/destination.jpg"
          alt="Destination"
          fill
          className="object-cover transition-transform duration-1000 group-hover:scale-110"
        />

        {/* Dynamic Overlay */}
        <div className="absolute inset-0 bg-gradient-to-t from-slate-950 via-slate-950/20 to-transparent opacity-80" />
        
        <div className="absolute inset-0 p-8 flex flex-col justify-between text-white z-10">
          <div className="flex justify-between items-start">
            <span className="bg-white/10 backdrop-blur-sm px-4 py-1.5 rounded-full text-[10px] font-black border border-white/20 uppercase tracking-[0.2em]">
              {flight.flight_number}
            </span>

            <div className="flex items-center gap-2 bg-black/20 backdrop-blur-sm px-3 py-1 rounded-full border border-white/10">
              <span className="text-[10px] font-black uppercase tracking-widest">{flight.origin.code}</span>
              <Plane className="size-3 text-primary rotate-45" />
              <span className="text-[10px] font-black uppercase tracking-widest">{flight.destination.code}</span>
            </div>
          </div>

          <div className="space-y-2 translate-y-4 group-hover:translate-y-0 transition-transform duration-300">
            <h3 className="text-5xl font-black tracking-tighter leading-none group-hover:text-primary transition-colors">
              {flight.destination.city}
            </h3>
            <div className="flex items-center gap-2 text-white/60 font-bold uppercase tracking-[0.15em] text-[10px]">
              <span className="w-8 h-px bg-white/20" />
              {new Date(flight.departure_time).toLocaleDateString("id-ID", {
                weekday: "short",
                day: "numeric",
                month: "short"
              })}
            </div>
          </div>
        </div>
      </Card>

      {/* Floating Info Layer */}
      <div className="absolute -bottom-8 left-1/2 -translate-x-1/2 w-[85%] transition-all duration-300 group-hover:w-[90%] -z-10 group-hover:z-20">
        <Link href={`/flight/${flight.id}`}>
          <Card className="rounded-[2.5rem] shadow-xl border border-white/40 bg-white/95 dark:bg-slate-900/95 backdrop-blur-md overflow-hidden hover:bg-white transition-colors">
            <CardContent className="p-6 flex justify-between items-center">
              <div className="space-y-0.5">
                <p className="text-[9px] text-slate-400 font-black uppercase tracking-[0.3em]">
                  Start from
                </p>
                <div className="flex items-baseline gap-1">
                  <span className="text-xs font-black text-primary/60">IDR</span>
                  <p className="text-2xl font-black text-slate-900 dark:text-white tracking-tighter">
                    {minPrice.toLocaleString("id-ID")}
                  </p>
                </div>
              </div>

              <div className="h-14 w-14 rounded-2xl bg-primary text-white flex items-center justify-center shadow-lg shadow-primary/30 transition-all duration-300 group-hover:scale-105">
                <ArrowRight className="size-6" />
              </div>
            </CardContent>
          </Card>
        </Link>
      </div>
    </div>
  )
}