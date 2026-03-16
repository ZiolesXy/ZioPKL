"use client"
import { useEffect } from "react"
import { useInView } from "react-intersection-observer"
import { useInfiniteQuery } from "@tanstack/react-query"
import { CardTicket } from "@/components/CardTicket"
import { getFlights } from "@/lib/api/FlightApi"
import { Flight } from "@/lib/type/flight";

export function FlightList({ initialData }: { initialData: any }) {
  const { ref, inView } = useInView()

  const { data, fetchNextPage, hasNextPage, isFetchingNextPage, status } = useInfiniteQuery({
    queryKey: ["flights"],
    queryFn: ({ pageParam }) => getFlights({ pageParam }),
    initialPageParam: 1,
    // Sinkronisasi data awal dari Server Component
    initialData: {
      pages: [initialData],
      pageParams: [1],
    },
    getNextPageParam: (lastPage) => {
      const { page, total, limit } = lastPage.meta;
      return page * limit < total ? page + 1 : undefined;
    },
  })

  useEffect(() => {
    if (inView && hasNextPage && !isFetchingNextPage) {
      fetchNextPage()
    }
  }, [inView, hasNextPage, isFetchingNextPage, fetchNextPage])

  return (
    <div className="space-y-10">
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-y-20 gap-x-6">
        {data.pages.map((page) =>
          page.data.map((flight: Flight) => (
            <CardTicket key={flight.id} flight={flight} />
          ))
        )}
      </div>

      <div ref={ref} className="h-20 flex items-center justify-center">
        {isFetchingNextPage ? (
          <p className="text-indigo-600 font-medium">Loading more flights...</p>
        ) : hasNextPage ? (
          <p className="text-slate-400 text-sm">Scroll down to see more</p>
        ) : (
          <p className="text-slate-400 text-sm italic">All flights loaded.</p>
        )}
      </div>
    </div>
  )
}