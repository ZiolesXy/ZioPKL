import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label" // Menggunakan Label standar shadcn

export function RegisterForm({ className, ...props }: React.ComponentProps<"form">) {
  return (
    <form className={cn("flex flex-col gap-10 bg-white/40 dark:bg-slate-900/40 backdrop-blur-3xl p-10 md:p-14 rounded-[3.5rem] border border-white/20 dark:border-white/5 shadow-[0_32px_80px_-16px_rgba(0,0,0,0.1)]", className)} {...props}>
      <div className="flex flex-col items-center gap-4 text-center">
        <h1 className="text-4xl font-black tracking-tighter text-slate-900 dark:text-white">Create premium account</h1>
        <p className="text-slate-500 font-medium text-lg leading-relaxed max-w-sm">
          Join VocaPlane to elevate your travel standards.
        </p>
      </div>

      <div className="grid gap-8">
        {/* Full Name Field */}
        <div className="grid gap-3 text-left">
          <Label htmlFor="full-name" className="text-[11px] font-black uppercase tracking-[0.2em] text-slate-400 ml-2">Passport Full Name</Label>
          <Input id="full-name" placeholder="ALEXANDER SMITH" required className="h-16 px-8 bg-slate-100/50 dark:bg-white/5 border-none rounded-[1.5rem] focus:ring-4 focus:ring-primary/20 text-slate-900 dark:text-white font-bold text-lg placeholder:text-slate-300 dark:placeholder:text-slate-700 transition-all outline-none" />
        </div>

        {/* Email Field */}
        <div className="grid gap-3 text-left">
          <Label htmlFor="email" className="text-[11px] font-black uppercase tracking-[0.2em] text-slate-400 ml-2">Verified Email</Label>
          <Input id="email" type="email" placeholder="alex@premium-travel.com" required className="h-16 px-8 bg-slate-100/50 dark:bg-white/5 border-none rounded-[1.5rem] focus:ring-4 focus:ring-primary/20 text-slate-900 dark:text-white font-bold text-lg placeholder:text-slate-300 dark:placeholder:text-slate-700 transition-all outline-none" />
        </div>

        {/* Password Field */}
        <div className="grid gap-3 text-left">
          <Label htmlFor="password" className="text-[11px] font-black uppercase tracking-[0.2em] text-slate-400 ml-2">Security Key</Label>
          <Input id="password" type="password" required className="h-16 px-8 bg-slate-100/50 dark:bg-white/5 border-none rounded-[1.5rem] focus:ring-4 focus:ring-primary/20 text-slate-900 dark:text-white font-bold text-lg transition-all outline-none" placeholder="••••••••••••" />
          <p className="text-[10px] font-bold text-slate-400 tracking-widest ml-2 uppercase">Min. 8 characters</p>
        </div>

        <Button className="w-full h-16 bg-primary hover:bg-primary/90 text-white font-black text-lg rounded-[1.5rem] shadow-2xl shadow-primary/30 transition-all hover:scale-[1.02] active:scale-95 mt-4">
          Initiate Account
        </Button>
      </div>

      <div className="relative py-2">
        <div className="absolute inset-0 flex items-center"><span className="w-full border-t border-slate-100 dark:border-white/5"></span></div>
        <div className="relative flex justify-center text-[10px] uppercase font-black tracking-[0.3em]"><span className="bg-transparent px-6 text-slate-400">Elite authentication</span></div>
      </div>

      <div className="grid grid-cols-2 gap-4">
        <Button variant="outline" className="h-14 border-2 border-slate-100 dark:border-white/5 text-slate-600 dark:text-slate-400 font-bold rounded-[1.25rem] hover:bg-slate-50 dark:hover:bg-white/5 transition-all">
          Google
        </Button>
        <Button variant="outline" className="h-14 border-2 border-slate-100 dark:border-white/5 text-slate-600 dark:text-slate-400 font-bold rounded-[1.25rem] hover:bg-slate-50 dark:hover:bg-white/5 transition-all">
          Apple ID
        </Button>
      </div>

      <p className="text-center text-sm font-medium text-slate-500">
        Already an elite member?{" "}
        <a href="/login" className="font-black text-primary hover:tracking-[0.05em] transition-all">
          Sign In
        </a>
      </p>
    </form>
  )
}