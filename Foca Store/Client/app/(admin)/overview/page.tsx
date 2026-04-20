import { GetOverview } from '@/lib/api/system'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { 
  Users, 
  Package, 
  Layers, 
  Ticket, 
  ShoppingCart, 
  Banknote, 
  Clock, 
  CheckCircle2, 
  XCircle 
} from 'lucide-react'

//nanti pakai format di lib
const formatCurrency = (value: number) => {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0,
  }).format(value)
}

async function OverviewPage() {
  const response = await GetOverview()
  const data = response.data

  const stats = [
    {
      title: "Total Revenue",
      value: formatCurrency(data.total_revenue),
      icon: Banknote,
      description: "Total pendapatan kotor",
      color: "text-green-600"
    },
    {
      title: "Total Orders",
      value: data.total_orders,
      icon: ShoppingCart,
      description: "Semua transaksi masuk",
      color: "text-blue-600"
    },
    {
      title: "Total Products",
      value: data.total_products,
      icon: Package,
      description: "Produk aktif di toko",
      color: "text-orange-600"
    },
    {
      title: "Total Users",
      value: data.total_users,
      icon: Users,
      description: "Pelanggan terdaftar",
      color: "text-purple-600"
    },
  ]

  const orderStatus = [
    { title: "Pending", value: data.pending_orders, icon: Clock, color: "text-yellow-500" },
    { title: "Accepted", value: data.accepted_orders, icon: CheckCircle2, color: "text-emerald-500" },
    { title: "Declined", value: data.decline_orders, icon: XCircle, color: "text-destructive" },
    { title: "Categories", value: data.total_category, icon: Layers, color: "text-slate-500" },
    { title: "Coupons", value: data.total_coupons, icon: Ticket, color: "text-pink-500" },
  ]

  return (
    <div className="flex flex-col gap-6">
      <div>
        <h2 className="text-3xl font-bold tracking-tight">Dashboard</h2>
        <p className="text-muted-foreground">Ringkasan performa sistem Voca Store saat ini.</p>
      </div>

      {/* Main Stats Grid */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        {stats.map((stat) => (
          <Card key={stat.title}>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">{stat.title}</CardTitle>
              <stat.icon className={`h-4 w-4 ${stat.color}`} />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{stat.value}</div>
              <p className="text-xs text-muted-foreground">{stat.description}</p>
            </CardContent>
          </Card>
        ))}
      </div>

      {/* Secondary Stats/Status */}
      <div className="grid gap-4 grid-cols-2 md:grid-cols-3 lg:grid-cols-5">
        {orderStatus.map((status) => (
          <Card key={status.title} className="bg-muted/50">
            <CardContent className="p-4 flex items-center gap-4">
              <div className={`p-2 rounded-full bg-background ${status.color}`}>
                <status.icon className="h-4 w-4" />
              </div>
              <div>
                <p className="text-sm font-medium text-muted-foreground">{status.title}</p>
                <p className="text-xl font-bold">{status.value}</p>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>
      
      {/* Kamu bisa menambahkan grafik di bawah sini nantinya */}
    </div>
  )
}

export default OverviewPage