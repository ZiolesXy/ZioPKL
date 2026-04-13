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
    <form onSubmit={handleLogin} className={cn("flex flex-col gap-10 bg-white/40 dark:bg-slate-900/40 backdrop-blur-3xl p-10 md:p-14 rounded-[3.5rem] border border-white/20 dark:border-white/5 shadow-[0_32px_80px_-16px_rgba(0,0,0,0.1)]", className)} {...props}>
      <div className="flex flex-col items-center gap-4 text-center">
        <div className="bg-primary/20 p-5 rounded-[2rem] mb-2 shadow-inner">
          <Plane className="size-10 rotate-45 text-primary" />
        </div>
        <h1 className="text-4xl font-black tracking-tighter text-slate-900 dark:text-white">Welcome back</h1>
        <p className="text-slate-500 font-medium text-lg leading-relaxed max-w-sm">Securely log in to your premium travel account and manage your elite journeys.</p>
        {error && (
          <div className="animate-in fade-in slide-in-from-top-4 duration-300 w-full mt-4">
            <p className="text-xs font-black uppercase tracking-widest text-red-500 bg-red-50/50 border border-red-100 p-4 rounded-2xl w-full text-center">
              {error}
            </p>
          </div>
        )}
      </div>

      <FieldGroup className="gap-8">
        <Field className="space-y-3">
          <FieldLabel className="text-[11px] font-black uppercase tracking-[0.2em] text-slate-400 ml-2">Verified Email</FieldLabel>
          <Input
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            className="h-16 px-8 bg-slate-100/50 dark:bg-white/5 border-none rounded-[1.5rem] focus:ring-4 focus:ring-primary/20 text-slate-900 dark:text-white font-bold text-lg placeholder:text-slate-300 dark:placeholder:text-slate-700 transition-all outline-none"
            id="email" type="email" placeholder="alex@premium-travel.com"
            required
          />
        </Field>

        <Field className="space-y-3">
          <div className="flex items-center justify-between ml-2">
            <FieldLabel className="text-[11px] font-black uppercase tracking-[0.2em] text-slate-400">Security Key</FieldLabel>
            <a href="#" className="text-xs font-black uppercase tracking-widest text-primary hover:tracking-[0.2em] transition-all">
              Lost access?
            </a>
          </div>
          <Input 
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            className="h-16 px-8 bg-slate-100/50 dark:bg-white/5 border-none rounded-[1.5rem] focus:ring-4 focus:ring-primary/20 text-slate-900 dark:text-white font-bold text-lg placeholder:text-slate-300 dark:placeholder:text-slate-700 transition-all outline-none" 
            id="password" 
            type="password" 
            placeholder="••••••••••••"
            required
          />
        </Field>

        <Button 
          type="submit" 
          disabled={loading}
          className="h-16 bg-primary hover:bg-primary/90 text-white font-black text-lg rounded-[1.5rem] shadow-2xl shadow-primary/30 transition-all hover:scale-[1.02] active:scale-95 disabled:scale-100 disabled:opacity-50"
        >
          {loading ? (
            <div className="flex items-center gap-3">
              <span className="w-2.5 h-2.5 rounded-full bg-white animate-bounce" />
              <span className="w-2.5 h-2.5 rounded-full bg-white animate-bounce [animation-delay:0.2s]" />
              <span className="w-2.5 h-2.5 rounded-full bg-white animate-bounce [animation-delay:0.4s]" />
            </div>
          ) : "Enter Dashboard"}
        </Button>

        <div className="relative my-2">
          <div className="absolute inset-0 flex items-center">
            <span className="w-full border-t border-slate-100 dark:border-white/5" />
          </div>
          <div className="relative flex justify-center text-[10px] uppercase font-black tracking-[0.3em]">
            <span className="bg-transparent px-6 text-slate-400">Identity check</span>
          </div>
        </div>

        <Link href="/register" className="w-full">
          <Button variant="ghost" className="w-full h-16 border-2 border-slate-100 dark:border-white/5 text-slate-600 dark:text-slate-400 font-black text-lg rounded-[1.5rem] hover:bg-slate-50 dark:hover:bg-white/5 transition-all">
            Initiate Account
          </Button>
        </Link>
      </FieldGroup>
    </form>
  );
}