"use client"

import { useState, useRef } from "react"
import { useRouter } from "next/navigation"
import { zodResolver } from "@hookform/resolvers/zod"
import { useForm } from "react-hook-form"
import * as z from "zod"
import {
  Camera,
  User,
  Mail,
  Phone,
  Lock,
  ShieldCheck,
  Loader2,
  Save,
  ChevronLeft,
} from "lucide-react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage, FormDescription } from "@/components/ui/form"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { Badge } from "@/components/ui/badge"
import { notifySuccess, notifyError } from "@/lib/toast"
import { updateUser } from "@/lib/api/user"
import { User as UserType } from "@/types/auth"
import PasswordForm from "./PasswordForm"
import AddressTab from "@/components/address/AddressTab"

const profileSchema = z.object({
  name: z.string().min(2, "Nama minimal 2 karakter"),
  telephone_number: z.string().min(10, "Nomor HP minimal 10 karakter"),
})


export default function EnhancedProfileForm({ user, backUrl }: { user: UserType; backUrl: string }) {

  const router = useRouter()
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [preview, setPreview] = useState<string | null>(null)
  const [file, setFile] = useState<File | null>(null)
  const fileInputRef = useRef<HTMLInputElement>(null)

  const form = useForm<z.infer<typeof profileSchema>>({
    resolver: zodResolver(profileSchema),
    defaultValues: {
      name: user?.name ?? "",
      telephone_number: user?.telephone_number ?? "",
    },
  })

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const selectedFile = e.target.files?.[0]
    if (selectedFile) {
      setFile(selectedFile)
      setPreview(URL.createObjectURL(selectedFile))
    }
  }

  async function onSubmit(values: z.infer<typeof profileSchema>) {
    setIsSubmitting(true)
    try {
      const formData = new FormData()
      formData.append("name", values.name)
      formData.append("telephone_number", values.telephone_number || "")
      if (file) formData.append("profile_image", file)

      await updateUser(formData)
      notifySuccess("Profil diperbarui", "Perubahan berhasil disimpan.")
      router.refresh()
    } catch {
      notifyError("Gagal", "Cek koneksi internet Anda.")
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <div className="max-w-300px mx-auto py-10 px-4 md:px-8">
      <div className="flex items-center gap-4 mb-8">
        <Button
          variant="ghost"
          size="icon"
          onClick={() => router.push(backUrl)}
          className="rounded-full bg-white shadow-sm border"
        >
          <ChevronLeft size={20} />
        </Button>
        <div>
          <h1 className="text-3xl font-bold text-gray-900">Pengaturan Profil</h1>
          <p className="text-gray-500">Kelola informasi akun dan pengaturan keamanan Anda</p>
        </div>
      </div>

      <div className="grid lg:grid-cols-12 gap-8">
        {/* Left Column: Profile Card */}
        <div className="lg:col-span-4 space-y-6">
          <Card className="border-none shadow-xl overflow-hidden rounded-3xl bg-teal-600 text-white relative">
            <div className="absolute top-0 right-0 p-4">
              <Badge variant="secondary" className="bg-white/20 text-white border-none backdrop-blur-md">
                Platinum
              </Badge>
            </div>
            <CardContent className="pt-12 pb-10 flex flex-col items-center text-center">
              <div className="relative mb-6">
                <Avatar className="h-32 w-32 border-4 border-white/30 shadow-2xl">
                  <AvatarImage src={preview || user.profile_image_url} className="object-cover" />
                  <AvatarFallback className="bg-white text-teal-600 text-4xl font-extrabold">
                    {user?.name?.[0]}
                  </AvatarFallback>
                </Avatar>
                <button
                  onClick={() => fileInputRef.current?.click()}
                  className="absolute bottom-0 right-0 rounded-full bg-white p-3 text-teal-600 shadow-xl hover:scale-110 transition-transform ring-4 ring-teal-600"
                >
                  <Camera size={18} />
                </button>
                <input type="file" ref={fileInputRef} className="hidden" accept="image/*" onChange={handleFileChange} />
              </div>

              <h2 className="text-2xl font-bold">{user?.name}</h2>
              <p className="text-teal-100/80 text-sm mt-1 mb-4">{user?.email}</p>

              <div className="flex items-center gap-2 px-4 py-2 bg-white/10 rounded-full text-xs font-semibold backdrop-blur-sm border border-white/10">
                <ShieldCheck size={14} className="text-teal-200" /> Akun Terverifikasi
              </div>
            </CardContent>

            {/* Decorative elements */}
            <div className="absolute -right-16 -bottom-16 h-48 w-48 rounded-full bg-white/10 blur-3xl" />
            <div className="absolute -left-16 -top-16 h-48 w-48 rounded-full bg-black/10 blur-3xl" />
          </Card>

          <Card className="border-none shadow-md rounded-3xl overflow-hidden bg-gray-50/50">
            <CardContent className="p-6">
              <h3 className="font-bold text-gray-900 mb-4">Ringkasan Akun</h3>
              <div className="space-y-4">
                <div className="flex justify-between items-center text-sm">
                  <span className="text-gray-500">Terdaftar sejak</span>
                  <span className="font-medium text-gray-900">
                    {user?.created_at ? new Date(user.created_at).toLocaleDateString("id-ID", { year: 'numeric', month: 'long', day: 'numeric' }) : "-"}
                  </span>
                </div>
                <div className="flex justify-between items-center text-sm">
                  <span className="text-gray-500">Status Akun</span>
                  <Badge className="bg-green-500/10 text-green-600 border-none shadow-none">Aktif</Badge>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Right Column: Tabs and Content */}
        <div className="lg:col-span-8">
          <Tabs defaultValue="general" className="w-full">
            <TabsList className="bg-gray-100/50 p-1 rounded-2xl w-full grid grid-cols-3 h-14 mb-8">
              <TabsTrigger value="general" className="rounded-xl data-[state=active]:bg-white data-[state=active]:shadow-sm data-[state=active]:text-teal-600 font-semibold transition-all">
                Biodata
              </TabsTrigger>
              <TabsTrigger value="address" className="rounded-xl data-[state=active]:bg-white data-[state=active]:shadow-sm data-[state=active]:text-teal-600 font-semibold transition-all">
                Alamat
              </TabsTrigger>
              <TabsTrigger value="security" className="rounded-xl data-[state=active]:bg-white data-[state=active]:shadow-sm data-[state=active]:text-teal-600 font-semibold transition-all">
                Keamanan
              </TabsTrigger>
            </TabsList>

            <TabsContent value="general" className="mt-0">
              <Card className="border-none shadow-lg rounded-3xl overflow-hidden">
                <CardHeader className="bg-gray-50/50 border-b px-8 py-6">
                  <CardTitle className="text-xl font-bold">Informasi Personal</CardTitle>
                  <CardDescription>Perbarui informasi profil Anda untuk pengalaman belanja yang lebih baik.</CardDescription>
                </CardHeader>
                <CardContent className="p-8">
                  <Form {...form}>
                    <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
                      <div className="grid md:grid-cols-2 gap-8">
                        <FormField
                          control={form.control}
                          name="name"
                          render={({ field }) => (
                            <FormItem>
                              <FormLabel className="font-semibold text-gray-700">Nama Lengkap</FormLabel>
                              <FormControl>
                                <div className="relative">
                                  <User className="absolute left-3 top-1/2 -translate-y-1/2 text-teal-600/50" size={18} />
                                  <Input {...field} className="pl-10 h-12 rounded-xl focus-visible:ring-teal-600 border-gray-200" />
                                </div>
                              </FormControl>
                              <FormMessage />
                            </FormItem>
                          )}
                        />
                        <FormItem>
                          <FormLabel className="font-semibold text-gray-700">Email Utama</FormLabel>
                          <FormControl>
                            <div className="relative">
                              <Mail className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" size={18} />
                              <Input value={user?.email} disabled className="pl-10 h-12 rounded-xl bg-gray-50 border-gray-100 text-gray-500" />
                            </div>
                          </FormControl>
                          <FormDescription>Email tidak dapat diubah demi keamanan akun.</FormDescription>
                        </FormItem>
                        <FormField
                          control={form.control}
                          name="telephone_number"
                          render={({ field }) => (
                            <FormItem>
                              <FormLabel className="font-semibold text-gray-700">Nomor Handphone</FormLabel>
                              <FormControl>
                                <div className="relative">
                                  <Phone className="absolute left-3 top-1/2 -translate-y-1/2 text-teal-600/50" size={18} />
                                  <Input
                                    {...field}
                                    placeholder="08xxxxxxx"
                                    className="pl-10 h-12 rounded-xl focus-visible:ring-teal-600 border-gray-200"
                                  />
                                </div>
                              </FormControl>
                              <FormMessage />
                            </FormItem>
                          )}
                        />
                      </div>
                      <div className="flex justify-end pt-4 border-t gap-4">
                        <Button type="button" variant="ghost" onClick={() => router.push(backUrl)} className="rounded-xl px-6">
                          Batalkan
                        </Button>
                        <Button type="submit" className="bg-teal-600 hover:bg-teal-700 rounded-xl px-10 h-12 font-bold shadow-lg shadow-teal-600/20" disabled={isSubmitting}>
                          {isSubmitting ? <Loader2 className="animate-spin" /> : <><Save className="mr-2 h-5 w-5" /> Simpan Perubahan</>}
                        </Button>
                      </div>
                    </form>
                  </Form>
                </CardContent>
              </Card>
            </TabsContent>

            <TabsContent value="address" className="mt-0">
              <AddressTab />
            </TabsContent>

            <TabsContent value="security" className="mt-0">
              <Card className="border-none shadow-lg rounded-3xl overflow-hidden">
                <CardHeader className="bg-gray-50/50 border-b px-8 py-6">
                  <div className="flex items-center gap-3">
                    <div className="bg-teal-100 p-2 rounded-xl">
                      <Lock className="text-teal-600" size={24} />
                    </div>
                    <div>
                      <CardTitle className="text-xl font-bold">Keamanan Akun</CardTitle>
                      <CardDescription>Ganti kata sandi Anda secara berkala untuk menjaga keamanan akun.</CardDescription>
                    </div>
                  </div>
                </CardHeader>
                <CardContent className="p-8">
                  <PasswordForm />
                </CardContent>
              </Card>
            </TabsContent>
          </Tabs>
        </div>
      </div>
    </div>
  )
}

