"use client"

import { UseFormReturn } from "react-hook-form"
import { MapPin, Phone, User, Home, Building, Hash } from "lucide-react"
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form"
import { Input } from "@/components/ui/input"
import { Checkbox } from "@/components/ui/checkbox"
import type { AddressFormValues } from "@/lib/validators/address"

interface AddressFormFieldsProps {
  form: UseFormReturn<AddressFormValues>
}

export default function AddressFormFields({ form }: AddressFormFieldsProps) {
  return (
    <>
      <div className="grid grid-cols-2 gap-4">
        <FormField
          control={form.control}
          name="label"
          render={({ field }) => (
            <FormItem>
              <FormLabel className="flex items-center gap-2">
                <Home size={14} /> Label Alamat
              </FormLabel>
              <FormControl>
                <Input
                  placeholder="Rumah / Kantor"
                  {...field}
                  className="focus-visible:ring-teal-600"
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="recipient_name"
          render={({ field }) => (
            <FormItem>
              <FormLabel className="flex items-center gap-2">
                <User size={14} /> Nama Penerima
              </FormLabel>
              <FormControl>
                <Input
                  placeholder="Nama Lengkap"
                  {...field}
                  className="focus-visible:ring-teal-600"
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
      </div>

      <FormField
        control={form.control}
        name="phone"
        render={({ field }) => (
          <FormItem>
            <FormLabel className="flex items-center gap-2">
              <Phone size={14} /> Nomor HP Penerima
            </FormLabel>
            <FormControl>
              <Input
                placeholder="08xxxxxxxxxx"
                {...field}
                className="focus-visible:ring-teal-600"
              />
            </FormControl>
            <FormMessage />
          </FormItem>
        )}
      />

      <FormField
        control={form.control}
        name="address_line"
        render={({ field }) => (
          <FormItem>
            <FormLabel className="flex items-center gap-2">
              <MapPin size={14} /> Alamat Lengkap
            </FormLabel>
            <FormControl>
              <Input
                placeholder="Jl. Nama Jalan, No. Rumah, RT/RW"
                {...field}
                className="focus-visible:ring-teal-600"
              />
            </FormControl>
            <FormMessage />
          </FormItem>
        )}
      />

      <div className="grid grid-cols-2 gap-4">
        <FormField
          control={form.control}
          name="city"
          render={({ field }) => (
            <FormItem>
              <FormLabel className="flex items-center gap-2">
                <Building size={14} /> Kota
              </FormLabel>
              <FormControl>
                <Input
                  placeholder="Kota / Kabupaten"
                  {...field}
                  className="focus-visible:ring-teal-600"
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="province"
          render={({ field }) => (
            <FormItem>
              <FormLabel className="flex items-center gap-2">
                <MapPin size={14} /> Provinsi
              </FormLabel>
              <FormControl>
                <Input
                  placeholder="Provinsi"
                  {...field}
                  className="focus-visible:ring-teal-600"
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
      </div>

      <FormField
        control={form.control}
        name="postal_code"
        render={({ field }) => (
          <FormItem className="w-1/2">
            <FormLabel className="flex items-center gap-2">
              <Hash size={14} /> Kode Pos
            </FormLabel>
            <FormControl>
              <Input
                placeholder="12345"
                {...field}
                className="focus-visible:ring-teal-600"
              />
            </FormControl>
            <FormMessage />
          </FormItem>
        )}
      />

      <FormField
        control={form.control}
        name="is_primary"
        render={({ field }) => (
          <FormItem className="flex flex-row items-start space-x-3 space-y-0 rounded-2xl border p-4 shadow-sm bg-gray-50/50">
            <FormControl>
              <Checkbox
                checked={field.value}
                onCheckedChange={field.onChange}
                className="mt-1"
              />
            </FormControl>
            <div className="space-y-1 leading-none">
              <FormLabel className="font-semibold text-teal-900">
                Jadikan Alamat Utama
              </FormLabel>
              <p className="text-sm text-gray-500">
                Alamat ini akan digunakan secara default untuk setiap pesanan
                Anda.
              </p>
            </div>
          </FormItem>
        )}
      />
    </>
  )
}
