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
      <Card className="relative h-64 rounded-3xl overflow-hidden border-none shadow-md transition-all duration-300 group-hover:shadow-lg">
        <Image
          src="/destination.jpg"
          alt="Destination"
          fill
          sizes="(max-width: 768px) 100vw, 33vw"
          className="object-cover transition-transform duration-500 group-hover:scale-105"
          priority={false}
        />

        {/* Dynamic Overlay - Optimized to simple gradient */}
        <div className="absolute inset-0 bg-gradient-to-t from-slate-950/90 via-slate-950/30 to-transparent" />
        
        <div className="absolute inset-0 p-6 flex flex-col justify-between text-white z-10">
          <div className="flex justify-between items-start">
            <span className="bg-white/20 px-3 py-1 rounded-full text-[10px] font-bold border border-white/20 uppercase tracking-wider">
              {flight.flight_number}
            </span>

            <div className="flex items-center gap-2 bg-black/40 px-3 py-1 rounded-full border border-white/10">
              <span className="text-[10px] font-bold uppercase tracking-wider">{flight.origin.code}</span>
              <Plane className="size-3 text-primary rotate-45" />
              <span className="text-[10px] font-bold uppercase tracking-wider">{flight.destination.code}</span>
            </div>
          </div>

          <div className="space-y-1">
            <h3 className="text-3xl font-bold tracking-tight leading-tight">
              {flight.destination.city}
            </h3>
            <div className="flex items-center gap-2 text-white/70 font-medium uppercase tracking-wider text-[10px]">
              <span className="w-6 h-px bg-white/30" />
              {new Date(flight.departure_time).toLocaleDateString("id-ID", {
                weekday: "short",
                day: "numeric",
                month: "short"
              })}
            </div>
          </div>
        </div>
      </Card>

      {/* Floating Info Layer - Simplified */}
      <div className="absolute -bottom-6 left-1/2 -translate-x-1/2 w-[90%] z-20">
        <Link href={`/flight/${flight.id}`}>
          <Card className="rounded-2xl shadow-lg border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 transition-colors hover:bg-slate-50 dark:hover:bg-slate-800">
            <CardContent className="p-4 flex justify-between items-center">
              <div className="space-y-0.5">
                <p className="text-[9px] text-slate-500 font-bold uppercase tracking-wider">
                  Start from
                </p>
                <div className="flex items-baseline gap-1">
                  <span className="text-xs font-bold text-primary/80">IDR</span>
                  <p className="text-xl font-bold text-slate-900 dark:text-white tracking-tight">
                    {minPrice.toLocaleString("id-ID")}
                  </p>
                </div>
              </div>

              <div className="h-10 w-10 rounded-xl bg-primary text-white flex items-center justify-center shadow-md transition-transform duration-200 group-hover:scale-105">
                <ArrowRight className="size-5" />
              </div>
            </CardContent>
          </Card>
        </Link>
      </div>
    </div>
  )
}