import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { ChevronRight, Calendar, Ticket, Plane, CreditCard, ArrowRight } from "lucide-react";
import Link from "next/link";
import { getTransactionByUser } from "@/lib/api/TransactionApi";

export default async function MyBookingsPage() {
    const transactions = await getTransactionByUser();
    const data = transactions.data || [];

    return (
        <div className="bg-slate-50/50 min-h-screen pb-20">
            {/* Minimalist Header */}
            <div className="bg-white border-b border-slate-200 mb-8">
                <div className="container mx-auto px-4 max-w-5xl py-8">
                    <div className="flex items-center gap-3">
                        <div className="bg-blue-600 p-2 rounded-lg">
                            <Ticket className="size-5 text-white" />
                        </div>
                        <div>
                            <h1 className="text-2xl font-bold text-slate-900">Pesanan Saya</h1>
                            <p className="text-slate-500 text-sm">Daftar riwayat dan status reservasi penerbangan Anda</p>
                        </div>
                    </div>
                </div>
            </div>

            <div className="container mx-auto px-4 max-w-5xl">
                {data.length === 0 ? (
                    <div className="bg-white rounded-2xl shadow-sm p-16 text-center border border-slate-200">
                        <div className="bg-slate-50 p-6 rounded-full mb-4 inline-block">
                            <Plane className="size-12 text-slate-300" />
                        </div>
                        <h3 className="text-xl font-bold text-slate-800">Belum ada penerbangan</h3>
                        <p className="text-slate-500 mb-6">Mulai petualangan Anda dengan memesan tiket pertama.</p>
                        <Link href="/flight" className="bg-blue-600 hover:bg-blue-700 text-white font-semibold py-2.5 px-6 rounded-lg transition-all inline-flex items-center gap-2">
                            Cari Penerbangan
                        </Link>
                    </div>
                ) : (
                    <div className="grid gap-4">
                        {data.map((booking: any) => {
                            const flight = booking.flight;
                            const isPending = booking.payment_status === 'PENDING';
                            const isSuccess = ['PAID', 'SUCCESS'].includes(booking.payment_status);

                            return (
                                <Card key={booking.code} className="group hover:border-blue-200 transition-all duration-200 shadow-sm border-slate-200 overflow-hidden">
                                    <CardContent className="p-0">
                                        <div className="flex flex-col sm:flex-row items-stretch">
                                            {/* Info Utama - Route & Date */}
                                            <Link href={`/my-bookings/${booking.code}`} className="flex-1 p-5 flex flex-col md:flex-row md:items-center gap-6">
                                                <div className="min-w-[120px]">
                                                    <span className="text-[10px] font-bold text-slate-400 uppercase tracking-wider block mb-1">Booking Code</span>
                                                    <span className="font-mono font-bold text-blue-600">{booking.code.substring(0, 8)}</span>
                                                </div>

                                                <div className="flex flex-1 items-center gap-4">
                                                    <div className="flex flex-col items-center">
                                                        <span className="text-lg font-bold text-slate-800">{flight?.origin?.city || 'Origin'}</span>
                                                        <span className="text-xs text-slate-400 uppercase font-medium">{flight?.origin?.code || '---'}</span>
                                                    </div>

                                                    <div className="flex flex-col items-center flex-1 px-4 max-w-[100px]">
                                                        <Plane className="size-4 text-slate-300 mb-1" />
                                                        <div className="h-px w-full bg-slate-100 relative">
                                                            <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 size-1.5 bg-slate-200 rounded-full" />
                                                        </div>
                                                    </div>

                                                    <div className="flex flex-col items-center">
                                                        <span className="text-lg font-bold text-slate-800">{flight?.destination?.city || 'Dest'}</span>
                                                        <span className="text-xs text-slate-400 uppercase font-medium">{flight?.destination?.code || '---'}</span>
                                                    </div>
                                                </div>

                                                <div className="md:border-l border-slate-100 md:pl-6 flex flex-col justify-center">
                                                    <div className="flex items-center gap-2 text-slate-600 text-sm mb-1">
                                                        <Calendar className="size-3.5" />
                                                        <span className="font-medium">
                                                            {new Date(booking.created_at).toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' })}
                                                        </span>
                                                    </div>
                                                    <Badge variant="outline" className={`w-fit text-[10px] h-5 border-none ${isSuccess ? 'bg-emerald-50 text-emerald-600' :
                                                            isPending ? 'bg-amber-50 text-amber-600' : 'bg-red-50 text-red-600'
                                                        }`}>
                                                        {booking.payment_status}
                                                    </Badge>
                                                </div>
                                            </Link>

                                            {/* Action/Price Section */}
                                            <div className="bg-slate-50/50 border-t sm:border-t-0 sm:border-l border-slate-100 p-5 flex sm:flex-col items-center justify-between sm:justify-center min-w-[180px] gap-3">
                                                <div className="text-left sm:text-center">
                                                    <p className="text-[10px] text-slate-400 font-bold uppercase mb-0.5">Total Bayar</p>
                                                    <p className="font-bold text-slate-900">
                                                        Rp {booking.total_price.toLocaleString('id-ID')}
                                                    </p>
                                                </div>

                                                {isPending && booking.payment_url ? (
                                                    <Link
                                                        href={booking.payment_url}
                                                        target="_blank"
                                                        className="bg-amber-500 hover:bg-amber-600 text-white text-xs font-bold py-2 px-4 rounded-md transition-colors flex items-center gap-2"
                                                    >
                                                        Bayar Sekarang <CreditCard className="size-3" />
                                                    </Link>
                                                ) : (
                                                    <Link href={`/my-bookings/${booking.code}`} className="text-blue-600 text-xs font-bold hover:underline flex items-center gap-1">
                                                        Detail <ChevronRight className="size-3" />
                                                    </Link>
                                                )}
                                            </div>
                                        </div>
                                    </CardContent>
                                </Card>
                            );
                        })}
                    </div>
                )}
            </div>
        </div>
    );
}