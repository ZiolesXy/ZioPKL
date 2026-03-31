import { TableCell } from "@/components/ui/table"
import getBarangs from '@/lib/api/resources/barangs-native'
import PutButton from '../PutButton'
import TableCard from "../TableCard"

async function PutCardTableBarang() {
  const barangs = await getBarangs()

  return (
    <TableCard
      title="List of Barang"
      columns={[
        { label: "ID", className: "w-24" },
        { label: "Name", className: "w-auto" },
        { label: "Action" },
      ]}
      items={barangs as Array<{ id: number | string; name?: string }>}
      getRowKey={(barang) => barang.id}
      renderRow={(barang) => (
        <>
          <TableCell className="font-medium">{barang.id}</TableCell>
          <TableCell className="truncate">{barang.name}</TableCell>
          <TableCell>
            <PutButton
              id={barang.id}
              endpoint="/barang"
              defaultValues={{
                name: barang.name,
              }}
            />
          </TableCell>
        </>
      )}
    />
  )
}

export default PutCardTableBarang
