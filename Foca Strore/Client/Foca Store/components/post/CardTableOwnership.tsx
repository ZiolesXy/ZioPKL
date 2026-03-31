import getOwnership from "@/lib/api/resources/ownership-native"
import getUsers from "@/lib/api/resources/users-native"
import getBarangs from "@/lib/api/resources/barangs-native"
import { TableCell } from "@/components/ui/table"
import { PostOwnershipButton } from "../PostOwnershipButton"
import TableCard from "../TableCard"
import { buildIdNameLookups, getOwnershipDisplay } from "@/lib/ownership-lookup"

async function PostCardTableOwnership() {
  const ownership = await getOwnership()

  const needsJoin = (ownership as unknown[]).some((o) => {
    const rec = (o ?? {}) as Record<string, unknown>
    return !rec.barang_name || !rec.user_name
  })

  const lookups = needsJoin
    ? buildIdNameLookups({
        users: (await getUsers()) as unknown[],
        barangs: (await getBarangs()) as unknown[],
      })
    : buildIdNameLookups({ users: [], barangs: [] })

  return (
    <TableCard
      title="List of Ownership"
      headerRight={<PostOwnershipButton />}
      columns={[
        { label: "Barang", className: "w-56" },
        { label: "User", className: "w-56" },
      ]}
      items={ownership as unknown[]}
      getRowKey={(o, idx) => {
        const rec = (o ?? {}) as Record<string, unknown>
        return `${String(rec.user_id ?? rec.user_name)}-${String(rec.barang_id ?? rec.barang_name)}-${idx}`
      }}
      renderRow={(o) => {
        const display = getOwnershipDisplay({
          row: o,
          lookups: {
            userNameById: lookups.userNameById,
            barangNameById: lookups.barangNameById,
          },
        })

        return (
          <>
            <TableCell className="truncate">{display.barangLabel}</TableCell>
            <TableCell className="truncate">{display.userLabel}</TableCell>
          </>
        )
      }}
    />
  )
}

export default PostCardTableOwnership
