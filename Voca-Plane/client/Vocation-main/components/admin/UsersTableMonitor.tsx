"use client"

import { useEffect, useState } from "react"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { getUserAdmin, updateUserAdmin } from "@/lib/api/MonitorApi" 
import { UpsertForm, FormField } from "./UpsertForm" 
import { toast } from "sonner"
import { Trash2, ShieldCheck } from "lucide-react"

// Konfigurasi field hanya untuk ROLE
const roleFields: FormField[] = [
  { 
    name: "role", 
    label: "Privilege Level", 
    type: "select", 
    options: [
      { label: "Super Admin", value: "super_admin" },
      { label: "Admin", value: "admin" },
      { label: "User", value: "user" },
    ],
    required: true,
    placeholder: "Select new access level"
  },
]

export function UserDataTable() {
  const [users, setUsers] = useState<any[]>([])
  const [isLoading, setIsLoading] = useState(true)

  const fetchData = async () => {
    setIsLoading(true)
    try {
      const response = await getUserAdmin()
      setUsers(response.data)
    } catch (error) {
      console.error("Gagal mengambil data user:", error)
      toast.error("Gagal memuat data pengguna")
    } finally {
      setIsLoading(false)
    }
  }

  useEffect(() => {
    fetchData()
  }, [])

  // Fungsi khusus untuk merubah role
  const handleUpdateRole = async (formData: any) => {
    try {
      // Hanya ambil ID dan Role baru
      const { id, role } = formData
      
      if (!id) throw new Error("ID Pengguna tidak ditemukan")

      // Kirim payload minimalis ke API
      await updateUserAdmin(id.toString(), { role })
      
    } catch (error: any) {
      const errorMsg = error.response?.data?.message || error.message || "Gagal memperbarui hak akses"
      throw new Error(errorMsg)
    }
  }

  const getRoleBadge = (role: string) => {
    switch (role) {
      case 'super_admin': 
        return <Badge className="bg-purple-600 font-bold text-white border-none shadow-sm">Super Admin</Badge>
      case 'admin': 
        return <Badge className="bg-blue-600 font-bold text-white border-none shadow-sm">Admin</Badge>
      default: 
        return <Badge className="bg-slate-500 font-bold text-white border-none shadow-sm">User</Badge>
    }
  }

  if (isLoading) {
    return (
      <div className="flex flex-col items-center justify-center p-20 space-y-4">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
        <p className="text-[10px] font-black uppercase tracking-[0.2em] text-slate-400">Authenticating Nodes...</p>
      </div>
    )
  }

  return (
    <div className="rounded-[2.5rem] border border-slate-100 dark:border-white/5 bg-white dark:bg-slate-900/50 overflow-hidden shadow-sm">
      <Table>
        <TableHeader className="bg-slate-50/50 dark:bg-white/2">
          <TableRow className="hover:bg-transparent border-slate-100 dark:border-white/5">
            <TableHead className="w-[80px] font-black text-[10px] uppercase tracking-widest py-6 px-6">ID</TableHead>
            <TableHead className="font-black text-[10px] uppercase tracking-widest py-6">User Identity</TableHead>
            <TableHead className="font-black text-[10px] uppercase tracking-widest py-6">Current Role</TableHead>
            <TableHead className="font-black text-[10px] uppercase tracking-widest py-6">Registered</TableHead>
            <TableHead className="text-right font-black text-[10px] uppercase tracking-widest py-6 px-6">Management</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {users.length > 0 ? (
            users.map((user) => (
              <TableRow key={user.id} className="group border-slate-100 dark:border-white/5 hover:bg-slate-50/50 dark:hover:bg-white/2 transition-colors">
                <TableCell className="font-bold text-slate-400 px-6 italic">#{user.id}</TableCell>
                <TableCell>
                  <div className="flex flex-col">
                    <span className="font-bold text-slate-900 dark:text-white tracking-tight">{user.name}</span>
                    <span className="text-xs text-slate-400 font-medium">{user.email}</span>
                  </div>
                </TableCell>
                <TableCell>{getRoleBadge(user.role)}</TableCell>
                <TableCell className="text-slate-500 text-sm font-bold">
                   {new Date(user.created_at).toLocaleDateString('id-ID', {
                    day: '2-digit', month: 'short', year: 'numeric'
                  })}
                </TableCell>
                <TableCell className="text-right px-6">
                  <div className="flex justify-end items-center gap-3">
                    
                    {/* FORM EDIT ROLE */}
                    <UpsertForm
                      title="Access Role"
                      description={`Modifikasi tingkat otorisasi untuk ${user.name}.`}
                      fields={roleFields}
                      initialData={user} 
                      onSubmit={handleUpdateRole}
                      onSuccess={fetchData} 
                      triggerVariant="icon"
                      triggerIcon={<ShieldCheck className="h-4 w-4" />} // Icon berbeda khusus Role
                      maxWidth="sm:max-w-[450px]"
                    />

                    <Button 
                      variant="ghost" 
                      size="icon" 
                      className="h-10 w-10 text-slate-300 hover:text-red-500 hover:bg-red-50 dark:hover:bg-red-500/10 rounded-xl transition-all opacity-0 group-hover:opacity-100"
                    >
                      <Trash2 className="size-4" />
                    </Button>
                  </div>
                </TableCell>
              </TableRow>
            ))
          ) : (
            <TableRow>
              <TableCell colSpan={5} className="text-center py-20 text-slate-400 font-bold text-xs uppercase tracking-[0.2em]">
                Empty Directory
              </TableCell>
            </TableRow>
          )}
        </TableBody>
      </Table>
    </div>
  )
}