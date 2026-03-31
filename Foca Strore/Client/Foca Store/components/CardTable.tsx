import { Card, CardContent, CardHeader, CardTitle } from './ui/card'
import getUsers from '@/lib/api/resources/users-native'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"

async function CardTable(){
  type Barang = {
  id: number
  name: string
}

// type Role = {
//   id: number
//   name: string
// }

// type User = {
//   id: number
//   name: string
//   role?: Role
//   barangs?: Barang[]
// }

  const users = await getUsers()
  return (
    <Card className="w-full h-full">
      <CardHeader className='border-b'>
        <CardTitle>List of Users</CardTitle>
      </CardHeader>
      <CardContent className="flex-1 overflow-hidden">
        <div className="w-full overflow-x-auto">
          <Table className='min-w-full text-lg table-fixed'>
            <TableHeader>
              <TableRow>
                <TableHead className='w-20 font-bold'>ID</TableHead>
                <TableHead className='w-45 font-bold'>Name</TableHead>
                <TableHead className='w-40 font-bold'>Role</TableHead>
                <TableHead className='font-bold'>Barang</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {users.map((u: any, idx: number) => {
                const barangNames: string[] = (u.barangs ?? []).map((b: Barang) => b.name).filter(Boolean)
                const preview = barangNames.slice(0, 2)
                const rest = barangNames.slice(2)

                return (
                  <TableRow key={u.id ?? idx}>
                    <TableCell className="whitespace-nowrap">{u.id}</TableCell>
                    <TableCell className="whitespace-nowrap">{u.name}</TableCell>
                    <TableCell className="whitespace-nowrap">{u.role?.name}</TableCell>
                    <TableCell className="min-w-0">
                      <div className="flex flex-wrap items-center gap-1 min-w-0">
                        {preview.map((name) => (
                          <span key={name} className="bg-muted text-muted-foreground rounded-md px-2 py-0.5 text-xs">
                            {name}
                          </span>
                        ))}
                        {rest.length > 0 && (
                          <details className="relative">
                            <summary className="cursor-pointer select-none text-xs text-muted-foreground underline underline-offset-2">
                              +{rest.length}
                            </summary>
                            <div className="absolute right-0 z-10 mt-2 w-72 rounded-md border bg-background p-2 shadow-md">
                              <div className="flex flex-wrap gap-1">
                                {rest.map((name) => (
                                  <span key={name} className="bg-muted text-muted-foreground rounded-md px-2 py-0.5 text-xs">
                                    {name}
                                  </span>
                                ))}
                              </div>
                            </div>
                          </details>
                        )}
                      </div>
                    </TableCell>
                  </TableRow>
                )
              })}
            </TableBody>
          </Table>  
        </div>
      </CardContent>
    </Card>
      
  )
}

export default CardTable
