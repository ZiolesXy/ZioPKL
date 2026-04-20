"use client"
import Link from "next/link"
import { Separator } from "@/components/ui/separator"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Facebook, Instagram, Twitter, Youtube } from "lucide-react"
import { usePathname } from "next/navigation"

export default function Footer() {
  const pathname = usePathname()
  if (pathname === '/chat') return null;

  return (
    <footer className="w-full bg-white border-t">
      <div className="container mx-auto px-4 py-12">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-8">

          <div className="space-y-4">
            <h2 className="text-2xl font-bold text-teal-600">Voca</h2>
            <p className="text-sm text-muted-foreground">
              E-commerce nomor 1 untuk kebutuhan Sehari-hari di Indonesia.
            </p>
          </div>

          <div>
            <h3 className="font-semibold mb-4">Layanan</h3>
            <ul className="space-y-2 text-sm text-muted-foreground">
              <li><Link href="/product" className="hover:text-teal-600">Semua Produk</Link></li>
              <li><Link href="/category" className="hover:text-teal-600">Kategori</Link></li>
              <li><Link href="/" className="hover:text-teal-600">Rekomendasi</Link></li>
            </ul>
          </div>

          <div>
            <h3 className="font-semibold mb-4">Bantuan</h3>
            <ul className="space-y-2 text-sm text-muted-foreground">
              <li><Link href="/faq" className="hover:text-teal-600">FAQ</Link></li>
              <li><Link href="/shipping" className="hover:text-teal-600">Info Pengiriman</Link></li>
              <li><Link href="/contact" className="hover:text-teal-600">Hubungi Kami</Link></li>
            </ul>
          </div>

          <div className="space-y-4">
            <h3 className="font-semibold">Langganan Newsletter</h3>
            <div className="flex gap-2">
              <Input placeholder="Email anda..." className="max-w-50" />
              <Button className="bg-teal-600 hover:bg-teal-700">Ikuti</Button>
            </div>
          </div>
        </div>

        <Separator className="my-8" />

        <div className="flex flex-col md:flex-row justify-between items-center gap-4">
          <p className="text-sm text-muted-foreground">
            © 2026 Voca Store All rights reserved.
            <Link href="/privacy">| Privacy </Link>
            <Link href="/terms"> | Terms </Link>
            <Link href="/refund"> | Refund</Link>


          </p>
          <div className="flex gap-4">
            <Link target="_blank" href="https://x.com/search?q=dreamybull"><Twitter className="w-5 h-5 text-muted-foreground hover:text-teal-600" /></Link>
            <Link target="_blank" href="https://www.instagram.com/popular/msbreewc-official/"><Instagram className="w-5 h-5 text-muted-foreground hover:text-teal-600" /></Link>
            <Link target="_blank" href="https://www.youtube.com/@AmrulsiGuruDigital"><Youtube className="w-5 h-5 text-muted-foreground hover:text-teal-600" /></Link>
          </div>
        </div>
      </div>
    </footer>
  )
}