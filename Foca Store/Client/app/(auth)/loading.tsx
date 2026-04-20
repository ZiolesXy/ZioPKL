export default function Loading() {
    return (
        <div className="flex flex-col items-center justify-center min-h-[400px] gap-4">
            <div className="w-8 h-8 border-2 border-teal-600 border-t-transparent rounded-full animate-spin"></div>
            <p className="text-sm font-medium text-muted-foreground">Menyiapkan autentikasi...</p>
        </div>
    )
}
