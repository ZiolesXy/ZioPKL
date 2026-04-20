import { Ticket } from "lucide-react";

export function CouponHeader() {
    return (
        <div className="mb-10 text-center md:text-left space-y-4">
            <div className="inline-flex items-center justify-center p-3 bg-teal-500/10 rounded-full mb-4">
                <Ticket className="w-8 h-8 text-teal-500" />
            </div>
            <h1 className="text-3xl md:text-4xl font-extrabold tracking-tight text-teal-500">
                Pusat Kupon
            </h1>
            <p className="text-muted-foreground text-lg max-w-2xl">
                Klaim kupon yang tersedia dan gunakan untuk mendapatkan diskon terbaik saat checkout.
            </p>
        </div>
    );
}
