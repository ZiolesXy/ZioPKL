import React from 'react'
import { UserDataTable } from '@/components/admin/UsersTableMonitor'
import { Button } from '@/components/ui/button'
import { ChevronRight, PlusCircle, Users } from 'lucide-react'
import Link from 'next/link'
import { Separator } from '@/components/ui/separator'
import { info } from 'console'
import { getInformation } from '@/lib/api/MonitorApi'

const userStats = await getInformation();

async function UsersMonitoringPage() {
  return (
    <div className="flex-1 space-y-8 pt-2">
      <div className="flex flex-col space-y-2">
        <nav className="flex items-center space-x-2 text-sm text-muted-foreground">
          <Link href="/dashboard" className="hover:text-primary transition-colors">
            Dashboard
          </Link>
          <ChevronRight className="h-4 w-4" />
          <span className="font-medium text-foreground">Users</span>
        </nav>

        <div className="flex items-center justify-between">
          <div>
            <h2 className="text-3xl font-bold tracking-tight">User Management</h2>
            <p className="text-muted-foreground">
              Monitor user accounts, roles, and access permissions for the flight booking system.
            </p>
          </div>
          <Button className="gap-2">
            <PlusCircle className="h-4 w-4" />
            Add New User
          </Button>
        </div>
      </div>

      <Separator />

      {/* Stats Overview (Optional but Recommended for Admin) */}
      <div className="grid gap-4 md:grid-cols-3">
        <div className="rounded-xl border bg-card p-6 text-card-foreground shadow-sm">
          <div className="flex flex-row items-center justify-between space-y-0 pb-2">
            <h3 className="text-sm font-medium">Total Users</h3>
            <Users className="h-4 w-4 text-muted-foreground" />
          </div>
          <div className="text-2xl font-bold">{userStats?.total_users || 0}</div> {/* Sesuai total di JSON Meta  */}
          <p className="text-xs text-muted-foreground">Active accounts in database</p>
        </div>
      </div>

      {/* Table Section */}
      <div className="rounded-xl border bg-background shadow-sm">
        <div className="p-6">
          <div className="flex items-center justify-between pb-4">
            <h3 className="text-lg font-semibold">Database Records</h3>
          </div>
          <UserDataTable />
        </div>
      </div>
    </div>
  )
}

export default UsersMonitoringPage