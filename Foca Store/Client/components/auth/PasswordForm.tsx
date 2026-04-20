"use client"

import { useState } from "react"
import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import * as z from "zod"
import { Loader2, Save } from "lucide-react"
import { Button } from "@/components/ui/button"
import {
  Form,
  FormField,
  FormItem,
  FormLabel,
  FormControl,
  FormMessage,
} from "@/components/ui/form"
import { notifySuccess, notifyError } from "@/lib/toast"
import { changePassword } from "@/lib/api/user"
import PasswordInput from "@/components/auth/PasswordInput"

const passwordSchema = z
  .object({
    old_password: z.string().min(6, "Minimal 6 karakter"),
    new_password: z.string().min(6, "Minimal 6 karakter"),
    confirm_password: z.string().min(6, "Minimal 6 karakter"),
  })
  .refine((data) => data.new_password === data.confirm_password, {
    message: "Password tidak sama",
    path: ["confirm_password"],
  })

export default function PasswordForm() {
  const [isSubmitting, setIsSubmitting] = useState(false)

  const form = useForm<z.infer<typeof passwordSchema>>({
    resolver: zodResolver(passwordSchema),
  })

  async function onSubmit(values: z.infer<typeof passwordSchema>) {
    setIsSubmitting(true)
    try {
      await changePassword(values)
      notifySuccess("Password berhasil diubah", "Silakan login ulang jika diperlukan.")
      form.reset()
    } catch (error: unknown) {
      const message =
        error instanceof Error ? error.message : "Terjadi kesalahan tak terduga"
      notifyError("Gagal", message)
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
        <FormField
          control={form.control}
          name="old_password"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Password Lama</FormLabel>
              <FormControl>
                <PasswordInput {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="new_password"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Password Baru</FormLabel>
              <FormControl>
                <PasswordInput {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="confirm_password"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Konfirmasi Password</FormLabel>
              <FormControl>
                <PasswordInput {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <Button
          type="submit"
          className="bg-teal-600 hover:bg-teal-700 w-full"
          disabled={isSubmitting}
        >
          {isSubmitting ? (
            <Loader2 className="animate-spin" />
          ) : (
            <>
              <Save className="mr-2 h-4 w-4" /> Simpan Password
            </>
          )}
        </Button>
      </form>
    </Form>
  )
}
