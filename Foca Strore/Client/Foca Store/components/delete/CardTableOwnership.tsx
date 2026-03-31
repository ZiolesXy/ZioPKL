import getOwnership from "@/lib/api/resources/ownership-native"
import getUsers from "@/lib/api/resources/users-native"
import getBarangs from "@/lib/api/resources/barangs-native"
import { TableCell } from "@/components/ui/table"
import DeleteButton from "../DeleteButton"
import TableCard from "../TableCard"
import { buildIdNameLookups, getOwnershipDisplay, resolveOwnershipIds } from "@/lib/ownership-lookup"

async function CardTableOwnership() {
  const [ownership, users, barangs] = await Promise.all([
    getOwnership(),
    getUsers(),
    getBarangs(),
  ])

  const lookups = buildIdNameLookups({
    users: users as unknown[],
    barangs: barangs as unknown[],
  })

  return (
    <TableCard
      title="Ownership"
      columns={[
        { label: "Barang", className: "w-56" },
        { label: "User", className: "w-56" },
        { label: "Action", className: "w-28" },
      ]}
      items={ownership as unknown[]}
      getRowKey={(o, idx) => {
        const rec = (o ?? {}) as Record<string, unknown>
        return `${String(rec.user_id ?? rec.user_name)}-${String(rec.barang_id ?? rec.barang_name)}-${idx}`
      }}
      renderRow={(o, idx) => {
        const display = getOwnershipDisplay({
          row: o,
          lookups: {
            userNameById: lookups.userNameById,
            barangNameById: lookups.barangNameById,
          },
        })

        const resolved = resolveOwnershipIds({
          row: o,
          lookups: {
            userIdByName: lookups.userIdByName,
            barangIdByName: lookups.barangIdByName,
          },
        })

        return (
          <>
            <TableCell className="truncate">{display.barangLabel}</TableCell>
            <TableCell className="truncate">{display.userLabel}</TableCell>
            <TableCell>
              {typeof resolved.userId === "number" && typeof resolved.barangId === "number" ? (
                <DeleteButton
                  type="ownership"
                  userId={resolved.userId}
                  barangId={resolved.barangId}
                  userName={display.userLabel}
                  barangName={display.barangLabel}
                  label={`Hapus relasi ${resolved.userId}-${resolved.barangId}-${idx}`}
                />
              ) : null}
            </TableCell>
          </>
        )
      }}
    />
  )
}

export default CardTableOwnership
