"use client"
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { useRouter } from "next/navigation"
import { User, LogOut, Loader2, Ticket } from "lucide-react"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"

import { useUser } from "@/hooks/useUser"

export function UserItem() {
    const router = useRouter()
    const { data: user, isLoading, isError } = useUser()

    const handleLogout = async () => {
        try {
            await fetch("/api/logout", {
                method: "POST",
            })
        } finally {
            router.push("/")
            router.refresh()
        }
    }

    return (
        <DropdownMenu>
            <DropdownMenuTrigger asChild>
                <button className="relative h-10 w-10 rounded-full hover:bg-muted p-1 flex items-center justify-center border">
                    {isLoading ? (
                        <Loader2 className="h-5 w-5 animate-spin text-muted-foreground" />
                    ) : (
                        <Avatar className="h-10 w-10">
                            <AvatarImage src={user?.profile_image_url} alt={user?.name} />
                            <AvatarFallback>
                                {user?.name?.charAt(0) || "U"}
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
                            <p className="text-xs text-destructive">Failed to load user</p>
                        ) : (
                            <>
                                <p className="text-sm font-medium leading-none">{user?.name || "User"}</p>
                                <p className="text-xs leading-none text-muted-foreground">{user?.email || ""}</p>
                            </>
                        )}
                    </div>
                </DropdownMenuLabel>
                <DropdownMenuSeparator />
                <DropdownMenuItem onClick={() => router.push('/update-profile')}>
                    <User className="mr-2 h-4 w-4" />
                    <span>Edit Profile</span>
                </DropdownMenuItem>
                <DropdownMenuItem onClick={() => router.push('/coupons')}>
                    <Ticket className="mr-2 h-4 w-4" />
                    <span>Kupon</span>
                </DropdownMenuItem>
                <DropdownMenuSeparator />
                <DropdownMenuItem className="text-red-600 focus:text-red-600" onClick={handleLogout}>
                    <LogOut className="mr-2 h-4 w-4" />
                    <span>Log out</span>
                </DropdownMenuItem>
            </DropdownMenuContent>
        </DropdownMenu>
    )
}
