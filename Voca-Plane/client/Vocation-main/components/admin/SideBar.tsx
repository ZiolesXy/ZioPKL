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
      { title: "Promo Codes", url: "/promos", icon: TicketPercent }, // [cite: 82]
    ],
  },
  {
    label: "User Management",
    items: [
      { title: "User Access", url: "/users-monitoring", icon: UserCog }, // [cite: 4]
      { title: "Admin Tools", url: "/admin-tool", icon: Terminal }, // Untuk POST Seed & GET Test dari gambar
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
    <Sidebar collapsible="icon" className="border-r border-slate-200/50 bg-slate-50/30 backdrop-blur-xl">
      <SidebarHeader className="py-6 border-b border-slate-200/50" >
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton size="lg" asChild className="hover:bg-transparent">
              <Link href="/dashboard" className="flex items-center gap-3">
                <div className="flex aspect-square size-10 items-center justify-center rounded-2xl bg-primary shadow-lg shadow-primary/20 text-white">
                  <Plane className="size-6 rotate-45" />
                </div>
                <div className="flex flex-col gap-0.5 leading-none group-data-[collapsible=icon]:hidden animate-in fade-in slide-in-from-left-4">
                  <span className="font-black text-xl tracking-tighter text-slate-900">Voca<span className="text-primary">Admin</span></span>
                  <span className="text-[10px] text-slate-400 font-bold uppercase tracking-widest">Enterprise Suite</span>
                </div>
              </Link>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>

      <SidebarContent className="px-2 pt-4">
        {menuGroups.map((group) => (
          <SidebarGroup key={group.label} className="py-4">
            <SidebarGroupLabel className="text-slate-400 font-black px-4 py-2 uppercase text-[10px] tracking-[0.2em] mb-2 opacity-50 group-data-[collapsible=icon]:hidden">
              {group.label}
            </SidebarGroupLabel>
            <SidebarGroupContent>
              <SidebarMenu className="gap-1">
                {group.items.map((item) => {
                  const isActive = pathname === item.url
                  
                  return (
                    <SidebarMenuItem key={item.title}>
                      <SidebarMenuButton 
                        asChild 
                        tooltip={item.title}
                        className={cn(
                          "transition-all duration-300 py-6 rounded-2xl group/btn",
                          isActive 
                            ? "bg-primary/10 text-primary hover:bg-primary/15 font-bold shadow-sm" 
                            : "text-slate-500 hover:bg-slate-100/50 hover:text-slate-900"
                        )}
                      >
                        <Link href={item.url} className="flex items-center gap-3 w-full px-4">
                          <item.icon className={cn("size-5 transition-transform group-hover/btn:scale-110", isActive ? "text-primary" : "text-slate-400 group-hover/btn:text-slate-600")} />
                          <span className="group-data-[collapsible=icon]:hidden">{item.title}</span>
                          {isActive && (
                            <div className="ml-auto size-1.5 rounded-full bg-primary shadow-[0_0_8px_rgba(var(--primary),0.5)] group-data-[collapsible=icon]:hidden" />
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

      <SidebarFooter className="p-4 border-t border-slate-200/50 bg-slate-100/30">
        <div className="flex items-center gap-3 px-2 group-data-[collapsible=icon]:justify-center">
          <Avatar className="h-10 w-10 border-2 border-white shadow-md">
            <AvatarImage src="https://github.com/shadcn.png" />
            <AvatarFallback className="bg-primary text-white font-black">AD</AvatarFallback>
          </Avatar>
          <div className="flex flex-col gap-0.5 leading-none group-data-[collapsible=icon]:hidden overflow-hidden">
            <p className="font-bold text-sm text-slate-900 truncate">Administrator</p>
            <p className="text-[10px] text-slate-400 font-medium truncate">admin@voca-plane.com</p>
          </div>
          <Button variant="ghost" size="icon" className="ml-auto rounded-xl hover:bg-red-50 hover:text-red-500 transition-colors group-data-[collapsible=icon]:hidden">
            <LogOut className="size-5" />
          </Button>
        </div>
      </SidebarFooter>
    </Sidebar>
  )
}