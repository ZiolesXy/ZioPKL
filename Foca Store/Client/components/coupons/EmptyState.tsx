import { Ticket } from "lucide-react";

interface EmptyStateProps {
    title: string;
    description: string;
}

export function EmptyState({ title, description }: EmptyStateProps) {
    return (
        <div className="text-center py-20 bg-muted/30 rounded-2xl border border-dashed border-border">
            <Ticket className="w-12 h-12 text-muted-foreground/50 mx-auto mb-4" />
            <h3 className="text-lg font-semibold mb-1">{title}</h3>
            <p className="text-muted-foreground">{description}</p>
        </div>
    );
}
