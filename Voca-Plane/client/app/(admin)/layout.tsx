import { SidebarProvider, SidebarTrigger } from "@/components/ui/sidebar"
import { AppSidebar } from "@/components/admin/SideBar"
import { AdminSearch } from "@/components/AdminSearch"
import { ThemeToggle } from "@/components/ThemeToggle"
export default function Layout({ children }: { children: React.ReactNode }) {
  return (

    <SidebarProvider>
      <AppSidebar />
      <main className="w-full flex flex-col min-h-screen bg-slate-50 dark:bg-slate-950/50">
        <header className="flex h-20 shrink-0 items-center justify-between gap-4 px-8 bg-white/80 dark:bg-slate-900/80 backdrop-blur-md border-b border-slate-200/50 dark:border-white/5 sticky top-0 z-10">
          <div className="flex items-center gap-6">
            <div className="p-2 rounded-xl bg-slate-100 dark:bg-white/5 hover:bg-slate-200 dark:hover:bg-white/10 transition-colors">
              <SidebarTrigger className="text-slate-500 hover:text-primary transition-colors h-5 w-5" />
            </div>
            <div className="h-6 w-[1px] bg-slate-200 dark:bg-white/10" />
            <AdminSearch />
          </div>
          <div className="flex items-center gap-4">
             <div className="bg-primary/5 text-primary text-[10px] font-black uppercase tracking-[0.2em] px-4 py-1.5 rounded-full border border-primary/10">
               Production Ready
             </div>
             <div className="h-8 w-[1px] bg-slate-200 dark:bg-white/10 mx-2" />
             <ThemeToggle />
          </div>
        </header>
        <div className="flex-1 p-4 lg:p-6">
          {children}
        </div>
      </main>
    </SidebarProvider>

  )
}