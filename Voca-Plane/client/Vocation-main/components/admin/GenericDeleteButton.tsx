"use client"

import { useState } from "react"
import { Trash2 } from "lucide-react"
import { toast } from "sonner"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import { Button } from "@/components/ui/button"

interface GenericDeleteButtonProps {
  id: number | string
  name: string
  deleteApi: (id: any) => Promise<any>
  onSuccess?: () => void
}

export function GenericDeleteButton({
  id,
  name,
  deleteApi,
  onSuccess,
}: GenericDeleteButtonProps) {
  const [open, setOpen] = useState(false)
  const [isDeleting, setIsDeleting] = useState(false)

  const handleDelete = async () => {
    setIsDeleting(true)
    try {
      await deleteApi(id)
      toast.success(`${name} berhasil dihapus`)
      setOpen(false)
      if (onSuccess) onSuccess()
    } catch (error: any) {
      console.error("Gagal menghapus:", error)
      toast.error(`Gagal menghapus ${name}: ${error.message || "Terjadi kesalahan"}`)
    } finally {
      setIsDeleting(false)
    }
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button
          variant="ghost"
          size="sm"
          className="h-8 px-3 text-xs text-destructive hover:text-destructive hover:bg-destructive/10"
        >
          <Trash2 className="mr-2 h-3.5 w-3.5" />
          Hapus
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Konfirmasi Hapus</DialogTitle>
          <DialogDescription>
            Apakah Anda yakin ingin menghapus <strong>{name}</strong>? Tindakan ini tidak dapat dibatalkan.
          </DialogDescription>
        </DialogHeader>
        <DialogFooter className="gap-2 sm:gap-0">
          <Button
            variant="outline"
            onClick={() => setOpen(false)}
            disabled={isDeleting}
          >
            Batal
          </Button>
          <Button
            variant="destructive"
            onClick={handleDelete}
            disabled={isDeleting}
          >
            {isDeleting ? "Menghapus..." : "Ya, Hapus"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
