"use client"

import Link from "next/link"
import { Plane } from "lucide-react"
import { Button } from "@/components/ui/button"
import UserItem from "./UserItem"
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
    <header className="sticky top-0 z-50 w-full border-b border-white/20 bg-white/70 backdrop-blur-xl dark:bg-slate-950/70">
      <div className="container mx-auto flex h-16 items-center justify-between px-6">
        <Link href="/" className="flex items-center gap-2 group">
          <div className="bg-primary/10 p-2 rounded-xl group-hover:bg-primary/20 transition-colors">
            <Plane className="size-6 rotate-45 text-primary" />
          </div>
          <span className="font-bold text-2xl tracking-tight bg-clip-text text-transparent bg-gradient-to-r from-primary to-primary/60">
            Vocation
          </span>
        </Link>

        <nav className="hidden md:flex items-center gap-8 text-sm font-semibold">
          <Link href="/" className="text-slate-600 hover:text-primary transition-all relative group/link">
            Home
            <span className="absolute -bottom-1 left-0 w-0 h-0.5 bg-primary transition-all group-hover/link:w-full" />
          </Link>
          <Link href="/promos" className="text-slate-600 hover:text-primary transition-all relative group/link">
            Promos
            <span className="absolute -bottom-1 left-0 w-0 h-0.5 bg-primary transition-all group-hover/link:w-full" />
          </Link>
          <Link href="/my-bookings" className="text-slate-600 hover:text-primary transition-all relative group/link">
            My Bookings
            <span className="absolute -bottom-1 left-0 w-0 h-0.5 bg-primary transition-all group-hover/link:w-full" />
          </Link>
        </nav>

        <div className="flex items-center gap-4">
          {isAuthenticated ? (
            <UserItem />
          ) : isLoading ? (
            <div className="h-10 w-24 bg-slate-100 animate-pulse rounded-full" />
          ) : (
            <>
              <Link href="/login">
                <Button variant="ghost" className="text-slate-600 hover:text-primary font-semibold rounded-full px-6">
                  Log In
                </Button>
              </Link>
              <Link href="/register">
                <Button className="bg-primary hover:bg-primary/90 text-white shadow-lg shadow-primary/20 rounded-full px-8 transition-all hover:-translate-y-0.5">
                  Sign Up
                </Button>
              </Link>
            </>
          )}
        </div>
      </div>
    </header>
  )
}