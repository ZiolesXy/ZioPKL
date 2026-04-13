"use client"

import { useState, useEffect, useMemo } from "react"
import { getAirlines } from "@/lib/api/AirlineApi"
import { getAirports } from "@/lib/api/AirportApi"
import { FormField } from "@/components/admin/UpsertForm"

export function useFlightFields() {
  const [airlineOptions, setAirlineOptions] = useState<{ label: string; value: number }[]>([])
  const [airportOptions, setAirportOptions] = useState<{ label: string; value: number }[]>([])
  const [isDataLoading, setIsDataLoading] = useState(true)

  const flightFields: FormField[] = useMemo(() => [
    { name: "airline_id", label: "Maskapai", type: "select", required: true, options: airlineOptions, placeholder: "Pilih Maskapai" },
    { name: "flight_number", label: "Nomor Penerbangan", required: true, placeholder: "Contoh: GA-123" },
    { name: "origin_id", label: "Asal", type: "select", required: true, options: airportOptions, placeholder: "Pilih Bandara Asal" },
    { name: "destination_id", label: "Tujuan", type: "select", required: true, options: airportOptions, placeholder: "Pilih Bandara Tujuan" },
    { name: "departure_time", label: "Waktu Berangkat", type: "datetime", required: true },
    { name: "arrival_time", label: "Waktu Tiba", type: "datetime", required: true },
    { name: "total_rows", label: "Baris", type: "number", required: true, placeholder: "Jumlah baris kursi" },
    { name: "total_columns", label: "Kolom", type: "number", required: true, placeholder: "Jumlah kolom kursi" },
    { 
      name: "class_prices", 
      label: "Pengaturan Harga Kelas", 
      type: "dynamic-list" 
    },
  ], [airlineOptions, airportOptions]);

  useEffect(() => {
    const loadDependencies = async () => {
      setIsDataLoading(true)
      try {
        const resAirlines = await getAirlines(1, 100); 
        const resAirports = await getAirports(1, 100);
        
        setAirlineOptions(resAirlines.data.map((a: any) => ({
          label: a.name,
          value: a.id,
        })));
        
        setAirportOptions(resAirports.data.map((a: any) => ({
          label: `${a.code} - ${a.city}`,
          value: a.id,
        })));
      } catch (error) {
        console.error("Gagal memuat data pendukung:", error);
      } finally {
        setIsDataLoading(false)
      }
    };
    loadDependencies();
  }, []);

  return { flightFields, isDataLoading }
}
