"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { Plane } from "lucide-react";
import Cookies from "js-cookie"; // Install: npm install js-cookie @types/js-cookie
import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import { Field, FieldGroup, FieldLabel } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Login } from "@/lib/api/auth/login"; 
import { getProfile } from "@/lib/api/UserApi";

export function LoginForm({ className, ...props }: React.ComponentProps<"form">) {
  const router = useRouter();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError("");

    try {
      const res = await Login(email, password);

      if (res.success) {
        // 1. Simpan Access Token ke Cookie
        Cookies.set("token", res.data.access_token, { expires: 1 });
        Cookies.set("access_token", res.data.access_token, { expires: 1 });

        const roleFromLogin = res?.data?.role;
        if (roleFromLogin) {
          Cookies.set("role", roleFromLogin, { expires: 1 });
        } else {
          try {
            const profile = await getProfile();
            const roleFromProfile = profile?.data?.role;
            if (roleFromProfile) {
              Cookies.set("role", roleFromProfile, { expires: 1 });
            }
          } catch {
          }
        }

        // 3. Redirect ke root (user group) sesuai flow kamu
        router.push("/");
        router.refresh(); // Refresh agar middleware mendeteksi cookie baru
      }
    } catch (err: any) {
      setError(err.toString());
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={handleLogin} className={cn("flex flex-col gap-8 bg-white/50 backdrop-blur-xl p-8 rounded-[2.5rem] border border-white/20 shadow-2xl", className)} {...props}>
      <div className="flex flex-col items-center gap-3 text-center mb-4">
        <div className="bg-primary/10 p-4 rounded-3xl mb-2">
          <Plane className="size-8 rotate-45 text-primary" />
        </div>
        <h1 className="text-3xl font-black tracking-tight text-slate-900">Welcome Back</h1>
        <p className="text-slate-500 font-medium leading-relaxed">Enter your details to access your booking and explore more destinations.</p>
        {error && (
          <div className="animate-in fade-in slide-in-from-top-4 duration-300 w-full">
            <p className="text-xs font-bold text-red-500 bg-red-50/50 border border-red-100 p-3 rounded-2xl w-full text-center">
              {error}
            </p>
          </div>
        )}
      </div>

      <FieldGroup className="gap-6">
        <Field className="space-y-2">
          <FieldLabel className="text-xs font-black uppercase tracking-widest text-slate-400 ml-1">Email Address</FieldLabel>
          <Input
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            className="h-14 px-6 bg-slate-50 border-none rounded-2xl focus:ring-4 focus:ring-primary/10 text-slate-900 font-semibold placeholder:text-slate-300 transition-all"
            id="email" type="email" placeholder="user@flightbooking.com"
            required
          />
        </Field>

        <Field className="space-y-2">
          <div className="flex items-center justify-between ml-1">
            <FieldLabel className="text-xs font-black uppercase tracking-widest text-slate-400">Password</FieldLabel>
            <a href="#" className="text-xs font-bold text-primary hover:underline underline-offset-4 decoration-2">
              Forgot password?
            </a>
          </div>
          <Input 
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            className="h-14 px-6 bg-slate-50 border-none rounded-2xl focus:ring-4 focus:ring-primary/10 text-slate-900 font-semibold placeholder:text-slate-300 transition-all" 
            id="password" 
            type="password" 
            required
          />
        </Field>

        <Button 
          type="submit" 
          disabled={loading}
          className="h-14 bg-primary hover:bg-primary/90 text-white font-black rounded-2xl shadow-xl shadow-primary/20 transition-all hover:-translate-y-1 active:scale-95 disabled:scale-100 disabled:opacity-50"
        >
          {loading ? (
            <div className="flex items-center gap-2">
              <span className="w-2 h-2 rounded-full bg-white animate-bounce" />
              <span className="w-2 h-2 rounded-full bg-white animate-bounce [animation-delay:0.2s]" />
              <span className="w-2 h-2 rounded-full bg-white animate-bounce [animation-delay:0.4s]" />
            </div>
          ) : "Sign In to Account"}
        </Button>

        <div className="relative my-4">
          <div className="absolute inset-0 flex items-center">
            <span className="w-full border-t border-slate-100" />
          </div>
          <div className="relative flex justify-center text-[10px] uppercase font-black tracking-widest">
            <span className="bg-white/0 px-4 text-slate-300">New around here?</span>
          </div>
        </div>

        <Link href="/register" className="w-full">
          <Button variant="ghost" className="w-full h-14 border-2 border-slate-50 text-slate-600 font-black rounded-2xl hover:bg-slate-50 transition-all">
            Create an Account
          </Button>
        </Link>
      </FieldGroup>
    </form>
  );
}