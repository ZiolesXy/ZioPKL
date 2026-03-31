import Sidebar from "@/components/SideBar";
export default function MainLayout({ children }: { children: React.ReactNode }) {
    return (
        <div className="flex min-h-screen">
            <div className="border-r">
                <Sidebar />
            </div>
            <main className="flex-1 p-4">
                {children}
            </main>
        </div>
    );
}