"use client"
import React, { useState } from 'react'
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { 
  Database, 
  Layers, 
  Trash2, 
  RefreshCcw, 
  Image as ImageIcon, 
  FileWarning, 
  ChevronRight,
  Loader2,
  CloudOff
} from 'lucide-react'
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog"
import { Input } from "@/components/ui/input"
import { toast } from "sonner"
import { 
  ResetAllNoProCa, Seed, SeedProduct,
  SeedProductWithAssets, SeedProductWithoutAssets, SyncProductAssets,
  DeleteSyncAssets, MigrateDatabase, ResetDatabase, ResetDatabaseFull,
  DeleteCloudinaryAssets
} from "@/lib/api/system"

export default function SystemAdminPage() {
  const [loading, setLoading] = useState<string | null>(null)
  const [password, setPassword] = useState<string>("")

  const executeAction = async (id: string) => {
    if (!password) {
      toast.error("Password wajib diisi")
      return
    }
    setLoading(id)
    try {
      let res;
      switch (id) {
        case 'migrate': res = await MigrateDatabase(password); break;
        case 'seed-basic': res = await Seed(password); break;
        case 'seed-prod': res = await SeedProduct(password); break;
        case 'seed-img': res = await SeedProductWithAssets(password); break;
        case 'seed-no-img': res = await SeedProductWithoutAssets(password); break;
        case 'sync-asset': res = await SyncProductAssets(password); break;
        case 'reset-soft': res = await ResetAllNoProCa(password); break;
        case 'del-sync': res = await DeleteSyncAssets(password); break;
        case 'reset-basic': res = await ResetDatabase(password); break;
        case 'purge-cloud': res = await DeleteCloudinaryAssets(password); break;
        case 'reset-full': res = await ResetDatabaseFull(password); break;
      }
      
      const resAny = res as any
      if (resAny?.data?.success || resAny?.status === 200 || resAny?.status === 201) {
        toast.success(resAny?.data?.message || "Berhasil dieksekusi")
      } else if (resAny?.data) {
        toast.error(resAny?.data?.message || "Gagal mengeksekusi")
      }
    } catch (e) {
      toast.error("Terjadi kesalahan pada sistem")
    } finally {
      setLoading(null)
    }
  }

  return (
    <div className="p-8 max-w-7xl space-y-10 min-h-screen">
      <header className="space-y-4">
        <div>
          <h1 className="text-3xl font-extrabold text-slate-900 tracking-tight">System Engine</h1>
          <p className="text-slate-500">Utility untuk seeding, migrasi, dan pembersihan database.</p>
        </div>
        <div className="max-w-md">
          <Input 
            type="password" 
            placeholder="System Admin Password" 
            value={password}
            onChange={(e: React.ChangeEvent<HTMLInputElement>) => setPassword(e.target.value)}
            className="border-slate-300 shadow-sm"
          />
        </div>
      </header>

      {/* SECTION 1: DATABASE & STRUCTURE */}
      <section className="space-y-4">
        <div className="flex items-center gap-2 border-b pb-2">
          <Database className="w-5 h-5 text-indigo-600" />
          <h2 className="text-lg font-bold">Database & Migration</h2>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          <ActionCard 
            title="Migrate Database" 
            desc="Singkronisasi skema database dengan model Prisma/ORM terbaru."
            method="POST"
            icon={<Layers className="w-4 h-4" />}
            onRun={() => executeAction('migrate')}
            isLoading={loading === 'migrate'}
          />
          <ActionCard 
            title="Basic Seeder" 
            desc="Isi data dasar: User admin, kategori utama, alamat, dan kupon."
            method="POST"
            onRun={() => executeAction('seed-basic')}
            isLoading={loading === 'seed-basic'}
          />
          <ActionCard 
            title="Seeder + Product" 
            desc="Isi data dasar lengkap beserta list produk standar."
            method="POST"
            onRun={() => executeAction('seed-prod')}
            isLoading={loading === 'seed-prod'}
          />
        </div>
      </section>

      {/* SECTION 2: ASSET & SEEDER MANAGEMENT */}
      <section className="space-y-4">
        <div className="flex items-center gap-2 border-b pb-2">
          <ImageIcon className="w-5 h-5 text-teal-600" />
          <h2 className="text-lg font-bold">Asset & Product Seeding</h2>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          <ActionCard 
            title="Seed Product (Images)" 
            desc="Generate produk contoh yang sudah memiliki asset gambar."
            method="POST"
            onRun={() => executeAction('seed-img')}
            isLoading={loading === 'seed-img'}
          />
          <ActionCard 
            title="Seed Product (No Images)" 
            desc="Generate produk tanpa gambar untuk testing kecepatan load."
            method="POST"
            onRun={() => executeAction('seed-no-img')}
            isLoading={loading === 'seed-no-img'}
          />
          <ActionCard 
            title="Sync Assets Seed" 
            desc="Menghubungkan file asset yang menggantung ke data produk."
            method="POST"
            onRun={() => executeAction('sync-asset')}
            isLoading={loading === 'sync-asset'}
          />
        </div>
      </section>

      {/* SECTION 3: DANGER ZONE */}
      <section className="space-y-4 p-6 bg-red-50/50 border border-red-100 rounded-2xl">
        <div className="flex items-center gap-2 border-b border-red-200 pb-2">
          <FileWarning className="w-5 h-5 text-red-600" />
          <h2 className="text-lg font-bold text-red-800">Dangerous Zone (Destructive)</h2>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          <ActionCard 
            title="Reset Project Soft" 
            desc="Hapus semua data kecuali Produk dan Kategori (User/Order/Cart)."
            method="POST"
            isDanger
            onRun={() => executeAction('reset-soft')}
            isLoading={loading === 'reset-soft'}
          />
          <ActionCard 
            title="Delete Sync Seed" 
            desc="Memutus hubungan antara asset gambar dan data produk."
            method="DEL"
            isDanger
            onRun={() => executeAction('del-sync')}
            isLoading={loading === 'del-sync'}
          />
          <ActionCard 
            title="Reset Database Basic" 
            desc="Reset Database Tanpa nambah produk."
            method="POST"
            isDanger
            onRun={() => executeAction('reset-basic')}
            isLoading={loading === 'reset-basic'}
          />
          <ActionCard 
            title="Cloudinary Purge" 
            desc="Hapus SEMUA file di Cloudinary storage secara permanen."
            method="DEL"
            isDanger
            icon={<CloudOff className="w-4 h-4" />}
            onRun={() => executeAction('purge-cloud')}
            isLoading={loading === 'purge-cloud'}
          />
          <ActionCard 
            title="Reset Database Full" 
            desc="Wipe out semua tabel dan buat ulang. Data akan hilang total!"
            method="POST"
            isDanger
            icon={<Trash2 className="w-4 h-4" />}
            onRun={() => executeAction('reset-full')}
            isLoading={loading === 'reset-full'}
          />
        </div>
      </section>
    </div>
  )
}

