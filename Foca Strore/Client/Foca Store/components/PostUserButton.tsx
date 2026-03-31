"use client"

import FormDialog from "@/components/FormDialog"
import { postJson } from "@/lib/api/resources/post"

export function PostButton() {
  return (
    <FormDialog<{ name: string; role_id: number }>
      triggerLabel="Post User"
      title="Create User"
      description="Tambah user baru. Setelah submit, table akan refresh otomatis."
      submitLabel="Create"
      fields={[
        { name: "name", label: "Name", type: "text", required: true },
        { name: "role_id", label: "Role ID", type: "number", required: true },
      ]}
      onSubmit={async (values) => {
        await postJson("/users", values)
      }}
    />
  )
}
