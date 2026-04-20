"use client"

import { Button } from "@/components/ui/button"

export default function GlobalError({
  error,
  reset,
}: {
  error: Error & { digest?: string }
  reset: () => void
}) {
  return (
    <html lang="id">
      <body className="flex min-h-screen flex-col items-center justify-center gap-6 p-6 text-center">
        <div className="space-y-2">
          <h1 className="text-4xl font-extrabold tracking-tight text-destructive">
            500
          </h1>
          <h2 className="text-2xl font-bold">Kesalahan Fatal Sistem</h2>
          <p className="text-muted-foreground max-w-sm mx-auto">
            Terjadi kesalahan kritis pada aplikasi. Mohon maaf atas
            ketidaknyamanan ini.
          </p>
        </div>
        {process.env.NODE_ENV === "development" && (
          <div className="p-4 bg-muted rounded-lg text-left text-xs font-mono max-w-lg overflow-auto border">
            {error.message}
          </div>
        )}
        <Button
          size="lg"
          onClick={() => reset()}
          className="bg-teal-600 hover:bg-teal-700"
        >
          Muat Ulang Aplikasi
        </Button>
      </body>
    </html>
  )
}
