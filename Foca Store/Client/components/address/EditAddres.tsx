"use client"

import { useState } from "react"
import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { Pencil, Loader2 } from "lucide-react"
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
import { updateAddress, getAddresses } from "@/lib/api/user"
import { AddressData } from "@/types/auth"
import {
  addressSchema,
  type AddressFormValues,
} from "@/lib/validators/address"
import AddressFormFields from "@/components/address/AddressFormFields"

interface EditAddresProps {
  address: AddressData
  onSuccess?: () => void
}

export default function EditAddres({ address, onSuccess }: EditAddresProps) {
  const [open, setOpen] = useState(false)
  const [isSubmitting, setIsSubmitting] = useState(false)

  const form = useForm<AddressFormValues>({
    resolver: zodResolver(addressSchema),
    defaultValues: {
      label: address.label,
      recipient_name: address.recipient_name,
      phone: address.phone,
      address_line: address.address_line,
      city: address.city,
      province: address.province,
      postal_code: address.postal_code,
      is_primary: address.is_primary,
    },
  })

  async function onSubmit(values: AddressFormValues) {
    setIsSubmitting(true)
    try {
      // Jika memilih sebagai primary dan sebelumnya bukan primary
      if (values.is_primary && !address.is_primary) {
        const existingAddresses = await getAddresses()
        const currentPrimary = existingAddresses.find(
          (addr) => addr.is_primary
        )
        if (currentPrimary) {
          await updateAddress(currentPrimary.uid, { is_primary: false })
        }
      }

      await updateAddress(address.uid, values)
      notifySuccess("Alamat diperbarui", "Perubahan alamat berhasil disimpan.")
      setOpen(false)
      if (onSuccess) onSuccess()
    } catch {
      notifyError("Gagal", "Terjadi kesalahan saat memperbarui alamat.")
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <button className="text-gray-400 hover:text-teal-600 transition-colors">
          <Pencil size={16} />
        </button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[550px] rounded-3xl">
        <DialogHeader>
          <DialogTitle className="text-2xl font-bold text-teal-900 flex items-center gap-2">
            <Pencil className="text-teal-600" size={20} /> Edit Alamat
          </DialogTitle>
          <DialogDescription>
            Perbarui detail alamat pengiriman Anda di bawah ini.
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
                  "Simpan Perubahan"
                )}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  )
}

