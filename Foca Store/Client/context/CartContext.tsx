"use client"

import { createContext, useContext, useState, useCallback, useEffect, ReactNode } from "react"
import { getCart, addToCart as apiAddToCart } from "@/lib/api/cart"
import { Cart } from "@/types/cart"

interface CartContextType {
    cart: Cart | null
    itemCount: number
    refreshCart: () => Promise<void>
    addToCart: (productId: number, quantity: number) => Promise<void>
    isAdding: boolean
}

const CartContext = createContext<CartContextType | undefined>(undefined)

export function CartProvider({ children }: { children: ReactNode }) {
    const [cart, setCart] = useState<Cart | null>(null)
    const [isAdding, setIsAdding] = useState(false)

    const refreshCart = useCallback(async () => {
        try {
            const data = await getCart()
            setCart(data ?? null)
        } catch (err) {
            console.error("Failed to fetch cart:", err)
        }
    }, [])

    const addToCart = useCallback(async (productId: number, quantity: number) => {
        setIsAdding(true)
        try {
            await apiAddToCart(productId, quantity)
            await refreshCart()
        } catch (err) {
            console.error("Failed to add to cart:", err)
            throw err
        } finally {
            setIsAdding(false)
        }
    }, [refreshCart])

    // Fetch cart on initial mount
    useEffect(() => {
        refreshCart()
    }, [refreshCart])

    const itemCount = cart?.items?.length || 0

    return (
        <CartContext.Provider value={{ cart, itemCount, refreshCart, addToCart, isAdding }}>
            {children}
        </CartContext.Provider>
    )
}

export function useCart() {
    const context = useContext(CartContext)
    if (context === undefined) {
        throw new Error("useCart must be used within a CartProvider")
    }
    return context
}
