"use client"

import FormDialog from "@/components/FormDialog"
import { postJson } from "@/lib/api/resources/post"

export function PostBarangButton() {
  return (
    <FormDialog<{ name: string; role_id: number }>
      triggerLabel="Post Barang"
      title="Create Barang"
      description="Tambah Barang baru. Setelah submit, table akan refresh otomatis."
      submitLabel="Create"
      fields={[
        { name: "name", label: "Barang Name", required: true },
      ]}
      onSubmit={async (values) => {
        await postJson("/barangs", values)
      }}
    />
  )
}
