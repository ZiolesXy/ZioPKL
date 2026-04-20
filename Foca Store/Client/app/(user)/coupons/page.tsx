"use client";

import { useEffect, useState } from "react";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { AvailableCoupon, MyCoupon } from "@/types/coupon";
import { Loader2 } from "lucide-react";
import { toast } from "sonner";
import { getCoupon, getMyCoupon, claimCoupon } from "@/lib/api/coupon";
import { CouponHeader } from "@/components/coupons/CouponHeader";
import { AvailableCouponCard } from "@/components/coupons/AvailableCouponCard";
import { MyCouponCard } from "@/components/coupons/MyCouponCard";
import { EmptyState } from "@/components/coupons/EmptyState";

export default function CouponsPage() {
    const [activeTab, setActiveTab] = useState("available");
    const [availableCoupons, setAvailableCoupons] = useState<AvailableCoupon[]>([]);
    const [myCoupons, setMyCoupons] = useState<MyCoupon[]>([]);
    const [loading, setLoading] = useState(false);
    const [claimedStatus, setClaimedStatus] = useState<Record<number, boolean>>({})

    const fetchAvailable = async () => {
        try {
            setLoading(true);
            const data = await getCoupon();
            setAvailableCoupons(data || []);
        } catch (err) {
            console.error(err);
            toast.error("Gagal mengambil data kupon tersedia");
        } finally {
            setLoading(false);
        }
    };

    const fetchMine = async () => {
        try {
            setLoading(true);
            const data = await getMyCoupon();
            setMyCoupons(data || []);
        } catch (err) {
            console.error(err);
            toast.error("Gagal mengambil data kupon saya");
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        if (activeTab === "available") {
            fetchAvailable();
        } else {
            fetchMine();
        }
    }, [activeTab]);

    const handleCopyCode = (code: string) => {
        navigator.clipboard.writeText(code);
        toast.success("Kode kupon berhasil disalin!");
    };

    const handleClaim = async (id: number, code: string) => {
        try {
            setClaimedStatus((prev) => ({ ...prev, [id]: true })) // optimistik
            await claimCoupon(id)
            toast.success(`Kupon ${code} berhasil diklaim!`)
            fetchMine()
        } catch {
            setClaimedStatus((prev) => ({ ...prev, [id]: false }))
            toast.error("Gagal mengklaim kupon")
        }
    }

    const formatExpiryDate = (dateString: string | null) => {
        if (!dateString) return "Tanpa Batas Waktu";
        return new Date(dateString).toLocaleDateString("id-ID", {
            day: 'numeric',
            month: 'short',
            year: 'numeric'
        });
    };

    return (
        <div className="container max-w-5xl mx-auto p-4 md:p-6 lg:p-8 min-h-screen">
            <CouponHeader />

            <Tabs value={activeTab} onValueChange={setActiveTab} className="w-full">
                <TabsList className="grid w-full max-w-md grid-cols-2 mb-8 bg-muted/50 p-1 rounded-xl">
                    <TabsTrigger value="available" className="rounded-lg text-base font-medium transition-all">
                        Kupon Tersedia
                    </TabsTrigger>
                    <TabsTrigger value="my-coupons" className="rounded-lg text-base font-medium transition-all">
                        Kupon Saya
                    </TabsTrigger>
                </TabsList>

                {loading ? (
                    <div className="flex flex-col items-center justify-center py-20 space-y-4">
                        <Loader2 className="w-10 h-10 text-teal-500 animate-spin" />
                        <p className="text-muted-foreground">Memuat data...</p>
                    </div>
                ) : (
                    <>
                        <TabsContent value="available" className="animate-in fade-in-50 slide-in-from-bottom-2 duration-500">
                            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                                {availableCoupons.map((coupon) => (
                                    <AvailableCouponCard
                                        key={coupon.id}
                                        coupon={coupon}
                                        isClaimed={claimedStatus[coupon.id]}
                                        onClaim={handleClaim}
                                        formatExpiryDate={formatExpiryDate}
                                    />
                                ))}
                            </div>
                            {availableCoupons.length === 0 && (
                                <EmptyState
                                    title="Tidak ada kupon baru"
                                    description="Saat ini belum ada kupon baru yang bisa diklaim."
                                />
                            )}
                        </TabsContent>

                        <TabsContent value="my-coupons" className="animate-in fade-in-50 slide-in-from-bottom-2 duration-500">
                            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                                {myCoupons.map((coupon) => (
                                    <MyCouponCard
                                        key={coupon.id}
                                        coupon={coupon}
                                        onCopy={handleCopyCode}
                                    />
                                ))}
                            </div>
                            {myCoupons.length === 0 && (
                                <EmptyState
                                    title="Belum ada kupon"
                                    description="Anda belum memiliki atau belum pernah menukarkan kupon."
                                />
                            )}
                        </TabsContent>
                    </>
                )}
            </Tabs>
        </div>
    );
}
