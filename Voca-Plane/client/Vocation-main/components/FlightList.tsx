"use client"
import { useInfiniteQuery } from "@tanstack/react-query"
import { VList } from "virtua"
import { CardTicket } from "@/components/CardTicket"
import { getFlights } from "@/lib/api/FlightApi"
import { Flight } from "@/lib/type/flight";

import { Plane } from "lucide-react"

export function FlightList({ initialData, filters = {} }: { initialData: any, filters?: any }) {
  const { data, fetchNextPage, hasNextPage, isFetchingNextPage, isLoading } = useInfiniteQuery({
    queryKey: ["flights", filters],
    queryFn: ({ pageParam }) => getFlights(pageParam as number, 10, filters),
    initialPageParam: 1,
    initialData: {
      pages: [initialData],
      pageParams: [1],
    },
    getNextPageParam: (lastPage: any) => {
      const { page, total, limit } = lastPage.meta;
      return page * limit < total ? page + 1 : undefined;
    },
  })

  // Perataan data (flattening) untuk virtualisasi
  const allFlights: Flight[] = data?.pages.flatMap((page: any) => page?.data || []) || [];

  if (!isLoading && allFlights.length === 0) {
    return (
      <div className="w-full h-[60vh] flex flex-col items-center justify-center space-y-6 animate-in fade-in zoom-in duration-500">
        <div className="relative">
          <div className="absolute inset-0 bg-primary/10 rounded-full blur-3xl" />
          <div className="relative size-24 rounded-full bg-slate-50 dark:bg-slate-800 border border-slate-100 dark:border-slate-700 flex items-center justify-center shadow-xl">
            <Plane className="size-10 text-slate-300 dark:text-slate-600 -rotate-45" />
          </div>
        </div>
        <div className="text-center space-y-2 max-w-sm mx-auto">
          <h3 className="text-xl font-bold text-slate-900 dark:text-white">No flights found</h3>
          <p className="text-sm text-slate-500 dark:text-slate-400 font-medium leading-relaxed">
            We couldn't find any flights matching your search. Try adjusting your filters or searching for a different route.
          </p>
        </div>
      </div>
    )
  }

  return (
    <div className="w-full h-[80vh] overflow-hidden">
      <VList
        data={allFlights}
        onScrollEnd={() => {
          if (hasNextPage && !isFetchingNextPage) {
            fetchNextPage();
          }
        }}
      >
        {(flight: Flight) => (
          <div key={flight.id} className="pb-24 px-1">
            <CardTicket flight={flight} />
          </div>
        )}
      </VList>

      <div className="h-20 flex items-center justify-center">
        {isFetchingNextPage ? (
          <p className="text-primary font-medium text-sm animate-pulse">Scanning the skies...</p>
        ) : !hasNextPage && allFlights.length > 0 ? (
          <p className="text-slate-400 text-xs italic">All available flights have been found.</p>
        ) : null}
      </div>
    </div>
  )
}