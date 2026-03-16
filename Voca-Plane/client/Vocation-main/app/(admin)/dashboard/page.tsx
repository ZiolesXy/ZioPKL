import { 
  Users, 
  CreditCard, 
  DollarSign, 
  Plane,
  CheckCircle2,
  Clock
} from "lucide-react"
import { 
  Card, 
  CardContent, 
  CardHeader, 
  CardTitle, 
  CardDescription 
} from "@/components/ui/card"
import { 
  Table, 
  TableBody, 
  TableCell, 
  TableHead, 
  TableHeader, 
  TableRow 
} from "@/components/ui/table"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { getInformation } from "@/lib/api/MonitorApi"

export default async function DashboardPage() {
  const response = await getInformation();
  const info = response.data;

  const stats = [
    {
      title: "Total Users",
      value: info?.total_users || 0,
      desc: "Registered users",
      icon: <Users className="size-5 text-indigo-600" />,
      color: "bg-indigo-100",
    },
    {
      title: "Total Flights",
      value: info?.total_flights || 0,
      desc: "Available flights",
      icon: <Plane className="size-5 text-emerald-600" />,
      color: "bg-emerald-100",
    },
    {
      title: "Total Transactions",
      value: info?.total_transactions || 0,
      desc: "Overall transactions",
      icon: <CreditCard className="size-5 text-amber-600" />,
      color: "bg-amber-100",
    },
    {
      title: "Total Revenue",
      value: `Rp ${(info?.total_revenue || 0).toLocaleString("id-ID")}`,
      desc: "Overall revenue",
      icon: <DollarSign className="size-5 text-rose-600" />,
      color: "bg-rose-100",
    },
    {
      title: "Completed Bookings",
      value: info?.completed_bookings || 0,
      desc: "Successfully booked",
      icon: <CheckCircle2 className="size-5 text-teal-600" />,
      color: "bg-teal-100",
    },
    {
      title: "Pending Payments",
      value: info?.pending_payments || 0,
      desc: "Awaiting payment",
      icon: <Clock className="size-5 text-orange-600" />,
      color: "bg-orange-100",
    },
  ];

  return (
    <div className="p-5 space-y-8 bg-slate-50/50 min-h-screen">
      <div className="flex justify-between items-end">
        <div>
          <h1 className="text-3xl font-bold tracking-tight text-slate-900">Dashboard</h1>
          <p className="text-slate-500">Welcome back! Here's what's happening today.</p>
        </div>
        <Button className="bg-indigo-600 hover:bg-indigo-700 shadow-md">
          Export Report
        </Button>
      </div>

      <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-6">
        {stats.map((stat, i) => (
          <Card key={i} className="border-none shadow-sm hover:shadow-md transition-shadow">
            <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
              <CardTitle className="text-sm font-medium text-slate-500">{stat.title}</CardTitle>
              <div className={`p-2 rounded-lg ${stat.color}`}>
                {stat.icon}
              </div>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{stat.value}</div>
              <p className="text-xs text-slate-400 mt-1">{stat.desc}</p>
            </CardContent>
          </Card>
        ))}
      </div>

      <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-7">
        <Card className="lg:col-span-4 border-none shadow-sm">
          <CardHeader>
            <div className="flex justify-between items-center">
              <div>
                <CardTitle>Recent Transactions</CardTitle>
                <CardDescription>Latest flight bookings across the platform.</CardDescription>
              </div>
              <Button variant="outline" size="sm">View All</Button>
            </div>
          </CardHeader>
          <CardContent>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Customer</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead className="text-right">Amount</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {/* Keep commented out placeholder if no transactions yet */}
              </TableBody>
            </Table>
          </CardContent>
        </Card>

        <Card className="lg:col-span-3 border-none shadow-sm">
          <CardHeader>
            <CardTitle>Active Flights Status</CardTitle>
            <CardDescription>Current fleet operational status.</CardDescription>
          </CardHeader>
          <CardContent className="space-y-6">
            {/* Kept empty for now until flights data is implemented */}
          </CardContent>
        </Card>
      </div>
      
    </div>
  )
}
