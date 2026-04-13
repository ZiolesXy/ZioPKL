import { PlaneTakeoff } from "lucide-react"
import Image from "next/image"
import Link from "next/link"
import { LoginForm } from "@/components/LoginForm"

export default function LoginPage() {
  return (
    <div className="min-h-svh flex items-center justify-center bg-slate-50 relative overflow-hidden">
      {/* Background Decorative Elements */}
      <div className="absolute top-0 left-0 w-full h-full overflow-hidden pointer-events-none">
        <div className="absolute top-[-10%] left-[-10%] w-[40%] h-[40%] bg-primary/5 rounded-full blur-[120px]" />
        <div className="absolute bottom-[-10%] right-[-10%] w-[40%] h-[40%] bg-primary/5 rounded-full blur-[120px]" />
      </div>

      <div className="relative z-10 w-full max-w-md p-6 md:p-10">
        <div className="flex flex-col gap-8">
          <div className="flex justify-center flex-col items-center gap-4">
            <Link href="/" className="flex items-center gap-3 group">
              <div className="flex size-12 items-center justify-center rounded-2xl bg-primary text-white shadow-xl shadow-primary/20 transition-transform group-hover:scale-110">
                <PlaneTakeoff className="size-6 rotate-0 transition-transform group-hover:rotate-12" />
              </div>
              <span className="font-black text-3xl tracking-tighter text-slate-900">Voca<span className="text-primary">Plane</span></span>
            </Link>
          </div>
          
          <div className="animate-in fade-in slide-in-from-bottom-8 duration-700">
            <LoginForm />
          </div>

          <p className="text-center text-slate-400 text-xs font-medium px-8 leading-relaxed">
            By signing in, you agree to our <a href="#" className="underline underline-offset-4 hover:text-primary transition-colors">Terms of Service</a> and <a href="#" className="underline underline-offset-4 hover:text-primary transition-colors">Privacy Policy</a>.
          </p>
        </div>
      </div>
    </div>
  )
}