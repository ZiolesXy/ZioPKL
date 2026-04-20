"use client"
import { useState } from "react"
import { useRouter } from "next/navigation"
import { Button } from "@/components/ui/button"
import { Field, FieldGroup, FieldLabel } from "@/components/ui/field"
import { Input } from "@/components/ui/input"
import AuthHeroPanel from "@/components/auth/AuthHeroPanel"
import { forgotPassword } from "@/lib/api/user"
import Link from "next/link"

export default function ForgotPasswordPage() {
  const [email, setEmail] = useState("")
  const [isPending, setIsPending] = useState(false)
  const [error, setError] = useState("")
  const [success, setSuccess] = useState(false)
  const router = useRouter()

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setIsPending(true)
    setError("")

    try {
      // lib/api/user.ts: forgotPassword handles the clientApi call to /forgot-password
      const res = await forgotPassword(email)
      
      // If we got a response (usually backend returns 200 OK for valid request)
      if (res) {
        setSuccess(true)
        // Redirect after a short delay or immediately
        setTimeout(() => {
          router.push(`/reset-password?email=${encodeURIComponent(email)}`)
        }, 1500)
      } else {
        setError("Gagal mengirim OTP. Mohon cek kembali email Anda.")
      }
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : "Terjadi kesalahan saat memproses permintaan.")
    } finally {
      setIsPending(false)
    }
  }

  return (
    <div className="grid w-full lg:grid-cols-2 min-h-screen">
      <AuthHeroPanel />

      <div className="flex flex-col justify-center p-8 md:p-12">
        <div className="mx-auto w-full max-w-sm space-y-6">
          <div className="space-y-2 text-center lg:text-left">
            <h1 className="text-3xl font-bold">Lupa Password</h1>
            <p className="text-muted-foreground">
              Masukkan email Anda. Kami akan mengirimkan kode OTP 6 digit untuk mengatur ulang kata sandi Anda.
            </p>
          </div>

          {error && (
            <div className="bg-red-50 border border-red-200 text-red-600 p-3 rounded-lg text-sm flex items-center gap-2">
              <span className="w-1.5 h-1.5 rounded-full bg-red-600 shrink-0" />
              {error}
            </div>
          )}

          {success && (
            <div className="bg-emerald-50 border border-emerald-200 text-emerald-600 p-3 rounded-lg text-sm flex items-center gap-2">
              <span className="w-1.5 h-1.5 rounded-full bg-emerald-600 shrink-0" />
              OTP Berhasil dikirim! Mengalihkan...
            </div>
          )}

          <form onSubmit={handleSubmit} className="grid gap-4">
            <FieldGroup className="grid gap-4">
              <Field>
                <FieldLabel htmlFor="email">Email</FieldLabel>
                <Input
                  id="email"
                  type="email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  placeholder="name@example.com"
                  required
                  className="bg-white"
                  disabled={isPending || success}
                />
              </Field>
              <Button 
                type="submit" 
                disabled={isPending || success} 
                className="w-full bg-teal-600 hover:bg-teal-700 h-11 transition-all active:scale-[0.98]"
              >
                {isPending ? "Sedang Mengirim..." : "Kirim Kode OTP"}
              </Button>
            </FieldGroup>
          </form>

          <div className="text-center text-sm text-muted-foreground pt-4 border-t">
            Kembali ke halaman {" "}
            <Link href="/login" className="font-semibold text-teal-600 hover:text-teal-700 underline underline-offset-4">
              Login
            </Link>
          </div>
        </div>
      </div>
    </div>
  )
}

