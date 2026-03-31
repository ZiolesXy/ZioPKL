import api from "../../axios"

type DeleteType = "user" | "barang" | "role" | "ownership"

interface DeletePayload {
  type: DeleteType
  userId?: number
  barangId?: number
  userName?: string
  barangName?: string
  roleId?: number
}

export async function deleteResource(payload: DeletePayload) {
  const { type, userId, barangId, roleId } = payload

  switch (type) {
    case "user":
      if (!userId) throw new Error("userId required")
      return api.delete(`/user/${userId}`)

    case "barang":
      if (!barangId) throw new Error("barangId required")
      return api.delete(`/barang/${barangId}`)

    case "role":
      if (!roleId) throw new Error("roleId required")
      return api.delete(`/role/${roleId}`)

    case "ownership":
      if (!userId || !barangId) throw new Error("userId & barangId required")
      return api.delete(`/user/${userId}/barang/${barangId}`)

    default:
      throw new Error("Invalid delete type")
  }
}
