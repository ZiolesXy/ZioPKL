"use client"
import { 
  Plane, 
  Settings, 
  LayoutDashboard,
  LogOut,
  Building2,
  CalendarDays,
  MapPin,
  PlaneTakeoff,
  ReceiptText,
  Terminal,
  TicketPercent,
  UserCog
} from "lucide-react"
import { usePathname } from "next/navigation"
import Link from "next/link"
import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"

import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar"
import { Avatar, AvatarFallback, AvatarImage } from "../ui/avatar"

// 1. Kelompokkan menu agar lebih terorganisir
const menuGroups = [
  {
    label: "Main",
    items: [
      { title: "Dashboard", url: "/dashboard", icon: LayoutDashboard },
      { title: "Flight Schedule", url: "/flights-schedule", icon: CalendarDays },
    ],
  },
  {
    label: "Flight Operations",
    items: [
      { title: "Flights Master", url: "/flights/manage", icon: PlaneTakeoff }, // Kelola flight_classes & seats [cite: 32, 52]
      { title: "Airlines", url: "/airlines", icon: Building2 }, // [cite: 9]
      { title: "Airports", url: "/airports", icon: MapPin }, // [cite: 11]
    ],
  },
  {
    label: "Sales & Marketing",
    items: [
      { title: "Transactions", url: "/transactions", icon: ReceiptText }, // [cite: 77]
      // { title: "Promo Codes", url: "/promos", icon: TicketPercent }, // [cite: 82]
    ],
  },
  {
    label: "User Management",
    items: [
      { title: "User Access", url: "/users-monitoring", icon: UserCog }, // [cite: 4]
      // { title: "Admin Tools", url: "/admin-tool", icon: Terminal }, // Untuk POST Seed & GET Test dari gambar
    ],
  },
  {
    label: "System",
    items: [
      { title: "Settings", url: "/settings", icon: Settings },
    ],
  },
]

export function AppSidebar() {
  const pathname = usePathname() // Untuk mendeteksi halaman aktif

  return (
    <Sidebar collapsible="icon" className="border-r border-slate-200/50 dark:border-white/5 bg-white/80 dark:bg-slate-950/80 backdrop-blur-md">
      <SidebarHeader className="py-8 border-b border-slate-200/50 dark:border-white/5" >
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton size="lg" asChild className="hover:bg-transparent">
              <Link href="/dashboard" className="flex items-center gap-4">
                <div className="flex aspect-square size-12 items-center justify-center rounded-[1.25rem] bg-gradient-to-br from-primary to-primary/80 shadow-2xl shadow-primary/40 text-white transition-transform hover:scale-105 active:scale-95">
                  <Plane className="size-6 rotate-45" />
                </div>
                <div className="flex flex-col gap-0.5 leading-none group-data-[collapsible=icon]:hidden">
                  <span className="font-black text-2xl tracking-tighter text-slate-900 dark:text-white">Voca<span className="text-primary italic">Admin</span></span>
                  <span className="text-[10px] text-slate-400 dark:text-slate-500 font-black uppercase tracking-[0.25em]">Cloud Systems</span>
                </div>
              </Link>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>

      <SidebarContent className="px-4 pt-6">
        {menuGroups.map((group) => (
          <SidebarGroup key={group.label} className="py-5">
            <SidebarGroupLabel className="text-slate-400 dark:text-slate-500 font-black px-4 py-2 uppercase text-[10px] tracking-[0.3em] mb-3 opacity-40 group-data-[collapsible=icon]:hidden">
              {group.label}
            </SidebarGroupLabel>
            <SidebarGroupContent>
              <SidebarMenu className="gap-2">
                {group.items.map((item) => {
                  const isActive = pathname === item.url
                  
                  return (
                    <SidebarMenuItem key={item.title}>
                      <SidebarMenuButton 
                        asChild 
                        tooltip={item.title}
                        className={cn(
                          "transition-all duration-200 py-7 rounded-[1.5rem] group/btn",
                          isActive 
                            ? "bg-primary text-white hover:bg-primary shadow-lg shadow-primary/20 font-bold" 
                            : "text-slate-500 dark:text-slate-400 hover:bg-slate-100/50 dark:hover:bg-white/5 hover:text-slate-900 dark:hover:text-white"
                        )}
                      >
                        <Link href={item.url} className="flex items-center gap-4 w-full px-5">
                          <item.icon className={cn("size-5 transition-all duration-300", isActive ? "text-white scale-110" : "text-slate-400 group-hover/btn:text-primary group-hover/btn:scale-110")} />
                          <span className="group-data-[collapsible=icon]:hidden text-sm uppercase tracking-wider">{item.title}</span>
                          {isActive && (
                            <div className="ml-auto size-2 rounded-full bg-white animate-pulse group-data-[collapsible=icon]:hidden" />
                          )}
                        </Link>
                      </SidebarMenuButton>
                    </SidebarMenuItem>
                  )
                })}
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>
        ))}
      </SidebarContent>

      <SidebarFooter className="p-6 border-t border-slate-200/50 dark:border-white/5 bg-slate-50/50 dark:bg-white/2">
        <div className="flex items-center gap-4 px-2 group-data-[collapsible=icon]:justify-center">
          <div className="relative group">
            <div className="absolute -inset-1 bg-gradient-to-r from-primary to-indigo-600 rounded-full blur opacity-25 group-hover:opacity-50 transition duration-1000 group-hover:duration-200"></div>
            <Avatar className="h-11 w-11 border-2 border-white dark:border-slate-800 shadow-xl relative">
              <AvatarImage src="https://github.com/shadcn.png" />
              <AvatarFallback className="bg-primary text-white font-black text-xs">AD</AvatarFallback>
            </Avatar>
          </div>
          <div className="flex flex-col gap-0.5 leading-none group-data-[collapsible=icon]:hidden overflow-hidden">
            <p className="font-black text-sm text-slate-900 dark:text-white tracking-tight truncate">Administrator</p>
            <p className="text-[10px] text-slate-400 dark:text-slate-500 font-bold uppercase tracking-widest truncate">System Root</p>
          </div>
          <Button variant="ghost" size="icon" className="ml-auto rounded-2xl hover:bg-red-50 dark:hover:bg-red-500/10 hover:text-red-500 transition-all group-data-[collapsible=icon]:hidden">
            <LogOut className="size-5" />
          </Button>
        </div>
      </SidebarFooter>
    </Sidebar>
  )
}