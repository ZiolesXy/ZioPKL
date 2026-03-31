import { TableCell } from "@/components/ui/table"
import getBarangs from '@/lib/api/resources/barangs-native'
import DeleteButton from '../DeleteButton'
import TableCard from "../TableCard"

async function CardTableBarang() {
  const barangs = await getBarangs()

  return (
    <TableCard
      title="Barang"
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
            <DeleteButton type="barang" barangId={Number(barang.id)} />
          </TableCell>
        </>
      )}
    />
  )
}

export default CardTableBarang
