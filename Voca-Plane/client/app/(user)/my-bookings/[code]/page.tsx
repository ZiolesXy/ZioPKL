import { Button } from '@/components/ui/button';
import { Plane, QrCode, Download, Mail, ChevronLeft, Calendar, User, Armchair, Clock, Info } from 'lucide-react';
import { Separator } from '@/components/ui/separator';
import { getTransactionByCode } from '@/lib/api/TransactionApi';
import Link from 'next/link';
import { Transaction } from '@/lib/type/transaction';
import { Badge } from '@/components/ui/badge';

export default async function BookingDetailPage({ params }: { params: Promise<{ code: string }> }) {
  const { code: bookingCode } = await params;
  const response = await getTransactionByCode(bookingCode);
  const booking: Transaction = response.data;

  if (!booking) return (
    <div className="bg-slate-50 min-h-screen py-20 px-4 flex justify-center items-center">
      <div className="text-center">
        <h2 className="text-2xl font-bold text-slate-800 mb-2">Transaksi Tidak Ditemukan</h2>
        <Link href="/my-bookings" className="text-blue-600 hover:underline flex items-center justify-center gap-2">
          <ChevronLeft className="size-4" /> Kembali ke Pesanan Saya
        </Link>
      </div>
    </div>
  );

  const isPending = booking.payment_status === 'PENDING';
  const isSuccess = ['PAID', 'SUCCESS'].includes(booking.payment_status);

  return (
    <div className="bg-slate-50 min-h-screen py-10 px-4">
      <div className="container mx-auto max-w-3xl">
        <Link href="/my-bookings" className="inline-flex items-center gap-2 text-slate-500 hover:text-blue-600 mb-6 transition-colors font-medium">
          <ChevronLeft className="size-4" /> Kembali ke Daftar Pesanan
        </Link>

        {/* 1. PAYMENT ALERT (Untuk Pending) */}
        {isPending && (
          <div className="bg-amber-50 border border-amber-200 rounded-2xl p-4 mb-6 flex items-start gap-4 shadow-sm">
            <div className="bg-amber-100 p-2 rounded-xl text-amber-600">
              <Clock className="size-5" />
            </div>
            <div>
              <p className="text-amber-800 font-bold text-sm">Selesaikan Pembayaran Segera</p>
              <p className="text-amber-700 text-xs">
                Tiket ini akan otomatis dibatalkan pada: <span className="font-bold">{new Date(booking.expires_at).toLocaleString('id-ID')}</span> 
              </p>
            </div>
          </div>
        )}

        <div className="bg-white rounded-[2rem] shadow-2xl overflow-hidden border border-slate-100">
          {/* 2. HEADER TIKET */}
          <div className={`p-8 text-white relative ${isSuccess ? 'bg-blue-600' : isPending ? 'bg-slate-800' : 'bg-red-600'}`}>
             <div className="flex justify-between items-center relative z-10">
                <div className="flex flex-col items-start gap-2">
                    <h1 className="text-3xl font-black tracking-tight">E-Ticket</h1>
                    <Badge variant="outline" className={`px-2 py-0.5 text-[10px] font-bold uppercase tracking-wider border-none ${
                        isSuccess ? 'bg-emerald-500/20 text-white' :
                        isPending ? 'bg-amber-500/20 text-white' : 'bg-white/20 text-white'
                    }`}>
                        {booking.payment_status}
                    </Badge>
                </div>
                <div className="text-right">
                    <p className="text-white/60 text-[10px] font-bold uppercase tracking-widest">Booking ID</p>
                    <p className="text-xl font-mono font-bold">#{booking.code.substring(0, 8)} </p>
                </div>
             </div>
          </div>

          <div className="p-8">
            {/* 3. FLIGHT ROUTE */}
            <div className="flex justify-between items-center mb-10">
              <div className="flex-1">
                <p className="text-4xl font-black text-slate-900">{booking.flight?.origin?.code || '---'}</p>
                <p className="text-slate-500 font-medium">{booking.flight?.origin?.city || 'Origin'} </p>
              </div>
              <div className="flex flex-col items-center px-6">
                <Plane className="size-6 text-blue-600 mb-2" />
                <div className="h-[2px] w-24 bg-slate-100 relative">
                  <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 size-2 bg-blue-600 rounded-full" />
                </div>
                <p className="text-[10px] font-bold text-slate-400 mt-2 uppercase">{booking.flight_number} </p>
              </div>
              <div className="flex-1 text-right">
                <p className="text-4xl font-black text-slate-900">{booking.flight?.destination?.code || '---'}</p>
                <p className="text-slate-500 font-medium">{booking.flight?.destination?.city || 'Destination'} </p>
              </div>
            </div>

            {/* 4. PASSENGER LIST (Dynamic Loop) */}
            <div className="mb-10">
              <p className="text-[11px] font-bold text-slate-400 uppercase tracking-[0.2em] mb-4">Informasi Penumpang [cite: 79]</p>
              <div className="space-y-4">
                {booking.transactions_passanger?.map((passenger, index) => (
                  <div key={index} className="flex items-center justify-between p-4 bg-slate-50 rounded-2xl border border-slate-100 group hover:border-blue-200 transition-colors">
                    <div className="flex items-center gap-4">
                      <div className="bg-white p-2.5 rounded-xl border border-slate-200 text-slate-400 group-hover:text-blue-600 transition-colors">
                        <User className="size-5" />
                      </div>
                      <div>
                        <p className="font-bold text-slate-800 capitalize">{passenger.passenger_name} </p>
                        <p className="text-xs text-slate-500 font-medium uppercase tracking-tight">{passenger.nationality} • {passenger.passport_no} </p>
                      </div>
                    </div>
                    <div className="text-right">
                      <p className="text-[10px] font-bold text-slate-400 uppercase">Seat</p>
                       <p className="text-xs font-bold text-slate-400 uppercase">{passenger.class_name}</p>
                      <p className="font-black text-xl text-blue-600">{passenger.seat_number} </p>
                    </div>
                  </div>
                ))}
              </div>
            </div>
            
            {/* 4.5 BUYER INFO (Ordered By) */}
            {booking.user && (
              <div className="mb-10 p-4 bg-slate-50/50 rounded-2xl border border-slate-100 flex items-center gap-4">
                 <div className="bg-white p-2.5 rounded-xl text-slate-400 border border-slate-100">
                    <User className="size-4" />
                 </div>
                 <div>
                    <p className="text-[10px] text-slate-400 font-bold uppercase tracking-widest leading-none mb-1">Pemesanan Oleh</p>
                    <p className="font-bold text-slate-700 leading-none">{booking.user.name} <span className="text-slate-400 font-medium ml-1">· {booking.user.email}</span></p>
                 </div>
              </div>
            )}

            <Separator className="my-8 border-dashed border-2" />

            {/* 5. PRICE BREAKDOWN */}
            <div className="bg-slate-50 rounded-3xl p-6 mb-8">
               <p className="text-[11px] font-bold text-slate-400 uppercase tracking-[0.2em] mb-4">Rincian Harga </p>
               <div className="space-y-3">
                  <div className="flex justify-between text-sm">
                    <span className="text-slate-500">Harga Tiket ({booking.transactions_passanger?.length} Penumpang) [cite: 79]</span>
                    <span className="font-semibold text-slate-900">Rp {(booking.total_price + booking.discount).toLocaleString('id-ID')}</span>
                  </div>
                  {booking.discount > 0 && (
                    <div className="flex justify-between text-sm text-emerald-600 font-medium">
                      <span>Potongan Promo ({booking.promo_code}) [cite: 83]</span>
                      <span>- Rp {booking.discount.toLocaleString('id-ID')}</span>
                    </div>
                  )}
                  <Separator className="bg-slate-200" />
                  <div className="flex justify-between items-center">
                    <span className="font-bold text-slate-900">Total Bayar </span>
                    <span className="text-2xl font-black text-blue-600">Rp {booking.total_price.toLocaleString('id-ID')}</span>
                  </div>
               </div>
            </div>

            {/* 6. QR CODE (Only on Success) */}
            {isSuccess && (
              <div className="flex flex-col items-center justify-center p-6 border-2 border-dashed border-slate-200 rounded-3xl">
                <QrCode className="size-24 text-slate-800 mb-3" />
                <p className="text-[10px] font-bold text-slate-400 uppercase text-center">Tunjukkan QR Code ini kepada petugas check-in bandara</p>
              </div>
            )}
          </div>
        </div>

        {/* 7. ACTIONS */}
        <div className="mt-8 grid grid-cols-1 sm:grid-cols-2 gap-4">
          {isSuccess ? (
            <>
                <Button variant="outline" className="h-14 rounded-2xl font-bold bg-white border-slate-200">
                    <Download className="mr-2 size-5" /> Download E-Ticket
                </Button>
                <Button className="h-14 rounded-2xl font-bold bg-blue-600">
                    <Mail className="mr-2 size-5" /> Kirim Email
                </Button>
            </>
          ) : isPending ? (
            <Link href={booking.payment_url || '#'} className="col-span-2">
                <Button className="w-full h-14 rounded-2xl font-bold bg-amber-500 hover:bg-amber-600 shadow-lg shadow-amber-200">
                    Bayar Sekarang Rp {booking.total_price.toLocaleString('id-ID')} 
                </Button>
            </Link>
          ) : null}
        </div>
      </div>
    </div>
  );
}