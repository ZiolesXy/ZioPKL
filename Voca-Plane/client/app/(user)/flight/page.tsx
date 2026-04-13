import { Search, Filter, ArrowUpDown } from "lucide-react"
import { Button } from "@/components/ui/button"
import { InputGroupInput } from "@/components/ui/input-group"
import { Checkbox } from "@/components/ui/checkbox"
import { Slider } from "@/components/ui/slider"
import { Label } from "@/components/ui/label"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { getFlights } from "@/lib/api/FlightApi"
import { FlightList } from "@/components/FlightList"
import Link from "next/link"

import { Suspense } from "react"
import { FlightListSkeleton } from "@/components/FlightSkeleton"

async function FlightResults({ params }: { params: any }) {
  const data = await getFlights(1, 10, params);
  return <FlightList initialData={data} filters={params} />;
}

async function FlightPage({ searchParams }: { searchParams: Promise<any> }) {
  const params = await searchParams;

  return (
    <div className="container mx-auto py-8 px-4 max-w-7xl min-h-screen">
      {/* Search Header - Simplified */}
      <div className="bg-white dark:bg-slate-900 p-5 rounded-2xl shadow-sm border border-slate-100 dark:border-slate-800 mb-8">
          <form method="GET" className="grid grid-cols-1 md:grid-cols-3 gap-4 group">
            {Object.entries(params as Record<string, string | string[] | undefined>).map(([key, value]) => {
              if (key === "origin" || key === "destination" || !value) return null;
              const stringValue = Array.isArray(value) ? value[0] : value;
              return <input key={key} type="hidden" name={key} value={stringValue} />;
            })}
            <div className="space-y-1.5">
              <label className="text-[10px] font-black text-slate-400 uppercase tracking-[0.2em] ml-1">Origin</label>
              <InputGroupInput 
                name="origin" 
                defaultValue={params.origin || ""} 
                placeholder="Where from?" 
                className="rounded-xl h-12 border-slate-200 focus:border-primary transition-all bg-slate-50/50" 
              />
            </div>
            <div className="space-y-1.5">
              <label className="text-[10px] font-black text-slate-400 uppercase tracking-[0.2em] ml-1">Destination</label>
              <InputGroupInput 
                name="destination" 
                defaultValue={params.destination || ""} 
                placeholder="Where to?" 
                className="rounded-xl h-12 border-slate-200 focus:border-primary transition-all bg-slate-50/50" 
              />
            </div>
            <div className="flex items-end">
              <Button type="submit" className="w-full bg-primary hover:bg-primary/90 text-white h-12 rounded-xl font-bold transition-transform active:scale-[0.98] shadow-lg shadow-primary/20">
                <Search className="mr-2 size-5" /> Search Flights
              </Button>
            </div>
          </form>
      </div>

      <div className="flex flex-col lg:flex-row gap-10">
        {/* Simplified Filters for Performance */}
        <aside className="hidden lg:block w-72 shrink-0">
          <div className="sticky top-24 space-y-10">
            <div className="flex items-center justify-between">
              <h3 className="font-bold text-xl flex items-center gap-2">
                <Filter className="size-5 text-primary" /> Filters
              </h3>
              <Link href="/flight" className="text-xs text-primary font-bold hover:underline transition-opacity">Reset All</Link>
            </div>

            <div className="space-y-8">
              <section>
                <p className="text-[10px] font-black text-slate-400 uppercase tracking-[0.2em] mb-5">Airlines</p>
                <div className="space-y-4">
                  {["Garuda Indonesia", "AirAsia", "Lion Air"].map((airline) => {
                    const isActive = params.airline === airline;
                    return (
                      <Link 
                        key={airline} 
                        href={`?${new URLSearchParams({...params, airline: isActive ? undefined : airline}).toString()}`}
                        className={`flex items-center gap-3 text-sm font-medium transition-all p-2 rounded-xl hover:bg-slate-50 dark:hover:bg-slate-800 ${isActive ? 'text-primary bg-primary/5' : 'text-slate-600 dark:text-slate-400'}`}
                        scroll={false}
                      >
                        <div className={`size-5 rounded-md border-2 flex items-center justify-center transition-colors ${isActive ? 'bg-primary border-primary' : 'border-slate-200'}`}>
                          {isActive && <div className="size-2 bg-white rounded-full" />}
                        </div>
                        {airline}
                      </Link>
                    );
                  })}
                </div>
              </section>

              <section>
                <p className="text-[10px] font-black text-slate-400 uppercase tracking-[0.2em] mb-5">Cabin Class</p>
                <div className="grid grid-cols-2 gap-2">
                  {["Economy", "Business", "First"].map((cls) => {
                    const isActive = params.class === cls;
                    return (
                      <Link
                        key={cls}
                        href={`?${new URLSearchParams({...params, class: isActive ? undefined : cls}).toString()}`}
                        className={`px-3 py-2.5 rounded-xl text-xs font-bold text-center transition-all ${isActive ? 'bg-primary text-white shadow-md shadow-primary/20' : 'bg-slate-50 dark:bg-slate-800 text-slate-500 hover:bg-slate-100'}`}
                        scroll={false}
                      >
                        {cls}
                      </Link>
                    );
                  })}
                </div>
              </section>
            </div>
          </div>
        </aside>

        <main className="flex-1 min-w-0">
          <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center mb-8 gap-4">
            <div>
              <h1 className="text-2xl font-bold tracking-tight text-slate-900 dark:text-white">Skybound Results</h1>
              <p className="text-sm text-slate-500 font-medium mt-1">Found amazing flights for your next adventure</p>
            </div>
            
            <div className="flex items-center gap-2 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 p-1.5 rounded-2xl shadow-sm">
              <ArrowUpDown className="size-4 text-slate-400 ml-2" />
              <Select defaultValue={params.sort || "cheapest"}>
                <SelectTrigger className="w-40 border-none bg-transparent focus:ring-0 shadow-none font-bold text-sm h-9">
                  <SelectValue placeholder="Sort by" />
                </SelectTrigger>
                <SelectContent className="rounded-2xl border-slate-100 shadow-xl">
                  <SelectItem value="cheapest" className="rounded-xl">Lowest Price</SelectItem>
                  <SelectItem value="highest" className="rounded-xl">Highest Price</SelectItem>
                  <SelectItem value="earliest" className="rounded-xl">Earliest</SelectItem>
                  <SelectItem value="latest" className="rounded-xl">Latest</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>

          <Suspense fallback={<FlightListSkeleton />}>
             <FlightResults params={params} />
          </Suspense>
        </main>
      </div>
    </div>
  )
}

export default FlightPage