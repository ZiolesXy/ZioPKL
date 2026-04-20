"use client"

import { useEffect, useState } from "react"
import { useParams } from "next/navigation"
import { getCheckoutDetail } from "@/lib/api/user"
import { CheckoutDetail } from "@/types/checkout"
import { formatRupiah } from "@/lib/utils"
import { Package, MapPin, Clock, CheckCircle2, MessageCircle, XCircle } from "lucide-react"
import Link from "next/link"
import Script from "next/script"


export default function OrderDetailPage() {
  const { uid } = useParams()
  const [order, setOrder] = useState<CheckoutDetail | null>(null)
  console.log(process.env.NEXT_PUBLIC_MIDTRANS_CLIENT_KEY)
  useEffect(() => {
    if (uid && typeof uid === "string") {
      getCheckoutDetail(uid).then(setOrder)
    }
  }, [uid])

  const handlePayment = () => {
    if (!window.snap) {
      alert("Snap belum siap")
      return
    }

    if (!order?.snap_token) {
      alert("Snap token tidak tersedia")
      return
    }

    window.snap.pay(order.snap_token, {
      onSuccess: function () {
        window.location.reload()
      },
      onPending: function () {
        window.location.reload()
      },
      onError: function () {
        alert("Pembayaran gagal")
      },
    })
  }

  if (!order) return <div className="p-10 text-center animate-pulse">Memuat data...</div>

  return (
    <div className="max-w-5xl mx-auto py-10 px-4 space-y-6">
      <div className={`p-6 rounded-2xl flex items-start gap-4 border ${order.status === 'pending'
        ? 'bg-amber-50 border-amber-200 text-amber-800'
        : order.status === 'approved'
          ? 'bg-teal-50 border-teal-200 text-teal-800'
          : 'bg-red-50 border-red-200 text-red-800'
        }`}>
        {order.status === 'pending' && <Clock className="mt-1" />}
        {order.status === 'approved' && <CheckCircle2 className="mt-1" />}
        {order.status === 'rejected' && <XCircle className="mt-1" />}

        <div>
          <h2 className="font-bold text-lg">
            {order.status === 'pending' && 'Pesanan Sedang Diverifikasi'}
            {order.status === 'approved' && 'Pesanan Disetujui'}
            {order.status === 'rejected' && 'Pesanan Ditolak'}
          </h2>
          <p className="text-sm opacity-90">
            {order.status === 'pending' && 'Admin akan segera mengecek pesanan Anda. Mohon tunggu konfirmasi selanjutnya.'}
            {order.status === 'approved' && 'Pesanan Anda telah disetujui oleh admin dan sedang dalam tahap proses.'}
            {order.status === 'rejected' && 'Maaf, pesanan Anda telah ditolak oleh admin. Silakan hubungi admin untuk informasi lebih lanjut.'}
          </p>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <div className="lg:col-span-2 space-y-6">
          {/* Item List */}
         {/* Daftar Produk */}
          <section className="bg-white border rounded-xl overflow-hidden shadow-sm">
            <div className="px-5 py-4 border-b font-bold flex items-center gap-2 bg-gray-50/50">
              <Package size={18} className="text-teal-600" /> 
              <span>Daftar Produk</span>
            </div>
            <div className="p-5 divide-y">
              {order.items.map((item) => (
                <div key={item.id} className="py-4 first:pt-0 last:pb-0 flex justify-between items-center">
                  <div>
                    <p className="font-semibold text-gray-800">{item.product.name}</p>
                    <p className="text-sm text-gray-500">{item.quantity} Unit x {formatRupiah(item.price)}</p>
                  </div>
                  <p className="font-bold text-gray-900">{formatRupiah(item.price * item.quantity)}</p>
                </div>
              ))}
            </div>
          </section>

          <section className="bg-white border rounded-xl p-5 shadow-sm space-y-3">
            <h3 className="font-bold flex items-center gap-2 border-b pb-3">
              <MapPin size={18} /> Lokasi Pengiriman
            </h3>
            <div className="text-sm text-gray-600">
              <p className="font-bold text-gray-900">{order.address.recipient_name}</p>
              <p>{order.address.phone}</p>
              <p className="mt-2 italic">{order.address.full_address}</p>
            </div>
          </section>
        </div>

       <aside className="space-y-4">
          <div className="bg-white border rounded-2xl p-6 shadow-sm sticky top-6">
            <h3 className="font-bold mb-4 text-lg">Ringkasan Pesanan</h3>
            <div className="space-y-3 pb-4 border-b text-sm">
              <div className="flex justify-between text-gray-500">
                <span>Total Item</span>
                <span className="text-gray-900 font-medium">{order.items.length} Produk</span>
              </div>
              <div className="flex justify-between text-gray-500">
                <span>Total Harga</span>
                <span>{formatRupiah(order.total_price)}</span>
              </div>
            </div>
            
            <div className="flex justify-between py-5">
              <span className="font-bold">Total Tagihan</span>
              <span className="font-bold text-teal-600 text-2xl">{formatRupiah(order.total_price)}</span>
            </div>

            <div className="space-y-3">
              {order.status === "pending" && (
                <button
                  onClick={handlePayment}
                  className="w-full bg-teal-600 hover:bg-teal-700 text-white rounded-xl py-3.5 font-bold shadow-lg shadow-teal-100 transition-all active:scale-[0.98]"
                >
                  Bayar Sekarang
                </button>
              )}

              <Link href="https://wa.link/lswl9x" target="_blank" className="block">
                <button className="w-full bg-white border border-gray-200 text-gray-700 rounded-xl py-3 flex items-center justify-center gap-2 hover:bg-gray-50 transition-all font-medium">
                  <MessageCircle size={18} /> Hubungi Admin
                </button>
              </Link>
            </div>

            <p className="text-[10px] text-gray-400 mt-6 text-center uppercase tracking-widest">
              Order UID: {order.uid}
            </p>
          </div>
        </aside>
      </div>
      <Script
        src="https://app.sandbox.midtrans.com/snap/snap.js"
        data-client-key={process.env.NEXT_PUBLIC_MIDTRANS_CLIENT_KEY}
        strategy="afterInteractive"
      />
    </div>
  )
}