"use client"

import Link from "next/link"
import { Plane } from "lucide-react"
import { Button } from "@/components/ui/button"
import UserItem from "./UserItem"
import { ThemeToggle } from "./ThemeToggle"
import Cookies from "js-cookie"
import { useQuery } from "@tanstack/react-query"
import { getProfile } from "@/lib/api/UserApi"

export function Header() {
  const token = Cookies.get("access_token") || Cookies.get("token")

  const { data, isLoading, isError } = useQuery({
    queryKey: ["user-profile"],
    queryFn: getProfile,
    retry: false,
    enabled: Boolean(token),
  })

  const isAuthenticated = Boolean(token) && Boolean(data?.data) && !isError

  return (
    <header className="absolute top-4 left-1/2 -translate-x-1/2 z-50 w-[95%] max-w-7xl">
      <div className="bg-white/80 dark:bg-slate-900/80 backdrop-blur-md border border-white/20 dark:border-white/5 shadow-xl rounded-[2.5rem] px-8 h-20 flex items-center justify-between transition-all">
        <Link href="/" className="flex items-center gap-3 group">
          <div className="bg-primary p-2.5 rounded-2xl shadow-lg shadow-primary/20 group-hover:scale-110 transition-transform">
            <Plane className="size-6 rotate-45 text-white" />
          </div>
          <span className="font-black text-2xl tracking-tighter bg-clip-text text-transparent bg-gradient-to-r from-primary to-primary/60">
            VocaPlane
          </span>
        </Link>

        <nav className="hidden lg:flex items-center gap-10">
          <Link href="/" className="text-sm font-bold text-slate-600 hover:text-primary dark:text-slate-300 dark:hover:text-primary transition-colors relative group/link">
            Home
            <span className="absolute -bottom-1 left-0 w-0 h-0.5 bg-primary transition-all group-hover/link:w-full" />
          </Link>
          <Link href="/promos" className="text-sm font-bold text-slate-600 hover:text-primary dark:text-slate-300 dark:hover:text-primary transition-colors relative group/link">
            Promos
            <span className="absolute -bottom-1 left-0 w-0 h-0.5 bg-primary transition-all group-hover/link:w-full" />
          </Link>
          <Link href="/my-bookings" className="text-sm font-bold text-slate-600 hover:text-primary dark:text-slate-300 dark:hover:text-primary transition-colors relative group/link">
            My Bookings
            <span className="absolute -bottom-1 left-0 w-0 h-0.5 bg-primary transition-all group-hover/link:w-full" />
          </Link>
        </nav>

        <div className="flex items-center gap-6">
          <ThemeToggle />
          {isAuthenticated ? (
            <UserItem />
          ) : isLoading ? (
            <div className="h-10 w-24 bg-slate-200/50 animate-pulse rounded-full" />
          ) : (
            <>
              <Link href="/login">
                <Button variant="ghost" className="text-slate-600 hover:text-primary font-bold rounded-full px-6 transition-colors">
                  Log In
                </Button>
              </Link>
              <Link href="/register">
                <Button className="bg-primary hover:bg-primary/90 text-white shadow-xl shadow-primary/20 rounded-full px-8 font-black transition-all hover:-translate-y-1 active:scale-95">
                  Get Started
                </Button>
              </Link>
            </>
          )}
        </div>
      </div>
    </header>
  )
}