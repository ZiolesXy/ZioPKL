import React from 'react'

import { TransactionTableData } from "@/components/admin/TransactionTableData"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { ChevronRight } from "lucide-react"
import Link from "next/link"

function TransactionsPage() {
  return (
    <div className="p-5 space-y-4">
      <div className="flex flex-col space-y-2">
        <nav className="flex items-center space-x-2 text-sm text-muted-foreground">
          <Link href="/dashboard" className="hover:text-primary transition-colors">
            Dashboard
          </Link>
          <ChevronRight className="h-4 w-4" />
          <span className="font-medium text-foreground">Transaction Monitoring</span>
        </nav>

        <div>
          <h2 className="text-3xl font-bold tracking-tight">Transaction Monitoring</h2>
          <p className="text-muted-foreground">
            Monitor and manage all flight booking transactions.
          </p>
        </div>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Daftar Transaksi</CardTitle>
        </CardHeader>
        <CardContent>
          <TransactionTableData />
        </CardContent>
      </Card>
    </div>
  )
}

export default TransactionsPage