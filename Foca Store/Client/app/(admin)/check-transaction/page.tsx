"use client"

import React, { useState, useEffect, useCallback } from 'react'
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { formatRupiah } from '@/lib/utils'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { ShoppingCart, CheckCircle2, AlertCircle, Loader2, RefreshCw } from "lucide-react"
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog"
import { getCheckout } from '@/lib/api/verify'

type OrderStatus = 'pending' | 'approved' | 'rejected'

interface ApiProduct {
    id: number
    name: string
}

interface ApiOrderItem {
    id: number
    quantity: number
    price: number
    product: ApiProduct
}

interface ApiUser {
    id: number
    name: string
    email: string
}

interface ApiOrder {
    id: number
    user: ApiUser
    total_price: number
    status: OrderStatus
    items: ApiOrderItem[]
    CreatedAt: string
    UpdatedAt: string
}

export default function CheckProductsPage() {
    const [orders, setOrders] = useState<ApiOrder[]>([])
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState<string | null>(null)

    const fetchOrders = useCallback(async () => {
        try {
            setLoading(true)
            setError(null)
            const data = await getCheckout()
            setOrders((data ?? []) as ApiOrder[])
        } catch (err: unknown) {
            if(err instanceof Error) {
                setError(err?.message)
            } else {
                setError("Gagal memuat data pesanan.")
            }
        } finally {
            setLoading(false)
        }
    }, [])

    useEffect(() => {
        fetchOrders()
    }, [fetchOrders])


    const formatDate = (dateStr: string) => {
        return new Date(dateStr).toLocaleDateString('id-ID', {
            day: '2-digit',
            month: 'short',
            year: 'numeric',
            hour: '2-digit',
            minute: '2-digit',
        })
    }
    const pendingCount = orders.filter(o => o.status === 'pending').length
    const approvedCount = orders.filter(o => o.status === 'approved').length
    const totalRevenue = orders
        .filter(o => o.status === 'approved')
        .reduce((acc, curr) => acc + curr.total_price, 0)

    const statusMap = {
        pending: {
            label: "Menunggu",
            className: "bg-yellow-100 text-yellow-800 hover:bg-yellow-100"
        },
        approved: {
            label: "Disetujui",
            className: "bg-teal-100 text-teal-800 hover:bg-teal-100"
        },
        rejected: {
            label: "Ditolak",
            className: "bg-red-100 text-red-800 hover:bg-red-100"
        }
    }

    return (
       <div className="p-4 md:p-8 lg:p-12 xl:p-8 flex flex-col gap-10 md:gap-12 lg:gap-16 bg-white min-h-screen w-full mx-auto">
            <div className="flex justify-between items-center">
                <div>
                    <h1 className="text-3xl font-bold text-slate-800 tracking-tight">Verifikasi Pesanan</h1>
                    <p className="text-slate-500 mt-1">Kelola dan tinjau pesanan masuk dari pelanggan.</p>
                </div>
                <div className="flex gap-2">
                    <Button variant="outline">Ekspor Data</Button>
                    <Button
                        className="bg-teal-600 hover:bg-teal-700"
                        onClick={fetchOrders}
                        disabled={loading}
                    >
                        {loading ? (
                            <Loader2 className="h-4 w-4 animate-spin mr-2" />
                        ) : (
                            <RefreshCw className="h-4 w-4 mr-2" />
                        )}
                        Refresh Data
                    </Button>
                </div>
            </div>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">

                <Card className="border-l-4 border-l-yellow-500 shadow-sm hover:shadow-md transition-shadow">
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium text-slate-600">Menunggu Verifikasi</CardTitle>
                        <AlertCircle className="h-4 w-4 text-yellow-500" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold text-slate-800">{pendingCount}</div>
                        <p className="text-xs text-slate-500">Pesanan perlu tindakan</p>
                    </CardContent>
                </Card>

                <Card className="border-l-4 border-l-teal-600 shadow-sm hover:shadow-md transition-shadow">
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium text-slate-600">Pesanan Disetujui</CardTitle>
                        <CheckCircle2 className="h-4 w-4 text-teal-600" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold text-slate-800">{approvedCount}</div>
                        <p className="text-xs text-slate-500">Total pesanan sukses</p>
                    </CardContent>
                </Card>

                <Card className="border-l-4 border-l-blue-600 shadow-sm hover:shadow-md transition-shadow">
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium text-slate-600">Total Pendapatan (Appr.)</CardTitle>
                        <ShoppingCart className="h-4 w-4 text-blue-600" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold text-slate-800">{formatRupiah(totalRevenue)}</div>
                        <p className="text-xs text-slate-500">Dari pesanan yang disetujui</p>
                    </CardContent>
                </Card>
            </div>

            <Card className="shadow-md border-slate-100">
                <CardHeader>
                    <CardTitle>Daftar Pesanan Terbaru</CardTitle>
                    <CardDescription>
                        Daftar lengkap transaksi yang masuk ke sistem.
                    </CardDescription>
                </CardHeader>
                <CardContent>
                    {loading && (
                        <div className="flex flex-col items-center justify-center py-16 gap-3">
                            <Loader2 className="h-8 w-8 animate-spin text-teal-600" />
                            <p className="text-slate-500 text-sm">Memuat data pesanan...</p>
                        </div>
                    )}

                    {!loading && error && (
                        <div className="flex flex-col items-center justify-center py-16 gap-3">
                            <AlertCircle className="h-8 w-8 text-red-500" />
                            <p className="text-red-500 text-sm">{error}</p>
                            <Button variant="outline" size="sm" onClick={fetchOrders}>
                                <RefreshCw className="h-4 w-4 mr-2" />
                                Coba Lagi
                            </Button>
                        </div>
                    )}

                    {!loading && !error && orders.length === 0 && (
                        <div className="flex flex-col items-center justify-center py-16 gap-3">
                            <ShoppingCart className="h-8 w-8 text-slate-300" />
                            <p className="text-slate-400 text-sm">Belum ada pesanan masuk.</p>
                        </div>
                    )}

                    {!loading && !error && orders.length > 0 && (
                        <Table>
                            <TableHeader>
                                <TableRow className="bg-slate-50/50">
                                    <TableHead className="w-25">ID Order</TableHead>
                                    <TableHead>Pelanggan</TableHead>
                                    <TableHead>Detail Item</TableHead>
                                    <TableHead>Total</TableHead>
                                    <TableHead>Tanggal</TableHead>
                                    <TableHead>Status</TableHead>
                                </TableRow>
                            </TableHeader>
                            <TableBody>

                                {orders.length === 0 ? (
                                    <TableRow>
                                        <TableCell colSpan={7} className="text-center py-10 text-slate-400">
                                            Tidak ada pesanan yang perlu diverifikasi.
                                        </TableCell>
                                    </TableRow>
                                ) :
                                    orders.map((order) => (
                                        <TableRow key={order.id} className="hover:bg-slate-50/50 transition-colors">
                                            <TableCell className="font-medium text-slate-700">#{order.id}</TableCell>
                                            <TableCell>
                                                <div className="flex flex-col">
                                                    <span className="font-medium text-slate-800">{order.user.name}</span>
                                                    <span className="text-xs text-slate-400">{order.user.email}</span>
                                                </div>
                                            </TableCell>
                                            <TableCell>
                                                <Dialog>
                                                    <DialogTrigger asChild>
                                                        <div className="cursor-pointer text-sm text-slate-600 space-y-1 hover:text-teal-600">
                                                            {order.items.slice(0, 3).map((item) => (
                                                                <div key={item.id} className="text-xs">
                                                                    {item.product.name} (x{item.quantity})
                                                                </div>
                                                            ))}

                                                            {order.items.length > 3 && (
                                                                <div className="text-xs text-slate-400 italic">
                                                                    +{order.items.length - 3} item lainnya (klik)
                                                                </div>
                                                            )}
                                                        </div>
                                                    </DialogTrigger>

                                                    <DialogContent>
                                                        <DialogHeader>
                                                            <DialogTitle>Semua Item — Order #{order.id}</DialogTitle>
                                                        </DialogHeader>

                                                        <div className="space-y-2 max-h-100 overflow-y-auto">
                                                            {order.items.map((item) => (
                                                                <div key={item.id} className="flex justify-between border p-2 rounded">
                                                                    <span>{item.product.name} (x{item.quantity})</span>
                                                                    <span>{formatRupiah(item.price)}</span>
                                                                </div>
                                                            ))}
                                                        </div>
                                                    </DialogContent>
                                                </Dialog>
                                            </TableCell>

                                            <TableCell className="font-bold text-teal-600">
                                                {formatRupiah(order.total_price)}
                                            </TableCell>
                                            <TableCell className="text-slate-500">{formatDate(order.CreatedAt)}</TableCell>
                                            <TableCell>
                                                <Badge className={statusMap[order.status].className}>
                                                    {statusMap[order.status].label}
                                                </Badge>
                                            </TableCell>
                                
                                        </TableRow>
                                    ))}
                            </TableBody>
                        </Table>
                    )}
                </CardContent>
            </Card>
        </div>
    )
}