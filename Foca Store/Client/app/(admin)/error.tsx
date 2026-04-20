"use client"
import ErrorBoundary from "@/components/shared/ErrorBoundary"

export default function Error(props: {
  error: Error & { digest?: string }
  reset: () => void
}) {
  return (
    <ErrorBoundary
      {...props}
      title="Panel Admin Bermasalah"
      description="Gagal memuat data admin. Ini mungkin masalah koneksi atau sesi login Anda kadaluarsa."
      backUrl="/overview"
      backLabel="Kembali ke Kontrol"
    />
  )
}

