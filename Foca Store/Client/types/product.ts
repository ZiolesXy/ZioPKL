import { Category } from "./category"

export interface Product {
    id: number
    name: string
    slug: string
    description: string
    image_url: string 
    price: number
    stock: number
    category: Category | string
    created_at?: string
    updated_at?: string
}

export type ProductCardProps = Product
