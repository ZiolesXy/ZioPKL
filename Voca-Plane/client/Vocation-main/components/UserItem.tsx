"use client"

import { useQuery } from "@tanstack/react-query" // Sangat disarankan
import { useRouter } from "next/navigation" // Gunakan navigation, bukan router
import Cookies from "js-cookie"
import { getProfile } from "@/lib/api/UserApi"
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { User, LogOut, Loader2 } from "lucide-react"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"

function UserItem() {
    const router = useRouter()

    // Mengambil data menggunakan React Query
    const { data, isLoading, isError } = useQuery({
        queryKey: ["user-profile"],
        queryFn: getProfile,
        retry: false
    })

    const user = data?.data // Sesuaikan dengan struktur response API kamu

    const handleLogout = () => {
        Cookies.remove("token")
        Cookies.remove("access_token")
        Cookies.remove("role")
        router.push("/login")
        router.refresh()
    }

    return (
        <DropdownMenu>
            <DropdownMenuTrigger asChild>
                <button className="relative h-10 w-10 rounded-full hover:bg-muted p-1 flex items-center justify-center border">
                    {isLoading ? (
                        <Loader2 className="h-5 w-5 animate-spin text-muted-foreground" />
                    ) : (
                        <Avatar className="h-10 w-10">
                            <AvatarImage src={user?.name} alt={"none"} />
                            <AvatarFallback>
                                {user?.name?.charAt(0) || "H"} 
                            </AvatarFallback>
                        </Avatar>
                    )}
                </button>
            </DropdownMenuTrigger>

            <DropdownMenuContent className="w-56" align="end" forceMount>
                <DropdownMenuLabel className="font-normal">
                    <div className="flex flex-col space-y-1">
                        {isLoading ? (
                            <>
                                <div className="h-4 w-24 bg-muted animate-pulse rounded" />
                                <div className="h-3 w-32 bg-muted animate-pulse rounded" />
                            </>
                        ) : isError ? (
                            <p className="text-xs text-destructive">Session expired</p>
                        ) : (
                            <>
                                <p className="text-sm font-medium leading-none">{user?.name || "User"} </p>
                                <p className="text-xs leading-none text-muted-foreground">{user?.email || ""} </p>
                            </>
                        )}
                    </div>
                </DropdownMenuLabel>
                <DropdownMenuSeparator />
                
                {/* Route Update Profile bisa diakses semua role */}
                <DropdownMenuItem onClick={() => router.push('/update-profile')}>
                    <User className="mr-2 h-4 w-4" />
                    <span>Edit Profile</span>
                </DropdownMenuItem>

                <DropdownMenuSeparator />
                <DropdownMenuItem onClick={handleLogout} className="text-red-600 focus:text-red-600">
                    <LogOut className="mr-2 h-4 w-4" />
                    <span>Log out</span>
                </DropdownMenuItem>
            </DropdownMenuContent>
        </DropdownMenu>
    )
}

export default UserItem