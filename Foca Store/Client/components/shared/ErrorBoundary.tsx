"use client"

import { useEffect } from "react"
import { Button } from "@/components/ui/button"

interface ErrorBoundaryProps {
  error: Error & { digest?: string }
  reset: () => void
  title?: string
  description?: string
  backUrl?: string
  backLabel?: string
}

export default function ErrorBoundary({
  error,
  reset,
  title = "Terjadi kesalahan!",
  description,
  backUrl = "/",
  backLabel = "Kembali ke Beranda",
}: ErrorBoundaryProps) {
  useEffect(() => {
    console.error(error)
  }, [error])

  return (
    <div className="flex flex-col items-center justify-center min-h-[400px] p-6 text-center">
      <div className="bg-destructive/10 p-4 rounded-full mb-4">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
          strokeWidth={1.5}
          stroke="currentColor"
          className="w-8 h-8 text-destructive"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            d="M12 9v3.75m9-.75a9 9 0 11-18 0 9 9 0 0118 0zm-9 3.75h.008v.008H12v-.008z"
          />
        </svg>
      </div>
      <h2 className="text-2xl font-bold tracking-tight mb-2">{title}</h2>
      <p className="text-muted-foreground mb-6 max-w-md">
        {description ??
          error.message ??
          "Mohon maaf, halaman ini gagal dimuat. Silakan coba lagi nanti."}
      </p>
      <div className="flex gap-4">
        <Button
          variant="default"
          className="bg-teal-600 hover:bg-teal-700"
          onClick={() => reset()}
        >
          Coba Lagi
        </Button>
        <Button
          variant="outline"
          onClick={() => (window.location.href = backUrl)}
        >
          {backLabel}
        </Button>
      </div>
    </div>
  )
}
