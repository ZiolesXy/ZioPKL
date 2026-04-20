import { useQuery } from "@tanstack/react-query"
import { getMyCheckouts } from "@/lib/api/user"
import type { CheckoutNotification } from "@/types/checkout"

export function useNotifications() {
  return useQuery<CheckoutNotification[]>({
    queryKey: ["my-checkouts"],
    queryFn: getMyCheckouts,
    retry: false,
    refetchOnWindowFocus: false,
  })
}