import { TableCell } from "@/components/ui/table"
import getRoles from '@/lib/api/resources/roles-native'
import { PostRoleButton } from '../PostRoleButton'
import TableCard from "../TableCard"

async function PostCardTableRole() {
  const roles = await getRoles()
  return (
    <TableCard
      title="List of Roles"
      headerRight={<PostRoleButton />}
      columns={[
        { label: "ID", className: "w-24" },
        { label: "Name", className: "w-auto" },
      ]}
      items={roles as Array<{ id: number | string; name?: string }>}
      getRowKey={(role) => role.id}
      renderRow={(role) => (
        <>
          <TableCell className="font-medium">{role.id}</TableCell>
          <TableCell className="truncate">{role.name}</TableCell>
        </>
      )}
    />
  )
}

export default PostCardTableRole
