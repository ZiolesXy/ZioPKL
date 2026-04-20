import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function formatRupiah(value: number) {
  return new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency: "IDR",
  }).format(value)
}

export class ApiError extends Error {
  status?: number

  constructor(message: string, status?: number){
    super(message)
    this.status = status
  }
}

export function handleApiError(error: unknown): never {
  if (error instanceof Error) {
    throw error
  }
  throw new ApiError("Terjadi kesalahan tak terduga")
}

/**
 * Safely extract error message from unknown errors (axios-style or standard).
 */
export function getErrorMessage(error: unknown, fallback = "Terjadi kesalahan"): string {
  if (error instanceof Error) {
    // Check for axios-style response.data.message
    const errObj = error as unknown as { response?: { data?: { message?: string } } }
    if (typeof errObj.response?.data?.message === "string") {
      return errObj.response.data.message
    }
    return error.message
  }
  return fallback
}