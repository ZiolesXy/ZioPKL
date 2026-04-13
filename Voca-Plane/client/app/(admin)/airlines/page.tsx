"use client"
import { useState } from "react"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Search, ChevronRight, PlusCircle } from "lucide-react"
import { Button } from "@/components/ui/button"
import { AirlinesTableData } from "@/components/admin/AirlinesTableData"
import Link from "next/link"
import { Separator } from "@/components/ui/separator"
import { UpsertForm, FormField } from "@/components/admin/UpsertForm"
import { createAirline } from "@/lib/api/AirlineApi"

export default function AirlineMonitoringPage() {
  const [refreshKey, setRefreshKey] = useState(0)

  const airlineFields: FormField[] = [
    { name: "name", label: "Nama Maskapai", placeholder: "Contoh: Garuda Indonesia", required: true },
    { name: "code", label: "Kode IATA", placeholder: "Contoh: GA", required: true },
    { 
      name: "logo", 
      label: "Logo Maskapai", 
      type: "file", 
      placeholder: "Pilih gambar logo" 
    },
  ]

  return (
    <div className="space-y-6">
      <div className="flex flex-col space-y-2">
        <nav className="flex items-center space-x-2 text-sm text-muted-foreground">
          <Link href="/dashboard" className="hover:text-primary transition-colors">
            Dashboard
          </Link>
          <ChevronRight className="h-4 w-4" />
          <span className="font-medium text-foreground">Airline Monitoring</span>
        </nav>

        <div className="flex items-center justify-between">
          <div>
            <h2 className="text-3xl font-bold tracking-tight">Airline Monitoring</h2>
            <p className="text-muted-foreground">
              Monitor Airline for the flight booking system.
            </p>
          </div>
          <UpsertForm
            title="Airline"
            description="Authorize a new carrier for the Voca-Plane fleet."
            fields={airlineFields}
            triggerLabel="Add New Airline"
            triggerIcon={<PlusCircle className="h-4 w-4" />}
            onSubmit={createAirline}
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
          <AirlinesTableData refreshKey={refreshKey} />
        </CardContent>
      </Card>
    </div>
  )
}