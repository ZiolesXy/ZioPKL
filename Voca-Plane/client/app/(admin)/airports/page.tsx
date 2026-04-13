"use client"
import { useState } from "react"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Search, ChevronRight, PlusCircle } from "lucide-react"
import { Button } from "@/components/ui/button"
import Link from "next/link"
import { Separator } from "@/components/ui/separator"
import { AirportsTableData } from "@/components/admin/AirportTableData"
import { UpsertForm, FormField } from "@/components/admin/UpsertForm"
import { createAirport } from "@/lib/api/AirportApi"

export default function AirportMonitoringPage() {
  const [refreshKey, setRefreshKey] = useState(0)

  const airportFields: FormField[] = [
    { name: "code", label: "Kode Bandara", placeholder: "Contoh: CGK", required: true },
    { name: "name", label: "Nama Bandara", placeholder: "Contoh: Soekarno-Hatta", required: true },
    { name: "city", label: "Kota", placeholder: "Contoh: Jakarta", required: true },
  ]

  return (
    <div className="space-y-6">
      <div className="flex flex-col space-y-2">
        <nav className="flex items-center space-x-2 text-sm text-muted-foreground">
          <Link href="/dashboard" className="hover:text-primary transition-colors">
            Dashboard
          </Link>
          <ChevronRight className="h-4 w-4" />
          <span className="font-medium text-foreground">Airport Monitoring</span>
        </nav>

        <div className="flex items-center justify-between">
          <div>
            <h2 className="text-3xl font-bold tracking-tight">Airport Monitoring</h2>
            <p className="text-muted-foreground">
              Monitor Airport for the flight booking system.
            </p>
          </div>
          <UpsertForm
            title="Airport"
            description="Initialize a new airport node in the global aviation network."
            fields={airportFields}
            triggerLabel="Add New Airport"
            triggerIcon={<PlusCircle className="h-4 w-4" />}
            onSubmit={createAirport}
            onSuccess={() => setRefreshKey(prev => prev + 1)}
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
          <AirportsTableData refreshKey={refreshKey} />
        </CardContent>
      </Card>
    </div>
  )
}