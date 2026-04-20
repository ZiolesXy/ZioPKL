import { SidebarProvider, SidebarTrigger } from "@/components/ui/sidebar"
import { AdminSidebar } from "@/components/layout/SideBar"
import { Separator } from "@/components/ui/separator" 
import { SearchBar } from "@/components/admin/SearchBar"
import Link from "next/link"
import { DynamicBreadcrumb } from "@/components/admin/PathBreadCrumb"
import { UserItem } from "@/components/layout/UserItem"
import { ThemeToggle } from "@/components/shared/ThemeToggle"

export default function Layout({ children }: { children: React.ReactNode }) {
  return (
    <SidebarProvider>
      <AdminSidebar />
      <main className="flex-1 flex flex-col min-h-screen">
        <header className="sticky top-0 z-10 flex h-16 shrink-0 items-center border-b bg-background/95 backdrop-blur px-4">
          <div className="flex items-center justify-between w-full">
            
            <div className="flex items-center gap-2">
              <SidebarTrigger className="-ml-1" />
              <Separator orientation="vertical" className="mx-2 h-4" />
              
              <nav className="flex items-center gap-2 text-sm font-medium">
                <Link 
                  href="/overview" 
                  className="text-muted-foreground hover:text-foreground transition-colors"
                >
                  Control
                </Link> 
                <DynamicBreadcrumb />
              </nav>
            </div>

            <div className="flex items-center gap-3 sm:gap-4">
              <div className="hidden md:block">
                <SearchBar />
              </div>
              
              <div className="flex items-center border-l pl-3 sm:pl-4 gap-3 sm:gap-4 border-border">
                <ThemeToggle />
                <UserItem />
              </div>
            </div>
          </div>
        </header>

        <div className="flex flex-1 flex-col gap-4 p-6 pt-4">
          {children}
        </div>
      </main>
    </SidebarProvider>
  )
}