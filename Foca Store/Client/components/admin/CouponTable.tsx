"use client"
import { useState } from "react"
import { Edit, Trash2, TicketPercent, Plus, Calendar } from "lucide-react"
import { useRouter } from "next/navigation"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import {Card,CardContent,CardDescription,CardHeader,CardTitle,
} from "@/components/ui/card"
import {Table,TableBody,TableCell,TableHead,TableHeader,TableRow,} from "@/components/ui/table"
import { GenericDelete } from "@/components/admin/GenericDelete"
import { GenericEditDialog, type GenericEditField } from "@/components/admin/GenericEdit"
import { formatRupiah } from "@/lib/utils"
import { addCoupon, deleteCoupon, editCoupon } from "@/lib/api/coupon"
import { AvailableCoupon } from "@/types/coupon"

export default function CouponTable({ coupons }: { coupons: AvailableCoupon[] }) {
  const [editItem, setEditItem] = useState<AvailableCoupon | null>(null)
  const [createOpen, setCreateOpen] = useState(false)
  const router = useRouter()

  type CouponEditValues = {
    code: string
    type: "percentage" | "fixed" | string
    value: string | number
    quota: string | number
    minimum_purchase: string | number
    expires_at: string
    is_active: "true" | "false" | boolean
  }

  const fields: Array<GenericEditField<CouponEditValues>> = [
    {
      name: "code",
      label: "Kode Kupon",
      type: "text",
      placeholder: "Contoh: HEMAT50",
      required: true,
    },
    {
      name: "type",
      label: "Tipe Kupon",
      type: "select",
      options: [
        { value: "percentage", label: "Persentase" },
        { value: "fixed", label: "Nominal" },
      ],
      placeholder: "percentage / fixed",
      required: true,
    },
    {
      name: "value",
      label: "Nilai",
      type: "number",
      placeholder: "Masukan angka saja",
      required: true,
      min: 0,
    },
    {
      name: "quota",
      label: "Kuota",
      type: "number",
      placeholder: "0",
      required: true,
      min: 0,
    },
    {
      name: "minimum_purchase",
      label: "Minimum Pembelian",
      type: "number",
      placeholder: "0 jika tidak ada",
      required: true,
      min: 0,
    },
    {
      name: "expires_at",
      label: "Tanggal Kadaluarsa",
      type: "date",
      required: true,
    },
    {
      name: "is_active",
      label: "Status",
      type: "select",
      options: [
        { value: "true", label: "Aktif" },
        { value: "false", label: "Nonaktif" },
      ],
      required: true,
    },
  ]

  return (
    <div className="p-8 space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight text-slate-800">Manajemen Kupon</h1>
          <p className="text-muted-foreground">Kelola diskon dan promosi toko Anda.</p>
        </div>
        <Button className="bg-teal-600 hover:bg-teal-700" onClick={() => setCreateOpen(true)}>
          <Plus className="mr-2 h-4 w-4" /> Tambah Kupon
        </Button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <Card className="border-l-4 border-l-teal-500">
          <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
            <CardTitle className="text-sm font-medium">Total Kupon Aktif</CardTitle>
            <TicketPercent className="h-4 w-4 text-teal-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{coupons.filter((c) => c.is_active === true).length}</div>
          </CardContent>
        </Card>
      </div>

      {/* Table Section */}
      <Card className="shadow-sm border-none ring-1 ring-slate-200">
        <CardHeader>
          <CardTitle>Daftar Kupon</CardTitle>
          <CardDescription>Kupon yang dapat digunakan pelanggan di halaman checkout.</CardDescription>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow className="bg-slate-50/50">
                <TableHead>Kode</TableHead>
                <TableHead>Potongan</TableHead>
                <TableHead>Min. Belanja</TableHead>
                <TableHead>Quota</TableHead>
                <TableHead>Terpakai</TableHead>
                <TableHead>Berlaku Sampai</TableHead>
                <TableHead>Status</TableHead>
                <TableHead className="text-center">Aksi</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {coupons.map((coupon) => (
                <TableRow key={coupon.id} className="hover:bg-slate-50/50 transition-colors">
                  <TableCell className="font-bold font-mono tracking-wider">
                    {coupon.code}
                  </TableCell>
                  <TableCell>
                    {coupon.type === "percentage"
                      ? `${coupon.value}%`
                      : formatRupiah(coupon.value)}
                  </TableCell>
                  <TableCell>{formatRupiah(coupon.minimum_purchase)}</TableCell>
                  <TableCell>{coupon.quota}</TableCell>
                  <TableCell>{coupon.used_count}</TableCell>
                  <TableCell className="text-muted-foreground flex items-center gap-1.5">
                    <Calendar size={14} />
                    {new Date(coupon.expires_at).toLocaleDateString("id-ID")}
                  </TableCell>
                  <TableCell>
                    <Badge
                      variant={coupon.is_active === true ? "outline" : "destructive"}
                      className={coupon.is_active === true ? "bg-emerald-50 text-emerald-600 border-emerald-200" : ""}
                    >
                      {coupon.is_active === true ? "Aktif" : "Kedaluwarsa"}
                    </Badge>
                  </TableCell>
                  <TableCell className="text-center space-x-2">
                    <Button
                      variant="outline"
                      size="icon"
                      className="h-8 w-8 text-slate-600"
                      onClick={() => setEditItem(coupon)}
                    >
                      <Edit size={14} />
                    </Button>
                    <GenericDelete
                      id={coupon.id}
                      name={coupon.code}
                      entityName="Kupon"
                      onDelete={async (id) => {
                        await deleteCoupon(Number(id))
                        router.refresh()
                      }}
                      trigger={
                        <Button
                          variant="outline"
                          size="icon"
                          className="h-8 w-8 text-red-600 hover:bg-red-50 border-red-100"
                        >
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
        <GenericEditDialog<CouponEditValues>
          open={!!editItem}
          onOpenChange={(open) => !open && setEditItem(null)}
          title="Edit Detail Kupon"
          description="Perbarui informasi kupon diskon Anda."
          initialValues={{
            code: editItem.code,
            type: editItem.type,
            value: String(editItem.value),
            quota: String(editItem.quota),
            minimum_purchase: String(editItem.minimum_purchase),
            expires_at: editItem.expires_at,
            is_active: editItem.is_active ? "true" : "false",
          }}
          fields={fields}
          onSubmit={async (values) => {
            const parsed = parseCouponPayload(values)
            if (!parsed.success) return parsed
            try {
              await editCoupon(editItem.id, parsed.payload)
              return { success: true }
            } catch (e) {
              return { success: false, message: e instanceof Error ? e.message : "Gagal mengedit kupon" }
            }
          }}
          onSuccess={() => {
            setEditItem(null)
            router.refresh()
          }}
        />
      )}

      <GenericEditDialog<CouponEditValues>
        open={createOpen}
        onOpenChange={setCreateOpen}
        title="Tambah Kupon"
        description="Isi detail kupon baru."
        initialValues={{
          code: "",
          type: "percentage",
          value: "",
          quota: "",
          minimum_purchase: "",
          expires_at: "",
          is_active: "true",
        }}
        fields={fields}
        onSubmit={async (values) => {
          const parsed = parseCouponPayload(values)
          if (!parsed.success) return parsed
          try {
            await addCoupon(parsed.payload)
            return { success: true }
          } catch (e) {
            return { success: false, message: e instanceof Error ? e.message : "Gagal menambahkan kupon" }
          }
        }}
        onSuccess={() => {
          setCreateOpen(false)
          router.refresh()
        }}
      />
    </div>
  )
}

function parseCouponPayload(values: {
  code: string
  type: string
  value: string | number
  quota: string | number
  minimum_purchase: string | number
  expires_at: string
  is_active: string | boolean
}) {
  if (!values.expires_at) {
    return { success: false as const, message: "Tanggal kadaluarsa wajib diisi" }
  }
  return {
    success: true as const,
    payload: {
      code: values.code,
      type: values.type as "percentage" | "fixed",
      value: Number(values.value),
      quota: Number(values.quota),
      minimum_purchase: Number(values.minimum_purchase),
      expires_at: values.expires_at,
      is_active: values.is_active === true || values.is_active === "true",
    },
  }
}

