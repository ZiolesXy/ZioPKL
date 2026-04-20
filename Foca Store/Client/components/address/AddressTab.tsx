"use client"

import { useEffect, useState, useCallback } from "react"
import { getAddresses } from "@/lib/api/user"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Loader2, MapPin, Phone as PhoneIcon, Home } from "lucide-react"
import { AddressData } from "@/types/auth"
import AddAddresses from "./AddAddresses"
import { Trash2 } from "lucide-react"
import { deleteAddress } from "@/lib/api/user"
import EditAddres from "./EditAddres"


export default function AddressTab() {

    const [addresses, setAddresses] = useState<AddressData[]>([])
    const [loading, setLoading] = useState(true)

    const fetchAddress = useCallback(async () => {
        setLoading(true)
        try {
            const data = await getAddresses()
            setAddresses(data || [])
        } catch (error) {
            console.error(error)
        } finally {
            setLoading(false)
        }
    }, [])

    useEffect(() => {
        fetchAddress()
    }, [fetchAddress])

    const handleDelete = async (uid: string) => {
        const confirmed = window.confirm("Hapus alamat ini?")
        if (!confirmed) return

        try {
            await deleteAddress(uid)
            fetchAddress()
        } catch (error) {
            console.error(error)
        }
    }

    return (
        <Card className="border-none shadow-md overflow-hidden rounded-3xl">
            <CardHeader className="flex flex-row items-center justify-between border-b bg-gray-50/50 px-6 py-4">
                <div>
                    <CardTitle className="text-xl font-bold text-teal-900 flex items-center gap-2">
                        <Home className="text-teal-600" size={20} /> Daftar Alamat
                    </CardTitle>
                    <p className="text-sm text-gray-500 mt-1">Kelola alamat pengiriman belanja Anda</p>
                </div>
                <AddAddresses onSuccess={fetchAddress} />
            </CardHeader>
            <CardContent className="p-6">
                {loading ? (
                    <div className="flex flex-col items-center justify-center py-20 gap-4">
                        <Loader2 className="animate-spin text-teal-600 h-8 w-8" />
                        <p className="text-sm text-gray-500 animate-pulse">Memuat data alamat...</p>
                    </div>
                ) : addresses.length === 0 ? (
                    <div className="text-center py-16 border-2 border-dashed border-gray-100 rounded-3xl bg-gray-50/30">
                        <div className="bg-gray-100 w-16 h-16 rounded-full flex items-center justify-center mx-auto mb-4">
                            <MapPin className="text-gray-400" size={32} />
                        </div>
                        <h3 className="text-lg font-semibold text-gray-900">Belum Ada Alamat</h3>
                        <p className="text-sm text-gray-500 mt-1 max-w-xs mx-auto">
                            Anda belum menambahkan alamat pengiriman. Tambahkan alamat sekarang untuk mempermudah checkout.
                        </p>
                    </div>
                ) : (
                    <div className="grid gap-4 sm:grid-cols-2">
                        {addresses.map((addr) => (
                            <div
                                key={addr.uid}
                                className={`group border-2 transition-all p-5 rounded-3xl relative hover:shadow-lg ${addr.is_primary
                                    ? "border-teal-500 bg-teal-50/50"
                                    : "border-gray-100 bg-white hover:border-teal-200"
                                    }`}
                            >
                                <div className="flex justify-between items-start mb-3">
                                    <Badge
                                        variant={addr.is_primary ? "default" : "outline"}
                                        className={addr.is_primary ? "bg-teal-600" : "text-gray-500 border-gray-200"}
                                    >
                                        {addr.label}
                                    </Badge>
                                    <div className="flex items-center gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
                                        <EditAddres address={addr} onSuccess={fetchAddress} />
                                        <button
                                            onClick={() => handleDelete(addr.uid)}
                                            className="text-gray-400 hover:text-red-600 transition-colors"
                                        >
                                            <Trash2 size={16} />
                                        </button>
                                    </div>

                                </div>

                                <p className="font-bold text-gray-900 text-lg mb-1">
                                    {addr.recipient_name}
                                </p>

                                <div className="flex items-center gap-2 text-gray-600 text-sm mb-3">
                                    <PhoneIcon size={14} className="text-gray-400" />
                                    {addr.phone}
                                </div>

                                <div className="flex items-start gap-2 text-gray-600 text-sm">
                                    <MapPin size={16} className="text-teal-600 mt-0.5 shrink-0" />
                                    <p className="leading-relaxed">
                                        {addr.address_line}<br />
                                        <span className="font-medium text-gray-900">
                                            {addr.city}, {addr.province}, {addr.postal_code}
                                        </span>
                                    </p>
                                </div>
                            </div>
                        ))}
                    </div>
                )}
            </CardContent>
        </Card>
    )
}
