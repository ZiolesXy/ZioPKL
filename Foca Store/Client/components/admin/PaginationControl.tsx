"use client"

import { number } from "zod" // eslint-disable-line @typescript-eslint/no-unused-vars
import { Button } from "@/components/ui/button"

export default function PaginationControl({
    page,
    totalPages,
    onPageChange,
}: {
    page: number
    totalPages: number
    onPageChange: (page: number) => void
}) {
    return (
        <div className="flex justify-center gap-2 mt-6">
           <Button
        variant="outline"
        disabled={page === 1}
        onClick={() => onPageChange(page - 1)}
      >
        Prev
      </Button>

      <span className="px-3 py-2 text-sm">
        Halaman {page} / {totalPages}
      </span>

      <Button
        variant="outline"
        disabled={page === totalPages}
        onClick={() => onPageChange(page + 1)}
      >
        Next
      </Button>
        </div>
    )
}


