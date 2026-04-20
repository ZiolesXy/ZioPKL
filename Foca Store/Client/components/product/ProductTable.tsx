"use client"
import { useState } from "react"
import { Pencil, Trash2, MoreHorizontal } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table"
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import EditProduct from "@/components/product/EditProduct"
import { GenericDelete } from "@/components/admin/GenericDelete"
import { useRouter } from "next/navigation"
import Image from "next/image"
import { formatRupiah } from "@/lib/utils"
import { Product } from "@/types/product"
import { deleteProduct } from "@/lib/api/product"
import { div } from "three/src/nodes/math/OperatorNode.js"

export default function ProductTable({ products }: { products: Product[] }) {
    const [editProduct, setEditProduct] = useState<Product | null>(null)
    const router = useRouter()

    return (
        <div className="space-y-5">
           
            <div className="rounded-md border bg-white">
                <Table>
                    <TableHeader>
                        <TableRow className="bg-muted/50">
                            <TableHead className="w-20">Foto</TableHead>
                            <TableHead>Nama Produk</TableHead>
                            <TableHead>Harga</TableHead>
                            <TableHead>Stok</TableHead>
                            <TableHead>Status</TableHead>
                            <TableHead className="text-right">Aksi</TableHead>
                        </TableRow>
                    </TableHeader>
                    <TableBody>
                        {products.map((item) => (
                            <TableRow key={item.id} className="hover:bg-muted/20">
                                <TableCell>
                                    <Image
                                        src={item.image_url ?? "/placeholder.png"}
                                        alt={item.name}
                                        width={48}
                                        height={48}
                                        className="h-12 w-12 rounded-md object-cover border"
                                    />
                                </TableCell>

                                <TableCell className="font-medium">{item.name}</TableCell>
                                <TableCell>
                                    {formatRupiah(item.price)}
                                </TableCell>
                                <TableCell>{item.stock} Unit</TableCell>
                                <TableCell>
                                    {item.stock > 0 ? (
                                        <Badge
                                            variant="outline"
                                            className="text-emerald-600 border-emerald-200 bg-emerald-50"
                                        >
                                            Tersedia
                                        </Badge>
                                    ) : (
                                        <Badge variant="destructive">Habis</Badge>
                                    )}
                                </TableCell>
                                <TableCell className="text-right">
                                    <DropdownMenu>
                                        <DropdownMenuTrigger asChild>
                                            <Button variant="ghost" className="h-8 w-8 p-0">
                                                <MoreHorizontal className="h-4 w-4" />
                                            </Button>
                                        </DropdownMenuTrigger>
                                        <DropdownMenuContent align="end">
                                            <DropdownMenuLabel>Opsi</DropdownMenuLabel>
                                            <DropdownMenuItem
                                                className="cursor-pointer"
                                                onClick={() => setEditProduct(item)}
                                            >
                                                <Pencil className="mr-2 h-4 w-4" /> Edit
                                            </DropdownMenuItem>
                                            <DropdownMenuSeparator />
                                            <GenericDelete
                                                id={item.id}
                                                name={item.name}
                                                entityName="Produk"
                                                onDelete={async (id) => {
                                                    await deleteProduct(Number(id))
                                                    router.refresh()
                                                }}
                                                trigger={
                                                    <DropdownMenuItem
                                                        className="text-red-600 cursor-pointer focus:bg-red-50 focus:text-red-600"
                                                        onSelect={(e) => e.preventDefault()}
                                                    >
                                                        <Trash2 className="mr-2 h-4 w-4" /> Hapus
                                                    </DropdownMenuItem>
                                                }
                                            />
                                        </DropdownMenuContent>
                                    </DropdownMenu>
                                </TableCell>
                            </TableRow>
                        ))}
                    </TableBody>
                </Table>
            </div>

            {editProduct && (
                <EditProduct
                    product={editProduct}
                    open={!!editProduct}
                    onOpenChange={(open) => {
                        if (!open) setEditProduct(null)
                    }}
                    onSuccess={() => {
                        setEditProduct(null)
                        router.refresh()
                    }}
                />
            )}
        </div>
    )
}

