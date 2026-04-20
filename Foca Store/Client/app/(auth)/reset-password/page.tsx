"use client"
import { useState, useEffect, Suspense } from "react"
import { useRouter, useSearchParams } from "next/navigation"
import { Button } from "@/components/ui/button"
import { Field, FieldGroup, FieldLabel } from "@/components/ui/field"
import { Input } from "@/components/ui/input"
import {
  InputOTP,
  InputOTPGroup,
  InputOTPSeparator,
  InputOTPSlot,
} from "@/components/ui/input-otp"
import AuthHeroPanel from "@/components/auth/AuthHeroPanel"
import { resetPassword } from "@/lib/api/user"
import Link from "next/link"

function ResetPasswordForm() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const [email, setEmail] = useState("")
  const [otp, setOtp] = useState("")
  const [newPassword, setNewPassword] = useState("")
  const [confirmPassword, setConfirmPassword] = useState("")
  const [isPending, setIsPending] = useState(false)
  const [error, setError] = useState("")
  const [success, setSuccess] = useState(false)

  useEffect(() => {
    const emailParam = searchParams.get("email")
    if (emailParam) {
      setEmail(emailParam)
    }
  }, [searchParams])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setIsPending(true)
    setError("")

    if (newPassword !== confirmPassword) {
      setError("Konfirmasi password tidak cocok.")
      setIsPending(false)
      return
    }

    if (otp.length !== 6) {
      setError("Kode OTP harus 6 digit.")
      setIsPending(false)
      return
    }

    try {
      // lib/api/user.ts: resetPassword handles the clientApi call to /verify-otp
      const res = await resetPassword({ 
        email, 
        otp, 
        new_password: newPassword 
      })

      if (res) {
        setSuccess(true)
        setTimeout(() => {
          router.push("/login")
        }, 2000)
      } else {
        setError("Gagal mereset password. Mohon cek kembali kode OTP Anda.")
      }
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : "Terjadi kesalahan saat mereset password.")
    } finally {
      setIsPending(false)
    }
  }

  return (
    <div className="grid w-full lg:grid-cols-2 min-h-screen bg-white">
      <AuthHeroPanel />

      <div className="flex flex-col justify-center p-8 md:p-12">
        <div className="mx-auto w-full max-w-sm space-y-8">
          <div className="space-y-2 text-center lg:text-left">
            <h1 className="text-3xl font-bold tracking-tight text-slate-900">Atur Ulang Password</h1>
            <p className="text-muted-foreground text-sm leading-relaxed">
              Kami telah mengirimkan kode OTP ke <span className="font-semibold text-slate-900">{email || "email Anda"}</span>. Silakan masukkan kode tersebut dan buat password baru.
            </p>
          </div>

          {error && (
            <div className="bg-red-50 border border-red-200 text-red-600 p-3 rounded-lg text-sm flex items-center gap-3 animate-in fade-in slide-in-from-top-1">
              <div className="w-1.5 h-1.5 rounded-full bg-red-600 shrink-0" />
              {error}
            </div>
          )}

          {success && (
            <div className="bg-emerald-50 border border-emerald-200 text-emerald-600 p-3 rounded-lg text-sm flex items-center gap-3 animate-in fade-in slide-in-from-top-1">
              <div className="w-1.5 h-1.5 rounded-full bg-emerald-600 shrink-0" />
              Password berhasil diperbarui! Mengalihkan ke halaman Login...
            </div>
          )}

          <form onSubmit={handleSubmit} className="grid gap-6">
            <FieldGroup className="grid gap-6">
              <Field className="hidden"> 
                {/* Hidden email field to satisfy form expectations, but UI shows it in text above */}
                <Input
                  id="email"
                  type="email"
                  value={email}
                  readOnly
                  required
                />
              </Field>

              <div className="space-y-4">
                <FieldLabel className="text-[11px] font-bold uppercase tracking-widest text-muted-foreground/80">Kode 6-Digit OTP</FieldLabel>
                <div className="flex justify-center lg:justify-start">
                    <InputOTP 
                        maxLength={6} 
                        value={otp} 
                        onChange={setOtp}
                        disabled={isPending || success}
                        autoFocus
                    >
                        <InputOTPGroup className="gap-2">
                             <InputOTPSlot index={0} className="rounded-md border-slate-200 h-12 w-10 text-lg shadow-none focus:border-teal-500" />
                             <InputOTPSlot index={1} className="rounded-md border-slate-200 h-12 w-10 text-lg shadow-none focus:border-teal-500" />
                             <InputOTPSlot index={2} className="rounded-md border-slate-200 h-12 w-10 text-lg shadow-none focus:border-teal-500" />
                        </InputOTPGroup>
                        <InputOTPSeparator className="text-slate-300" />
                        <InputOTPGroup className="gap-2">
                             <InputOTPSlot index={3} className="rounded-md border-slate-200 h-12 w-10 text-lg shadow-none focus:border-teal-500" />
                             <InputOTPSlot index={4} className="rounded-md border-slate-200 h-12 w-10 text-lg shadow-none focus:border-teal-500" />
                             <InputOTPSlot index={5} className="rounded-md border-slate-200 h-12 w-10 text-lg shadow-none focus:border-teal-500" />
                        </InputOTPGroup>
                    </InputOTP>
                </div>
              </div>

              <div className="space-y-4 pt-2">
                <Field>
                    <FieldLabel htmlFor="new_password" className="text-[11px] font-bold uppercase tracking-widest text-muted-foreground/80">Password Baru</FieldLabel>
                    <Input
                    id="new_password"
                    type="password"
                    value={newPassword}
                    onChange={(e) => setNewPassword(e.target.value)}
                    required
                    placeholder="â€¢â€¢â€¢â€¢â€¢â€¢â€¢â€¢"
                    className="h-11 bg-slate-50/50 border-slate-200 shadow-none focus:bg-white transition-colors"
                    disabled={isPending || success}
                    />
                </Field>

                <Field>
                    <FieldLabel htmlFor="confirm_password" className="text-[11px] font-bold uppercase tracking-widest text-muted-foreground/80">Konfirmasi Password Baru</FieldLabel>
                    <Input
                    id="confirm_password"
                    type="password"
                    value={confirmPassword}
                    onChange={(e) => setConfirmPassword(e.target.value)}
                    required
                    placeholder="â€¢â€¢â€¢â€¢â€¢â€¢â€¢â€¢"
                    className="h-11 bg-slate-50/50 border-slate-200 shadow-none focus:bg-white transition-colors"
                    disabled={isPending || success}
                    />
                </Field>
              </div>

              <Button 
                type="submit" 
                disabled={isPending || success || otp.length !== 6} 
                className="w-full bg-teal-600 hover:bg-teal-700 h-12 text-base font-semibold shadow-xl shadow-teal-600/10 transition-all active:scale-[0.98] mt-2 group"
              >
                {isPending ? "Sedang Memproses..." : "Atur Ulang Sekarang"}
              </Button>
            </FieldGroup>
          </form>

          <div className="text-center text-sm pt-6 border-t border-slate-100">
            Tiba-tiba teringat?{" "}
            <Link href="/login" className="font-bold text-teal-600 hover:text-teal-700 underline underline-offset-4 decoration-2">
              Login ke Akun
            </Link>
          </div>
        </div>
      </div>
    </div>
  )
}

export default function ResetPasswordPage() {
    return (
        <Suspense fallback={null}>
            <ResetPasswordForm />
        </Suspense>
    )
}

