"use client"

import { GenericEditDialog, type GenericEditField } from "@/components/admin/GenericEdit"
import { editProductAction } from "@/app/(admin)/overview/actions"
import { Product } from "@/types/product"

interface EditProductProps {
  product: Product
  open: boolean
  onOpenChange: (open: boolean) => void
  onSuccess: () => void
}

type ProductEditValues = {
  name: string
  description: string
  price: string
  stock: string
  image: File | string | null
}

const FIELDS: Array<GenericEditField<ProductEditValues>> = [
  {
    name: "image",
    label: "Foto Produk",
    type: "file",
    accept: "image/*",
    hint: "Kosongkan jika tidak ingin mengubah gambar.",
    required: false,
  },
  {
    name: "name",
    label: "Nama Produk",
    type: "text",
    placeholder: "Nama produk",
    required: true,
  },
  {
    name: "description",
    label: "Deskripsi",
    type: "textarea",
    placeholder: "Deskripsi produk...",
    required: true,
  },
  {
    name: "price",
    label: "Harga (IDR)",
    type: "number",
    placeholder: "0",
    min: 0,
    required: true,
  },
  {
    name: "stock",
    label: "Stok",
    type: "number",
    placeholder: "0",
    min: 0,
    required: true,
  },
]

export default function EditProduct({
  product,
  open,
  onOpenChange,
  onSuccess,
}: EditProductProps) {
  const initialValues: ProductEditValues = {
    name: product.name,
    description: product.description,
    price: String(product.price),
    stock: String(product.stock),
    image: product.image_url,
  }

  return (
    <GenericEditDialog<ProductEditValues>
      open={open}
      onOpenChange={onOpenChange}
      title="Edit Produk"
      description="Ubah detail produk di bawah ini. Klik simpan setelah selesai."
      initialValues={initialValues}
      fields={FIELDS}
      onSubmit={async (values) => {
        try {
          const formData = new FormData()
          formData.append("name", values.name)
          formData.append("description", values.description)
          formData.append("price", values.price)
          formData.append("stock", values.stock)
          if (values.image instanceof File) {
            formData.append("image", values.image)
          }
          const result = await editProductAction(product.id, formData)
          if (!result.success) {
            return { success: false, message: result.message ?? "Gagal update produk" }
          }
          onSuccess()
          return { success: true }
        } catch {
          return {
            success: false,
            message: "Terjadi kesalahan tidak dikenal.",
          }
        }
      }}
      onSuccess={onSuccess}
    />
  )
}

