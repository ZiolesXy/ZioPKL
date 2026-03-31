import { TableCell } from "@/components/ui/table"
import getBarangs from '@/lib/api/resources/barangs-native'
import { PostBarangButton } from '../PostBarangButton'
import TableCard from "../TableCard"

async function PostCardTableBarang() {
  const barangs = await getBarangs()

  return (
    <TableCard
      title="List of Barang"
      headerRight={<PostBarangButton />}
      columns={[
        { label: "ID", className: "w-24" },
        { label: "Name", className: "w-auto" },
      ]}
      items={barangs as Array<{ id: number | string; name?: string }>}
      getRowKey={(barang) => barang.id}
      renderRow={(barang) => (
        <>
          <TableCell className="font-medium">{barang.id}</TableCell>
          <TableCell className="truncate">{barang.name}</TableCell>
        </>
      )}
    />
  )
}

export default PostCardTableBarang
