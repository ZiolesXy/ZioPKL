import * as z from "zod"

export const addressSchema = z.object({
  label: z.string().min(2, "Label minimal 2 karakter (contoh: Rumah, Kantor)"),
  recipient_name: z.string().min(2, "Nama penerima minimal 2 karakter"),
  phone: z.string().min(10, "Nomor HP minimal 10 karakter"),
  address_line: z.string().min(5, "Alamat lengkap minimal 5 karakter"),
  city: z.string().min(2, "Kota minimal 2 karakter"),
  province: z.string().min(2, "Provinsi minimal 2 karakter"),
  postal_code: z.string().min(5, "Kode pos minimal 5 karakter"),
  is_primary: z.boolean(),
})

export type AddressFormValues = z.infer<typeof addressSchema>

export const ADDRESS_DEFAULT_VALUES: AddressFormValues = {
  label: "",
  recipient_name: "",
  phone: "",
  address_line: "",
  city: "",
  province: "",
  postal_code: "",
  is_primary: false,
}