interface CardProps {
  title: string
  desc: string
  method: "GET" | "POST" | "PUT" | "DEL"
  onRun: () => void
  isLoading?: boolean
  isDanger?: boolean
  icon?: React.ReactNode
}

function ActionCard({ title, desc, method, onRun, isLoading, isDanger, icon }: CardProps) {
  const methodColor: any = {
    POST: "bg-indigo-100 text-indigo-700 border-indigo-200",
    DEL: "bg-red-100 text-red-700 border-red-200",
  }

  const ActionButton = (
    <Button 
      variant={isDanger ? "destructive" : "default"} 
      className={`w-full group ${!isDanger && 'bg-white text-slate-900 border hover:bg-slate-50'}`}
      disabled={isLoading}
    >
      {isLoading ? (
        <Loader2 className="w-4 h-4 animate-spin mr-2" />
      ) : (
        icon || <RefreshCcw className={`w-4 h-4 mr-2 ${!isDanger && 'text-indigo-600'}`} />
      )}
      {isLoading ? "Running..." : "Execute"}
    </Button>
  )

  return (
    <Card className={`overflow-hidden transition-all hover:shadow-md ${isDanger ? 'border-red-100' : 'border-slate-200'}`}>
      <CardHeader className="p-4 space-y-2">
        <div className="flex justify-between items-center">
          <Badge variant="outline" className={`${methodColor[method]} text-[10px] font-bold`}>
            {method}
          </Badge>
          {isDanger && <Badge className="bg-red-600 text-[10px]">RISK</Badge>}
        </div>
        <CardTitle className="text-md font-bold leading-tight">{title}</CardTitle>
        <CardDescription className="text-xs line-clamp-2 min-h-[32px]">
          {desc}
        </CardDescription>
      </CardHeader>
      <CardContent className="p-4 pt-0">
        {isDanger ? (
          <AlertDialog>
            <AlertDialogTrigger asChild>
              {ActionButton}
            </AlertDialogTrigger>
            <AlertDialogContent>
              <AlertDialogHeader>
                <AlertDialogTitle>Konfirmasi Tindakan Berbahaya</AlertDialogTitle>
                <AlertDialogDescription>
                  Anda akan menjalankan <strong>{title}</strong>. 
                  Tindakan ini bersifat destruktif dan data tidak dapat dikembalikan. Lanjutkan?
                </AlertDialogDescription>
              </AlertDialogHeader>
              <AlertDialogFooter>
                <AlertDialogCancel>Batal</AlertDialogCancel>
                <AlertDialogAction onClick={onRun} className="bg-red-600 hover:bg-red-700">
                  Ya, Eksekusi
                </AlertDialogAction>
              </AlertDialogFooter>
            </AlertDialogContent>
          </AlertDialog>
        ) : (
          <div onClick={onRun}>{ActionButton}</div>
        )}
      </CardContent>
    </Card>
  )
}