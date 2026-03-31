"use client"

import FormDialog from "@/components/FormDialog"
import { postJson } from "@/lib/api/resources/post"

export function PostRoleButton() {
    return (
        <FormDialog<{ name: string }>
            triggerLabel="Post Role"
            title="Create Role"
            fields={[{ name: "name", label: "Role Name", required: true }]}
            onSubmit={async (values) => {
                await postJson("/roles", values)
            }}
        />)
}
