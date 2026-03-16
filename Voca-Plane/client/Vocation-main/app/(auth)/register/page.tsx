import { Plane, Sparkles } from "lucide-react"
import { RegisterForm } from "@/components/RegisterForm"

export default function RegisterPage() {
  return (
    <div className="grid min-h-svh lg:grid-cols-2 bg-white">
      {/* Bagian Visual (Kiri) */}
      <div className="relative hidden lg:block overflow-hidden bg-indigo-900">
        <img
          src="https://images.unsplash.com/photo-1542296332-2e4473faf563?q=80&w=2070" 
          alt="Airplane Cabin"
          className="absolute inset-0 h-full w-full object-cover opacity-60 mix-blend-overlay"
        />
        
        {/* Overlay SVG & Content */}
        <div className="absolute inset-0 flex flex-col justify-between p-12 text-white">
          <div className="flex items-center gap-2 font-bold text-xl">
            <Plane className="size-6 rotate-45" />
            SkyPass.
          </div>
          
          <div className="max-w-md">
            <div className="inline-flex items-center gap-2 bg-white/20 backdrop-blur-md px-3 py-1 rounded-full text-xs font-medium mb-4">
              <Sparkles className="size-3 text-yellow-300" />
              New Member Privilege
            </div>
            <h2 className="text-4xl font-bold leading-tight">Travel the world with more rewards.</h2>
            <p className="mt-4 text-indigo-100">Get 10% off on your first flight booking and exclusive access to airport lounges.</p>
          </div>
          
          <div className="text-sm text-indigo-200/60">
            © 2026 SkyPass International. All flights reserved.
          </div>
        </div>
      </div>

      {/* Bagian Form (Kanan) */}
      <div className="flex flex-col p-6 md:p-10 justify-center items-center bg-slate-50/50">
        <div className="w-full max-w-sm">
          <RegisterForm />
        </div>
      </div>
    </div>
  )
}