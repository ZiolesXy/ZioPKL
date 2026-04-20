import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { AvailableCoupon } from "@/types/coupon";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Percent, Ticket, CheckCircle2, AlertCircle, Clock, Sparkles } from "lucide-react";

interface AvailableCouponCardProps {
    coupon: AvailableCoupon;
    isClaimed: boolean;
    onClaim: (id: number, code: string) => void;
    formatExpiryDate: (date: string | null) => string;
}

export function AvailableCouponCard({
    coupon,
    isClaimed,
    onClaim,
    formatExpiryDate
}: AvailableCouponCardProps) {
    const isExpired = new Date(coupon.expires_at) < new Date();

    return (
        <Card className={`relative overflow-hidden group transition-all duration-300 border-border/50 hover:shadow-lg hover:shadow-teal-500/5 hover:-translate-y-1 bg-card/50 backdrop-blur-sm ${isClaimed || isExpired ? "opacity-75 grayscale-[0.2]" : ""}`}>
            <div className="absolute top-0 right-0 w-24 h-24 bg-teal-500/10 rounded-bl-full -z-10 transition-transform group-hover:scale-110 flex items-start justify-end p-4">
                {coupon.type === "percentage" ? (
                    <Percent className="w-6 h-6 text-teal-500 opacity-20" />
                ) : (
                    <Ticket className="w-6 h-6 text-teal-500 opacity-20" />
                )}
            </div>
            <CardHeader className="pb-3">
                <div className="flex justify-between items-start mb-2">
                    <Badge variant="outline" className="bg-teal-500/5 text-teal-600 border-teal-500/20">
                        {coupon.type === "percentage" ? `${coupon.value}%` : `Rp ${coupon.value.toLocaleString("id-ID")}`}
                    </Badge>
                    {isClaimed && (
                        <Badge variant="secondary" className="bg-green-500/10 text-green-600">
                            <CheckCircle2 className="w-3 h-3 mr-1" /> Diklaim
                        </Badge>
                    )}
                </div>
                <CardTitle className="text-xl">{coupon.code}</CardTitle>
                <CardDescription className="line-clamp-2">Sisa Kuota: {coupon.remaining} kupon</CardDescription>
            </CardHeader>
            <CardContent className="pb-4">
                <div className="space-y-2 text-sm text-muted-foreground">
                    <div className="flex items-center gap-2">
                        <AlertCircle className="w-4 h-4" />
                        <span>Min. Belanja: <span className="font-medium text-foreground">{coupon.minimum_purchase}</span></span>
                    </div>
                    <div className="flex items-center gap-2">
                        <Clock className="w-4 h-4" />
                        <span>
                            Berlaku hingga:{" "}
                            <span className="font-medium text-foreground">
                                {formatExpiryDate(coupon.expires_at)}
                            </span>
                        </span>
                    </div>
                </div>
            </CardContent>
            <CardFooter className="pt-0">
                <Button
                    onClick={() => onClaim(coupon.id, coupon.code)}
                    disabled={isClaimed || isExpired}
                    className="w-full bg-teal-500 hover:bg-teal-600 text-white"
                >
                    <span className="relative z-10 flex items-center gap-2">
                        {isClaimed ? "Sudah Diklaim" : isExpired ? "Kadaluarsa" : "Klaim Kupon"}
                        {!isClaimed && !isExpired && <Sparkles className="w-4 h-4" />}
                    </span>
                </Button>
            </CardFooter>
        </Card>
    );
}
