import React from "react"
import { Checkbox } from "@/components/ui/checkbox"
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger
} from "@/components/ui/tooltip"
import { cn } from "@/lib/utils"

interface SeatGridProps {
  flightData: any
  selectedSeats: string[]
  onSeatChange: (seatCode: string) => void
}

export function SeatGrid({ flightData, selectedSeats, onSeatChange }: SeatGridProps) {
  const { flight_seats, total_columns } = flightData

  const aisleIndex = Math.floor(total_columns / 2)

  // =============================
  // GROUP SEAT BERDASARKAN ROW
  // =============================
  const seatMap: Record<number, any[]> = {}

  flight_seats.forEach((seat: any) => {
    const row = parseInt(seat.seat_code)

    if (!seatMap[row]) {
      seatMap[row] = []
    }

    seatMap[row].push(seat)
  })

  // =============================
  // SORT ROW
  // =============================
  const rows = Object.keys(seatMap)
    .map(Number)
    .sort((a, b) => a - b)

  return (
    <div className="p-6 bg-slate-50 rounded-xl border border-slate-100 w-fit mx-auto">
      <div className="text-center mb-8">
        <span className="text-[10px] font-black text-slate-400 uppercase tracking-[0.2em] bg-white px-3 py-1 rounded-full border">
          Aircraft Front
        </span>
      </div>

      <TooltipProvider delayDuration={0}>
        <div
          className="grid gap-y-4 gap-x-2 items-center"
          style={{
            gridTemplateColumns: `repeat(${total_columns + 1}, minmax(0, 1fr))`
          }}
        >
          {rows.map((row) => {
            const rowSeats = seatMap[row].sort((a, b) =>
              a.seat_code.localeCompare(b.seat_code)
            )

            return (
              <React.Fragment key={`row-${row}`}>
                {rowSeats.map((seat: any, colIndex: number) => (
                  <React.Fragment key={seat.id}>
                    
                    {/* AISLE */}
                    {colIndex === aisleIndex && (
                      <div className="flex justify-center items-center text-[10px] font-bold text-slate-300">
                        {row}
                      </div>
                    )}

                    <Tooltip>
                      <TooltipTrigger asChild>
                        <div className="relative group">
                          {(() => {
                            const type = seat.class_type?.toLowerCase()
                            const isSelected = selectedSeats.includes(seat.seat_code)

                            let borderClass = "border-slate-200"
                            let activeBgClass = "data-[state=checked]:bg-slate-600"

                            if (type === "first") {
                              borderClass = "border-amber-400"
                              activeBgClass = "data-[state=checked]:bg-amber-500"
                            } 
                            else if (type === "business") {
                              borderClass = "border-blue-400"
                              activeBgClass = "data-[state=checked]:bg-blue-500"
                            } 
                            else if (type === "economy") {
                              borderClass = "border-indigo-400"
                              activeBgClass = "data-[state=checked]:bg-indigo-600"
                            }

                            return (
                              <Checkbox
                                id={seat.seat_code}
                                checked={isSelected}
                                onCheckedChange={() => onSeatChange(seat.seat_code)}
                                disabled={!seat.is_available}
                                className={cn(
                                  "size-9 rounded-md border-2 transition-all",
                                  borderClass,
                                  activeBgClass,
                                  !seat.is_available &&
                                    "bg-slate-200 border-slate-200 opacity-40"
                                )}
                              />
                            )
                          })()}

                          <label className="absolute inset-0 flex items-center justify-center text-[9px] font-bold pointer-events-none">
                            {seat.seat_code}
                          </label>
                        </div>
                      </TooltipTrigger>

                      <TooltipContent>
                        <p className="text-xs">
                          {seat.class_type.toUpperCase()} -{" "}
                          {seat.is_available ? "Available" : "Booked"}
                        </p>
                      </TooltipContent>
                    </Tooltip>

                  </React.Fragment>
                ))}
              </React.Fragment>
            )
          })}
        </div>
      </TooltipProvider>

      <div className="mt-8 flex justify-center gap-4 border-t pt-4">
        <LegendItem color="bg-amber-400" border="border-amber-400" label="First" />
        <LegendItem color="bg-blue-400" border="border-blue-400" label="Business" />
        <LegendItem color="bg-indigo-400" border="border-indigo-400" label="Economy" />
        <LegendItem color="bg-slate-200" border="border-slate-200" label="Full" />
      </div>
    </div>
  )
}

function LegendItem({
  color,
  border,
  label
}: {
  color: string
  border: string
  label: string
}) {
  return (
    <div className="flex items-center gap-1.5">
      <div className={cn("size-3 rounded border", color, border)} />
      <span className="text-[10px] text-slate-500 font-bold uppercase tracking-tight">
        {label}
      </span>
    </div>
  )
}