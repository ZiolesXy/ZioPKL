"use client" // Wajib jika menggunakan hooks

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
import { getUserAdmin } from "@/lib/api/MonitorApi"

export function UserDataTable() {
  const [users, setUsers] = useState<any[]>([])
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await getUserAdmin()
        setUsers(response.data) 
      } catch (error) {
        console.error("Gagal mengambil data user:", error)
      } finally {
        setIsLoading(false)
      }
    }
    fetchData()
  }, [])

  const getRoleBadge = (role: string) => {
    switch (role) {
      case 'super_admin': return <Badge className="bg-purple-600 font-bold text-white">Super Admin</Badge>;
      case 'admin': return <Badge className="bg-blue-600 font-bold text-white">Admin</Badge>;
      default: return <Badge className="bg-slate-200 font-bold text-white">User</Badge>;
    }
  }

  if (isLoading) return <p className="p-4 text-center">Memuat data pengguna...</p>

  return (
    <div className="rounded-md border">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead className="w-[80px]">ID</TableHead>
            <TableHead>Nama Pengguna</TableHead>
            <TableHead>Email</TableHead>
            <TableHead>Role</TableHead>
            <TableHead>Terdaftar Pada</TableHead>
            <TableHead className="text-right">Aksi</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {users.length > 0 ? (
            users.map((user) => (
              <TableRow key={user.id}>
                <TableCell className="font-medium">{user.id}</TableCell>
                <TableCell>{user.name}</TableCell>
                <TableCell>{user.email}</TableCell>
                <TableCell>{getRoleBadge(user.role)}</TableCell>
                <TableCell>
                  {new Date(user.created_at).toLocaleDateString('id-ID', {
                    day: '2-digit', month: 'short', year: 'numeric'
                  })}
                </TableCell>
                <TableCell className="text-right space-x-2">
                  <Button variant="ghost" size="sm">Edit</Button>
                  <Button variant="destructive" size="sm">Suspend</Button>
                </TableCell>
              </TableRow>
            ))
          ) : (
            <TableRow>
              <TableCell colSpan={6} className="text-center">Tidak ada data.</TableCell>
            </TableRow>
          )}
        </TableBody>
      </Table>
    </div>
  )
}