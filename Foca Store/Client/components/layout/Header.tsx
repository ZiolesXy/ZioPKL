"use client"
import { useEffect, useState } from 'react'
import { CommandDemo } from '@/components/shared/Command'
import { User } from '@/types/auth'
import { CartSheet } from "@/components/cart/CartSheet"
import { UserItem } from './UserItem'
import Link from 'next/link'
import { clientApi } from '@/lib/client-api'
import { Product } from '@/types/product'
import NotificationBell from './Notification'
import { Menu } from 'lucide-react'
import { Sheet, SheetContent, SheetTrigger, SheetHeader, SheetTitle } from '@/components/ui/sheet'
import { Button } from '@/components/ui/button'

function Header({ user }: { user?: User }) {
    const [products, setProducts] = useState<Product[]>([])
    const isAuthenticated = Boolean(user)
    const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false)

    useEffect(() => {
        clientApi.get<{ data: { entries: Product[] } }>("/products")
            .then(res => {
                const entries = res.data?.entries ?? []
                setProducts(entries)
            })
            .catch(err => {
                console.error("Gagal fetch products untuk search:", err)
            })
    }, [])


    return (
        <header className="sticky top-0 z-50 w-full border-b border-slate-200 bg-white/80 backdrop-blur-md shadow-sm transition-all animate-in fade-in slide-in-from-top-4 duration-500">
            <div className="container mx-auto px-4 h-16 flex items-center justify-between gap-4">

                {/* Logo & Mobile Menu */}
                <div className="flex items-center gap-3">
                    <Sheet open={isMobileMenuOpen} onOpenChange={setIsMobileMenuOpen}>
                        <SheetTrigger asChild>
                            <Button variant="ghost" size="icon" className="md:hidden text-slate-600 hover:bg-slate-100">
                                <Menu className="h-6 w-6" />
                                <span className="sr-only">Toggle menu</span>
                            </Button>
                        </SheetTrigger>
                        <SheetContent side="left" className="w-[300px] sm:w-[400px] "
                        onOpenAutoFocus={(e) => e.preventDefault()}>
                            <SheetHeader>
                                <SheetTitle className="text-left text-teal-600 text-2xl font-black">Voca</SheetTitle>
                            </SheetHeader>
                            <div className="py-6 space-y-6">
                                <div className="block md:hidden">
                                    <CommandDemo products={products} />
                                </div>
                                <div className="flex flex-col gap-4 ml-2.5">
                                    <Link href="/" onClick={() => setIsMobileMenuOpen(false)} className="text-lg font-medium text-slate-700 hover:text-teal-600 transition-colors">Home</Link>
                                    <Link href="/category" onClick={() => setIsMobileMenuOpen(false)} className="text-lg font-medium text-slate-700 hover:text-teal-600 transition-colors">Kategori</Link>
                                    <Link href="/product" onClick={() => setIsMobileMenuOpen(false)} className="text-lg font-medium text-slate-700 hover:text-teal-600 transition-colors">Produk</Link>
                                </div>
                                {!isAuthenticated && (
                                    <div className="flex flex-col gap-3 pt-6 border-t border-slate-100">
                                        <Link
                                            href="/login"
                                            onClick={() => setIsMobileMenuOpen(false)}
                                            className="w-full rounded-md bg-teal-600 px-3 py-3 text-center text-sm font-medium text-white hover:bg-teal-700 transition shadow-sm"
                                        >
                                            Login
                                        </Link>
                                        <Link
                                            href="/register"
                                            onClick={() => setIsMobileMenuOpen(false)}
                                            className="w-full rounded-md border border-slate-200 px-3 py-3 text-center text-sm font-medium text-slate-700 hover:bg-slate-50 transition"
                                        >
                                            Register
                                        </Link>
                                    </div>
                                )}
                            </div>
                        </SheetContent>
                    </Sheet>

                    <Link href="/" className="font-extrabold text-2xl tracking-tighter text-teal-600 hover:opacity-90 transition-opacity">
                        Voca
                    </Link>
                </div>

                {/* Desktop Search */}
                <div className="hidden md:block flex-1 max-w-xl px-4 animate-in fade-in duration-700">
                    <CommandDemo products={products} />
                </div>

                {/* Right Actions */}
                <div className="flex items-center gap-1 sm:gap-2">
                    {isAuthenticated ? (
                        <>
                            <NotificationBell />
                            <CartSheet />
                            <UserItem />
                        </>
                    ) : (
                        <div className="hidden md:flex items-center gap-2">
                            <Link
                                href="/login"
                                className="rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90 transition-colors shadow-sm"
                            >
                                Login
                            </Link>
                            <Link
                                href="/register"
                                className="rounded-md border border-slate-200 bg-background px-4 py-2 text-sm font-medium text-slate-700 hover:bg-slate-50 transition-colors"
                            >
                                Register
                            </Link>
                        </div>
                    )}
                </div>
            </div>
        </header>
    )
}

export default Header
