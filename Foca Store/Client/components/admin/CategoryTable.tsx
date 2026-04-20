"use client"

import { useState } from "react"
import { Edit, Trash2, Layers, Plus } from "lucide-react"
import { useRouter } from "next/navigation"
import { Button } from "@/components/ui/button"
import Image from "next/image"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { GenericDelete } from "@/components/admin/GenericDelete"
import { GenericEditDialog, type GenericEditField } from "@/components/admin/GenericEdit"
import type { Category } from "@/types/category"
import { deleteCategory, editCategory, addCategory } from "@/lib/api/category"

export default function CategoryTable({ categories }: { categories: Category[] }) {
  
  const [editItem, setEditItem] = useState<Category | null>(null)
  const [createOpen, setCreateOpen] = useState(false)
  const router = useRouter()

  type CategoryFormValues = {
    name: string
    icon: File | null
  }
  type FormResult =
    | { success: true }
    | { success: false; message: string }

  const fields: Array<GenericEditField<CategoryFormValues>> = [
    {
      name: "name",
      label: "Nama Kategori",
      type: "text",
      required: true,
    },
    {
      name: "icon",
      label: "Icon",
      type: "file",
      accept: "image/*",
      hint: "Format gambar: .png, .jpg, .jpeg, .webp",
    },
  ]
  return (
    <div className="p-8 space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Manajemen Kategori</h1>
          <p className="text-muted-foreground">Tambah, edit, atau hapus kategori produk.</p>
        </div>
        <Button
          className="bg-teal-600 hover:bg-teal-700"
          onClick={() => setCreateOpen(true)}
        >
          <Plus className="mr-2 h-4 w-4" /> Tambah Kategori
        </Button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
            <CardTitle className="text-sm font-medium">Total Kategori</CardTitle>
            <Layers className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{categories.length}</div>
          </CardContent>
        </Card>
      </div>

      <Card className="shadow-sm">
        <CardHeader>
          <CardTitle>Daftar Kategori</CardTitle>
          <CardDescription>Menampilkan semua kategori yang aktif di toko.</CardDescription>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>ID</TableHead>
                <TableHead>Ikon</TableHead>
                <TableHead>Nama Kategori</TableHead>
                <TableHead>Jumlah Produk</TableHead>
                <TableHead className="text-center">Aksi</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {categories.map((category) => (
                <TableRow key={category.id}>
                  <TableCell>{category.id}</TableCell>
                  <TableCell>
                  {category.icon_url ? (
                      <Image
                        src={category.icon_url}
                        alt={category.name}
                        width={32}
                        height={32}
                        className="w-8 h-8 object-contain"
                      />
                    ) : (
                      <div className="w-8 h-8 bg-gray-200 rounded" />
                    )}
                  </TableCell>
                  <TableCell>{category.name}</TableCell>
                  <TableCell>{category.product_count}</TableCell>
                  <TableCell className="text-center space-x-2">
                    <Button
                      variant="outline"
                      size="icon"
                      onClick={() => setEditItem(category)}
                    >
                      <Edit size={14} />
                    </Button>
                    <GenericDelete
                      id={category.id}
                      name={category.name}
                      entityName="Kategori"
                      onDelete={async (id) => {
                        await deleteCategory(Number(id))
                        router.refresh()
                      }}
                      trigger={
                        <Button variant="outline" size="icon" className="text-red-600">
                          <Trash2 size={14} />
                        </Button>
                      }
                    />
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </CardContent>
      </Card>

      {editItem && (
        <GenericEditDialog<CategoryFormValues>
          open={!!editItem}
          onOpenChange={(open) => {
            if (!open) setEditItem(null)
          }}
          title="Edit Kategori"
          description="Ubah detail kategori di bawah ini. Klik simpan setelah selesai."
          initialValues={{ name: editItem.name, icon: null }}
          fields={fields}
          onSubmit={async (values) => {
            try {
              const formData = new FormData()
              formData.append("name", values.name)
              if (values.icon) {
                formData.append("icon", values.icon)
              }
              await editCategory(editItem.id, formData)
              return { success: true }
            } catch (e: unknown) {
              return {
                success: false,
                message: e instanceof Error ? e.message : "Gagal melakukan operasi kategori",
              }
            }
          }}
          onSuccess={() => {
            setEditItem(null)
            router.refresh()
          }}
        />
      )}
      {createOpen && (
        <GenericEditDialog<CategoryFormValues>
          open={createOpen}
          onOpenChange={setCreateOpen}
          title="Tambah Kategori"
          description="Isi detail kategori baru."
          initialValues={{ name: "", icon: null }}
          fields={fields}
          onSubmit={async (values): Promise<FormResult> => {
            try {
              await addCategory({ name: values.name, icon: values.icon! })
              return { success: true }
            } catch (e: unknown) {
              return {
                success: false,
                message: e instanceof Error ? e.message : "Gagal menambahkan kategori",
              }
            }
          }}
          onSuccess={() => {
            setCreateOpen(false)
            router.refresh()
          }}
        />
      )}
    </div>
  )
}

