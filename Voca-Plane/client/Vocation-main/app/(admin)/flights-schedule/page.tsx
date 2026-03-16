"use client"

import { useState } from "react"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Search, ChevronRight, PlusCircle } from "lucide-react"
import { Button } from "@/components/ui/button"
import { FlightsTableData } from "@/components/admin/FlightsTableData"
import Link from "next/link"
import { Separator } from "@/components/ui/separator"
import { UpsertForm } from "@/components/admin/UpsertForm"
import { useFlightFields } from "@/hooks/useFlightFields"
import { createFlight } from "@/lib/api/FlightApi"

export default function FlightSchedulePage() {
  const [refreshKey, setRefreshKey] = useState(0)
  const { flightFields } = useFlightFields()

  const handleRefresh = () => {
    setRefreshKey((prev) => prev + 1)
  }

  return (
    <div className="space-y-6">
      <div className="flex flex-col space-y-2">
        <nav className="flex items-center space-x-2 text-sm text-muted-foreground">
          <Link href="/dashboard" className="hover:text-primary transition-colors">
            Dashboard
          </Link>
          <ChevronRight className="h-4 w-4" />
          <span className="font-medium text-foreground">Flights Schedule</span>
        </nav>

        <div className="flex items-center justify-between">
          <div>
            <h2 className="text-3xl font-bold tracking-tight">Flights Schedule Management</h2>
            <p className="text-muted-foreground">
              Monitor flights schedule for the flight booking system.
            </p>
          </div>
          
          <UpsertForm
            title="Flight"
            description="Tambahkan data penerbangan baru"
            fields={flightFields}
            columns={2}
            maxWidth="sm:max-w-[1100px]"
            triggerLabel="Add New Flight"
            triggerIcon={<PlusCircle className="h-4 w-4" />}
            onSubmit={async (data) => {
              const payload = {
                ...data,
                airline_id: Number(data.airline_id),
                origin_id: Number(data.origin_id),
                destination_id: Number(data.destination_id),
                total_rows: Number(data.total_rows),
                total_columns: Number(data.total_columns),
                total_seats: Number(data.total_rows) * Number(data.total_columns),
                class_count: data.class_prices?.length || 0,
                class_prices: data.class_prices?.map((cp: any) => ({
                  ...cp,
                  price: Number(cp.price)
                }))
              };
              await createFlight(payload);
            }}
            onSuccess={handleRefresh}
          />
        </div>
      </div>

      <Separator />

      <Card>
        <CardHeader>
          <div className="flex items-center gap-4">
            <div className="relative flex-1">
              <Search className="absolute left-2 top-2.5 h-4 w-4 text-muted-foreground" />
              <Input placeholder="Search flight number or city..." className="pl-8" />
            </div>
          </div>
        </CardHeader>
        <CardContent>
          <Tabs defaultValue="all">
            <TabsList>
              <TabsTrigger value="all">All Flights</TabsTrigger>
              <TabsTrigger value="active">Active</TabsTrigger>
              <TabsTrigger value="archived">Archived (Soft Deleted)</TabsTrigger>
            </TabsList>
            <TabsContent value="all">
              <FlightsTableData refreshKey={refreshKey} />
            </TabsContent>
          </Tabs>
        </CardContent>
      </Card>
    </div>
  )
}