import { SidebarProvider, SidebarTrigger } from "@/components/ui/sidebar"
import { AppSidebar } from "@/components/admin/SideBar"
import { AdminSearch } from "@/components/AdminSearch"
import { ThemeToggle } from "@/components/ThemeToggle"
export default function Layout({ children }: { children: React.ReactNode }) {
  return (

    <SidebarProvider>
      <AppSidebar />
      <main className="w-full flex flex-col min-h-screen">
        <header className="flex h-16 shrink-0 items-center justify-between gap-2 border-b px-4 bg-white dark:bg-slate-900 sticky top-0 z-10">
          <div className="flex items-center gap-2">
            <SidebarTrigger className="-ml-1 text-slate-500 hover:text-indigo-600" />
            <AdminSearch />
          </div>
          <ThemeToggle />
        </header>
        <div className="flex-1 p-6 bg-slate-50/50">
          {children}
        </div>
      </main>
    </SidebarProvider>

  )
}