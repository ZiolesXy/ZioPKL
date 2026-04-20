"use client"
import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"
import {
  Field,
  FieldGroup,
  FieldLabel,
} from "@/components/ui/field"
import { Input } from "@/components/ui/input"
import Link from "next/link"
import { useState } from "react"
import { useRouter, useSearchParams } from "next/navigation"
import AuthHeroPanel from "@/components/auth/AuthHeroPanel"

// Helper: validasi redirect URL hanya path internal
function isValidRedirect(url: string): boolean {
  // Hanya terima path relatif (mulai /, tanpa //)
  return url.startsWith("/") && !url.startsWith("//")
}

export function LoginForm({ className, ...props }: React.ComponentProps<"div">) {
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  const [error, setError] = useState("")
  const [isPending, setIsPending] = useState(false)
  const router = useRouter()
  const searchParams = useSearchParams()

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError("")
    setIsPending(true)

    try {
      const res = await fetch("/api/auth/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email, password }),
      })

      const data = await res.json()

      if (!res.ok) {
        setError(data.message || "Login gagal")
        return
      }

      const user = data.data.user

      // Validasi redirect URL (lindungi dari open redirect)
      const next = searchParams.get("next")
      if (next && isValidRedirect(next)) {
        router.push(next)
        return
      }

      if (user.role === "Admin") {
        router.push("/overview")
        return
      }

      router.push("/")
    } catch {
      setError("Terjadi kesalahan jaringan")
    } finally {
      setIsPending(false)
    }
  }

  return (
    <div className={cn("grid w-full lg:grid-cols-2", className)} {...props}>
      <AuthHeroPanel />

      <div className="flex flex-col justify-center p-8 md:p-12">
        <div className="mx-auto w-full max-w-sm space-y-6">
          <div className="space-y-2 text-center lg:text-left">
            <h1 className="text-3xl font-bold">Login</h1>
            <p className="text-balance text-muted-foreground">
              Enter your email below to login to your account
            </p>
          </div>

          {error && (
            <div className="bg-red-100 text-red-700 p-3 rounded text-sm">
              {error}
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
                  required
                  autoComplete="email"
                />
              </Field>
              <Field>
                <div className="flex items-center justify-between">
                  <FieldLabel htmlFor="password">Password</FieldLabel>
                  <Link href="/forgot-password" className="text-sm underline">Forgot password?</Link>
                </div>
                <Input
                  id="password"
                  type="password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  required
                  autoComplete="current-password"
                />
              </Field>
              <Button type="submit" disabled={isPending} className="w-full bg-teal-600 hover:bg-teal-700">
                {isPending ? "Logging in..." : "Login"}
              </Button>
              <Button variant="outline" className="w-full">
                Login with Google
              </Button>
            </FieldGroup>
          </form>

          <div className="text-center text-sm">
            Don&apos;t have an account?{" "}
            <Link href="/register" className="underline">Sign up</Link>
          </div>
        </div>
      </div>
    </div>
  )
}
