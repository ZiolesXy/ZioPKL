"use client"

import { Button } from "@/components/ui/button"
import { Trash2 } from "lucide-react"
import { deleteCartItem } from "@/lib/api/cart"

export default function DeleteCartItemButton({ itemId, onDeleted }: { itemId: number; onDeleted?: () => void }) {

  async function handleDelete() {

    try {
      await deleteCartItem(itemId)
      if (onDeleted) {
        onDeleted()
      } 
    } catch (error: unknown) {
      console.error(error)
      alert("Gagal hapus")
    }

  }

  return (
    <Button
      variant="destructive"
      size="icon"
      onClick={handleDelete}
    >
      <Trash2 className="w-4 h-4" />
    </Button>
  )
}

