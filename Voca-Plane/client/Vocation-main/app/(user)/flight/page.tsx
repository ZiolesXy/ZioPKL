import { Search, Filter, ArrowUpDown } from "lucide-react"
import { Button } from "@/components/ui/button"
import { InputGroupInput } from "@/components/ui/input-group"
import { Checkbox } from "@/components/ui/checkbox"
import { Slider } from "@/components/ui/slider"
import { Label } from "@/components/ui/label"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { getFlights } from "@/lib/api/FlightApi"
import { FlightList } from "@/components/FlightList" // Import komponen baru

async function FlightPage() {
  // Ambil data pertama di server untuk SEO dan kecepatan
  const initialFlights = await getFlights({ pageParam: 1 });

  return (
    <div className="container mx-auto py-10 px-4">
      {/* --- Bagian Search Bar Tetap Sama --- */}
      <div className="bg-white p-6 rounded-xl shadow-sm border mb-10 flex flex-wrap gap-4 items-end">
        <div className="flex-1 min-w-50">
          <label className="text-xs font-bold text-slate-500 uppercase">From</label>
          <InputGroupInput placeholder="Origin City..." className="mt-1" />
        </div>
        <div className="flex-1 min-w-50">
          <label className="text-xs font-bold text-slate-500 uppercase">To</label>
          <InputGroupInput placeholder="Destination City..." className="mt-1" />
        </div>
        <Button className="bg-indigo-600 h-11 px-8">
          <Search className="mr-2 size-4" /> Search
        </Button>
      </div>

      <div className="flex flex-col md:flex-row gap-8">
        {/* --- Bagian Aside Filters Tetap Sama --- */}
        <aside className="hidden md:block w-64 space-y-8">
          <div>
            <h3 className="font-bold mb-4 flex items-center gap-2">
              <Filter className="size-4" /> Filters
            </h3>

            <div className="space-y-4 mb-6">
              <p className="text-sm font-semibold text-slate-700">Airlines</p>
              <div className="space-y-2">
                {["Garuda Indonesia", "AirAsia", "Lion Air"].map((airline) => (
                  <div key={airline} className="flex items-center space-x-2">
                    <Checkbox id={airline} />
                    <label htmlFor={airline} className="text-sm text-slate-600 cursor-pointer">
                      {airline} [cite: 10]
                    </label>
                  </div>
                ))}
              </div>
            </div>

            <div className="space-y-4 mb-6">
              <p className="text-sm font-semibold text-slate-700">Price Range</p>
              <Slider defaultValue={[500000]} max={5000000} step={100000} />
              <div className="flex justify-between text-xs text-slate-500">
                <span>Rp 500rb</span>
                <span>Rp 5jt [cite: 51]</span>
              </div>
            </div>

            <div className="space-y-4">
              <p className="text-sm font-semibold text-slate-700">Class</p>
              <div className="space-y-2">
                {["Economy", "Business", "First"].map((cls) => (
                  <div key={cls} className="flex items-center space-x-2">
                    <Checkbox id={cls} />
                    <Label htmlFor={cls} className="font-normal">{cls} [cite: 47]</Label>
                  </div>
                ))}
              </div>
            </div>
            <div className="space-y-4 mt-6">
              <Button variant="outline" className="w-full">Reset Filters</Button>
            </div>
            <div className="space-y-4 mt-6">
              <Button variant="default" className="w-full">Apply Filters</Button>
            </div>
          </div>
        </aside>

        <main className="flex-1">
          <div className="flex justify-between items-center mb-6">
            <h1 className="text-xl font-bold">Showing {initialFlights.meta.total} Flights</h1>
            <div className="flex items-center gap-2">
              <ArrowUpDown className="size-4 text-slate-400" />
              <Select defaultValue="cheapest">
                <SelectTrigger className="w-45">
                  <SelectValue placeholder="Sort by" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="cheapest">Lowest Price</SelectItem>
                   <SelectItem value="highest">Highest Price</SelectItem>
                  <SelectItem value="earliest">Earliest Flight</SelectItem>
                  <SelectItem value="latest">Latest Flight</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>

          {/* GANTI GRID LAMA DENGAN INI */}
          <FlightList initialData={initialFlights} />
        </main>
      </div>
    </div>
  )
}

export default FlightPage