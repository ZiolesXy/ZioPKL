export default function Loading() {
    return (
        <div className="flex flex-col items-center justify-center min-h-[400px] gap-4 w-full">
            <div className="relative w-12 h-12">
                <div className="absolute top-0 left-0 w-full h-full border-4 border-slate-100 rounded-full"></div>
                <div className="absolute top-0 left-0 w-full h-full border-4 border-slate-800 rounded-full border-t-transparent animate-spin"></div>
            </div>
            <p className="text-sm font-medium text-slate-600 animate-pulse">Menyiapkan Panel Admin...</p>
        </div>
    )
}
