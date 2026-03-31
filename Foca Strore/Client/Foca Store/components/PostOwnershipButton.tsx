"use client"

import FormDialog from "@/components/FormDialog"
import { postJson } from "@/lib/api/resources/post"

export function PostOwnershipButton() {
  return (
    <FormDialog<{ user_id: number; barang_id: number }>
      triggerLabel="Add Ownership"
      title="Create Ownership"
      description="Buat relasi user dengan barang. Setelah submit, table akan refresh otomatis."
      submitLabel="Create"
      fields={[
        { name: "user_id", label: "User ID", type: "number", required: true },
        { name: "barang_id", label: "Barang ID", type: "number", required: true },
      ]}
      onSubmit={async (values) => {
        await postJson(`/user/${values.user_id}/barang/${values.barang_id}`, {})
      }}
    />
  )
}
