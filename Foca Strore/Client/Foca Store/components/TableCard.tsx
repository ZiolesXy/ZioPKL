import { Card, CardContent, CardHeader, CardTitle } from "./ui/card"
import {
  Table,
  TableBody,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import type { ReactNode } from "react"

export default function TableCard<TItem>({
  title,
  headerRight,
  className,
  columns,
  items,
  getRowKey,
  renderRow,
}: {
  title: string
  headerRight?: ReactNode
  className?: string
  columns: Array<{ label: string; className?: string }>
  items: TItem[]
  getRowKey: (item: TItem, index: number) => string | number
  renderRow: (item: TItem, index: number) => ReactNode
}) {
  return (
    <Card className={className ?? "w-full h-105"}>
      <CardHeader className="border-b flex justify-between">
        <CardTitle>{title}</CardTitle>
        {headerRight}
      </CardHeader>
      <CardContent className="h-85 overflow-hidden">
        <div className="h-full w-full overflow-auto">
          <Table className="w-full table-fixed">
            <TableHeader>
              <TableRow>
                {columns.map((col) => (
                  <TableHead key={col.label} className={col.className}>
                    {col.label}
                  </TableHead>
                ))}
              </TableRow>
            </TableHeader>
            <TableBody>
              {items.map((item, idx) => (
                <TableRow key={getRowKey(item, idx)}>{renderRow(item, idx)}</TableRow>
              ))}
            </TableBody>
          </Table>
        </div>
      </CardContent>
    </Card>
  )
}
