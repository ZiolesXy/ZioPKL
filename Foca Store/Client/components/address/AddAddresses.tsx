"use client"

import { useState } from "react"
import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { Plus, Loader2, MapPin } from "lucide-react"
import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import { Form } from "@/components/ui/form"
import { notifySuccess, notifyError } from "@/lib/toast"
import { addAddresses, getAddresses, updateAddress } from "@/lib/api/user"
import { CreateAddressData } from "@/types/auth"
import {
  addressSchema,
  ADDRESS_DEFAULT_VALUES,
  type AddressFormValues,
} from "@/lib/validators/address"
import AddressFormFields from "@/components/address/AddressFormFields"

interface AddAddressesProps {
  onSuccess?: () => void
}

export default function AddAddresses({ onSuccess }: AddAddressesProps) {
  const [open, setOpen] = useState(false)
  const [isSubmitting, setIsSubmitting] = useState(false)

  const form = useForm<AddressFormValues>({
    resolver: zodResolver(addressSchema),
    defaultValues: ADDRESS_DEFAULT_VALUES,
  })

  async function onSubmit(values: AddressFormValues) {
    setIsSubmitting(true)
    try {
      // Jika input alamat baru sebagai primary, nonaktifkan primary lama
      if (values.is_primary) {
        const existingAddresses = await getAddresses()
        const currentPrimary = existingAddresses.find(
          (addr) => addr.is_primary
        )
        if (currentPrimary) {
          await updateAddress(currentPrimary.uid, { is_primary: false })
        }
      }

      await addAddresses(values as CreateAddressData)
      notifySuccess("Alamat ditambahkan", "Alamat baru berhasil disimpan.")
      setOpen(false)
      form.reset()
      if (onSuccess) onSuccess()
    } catch {
      notifyError("Gagal", "Terjadi kesalahan saat menyimpan alamat.")
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button className="bg-teal-600 hover:bg-teal-700 text-white rounded-xl shadow-md transition-all hover:text-red-500 hover:scale-105">
          <Plus className="mr-2 h-4 w-4" /> Tambah Alamat Baru
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[550px] rounded-3xl">
        <DialogHeader>
          <DialogTitle className="text-2xl font-bold text-teal-900 flex items-center gap-2">
            <MapPin className="text-teal-600" /> Tambah Alamat
          </DialogTitle>
          <DialogDescription>
            Lengkapi detail alamat pengiriman Anda di bawah ini.
          </DialogDescription>
        </DialogHeader>

        <Form {...form}>
          <form
            onSubmit={form.handleSubmit(onSubmit)}
            className="space-y-4 py-4"
          >
            <AddressFormFields form={form} />

            <DialogFooter className="pt-4">
              <Button
                type="button"
                variant="outline"
                onClick={() => setOpen(false)}
                className="rounded-xl"
              >
                Batal
              </Button>
              <Button
                type="submit"
                disabled={isSubmitting}
                className="bg-teal-600 hover:bg-teal-700 rounded-xl px-8"
              >
                {isSubmitting ? (
                  <>
                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                    Menyimpan...
                  </>
                ) : (
                  "Simpan Alamat"
                )}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  )
}

