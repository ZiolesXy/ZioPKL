"use client"

import { useState, useEffect, useRef } from "react"
import { useActionState } from "react"
import { createProduct } from "./actions"
import { getCategories } from "@/lib/api/category"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Textarea } from "@/components/ui/textarea"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { ImageIcon } from "lucide-react"
import Image from "next/image"
import { Category } from "@/types/category"
import { toast } from "sonner"
import { useRouter } from "next/navigation"
const initialState = { message: "" }

export default function AddProductPage() {
  const [state, formAction, isPending] = useActionState(createProduct, initialState)
  const [preview, setPreview] = useState<string | null>(null)
  const [categories, setCategories] = useState<Category[]>([])
  const [selectedCategoryId, setSelectedCategoryId] = useState<string>("")
  const router = useRouter()
  const prevMessageRef = useRef(state?.message)
  const formRef = useRef<HTMLFormElement>(null)

  // Only perform side effects (toast/navigation), no setState calls
  useEffect(() => {
    if (state?.message && state.message !== prevMessageRef.current) {
      prevMessageRef.current = state.message
      if (state.message === "success") {
        toast.success("Product created successfully")
        // Reset form via DOM (not setState)
        formRef.current?.reset()
        router.push("/add-products")
      } else {
        toast.error("Error", { description: state.message })
      }
    }
  }, [state, router])

  useEffect(() => {
    async function fetchCategories() {
      const data = await getCategories()
      setCategories(data)
    }
    fetchCategories()
  }, [])

  const handleImageChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (file) {
      setPreview(URL.createObjectURL(file))
    }
  }

  return (
    <div className="p-8 bg-muted/20 min-h-screen">
      <div className="max-w-5xl mx-auto space-y-6">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold tracking-tight">Tambahkan Produk</h1>
            <p className="text-muted-foreground">Tambahkan item baru ke inventaris toko Anda.</p>
          </div>
        </div>

        <form ref={formRef} action={formAction} className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <div className="md:col-span-2 space-y-6">
            <Card>
              <CardHeader>
                <CardTitle>General Information</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="space-y-2">
                  <Label htmlFor="name">Product Name</Label>
                  <Input id="name" name="name" placeholder="e.g. Lenovo LOQ RTX 3060" required />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="description">Description</Label>
                  <Textarea id="description" name="description" placeholder="Write something about this product..." className="min-h-[200px]" required />
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Inventory & Pricing</CardTitle>
              </CardHeader>
              <CardContent className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="price">Price (IDR)</Label>
                  <Input id="price" name="price" type="number" placeholder="0" required />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="stock">Stock Quantity</Label>
                  <Input id="stock" name="stock" type="number" placeholder="0" required />
                </div>
              </CardContent>
            </Card>
          </div>

          <div className="space-y-6">
            <Card>
              <CardHeader>
                <CardTitle>Category</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-2">
                  <Label htmlFor="category_id">Product Category</Label>
                  <Select
                    value={selectedCategoryId}
                    onValueChange={(value) => setSelectedCategoryId(value)}
                  >
                    <SelectTrigger className="w-full">
                      <SelectValue placeholder="Select a category" />
                    </SelectTrigger>
                    <SelectContent>
                      {categories.map((category) => (
                        <SelectItem key={category.id} value={String(category.id)}>
                          {category.name}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                  <input type="hidden" name="category_id" value={selectedCategoryId} />
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Product Media</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="aspect-square rounded-md border-2 border-dashed flex flex-col items-center justify-center relative overflow-hidden bg-muted/50">
                  {preview ? (
                    <Image src={preview} alt="Preview" width={500} height={500} className="object-cover w-full h-full" />
                  ) : (
                    <div className="flex flex-col items-center text-muted-foreground">
                      <ImageIcon className="h-10 w-10 mb-2" />
                      <span className="text-xs">No image selected</span>
                    </div>
                  )}
                </div>
                <div className="mt-4">
                  <Input id="image" name="image" type="file" accept="image/*" onChange={handleImageChange} required className="cursor-pointer" />
                </div>
              </CardContent>
            </Card>

            <div className="space-y-3">
              <Button type="submit" className="w-full bg-teal-600 hover:bg-teal-700" disabled={isPending}>
                {isPending ? "Publishing..." : "Publish Product"}
              </Button>
              <Button variant="outline" type="button" className="w-full">
                Save as Draft
              </Button>
            </div>

            
          </div>
        </form>
      </div>
    </div>
  )
}