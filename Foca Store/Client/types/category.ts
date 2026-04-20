import { Product } from "./product"

export interface Category {
    id: number
    name: string
    slug: string
    icon_url?: string
    product_count?: number
    created_at?: string
    updated_at?: string
}

export interface CategoryWithProducts extends Category {
    products: Product[]
}