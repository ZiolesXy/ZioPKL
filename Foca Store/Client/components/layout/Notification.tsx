"use client"

import { BellIcon, CheckCircle2, Clock, XCircle, Package, CheckCheck } from "lucide-react"
import { Button } from "@/components/ui/button"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { ScrollArea } from "@/components/ui/scroll-area"
import { useNotifications } from "@/hooks/useNotifications"
import { CheckoutNotification } from "@/types/checkout"
import { formatDistanceToNow } from "date-fns"
import { id } from "date-fns/locale"
import Link from "next/link"

const STATUS_MAP = {
  approved: {
    icon: <CheckCircle2 className="h-4 w-4 text-emerald-500" />,
    bg: "bg-emerald-50",
    label: "Disetujui",
  },
  rejected: {
    icon: <XCircle className="h-4 w-4 text-red-500" />,
    bg: "bg-red-50",
    label: "Ditolak",
  },
  pending: {
    icon: <Clock className="h-4 w-4 text-amber-500" />,
    bg: "bg-amber-50",
    label: "Menunggu",
  },
} as const

const getStatusDetails = (status: string) =>
  STATUS_MAP[status as keyof typeof STATUS_MAP] ?? STATUS_MAP.pending


export default function NotificationBell() {

  const { data, isLoading } = useNotifications()
  const notifications = data ?? []

  const unreadCount = notifications.filter(item => item.status === "pending").length

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="outline" size="icon" className="relative hover:bg-slate-50 transition-colors">
          <BellIcon className="h-[1.2rem] w-[1.2rem]" />
          {unreadCount > 0 && (
            <span className="absolute -top-1 -right-1 flex h-4 w-4 items-center justify-center rounded-full bg-red-500 text-[10px] font-bold text-white ring-2 ring-white">
              {unreadCount}
            </span>
          )}
        </Button>
      </DropdownMenuTrigger>

      <DropdownMenuContent align="end" className="w-80 p-0 shadow-xl border-slate-200">
        <div className="flex items-center justify-between p-4 border-b bg-slate-50/50">
          <h4 className="text-sm font-semibold text-slate-900">Notifikasi</h4>
          <Button variant="ghost" size="sm" className="h-auto p-0 text-xs text-teal-600 hover:text-teal-700 hover:bg-transparent">
            <CheckCheck className="mr-1 h-3 w-3" /> Tandai semua dibaca
          </Button>
        </div>

        <ScrollArea className="h-80">
          {isLoading && (
            <div className="flex flex-col items-center justify-center py-10 gap-2">
              <div className="h-5 w-5 animate-spin rounded-full border-2 border-teal-500 border-t-transparent" />
              <p className="text-xs text-muted-foreground">Memuat data...</p>
            </div>
          )}

          {!isLoading && notifications.length === 0 && (
            <div className="flex flex-col items-center justify-center py-12 px-4 text-center">
              <div className="h-12 w-12 rounded-full bg-slate-100 flex items-center justify-center mb-3">
                <Package className="h-6 w-6 text-slate-400" />
              </div>
              <p className="text-sm font-medium text-slate-900">Belum ada kabar</p>
              <p className="text-xs text-muted-foreground">Semua notifikasi pesananmu akan muncul di sini.</p>
            </div>
          )}

          {notifications.map((item: CheckoutNotification) => {
            const details = getStatusDetails(item.status)
            return (
              <DropdownMenuItem
                key={item.id}
                asChild
                className={`cursor-pointer border-b last:border-0 p-4 focus:bg-slate-50 transition-colors ${item.status === 'pending' ? 'bg-teal-50/30' : ''}`}
              >
                <Link href={`/orders/${item.uid}`} className="flex gap-3 items-start w-full">
                  <div className={`p-2 rounded-full shrink-0 ${details.bg}`}>
                    {details.icon}
                  </div>
                  <div className="flex flex-col gap-1 flex-1">
                    <p className="text-sm leading-none">
                      <span className="font-semibold text-slate-900">Pesanan {item.uid}</span>
                      <span className="text-slate-600"> telah {details.label.toLowerCase()}</span>
                    </p>
                    <p className="text-[11px] text-muted-foreground flex items-center">
                      {formatDistanceToNow(new Date(item.created_at), { addSuffix: true, locale: id })}
                    </p>
                  </div>
                  {item.status === 'pending' && (
                    <div className="h-2 w-2 bg-teal-500 rounded-full mt-1" />
                  )}
                </Link>
              </DropdownMenuItem>
            )
          })}
        </ScrollArea>

        <Link href="/notifications" className="block p-3 text-center text-xs font-medium text-slate-500 hover:text-teal-600 hover:bg-slate-50 border-t">
          Lihat Semua Aktivitas
        </Link>
      </DropdownMenuContent>
    </DropdownMenu>
  )
}