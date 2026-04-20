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
import { useRouter } from "next/navigation"
import { toast } from "sonner"
import AuthHeroPanel from "@/components/auth/AuthHeroPanel"

export function RegisterForm({ className, ...props }: React.ComponentProps<"div">) {
  const router = useRouter()
  const [, setLoading] = useState(false)

  async function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault()
    setLoading(true)

    const formData = new FormData(e.currentTarget)

    const payload = {
      name: formData.get("name"),
      email: formData.get("email"),
      password: formData.get("password"),
    }

    try {
      const res = await fetch("/api/auth/register", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      })

      const data = await res.json()

      if (!res.ok) {
        throw new Error(data.message)
      }

      toast.success("Berhasil", {
        description: "Akun berhasil dibuat. Silakan login.",
      })

      router.push("/login")
    } catch (err: unknown) {
      const message =
        err instanceof Error ? err.message : "Terjadi kesalahan"
      toast.error("Gagal", {
        description: message,
      })
    } finally {
      setLoading(false)
    }
  }
  return (
    <div className={cn("grid w-full lg:grid-cols-2", className)} {...props}>
      <AuthHeroPanel />

      <div className="flex flex-col justify-center p-8 md:p-12">
        <div className="mx-auto w-full max-w-sm space-y-6">
          <div className="space-y-2 text-center lg:text-left">
            <h1 className="text-3xl font-bold">Register</h1>
            <p className="text-balance text-muted-foreground">
              Enter your email below to register to your account
            </p>
          </div>

          <form onSubmit={handleSubmit} className="grid gap-4">
            <FieldGroup className="grid gap-4">
              <Field>
                <FieldLabel htmlFor="email">Email</FieldLabel>
                <Input name="email" id="email" type="email" placeholder="m@example.com" required />
              </Field>
              <Field>
                <FieldLabel htmlFor="name">Name</FieldLabel>
                <Input name="name" id="name" type="text" placeholder="John Doe" required />
              </Field>
              <Field>
                <div className="flex items-center justify-between">
                  <FieldLabel htmlFor="password">Password</FieldLabel>
                  <Link href="#" className="text-sm underline">Forgot password?</Link>
                </div>
                <Input name="password" id="password" type="password" required />
              </Field>
              <Button type="submit" className="w-full bg-teal-600 hover:bg-teal-700">
                Register
              </Button>
              <Button variant="outline" className="w-full">
                Register with Google
              </Button>
            </FieldGroup>
          </form>

          <div className="text-center text-sm">
            Already have an account?{" "}
            <Link href="/login" className="underline hover:text-teal-600 transition-colors font-medium">Sign In</Link>
          </div>
        </div>
      </div>
    </div>
  )
}

