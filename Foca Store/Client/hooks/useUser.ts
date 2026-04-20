import { useQuery } from "@tanstack/react-query"
import { getUserClient } from "@/lib/api/user"

export const useUser = () => {
    return useQuery({
        queryKey: ["user-profile"],
        queryFn: getUserClient,
        retry: false,
        refetchOnWindowFocus: false,
        staleTime: 1000 * 60 * 5, // 5 minutes
    })
}
