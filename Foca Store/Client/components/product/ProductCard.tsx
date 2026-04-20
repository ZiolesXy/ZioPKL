"use client"
import { Card, CardContent, CardFooter } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { ShoppingCart, Package, Loader2, Tag, Plus, Minus } from 'lucide-react'
import { useCart } from '@/context/CartContext'
import { notifyError, notifySuccess } from '@/lib/toast'
import Image from 'next/image'
import Link from 'next/link'
import { useState } from 'react'
import { formatRupiah } from '@/lib/utils'
import { CartProductCardProps } from "@/types/cart"
import { useRouter } from "next/navigation"

export default function ProductCard({ id, slug, name, category, price, stock, image }: CartProductCardProps) {
    const { addToCart } = useCart()
    const [qty, setQty] = useState(1)
    const [isLoading, setIsLoading] = useState(false)
    const [isAdding, setIsAdding] = useState(false)
    const router = useRouter()
    async function handleConfirmAdd() {
        setIsLoading(true)
        try {
            await addToCart(id, qty)
            notifySuccess(`${name} berhasil ditambah!`)
            setIsAdding(false)
            setQty(1)
        } catch (err) {
            console.error("Gagal menambahkan ke keranjang:", err)
            notifyError("Gagal menambahkan ke keranjang")
            router.push("/login")
        } finally {
            setIsLoading(false)
        }
    }

    return (
        <Card className='group overflow-hidden border-none shadow-md transition-all hover:shadow-xl bg-white flex flex-col'>
            <Link href={`/product/${slug}`} className="block">
                <div className='relative aspect-square w-full overflow-hidden bg-gray-50 p-3 sm:p-4'>
                    <Image
                        src={image}
                        alt={name}
                        fill
                        className='object-contain transition-transform duration-500 group-hover:scale-110'
                        sizes='(max-width: 786px) 100vw, 33vw'
                    />
                </div>

                <CardContent className="p-3 sm:p-4 space-y-3 flex-1 flex flex-col justify-between">
                    <div className="space-y-2">
                        <div className="flex items-center gap-1">
                            <Badge variant="outline" className="text-[10px] font-semibold text-teal-600 bg-teal-50 border-teal-200">
                                {typeof category === "string" ? category : category?.name ?? "Umum"}
                            </Badge>
                        </div>

                        <h3 className="font-bold text-sm sm:text-base line-clamp-2 text-slate-800 leading-tight group-hover:text-teal-700 transition-colors">
                            {name}
                        </h3>
                    </div>

                    <div className="flex flex-col sm:flex-row sm:items-end sm:justify-between gap-1 pt-2">
                        <p className="text-lg font-black tracking-tight text-slate-900">
                            {formatRupiah(price)}
                        </p>
                        <Badge
                            variant={stock > 0 ? "secondary" : "destructive"}
                            className="text-[10px] uppercase font-bold tracking-wider bg-slate-100 text-slate-600 w-fit"
                        >
                            {stock > 0 ? `Sisa ${stock}` : "Habis"}
                        </Badge>
                    </div>
                </CardContent>
            </Link>

            <CardFooter className="p-3 sm:p-4 pt-0">
                {!isAdding ? (
                    <Button
                        disabled={stock === 0 || isLoading}
                        onClick={() => setIsAdding(true)}
                        className="w-full bg-teal-600 hover:bg-teal-700 transition-all duration-300 shadow-sm"
                    >
                        <ShoppingCart className="w-4 h-4 mr-2" />
                        {stock > 0 ? "Tambah" : "Stok Kosong"}
                    </Button>
                ) : (
                    <div className="flex flex-col w-full gap-2 animate-in fade-in slide-in-from-bottom-2 duration-300">
                        <div className="flex items-center justify-between border rounded-lg p-1 bg-gray-50">
                            <Button
                                variant="ghost"
                                size="icon"
                                className="h-9 w-9 sm:h-8 sm:w-8 text-teal-600"
                                onClick={() => setQty(prev => Math.max(1, prev - 1))}
                            >
                                <Minus className="w-4 h-4" />
                            </Button>

                            <span className="font-bold text-sm">{qty}</span>

                            <Button
                                variant="ghost"
                                size="icon"
                                className="h-9 w-9 sm:h-8 sm:w-8 text-teal-600"
                                onClick={() => setQty(prev => Math.min(stock, prev + 1))}
                            >
                                <Plus className="w-4 h-4" />
                            </Button>
                        </div>

                        <div className="flex gap-2">
                            <Button
                                variant="outline"
                                className="flex-1 text-xs h-9"
                                onClick={() => {
                                    setIsAdding(false)
                                    setQty(1)
                                }}
                            >
                                Batal
                            </Button>
                            <Button
                                className="flex-1 bg-teal-600 hover:bg-teal-700 text-xs h-9"
                                onClick={handleConfirmAdd}
                                disabled={isLoading}
                            >
                                {isLoading ? <Loader2 className="w-4 h-4 animate-spin" /> : "Konfirmasi"}
                            </Button>
                        </div>
                    </div>
                )}
            </CardFooter>
        </Card>
    )
}