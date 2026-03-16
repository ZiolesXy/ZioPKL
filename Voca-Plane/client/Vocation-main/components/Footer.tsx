import Link from "next/link"
import { Facebook, Instagram, Twitter, Plane } from "lucide-react"

export function Footer() {
  return (
    <footer className="relative bg-slate-900 text-white pt-24 pb-12 overflow-hidden">
      {/* Visual Accent */}
      <div className="absolute top-0 left-0 w-full h-px bg-gradient-to-r from-transparent via-primary/50 to-transparent" />
      
      <div className="container mx-auto px-8 relative z-10">
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-12 mb-20">
          <div className="space-y-6">
            <Link href="/" className="flex items-center gap-3">
              <div className="bg-primary p-2.5 rounded-2xl shadow-lg shadow-primary/20">
                <Plane className="size-6 rotate-45 text-white" />
              </div>
              <span className="font-black text-2xl tracking-tighter">
                VocaPlane
              </span>
            </Link>
            <p className="text-slate-400 font-medium leading-relaxed max-w-xs">
              Elevating your travel experience with world-class service and seamless booking.
            </p>
            <div className="flex gap-4">
              <Instagram className="size-5 cursor-pointer text-slate-500 hover:text-primary transition-colors" />
              <Twitter className="size-5 cursor-pointer text-slate-500 hover:text-primary transition-colors" />
              <Facebook className="size-5 cursor-pointer text-slate-500 hover:text-primary transition-colors" />
            </div>
          </div>

          <div>
            <h4 className="text-xs font-black uppercase tracking-[0.2em] mb-8 text-primary/80">Services</h4>
            <ul className="space-y-4">
              <li><Link href="#" className="text-slate-400 hover:text-white transition-colors font-semibold text-sm">Flight Schedule</Link></li>
              <li><Link href="#" className="text-slate-400 hover:text-white transition-colors font-semibold text-sm">Airport Info</Link></li>
              <li><Link href="#" className="text-slate-400 hover:text-white transition-colors font-semibold text-sm">Refund Policy</Link></li>
            </ul>
          </div>

          <div>
            <h4 className="text-xs font-black uppercase tracking-[0.2em] mb-8 text-primary/80">Support</h4>
            <ul className="space-y-4">
              <li><Link href="#" className="text-slate-400 hover:text-white transition-colors font-semibold text-sm">Help Center</Link></li>
              <li><Link href="#" className="text-slate-400 hover:text-white transition-colors font-semibold text-sm">Safe Travel</Link></li>
              <li><Link href="#" className="text-slate-400 hover:text-white transition-colors font-semibold text-sm">Contact Us</Link></li>
            </ul>
          </div>

          <div>
            <h4 className="text-xs font-black uppercase tracking-[0.2em] mb-8 text-primary/80">Newsletter</h4>
            <p className="text-slate-400 font-medium mb-6 text-sm leading-relaxed">Stay updated with our latest offers and premium routes.</p>
            <div className="flex gap-2">
              <input type="email" placeholder="Email Address" className="flex-1 bg-white/5 border border-white/10 rounded-xl px-4 py-2 text-xs font-medium focus:outline-none focus:border-primary transition-colors placeholder:text-slate-600" />
              <button className="bg-primary px-4 py-2 rounded-xl text-xs font-black hover:bg-primary/90 transition-colors shadow-lg shadow-primary/20">Join</button>
            </div>
          </div>
        </div>

        <div className="pt-12 border-t border-white/5 flex flex-col md:flex-row justify-between items-center gap-6">
          <p className="text-slate-500 text-[10px] font-black uppercase tracking-widest">
            © 2026 VocaPlane. All Rights Reserved.
          </p>
          <div className="flex gap-8 text-slate-500 text-[10px] font-black uppercase tracking-widest">
            <Link href="/" className="hover:text-white transition-colors">Privacy Policy</Link>
            <Link href="/" className="hover:text-white transition-colors">Terms of Service</Link>
          </div>
        </div>
      </div>
    </footer>
  )
}