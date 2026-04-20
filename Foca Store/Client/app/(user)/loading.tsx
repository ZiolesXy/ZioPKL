export default function Loading() {
    return (
        <div className="flex flex-col items-center justify-center min-h-[400px] gap-4">
            <div className="relative w-16 h-16">
                <div className="absolute top-0 left-0 w-full h-full border-4 border-teal-100 rounded-full"></div>
                <div className="absolute top-0 left-0 w-full h-full border-4 border-teal-600 rounded-full border-t-transparent animate-spin"></div>
            </div>
            <div className="space-y-2 text-center">
                <p className="text-lg font-medium animate-pulse text-teal-800">Loading...</p>
                <p className="text-sm text-muted-foreground">Sedang mengambil data terbaru untuk Anda.</p>
            </div>
        </div>
    )
}
