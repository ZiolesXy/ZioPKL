export type IdNameLookups = {
  userNameById: Map<number, string>
  barangNameById: Map<number, string>
  userIdByName: Map<string, number>
  barangIdByName: Map<string, number>
}

export function toNumberId(value: unknown): number | null {
  if (typeof value === "number" && Number.isFinite(value)) return value
  if (typeof value === "string") {
    const n = Number(value)
    if (Number.isFinite(n)) return n
  }
  return null
}

export function pickName(obj: unknown): string | null {
  if (obj === null || typeof obj !== "object") return null
  const rec = obj as Record<string, unknown>
  const candidate = rec.name ?? rec.nama ?? rec.title
  return typeof candidate === "string" && candidate.trim() ? candidate : null
}

export function buildIdNameLookups(params: {
  users: unknown[]
  barangs: unknown[]
}): IdNameLookups {
  const userNameById = new Map<number, string>()
  const barangNameById = new Map<number, string>()
  const userIdByName = new Map<string, number>()
  const barangIdByName = new Map<string, number>()

  for (const u of params.users) {
    const rec = u as Record<string, unknown>
    const id = toNumberId(rec.id ?? rec.user_id)
    const name = pickName(u)
    if (id !== null && name) {
      userNameById.set(id, name)
      userIdByName.set(name, id)
    }
  }

  for (const b of params.barangs) {
    const rec = b as Record<string, unknown>
    const id = toNumberId(rec.id ?? rec.barang_id)
    const name = pickName(b)
    if (id !== null && name) {
      barangNameById.set(id, name)
      barangIdByName.set(name, id)
    }
  }

  return { userNameById, barangNameById, userIdByName, barangIdByName }
}

export function getOwnershipDisplay(params: {
  row: unknown
  lookups: Pick<IdNameLookups, "userNameById" | "barangNameById">
}): { userLabel: string; barangLabel: string } {
  const rec = (params.row ?? {}) as Record<string, unknown>

  const barangName = typeof rec.barang_name === "string" ? rec.barang_name : null
  const userName = typeof rec.user_name === "string" ? rec.user_name : null

  const barangId = typeof rec.barang_id === "number" ? rec.barang_id : null
  const userId = typeof rec.user_id === "number" ? rec.user_id : null

  const barangLabel =
    barangName ?? (typeof barangId === "number" ? params.lookups.barangNameById.get(barangId) ?? "-" : "-")

  const userLabel =
    userName ?? (typeof userId === "number" ? params.lookups.userNameById.get(userId) ?? "-" : "-")

  return { userLabel, barangLabel }
}

export function resolveOwnershipIds(params: {
  row: unknown
  lookups: Pick<IdNameLookups, "userIdByName" | "barangIdByName">
}): { userId?: number; barangId?: number } {
  const rec = (params.row ?? {}) as Record<string, unknown>

  const userId = typeof rec.user_id === "number" ? rec.user_id : undefined
  const barangId = typeof rec.barang_id === "number" ? rec.barang_id : undefined

  const resolvedUserId =
    typeof userId === "number"
      ? userId
      : typeof rec.user_name === "string"
        ? params.lookups.userIdByName.get(rec.user_name)
        : undefined

  const resolvedBarangId =
    typeof barangId === "number"
      ? barangId
      : typeof rec.barang_name === "string"
        ? params.lookups.barangIdByName.get(rec.barang_name)
        : undefined

  return { userId: resolvedUserId, barangId: resolvedBarangId }
}
