"use client"

import { ShoppingCart } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Separator } from "@/components/ui/separator"
import {
    Sheet,
    SheetContent,
    SheetDescription,
    SheetFooter,
    SheetHeader,
    SheetTitle,
    SheetTrigger,
} from "@/components/ui/sheet"
import Image from "next/image"
import DeleteCartItemButton from "@/components/admin/DeleteItem"
import { useCart } from "@/context/CartContext"
import { Checkbox } from "@/components/ui/checkbox"
import { useState, useEffect } from "react"

import { CartItem } from "@/types/cart"
import { checkoutCart, deleteCartItemsMany } from "@/lib/api/cart"
import { getMyCoupon } from "@/lib/api/coupon"
import { notifySuccess, notifyError } from "@/lib/toast"
import { formatRupiah } from "@/lib/utils"
import { useUser } from "@/hooks/useUser"
import { useRouter } from "next/navigation"
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { CheckoutDetail } from "@/types/checkout"
import { useQueryClient } from "@tanstack/react-query"

export function CartSheet() {
    const { cart, itemCount, refreshCart } = useCart()
    const { data: user } = useUser()
    const router = useRouter() // eslint-disable-line @typescript-eslint/no-unused-vars
    const queryClient = useQueryClient()
    const [selectedItems, setSelectedItems] = useState<number[]>([])
    const [selectedAddressUid, setSelectedAddressUid] = useState<string>("")
    const [couponCode, setCouponCode] = useState<string>("")
    const [checkoutData, setCheckoutData] =
        useState<CheckoutDetail | null>(null)

    const allItemId = cart?.items?.map(item => item.id) || []

    const isAllSelected =
        allItemId.length > 0 &&
        allItemId.every(id => selectedItems.includes(id))

    const toggleSelectAll = () => {
        if (isAllSelected) {
            setSelectedItems([])
        } else {
            setSelectedItems(allItemId)
        }
    }

    const toggleItem = (id: number) => {
        setSelectedItems(prev =>
            prev.includes(id)
                ? prev.filter(itemId => itemId !== id)
                : [...prev, id]
        )
    }
    useEffect(() => {
        if (selectedItems.length === 0) {
            setCheckoutData(null);
        }
    }, [selectedItems]);
    const total =
        cart?.items
            ?.filter((item: CartItem) => selectedItems.includes(item.id))
            .reduce(
                (sum: number, item: CartItem) =>
                    sum + item.product.price * item.quantity,
                0
            ) || 0

    const subtotalDisplay = checkoutData?.subtotal ?? total
    const discountDisplay = checkoutData?.discount_amount ?? 0
    const finalTotalDisplay = checkoutData?.total_price ?? total
    const [isProcessing, setIsProcessing] = useState(false)

    const handleCheckout = async () => {
        if (selectedItems.length === 0) {
            notifyError("Pilih Terlebih Dahulu")
            return
        }

        if (!selectedAddressUid) {
            notifyError("Pilih Alamat Pengiriman Terlebih Dahulu")
            return
        }

       setIsProcessing(true)
        try {
            const response = await checkoutCart(
                selectedItems,
                selectedAddressUid,
                couponCode
            ) as { data: { data: CheckoutDetail } }

            const data: CheckoutDetail = response.data.data
            setCheckoutData(data)
            notifySuccess("Checkout berhasil ")
            queryClient.invalidateQueries({ queryKey: ["my-checkouts"] })
            
            refreshCart()
            setSelectedItems([])
            setCouponCode("")
        } catch (error) {
            console.error("Checkout error:", error)
            notifyError("Gagal melakukan checkout. Silakan coba lagi.")
        } finally {
            setIsProcessing(false)
        }
    }
    const handleApplyCoupon = async () => {
        if (!couponCode) return;

        setIsProcessing(true);
        try {
            const myCoupons = await getMyCoupon();
            const coupon = myCoupons?.find((c) => c.coupon_code === couponCode);

            if (!coupon) {
                notifyError("Kupon tidak ditemukan di koleksi Anda.");
                setCheckoutData(null);
                return;
            }

            if (coupon.used_at) {
                notifyError("Kupon ini sudah pernah digunakan.");
                setCheckoutData(null);
                return;
            }

            // Hitung diskon secara manual untuk preview
            let discountAmount = 0;
            if (coupon.coupon_type === "percentage") {
                discountAmount = Math.floor((total * coupon.value) / 100);
            } else {
                discountAmount = coupon.value;
            }

            // Simulasi data checkout untuk preview
            setCheckoutData({
                subtotal: total,
                discount_amount: discountAmount,
                total_price: Math.max(0, total - discountAmount),
                // Data dummy lainnya agar TS tidak komplain jika menggunakan type CheckoutDetail penuh
            } as CheckoutDetail);

            notifySuccess("Kupon berhasil diterapkan! ðŸŽ‰");
        } catch {
            setCheckoutData(null);
            notifyError("Gagal memeriksa kupon.");
        } finally {
            setIsProcessing(false);
        }
    };
    const handleDeleteSelected = async () => {
        if (selectedItems.length === 0) {
            notifyError("Pilih item dulu ")
            return
        }

        setIsProcessing(true)
        try {
            await deleteCartItemsMany(selectedItems)
            notifySuccess("Berhasil menghapus item terpilih")
            refreshCart()
            setSelectedItems([])
        } catch (error) {
            console.error("Delete many error:", error)
            notifyError("Gagal menghapus item. Silakan coba lagi.")
        } finally {
            setIsProcessing(false)
        }
    }


    return (
        <Sheet>
            <SheetTrigger asChild>
                <Button variant="outline" size="icon" className="relative">
                    <ShoppingCart className="h-4 w-4" />
                    {itemCount > 0 && (
                        <span className="absolute -top-1 -right-1 flex h-4 w-4 items-center justify-center rounded-full bg-teal-500 text-[10px] text-white">
                            {itemCount}
                        </span>
                    )}
                </Button>
            </SheetTrigger>

            <SheetContent className="flex w-full flex-col sm:max-w-md">
                <SheetHeader>
                    <SheetTitle>Keranjang Belanja</SheetTitle>
                    <SheetDescription>
                        Anda memiliki {itemCount} item di keranjang.
                    </SheetDescription>
                </SheetHeader>

                <div className="flex-1 overflow-y-auto py-4">

                    <div className="flex items-center gap-3 px-3 py-2 bg-slate-50 rounded-md mb-2">
                        <Checkbox
                            id="select-all"
                            checked={isAllSelected}
                            onCheckedChange={toggleSelectAll}
                        />
                        <label htmlFor="select-all" className="text-sm font-medium text-slate-600 cursor-pointer">
                            Pilih Semua Produk
                        </label>

                        <Button
                            variant="destructive"
                            onClick={handleDeleteSelected}
                            disabled={selectedItems.length === 0 || isProcessing}
                            className="ml-auto"
                        >
                            Hapus Item Terpilih
                        </Button>
                    </div>
                    {cart?.items?.map((item: CartItem) => (
                        <div key={item.id}>
                            <div className="p-3 flex items-center gap-4 py-3">
                                <Checkbox
                                    checked={selectedItems.includes(item.id)}
                                    onCheckedChange={() => toggleItem(item.id)}
                                    className="data-[state=checked]:bg-teal-600 border-slate-300"
                                />

                                <div className="relative h-16 w-16 shrink-0 rounded bg-slate-100 overflow-hidden border">
                                    <Image
                                        src={item.product.image_url || "/placeholder.png"}
                                        alt={item.product.name}
                                        fill
                                        className="object-cover"
                                        sizes="64px"
                                    />
                                </div>

                                <div className="flex flex-col gap-1 flex-1">
                                    <p className="text-sm font-medium line-clamp-1">{item.product.name}</p>
                                    <p className="text-xs text-muted-foreground">
                                        Qty: {item.quantity} x {formatRupiah(item.product.price)}
                                    </p>
                                </div>

                                <div className="flex shrink-0">
                                    <DeleteCartItemButton itemId={item.id} onDeleted={refreshCart} />
                                </div>
                            </div>
                            <Separator />
                        </div>
                    ))}
                </div>

                <SheetFooter className="flex-col gap-3 sm:flex-col">
                    <Separator />

                    <div className="flex flex-col gap-4 py-2">
                        <div className="space-y-2">
                            <Label htmlFor="address">Alamat Pengiriman</Label>
                            <Select onValueChange={setSelectedAddressUid} value={selectedAddressUid}>
                                <SelectTrigger id="address" className="w-full">
                                    <SelectValue placeholder="Pilih Alamat" />
                                </SelectTrigger>
                                <SelectContent>
                                    {user?.address && user.address.length > 0 ? (
                                        user.address.map((addr) => (
                                            <SelectItem key={addr.uid} value={addr.uid}>
                                                <div className="flex flex-col items-start">
                                                    <span className="font-semibold">{addr.label}</span>
                                                    <span className="text-xs text-muted-foreground line-clamp-1">
                                                        {addr.recipient_name} | {addr.address_line}, {addr.city}
                                                    </span>
                                                </div>
                                            </SelectItem>
                                        ))
                                    ) : (
                                        <div className="p-2 text-sm text-muted-foreground text-center">
                                            Belum ada alamat. Silakan tambah di profile.
                                        </div>
                                    )}
                                </SelectContent>
                            </Select>
                        </div>

                        <div className="space-y-2">
                            <Label htmlFor="coupon">Kode Kupon (Opsional)</Label>
                            <div className="flex gap-2">
                                <Input
                                    id="coupon"
                                    placeholder="Masukkan kode kupon"
                                    value={couponCode}
                                    onChange={(e) => setCouponCode(e.target.value)}
                                />
                                <Button
                                    type="button"
                                    variant="outline"
                                    onClick={handleApplyCoupon}
                                    disabled={!couponCode || selectedItems.length === 0 || isProcessing}
                                >
                                    Gunakan
                                </Button>
                            </div>
                        </div>
                    </div>

                    <Separator />

                    <div className="space-y-1">
                        <div className="flex justify-between text-sm">
                            <span>Subtotal</span>
                            <span>{formatRupiah(subtotalDisplay)}</span>
                        </div>

                        {discountDisplay > 0 && (
                            <div className="flex justify-between text-sm text-green-600">
                                <span>Diskon</span>
                                <span>- {formatRupiah(discountDisplay)}</span>
                            </div>
                        )}

                        <div className="flex justify-between font-bold text-lg">
                            <span>Total</span>
                            <span>{formatRupiah(finalTotalDisplay)}</span>
                        </div>
                    </div>

                    <Button
                        disabled={isProcessing}
                        onClick={handleCheckout}
                        className="cursor-pointer w-full bg-teal-600 hover:bg-teal-700"
                    >
                        {isProcessing ? "Memproses..." : "Lanjutkan ke Pembayaran"}
                    </Button>

                </SheetFooter>
            </SheetContent>
        </Sheet>
    )
}



