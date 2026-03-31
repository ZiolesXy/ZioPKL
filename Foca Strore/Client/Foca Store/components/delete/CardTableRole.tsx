import { TableCell } from "@/components/ui/table"
import getRoles from '@/lib/api/resources/roles-native'
import DeleteButton from '../DeleteButton'
import TableCard from "../TableCard"

async function CardTableRole() {
  const roles = await getRoles()
  return (
    <TableCard
      title="Roles"
      columns={[
        { label: "ID", className: "w-24" },
        { label: "Name", className: "w-auto" },
        { label: "Action" },
      ]}
      items={roles as Array<{ id: number | string; name?: string }>}
      getRowKey={(role) => role.id}
      renderRow={(role) => (
        <>
          <TableCell className="font-medium">{role.id}</TableCell>
          <TableCell className="truncate">{role.name}</TableCell>
          <TableCell>
            <DeleteButton type="role" roleId={Number(role.id)} />
          </TableCell>
        </>
      )}
    />
  )
}

export default CardTableRole
