import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { MyCoupon } from "@/types/coupon";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Percent, Copy } from "lucide-react";

interface MyCouponCardProps {
    coupon: MyCoupon;
    onCopy: (code: string) => void;
}

export function MyCouponCard({ coupon, onCopy }: MyCouponCardProps) {
    const isUsed = !!coupon.used_at;

    return (
        <Card className={`relative overflow-hidden group transition-all duration-300 border-border/50 ${isUsed ? 'opacity-60 bg-muted/30' : 'hover:shadow-lg hover:shadow-teal-500/5 hover:-translate-y-1 bg-card/50 backdrop-blur-sm'}`}>
            <div className="absolute top-0 right-0 w-24 h-24 bg-teal-500/5 text-teal-500 rounded-bl-full -z-10 flex items-start justify-end p-4">
                <Percent className="w-6 h-6 opacity-20" />
            </div>
            <CardHeader className="pb-3">
                <div className="flex justify-between items-start mb-2">
                    <Badge variant={isUsed ? "secondary" : "default"} className={!isUsed ? "bg-teal-500" : ""}>
                        {coupon.coupon_type === "percentage" ? `${coupon.value}%` : `Rp ${coupon.value.toLocaleString("id-ID")}`}
                    </Badge>
                    {isUsed && (
                        <Badge variant="outline" className="text-muted-foreground">Telah Digunakan</Badge>
                    )}
                </div>
                <CardTitle className="text-xl">{coupon.coupon_code}</CardTitle>
                <CardDescription>Diklaim pada {new Date(coupon.claimed_at).toLocaleDateString("id-ID")}</CardDescription>
            </CardHeader>
            <CardContent className="pb-4">
                <div className="flex items-center gap-2 bg-muted/50 p-2 lg:p-3 rounded-lg border border-border/50">
                    <code className="flex-1 font-mono text-center font-bold text-teal-600 tracking-wider">
                        {coupon.coupon_code}
                    </code>
                    <Button
                        size="icon"
                        variant="ghost"
                        className="h-8 w-8 hover:bg-teal-500/10 hover:text-teal-600"
                        onClick={() => onCopy(coupon.coupon_code)}
                        disabled={isUsed}
                    >
                        <Copy className="w-4 h-4" />
                    </Button>
                </div>
            </CardContent>
        </Card>
    );
}
