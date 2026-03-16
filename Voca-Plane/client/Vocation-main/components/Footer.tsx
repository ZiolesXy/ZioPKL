import Link from "next/link"
import { Facebook, Instagram, Twitter, Plane } from "lucide-react"

export function Footer() {
  return (
    <footer className="w-full border-t bg-slate-50">
      <div className="container mx-auto px-4 py-12">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-8">
          {/* Brand */}
          <div className="col-span-1 md:col-span-1">
            <Link href="/" className="flex items-center gap-2 font-bold text-xl text-indigo-600 mb-4">
              <Plane className="size-6 rotate-45" />
              Vocation
            </Link>
            <p className="text-sm text-slate-500 leading-relaxed">
              Making your journey easier, faster, and more affordable since 2026.
            </p>
          </div>

          {/* Quick Links */}
          <div>
            <h4 className="font-semibold text-slate-900 mb-4 text-sm uppercase tracking-wider">Services</h4>
            <ul className="space-y-2 text-sm text-slate-600">
              <li><Link href="#" className="hover:text-indigo-600">Flight Schedule</Link></li>
              <li><Link href="#" className="hover:text-indigo-600">Airport Info</Link></li>
              <li><Link href="#" className="hover:text-indigo-600">Refund Policy</Link></li>
            </ul>
          </div>

          {/* Company */}
          <div>
            <h4 className="font-semibold text-slate-900 mb-4 text-sm uppercase tracking-wider">Company</h4>
            <ul className="space-y-2 text-sm text-slate-600">
              <li><Link href="#" className="hover:text-indigo-600">About Us</Link></li>
              <li><Link href="#" className="hover:text-indigo-600">Partnership</Link></li>
              <li><Link href="#" className="hover:text-indigo-600">Contact</Link></li>
            </ul>
          </div>

          {/* Socials */}
          <div>
            <h4 className="font-semibold text-slate-900 mb-4 text-sm uppercase tracking-wider">Follow Us</h4>
            <div className="flex gap-4 text-slate-500">
              <Instagram className="size-5 cursor-pointer hover:text-indigo-600" />
              <Twitter className="size-5 cursor-pointer hover:text-indigo-600" />
              <Facebook className="size-5 cursor-pointer hover:text-indigo-600" />
            </div>
          </div>
        </div>

        <div className="mt-12 pt-8 border-t border-slate-200 text-center text-sm text-slate-400">
          <p>© 2026 SkyPass Inc. All rights reserved.</p>
        </div>
      </div>
    </footer>
  )
}