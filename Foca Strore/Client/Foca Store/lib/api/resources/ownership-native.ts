import { serverApi } from "../core/server"

export interface Ownership {
  barang_id?: number
  user_id?: number
  barang_name?: string
  user_name?: string
}

async function getOwnership(): Promise<Ownership[]> {
  return serverApi.getEntries<Ownership>("/user/barang", {
    // keys: ["ownership"],
    label: "ownership",
  })
}

export default getOwnership
