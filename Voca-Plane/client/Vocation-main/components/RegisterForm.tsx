import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label" // Menggunakan Label standar shadcn

export function RegisterForm({ className, ...props }: React.ComponentProps<"form">) {
  return (
    <form className={cn("flex flex-col gap-5", className)} {...props}>
      <div className="flex flex-col items-center gap-2 text-center">
        <h1 className="text-2xl font-bold tracking-tight text-slate-900">Create your account</h1>
        <p className="text-sm text-slate-500">
          Join SkyPass to start booking your next flight.
        </p>
      </div>

      <div className="grid gap-4">
        {/* Full Name Field */}
        <div className="grid gap-2 text-left">
          <Label htmlFor="full-name" className="text-slate-700">Full Name (as per Passport)</Label>
          <Input id="full-name" placeholder="John Doe" required className="border-slate-200 focus:ring-indigo-500" />
        </div>

        {/* Email Field */}
        <div className="grid gap-2 text-left">
          <Label htmlFor="email" className="text-slate-700">Email Address</Label>
          <Input id="email" type="email" placeholder="name@example.com" required className="border-slate-200 focus:ring-indigo-500" />
        </div>

        {/* Password Field */}
        <div className="grid gap-2 text-left">
          <Label htmlFor="password" className="text-slate-700">Password</Label>
          <Input id="password" type="password" required className="border-slate-200 focus:ring-indigo-500" />
          <p className="text-[10px] text-slate-400">Must be at least 8 characters long.</p>
        </div>

        <Button className="w-full bg-indigo-600 hover:bg-indigo-700 text-white shadow-md shadow-indigo-100 mt-2">
          Create Account
        </Button>
      </div>

      <div className="relative py-2">
        <div className="absolute inset-0 flex items-center"><span className="w-full border-t border-slate-200"></span></div>
        <div className="relative flex justify-center text-xs uppercase"><span className="bg-white px-2 text-slate-400">Or sign up with</span></div>
      </div>

      <div className="grid grid-cols-2 gap-4">
        <Button variant="outline" className="border-slate-200">
          Google
        </Button>
        <Button variant="outline" className="border-slate-200">
          Apple ID
        </Button>
      </div>

      <p className="text-center text-sm text-slate-500">
        Already have an account?{" "}
        <a href="/" className="font-semibold text-indigo-600 hover:underline">
          Sign In
        </a>
      </p>
    </form>
  )
}