"use client"

import { deleteResource } from "@/lib/api/resources/delete"
import { useTransition } from "react"
import { useRouter } from "next/navigation"
import { Button } from "./ui/button"

interface DeleteButtonProps {
  type: "user" | "barang" | "role" | "ownership"
  userId?: number
  barangId?: number
  userName?: string
  barangName?: string
  roleId?: number
  label?: string
  onSuccess?: () => void
}

export default function DeleteButton(props: DeleteButtonProps) {
  const [isPending, startTransition] = useTransition()
  const router = useRouter()

  const handleDelete = () => {
    if (!confirm("Sawit lu Aman?")) return

    startTransition(async () => {
      try {
        await deleteResource(props)
        props.onSuccess?.()
        router.refresh()
      } catch (err) {
        const message = err instanceof Error ? err.message : "Delete gagal"
        alert(message)
      }
    })
  }

  return (
    <Button
      onClick={handleDelete}
      disabled={isPending}
      variant="outline"
    >
      {props.label ?? "Delete"}
    </Button>
  )
}