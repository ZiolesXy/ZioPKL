import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"

import { Badge } from "@/components/ui/badge"

import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import LogoutButton from "@/components/Log-out"

export default function UserPage() {
  // contoh data (nanti bisa diganti hasil fetch API)
  const user = {
    name: "Narata",
    role: "User",
  }

  const items = [
    { id: 1, name: "Laptop", qty: 1 },
    { id: 2, name: "Mouse", qty: 2 },
    { id: 3, name: "Keyboard", qty: 1 },
  ]

  return (
    <div className="p-6 space-y-6">
      {/* Title */}
      <h1 className="text-2xl font-bold">Dashboard User</h1>

      {/* User Info */}
      <Card>
        <CardHeader>
          <CardTitle>Informasi User</CardTitle>
        </CardHeader>
        <CardContent className="space-y-2">
          <p>
            <span className="font-medium">Nama:</span> {user.name}
          </p>
          <p>
            <span className="font-medium">Role:</span>{" "}
            <Badge variant="secondary">{user.role}</Badge>
          </p>
        </CardContent>
      </Card>

      {/* User Items */}
      <Card>
        <CardHeader>
          <CardTitle>Barang Saya</CardTitle>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Nama Barang</TableHead>
                <TableHead>Jumlah</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {items.map((item) => (
                <TableRow key={item.id}>
                  <TableCell>{item.name}</TableCell>
                  <TableCell>{item.qty}</TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </CardContent>
        <LogoutButton />
      </Card>
    </div>
  )
}
