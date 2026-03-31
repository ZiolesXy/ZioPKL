import { TableCell } from "@/components/ui/table"
import getRoles from '@/lib/api/resources/roles-native'
import PutButton from '../PutButton'
import TableCard from "../TableCard"

async function PutCardTableRole() {
  const roles = await getRoles()
  return (
    <TableCard
      title="List of Role"
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
            <PutButton
              id={role.id}
              endpoint="/role"
              defaultValues={{
                name: role.name,
              }}
            />
          </TableCell>
        </>
      )}
    />
  )
}

export default PutCardTableRole
