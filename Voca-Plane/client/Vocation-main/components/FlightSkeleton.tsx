import { Card, CardContent } from "@/components/ui/card"

export function FlightSkeleton() {
  return (
    <div className="relative">
      <Card className="relative h-64 rounded-3xl overflow-hidden border-none bg-slate-100 dark:bg-slate-800 animate-pulse">
        <div className="absolute inset-0 p-6 flex flex-col justify-between">
          <div className="flex justify-between items-start">
            <div className="h-6 w-20 bg-slate-200 dark:bg-slate-700 rounded-full" />
            <div className="h-6 w-24 bg-slate-200 dark:bg-slate-700 rounded-full" />
          </div>
          <div className="space-y-3">
            <div className="h-10 w-3/4 bg-slate-200 dark:bg-slate-700 rounded-lg" />
            <div className="h-4 w-1/2 bg-slate-200 dark:bg-slate-700 rounded-lg" />
          </div>
        </div>
      </Card>
      <div className="absolute -bottom-6 left-1/2 -translate-x-1/2 w-[90%]">
        <Card className="rounded-2xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 shadow-sm animate-pulse">
          <CardContent className="p-4 flex justify-between items-center">
            <div className="space-y-2">
              <div className="h-2 w-16 bg-slate-100 dark:bg-slate-800 rounded" />
              <div className="h-6 w-32 bg-slate-100 dark:bg-slate-800 rounded" />
            </div>
            <div className="h-10 w-10 bg-slate-100 dark:bg-slate-800 rounded-xl" />
          </CardContent>
        </Card>
      </div>
    </div>
  )
}

export function FlightListSkeleton() {
  return (
    <div className="grid grid-cols-1 md:grid-cols-2 gap-y-20 gap-x-6">
      {Array.from({ length: 6 }).map((_, i) => (
        <FlightSkeleton key={i} />
      ))}
    </div>
  )
}
